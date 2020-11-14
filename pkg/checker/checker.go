package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/sdslabs/pinger/pkg/util/controller"
)

// This map stores all the checkers. The only way to add a new checker in
// this map is to use the `Register` method.
var checkers = map[string]newFunc{}

// Various errors thrown by a checker.
var (
	ErrValidation   = fmt.Errorf("validation error")
	ErrProvisioning = fmt.Errorf("provisioning error")
)

// newFunc is an alias for the function that can create a
// new checker.
type newFunc = func() Checker

// Register adds a new checker to the package. This does not throw an error,
// rather panics if the checker with the same name is already registered,
// hence a checker should be registered inside the init method of the
// package.
func Register(name string, fn newFunc) {
	if _, ok := checkers[name]; ok {
		panic(fmt.Errorf("checker with same name already exists: %s", name))
	}

	checkers[name] = fn
}

// Result is the checker result returned after a check exec has been
// completed.
//
// Result also implements the metrics.Metric interface so it can be exported.
type Result struct {
	Successful bool
	Timeout    bool
	StartTime  time.Time
	Duration   time.Duration
}

// Checker is something that probes the provided target and marks it's run
// as successful or not. This depends on the fact that whether the output
// received matches the desired output provided by the user or not.
type Checker interface {
	// Validate validates the check configuration.
	Validate(Check) error

	// Provision is run to set the checker fields after the check
	// configuration is validated.
	Provision(Check) error

	// Execute executes the check and returns the result.
	Execute(context.Context) (*Result, error)
}

// getCheckerInstance returns a new checker from the check config.
func getCheckerInstance(check Check) (string, Checker, error) {
	// from the check config, we get the check from the config's input type,
	// i.e., the input type is the checker's name.
	name := check.GetInput().GetType()
	newFunc, ok := checkers[name]
	if !ok {
		return "", nil, fmt.Errorf("checker with name not registered: %s", name)
	}

	return name, newFunc(), nil
}

// newControllerFunc creates, validates and provisions the checker and
// returns the controller function.
func newControllerFunc(check Check) (controller.RunnerFunc, error) {
	name, checker, err := getCheckerInstance(check)
	if err != nil {
		return nil, err
	}

	if err := checker.Validate(check); err != nil {
		return nil, fmt.Errorf("%s checker: %w: %v", name, ErrValidation, err)
	}

	if err := checker.Provision(check); err != nil {
		return nil, fmt.Errorf("%s checker: %w: %v", name, ErrProvisioning, err)
	}

	return func(ctx context.Context) (interface{}, error) {
		return checker.Execute(ctx)
	}, nil
}

// NewControllerOpts creates controller options for the check.
func NewControllerOpts(check Check) (*controller.Opts, error) {
	interval := check.GetInterval()
	if interval <= 0 {
		return nil, fmt.Errorf("interval should be > 0")
	}

	id := check.GetID()
	name := check.GetName()

	fn, err := newControllerFunc(check)
	if err != nil {
		return nil, err
	}

	return &controller.Opts{
		ID:       id,
		Name:     name,
		Interval: interval,
		Func:     fn,
	}, nil
}

// Validate validates the check config. This method is meant to be used only
// when needed to validate if the check conf is correct.
func Validate(check Check) error {
	_, checker, err := getCheckerInstance(check)
	if err != nil {
		return err
	}

	return checker.Validate(check)
}
