// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const (
	probeLatencyLabel = "probe_latency"
	probeStatusLabel  = "probe_status"
	exporterName      = "prometheus"
)

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// Exporter for exporting metrics to prometheus db.
type Exporter struct {
	Metrics map[string]*prometheus.Desc
	checks  []checker.Metric
}

// NewExporter creates an empty but not-nil `*Exporter`.
func NewExporter() *Exporter {
	metrics := make(map[string]*prometheus.Desc)

	// Probe latency metrics descriptor.
	metrics[probeLatencyLabel] = prometheus.NewDesc(
		// Name of the metrics defined by the descriptor
		probeLatencyLabel,
		// Help message for the metrics
		"Time in micro seconds which measures the latency of the probe defined by the controller",
		// Metrics variable level dimensions
		[]string{"check_name"},
		// Metrics constant label dimensions.
		nil,
	)

	// Status metrics descriptor.
	//
	// Value is 0.0 if the probe failed,
	// -1.0 if it failed due to timeout and +1.0 if it succeeded.
	metrics[probeStatusLabel] = prometheus.NewDesc(
		probeStatusLabel,
		"Status of the probe, value is 0.0 on probe failure, -1.0 on failure due to timeout and +1.0 on success",
		[]string{"check_name"},
		nil,
	)

	return &Exporter{
		Metrics: metrics,
		checks:  nil,
	}
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.Metrics {
		ch <- desc
	}
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Info("Starting to collect prometheus metrics.")

	// Iterate over each check and send them to prometheus.
	for _, check := range e.checks {
		latencyMetric := prometheus.MustNewConstMetric(
			e.Metrics[probeLatencyLabel],
			prometheus.GaugeValue,
			float64(check.GetDuration()/1e3),

			// variable labels for the metric.
			check.GetCheckName(),
		)
		ch <- prometheus.NewMetricWithTimestamp(check.GetStartTime(), latencyMetric)

		var statusVal float64 = 1.0
		if check.IsTimeout() {
			statusVal = -1.0
		} else if !check.IsSuccessful() {
			statusVal = 0.0
		}
		statusMetric := prometheus.MustNewConstMetric(
			e.Metrics[probeStatusLabel],
			prometheus.GaugeValue,
			statusVal,

			check.GetCheckName(),
		)
		ch <- prometheus.NewMetricWithTimestamp(check.GetStartTime(), statusMetric)
	}
}

// Provision sets e's configuration.
func (e *Exporter) Provision(ctx *appcontext.Context, provider exporter.Provider) error {
	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", provider.GetPort()),
	}
	defer httpServer.Close() //nolint:errcheck

	http.Handle("/metrics", promhttp.Handler())
	log.Infoln("Beginning to serve prometheus metrics on port", provider.GetPort())

	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(_ context.Context, metrics []checker.Metric) error {
	ch1 := make(chan<- *prometheus.Desc)
	e.Describe(ch1)

	e.checks = metrics
	ch2 := make(chan<- prometheus.Metric)
	e.Collect(ch2)

	return nil
}

// GetMetrics get the metrics of the given checks.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string,
) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}

	return nil, nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
