package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// exporterCollector exporter's internal metrics
type exporterCollector struct {
	collectorError *prometheus.Desc
}

func NewExporterCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger) prometheus.Collector {
	return &exporterCollector{
		collectorError: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "", "collector_error"),
			"Collector errors. 0 = no error, 1 = error occurred",
			[]string{"collector"}, nil,
		),
	}
}

func (c *exporterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.collectorError
}

func (c *exporterCollector) Collect(ch chan<- prometheus.Metric) {}
