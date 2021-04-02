package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/robinmarechal/asteriskk_exporter/cmd"
	"github.com/robinmarechal/asteriskk_exporter/util"
)

type asteriskCollector struct {
	cmdRunner *cmd.CmdRunner

	logger *log.Logger

	totalActiveChannels     *prometheus.Desc
	totalActiveCalls        *prometheus.Desc
	totalCallsProcessed     *prometheus.Desc
	systemUptimeSeconds     *prometheus.Desc
	lastReloadSeconds       *prometheus.Desc
	totalSipPeers           *prometheus.Desc
	totalMonitoredOnline    *prometheus.Desc
	totalMonitoredOffline   *prometheus.Desc
	totalUnmonitoredOnline  *prometheus.Desc
	totalUnmonitoredOffline *prometheus.Desc
	totalThreadsListed      *prometheus.Desc
	totalSipStatusUnknown   *prometheus.Desc
	totalSipStatusQualified *prometheus.Desc

	// Agents Infos
	agentsDefined *prometheus.Desc
	agentsLogged  *prometheus.Desc
	agentsTalking *prometheus.Desc

	// BridgeTechnologies
	bridgeTechnologiesInfo *prometheus.Desc
	bridgesInfo            *prometheus.Desc

	calendarsCount *prometheus.Desc

	channelsActive      *prometheus.Desc
	channelsIndications *prometheus.Desc
	channelsTransfer    *prometheus.Desc

	confBridgeInfo *prometheus.Desc // labels: type, name

	iaxChannelActive *prometheus.Desc

	imagesRegistered *prometheus.Desc

	modulesCount *prometheus.Desc

	onlineAgentsDefined *prometheus.Desc
	onlineAgentsLogged  *prometheus.Desc
	onlineAgentsTalking *prometheus.Desc

	sipDialogsActive       *prometheus.Desc
	sipSubscriptionsActive *prometheus.Desc
	sipChannelsActive      *prometheus.Desc

	systemTotalMemoryBytes  *prometheus.Desc
	systemFreeMemoryBytes   *prometheus.Desc
	systemBufferMemoryBytes *prometheus.Desc
	systemTotalSwapBytes    *prometheus.Desc
	systemFreeSwapBytes     *prometheus.Desc
	systemProcesses         *prometheus.Desc

	tasksProcessors          *prometheus.Desc
	tasksProcessedTasksTotal *prometheus.Desc
	tasksProcessesInQueue    *prometheus.Desc

	users *prometheus.Desc

	version *prometheus.Desc // label version

	collectErrors *prometheus.Desc
}

type asteriskMetrics struct {
	ChannelsInfo *cmd.ChannelsInfo
	UptimeInfo   *cmd.UptimeInfo
	PeersInfo    *cmd.PeersInfo
	ThreadsInfo  *cmd.ThreadsInfo
	//
	AgentsInfo             *cmd.AgentsInfo
	BridgeTechnologiesInfo *cmd.BridgeTechnologiesInfo
	BridgesInfo            *cmd.BridgesInfo
	CalendarsInfo          *cmd.CalendarsInfo
	ChannelTypesInfo       *cmd.ChannelTypesInfo
	ConfBridgeInfo         *cmd.ConfBridgeInfo
	IaxChannelsInfo        *cmd.IaxChannelsInfo
	ImagesInfo             *cmd.ImagesInfo
	ModulesInfo            *cmd.ModulesInfo
	OnlineAgentsInfo       *cmd.OnlineAgentsInfo
	SipChannelsInfo        *cmd.SipChannelsInfo
	SystemInfo             *cmd.SystemInfo
	TaskProcessorsInfo     *cmd.TaskProcessorsInfo
	UsersInfo              *cmd.UsersInfo
	VersionInfo            *cmd.VersionInfo
}

