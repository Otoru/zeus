package zeus

import "fmt"

type NotAFunctionError struct{}

func (e NotAFunctionError) Error() string {
	return "provided object is not a function"
}

type InvalidFactoryReturnError struct {
	NumReturns int
}

func (e InvalidFactoryReturnError) Error() string {
	return fmt.Sprintf("factory must return 1 or 2 values, got %d", e.NumReturns)
}

type UnexpectedReturnTypeError struct {
	TypeName string
}

func (e UnexpectedReturnTypeError) Error() string {
	return fmt.Sprintf("unexpected return type: %s", e.TypeName)
}

type FactoryAlreadyProvidedError struct {
	TypeName string
}

func (e FactoryAlreadyProvidedError) Error() string {
	return fmt.Sprintf("a factory for type %s has already been provided", e.TypeName)
}

type DependencyResolutionError struct {
	TypeName string
}

func (e DependencyResolutionError) Error() string {
	return fmt.Sprintf("failed to resolve dependency for type %s", e.TypeName)
}

type CyclicDependencyError struct {
	TypeName string
}

func (e CyclicDependencyError) Error() string {
	return fmt.Sprintf("cyclic dependency detected for type %s", e.TypeName)
}
