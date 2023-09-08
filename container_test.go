package zeus

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/otoru/zeus/errs"
	"gotest.tools/v3/assert"
)

func TestContainer(t *testing.T) {
	t.Parallel()

	t.Run("resolve", func(t *testing.T) {
		t.Parallel()

		t.Run("Cyclic dependency", func(t *testing.T) {
			c := New()
			c.Provide(func(s string) string { return s })
			_, got := c.resolve(reflect.TypeOf(""), []reflect.Type{reflect.TypeOf("")})
			expected := errs.CyclicDependencyError{TypeName: "string"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Unresolved dependency", func(t *testing.T) {
			c := New()
			_, got := c.resolve(reflect.TypeOf(0.0), nil)
			expected := errs.DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Successful resolution", func(t *testing.T) {
			c := New()
			c.Provide(func() int { return 42 })
			c.Provide(func(i int) string { return "Hello" })
			val, err := c.resolve(reflect.TypeOf(""), nil)

			assert.NilError(t, err)
			assert.Equal(t, val.String(), "Hello")
		})

		t.Run("Recursive Call Error - Unresolved Dependency", func(t *testing.T) {
			c := New()
			c.Provide(func(f float64) int { return int(f) })
			_, got := c.resolve(reflect.TypeOf(0), nil)
			expected := errs.DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Recursive call error - cyclic dependency", func(t *testing.T) {
			c := New()
			c.Provide(func(s string) int { return len(s) })
			c.Provide(func(i int) string { return fmt.Sprint(i) })
			_, err := c.resolve(reflect.TypeOf(0), nil)
			expected := errs.CyclicDependencyError{TypeName: "int"}

			assert.ErrorIs(t, err, expected)
		})

		t.Run("Factory returns a error", func(t *testing.T) {
			c := New()
			c.Provide(func() (int, error) { return 0, fmt.Errorf("some error") })
			_, err := c.resolve(reflect.TypeOf(0), nil)

			assert.ErrorContains(t, err, "some error")
		})

		t.Run("Shared Instance Between Dependencies", func(t *testing.T) {
			c := New()

			type ServiceC struct{}
			type ServiceB struct {
				C *ServiceC
			}
			type ServiceA struct {
				C *ServiceC
			}

			c.Provide(func() *ServiceC {
				return &ServiceC{}
			})
			c.Provide(func(c *ServiceC) *ServiceA {
				return &ServiceA{C: c}
			})
			c.Provide(func(c *ServiceC) *ServiceB {
				return &ServiceB{C: c}
			})

			aVal, _ := c.resolve(reflect.TypeOf(&ServiceA{}), nil)
			bVal, _ := c.resolve(reflect.TypeOf(&ServiceB{}), nil)

			a, ok := aVal.Interface().(*ServiceA)
			assert.Equal(t, ok, true)

			b, ok := bVal.Interface().(*ServiceB)
			assert.Equal(t, ok, true)

			assert.Equal(t, a.C, b.C)
		})

		t.Run("Hooks Injection", func(t *testing.T) {
			c := New()

			err := c.Provide(func(h Hooks) *strings.Builder {
				h.OnStart(func() error {
					return nil
				})

				return &strings.Builder{}
			})

			assert.NilError(t, err)

			val, err := c.resolve(reflect.TypeOf(&strings.Builder{}), nil)
			assert.NilError(t, err)

			_, ok := val.Interface().(*strings.Builder)
			assert.Assert(t, ok)
		})

	})

	t.Run("Provide", func(t *testing.T) {
		t.Parallel()

		t.Run("Not a function", func(t *testing.T) {
			c := New()
			got := c.Provide("string")
			expected := errs.NotAFunctionError{}
			assert.ErrorIs(t, got, expected)
		})

		t.Run("Invalid return count", func(t *testing.T) {
			c := New()
			got := c.Provide(func() (int, string, error) { return 0, "", nil })
			expected := errs.InvalidFactoryReturnError{NumReturns: 3}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Second return value not is a error", func(t *testing.T) {
			c := New()
			got := c.Provide(func() (int, string) { return 0, "" })
			expected := errs.UnexpectedReturnTypeError{TypeName: "string"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Valid factory", func(t *testing.T) {
			c := New()
			err := c.Provide(func() int { return 0 })

			assert.NilError(t, err)
		})

		t.Run("Duplicated factory", func(t *testing.T) {
			c := New()
			c.Provide(func() int { return 0 })
			got := c.Provide(func() int { return 1 })
			expected := errs.FactoryAlreadyProvidedError{TypeName: "int"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Hooks Injection", func(t *testing.T) {
			c := New()

			err := c.Provide(func(h Hooks) *strings.Builder {
				return &strings.Builder{}
			})
			assert.NilError(t, err)

			_, exists := c.providers[reflect.TypeOf(&strings.Builder{})]
			assert.Assert(t, exists)
		})
	})

	t.Run("Run", func(t *testing.T) {
		t.Parallel()

		t.Run("Not a function", func(t *testing.T) {
			c := New()
			got := c.Run("not a function")
			expected := errs.NotAFunctionError{}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Invalid return", func(t *testing.T) {
			c := New()
			got := c.Run(func() (int, string) { return 0, "" })
			expected := errs.InvalidFactoryReturnError{NumReturns: 2}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Function returns a non-error value", func(t *testing.T) {
			c := New()
			got := c.Run(func() int { return 42 })
			expected := errs.UnexpectedReturnTypeError{TypeName: "int"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Successful execution", func(t *testing.T) {
			c := New()
			c.Provide(func() int { return 42 })

			err := c.Run(func(i int) error {
				if i != 42 {
					return fmt.Errorf("expected 42, got %d", i)
				}
				return nil
			})

			assert.NilError(t, err)
		})

		t.Run("Function returns a error", func(t *testing.T) {
			c := New()
			c.Provide(func() int { return 42 })
			err := c.Run(func(i int) error {
				return fmt.Errorf("some error")
			})

			assert.ErrorContains(t, err, "some error")
		})

		t.Run("Dependency resolution error", func(t *testing.T) {
			c := New()
			got := c.Run(func(f float64) error { return nil })
			expected := errs.DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Successful Execution with Hooks", func(t *testing.T) {
			c := New()

			started := false
			stopped := false

			c.Provide(func(h Hooks) int {
				h.OnStart(func() error {
					started = true
					return nil
				})
				h.OnStop(func() error {
					stopped = true
					return nil
				})
				return 42
			})

			err := c.Run(func(number int) {})

			assert.NilError(t, err)
			assert.Assert(t, started)
			assert.Assert(t, stopped)
		})

		t.Run("Error in OnStart Hook", func(t *testing.T) {
			c := New()

			c.Provide(func(h Hooks) int {
				h.OnStart(func() error {
					return errors.New("start error")
				})
				return 42
			})

			err := c.Run(func(number int) {})
			assert.ErrorContains(t, err, "start error")
		})

		t.Run("Error in OnStop Hook", func(t *testing.T) {
			c := New()

			c.Provide(func(h Hooks) int {
				h.OnStop(func() error {
					return errors.New("stop error")
				})

				return 42
			})

			err := c.Run(func(number int) {})
			assert.ErrorContains(t, err, "stop error")
		})
	})

	t.Run("Merge", func(t *testing.T) {
		t.Run("Merge without conflicts", func(t *testing.T) {
			containerA := New()
			containerB := New()

			containerA.Provide(func() string { return "Hello" })
			containerB.Provide(func() int { return 42 })

			err := containerA.Merge(containerB)
			assert.NilError(t, err)
		})

		t.Run("Merge with identical factories", func(t *testing.T) {
			containerA := New()
			containerB := New()

			factory := func() string { return "Hello" }

			containerA.Provide(factory)
			containerB.Provide(factory)

			err := containerA.Merge(containerB)
			assert.NilError(t, err)
		})

		t.Run("Merge with conflicting factories", func(t *testing.T) {
			containerA := New()
			containerB := New()

			containerA.Provide(func() string { return "Hello" })
			containerB.Provide(func() string { return "World" })

			err := containerA.Merge(containerB)
			assert.ErrorIs(t, err, errs.FactoryAlreadyProvidedError{TypeName: "string"})
		})
	})
}
