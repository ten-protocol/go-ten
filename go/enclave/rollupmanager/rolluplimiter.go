package rollupmanager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrFailedToEncode   = errors.New("failed to encode data")
	ErrSizeExceedsLimit = errors.New("data size exceeds remaining limit")
)

// MaxTransactionSizeLimiter - configured to be close to what the ethereum clients
// have configured as the maximum size a transaction can have. Note that this isn't
// a protocol limit, but a miner imposed limit and it might be hard to find someone
// to include a transaction if it goes above it
// todo - figure out the best number, optimism uses 132KB
const MaxTransactionSizeLimiter = RollupLimiter(64 * 1024)

type RollupLimiter uint64

func (rl *RollupLimiter) Consume(encodable interface{}) error {
	encodedData, err := rlp.EncodeToBytes(encodable)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToEncode, err)
	}

	encodedSize := uint64(len(encodedData))
	if encodedSize > uint64(*rl) {
		return fmt.Errorf("%w: data size %d, remaining limit %d", ErrSizeExceedsLimit, encodedSize, *rl)
	}

	*rl -= RollupLimiter(encodedSize)
	return nil
}
