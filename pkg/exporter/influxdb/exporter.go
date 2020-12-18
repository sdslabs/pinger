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
	"github.com/sirupsen/logrus"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	config "github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/exporter"
	provider "github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "influxdb"

// Exporter for exporting metrics to influxdb.
type Exporter struct {
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	dbname   string
	log      *logrus.Logger
}

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// newClient creates a new client for the database.
func newClient(ctx *appcontext.Context, provider exporter.Provider) (client.Client, error) {
	addStr := fmt.Sprintf("http://%s", net.JoinHostPort(provider.GetHost(), strconv.Itoa(int(provider.GetPort()))))
	authStr := fmt.Sprintf("%s:%s", provider.GetUsername(), provider.GetPassword())
	c := client.NewClientWithOptions(
		addStr,
		authStr,
		// TODO(h3llix): Set this Batch size with yaml
		// currently batch size is set to 50
		client.DefaultOptions().SetBatchSize(50),
	)

	return c, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	errorsCh := e.writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			e.log.Printf("Write Error: %s\n", err.Error())
		}
	}()
	for _, metric := range metrics {
		tags := map[string]string{
			provider.CheckID: metric.GetCheckID(),
		}
		fields := map[string]interface{}{
			provider.StartTime:    metric.GetStartTime(),
			provider.Duration:     metric.GetDuration(),
			provider.CheckName:    metric.GetCheckName(),
			provider.IsSuccessful: metric.IsSuccessful(),
			provider.IsTimeout:    metric.IsTimeout(),
		}
		p := client.NewPoint("metrics", tags, fields, time.Now())
		//Non-blocking write client uses implicit batching.
		e.writeAPI.WritePoint(p)

	}
	e.writeAPI.Flush()

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
	query := fmt.Sprintf(`from(bucket:"%s")
	|> range(start: %s) 
	|> filter(fn: (r) => r._measurement == "metrics" and r.check_id =~ /%s/ )
	|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value")
	|> drop(columns:["_time", "_start","_stop"])`, bucketName, startTime, regexID(checkIDs))
	metrics := map[string][]checker.Metric{}
	result, err := e.queryAPI.Query(ctx, query)
	if err == nil {
		for result.Next() {
			// duration and time are parsed separately because influx converts time.Time to string
			// hence manually parsing and assigning is done here
			parsedDuration, err = time.ParseDuration(result.Record().ValueByKey("duration").(string))
			if err != nil {
				e.log.Printf("Cannot parse duration :%s", err)
			}
			parsedTime, err = time.Parse(time.RFC3339Nano, result.Record().ValueByKey("start_time").(string))
			if err != nil {
				e.log.Printf("Cannot parse start time :%s", err)
			}
			metric := config.Metric{
				CheckID:    result.Record().ValueByKey("check_id").(string),
				CheckName:  result.Record().ValueByKey("check_name").(string),
				StartTime:  parsedTime,
				Duration:   parsedDuration,
				Timeout:    result.Record().ValueByKey("is_timeout").(bool),
				Successful: result.Record().ValueByKey("is_successful").(bool),
			}
			if _, ok := metrics[metric.CheckID]; ok {
				metrics[metric.CheckID] = append(metrics[metric.CheckID], &metric)
			} else {
				metrics[metric.CheckID] = make([]checker.Metric, 0)
			}
		}
	} else {
		e.log.Printf("Query error: %s\n", err.Error())
	}
	return metrics, nil
}

// Formats the checkIDs for required flux query
func regexID(checkIDs []string) string {
	ids := "["
	for _, v := range checkIDs {
		ids += v
	}
	ids += "]"
	return ids
}

// GetMetrics get the metrics of the given checks.
func (e *Exporter) GetMetrics(
	ctx context.Context,
	time time.Duration,
	checkIDs ...string) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}
	ids := make([]string, len(checkIDs))
	copy(ids, checkIDs)

	//TODO(h3llix): to set bucket dynamically
	bucketName := "pinger"

	return e.getMetricsByChecksAndDuration(ctx, bucketName, ids, time)
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
	e.writeAPI = client.WriteAPI("", provider.GetDBName())
	e.queryAPI = client.QueryAPI("")
	e.dbname = provider.GetDBName()

	return nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
