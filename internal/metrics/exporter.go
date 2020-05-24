// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package metrics

import (
	"fmt"

	"github.com/sdslabs/status/internal/appcontext"
	"github.com/sdslabs/status/internal/checker"
	"github.com/sdslabs/status/internal/controller"
)

// This map stores all the exporters. The only way to add a new exporter in
// this map is to use the `Register` method.
var exporters = map[string]newFunc{}

// newFunc is an alias for the function that can create a new checker.
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
	// PrepareChecks lets exporter handle the checks registered with the
	// standalone mode.
	PrepareChecks([]checker.MutableCheck) error

	// Provision provisions the exporter. Creates database connection and
	// sets other configuration for the exporter.
	Provision(Provider) error

	// ExporterFunc returns the controller runner function that is run by
	// the exporter at regular intervals and actually does the exporting of
	// metrics.
	ExporterFunc(*appcontext.Context, *controller.Manager) (controller.RunnerFunc, error)
}

// Initialize method initializes the exporter.
func Initialize(
	ctx *appcontext.Context,
	manager *controller.Manager,
	provider Provider,
	checks []checker.MutableCheck,
) error {
	name := provider.GetBackend()
	newExporter, ok := exporters[name]
	if !ok {
		return fmt.Errorf("exporter with name does not exist: %s", name)
	}

	if provider.GetInterval() <= 0 {
		return fmt.Errorf("interval should be > 0")
	}

	exporter := newExporter()

	if err := exporter.PrepareChecks(checks); err != nil {
		return err
	}

	if err := exporter.Provision(provider); err != nil {
		return err
	}

	runnerFunc, err := exporter.ExporterFunc(ctx, manager)
	if err != nil {
		return err
	}

	ctrl, err := controller.NewController(ctx, &controller.Opts{
		Name:     fmt.Sprintf("exporter_%s", provider.GetBackend()),
		Interval: provider.GetInterval(),
		Func:     runnerFunc,
	})
	if err != nil {
		return err
	}

	ctrl.Start()
	return nil
}
