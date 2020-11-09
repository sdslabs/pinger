package alerter

import (
	"context"
	"fmt"

	"github.com/sdslabs/pinger/pkg/appcontext"
	"github.com/sdslabs/pinger/pkg/checker"
)

// This map stores all the alterters. The only way to add a new alerter in
// this map is to use the `Register` method.
var alerters = map[string]newFunc{}

// newFunc is an alias for the function that can create a new alterter.
type newFunc = func() Alerter

// Register adds a new alerter to the package. This does not throw an
// error, rather panics if the alerter with the same name is already
// registered, hence an alerter should be registered inside the init method
// of the package.
func Register(name string, fn newFunc) {
	if _, ok := alerters[name]; ok {
		panic(fmt.Errorf("alerter with same name already exists: %s", name))
	}

	alerters[name] = fn
}

// Alerter is responsible for sending alerts.
type Alerter interface {
	// Provision sets the alerter's configuration.
	Provision(*appcontext.Context, Provider) error

	// Alert actually sends the alert.
	//
	// It gets the metrics for which the alert needs to be sent (all metrics
	// are filtered before passing onto here). The last parameter it has is
	// a map of check IDs and corresponding alert to be sent.
	Alert(context.Context, []checker.Metric, map[string]Alert) error
}

// AlertFunc is the function that is used to alert the metrics into the
// provider.
type AlertFunc = func(context.Context, []checker.Metric, map[string]Alert) error

// Initialize method initializes the alerter and returns a function that
// alerts the metrics.
func Initialize(ctx *appcontext.Context, provider Provider) (AlertFunc, error) {
	name := provider.GetService()
	newAlerter, ok := alerters[name]
	if !ok {
		return nil, fmt.Errorf("alerter with name does not exist: %s", name)
	}

	alerter := newAlerter()

	if err := alerter.Provision(ctx, provider); err != nil {
		return nil, err
	}

	return alerter.Alert, nil
}
