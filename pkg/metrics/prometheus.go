package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sdslabs/status/pkg/controller"
)

type PrometheusExporter struct {
	Manager *controller.Manager
	Port    int

	Metrics map[string]*prometheus.Desc
}

func SetupPrometheusMetrics(config *ProviderConfig, manager *controller.Manager) error {
}

func runPrometheusMetricsServer(config *metrics.ProviderConfig) (*Manager, error) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :", config.Port)

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
}