// NewAsteriskCollector AsteriskCollector constructor
func NewAsteriskCollector(prefix string, asteriskPath *string, logger *log.Logger) prometheus.Collector {
	return &asteriskCollector{
		cmdRunner: cmd.NewCmdRunner(*asteriskPath, logger),
		logger:    logger,
		totalActiveChannels: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "active_channels"),
			"Number of currently active channels",
			nil, nil,
		),
		totalActiveCalls: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "active_calls"),
			"Number of currently active calls",
			nil, nil,
		),
		totalCallsProcessed: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "processed_calls"),
			"Number of processed calls",
			nil, nil,
		),
		systemUptimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_uptime_seconds"),
			"Number of seconds since system startup",
			nil, nil,
		),
		lastReloadSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "last_reload_seconds"),
			"Number of seconds since last reload",
			nil, nil,
		),
		totalSipPeers: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_peers"),
			"Number of SIP peers",
			nil, nil,
		),
		totalMonitoredOnline: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_monitored_online"),
			"Number of currently monitored online SIP",
			nil, nil,
		),
		totalMonitoredOffline: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_monitored_offline"),
			"Number of currently monitored offline SIP",
			nil, nil,
		),
		totalUnmonitoredOnline: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_unmonitored_online"),
			"Number of currently unmonitored online SIP",
			nil, nil,
		),
		totalUnmonitoredOffline: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_unmonitored_offline"),
			"Number of currently unmonitored offline SIP",
			nil, nil,
		),
		totalThreadsListed: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_threads"),
			"Number of threads",
			nil, nil,
		),
		totalSipStatusUnknown: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_unknown"),
			"Current number of unknown SIP",
			nil, nil,
		),
		totalSipStatusQualified: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "current_sip_qualified"),
			"Current number of qualified SIP",
			nil, nil,
		),
		collectErrors: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "collector_errors"),
			"Metrics collection errors. 0 = no error, 1 = some errors (check logs)",
			nil, nil,
		),
		agentsDefined: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_defined"),
			"Number of defined agents",
			nil, nil,
		),
		agentsLogged: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_logged"),
			"Number of logged agents",
			nil, nil,
		),
		agentsTalking: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_talking"),
			"Number of talking agents",
			nil, nil,
		),
		bridgeTechnologiesInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "bridge_technologies_info"),
			"Bridge technologies info",
			[]string{"name", "type", "priority", "suspended"}, nil,
		),
		bridgesInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "bridges_info"),
			"Bridges info",
			nil, nil,
		),
		calendarsCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "calendars"),
			"Number of calendars",
			nil, nil,
		),
		channelsActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "channels_active"),
			"Channels flag 'State'. 1 = 'Yes', 0 = 'No'",
			[]string{"type"}, nil,
		),
		channelsIndications: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "channels_indications"),
			"Channels flag 'Indications'. 1 = 'Yes', 0 = 'No'",
			[]string{"type"}, nil,
		),
		channelsTransfer: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "channels_transfer"),
			"Channels flag 'Transfer'. 1 = 'Yes', 0 = 'No'",
			[]string{"type"}, nil,
		),
		confBridgeInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "conf_bridge_info"),
			"ConfBridge information",
			[]string{"type", "name"}, nil,
		),
		iaxChannelActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "iax_channel_active"),
			"Number of IAX Active channels",
			nil, nil,
		),
		imagesRegistered: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "images_registered"),
			"Number of registered images",
			nil, nil,
		),
		modulesCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "modules"),
			"Number of installed modules",
			nil, nil,
		),
		onlineAgentsDefined: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_online_defined"),
			"Number of defined online agents",
			nil, nil,
		),
		onlineAgentsLogged: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_online_logged"),
			"Number of logged online agents",
			nil, nil,
		),
		onlineAgentsTalking: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "agents_online_talking"),
			"Number of talking online agents",
			nil, nil,
		),
		sipDialogsActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "sip_active_dialogs"),
			"Number of active SIP dialogs",
			nil, nil,
		),
		sipSubscriptionsActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "sip_active_subscriptions"),
			"Number of active SIP subscriptions",
			nil, nil,
		),
		sipChannelsActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "sip_active_channels"),
			"Number of active SIP channels",
			nil, nil,
		),
		systemTotalMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_memory_total_bytes"),
			"System total memory in bytes",
			nil, nil,
		),
		systemFreeMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_memory_free_bytes"),
			"System free memory in bytes",
			nil, nil,
		),
		systemBufferMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_memory_buffer_bytes"),
			"System buffer memory in bytes",
			nil, nil,
		),
		systemTotalSwapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_swap_total_bytes"),
			"System total swap in bytes",
			nil, nil,
		),
		systemFreeSwapBytes: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_swap_free_bytes"),
			"System free swap in bytes",
			nil, nil,
		),
		systemProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "system_processes"),
			"Number of system processes",
			nil, nil,
		),
		tasksProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "tasks_processors"),
			"Number of task processors",
			nil, nil,
		),
		tasksProcessedTasksTotal: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "tasks_processed_total"),
			"Number of processed tasks",
			nil, nil,
		),
		tasksProcessesInQueue: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "tasks_processes_in_queue"),
			"Task processes in queue",
			nil, nil,
		),
		users: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "users"),
			"Number of users",
			nil, nil,
		),
		version: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "version"),
			"Version info",
			[]string{"version"}, nil,
		),
	}
}

