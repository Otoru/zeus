package hooks

import "sync"

// Hooks defines an interface for lifecycle events.
// It provides methods to register functions that should be executed
// at the start and stop of the application.
type Hooks interface {
	OnStart(func() error)
	OnStop(func() error)
	Start() error
	Stop() error
}

// LifecycleHooks is the default implementation of the Hooks interface.
type LifecycleHooks struct {
	onStart []func() error
	onStop  []func() error
	mu      sync.Mutex
}

// OnStart adds a function to the list of functions to be executed at the start.
// Example:
//
//	hooks.OnStart(func() error {
//	   fmt.Println("Starting...")
//	   return nil
//	})
func (h *LifecycleHooks) OnStart(fn func() error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onStart = append(h.onStart, fn)
}

// OnStop adds a function to the list of functions to be executed at the stop.
// Example:
//
//	hooks.OnStop(func() error {
//	   fmt.Println("Stopping...")
//	   return nil
//	})
func (h *LifecycleHooks) OnStop(fn func() error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onStop = append(h.onStop, fn)
}

// Start executes all the registered OnStart hooks.
// It returns the first error encountered or nil if all hooks execute successfully.
// This method is internally used by the Container's Run function.
func (h *LifecycleHooks) Start() error {
	for _, hook := range h.onStart {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}

// Stop executes all the registered OnStop hooks.
// It returns the first error encountered or nil if all hooks execute successfully.
// This method is internally used by the Container's Run function.
func (h *LifecycleHooks) Stop() error {
	for _, hook := range h.onStop {
		if err := hook(); err != nil {
			return err
		}
	}

	return nil
}
