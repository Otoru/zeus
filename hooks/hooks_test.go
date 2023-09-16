package hooks

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHooksImpl(t *testing.T) {
	t.Run("OnStart", func(t *testing.T) {
		h := &LifecycleHooks{}

		t.Run("should add function to onStart slice", func(t *testing.T) {
			h.OnStart(func() error {
				return nil
			})
			assert.Equal(t, len(h.onStart), 1)
		})
	})

	t.Run("OnStop", func(t *testing.T) {
		h := &LifecycleHooks{}

		t.Run("should add function to onStop slice", func(t *testing.T) {
			h.OnStop(func() error {
				return nil
			})
			assert.Equal(t, len(h.onStop), 1)
		})
	})

	t.Run("Start", func(t *testing.T) {
		t.Run("should execute all onStart hooks without error", func(t *testing.T) {
			h := &LifecycleHooks{}
			h.OnStart(func() error {
				return nil
			})
			h.OnStart(func() error {
				return nil
			})
			err := h.Start()
			assert.NilError(t, err)
		})

		t.Run("should return error if any onStart hook fails", func(t *testing.T) {
			h := &LifecycleHooks{}
			h.OnStart(func() error {
				return nil
			})
			h.OnStart(func() error {
				return errors.New("start error")
			})
			err := h.Start()
			assert.ErrorContains(t, err, "start error")
		})
	})

	t.Run("Stop", func(t *testing.T) {
		t.Run("should execute all onStop hooks without error", func(t *testing.T) {
			h := &LifecycleHooks{}
			h.OnStop(func() error {
				return nil
			})
			h.OnStop(func() error {
				return nil
			})
			err := h.Stop()
			assert.NilError(t, err)
		})

		t.Run("should return error if any onStop hook fails", func(t *testing.T) {
			h := &LifecycleHooks{}
			h.OnStop(func() error {
				return nil
			})
			h.OnStop(func() error {
				return errors.New("stop error")
			})
			err := h.Stop()
			assert.ErrorContains(t, err, "stop error")
		})
	})
}
