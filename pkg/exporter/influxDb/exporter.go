// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package influxdb

import (
	"context"
	"fmt"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
	"github.com/sdslabs/pinger/pkg/exporter"
)

const exporterName = "influxdb"

// Exporter for exporting metrics to timescale db.
type Exporter struct {
	connection client.Client
}

func init() {
	exporter.Register(exporterName, func() exporter.Exporter { return new(Exporter) })
}

// newConn creates a new connection with the database.
func newClient(ctx *appcontext.Context, provider exporter.Provider) (client.Client, error) {
	addStr := fmt.Sprintf("http://%s:%s", provider.GetHost(), provider.GetPassword())

	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addStr,
		Username: provider.GetUsername(),
		Password: provider.GetPassword(),
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

// Export exports the metrics to the exporter.
func (e *Exporter) Export(ctx context.Context, metrics []checker.Metric) error {
	return e.createMetrics(ctx, metrics)
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
	e.connection = client

	return nil
}
