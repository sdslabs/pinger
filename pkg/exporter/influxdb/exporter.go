// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package influxdb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "influxdb"

// Exporter for exporting metrics to timescale db.
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
	addStr := fmt.Sprintf("http://%s:%s", provider.GetHost(), strconv.Itoa(int(provider.GetPort())))
	authStr := fmt.Sprintf("%s:%s", provider.GetUsername(), provider.GetPassword())
	c := client.NewClientWithOptions(
		addStr,
		authStr,
		client.DefaultOptions().SetBatchSize(50),
	)

	return c, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	errorsCh := e.writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			fmt.Printf("write error: %s\n", err.Error())
		}
	}()
	for _, metric := range metrics {
		tags := map[string]string{
			"check_id":    metric.GetCheckID(),
			"start_time ": metric.GetStartTime().String(),
		}
		fields := map[string]interface{}{
			"duration":      metric.GetDuration().String(),
			"check_name":    metric.GetCheckName(),
			"is_successful": metric.IsSuccessful(),
			"is_timeout":    metric.IsTimeout(),
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
	checkIDs []string,
	duration time.Duration,
) (map[string][]checker.Metric, error) {
	// metrics := map[string][]checker.Metric{}
	a := time.Duration(15) * time.Minute
	// fmt.Println(a)
	ids := "["
	for _, v := range checkIDs {
		ids += v
	}
	ids += "]"
	query := fmt.Sprintf(`from(bucket:"pinger")|> range(start: -%s) |> filter(fn: (r) => r._measurement == "metrics" and r.check_id =~ /%s/ )`, a.String(), ids)
	// fmt.Println(query)
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
	// Close client
	_, err = e.queryAPI.Query(ctx, query)
	if err != nil {
		fmt.Println("err")
	}
	return nil, nil
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
	IDs := make([]string, len(checkIDs))
	copy(IDs, checkIDs)

	return e.getMetricsByChecksAndDuration(ctx, IDs, time)
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
