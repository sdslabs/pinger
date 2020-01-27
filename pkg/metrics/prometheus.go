package metrics

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/controller"
)

// SetupPrometheusMetrics sets up the prometheus metric server for the status page.
func SetupPrometheusMetrics(config *ProviderConfig, manager *controller.Manager) {
	go runPrometheusMetricsServer(config.Port, manager)
}

// PrometheusExporter for exporting metrics to prometheus db.
type PrometheusExporter struct {
	Manager *controller.Manager
	Port    int

	Metrics map[string]*prometheus.Desc
}

// NewPrometheusExporter creates an empty but not-nil `*PrometheusExporter`.
func NewPrometheusExporter(port int, manager *controller.Manager) *PrometheusExporter {
	metrics := make(map[string]*prometheus.Desc)

	// Probe latency metrics descriptor.
	metrics["probe_latency"] = prometheus.NewDesc(
		// Name of the metrics defined by the descriptor
		"probe_latency",
		// Help message for the metrics
		"Time in micro seconds which measures the latency of the probe defined by the controller",
		// Metrics variable level dimensions
		[]string{"probe_type", "check_name"},

		// Metrics constant label dimensions.
		nil,
	)

	return &PrometheusExporter{
		Manager: manager,
		Port:    port,

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
	stats := e.Manager.PullLatestControllerStatistics()

	// Iterate over statistics of each controller and send them to prometheus.
	for _, cStats := range stats {
		fmt.Println("STATS for prometheus: ", cStats.StartTime, float64(cStats.Duration/1e3))

		m := prometheus.MustNewConstMetric(
			e.Metrics["probe_latency"],
			prometheus.GaugeValue,
			float64(cStats.Duration/1e3),

			// variable labels for the metric.
			cStats.Type,
			cStats.Name,
		)

		ch <- prometheus.NewMetricWithTimestamp(cStats.StartTime, m)
	}
}

func runPrometheusMetricsServer(port int, manager *controller.Manager) {
	exporter := NewPrometheusExporter(port, manager)
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve prometheus metrics on port :", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Error("Error while running prometheus metrics server, exitting: ", err)
		os.Exit(1)
	}
}
