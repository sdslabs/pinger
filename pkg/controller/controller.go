package controller

import (
	"time"
	"sync"
	"context"
	"reflect"
	"fmt"

	"github.com/sirupsen/logrus"
)

// logging field definitions
const (
	// fieldControllerName is the name of the controller
	fieldControllerName = "name"

	// fieldConsecutiveErrors is the number of consecutive errors of a controller
	fieldConsecutiveErrors = "consecutiveErrors"
)

// An interface satisfied by all the functions.
type Function interface{}

// An interface staisfied by all the function parameter types.
type FuncParam interface{}

// ControllerFunction represents an underyling function executed by the controller.
// A controller function is uniquely identified by its name.
type ControllerFunction struct {
	name string

	Function Function
	Params   []FuncParam
}

func (cf *ControllerFunction) Validate() error {
	return validateNewControllerParams(cf.Function, cf.Params)
}

// Run the actual function that ControllerFunction actually represents.
// The context passed as an argument to the Run function is propogated to the
// underlying Function.
func (cf *ControllerFunction) Run(ctx context.Context) error {
	function := reflect.ValueOf(cf.Function)

	var params []reflect.Value

	params = append(params, reflect.ValueOf(ctx))
	for _, param := range cf.Params {
		params = append(params, reflect.ValueOf(param))
	}

	// Make the actual call to the underlying function.
	// We for sure know the return structure of the function
	// as we have already validated it before adding the ControllerFunction.
	// This will panic if the ControllerFunction is not validated.
	values := function.Call(params)

	return reflect.ValueOf(values[0]).Interface().(error)
}


// Returns an instance of a new controller function.
func NewControllerFunction(function Function, params ...FuncParam) (*ControllerFunction, error) {
	err := validateNewControllerParams(function, params...)
	if err != nil {
		return nil, err
	}
	funcType := reflect.TypeOf(function)

	cf := &ControllerFunction{
		name:     funcType.Name(),

		Function: function,
		Params:   params[1:],
	}
	
	return cf, nil

}

// Validate function and its parameters, this returns an error if the
// provided function and parameters are not valid.
//
// A provided function for the validation must have the following form
// func(ctx context.Context, ...params) error {}
//
// One inbound and one outbound variable is required with types context.Context
// and error respectivly.
func validateNewControllerParams(function Function, params ...FuncParam) error {
	funcType := reflect.TypeOf(function)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("Provided function is not a valid function")
	}

	if funcType.NumIn() < 1 {
		return fmt.Errorf("Provided function is not valid, must have atleast one argument, the context")
	}

	if funcType.NumOut() != 1 {
		return fmt.Errorf("Provided function is not valid, must have an error as return value.")
	}

	if funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("The return type of the function is not valid, must return an error value.")
	}

	if funcType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		return fmt.Errorf("The first parameter to the Controller function must be a context.")
	}

	// First parameter to the function is the context and is not provided
	// while registering the function.
	if funcType.NumIn() != len(params) + 1 {
		return fmt.Errorf("Parameters not valid required %d given %d", funcType.NumIn() - 1, len(params))
	}

	for i, param := range params {
		// There is an offset of 1 here as the context is not provided
		// while creating a new ControllerFunction but is rather provided 
		// when running the function.
		typ := funcType.In(i + 1)
		if typ != reflect.TypeOf(param) {
			return fmt.Errorf("Parameter type for param %s is not valid", typ.Name())
		}
	}

	return nil
}

// ControllerParams contains all parameters of a controller, including the functions to
// run and other metadata related to runs.
type ControllerInternal struct {
	// DoFunc is the function that will be run until it succeeds and/or
	// using the interval RunInterval if not 0.
	// An unset DoFunc is an error and will be logged as one.
	DoFunc *ControllerFunction

	// StopFunc is called when the controller stops. It is intended to run any
	// clean-up tasks for the controller (e.g. deallocate/release resources)
	// It is guaranteed that DoFunc is called at least once before StopFunc is
	// called.
	// An unset StopFunc is not an error (and will be a no-op)
	// Note: Since this occurs on controller exit, error counts and tracking may
	// not be checked after StopFunc is run.
	StopFunc *ControllerFunction

	// If set to any other value than 0, will cause DoFunc to be run in the
	// specified interval. The interval starts from when the DoFunc has
	// returned last
	RunInterval time.Duration

	// ErrorRetryBaseDuration is the initial time to wait to run DoFunc
	// again on return of an error. On each consecutive error, this value
	// is multiplied by the number of consecutive errors to provide a
	// constant back off. The default is 1s.
	ErrorRetryBaseDuration time.Duration

	// Should we have a constant back off for retries during errors.
	RetryBackOff bool

	// NoErrorRetry when set to true, disabled retries on errors
	NoErrorRetry bool
}

// Controller is the actual underlying controller. Each controller is created for a specific task
// which is specified in `controller.internal`
type Controller struct {
	// Mutex for the controller to hold locks.
	mutex             sync.RWMutex
	
	// Name of the controller, used by manager.
	name              string
	
	internal          ControllerInternal
	successCount      int
	failureCount      int

	lastSuccessStamp  time.Time
	lastErrorStamp    time.Time
	
	consecutiveErrors int
	lastError         error
	lastDuration      time.Duration

	stop              chan struct{}
	update            chan struct{}

	ctxDoFunc         context.Context
	cancelDoFunc      context.CancelFunc

	// terminated is closed after the controller has been terminated
	terminated chan struct{}
}

