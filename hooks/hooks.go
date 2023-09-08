package hooks

import (
	"sync"

	"github.com/otoru/zeus/errs"
)

// Hooks defines an interface for lifecycle events.
// It provides methods to register functions that should be executed
// at the start and stop of the application.
type Hooks interface {
	OnStart(func() error)
	OnStop(func() error)
	Start() error
	Stop() error
}

// HooksImpl is the default implementation of the Hooks interface.
type HooksImpl struct {
	onStart []func() error
	onStop  []func() error
}

// OnStart adds a function to the list of functions to be executed at the start.
// Example:
//
//	hooks.OnStart(func() error {
//	   fmt.Println("Starting...")
//	   return nil
//	})
func (h *HooksImpl) OnStart(fn func() error) {
	h.onStart = append(h.onStart, fn)
}

// OnStop adds a function to the list of functions to be executed at the stop.
// Example:
//
//	hooks.OnStop(func() error {
//	   fmt.Println("Stopping...")
//	   return nil
//	})
func (h *HooksImpl) OnStop(fn func() error) {
	h.onStop = append(h.onStop, fn)
}

// Start executes all the registered OnStart hooks.
// It returns the first error encountered or nil if all hooks execute successfully.
// This method is internally used by the Container's Run function.
func (h *HooksImpl) Start() error {
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
func (h *HooksImpl) Stop() error {
	var wg sync.WaitGroup
	errorSet := &errs.ErrorSet{}

	for _, hook := range h.onStop {
		wg.Add(1)
		go func(hook func() error) {
			defer wg.Done()
			if err := hook(); err != nil {
				errorSet.Add(err)
			}
		}(hook)
	}

	wg.Wait()
	return errorSet.Result()
}
