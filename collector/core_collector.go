package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asterisk_exporter/cmd"
	"github.com/robinmarechal/asterisk_exporter/util"
)

// coreCollector collector for all 'core show ...' commands
type coreCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	totalActiveChannels      *prometheus.Desc
	totalActiveCalls         *prometheus.Desc
	totalCallsProcessed      *prometheus.Desc
	systemUptimeSeconds      *prometheus.Desc
	lastReloadSeconds        *prometheus.Desc
	imagesRegistered         *prometheus.Desc
	systemTotalMemoryBytes   *prometheus.Desc
	systemFreeMemoryBytes    *prometheus.Desc
	systemBufferMemoryBytes  *prometheus.Desc
	systemTotalSwapBytes     *prometheus.Desc
	systemFreeSwapBytes      *prometheus.Desc
	systemProcesses          *prometheus.Desc
	threadCount              *prometheus.Desc
	channelActive            *prometheus.Desc
	channelIndication        *prometheus.Desc
	channelTransfer          *prometheus.Desc
	tasksProcessors          *prometheus.Desc
	tasksProcessedTasksTotal *prometheus.Desc
	tasksProcessesInQueue    *prometheus.Desc
	version                  *prometheus.Desc

	collectorError *prometheus.Desc
}

type coreMetrics struct {
	UptimeInfo         *cmd.UptimeInfo
	ChannelsInfo       *cmd.ChannelsInfo
	ThreadsInfo        *cmd.ThreadsInfo
	ChannelTypesInfo   *cmd.ChannelTypesInfo
	ImagesInfo         *cmd.ImagesInfo
	SystemInfo         *cmd.SystemInfo
	TaskProcessorsInfo *cmd.TaskProcessorsInfo
	VersionInfo        *cmd.VersionInfo
}

func NewCoreCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &coreCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		totalActiveChannels: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "active_channels"),
			"Number of currently active channels",
			nil, nil,
		),
		totalActiveCalls: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "active_calls"),
			"Number of currently active calls",
			nil, nil,
		),
		totalCallsProcessed: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "processed_calls"),
			"Number of processed calls",
			nil, nil,
		),
		systemUptimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_uptime_seconds"),
			"Number of seconds since system startup",
			nil, nil,
		),
		lastReloadSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "last_reload_seconds"),
			"Number of seconds since last reload",
			nil, nil,
		),
		imagesRegistered: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "images_registered"),
			"Number of registered images",
			nil, nil,
		),
		systemTotalMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_memory_total_bytes"),
			"System total memory in bytes",
			nil, nil,
		),
		systemFreeMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_memory_free_bytes"),
			"System free memory in bytes",
			nil, nil,
		),
		systemBufferMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_memory_buffer_bytes"),
			"System buffer memory in bytes",
			nil, nil,
		),
		systemTotalSwapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_swap_total_bytes"),
			"System total swap in bytes",
			nil, nil,
		),
		systemFreeSwapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_swap_free_bytes"),
			"System free swap in bytes",
			nil, nil,
		),
		systemProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "system_processes"),
			"Number of system processes",
			nil, nil,
		),
		threadCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "thread_count"),
			"Number of threads",
			nil, nil,
		),
		channelActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "channel_active"),
			"Active flag of channels",
			[]string{"type"}, nil,
		),
		channelIndication: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "channel_indication"),
			"Indication flag of channels",
			[]string{"type"}, nil,
		),
		channelTransfer: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "channel_transfer"),
			"Transfer flag of channels",
			[]string{"type"}, nil,
		),
		tasksProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "tasks_processors"),
			"Number of task processors",
			nil, nil,
		),
		tasksProcessedTasksTotal: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "tasks_processed_total"),
			"Number of processed tasks",
			nil, nil,
		),
		tasksProcessesInQueue: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "tasks_processes_in_queue"),
			"Task processes in queue",
			nil, nil,
		),
		version: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "core", "version"),
			"Version info",
			[]string{"version"}, nil,
		),
	}
}

func (c *coreCollector) Name() string {
	return "core"
}

