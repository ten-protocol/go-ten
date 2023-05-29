package limiters

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

type BatchSizeLimiter interface {
	AcceptTransaction(tx *types.Transaction) error
	// ProcessReceipt(receipt *types.Receipt) error //todo @stefan add this again
}

var ErrInsufficientSpace = errors.New("insufficient space in BatchSizeLimiter")

type RollupLimiter interface {
	AcceptBatch(encodable interface{}) (bool, error)
}
