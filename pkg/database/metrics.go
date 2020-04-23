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

var standaloneUserID uint

type controllerFunc = func(context.Context) (controller.FunctionResult, error)

func getControllerDoFunc(ex *metrics.TimescaleExporter) controllerFunc {
	return func(context.Context) (controller.FunctionResult, error) {
		start := time.Now()

		stats := ex.PullLatestControllerStatistics()
		if len(stats) == 0 {
			return nil, nil
		}

		metricsToInsert := []Metric{}
		for _, st := range stats {
			checkID, err := strconv.Atoi(st.Name)
			if err != nil {
				log.WithFields(log.Fields{
					"check_name": st.Name,
					"check_type": st.Type,
				}).WithError(err).Errorln("cannot insert check")
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

	// set the standalone user ID that can be used globally since there
	// can be only one standalone user and this is updated only when
	// setting up the metrics.
	standaloneUserID = uid

	return &metrics.TimescaleExporter{
		Manager:  manager,
		Interval: conf.Interval,
		Quit:     make(chan bool),
	}, nil
}

func setupMetrics(ex *metrics.TimescaleExporter) {
	timescaleManager := controller.NewManager()
	doFunc, err := controller.NewControllerFunction(getControllerDoFunc(ex))
	if err != nil {
		log.WithError(err).Fatalln("cannot setup timescale metrics")
		return
	}
	executor := controller.Internal{
		DoFunc:      doFunc,
		RunInterval: ex.Interval,
	}
	if err := timescaleManager.UpdateController("timescale-exporter", "exporter", executor); err != nil {
		log.WithError(err).Fatalln("cannot setup timescale metrics")
		return
	}
	timescaleManager.Wait()
}

// SetupMetrics sets up the timescale metrics for agent.
func SetupMetrics(conf *metrics.ProviderConfig, manager *controller.Manager) {
	ex, err := setupMetricsExporter(conf, manager)
	if err != nil {
		log.WithError(err).Fatalln("cannot setup timescale metrics")
		return
	}

	go setupMetrics(ex)
}

// AddCheckToDB creates new check in database from config.CheckConf.
func AddCheckToDB(checkConfig *config.CheckConf) error {
	payloads := []Payload{}
	for _, payload := range checkConfig.GetPayloads() {
		payloads = append(payloads, Payload{
			Type:  payload.GetType(),
			Value: payload.GetValue(),
		})
	}

	check := &Check{
		OwnerID:     standaloneUserID,
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

	createdCheck, err := CreateCheck(check)
	if err != nil {
		return err
	}

	checkConfig.ID = createdCheck.ID
	return nil
}
