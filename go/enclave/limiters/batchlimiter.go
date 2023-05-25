package limiters

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// BatchSizeLimiter - Acts as a limiter for batches based
// the data from the transaction that we have to publish to the l1.
// Acts as a calldata reservation system that accounts for both
// transactions and cross chain messages.
type batchSizeLimiter struct {
	remainingSize uint64 // the available size in the limiter
}

// NewBatchSizeLimiter - Size is the total space available per batch for calldata in a rollup.
func NewBatchSizeLimiter(size uint64) BatchSizeLimiter {
	return &batchSizeLimiter{
		remainingSize: size,
	}
}

// AcceptTransaction - transaction is rlp encoded as it normally would be when publishing a rollup and
// its size is deducted from the remaining limit.
func (l *batchSizeLimiter) AcceptTransaction(tx *types.Transaction) error {
	rlpSize, err := getRlpSize(tx)
	if err != nil {
		return err
	}

	if uint64(rlpSize) > l.remainingSize {
		return ErrInsufficientSpace
	}

	l.remainingSize -= uint64(rlpSize)
	return nil
}

// todo (@stefan) figure out how to optimize the serialization out of the limiter
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
