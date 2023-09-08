package zeus

import (
	"github.com/otoru/zeus/hooks"
)

// Hooks is a facade for hooks.Hooks
type Hooks hooks.Hooks

// ErrorSet is a facade for errs.ErrorSet
type ErrorSet interface {
	IsEmpty() bool
	Result() error
	Error() string
	Errors() []error
	Add(err error)
}
