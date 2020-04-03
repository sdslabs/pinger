package database

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/metrics"
)

type controllerFunc = func(context.Context) (controller.FunctionResult, error)

func getControllerDoFunc(ex *metrics.TimescaleExporter) controllerFunc {
	return func(context.Context) (controller.FunctionResult, error) {
		start := time.Now()
		log.Infoln("Trying to insert metrics into DB")
		stats := ex.PullLatestControllerStatistics()
		if len(stats) == 0 {
			return nil, nil
		}
		metricsToInsert := []Metric{}
		for _, st := range stats {
			checkID, err := strconv.Atoi(st.Name)
			if err != nil {
				log.Errorf("Failed to insert metric for check=%s", st.Name)
				continue
			}
			metricsToInsert = append(metricsToInsert, Metric{
				CheckID:   uint(checkID),
				StartTime: &st.StartTime,
				Duration:  st.Duration,
				Timeout:   st.Timeout,
				Success:   st.Success,
			})
		}
		if err := CreateMetrics(metricsToInsert); err != nil {
			return nil, err
		}
		return &metrics.FunctionResult{
			Duration:  time.Since(start),
			StartTime: start,
			Success:   true,
			Timeout:   false,
		}, nil
	}
}

func setupMetricsExporter(conf *metrics.ProviderConfig, manager *controller.Manager) (*metrics.TimescaleExporter, error) {
	if err := initFromProvider(conf); err != nil {
		return nil, err
	}
	uid, err := createStandaloneUser()
	if err != nil {
		return nil, err
	}
	return &metrics.TimescaleExporter{
		Manager:  manager,
		UserID:   uid,
		Interval: conf.Interval,
		Quit:     make(chan bool),
	}, nil
}

func setupMetrics(ex *metrics.TimescaleExporter) {
	timescaleManager := controller.NewManager()
	doFunc, err := controller.NewControllerFunction(getControllerDoFunc(ex))
	if err != nil {
		log.Fatalf("Error setting up timescale metrics: %s", err.Error())
		return
	}
	executor := controller.Internal{
		DoFunc:      doFunc,
		RunInterval: ex.Interval,
	}
	if err := timescaleManager.UpdateController("timescale-exporter", "exporter", executor); err != nil {
		log.Fatalf("Error setting up timescale metrics: %s", err.Error())
		return
	}
	timescaleManager.Wait()
}

// SetupMetricsConf defines the setup configuration for timescale metrics.
type SetupMetricsConf struct {
	*metrics.ProviderConfig
	*controller.Manager

	Standalone bool
	Checks     []*config.CheckConf
}

// SetupMetrics sets up the timescale metrics for agent.
func SetupMetrics(conf *SetupMetricsConf) {
	ex, err := setupMetricsExporter(conf.ProviderConfig, conf.Manager)
	if err != nil {
		log.Fatalf("Cannot setup metrics: %s", err.Error())
		return
	}

	// Insert checks into DB before starting the metrics collector for standalone mode.
	if conf.Standalone {
		for _, checkConfig := range conf.Checks {
			payloads := []Payload{}
			for _, payload := range checkConfig.GetPayloads() {
				payloads = append(payloads, Payload{
					Type:  payload.GetType(),
					Value: payload.GetValue(),
				})
			}
			check := &Check{
				OwnerID:     ex.UserID,
				Interval:    time.Duration(checkConfig.GetInterval()),
				Timeout:     time.Duration(checkConfig.GetTimeout()),
				InputType:   checkConfig.GetInput().GetType(),
				InputValue:  checkConfig.GetInput().GetValue(),
				OutputType:  checkConfig.GetOutput().GetType(),
				OutputValue: checkConfig.GetOutput().GetValue(),
				TargetType:  checkConfig.GetTarget().GetType(),
				TargetValue: checkConfig.GetTarget().GetValue(),
				Title:       checkConfig.GetName(),
				Payloads:    payloads,
			}
			log.Infof("Inserting check '%s' inside DB", check.Title)
			createdCheck, err := CreateCheck(check)
			if err != nil {
				log.Errorf("Failed to insert check '%s'", check.Title)
				log.Errorln("Skipping and moving ahead...")
				continue
			}
			checkConfig.ID = createdCheck.ID
		}
	}

	go setupMetrics(ex)
}
