package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asterisk_exporter/cmd"
)

// calendarCollector collector for all 'calendar show ...' commands
type calendarCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	calendarsCount *prometheus.Desc
	collectorError *prometheus.Desc
}

type calendarMetrics struct {
	CalendarsInfo *cmd.CalendarsInfo
}

func NewCalendarCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &calendarCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		calendarsCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "calendars", "count"),
			"Number of calendars",
			nil, nil,
		),
	}
}

func (c *calendarCollector) Name() string {
	return "calendars"
}

func (c *calendarCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.calendarsCount
}

func (c *calendarCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting calendar metrics")
	metrics, err := collectCalendarMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "calendar metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectCalendarMetrics(c *cmd.CmdRunner) (*calendarMetrics, error) {
	metrics := &calendarMetrics{
		CalendarsInfo: c.CalendarsInfo(),
	}

	return metrics, nil
}

func (c *calendarCollector) updateMetrics(values *calendarMetrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.calendarsCount, prometheus.GaugeValue, float64(values.CalendarsInfo.Count))
	level.Debug(c.logger).Log("msg", "calendar metrics built")
}
