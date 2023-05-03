package syserr

import (
	"errors"
)

// InternalError is the implementation of the SystemError interface that's used for system consumption only
// represents errors at the enclave layer
type InternalError struct {
	msg string
	err error
}

// RPCError is the implementation of the SystemError interface that's used for system consumption only
// represents errors at the RPC layer
type RPCError struct {
	msg string
	err error
}

func NewRPCError(err error) error {
	return &RPCError{
		msg: err.Error(),
		err: err,
	}
}

func NewInternalError(err error) error {
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

func (e RPCError) Error() string {
	return e.msg
}

func (e RPCError) Unwrap() error {
	return e.err
}

func (e RPCError) Is(target error) bool {
	_, ok := target.(*RPCError) //nolint: errorlint
	return ok || errors.Is(e.err, target)
}
