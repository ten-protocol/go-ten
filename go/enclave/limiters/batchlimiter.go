package limiters

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// BatchSizeLimiter - Acts as a limiter for batches based
// the data from the transaction that we have to publish to the l1.
// Acts as a calldata reservation system that accounts for both
// transactions and cross chain messages.
type batchSizeLimiter struct {
	Size uint64
}

type BatchSizeLimiter interface {
	AcceptTransaction(tx *types.Transaction) error
	// ProcessReceipt(receipt *types.Receipt) error //todo @stefan add this again
}

var ErrInsufficientSpace = errors.New("insufficient space in BatchSizeLimiter")

// BatchMaxTransactionData - The number where we will cut off processing transactions inside the evm facade.
const BatchMaxTransactionData = 25_000

// NewBatchSizeLimiter - Size is the total space available per batch for calldata in a rollup.
// contractAddress - the address of the l2 message bus where cross chain events would originate from.
// topic - the event id of the cross chain message event.
func NewBatchSizeLimiter(size uint64) BatchSizeLimiter {
	return &batchSizeLimiter{
		Size: size,
	}
}

// AcceptTransaction - transaction is rlp encoded as it normally would be when publishing a rollup and
// its size is deducted from the remaining limit.
func (l *batchSizeLimiter) AcceptTransaction(tx *types.Transaction) error {
	rlpSize, err := getRlpSize(tx)
	if err != nil {
		return err
	}

	if uint64(rlpSize) > l.Size {
		return ErrInsufficientSpace
	}

	l.Size -= uint64(rlpSize)
	return nil
}

func getRlpSize(val interface{}) (int, error) {
	// todo (@stefan) - this should have a coefficient for compression
	enc, err := rlp.EncodeToBytes(val)
	if err != nil {
		return 0, err
	}

	return len(enc), nil
}

type unlimitedBatchSize struct{}

func NewUnlimitedSizePool() BatchSizeLimiter {
	return &unlimitedBatchSize{}
}

func (*unlimitedBatchSize) AcceptTransaction(*types.Transaction) error {
	return nil
}

func (*unlimitedBatchSize) ProcessReceipt(*types.Receipt) error {
	return nil
}
