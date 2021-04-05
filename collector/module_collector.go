package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asterisk_exporter/cmd"
)

// moduleCollector collector for all 'module show ...' commands
type moduleCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	modulesCount   *prometheus.Desc
	collectorError *prometheus.Desc
}

type moduleMetrics struct {
	ModulesInfo *cmd.ModulesInfo
}

func NewModuleCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &moduleCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		modulesCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "modules", "count"),
			"Number of installed modules",
			nil, nil,
		),
	}
}

func (c *moduleCollector) Name() string {
	return "modules"
}

func (c *moduleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.modulesCount
}

func (c *moduleCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting module metrics")
	metrics, err := collectModuleMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "module metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectModuleMetrics(c *cmd.CmdRunner) (*moduleMetrics, error) {
	metrics := &moduleMetrics{
		ModulesInfo: c.ModulesInfo(),
	}

	return metrics, nil
}

func (c *moduleCollector) updateMetrics(values *moduleMetrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.modulesCount, prometheus.GaugeValue, float64(values.ModulesInfo.ModuleCount))
	level.Debug(c.logger).Log("msg", "module metrics built")
}