// Describe implementation
func (c *asteriskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalActiveChannels
	ch <- c.totalActiveCalls
	ch <- c.totalCallsProcessed
	ch <- c.systemUptimeSeconds
	ch <- c.lastReloadSeconds
	ch <- c.totalSipPeers
	ch <- c.totalMonitoredOnline
	ch <- c.totalMonitoredOffline
	ch <- c.totalUnmonitoredOnline
	ch <- c.totalUnmonitoredOffline
	ch <- c.totalThreadsListed
	ch <- c.totalSipStatusUnknown
	ch <- c.totalSipStatusQualified
	ch <- c.collectErrors

	ch <- c.agentsDefined
	ch <- c.agentsLogged
	ch <- c.agentsTalking
	ch <- c.bridgeTechnologiesInfo
	ch <- c.bridgesInfo
	ch <- c.calendarsCount
	ch <- c.channelsActive
	ch <- c.channelsIndications
	ch <- c.channelsTransfer
	ch <- c.confBridgeInfo
	ch <- c.iaxChannelActive
	ch <- c.imagesRegistered
	ch <- c.modulesCount
	ch <- c.onlineAgentsDefined
	ch <- c.onlineAgentsLogged
	ch <- c.onlineAgentsTalking
	ch <- c.sipDialogsActive
	ch <- c.sipSubscriptionsActive
	ch <- c.sipChannelsActive
	ch <- c.systemTotalMemoryBytes
	ch <- c.systemFreeMemoryBytes
	ch <- c.systemBufferMemoryBytes
	ch <- c.systemTotalSwapBytes
	ch <- c.systemFreeSwapBytes
	ch <- c.systemProcesses
	ch <- c.tasksProcessors
	ch <- c.tasksProcessedTasksTotal
	ch <- c.tasksProcessesInQueue
	ch <- c.users
	ch <- c.version
}

// Collect implementation
func (c *asteriskCollector) Collect(ch chan<- prometheus.Metric) {

	level.Debug(*c.logger).Log("msg", "Running command "+c.cmdRunner.Cmd+"...", "cmd", c.cmdRunner.Cmd)
	metrics, err := collectMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.collectErrors, prometheus.GaugeValue, 1,
		)
		level.Error(*c.logger).Log("err", err)
		return
	}

	level.Info(*c.logger).Log("msg", "Metrics collected")

	ch <- prometheus.MustNewConstMetric(
		c.collectErrors, prometheus.GaugeValue, 0,
	)

	c.updateMetrics(metrics, ch)
}

func collectMetrics(c *cmd.CmdRunner) (*asteriskMetrics, error) {
	metrics := &asteriskMetrics{
		UptimeInfo:   c.UptimeInfos(),
		ChannelsInfo: c.ChannelsInfo(),
		PeersInfo:    c.PeersInfo(),
		ThreadsInfo:  c.ThreadsInfo(),
		//
		AgentsInfo:             c.AgentsInfo(),
		BridgeTechnologiesInfo: c.BridgeTechnologiesInfo(),
		BridgesInfo:            c.BridgesInfo(),
		CalendarsInfo:          c.CalendarsInfo(),
		ChannelTypesInfo:       c.ChannelTypesInfo(),
		ConfBridgeInfo:         c.ConfBridgeInfo(),
		IaxChannelsInfo:        c.IaxChannelsInfo(),
		ImagesInfo:             c.ImagesInfo(),
		ModulesInfo:            c.ModulesInfo(),
		OnlineAgentsInfo:       c.OnlineAgentsInfo(),
		SipChannelsInfo:        c.SipChannelsInfo(),
		SystemInfo:             c.SystemInfo(),
		TaskProcessorsInfo:     c.TaskProcessorsInfo(),
		UsersInfo:              c.UsersInfo(),
		VersionInfo:            c.VersionInfo(),
	}

	return metrics, nil
}

