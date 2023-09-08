package zeus

import (
	"reflect"

	"golang.org/x/exp/slices"
)

// Container holds the registered factories for dependency resolution.
type Container struct {
	providers map[reflect.Type]reflect.Value
	instances map[reflect.Type]reflect.Value
}

// New initializes and returns a new instance of the Container.
//
// Example:
//
//	c := zeus.New()
func New() *Container {
	providers := make(map[reflect.Type]reflect.Value, 0)
	instances := make(map[reflect.Type]reflect.Value, 0)

	container := new(Container)
	container.providers = providers
	container.instances = instances

	return container
}

// Provide registers a factory function for dependency resolution.
// It ensures that the factory is a function, has a valid return type, and checks for duplicate factories.
// Returns an error if any of these conditions are not met.
//
// Example:
//
//	c := zeus.New()
//	c.Provide(func() int { return 42 })
func (c *Container) Provide(factory interface{}) error {
	factoryType := reflect.TypeOf(factory)

	if factoryType.Kind() != reflect.Func {
		return NotAFunctionError{}
	}

	if numOut := factoryType.NumOut(); numOut < 1 || numOut > 2 {
		return InvalidFactoryReturnError{NumReturns: numOut}
	}

	if factoryType.NumOut() == 2 && factoryType.Out(1).Name() != "error" {
		return UnexpectedReturnTypeError{TypeName: factoryType.Out(1).Name()}
	}

	serviceType := factoryType.Out(0)

	if _, exists := c.providers[serviceType]; exists {
		return FactoryAlreadyProvidedError{TypeName: serviceType.Name()}
	}

	c.providers[serviceType] = reflect.ValueOf(factory)

	return nil
}

// resolve attempts to resolve a dependency of the given type.
// It checks for cyclic dependencies and ensures that all dependencies can be resolved.
// Returns the resolved value and any error encountered during resolution.
func (c *Container) resolve(t reflect.Type, stack []reflect.Type) (reflect.Value, error) {
	if slices.Contains(stack, t) {
		return reflect.Value{}, CyclicDependencyError{TypeName: t.Name()}
	}

	if instance, exists := c.instances[t]; exists {
		return instance, nil
	}

	provider, ok := c.providers[t]

	if !ok {
		return reflect.Value{}, DependencyResolutionError{TypeName: t.Name()}
	}

	providerType := provider.Type()
	dependencies := make([]reflect.Value, providerType.NumIn())

	for i := range dependencies {
		argType := providerType.In(i)
		argValue, err := c.resolve(argType, append(stack, t))

		if err != nil {
			return reflect.Value{}, err
		}

		dependencies[i] = argValue
	}

	results := provider.Call(dependencies)

	if len(results) == 2 && !results[1].IsNil() {
		return reflect.Value{}, results[1].Interface().(error)
	}

	c.instances[t] = results[0]

	return results[0], nil
}

// Run executes the provided function by resolving and injecting its dependencies.
// It ensures that the function has a valid signature and that all dependencies can be resolved.
// Returns an error if the function signature is invalid or if dependencies cannot be resolved.
//
// Example:
//
//	c := zeus.New()
//	c.Provide(func() int { return 42 })
//	c.Run(func(i int) {
//	    fmt.Println(i) // Outputs: 42
//	})
func (c *Container) Run(fn interface{}) error {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return NotAFunctionError{}
	}

	if numOut := fnType.NumOut(); numOut > 1 {
		return InvalidFactoryReturnError{NumReturns: numOut}
	}

	if fnType.NumOut() == 1 && fnType.Out(0).Name() != "error" {
		return UnexpectedReturnTypeError{TypeName: fnType.Out(0).Name()}
	}

	dependencies := make([]reflect.Value, fnType.NumIn())

	for i := range dependencies {
		argType := fnType.In(i)
		argValue, err := c.resolve(argType, nil)

		if err != nil {
			return err
		}

		dependencies[i] = argValue
	}

	results := reflect.ValueOf(fn).Call(dependencies)

	if fnType.NumOut() == 1 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}