// Start running the controller.
// TODO: improve this, currently it waits for the current request to finish and then waits for 
// interval duration to run the function again. This is not a constant interval check we are
// looking for, so wrap the runFunc inside a goroutine.
func (c *Controller) RunController() {
	errorRetries := 1

	c.mutex.RLock()
	internal := c.internal
	c.mutex.RUnlock()
	runFunc := true

	// Default Run Interval for a controller.
	interval := 10 * time.Minute

	for {
		var err error
		if runFunc {
			interval = internal.RunInterval

			start := time.Now()
			
			// Run the function.
			err = internal.DoFunc.Run(c.ctxDoFunc)
			duration := time.Since(start)
			c.mutex.Lock()

			c.lastDuration = duration
			c.getLogger().Debug("Controller func execution time: ", c.lastDuration)

			if err != nil {
				c.getLogger().WithField(fieldConsecutiveErrors, errorRetries).
					WithError(err).Debug("Controller run failed")
				c.recordError(err)

				if !internal.NoErrorRetry && c.internal.RetryBackOff {
					if internal.ErrorRetryBaseDuration != time.Duration(0) {
						interval = time.Duration(errorRetries) * internal.ErrorRetryBaseDuration
					} else {
						interval = time.Duration(errorRetries) * time.Second
					}

					errorRetries++
				}
			} else {
				c.recordSuccess()

				// reset error retries after successful attempt
				errorRetries = 1

				// If no run interval is specified, no further updates
				// are required.
				if interval == time.Duration(0) {
					// Don't exit the goroutine, since that only happens when the
					// controller is explicitly stopped. Instead, just wait for
					// the next update.
					c.getLogger().Debug("Controller run succeeded; waiting for next controller update or stop")
					runFunc = false
					interval = 10 * time.Minute
				}
			}

			c.mutex.Unlock()
		}

		select {
		case <-c.stop:
			goto shutdown

		case <-c.update:
			// If we receive a signal on both channels c.stop and c.update,
			// golang will pick either c.stop or c.update randomly.
			// This select will make sure we don't execute the controller
			// while we are shutting down.
			select {
			case <-c.stop:
				goto shutdown
			default:
			}
			// Pick up any changes to the parameters in case the controller has
			// been updated.
			c.mutex.RLock()
			internal = c.internal
			c.mutex.RUnlock()
			runFunc = true

		case <-time.After(interval):
		}

	}

shutdown:
	c.getLogger().Debug("Shutting down controller")

	if err := internal.StopFunc.Run(context.TODO()); err != nil {
		c.mutex.Lock()
		c.recordError(err)
		c.mutex.Unlock()
		c.getLogger().WithField(fieldConsecutiveErrors, errorRetries).
			WithError(err).Warn("Error on Controller stop")
	}

	close(c.terminated)
}

// updateParamsLocked sets the specified controller's parameters.
//
// If the RunInterval exceeds ControllerMaxInterval, it will be capped.
func (c *Controller) updateController(internal ControllerInternal) {
	c.internal = internal

	c.update <- struct{}{}
}

func (c *Controller) stopController() {
	if c.cancelDoFunc != nil {
		c.cancelDoFunc()
	}

	close(c.stop)
	close(c.update)
}

// logger returns a logrus object with controllerName and UUID fields.
func (c *Controller) getLogger() *logrus.Entry {
	return c.getLogger().WithFields(logrus.Fields{
		fieldControllerName: c.name,
	})
}

// recordError updates all statistic collection variables on error
// c.mutex must be held.
func (c *Controller) recordError(err error) {
	c.lastError = err
	c.lastErrorStamp = time.Now()
	c.failureCount++
	c.consecutiveErrors++
}

// recordSuccess updates all statistic collection variables on success
// c.mutex must be held.
func (c *Controller) recordSuccess() {
	c.lastError = nil
	c.lastSuccessStamp = time.Now()
	c.successCount++
	c.consecutiveErrors = 0
}

type ControllerStatus struct {
	Name string
	Configuration *ControllerConfigurationStatus

	Status *ControllerRunStatus
}

type ControllerConfigurationStatus struct {
	ErrorRetry bool
	ShouldBackOff bool
	Interval string
}

type ControllerRunStatus struct {
	SuccessCount int64
	LastSuccessStamp string
	FailureCount int64
	LastFailureStamp string
	ConsecutiveFailureCount int64
}

// Returns the current status of the controller.
func (c *Controller) status() *ControllerStatus {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return &ControllerStatus{
		Name: c.name,

		Configuration: &ControllerConfigurationStatus{
			ErrorRetry: !c.internal.NoErrorRetry,
			ShouldBackOff: !c.internal.RetryBackOff,
			Interval: c.internal.RunInterval.String(),
		},
		Status: &ControllerRunStatus{
			SuccessCount: int64(c.successCount),
			FailureCount: int64(c.failureCount),
			ConsecutiveFailureCount: int64(c.consecutiveErrors),
			LastSuccessStamp: c.lastSuccessStamp.String(),
			LastFailureStamp: c.lastErrorStamp.String(),
		},
	}
}
