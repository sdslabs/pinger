package questdb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

const exporterName = "questdb"

// Exporter for exporting metrics to quest db.
type Exporter struct {
	conn *pgxpool.Pool
}

// Metric model.
type Metric struct {
	CheckID   string
	CheckName string

	StartTime time.Time
	Duration  time.Duration
	Timeout   bool
	Success   bool
}

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
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
func newConn(ctx *appcontext.Context, provider exporter.Provider) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		`host=%s  user=%s dbname=%s password=%s`,
		provider.GetHost(),
		provider.GetUsername(),
		provider.GetDBName(),
		provider.GetPassword(),
	)

	if provider.GetPort() != 0 {
		connStr = fmt.Sprintf("%s port=%d",
			connStr,
			provider.GetPort(),
		)
	}

	if !provider.IsSSLMode() {
		connStr = fmt.Sprintf("%s sslmode=disable", connStr)
	}

	db, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	_, err1 := db.Exec(ctx, "CREATE TABLE IF NOT EXISTS metrics(check_id string, check_name string, start_time timestamp, duration long,timeout string, success string) timestamp(start_time);")
	if err1 != nil {
		return nil, err1
	}

	return db, nil
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

	cli, err := newConn(ctx, provider)
	if err != nil {
		return err
	}
	e.conn = cli

	return nil
}

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

func (e *Exporter) getMetricsByChecksAndDuration(
	ctx context.Context,
	checkIDs []string,
	duration time.Duration,
) (map[string][]checker.Metric, error) {
	startTime := time.Now().Add(-1 * duration).UTC().Format(time.RFC3339)
	metrics := map[string][]checker.Metric{}

	querystring := fmt.Sprintf(` SELECT * FROM metrics WHERE
	 ( check_id= '%s' `,
		checkIDs[0],
	)
	for _, value := range checkIDs {
		querystring += fmt.Sprintf(`OR check_id= '%s'`, value)
	}
	querystring += `)`
	querystring += fmt.Sprintf(`AND (start_time> '%s') ORDER BY start_time  `, startTime)

	fetched, err := e.conn.Query(ctx, querystring)
	if err != nil {
		fmt.Print(err)
	}

	for fetched.Next() {
		var CheckID string
		var CheckName string

		var StartTime time.Time
		var Duration time.Duration
		var Timeout string
		var Success string

		err = fetched.Scan(&CheckID, &CheckName, &StartTime, &Duration, &Timeout, &Success)
		if err != nil {
			return nil, err
		}

		timeout1, err1 := strconv.ParseBool(Timeout)
		if err1 != nil {
			return nil, err
		}

		success1, err2 := strconv.ParseBool(Success)
		if err2 != nil {
			return nil, err
		}

		m := Metric{CheckID, CheckName, StartTime, Duration, timeout1, success1}

		if _, ok := metrics[m.CheckID]; !ok {
			metrics[m.CheckID] = []checker.Metric{}
		}

		metrics[m.CheckID] = append(metrics[m.CheckID], m)

	}

	return metrics, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	return e.createMetrics(ctx, metrics)
}

// createMetrics inserts multiple metrics into TimescaleDB Hypertable.
func (e *Exporter) createMetrics(ctx context.Context, metrics []checker.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	for i := range metrics {

		batch.Queue("insert into metrics(check_id,check_name, start_time,duration,timeout,success) values($1, $2, $3, $4, $5,$6)",
			metrics[i].GetCheckID(),
			metrics[i].GetCheckName(),
			metrics[i].GetStartTime(),
			metrics[i].GetDuration(),
			strconv.FormatBool(metrics[i].IsTimeout()),
			strconv.FormatBool(metrics[i].IsSuccessful()),
		)
	}

	br := e.conn.SendBatch(ctx, batch)

	_, err := br.Exec()

	return err
}
