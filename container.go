package zeus

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type Container struct {
	providers map[reflect.Type]reflect.Value
}

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

func (c *Container) resolve(t reflect.Type, stack []reflect.Type) (reflect.Value, error) {
	if slices.Contains(stack, t) {
		return reflect.Value{}, CyclicDependencyError{TypeName: t.Name()}
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

	return results[0], nil
}

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

func New() *Container {
	providers := make(map[reflect.Type]reflect.Value, 0)

	container := new(Container)
	container.providers = providers

	return container
}
