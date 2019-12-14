package agent

import (
	"time"

	"github.com/sdslabs/status/pkg/check"
	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/metrics"

	log "github.com/sirupsen/logrus"
)

// ControllerManager is the global manager for the controller that comes with the
// agent. It is initialized when we run the GRPC servers.
var ControllerManager *controller.Manager

// RunStandaloneAgent runs a standalone status page agent with the provided agent config
// and the metrics config.
func RunStandaloneAgent(config *config.AgentConfig, metricsConfig *metrics.ProviderConfig) {
	log.Infof("Starting to run agent in standalone mode")

	ControllerManager = controller.NewManager()

	switch metricsConfig.PType {
	case metrics.PrometheusProviderType:
		metrics.SetupPrometheusMetrics(metricsConfig, ControllerManager)
	case metrics.TimeScaleProviderType:
	case metrics.EmptyProviderType:
	default:
	}

	log.Info("Creating contorllers for checks to be performed.")
	for _, checkConfig := range config.Checks {
		log.Debugf("Creating controller for check: %s", checkConfig.GetName())
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

		executor := controller.ControllerInternal{
			DoFunc:      cFunc,
			RunInterval: time.Second * time.Duration(checkConfig.GetInterval()),
		}
		err = ControllerManager.UpdateController(checkConfig.GetName(), checker.Type(), executor)
		if err != nil {
			log.Errorf("Error while creating controller: %s", err)
			log.Errorf("Skipping adding controller for check: %s", checkConfig.GetName())
			continue
		}

		log.Debugf("Controller Added for check %s", checkConfig.GetName())
	}

	ControllerManager.Wait()
}