func (c *asteriskCollector) updateMetrics(values *asteriskMetrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.totalActiveChannels, prometheus.GaugeValue, float64(values.ChannelsInfo.ActiveChannels))
	ch <- prometheus.MustNewConstMetric(c.totalActiveCalls, prometheus.GaugeValue, float64(values.ChannelsInfo.ActiveCalls))
	ch <- prometheus.MustNewConstMetric(c.totalCallsProcessed, prometheus.GaugeValue, float64(values.ChannelsInfo.ProcessedCalls))
	ch <- prometheus.MustNewConstMetric(c.systemUptimeSeconds, prometheus.GaugeValue, float64(values.UptimeInfo.SystemUptimeSeconds))
	ch <- prometheus.MustNewConstMetric(c.lastReloadSeconds, prometheus.GaugeValue, float64(values.UptimeInfo.LastReloadSeconds))
	ch <- prometheus.MustNewConstMetric(c.totalSipPeers, prometheus.GaugeValue, float64(values.PeersInfo.SipPeers))
	ch <- prometheus.MustNewConstMetric(c.totalMonitoredOnline, prometheus.GaugeValue, float64(values.PeersInfo.MonitoredOnline))
	ch <- prometheus.MustNewConstMetric(c.totalMonitoredOffline, prometheus.GaugeValue, float64(values.PeersInfo.MonitoredOffline))
	ch <- prometheus.MustNewConstMetric(c.totalUnmonitoredOnline, prometheus.GaugeValue, float64(values.PeersInfo.UnmonitoredOnline))
	ch <- prometheus.MustNewConstMetric(c.totalUnmonitoredOffline, prometheus.GaugeValue, float64(values.PeersInfo.UnmonitoredOffline))
	ch <- prometheus.MustNewConstMetric(c.totalThreadsListed, prometheus.GaugeValue, float64(values.ThreadsInfo.ThreadCount))
	ch <- prometheus.MustNewConstMetric(c.totalSipStatusUnknown, prometheus.GaugeValue, float64(values.PeersInfo.PeersStatusUnknown))
	ch <- prometheus.MustNewConstMetric(c.totalSipStatusQualified, prometheus.GaugeValue, float64(values.PeersInfo.PeersStatusQualified))

	ch <- prometheus.MustNewConstMetric(c.agentsDefined, prometheus.GaugeValue, float64(values.AgentsInfo.DefinedAgents))
	ch <- prometheus.MustNewConstMetric(c.agentsLogged, prometheus.GaugeValue, float64(values.AgentsInfo.LoggedAgents))
	ch <- prometheus.MustNewConstMetric(c.agentsTalking, prometheus.GaugeValue, float64(values.AgentsInfo.TalkingAgents))

	for _, btech := range values.BridgeTechnologiesInfo.BridgeTechnologies {
		ch <- prometheus.MustNewConstMetric(c.bridgeTechnologiesInfo, prometheus.GaugeValue, 1,
			btech.Name, btech.Type, btech.Priority, btech.Suspended)
	}

	ch <- prometheus.MustNewConstMetric(c.bridgesInfo, prometheus.GaugeValue, float64(values.BridgesInfo.Count))
	ch <- prometheus.MustNewConstMetric(c.calendarsCount, prometheus.GaugeValue, float64(values.CalendarsInfo.Count))

	for _, chanType := range values.ChannelTypesInfo.ChannelTypes {
		ch <- prometheus.MustNewConstMetric(c.channelsActive, prometheus.GaugeValue, util.BoolToFloat(chanType.DeviceState), chanType.Type)
		ch <- prometheus.MustNewConstMetric(c.channelsIndications, prometheus.GaugeValue, util.BoolToFloat(chanType.Indications), chanType.Type)
		ch <- prometheus.MustNewConstMetric(c.channelsTransfer, prometheus.GaugeValue, util.BoolToFloat(chanType.Transfer), chanType.Type)
	}

	for _, v := range values.ConfBridgeInfo.Menus {
		ch <- prometheus.MustNewConstMetric(c.confBridgeInfo, prometheus.GaugeValue, 1,
			"menu", v)
	}
	for _, v := range values.ConfBridgeInfo.Profiles {
		ch <- prometheus.MustNewConstMetric(c.confBridgeInfo, prometheus.GaugeValue, 1,
			"profile", v)
	}
	for _, v := range values.ConfBridgeInfo.Users {
		ch <- prometheus.MustNewConstMetric(c.confBridgeInfo, prometheus.GaugeValue, 1,
			"user", v)
	}

	ch <- prometheus.MustNewConstMetric(c.iaxChannelActive, prometheus.GaugeValue, float64(values.IaxChannelsInfo.ActiveCount))
	ch <- prometheus.MustNewConstMetric(c.imagesRegistered, prometheus.GaugeValue, float64(values.ImagesInfo.Registered))
	ch <- prometheus.MustNewConstMetric(c.modulesCount, prometheus.GaugeValue, float64(values.ModulesInfo.ModuleCount))
	ch <- prometheus.MustNewConstMetric(c.onlineAgentsDefined, prometheus.GaugeValue, float64(values.OnlineAgentsInfo.OnlineDefinedAgents))
	ch <- prometheus.MustNewConstMetric(c.onlineAgentsLogged, prometheus.GaugeValue, float64(values.OnlineAgentsInfo.OnlineLoggedAgents))
	ch <- prometheus.MustNewConstMetric(c.onlineAgentsTalking, prometheus.GaugeValue, float64(values.OnlineAgentsInfo.OnlineTalkingAgents))
	ch <- prometheus.MustNewConstMetric(c.sipDialogsActive, prometheus.GaugeValue, float64(values.SipChannelsInfo.ActiveSipDialogs))
	ch <- prometheus.MustNewConstMetric(c.sipSubscriptionsActive, prometheus.GaugeValue, float64(values.SipChannelsInfo.ActiveSipSubscriptions))
	ch <- prometheus.MustNewConstMetric(c.sipChannelsActive, prometheus.GaugeValue, float64(values.SipChannelsInfo.ActiveSipChannels))
	ch <- prometheus.MustNewConstMetric(c.systemTotalMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.TotalMemory))
	ch <- prometheus.MustNewConstMetric(c.systemFreeMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.FreeMemory))
	ch <- prometheus.MustNewConstMetric(c.systemBufferMemoryBytes, prometheus.GaugeValue, float64(values.SystemInfo.BufferMemory))
	ch <- prometheus.MustNewConstMetric(c.systemTotalSwapBytes, prometheus.GaugeValue, float64(values.SystemInfo.TotalSwap))
	ch <- prometheus.MustNewConstMetric(c.systemFreeSwapBytes, prometheus.GaugeValue, float64(values.SystemInfo.FreeSwap))
	ch <- prometheus.MustNewConstMetric(c.systemProcesses, prometheus.GaugeValue, float64(values.SystemInfo.ProcessCount))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessors, prometheus.GaugeValue, float64(values.TaskProcessorsInfo.ProcessorCounter))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessedTasksTotal, prometheus.CounterValue, float64(values.TaskProcessorsInfo.ProcessedTasksTotal))
	ch <- prometheus.MustNewConstMetric(c.tasksProcessesInQueue, prometheus.GaugeValue, float64(values.TaskProcessorsInfo.InQueue))
	ch <- prometheus.MustNewConstMetric(c.users, prometheus.GaugeValue, float64(values.UsersInfo.Users))
	ch <- prometheus.MustNewConstMetric(c.version, prometheus.GaugeValue, 1, values.VersionInfo.Version)

	level.Debug(*c.logger).Log("msg", "Metrics built")
}
