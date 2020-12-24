// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package log

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "log"

// metric keys
const (
	CheckID      = "check_id"
	CheckName    = "check_name"
	IsSuccessful = "is_successful"
	IsTimeout    = "is_timeout"
	StartTime    = "start_time"
	Duration     = "duration"
)

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
		CheckID:      metric.GetCheckID(),
		CheckName:    metric.GetCheckName(),
		IsSuccessful: metric.IsSuccessful(),
		IsTimeout:    metric.IsTimeout(),
		StartTime:    metric.GetStartTime(),
		Duration:     metric.GetDuration(),
	}).Infof("metrics for check (%s) %s", metric.GetCheckID(), metric.GetCheckName())
}

// GetMetrics returns error as could not be used with log exporter.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string,
) (map[string][]checker.Metric, error) {
	return nil, fmt.Errorf(
		`the method is not implemented for exporter: %s,
		log exporter is only meant to be used for debugging`,
		exporterName,
	)
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
