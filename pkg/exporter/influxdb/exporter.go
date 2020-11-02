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
			fmt.Println("write error: %s\n", err.Error())
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

	return nil, nil
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

	return nil
}

// Interface guard.
var _ exporter.Exporter = (*Exporter)(nil)
