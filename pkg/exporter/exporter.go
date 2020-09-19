// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package exporter

import (
	"context"
	"fmt"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
)

// This map stores all the exporters. The only way to add a new exporter in
// this map is to use the `Register` method.
var exporters = map[string]newFunc{}

// newFunc is an alias for the function that can create a new exporter.
type newFunc = func() Exporter

// Register adds a new exporter to the package. This does not throw an
// error, rather panics if the exporter with the same name is already
// registered, hence an exporter should be registered inside the init method
// of the package.
func Register(name string, fn newFunc) {
	if _, ok := exporters[name]; ok {
		panic(fmt.Errorf("exporter with same name already exists: %s", name))
	}

	exporters[name] = fn
}

// Exporter is anything that can export metrics into the database provider.
type Exporter interface {
	// Provision provisions the exporter. Creates database connection and
	// sets other configuration for the exporter.
	Provision(*appcontext.Context, Provider) error

	// Export is the function that does the actual exporting.
	Export(context.Context, []checker.Metric) error
}

// exportFunc is the function that is used to export the metrics into the
// provider.
type exportFunc = func(context.Context, []checker.Metric) error

// Initialize method initializes the exporter and returns a function that
// exports the metrics.
func Initialize(ctx *appcontext.Context, provider Provider) (exportFunc, error) {
	name := provider.GetBackend()
	newExporter, ok := exporters[name]
	if !ok {
		return nil, fmt.Errorf("exporter with name does not exist: %s", name)
	}

	exporter := newExporter()

	if err := exporter.Provision(ctx, provider); err != nil {
		return nil, err
	}

	return exporter.Export, nil
}
