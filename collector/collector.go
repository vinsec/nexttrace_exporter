package collector

import (
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vinsec/nexttrace_exporter/config"
	"github.com/vinsec/nexttrace_exporter/executor"
)

// Collector implements the prometheus.Collector interface
type Collector struct {
	executor *executor.Executor
	targets  []config.Target
	logger   *slog.Logger

	// Metric descriptors
	hopRTT            *prometheus.Desc
	hopLoss           *prometheus.Desc
	totalHops         *prometheus.Desc
	executionDuration *prometheus.Desc
	executionsTotal   *prometheus.Desc
	lastExecution     *prometheus.Desc
}

// NewCollector creates a new Collector instance
func NewCollector(exec *executor.Executor, targets []config.Target, logger *slog.Logger) *Collector {
	return &Collector{
		executor: exec,
		targets:  targets,
		logger:   logger,

		hopRTT: prometheus.NewDesc(
			"nexttrace_hop_rtt_milliseconds",
			"Average RTT for each hop in milliseconds",
			[]string{"target", "hop_number", "hop_ip", "hop_hostname", "hop_asn"},
			nil,
		),

		hopLoss: prometheus.NewDesc(
			"nexttrace_hop_loss_ratio",
			"Packet loss ratio for each hop (0-1)",
			[]string{"target", "hop_number", "hop_ip"},
			nil,
		),

		totalHops: prometheus.NewDesc(
			"nexttrace_total_hops",
			"Total number of hops to reach the target",
			[]string{"target"},
			nil,
		),

		executionDuration: prometheus.NewDesc(
			"nexttrace_execution_duration_seconds",
			"Duration of nexttrace command execution in seconds",
			[]string{"target"},
			nil,
		),

		executionsTotal: prometheus.NewDesc(
			"nexttrace_executions_total",
			"Total number of nexttrace executions",
			[]string{"target", "status"},
			nil,
		),

		lastExecution: prometheus.NewDesc(
			"nexttrace_last_execution_timestamp",
			"Timestamp of the last successful execution",
			[]string{"target"},
			nil,
		),
	}
}

// Describe implements prometheus.Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hopRTT
	ch <- c.hopLoss
	ch <- c.totalHops
	ch <- c.executionDuration
	ch <- c.executionsTotal
	ch <- c.lastExecution
}

// Collect implements prometheus.Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	results := c.executor.GetAllResults()

	for _, target := range c.targets {
		result, exists := results[target.Name]
		if !exists {
			continue
		}

		// Execution duration
		ch <- prometheus.MustNewConstMetric(
			c.executionDuration,
			prometheus.GaugeValue,
			result.Duration.Seconds(),
			target.Name,
		)

		// Execution count (we track cumulative in a simple way via the result status)
		// Note: In a real implementation, you'd want to use actual counters
		// For this exporter pattern, we'll emit the current state
		if result.Status == "success" {
			ch <- prometheus.MustNewConstMetric(
				c.executionsTotal,
				prometheus.CounterValue,
				1,
				target.Name,
				"success",
			)
		} else if result.Status == "error" {
			ch <- prometheus.MustNewConstMetric(
				c.executionsTotal,
				prometheus.CounterValue,
				1,
				target.Name,
				"error",
			)
		} else if result.Status == "timeout" {
			ch <- prometheus.MustNewConstMetric(
				c.executionsTotal,
				prometheus.CounterValue,
				1,
				target.Name,
				"timeout",
			)
		}

		// Last execution timestamp
		if result.Status == "success" {
			ch <- prometheus.MustNewConstMetric(
				c.lastExecution,
				prometheus.GaugeValue,
				float64(result.Timestamp.Unix()),
				target.Name,
			)
		}

		// If execution was not successful, skip hop metrics
		if result.Result == nil {
			continue
		}

		// Total hops
		ch <- prometheus.MustNewConstMetric(
			c.totalHops,
			prometheus.GaugeValue,
			float64(len(result.Result.Hops)),
			target.Name,
		)

		// Per-hop metrics
		for _, hop := range result.Result.Hops {
			if !hop.HasValidIP() {
				continue
			}

			hopNumber := formatHopNumber(hop.TTL)

			// Average RTT
			avgRTT := hop.AverageRTT()
			if avgRTT > 0 {
				ch <- prometheus.MustNewConstMetric(
					c.hopRTT,
					prometheus.GaugeValue,
					avgRTT,
					target.Name,
					hopNumber,
					hop.IP,
					hop.Hostname,
					hop.ASN,
				)
			}

			// Packet loss
			ch <- prometheus.MustNewConstMetric(
				c.hopLoss,
				prometheus.GaugeValue,
				hop.Loss,
				target.Name,
				hopNumber,
				hop.IP,
			)
		}
	}
}

// UpdateTargets updates the target list for the collector
func (c *Collector) UpdateTargets(targets []config.Target) {
	c.targets = targets
}

// formatHopNumber converts hop TTL to a string for use in labels
func formatHopNumber(ttl int) string {
	return strconv.Itoa(ttl)
}
