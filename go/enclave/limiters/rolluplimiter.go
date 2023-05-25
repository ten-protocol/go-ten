package limiters

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrFailedToEncode = errors.New("failed to encode data")
)

// MaxTransactionSizeLimiter - configured to be close to what the ethereum clients
// have configured as the maximum size a transaction can have. Note that this isn't
// a protocol limit, but a miner imposed limit and it might be hard to find someone
// to include a transaction if it goes above it
// todo - figure out the best number, optimism uses 132KB
const MaxTransactionSize = 64 * 1024

type rollupLimiter struct {
	remainingSize uint64
}

func NewRollupLimiter(size uint64) RollupLimiter {
	return &rollupLimiter{
		remainingSize: size,
	}
}

// todo (@stefan) figure out how to optimize the serialization out of the limiter
func (rl *rollupLimiter) AcceptBatch(encodable interface{}) (bool, error) {
	encodedData, err := rlp.EncodeToBytes(encodable)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrFailedToEncode, err)
	}

	encodedSize := uint64(len(encodedData))
	if encodedSize > rl.remainingSize {
		return true, nil
	}

	rl.remainingSize -= encodedSize
	return false, nil
}
