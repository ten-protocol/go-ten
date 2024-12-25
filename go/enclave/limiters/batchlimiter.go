package limiters

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common/compression"
)

// BatchSizeLimiter - Acts as a limiter for batches based
// the data from the transaction that we have to publish to the l1.
// Acts as a calldata reservation system that accounts for both
// transactions and cross chain messages.
type batchSizeLimiter struct {
	compressionService compression.DataCompressionService
	remainingSize      uint64 // the available size in the limiter
}

// NewBatchSizeLimiter - Size is the total space available per batch for calldata in a rollup.
func NewBatchSizeLimiter(size uint64, compressionService compression.DataCompressionService) BatchSizeLimiter {
	return &batchSizeLimiter{
		compressionService: compressionService,
		remainingSize:      size,
	}
}

// AcceptTransaction - transaction is rlp encoded as it normally would be when publishing a rollup and
// its size is deducted from the remaining limit.
func (l *batchSizeLimiter) AcceptTransaction(tx *types.Transaction) error {
	rlpSize, err := l.getCompressedSize(tx)
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
func (l *batchSizeLimiter) getCompressedSize(val interface{}) (int, error) {
	enc, err := rlp.EncodeToBytes(val)
	if err != nil {
		return 0, err
	}

	// compress the transaction. This is useless for small transactions, but might be useful for larger transactions such as deploying contracts
	// todo - keep a running compression of the current batch
	compr, err := l.compressionService.CompressBatch(enc)
	if err != nil {
		return 0, err
	}

	return len(compr), nil
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
