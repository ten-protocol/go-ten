package errutil

import (
	"errors"

	"github.com/ethereum/go-ethereum"
)

var (
	ErrNotFound = ethereum.NotFound
	ErrNoImpl   = errors.New("not implemented")
)
