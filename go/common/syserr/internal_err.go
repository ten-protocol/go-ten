package syserr

import "errors"

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