func (c *coreCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalActiveChannels
	ch <- c.totalActiveCalls
	ch <- c.totalCallsProcessed
	ch <- c.systemUptimeSeconds
	ch <- c.lastReloadSeconds
	ch <- c.imagesRegistered
	ch <- c.systemTotalMemoryBytes
	ch <- c.systemFreeMemoryBytes
	ch <- c.systemBufferMemoryBytes
	ch <- c.systemTotalSwapBytes
	ch <- c.systemFreeSwapBytes
	ch <- c.systemProcesses
	ch <- c.tasksProcessors
	ch <- c.tasksProcessedTasksTotal
	ch <- c.tasksProcessesInQueue
	ch <- c.version
}

func (c *coreCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting core metrics")
	metrics, err := collectCoreMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "core metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectCoreMetrics(c *cmd.CmdRunner) (*coreMetrics, error) {
	metrics := &coreMetrics{
		UptimeInfo:         c.UptimeInfos(),
		ChannelsInfo:       c.ChannelsInfo(),
		ThreadsInfo:        c.ThreadsInfo(),
		ChannelTypesInfo:   c.ChannelTypesInfo(),
		ImagesInfo:         c.ImagesInfo(),
		SystemInfo:         c.SystemInfo(),
		TaskProcessorsInfo: c.TaskProcessorsInfo(),
		VersionInfo:        c.VersionInfo(),
	}

	return metrics, nil
}

func (c *coreCollector) updateMetrics(values *coreMetrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.totalActiveChannels, prometheus.GaugeValue, float64(values.ChannelsInfo.ActiveChannels))
	ch <- prometheus.MustNewConstMetric(c.totalActiveCalls, prometheus.GaugeValue, float64(values.ChannelsInfo.ActiveCalls))
	ch <- prometheus.MustNewConstMetric(c.totalCallsProcessed, prometheus.GaugeValue, float64(values.ChannelsInfo.ProcessedCalls))
	ch <- prometheus.MustNewConstMetric(c.systemUptimeSeconds, prometheus.GaugeValue, float64(values.UptimeInfo.SystemUptimeSeconds))
	ch <- prometheus.MustNewConstMetric(c.lastReloadSeconds, prometheus.GaugeValue, float64(values.UptimeInfo.LastReloadSeconds))
	ch <- prometheus.MustNewConstMetric(c.imagesRegistered, prometheus.GaugeValue, float64(values.ImagesInfo.Registered))
	ch <- prometheus.MustNewConstMetric(c.systemTotalMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.TotalMemory))
	ch <- prometheus.MustNewConstMetric(c.systemFreeMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.FreeMemory))
	ch <- prometheus.MustNewConstMetric(c.systemBufferMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.BufferMemory))
	ch <- prometheus.MustNewConstMetric(c.systemTotalSwapBytes, prometheus.GaugeValue, float64(values.SystemInfo.TotalSwap))
	ch <- prometheus.MustNewConstMetric(c.systemFreeSwapBytes, prometheus.GaugeValue, float64(values.SystemInfo.FreeSwap))
	ch <- prometheus.MustNewConstMetric(c.systemProcesses, prometheus.GaugeValue, float64(values.SystemInfo.ProcessCount))
	ch <- prometheus.MustNewConstMetric(c.threadCount, prometheus.GaugeValue, float64(values.ThreadsInfo.ThreadCount))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessors, prometheus.GaugeValue, float64(values.TaskProcessorsInfo.ProcessorCounter))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessedTasksTotal, prometheus.CounterValue, float64(values.TaskProcessorsInfo.ProcessedTasksTotal))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessesInQueue, prometheus.GaugeValue, float64(values.TaskProcessorsInfo.InQueue))
	ch <- prometheus.MustNewConstMetric(c.version, prometheus.GaugeValue, 1, values.VersionInfo.Version)

	for _, typeInfo := range values.ChannelTypesInfo.ChannelTypes {
		ch <- prometheus.MustNewConstMetric(c.channelActive, prometheus.GaugeValue, util.BoolToFloat(typeInfo.DeviceState), typeInfo.Type)
		ch <- prometheus.MustNewConstMetric(c.channelTransfer, prometheus.GaugeValue, util.BoolToFloat(typeInfo.Transfer), typeInfo.Type)
		ch <- prometheus.MustNewConstMetric(c.channelIndication, prometheus.GaugeValue, util.BoolToFloat(typeInfo.Indications), typeInfo.Type)
	}

	level.Debug(c.logger).Log("msg", "core metrics built")
}
