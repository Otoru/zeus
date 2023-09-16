package zeus

import (
	"reflect"
	"slices"
	"sync"

	"github.com/otoru/zeus/errs"
	"github.com/otoru/zeus/hooks"
)

// Container holds the registered factories for dependency resolution.
type Container struct {
	providers map[reflect.Type]reflect.Value
	instances map[reflect.Type]reflect.Value
	mu        sync.RWMutex
	hooks     Hooks
}

// New initializes and returns a new instance of the Container.
//
// Example:
//
//	c := zeus.New()
func New() *Container {
	hooks := new(hooks.LifecycleHooks)
	providers := make(map[reflect.Type]reflect.Value)
	instances := make(map[reflect.Type]reflect.Value)

	container := new(Container)
	container.hooks = hooks
	container.providers = providers
	container.instances = instances

	return container
}

// resolve attempts to resolve a dependency of the given type.
// It checks for cyclic dependencies and ensures that all dependencies can be resolved.
// Returns the resolved value and any error encountered during resolution.
func (c *Container) resolve(t reflect.Type, stack []reflect.Type) (reflect.Value, error) {
	if slices.Contains(stack, t) {
		return reflect.Value{}, errs.CyclicDependencyError{TypeName: t.Name()}
	}

	c.mu.RLock()
	instance, hasInstance := c.instances[t]
	provider, hasProvider := c.providers[t]
	c.mu.RUnlock()

	if hasInstance {
		return instance, nil
	}

	if !hasProvider {
		return reflect.Value{}, errs.DependencyResolutionError{TypeName: t.Name()}
	}

	providerType := provider.Type()
	dependencies := make([]reflect.Value, providerType.NumIn())

	for i := range dependencies {
		argType := providerType.In(i)

		if argType.Implements(reflect.TypeOf((*Hooks)(nil)).Elem()) {
			dependencies[i] = reflect.ValueOf(c.hooks)
			continue
		}

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

// Provide registers a factory function for dependency resolution.
// It ensures that the factory is a function, has a valid return type, and checks for duplicate factories.
// Returns an error if any of these conditions are not met.
//
// Example:
//
//	c := zeus.New()
//	c.Provide(func() int { return 42 })
func (c *Container) Provide(factories ...interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, factory := range factories {
		factoryType := reflect.TypeOf(factory)

		if factoryType.Kind() != reflect.Func {
			return errs.NotAFunctionError{}
		}

		if numOut := factoryType.NumOut(); numOut < 1 || numOut > 2 {
			return errs.InvalidFactoryReturnError{NumReturns: numOut}
		}

		if factoryType.NumOut() == 2 {
			errorType := reflect.TypeOf((*error)(nil)).Elem()
			if !factoryType.Out(1).Implements(errorType) {
				return errs.UnexpectedReturnTypeError{TypeName: factoryType.Out(1).Name()}
			}
		}

		serviceType := factoryType.Out(0)

		if _, exists := c.providers[serviceType]; exists {
			return errs.FactoryAlreadyProvidedError{TypeName: serviceType.Name()}
		}

		c.providers[serviceType] = reflect.ValueOf(factory)
	}

	return nil
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
	errorSet := &errs.ErrorSet{}

	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return errs.NotAFunctionError{}
	}

	if numOut := fnType.NumOut(); numOut > 1 {
		return errs.InvalidFactoryReturnError{NumReturns: numOut}
	}

	if fnType.NumOut() == 1 && fnType.Out(0).Name() != "error" {
		return errs.UnexpectedReturnTypeError{TypeName: fnType.Out(0).Name()}
	}

	dependencies := make([]reflect.Value, fnType.NumIn())

	for i := range dependencies {
		argType := fnType.In(i)
		argValue, err := c.resolve(argType, nil)

		if err != nil {
			errorSet.Add(err)
			break
		}

		dependencies[i] = argValue
	}

	if !errorSet.IsEmpty() {
		return errorSet.Result()
	}

	if err := c.hooks.Start(); err != nil {
		errorSet.Add(err)
	}

	if !errorSet.IsEmpty() {
		return errorSet.Result()
	}

	results := reflect.ValueOf(fn).Call(dependencies)

	if fnType.NumOut() == 1 && !results[0].IsNil() {
		errorSet.Add(results[0].Interface().(error))
	}

	if err := c.hooks.Stop(); err != nil {
		errorSet.Add(err)
	}

	return errorSet.Result()
}

// Merge combines the factories of another container into the current container.
// If a factory from the other container conflicts with an existing factory in the current container,
// and they are not identical, a FactoryAlreadyProvidedError is returned.
//
// Example:
//
//	containerA := New()
//	containerB := New()
//
//	containerA.Provide(func() string { return "Hello" })
//	containerB.Provide(func() int { return 42 })
//
//	err := containerA.Merge(containerB)
//	if err != nil {
//	    // Handle merge error
//	}
func (c *Container) Merge(other *Container) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for t, factory := range other.providers {
		if existingFactory, exists := c.providers[t]; exists {
			if existingFactory.Pointer() != factory.Pointer() {
				return errs.FactoryAlreadyProvidedError{TypeName: t.Name()}
			}
			continue
		}

		c.providers[t] = factory
	}
	return nil
}
