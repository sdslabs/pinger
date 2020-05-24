// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package log

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/sdslabs/status/internal/appcontext"
	"github.com/sdslabs/status/internal/checker"
	"github.com/sdslabs/status/internal/controller"
	"github.com/sdslabs/status/internal/metrics"
)

const exporterName = "LOG"

func init() {
	metrics.Register(exporterName, func() metrics.Exporter { return new(Exporter) })
}

// Exporter logs the metrics to the console.
//
// Internally, it uses logrus as the logger. This is useful for testing or
// piping the logs into a file.
type Exporter struct {
	formatter logrus.Formatter
}

// PrepareChecks lets exporter handle the checks registered with the
// standalone mode.
func (e *Exporter) PrepareChecks([]checker.MutableCheck) error { return nil }

// Provision sets e's configuration.
func (e *Exporter) Provision(provider metrics.Provider) error {
	if provider.GetBackend() != exporterName {
		return fmt.Errorf(
			"invalid exporter name: expected '%s'; got '%s'",
			exporterName,
			provider.GetBackend(),
		)
	}

	if provider.GetInterval() <= 0 {
		return fmt.Errorf("interval should be > 0")
	}

	var formatter logrus.Formatter
	if provider.GetDBName() == "JSON" {
		formatter = new(logrus.JSONFormatter)
	} else {
		formatter = new(logrus.TextFormatter)
	}

	e.formatter = formatter
	return nil
}

// ExporterFunc returns the runner function that logs the metrics to the
// console.
func (e *Exporter) ExporterFunc(
	ctx *appcontext.Context,
	manager *controller.Manager,
) (controller.RunnerFunc, error) {
	logger := ctx.Logger()
	logger.SetFormatter(e.formatter)

	return func(c context.Context) (interface{}, error) {
		stats := manager.PullAllStats()

		for _, stat := range stats {
			for _, s := range stat {
				if err := s.Err; err != nil {
					logError(logger, s.ID, s.Name, err)
					continue
				}

				m, ok := s.Res.(*checker.Result)
				if !ok {
					er := fmt.Errorf("internal error: metrics not of correct type *checker.Result")
					logError(logger, s.ID, s.Name, er)
					continue
				}

				logMetric(logger, s.ID, s.Name, m)
			}
		}

		return nil, nil
	}, nil
}

// logError logs error to the console.
func logError(logger *logrus.Logger, checkID uint, checkName string, err error) {
	logger.WithError(err).Infof("metrics for check (%d) %s", checkID, checkName)
}

// logMetric logs the metric to the console.
func logMetric(logger *logrus.Logger, checkID uint, checkName string, metric *checker.Result) {
	logger.WithFields(logrus.Fields{
		"check_id":      checkID,
		"check_name":    checkName,
		"is_successful": metric.Successful,
		"is_timeout":    metric.Timeout,
		"start_time":    metric.StartTime,
		"duration":      metric.Duration,
	}).Infof("metrics for check (%d) %s", checkID, checkName)
}

// Interface guard.
var _ metrics.Exporter = (*Exporter)(nil)
