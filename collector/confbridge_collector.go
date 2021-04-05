package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// confbridgeCollector collector for all 'confbridge show ...' commands
type confbridgeCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	confBridgeInfo *prometheus.Desc
	collectorError *prometheus.Desc
}

type confbridgeMetrics struct {
	ConfBridgeInfo *cmd.ConfBridgeInfo
}

func NewConfbridgeCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &confbridgeCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		confBridgeInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "confbridges", "info"),
			"ConfBridge information",
			[]string{"type", "name"}, nil,
		),
	}
}

func (c *confbridgeCollector) Name() string {
	return "confbridges"
}

func (c *confbridgeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.confBridgeInfo
}

func (c *confbridgeCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting confbridge metrics")
	metrics, err := collectConfbridgeMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "confbridge metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectConfbridgeMetrics(c *cmd.CmdRunner) (*confbridgeMetrics, error) {
	metrics := &confbridgeMetrics{
		ConfBridgeInfo: c.ConfBridgeInfo(),
	}

	return metrics, nil
}

func (c *confbridgeCollector) updateMetrics(values *confbridgeMetrics, ch chan<- prometheus.Metric) {

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

	level.Debug(c.logger).Log("msg", "confbridge metrics built")
}
