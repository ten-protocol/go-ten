package errutil

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrNoImpl   = errors.New("not implemented")
)
