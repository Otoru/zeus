package zeus

import "fmt"

// NotAFunctionError indicates that the provided object is not a function.
type NotAFunctionError struct{}

// Error returns a string representation of the NotAFunctionError.
func (e NotAFunctionError) Error() string {
	return "provided object is not a function"
}

// InvalidFactoryReturnError indicates that the factory function has an invalid number of return values.
type InvalidFactoryReturnError struct {
	NumReturns int
}

// Error returns a string representation of the InvalidFactoryReturnError.
func (e InvalidFactoryReturnError) Error() string {
	return fmt.Sprintf("factory must return 1 or 2 values, got %d", e.NumReturns)
}

// UnexpectedReturnTypeError indicates that the return type of the factory function is unexpected.
type UnexpectedReturnTypeError struct {
	TypeName string
}

// Error returns a string representation of the UnexpectedReturnTypeError.
func (e UnexpectedReturnTypeError) Error() string {
	return fmt.Sprintf("unexpected return type: %s", e.TypeName)
}

// FactoryAlreadyProvidedError indicates that a factory for the given type has already been registered.
type FactoryAlreadyProvidedError struct {
	TypeName string
}

// Error returns a string representation of the FactoryAlreadyProvidedError.
func (e FactoryAlreadyProvidedError) Error() string {
	return fmt.Sprintf("a factory for type %s has already been provided", e.TypeName)
}

// DependencyResolutionError indicates that a dependency could not be resolved.
type DependencyResolutionError struct {
	TypeName string
}

// Error returns a string representation of the DependencyResolutionError.
func (e DependencyResolutionError) Error() string {
	return fmt.Sprintf("failed to resolve dependency for type %s", e.TypeName)
}

// CyclicDependencyError indicates that a cyclic dependency was detected.
type CyclicDependencyError struct {
	TypeName string
}

// Error returns a string representation of the CyclicDependencyError.
func (e CyclicDependencyError) Error() string {
	return fmt.Sprintf("cyclic dependency detected for type %s", e.TypeName)
}
