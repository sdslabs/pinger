package timescale

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
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

// GetCheckID returns the check ID.
func (m Metric) GetCheckID() string {
	return m.CheckID
}

// GetCheckName returns the check name.
func (m Metric) GetCheckName() string {
	return m.CheckName
}

// GetStartTime returns the start time.
func (m Metric) GetStartTime() time.Time {
	return m.StartTime
}

// GetDuration returns the duration.
func (m Metric) GetDuration() time.Duration {
	return m.Duration
}

// IsTimeout tells if the check timed out.
func (m Metric) IsTimeout() bool {
	return m.Timeout
}

// IsSuccessful tells if the check was successful.
func (m Metric) IsSuccessful() bool {
	return m.Success
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
	if len(metrics) == 0 {
		return nil
	}

	toInsert := make([]Metric, 0, len(metrics))
	for i := range metrics {
		m := metrics[i]
		toInsert = append(toInsert, Metric{
			CheckID:   m.GetCheckID(),
			CheckName: m.GetCheckName(),
			StartTime: m.GetStartTime(),
			Duration:  m.GetDuration(),
			Timeout:   m.IsTimeout(),
			Success:   m.IsSuccessful(),
		})
	}

	return e.connection.WithContext(ctx).Create(&toInsert).Error
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
	var fetched []Metric
	tx := e.connection.WithContext(ctx).Where("check_id IN (?) AND start_time > ?", checkIDs, startTime).
		Order("start_time DESC").
		Find(&fetched)

	metrics := map[string][]checker.Metric{}
	for i := range fetched {
		m := fetched[i]

		if _, ok := metrics[m.CheckID]; !ok {
			metrics[m.CheckID] = []checker.Metric{}
		}

		metrics[m.CheckID] = append(metrics[m.CheckID], m)
	}

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

	return e.getMetricsByChecksAndDuration(ctx, checkIDs, time)
}

// Interface guards.
var (
	_ exporter.Exporter = (*Exporter)(nil)
	_ checker.Metric    = Metric{}
)
