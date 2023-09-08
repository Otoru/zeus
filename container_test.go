package zeus

import (
	"fmt"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestContainer(t *testing.T) {
	t.Parallel()

	t.Run("Provide", func(t *testing.T) {
		t.Run("Not a function", func(t *testing.T) {
			c := New()
			got := c.Provide("string")
			expected := NotAFunctionError{}
			assert.ErrorIs(t, got, expected)
		})

		t.Run("Invalid return count", func(t *testing.T) {
			c := New()
			got := c.Provide(func() (int, string, error) { return 0, "", nil })
			expected := InvalidFactoryReturnError{NumReturns: 3}

			assert.ErrorIs(t, got, expected)
		})

		t.Run("Second return value not is a error", func(t *testing.T) {
			c := New()
			got := c.Provide(func() (int, string) { return 0, "" })
			expected := UnexpectedReturnTypeError{TypeName: "string"}

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
			expected := FactoryAlreadyProvidedError{TypeName: "int"}

			assert.ErrorIs(t, got, expected)
		})
	})

	t.Run("resolve", func(t *testing.T) {
		t.Run("Cyclic dependency", func(t *testing.T) {
			c := New()
			c.Provide(func(s string) string { return s })
			_, err := c.resolve(reflect.TypeOf(""), []reflect.Type{reflect.TypeOf("")})
			expected := CyclicDependencyError{TypeName: "string"}

			assert.ErrorIs(t, err, expected)
		})

		t.Run("Unresolved dependency", func(t *testing.T) {
			c := New()
			_, err := c.resolve(reflect.TypeOf(0.0), nil)
			expected := DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, err, expected)
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
			_, err := c.resolve(reflect.TypeOf(0), nil)
			expected := DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, err, expected)
		})

		t.Run("Recursive call error - cyclic dependency", func(t *testing.T) {
			c := New()
			c.Provide(func(s string) int { return len(s) })
			c.Provide(func(i int) string { return fmt.Sprint(i) })
			_, err := c.resolve(reflect.TypeOf(0), nil)
			expected := CyclicDependencyError{TypeName: "int"}

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

	})

	t.Run("Run", func(t *testing.T) {
		t.Run("Not a function", func(t *testing.T) {
			c := New()
			err := c.Run("not a function")
			expected := NotAFunctionError{}

			assert.ErrorIs(t, err, expected)
		})

		t.Run("Invalid return", func(t *testing.T) {
			c := New()
			err := c.Run(func() (int, string) { return 0, "" })
			expected := InvalidFactoryReturnError{NumReturns: 2}

			assert.ErrorIs(t, err, expected)
		})

		t.Run("Function returns a non-error value", func(t *testing.T) {
			c := New()
			err := c.Run(func() int { return 42 })
			expected := UnexpectedReturnTypeError{TypeName: "int"}

			assert.ErrorIs(t, err, expected)
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
			err := c.Run(func(f float64) error { return nil })
			expected := DependencyResolutionError{TypeName: "float64"}

			assert.ErrorIs(t, err, expected)
		})

	})
}
