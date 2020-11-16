// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package timescale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "timescale"

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// Exporter for exporting metrics to timescale db.
type Exporter struct {
	connection *gorm.DB
}

// Metric model.
type Metric struct {
	CheckID   string
	CheckName string

	StartTime time.Time     `gorm:"NOT NULL"`
	Duration  time.Duration `gorm:"NOT NULL"`
	Timeout   bool          `gorm:"NOT NULL"`
	Success   bool          `gorm:"NOT NULL"`
}

// newConn creates a new connection with the database.
func newConn(ctx *appcontext.Context, provider exporter.Provider) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		`host=%s port=%d user=%s dbname=%s password=%s`,
		provider.GetHost(),
		provider.GetPort(),
		provider.GetUsername(),
		provider.GetDBName(),
		provider.GetPassword(),
	)

	if !provider.IsSSLMode() {
		connStr = fmt.Sprintf("%s sslmode=disable", connStr)
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).AutoMigrate(
		&Metric{},
	)
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Exec("CREATE INDEX ON metrics (check_id, start_time DESC);").Error
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Exec(
		"SELECT create_hypertable('metrics', 'start_time', if_not_exists => TRUE, create_default_indexes => FALSE);",
	).Error
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createMetrics inserts multiple metrics into TimescaleDB Hypertable.
func (e *Exporter) createMetrics(ctx context.Context, metrics []checker.Metric) error {
	// We build a raw query since Gorm doesn't support bulk insert.
	// Since there is no `string` and there is no user input we can
	// safely build the raw query without worrying about injection.

	if len(metrics) == 0 {
		return nil
	}

	q := "INSERT INTO metrics (check_id, check_name, start_time, duration, timeout, success) VALUES %s;"
	timeFormat := "2006-01-02 15:04:05.000000-07:00" // Supported by PostgreSQL
	vals := []string{}
	for i := range metrics {
		val := fmt.Sprintf("(%s, %s, '%s', %d, %t, %t)",
			metrics[i].GetCheckID(),
			metrics[i].GetCheckName(),
			metrics[i].GetStartTime().Format(timeFormat),
			metrics[i].GetDuration(),
			metrics[i].IsTimeout(),
			metrics[i].IsSuccessful())
		vals = append(vals, val)
	}

	args := strings.Join(vals, ", ")
	query := fmt.Sprintf(q, args)

	return e.connection.WithContext(ctx).Exec(query).Error
}

// getMetricsByChecksAndDuration fetches metrics from the metrics hypertable
// for the given check IDs. It accepts a `duration` parameter that fetches
// metrics for the check in the past `duration time.Duration`.
func (e *Exporter) getMetricsByChecksAndDuration(
	ctx context.Context,
	checkIDs []string,
	duration time.Duration,
) (map[string][]checker.Metric, error) {
	startTime := time.Now().Add(-1 * duration)
	metrics := map[string][]checker.Metric{}
	tx := e.connection.WithContext(ctx).Where("check_id IN (?) AND start_time > ?", checkIDs, startTime).
		Order("start_time DESC").
		Group("check_id").
		Find(&metrics)

	return metrics, tx.Error
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

	conn, err := newConn(ctx, provider)
	if err != nil {
		return err
	}
	e.connection = conn

	return nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	return e.createMetrics(ctx, metrics)
}

// GetMetrics get the metrics of the given checks.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string,
) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}
	ids := make([]string, len(checkIDs))
	copy(ids, checkIDs)

	return e.getMetricsByChecksAndDuration(ctx, ids, time)
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
