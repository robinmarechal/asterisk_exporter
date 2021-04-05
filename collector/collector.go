package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// Custom collector interface
type Collector interface {
	prometheus.Collector

	Name() string
}

type CollectorFactory func(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, errorMetric *prometheus.Desc) Collector
