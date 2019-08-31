package controller

import (
	"sync"
	"time"
	"context"
	"fmt"
)

// A ControllerMap is the map of a controller name with the underlying Controller.
type ControllerMap map[string]*Controller

// Manager manages a ControllerMap and perform actions on it.
type Manager struct {
	controllers ControllerMap
	
	mutex sync.RWMutex
}

// Creates a new manager instance for the controller map.
func NewManager() Manager {
	return Manager{
		controllers: ControllerMap{},
	}
}

func NoopFunc(_ctx context.Context) error {
	return nil
}

// UpdateController installs or updates a controller in the manager. A
// controller is identified by its name. If a controller with the name already
// exists, the controller will be shut down and replaced with the provided
// controller. Updating a controller will cause the DoFunc to be run
// immediately regardless of any previous conditions. It will also cause any
// statistics to be reset.
func (m *Manager) UpdateController(name string, internal ControllerInternal) error {
	_, err := m.updateController(name, internal)

	return err
}

func (m *Manager) updateController(name string, internal ControllerInternal) (*Controller, error) {
	start := time.Now()

	// ensure the callbacks are valid
	if err := internal.DoFunc.Validate(); err != nil {
		return nil, err
	}
	if internal.StopFunc == nil {
		internal.StopFunc, _ = NewControllerFunction(NoopFunc)
	}

	m.mutex.Lock()

	if m.controllers == nil {
		m.controllers = ControllerMap{}
	}

	ctrl, exists := m.controllers[name]
	if exists {
		m.mutex.Unlock()

		ctrl.getLogger().Debug("Updating existing controller")
		ctrl.mutex.Lock()
		ctrl.updateController(internal)
		ctrl.mutex.Unlock()


		ctrl.getLogger().Debug("Controller update time: ", time.Since(start))
	} else {
		ctrl = &Controller{
			name:       name,
			stop:       make(chan struct{}),
			update:     make(chan struct{}, 1),
			terminated: make(chan struct{}),
		}
		ctrl.updateController(internal)
		ctrl.getLogger().Debug("Starting new controller")

		ctrl.ctxDoFunc, ctrl.cancelDoFunc = context.WithCancel(context.Background())
		m.controllers[ctrl.name] = ctrl
		m.mutex.Unlock()

		go ctrl.RunController()
	}

	return ctrl, nil
}

func (m *Manager) removeController(ctrl *Controller) {
	ctrl.stopController()
	delete(m.controllers, ctrl.name)

	ctrl.getLogger().Debug("Removed controller")
}

func (m *Manager) lookup(name string) *Controller {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if c, ok := m.controllers[name]; ok {
		return c
	}

	return nil
}

func (m *Manager) removeAndReturnController(name string) (*Controller, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.controllers == nil {
		return nil, fmt.Errorf("empty controller map")
	}

	oldCtrl, ok := m.controllers[name]
	if !ok {
		return nil, fmt.Errorf("unable to find controller %s", name)
	}

	m.removeController(oldCtrl)

	return oldCtrl, nil
}

// RemoveController stops and removes a controller from the manager. If DoFunc
// is currently running, DoFunc is allowed to complete in the background.
func (m *Manager) RemoveController(name string) error {
	_, err := m.removeAndReturnController(name)
	return err
}

// RemoveControllerAndWait stops and removes a controller using
// RemoveController() and then waits for it to run to completion.
func (m *Manager) RemoveControllerAndWait(name string) error {
	oldCtrl, err := m.removeAndReturnController(name)
	if err == nil {
		<-oldCtrl.terminated
	}

	return err
}

func (m *Manager) removeAll() []*Controller {
	ctrls := []*Controller{}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.controllers == nil {
		return ctrls
	}

	for _, ctrl := range m.controllers {
		m.removeController(ctrl)
		ctrls = append(ctrls, ctrl)
	}

	return ctrls
}

// RemoveAll stops and removes all controllers of the manager
func (m *Manager) RemoveAll() {
	m.removeAll()
}

// RemoveAllAndWait stops and removes all controllers of the manager and then
// waits for all controllers to exit
func (m *Manager) RemoveAllAndWait() {
	ctrls := m.removeAll()
	for _, ctrl := range ctrls {
		<-ctrl.terminated
	}
}
