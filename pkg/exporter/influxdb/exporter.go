// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package influxdb

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	write "github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/sirupsen/logrus"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "influxdb"

// metric keys
const (
	keyCheckID      = "check_id"
	keyCheckName    = "check_name"
	keyIsSuccessful = "is_successful"
	keyIsTimeout    = "is_timeout"
	keyStartTime    = "start_time"
	keyDuration     = "duration"
)

// Exporter for exporting metrics to influxdb.
type Exporter struct {
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
	dbname   string
	log      *logrus.Logger
}

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// newClient creates a new client for the database.
func newClient(ctx *appcontext.Context, provider exporter.Provider) (client.Client, error) {
	protocol := "http"
	if provider.IsSSLMode() {
		protocol = "https"
	}
	addStr := fmt.Sprintf("%s://%s",
		protocol,
		net.JoinHostPort(provider.GetHost(),
			strconv.Itoa(int(provider.GetPort()))))

	authStr := fmt.Sprintf("%s:%s", provider.GetUsername(), provider.GetPassword())
	c := client.NewClientWithOptions(
		addStr,
		authStr,
		client.DefaultOptions(),
	)

	return c, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	points := make([]*write.Point, 0)
	for _, metric := range metrics {
		tags := map[string]string{
			keyCheckID: metric.GetCheckID(),
		}
		fields := map[string]interface{}{
			keyStartTime:    metric.GetStartTime(),
			keyDuration:     metric.GetDuration(),
			keyCheckName:    metric.GetCheckName(),
			keyIsSuccessful: metric.IsSuccessful(),
			keyIsTimeout:    metric.IsTimeout(),
		}
		p := client.NewPoint("metrics", tags, fields, time.Now())
		points = append(points, p)

	}
	err := e.writeAPI.WritePoint(ctx, points...)
	if err != nil {
		return err
	}

	return nil
}

// getMetricsByChecksAndDuration fetches metrics from the metrics hypertable
// for the given check IDs. It accepts a `duration` parameter that fetches
// metrics for the check in the past `duration time.Duration`.
func (e *Exporter) getMetricsByChecksAndDuration(
	ctx context.Context,
	bucketName string,
	checkIDs []string,
	duration time.Duration,
) (map[string][]checker.Metric, error) {

	var parsedDuration time.Duration
	var parsedTime time.Time
	startTime := time.Now().Add(-1 * duration).Format(time.RFC3339Nano)

	// flux query where bucketName is the database name + the retention policy
	query := fmt.Sprintf(`from(bucket:%q)
	|> range(start: %s) 
	|> filter(fn: (r) => r._measurement == "metrics" and r.check_id =~ /%s/ )
	|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value")
	|> drop(columns:["_time", "_start","_stop"])`, bucketName, startTime, regexID(checkIDs))
	metrics := map[string][]checker.Metric{}
	result, err := e.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		// duration and time are parsed separately because influx converts time.Time to string
		// hence manually parsing and assigning is done here
		parsedDuration, err = time.ParseDuration(result.Record().ValueByKey(keyDuration).(string))
		if err != nil {
			return nil, err
		}
		parsedTime, err = time.Parse(time.RFC3339Nano, result.Record().ValueByKey(keyStartTime).(string))
		if err != nil {
			return nil, err
		}
		metric := config.Metric{
			CheckID:    result.Record().ValueByKey(keyCheckID).(string),
			CheckName:  result.Record().ValueByKey(keyCheckName).(string),
			StartTime:  parsedTime,
			Duration:   parsedDuration,
			Timeout:    result.Record().ValueByKey(keyIsTimeout).(bool),
			Successful: result.Record().ValueByKey(keyIsSuccessful).(bool),
		}

		if _, ok := metrics[metric.CheckID]; !ok {
			metrics[metric.CheckID] = make([]checker.Metric, 0)
		}
		metrics[metric.CheckID] = append(metrics[metric.CheckID], &metric)
	}

	return metrics, nil
}

// TODO(h3llix): validate variable to be alphanumeric
// Formats the checkIDs for required flux query
func regexID(checkIDs []string) string {
	ids := "[%s]"
	temp := ""
	for _, v := range checkIDs {
		temp += (v + "|")
	}
	return fmt.Sprintf(ids, temp)
}

// GetMetrics get the metrics of the given checks.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}

	return e.getMetricsByChecksAndDuration(ctx, e.dbname, checkIDs, time)
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

	client, err := newClient(ctx, provider)
	if err != nil {
		return err
	}
	e.log = ctx.Logger()
	e.writeAPI = client.WriteAPIBlocking("", provider.GetDBName())
	e.queryAPI = client.QueryAPI("")
	e.dbname = provider.GetDBName()

	return nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
