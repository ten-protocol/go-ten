package syserr

import (
	"errors"
)

// InternalError is the implementation of the SystemError interface that's used for system consumption only
// represents errors at the enclave layer
type InternalError struct {
	*BaseCustomError
}

// RPCError is the implementation of the SystemError interface that's used for system consumption only
// represents errors at the RPC layer
type RPCError struct {
	*BaseCustomError
}

func NewInternalError(err error) error {
	return &InternalError{
		BaseCustomError: &BaseCustomError{
			msg: err.Error(),
			err: err,
		},
	}
}

func NewRPCError(err error) error {
	return &RPCError{
		BaseCustomError: &BaseCustomError{
			msg: err.Error(),
			err: err,
		},
	}
}

func (e InternalError) Is(target error) bool {
	_, ok := target.(*InternalError) //nolint: errorlint
	return ok || errors.Is(e.err, target)
}

func (e RPCError) Is(target error) bool {
	_, ok := target.(*RPCError) //nolint: errorlint
	return ok || errors.Is(e.err, target)
}

// BaseCustomError is the base error for custom errors
type BaseCustomError struct {
	msg string
	err error
}

func (e BaseCustomError) Error() string {
	return e.msg
}

func (e BaseCustomError) Unwrap() error {
	return e.err
}
