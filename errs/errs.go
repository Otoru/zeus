package errs

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

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

// ErrorSet is a collection of errors.
// It can be used to accumulate errors and retrieve them as a single error or a list.
type ErrorSet struct {
	mu     sync.Mutex
	errors []error
}

// Add appends an error to the error set.
func (es *ErrorSet) Add(err error) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.errors = append(es.errors, err)
}

// Errors returns the list of errors in the error set.
func (es *ErrorSet) Errors() []error {
	slices.Reverse(es.errors)
	return es.errors
}

// Error implements the error interface.
// It returns a concatenated string of all error messages in the error set.
func (es *ErrorSet) Error() string {
	errMsgs := []string{}
	for _, err := range es.Errors() {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "; ")
}

// Result returns a single error if there's only one error in the set,
// the ErrorSet itself if there's more than one error, or nil if there are no errors.
// Example:
//
//	errSet := &ErrorSet{}
//	errSet.Add(errors.New("First error"))
//	errSet.Add(errors.New("Second error"))
//	err := errSet.Result()
//	fmt.Println(err) // Outputs: "First error; Second error"
func (es *ErrorSet) Result() error {
	if len(es.errors) == 1 {
		return es.errors[0]
	}

	if len(es.errors) > 1 {
		return es
	}

	return nil
}

// IsEmpty checks if the ErrorSet has no errors.
// It returns true if the ErrorSet is empty, otherwise false.
func (me *ErrorSet) IsEmpty() bool {
	return len(me.errors) == 0
}
