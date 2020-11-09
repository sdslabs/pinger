package controller

import (
	"context"
	"fmt"
	"sync"
)

// Manager manages multiple controllers and running at the same time.
type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.RWMutex

	controllers map[string]*Controller
}

// NewManager creates a new manager with no controllers.
func NewManager(ctx context.Context) *Manager {
	ctxt, cancel := context.WithCancel(ctx)

	return &Manager{
		ctx:    ctxt,
		cancel: cancel,

		mutex: sync.RWMutex{},

		controllers: make(map[string]*Controller),
	}
}

// UpdateController updates the controller if it exists with the same name
// or creates a new controller if it doesn't.
func (m *Manager) UpdateController(opts *Opts) error {
	if opts.Name == "" {
		return fmt.Errorf("cannot add a controller with empty name")
	}

	m.mutex.RLock()
	ctrl, ok := m.controllers[opts.ID]
	m.mutex.RUnlock()

	if ok {
		if opts.Interval > 0 && opts.Interval != ctrl.interval {
			if err := ctrl.UpdateInterval(opts.Interval); err != nil {
				return err
			}
		}

		if opts.Func != nil {
			if err := ctrl.UpdateFunc(opts.Func); err != nil {
				return err
			}
		}

		return nil
	}

	ctrl, err := NewController(m.ctx, opts)
	if err != nil {
		return err
	}

	ctrl.Start()

	m.mutex.Lock()
	m.controllers[opts.ID] = ctrl
	m.mutex.Unlock()

	return nil
}

// ListControllers lists all the controllers managed by the manager.
//
// This returns a map of controller ID with it's name.
func (m *Manager) ListControllers() map[string]string {
	list := map[string]string{}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, ctrl := range m.controllers {
		list[ctrl.ID()] = ctrl.Name()
	}

	return list
}

// RemoveController removes the controller from manager if it exists. If it
// doesn't, it does nothing. It does not wait for controller to stop.
func (m *Manager) RemoveController(id string) {
	ctrl, ok := m.removeCtrl(id)
	if !ok {
		return
	}

	ctrl.Stop()
}

// RemoveControllerAndWait is like RemoveController but waits for completion
// of controller.
func (m *Manager) RemoveControllerAndWait(id string) {
	ctrl, ok := m.removeCtrl(id)
	if !ok {
		return
	}

	ctrl.Stop()
	ctrl.Wait()
}

// removeCtrl removes the controller from the map. It returns the removed
// controller and whether it existed in the map or not.
func (m *Manager) removeCtrl(id string) (*Controller, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ctrl, ok := m.controllers[id]
	if !ok {
		return nil, false
	}
	delete(m.controllers, id)

	return ctrl, true
}

// RemoveAll removes all the controllers from the manager.
func (m *Manager) RemoveAll() {
	m.removeAll()
}

// RemoveAllAndWait removes all the controllers and waits
// for them to stop completely.
func (m *Manager) RemoveAllAndWait() {
	ctrls := m.removeAll()
	for _, c := range ctrls {
		c.Wait()
	}
}

// removeAll stops and removes all the controllers from the manager but does
// not wait for them to terminate.
func (m *Manager) removeAll() []*Controller {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ctrls := []*Controller{}

	for name, c := range m.controllers {
		c.Stop()
		ctrls = append(ctrls, c)
		delete(m.controllers, name)
	}

	return ctrls
}

// Wait waits for the manager to terminate.
func (m *Manager) Wait() {
	<-m.ctx.Done()
}

// Shutdown gracefully attempts to close the manager and terminate and
// remove all the controllers. In case of error, the manager is not closed,
// and hence `Close` is required to be called to terminate the manager.
func (m *Manager) Shutdown(ctx context.Context) error {
	waitChan := make(chan struct{})

	go func(w chan<- struct{}) {
		m.RemoveAllAndWait()
		w <- struct{}{}
	}(waitChan)

	select {
	case <-waitChan:
		m.cancel()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close closes all the controller irrespective of whether they shutdown
// properly.
func (m *Manager) Close() {
	m.RemoveAll()
	m.cancel()
}

// PullAllStats gets all the stats for all the controllers registered with
// the manager.
func (m *Manager) PullAllStats() map[string][]*RunStat {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string][]*RunStat)

	for id, ctrl := range m.controllers {
		cstats := []*RunStat{}

		pull := ctrl.PullAllStats()
		for _, p := range pull {
			cstats = append(cstats, p)
		}

		stats[id] = cstats
	}

	return stats
}

// PullLatestStats gets latest stat for each controller in the manager.
func (m *Manager) PullLatestStats() map[string]*RunStat {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]*RunStat)

	for id, ctrl := range m.controllers {
		stats[id] = ctrl.PullLatestStat()
	}

	return stats
}
