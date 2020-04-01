package agent

import (
	"fmt"
	"time"

	"github.com/sdslabs/status/pkg/check"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/database"
	"github.com/sdslabs/status/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

// ControllerManager is the global manager for the controller that comes with the
// agent. It is initialized when we run the GRPC servers.
var ControllerManager *controller.Manager

// RunStandaloneAgent runs a standalone status page agent with the provided agent config
// and the metrics config.
func RunStandaloneAgent(conf *config.AgentConfig) {
	log.Infof("Starting to run agent in standalone mode")

	ControllerManager = controller.NewManager()

	switch conf.Metrics.Backend {
	case metrics.PrometheusProviderType:
		metrics.SetupPrometheusMetrics(&conf.Metrics, ControllerManager)
	case metrics.TimeScaleProviderType:
		database.SetupMetrics(&database.SetupMetricsConf{
			ProviderConfig: &conf.Metrics,
			Manager:        ControllerManager,
			Standalone:     true,
			Checks:         conf.Checks,
		})
	case metrics.LogProviderType:
		metrics.SetupLogMetrics(&conf.Metrics, ControllerManager)
	default:
		log.Fatalf("Invalid metrics provider '%v'", conf.Metrics.Backend)
		return
	}

	log.Info("Creating contorllers for checks to be performed.")
	for _, checkConfig := range conf.Checks {
		checker, err := check.NewChecker(checkConfig)
		if err != nil {
			log.Errorf("Error while creating new checker: %s", err)
			log.Errorf("Skipping adding controller for check: %s", checkConfig.GetName())
			continue
		}

		cFunc, err := controller.NewControllerFunction(checker.ExecuteCheck)
		if err != nil {
			log.Errorf("Error while creating controller function: %s", err)
			log.Errorf("Skipping adding controller for check: %s", checkConfig.GetName())
			continue
		}

		executor := controller.Internal{
			DoFunc:      cFunc,
			RunInterval: time.Duration(checkConfig.GetInterval()),
		}
		controllerName := fmt.Sprint(checkConfig.GetId())
		if conf.Metrics.Backend != metrics.TimeScaleProviderType {
			controllerName = checkConfig.GetLabel()
		}
		err = ControllerManager.UpdateController(controllerName, checker.Type(), executor)
		if err != nil {
			log.Errorf("Error while creating controller: %s", err)
			log.Errorf("Skipping adding controller for check: %s", checkConfig.GetName())
			continue
		}
	}

	ControllerManager.Wait()
}
