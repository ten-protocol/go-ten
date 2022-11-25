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
