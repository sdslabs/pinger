package controller

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RunnerFunc is the function that runs after each tick.
type RunnerFunc func(context.Context) (interface{}, error)

// RunStat are the statistics for single run of the runner function.
type RunStat struct {
	ID   string
	Name string

	Err error
	Res interface{}
}

// Opts are the options required to create a new controller.
type Opts struct {
	ID       string
	Name     string
	Interval time.Duration
	Func     RunnerFunc
}

// Controller runs a specific operation infinitely until the context is
// canceled at regular intervals of time.
type Controller struct {
	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.RWMutex
	wg    sync.WaitGroup

	stats     map[time.Time]*RunStat
	latestRun time.Time

	interval time.Duration
	update   chan struct{}

	fn RunnerFunc

	id   string
	name string
}

// NewController creates a new controller with given name and type and runs
// the function at regular `interval`s of time without blocking.
func NewController(ctx context.Context, opts *Opts) (*Controller, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if err := validateControllerOpts(opts); err != nil {
		return nil, err
	}

	ctxt, cancel := context.WithCancel(ctx)

	return &Controller{
		ctx:    ctxt,
		cancel: cancel,

		mutex: sync.RWMutex{},
		wg:    sync.WaitGroup{},

		stats:     make(map[time.Time]*RunStat),
		latestRun: time.Time{},

		interval: opts.Interval,
		update:   make(chan struct{}, 1),

		fn: opts.Func,

		id:   opts.ID,
		name: opts.Name,
	}, nil
}

// validateControllerOpts validates all the options required to create a new
// controller.
func validateControllerOpts(opts *Opts) error {
	if opts.Interval <= 0 {
		return fmt.Errorf("controller interval should be > 0")
	}

	if opts.Func == nil {
		return fmt.Errorf("controller function cannot be nil")
	}

	return nil
}

// ID returns the ID of the controller.
func (c *Controller) ID() string {
	return c.id
}

// Name returns the name of the controller.
func (c *Controller) Name() string {
	return c.name
}

// Interval returns the interval after which the controller executes the
// function.
func (c *Controller) Interval() time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.interval
}

// Start starts the execution of the controller.
func (c *Controller) Start() {
	c.wg.Add(1)
	go func(ctrl *Controller) {
		defer ctrl.wg.Done()
		for {
			runCtx, runCancel := context.WithCancel(ctrl.ctx)
			defer runCancel()

			ctrl.run(runCtx)

			select {
			case <-ctrl.update:
				runCancel()
				continue

			case <-ctrl.ctx.Done():
				return
			}
		}
	}(c)
}

// run starts the ticker and executes the function with every tick.
func (c *Controller) run(runCtx context.Context) {
	c.wg.Add(1)
	go func(ctx context.Context, ctrl *Controller) {
		defer ctrl.wg.Done()

		// run the function on start once
		ctrl.runFunc()

		c.mutex.RLock()
		interval := ctrl.interval
		c.mutex.RUnlock()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctrl.runFunc()

			case <-ctx.Done():
				return
			}
		}
	}(runCtx, c)
}

// runFunc executes the runner function.
func (c *Controller) runFunc() {
	c.wg.Add(1)
	go func(ctrl *Controller) {
		defer ctrl.wg.Done()

		ctrl.mutex.RLock()
		fn := ctrl.fn
		ctrl.mutex.RUnlock()

		res, err := fn(ctrl.ctx)
		stat := &RunStat{
			ID:   ctrl.id,
			Name: ctrl.name,

			Err: err,
			Res: res,
		}

		ctrl.mutex.Lock()
		tnow := time.Now()
		ctrl.latestRun = tnow
		ctrl.stats[tnow] = stat
		ctrl.mutex.Unlock()
	}(c)
}

// Wait waits for execution of the controller to be complete.
func (c *Controller) Wait() {
	c.wg.Wait()
}

// UpdateInterval updates the interval of the controller.
func (c *Controller) UpdateInterval(interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("controller interval should be > 0")
	}

	c.mutex.Lock()
	c.interval = interval
	c.mutex.Unlock()

	c.update <- struct{}{}
	return nil
}

// UpdateFunc updates the controller function.
func (c *Controller) UpdateFunc(fn RunnerFunc) error {
	if fn == nil {
		return fmt.Errorf("controller function cannot be nil")
	}

	c.mutex.Lock()
	c.fn = fn
	c.mutex.Unlock()

	// no need to send the update signal since for each run the function is read
	// again and is protected with read lock.
	return nil
}

// Stop stops the execution of the controller. It does not wait for
// controller to complete. Use `Wait` for that.
func (c *Controller) Stop() {
	c.cancel()
}

// PullAllStats fetches the stats for the controller and also cleans up the
// stats.
func (c *Controller) PullAllStats() map[time.Time]*RunStat {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stats := c.stats
	c.stats = make(map[time.Time]*RunStat)
	return stats
}

// PullLatestStat fetches only the latest stat from the controller and
// cleans up the stats.
func (c *Controller) PullLatestStat() *RunStat {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stat, ok := c.stats[c.latestRun]
	if !ok {
		return nil
	}

	c.stats = make(map[time.Time]*RunStat)
	return stat
}
