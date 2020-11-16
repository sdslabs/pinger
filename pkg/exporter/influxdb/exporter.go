// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package influxdb

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
	provider "github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "influxdb"

// Exporter for exporting metrics to influxdb.
type Exporter struct {
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	dbname   string
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
			log.Printf("Write Error: %s\n", err.Error())
		}
	}()
	for _, metric := range metrics {
		tags := map[string]string{
			provider.CheckID:   metric.GetCheckID(),
			provider.StartTime: metric.GetStartTime().String(),
		}
		fields := map[string]interface{}{
			provider.Duration:     metric.GetDuration().String(),
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
	startTime := time.Now().Add(-1 * duration).Format(time.RFC3339Nano)

	query := fmt.Sprintf(`from(bucket:"%s")
	|> range(start: %s) 
	|> filter(fn: (r) => r._measurement == "metrics" and r.check_id =~ /%s/ )
	|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value") 
	`, bucketName, startTime, regexID(checkIDs))

	result, err := e.queryAPI.Query(ctx, query)
	if err == nil {
		for result.Next() {

			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	} else {
		fmt.Printf("Query error: %s\n", err.Error())
	}
	_, err = e.queryAPI.Query(ctx, query)
	if err != nil {
		fmt.Println("err")
	}
	return nil, nil
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
func (e *Exporter) GetMetrics(ctx context.Context, time time.Duration, checkIDs ...string) (map[string][]checker.Metric, error) {
	if len(checkIDs) == 0 {
		return nil, nil
	}
	ids := make([]string, len(checkIDs))
	copy(ids, checkIDs)

	// How to set bucket name dynamically ?
	// basically
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
	e.writeAPI = client.WriteAPI("", provider.GetDBName())
	e.queryAPI = client.QueryAPI("")
	e.dbname = provider.GetDBName()

	return nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
