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
	"github.com/sdslabs/status/internal/exporter"
)

const exporterName = "LOG"

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// Exporter logs the metrics to the console.
//
// Internally, it uses logrus as the logger. This is useful for testing or
// piping the logs into a file.
type Exporter struct {
	logger    *logrus.Logger
	formatter logrus.Formatter
}

// PrepareChecks lets exporter handle the checks registered with the
// standalone mode.
func (e *Exporter) PrepareChecks([]checker.MutableCheck) error { return nil }

// Provision sets e's configuration.
func (e *Exporter) Provision(ctx *appcontext.Context, provider exporter.Provider) error {
	if provider.GetBackend() != exporterName {
		return fmt.Errorf(
			"invalid exporter name: expected '%s'; got '%s'",
			exporterName,
			provider.GetBackend(),
		)
	}

	var formatter logrus.Formatter
	if provider.GetDBName() == "JSON" {
		formatter = new(logrus.JSONFormatter)
	} else {
		formatter = new(logrus.TextFormatter)
	}

	e.formatter = formatter
	e.logger = ctx.Logger()
	return nil
}

// Export logs the metrics onto the console.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	e.logger.SetFormatter(e.formatter)

	for _, metric := range metrics {
		e.logMetric(metric)
	}

	return nil
}

// logMetric logs the metric to the console.
func (e *Exporter) logMetric(metric checker.Metric) {
	e.logger.WithFields(logrus.Fields{
		"check_id":      metric.GetCheckID(),
		"check_name":    metric.GetCheckName(),
		"is_successful": metric.IsSuccessful(),
		"is_timeout":    metric.IsTimeout(),
		"start_time":    metric.GetStartTime(),
		"duration":      metric.GetDuration(),
	}).Infof("metrics for check (%d) %s", metric.GetCheckID(), metric.GetCheckName())
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
