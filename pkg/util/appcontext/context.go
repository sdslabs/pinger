package appcontext

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

// Context is the application context that carries a context.Context and
// some other utilities to be used across the application.
//
// Implements context.Context.
type Context struct {
	ctx   context.Context
	log   *logrus.Logger
	debug bool
}

// Deadline returns the time when work done on behalf of this context should
// be canceled.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled.
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err returns the error from the context.
func (c *Context) Err() error {
	return c.ctx.Err()
}

// Value returns the value associated with this context for key, or nil if
// no value is associated with key.
func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Logger returns the logger from context.
func (c *Context) Logger() *logrus.Logger {
	return c.log
}

// Debug tells if the application is running in debugging mode or not. This
// is set by setting the `DEBUG` environment variable to "on".
func (c *Context) Debug() bool {
	return c.debug
}

// Background returns an empty context with default logrus logger.
func Background() *Context {
	log := logrus.New()
	ctx := context.Background()
	return &Context{ctx: ctx, log: log}
}

// BackgroundDebug is same as Background with Debug returning true.
func BackgroundDebug() *Context {
	ctx := Background()
	ctx.debug = true
	return ctx
}

// WithCancel creates a Context from context.Context with a cancel function.
func WithCancel(parent *Context) (ctx *Context, cancel func()) {
	// check if the parent is already a *Context
	if parent == nil {
		return nil, nil
	}

	child, cancelFunc := context.WithCancel(parent.ctx)
	return &Context{ctx: child, log: parent.log, debug: parent.debug}, cancelFunc
}

// WithSignals creates a context that cancels on receiving the os.Signal.
func WithSignals(parent *Context, signals ...os.Signal) (ctx *Context, cancel func()) {
	ctx, cancel = WithCancel(parent)

	stream := make(chan os.Signal, 1)
	signal.Notify(stream, signals...)
	go func(recv <-chan os.Signal, c context.Context, cancelFunc func()) {
		select {
		case <-recv:
			cancelFunc()
		case <-c.Done():
			return
		}

		// If the signal is received again we want to exit the program and not just
		// wait for cancel at this point. There is no need to exit this go routine
		// since this is meant to be used in case of cmd exit and the thread will
		// be exited in case of successful cancel and shutdown.
		<-recv
		os.Exit(1)
	}(stream, ctx, cancel)

	return
}

// Interface guard.
var _ context.Context = (*Context)(nil)
