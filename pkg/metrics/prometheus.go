package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/controller"
)

const (
	probeLatencyLabel = "probe_latency"
	probeStatusLabel  = "probe_status"
)

// SetupPrometheusMetrics sets up the prometheus metric server for the status page.
func SetupPrometheusMetrics(config *ProviderConfig, manager *controller.Manager) {
	go runPrometheusMetricsServer(config.Port, config.Interval, manager)
}

// PrometheusExporter for exporting metrics to prometheus db.
type PrometheusExporter struct {
	Manager *controller.Manager
	Metrics map[string]*prometheus.Desc
}

// NewPrometheusExporter creates an empty but not-nil `*PrometheusExporter`.
func NewPrometheusExporter(manager *controller.Manager) *PrometheusExporter {
	metrics := make(map[string]*prometheus.Desc)

	// Probe latency metrics descriptor.
	metrics[probeLatencyLabel] = prometheus.NewDesc(
		// Name of the metrics defined by the descriptor
		"status_probe_latency",
		// Help message for the metrics
		"Time in micro seconds which measures the latency of the probe defined by the controller",
		// Metrics variable level dimensions
		[]string{"probe_type", "check_name"},
		// Metrics constant label dimensions.
		nil,
	)

	// Status metrics descriptor. Value is 0.0 if the probe failed,
	// -1.0 if it failed due to timeout and +1.0 if it succeeded.
	metrics[probeStatusLabel] = prometheus.NewDesc(
		"status_probe_status",
		"Status of the probe, value is 0.0 if the probe failed, -1.0 if it failed due to timeout and +1.0 if it succeeded",
		[]string{"probe_type", "check_name"},
		nil,
	)

	return &PrometheusExporter{
		Manager: manager,
		Metrics: metrics,
	}
}

// Describe implements prometheus.Collector.
func (e *PrometheusExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.Metrics {
		ch <- desc
	}
}

// Collect implements prometheus.Collector.
func (e *PrometheusExporter) Collect(ch chan<- prometheus.Metric) {
	log.Info("Starting to collect prometheus metrics.")
	// Take the current status gained by controller from manager.
	stats := e.Manager.PullOnlyLatestControllerStatistics()

	// Iterate over statistics of each controller and send them to prometheus.
	for _, cStats := range stats {
		latencyMetric := prometheus.MustNewConstMetric(
			e.Metrics[probeLatencyLabel],
			prometheus.GaugeValue,
			float64(cStats.Duration/1e3),

			// variable labels for the metric.
			cStats.Type,
			cStats.Name,
		)
		ch <- prometheus.NewMetricWithTimestamp(cStats.StartTime, latencyMetric)

		var statusVal float64 = 1.0
		if cStats.Timeout {
			statusVal = -1.0
		} else if !cStats.Success {
			statusVal = 0.0
		}
		statusMetric := prometheus.MustNewConstMetric(
			e.Metrics[probeStatusLabel],
			prometheus.GaugeValue,
			statusVal,

			cStats.Type,
			cStats.Name,
		)
		ch <- prometheus.NewMetricWithTimestamp(cStats.StartTime, statusMetric)
	}
}

func getPrometheusControllerDoFunc(manager *controller.Manager) controllerFunc {
	return func(context.Context) (controller.FunctionResult, error) {
		manager.CleanStats()
		return &FunctionResult{
			Duration:  0,
			StartTime: time.Now(),
			Timeout:   false,
			Success:   true,
		}, nil
	}
}

func runPrometheusMetricsServer(port int, interval time.Duration, manager *controller.Manager) {
	exporter := NewPrometheusExporter(manager)
	prometheus.MustRegister(exporter)

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	defer httpServer.Close() //nolint:errcheck

	http.Handle("/metrics", promhttp.Handler())
	log.Infoln("Beginning to serve prometheus metrics on port:", port)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error while running prometheus metrics server, exitting: %s", err.Error())
			return
		}
	}()

	// Clean stats regularly from prometheus to avoid over-flow of memory usage.
	prometheusManager := controller.NewManager()
	doFunc, err := controller.NewControllerFunction(getPrometheusControllerDoFunc(prometheusManager))
	if err != nil {
		log.Errorf("Error while starting prometheus exporter: %s", err.Error())
		return
	}
	executor := controller.Internal{
		DoFunc:      doFunc,
		RunInterval: interval,
	}
	if err := prometheusManager.UpdateController("prometheus-exporter", "exporter", executor); err != nil {
		log.Errorf("Error while starting prometheus exporter: %s", err.Error())
		return
	}
	prometheusManager.Wait()
}
