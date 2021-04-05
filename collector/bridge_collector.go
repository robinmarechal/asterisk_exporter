package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// bridgeCollector collector for all 'bridge show ...' commands
type bridgeCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	// BridgeTechnologies
	bridgeTechnologiesInfo *prometheus.Desc
	bridgesInfo            *prometheus.Desc

	collectorError *prometheus.Desc
}

type bridgeMetrics struct {
	BridgesInfo            *cmd.BridgesInfo
	BridgeTechnologiesInfo *cmd.BridgeTechnologiesInfo
}

func NewBridgeCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &bridgeCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		bridgeTechnologiesInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "bridges", "technologies_info"),
			"Bridge technologies info",
			[]string{"name", "type", "priority", "suspended"}, nil,
		),
		bridgesInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "bridges", "info"),
			"Bridges info",
			nil, nil,
		),
	}
}

func (c *bridgeCollector) Name() string {
	return "bridges"
}

func (c *bridgeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bridgeTechnologiesInfo
	ch <- c.bridgesInfo
}

func (c *bridgeCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting bridge metrics")
	metrics, err := collectBridgeMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "bridge metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectBridgeMetrics(c *cmd.CmdRunner) (*bridgeMetrics, error) {
	metrics := &bridgeMetrics{
		BridgeTechnologiesInfo: c.BridgeTechnologiesInfo(),
		BridgesInfo:            c.BridgesInfo(),
	}

	return metrics, nil
}

func (c *bridgeCollector) updateMetrics(values *bridgeMetrics, ch chan<- prometheus.Metric) {
	for _, btech := range values.BridgeTechnologiesInfo.BridgeTechnologies {
		ch <- prometheus.MustNewConstMetric(c.bridgeTechnologiesInfo, prometheus.GaugeValue, 1,
			btech.Name, btech.Type, btech.Priority, btech.Suspended)
	}

	ch <- prometheus.MustNewConstMetric(c.bridgesInfo, prometheus.GaugeValue, float64(values.BridgesInfo.Count))

	level.Debug(c.logger).Log("msg", "bridge metrics built")
}
