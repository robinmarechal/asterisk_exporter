package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robinmarechal/asteriskk_exporter/cmd"
)

// agentCollector collector for all 'agent show ...' commands
type agentCollector struct {
	cmdRunner *cmd.CmdRunner
	logger    log.Logger

	agentsDefined       *prometheus.Desc
	agentsLogged        *prometheus.Desc
	agentsTalking       *prometheus.Desc
	onlineAgentsLogged  *prometheus.Desc
	onlineAgentsTalking *prometheus.Desc

	collectorError *prometheus.Desc
}

type agentMetrics struct {
	AgentsInfo       *cmd.AgentsInfo
	OnlineAgentsInfo *cmd.OnlineAgentsInfo
}

func NewAgentCollector(prefix string, cmdRunner *cmd.CmdRunner, logger log.Logger, collectorError *prometheus.Desc) Collector {
	return &agentCollector{
		cmdRunner:      cmdRunner,
		logger:         logger,
		collectorError: collectorError,
		agentsDefined: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "agents", "defined"),
			"Number of defined agents",
			nil, nil,
		),
		agentsLogged: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "agents", "logged"),
			"Number of logged agents",
			nil, nil,
		),
		agentsTalking: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "agents", "talking"),
			"Number of talking agents",
			nil, nil,
		),
		onlineAgentsLogged: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "agents", "online_logged"),
			"Number of logged online agents",
			nil, nil,
		),
		onlineAgentsTalking: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "agents", "online_talking"),
			"Number of talking online agents",
			nil, nil,
		),
	}
}

func (c *agentCollector) Name() string {
	return "agents"
}

func (c *agentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.agentsDefined
	ch <- c.agentsLogged
	ch <- c.agentsTalking
	ch <- c.onlineAgentsLogged
	ch <- c.onlineAgentsTalking
}

func (c *agentCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "collecting agent metrics")
	metrics, err := collectAgentMetrics(c.cmdRunner)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 1, c.Name())
		level.Error(c.logger).Log("err", err)
		return
	}

	level.Debug(c.logger).Log("msg", "agent metrics collected")

	ch <- prometheus.MustNewConstMetric(c.collectorError, prometheus.GaugeValue, 0, c.Name())

	c.updateMetrics(metrics, ch)
}

func collectAgentMetrics(c *cmd.CmdRunner) (*agentMetrics, error) {
	metrics := &agentMetrics{
		AgentsInfo:       c.AgentsInfo(),
		OnlineAgentsInfo: c.OnlineAgentsInfo(),
	}

	return metrics, nil
}

func (c *agentCollector) updateMetrics(values *agentMetrics, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.agentsDefined, prometheus.GaugeValue, float64(values.AgentsInfo.DefinedAgents))
	ch <- prometheus.MustNewConstMetric(c.agentsLogged, prometheus.GaugeValue, float64(values.AgentsInfo.LoggedAgents))
	ch <- prometheus.MustNewConstMetric(c.agentsTalking, prometheus.GaugeValue, float64(values.AgentsInfo.TalkingAgents))

	ch <- prometheus.MustNewConstMetric(c.onlineAgentsLogged, prometheus.GaugeValue, float64(values.OnlineAgentsInfo.OnlineLoggedAgents))
	ch <- prometheus.MustNewConstMetric(c.onlineAgentsTalking, prometheus.GaugeValue, float64(values.OnlineAgentsInfo.OnlineTalkingAgents))

	level.Debug(c.logger).Log("msg", "agent metrics built")
}
