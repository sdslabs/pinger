package metrics

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/controller"
)

// SetupLogMetrics sets up logging of metrics to STDOUT.
func SetupLogMetrics(config *ProviderConfig, manager *controller.Manager) {
	go startMetricsLogger(config, manager)
}

// LogsType can either be "json" or "text".
type LogsType string

// Types of logs.
const (
	JSONLogs LogsType = "json"
	TextLogs LogsType = "text"
)

// LogsExporter for exporting the metrics as structured logs onto the console.
type LogsExporter struct {
	*controller.Manager
	LogsType
}

// SetupLogger initializes a logrus logger for the exporter.
func (le *LogsExporter) SetupLogger() {
	if le.LogsType == JSONLogs {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	} else {
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// Log outputs the metrics on the console.
func (le *LogsExporter) Log(stat *controller.ExecutionStat) {
	logrus.WithFields(logrus.Fields{
		"name":       stat.Name,
		"type":       stat.Type,
		"start_time": stat.StartTime,
		"duration":   stat.Duration,
		"timeout":    stat.Timeout,
		"success":    stat.Success,
	}).Infoln("check results")
}

func getLogDoFunc(ex *LogsExporter) controllerFunc {
	return func(context.Context) (controller.FunctionResult, error) {
		stats := ex.PullLatestControllerStatistics()

		for i := 0; i < len(stats); i++ {
			ex.Log(&stats[i])
		}

		return &FunctionResult{
			Duration:  0,
			StartTime: time.Now(),
			Timeout:   false,
			Success:   true,
		}, nil
	}
}

func startMetricsLogger(config *ProviderConfig, manager *controller.Manager) {
	exporter := &LogsExporter{Manager: manager}
	if config.DBName == string(JSONLogs) {
		exporter.LogsType = JSONLogs
	} else {
		exporter.LogsType = TextLogs
	}

	exporter.SetupLogger()

	logManager := controller.NewManager()
	doFunc, err := controller.NewControllerFunction(getLogDoFunc(exporter))
	if err != nil {
		logrus.WithError(err).Errorln("cannot run log metrics provider")
		return
	}
	executor := controller.Internal{
		DoFunc:      doFunc,
		RunInterval: config.Interval,
	}
	if err := logManager.UpdateController("log-exporter", "exporter", executor); err != nil {
		logrus.WithError(err).Errorln("cannot run log metrics provider")
		return
	}
	logManager.Wait()
}