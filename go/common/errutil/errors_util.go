package errutil

import (
	"errors"

	"github.com/ethereum/go-ethereum"
)

var (
	// ErrNotFound must equal Geth's not-found error. This is because some Geth components we use throw the latter, and
	// we want to be able to catch both types in a single error-check.
	ErrNotFound = ethereum.NotFound
	ErrNoImpl   = errors.New("not implemented")
)

// SystemError represents an error that's for system consumption only
type SystemError struct {
	msg string
	err error
}

func NewSystemErr(err error) error {
	return &SystemError{
		msg: err.Error(),
		err: err,
	}
}

func (e SystemError) Error() string {
	return e.msg
}

func (e SystemError) Unwrap() error {
	return e.err
}

func (e SystemError) Is(target error) bool {
	_, ok := target.(*SystemError) //nolint: errorlint
	return ok || errors.Is(e.err, target)
}
