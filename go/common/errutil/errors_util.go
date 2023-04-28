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

// InternalError is the implementation of the SystemError interface that's used for system consumption only
type InternalError struct {
	msg string
	err error
}

func NewInternalErr(err error) error {
	return &InternalError{
		msg: err.Error(),
		err: err,
	}
}

func (e InternalError) Error() string {
	return e.msg
}

func (e InternalError) Unwrap() error {
	return e.err
}

func (e InternalError) Is(target error) bool {
	_, ok := target.(*InternalError) //nolint: errorlint
	return ok || errors.Is(e.err, target)
}
