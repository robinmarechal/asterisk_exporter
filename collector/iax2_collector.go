package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// iax2Collector collector for all 'iax2 show ...' commands
type iax2Collector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	iaxChannelActive *prometheus.Desc
	collectorError   *prometheus.Desc
}

type iax2Metrics struct {
	IaxChannelsInfo *cmd.IaxChannelsInfo
}

func NewdIax2Collector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &iax2Collector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		iaxChannelActive: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "iax2", "channels_active"),
			"Number of IAX Active channels",
			nil, nil,
		),
	}
}

func (c *iax2Collector) Name() string {
	return "iax2"
}

func (c *iax2Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.iaxChannelActive
}

func (c *iax2Collector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting iax2 metrics")
	metrics, err := collectdIax2Metrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "iax2 metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectdIax2Metrics(c *cmd.CmdRunner) (*iax2Metrics, error) {
	metrics := &iax2Metrics{
		IaxChannelsInfo: c.IaxChannelsInfo(),
	}

	return metrics, nil
}

func (c *iax2Collector) updateMetrics(values *iax2Metrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.iaxChannelActive, prometheus.GaugeValue, float64(values.IaxChannelsInfo.ActiveCount))

	level.Debug(c.logger).Log("msg", "iax2 metrics built")
}
