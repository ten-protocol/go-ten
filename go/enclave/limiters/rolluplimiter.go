package limiters

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ethereum/go-ethereum/rlp"
)

const (
	// 85% is a very conservative number. It will most likely be 66% in practice.
	// We can lower it, once we have a mechanism in place to handle batches that don't actually compress to that.
	txCompressionFactor  = 0.85
	// Based on mainnet logs: 86,695 batches → 176,958 bytes (way over 90KB limit)
	// Need much more aggressive limiting: target ~15,000 batches for safety
	// Target: 90,000 ÷ 15,000 = 6 bytes/batch estimated
	// With 0.85 factor: 2.55 + 3.5 = 6.05 bytes/batch
	compressedHeaderSize = 4
)

type rollupLimiter struct {
	remainingSize uint64
}

func NewRollupLimiter(size uint64) RollupLimiter {
	return &rollupLimiter{
		remainingSize: size,
	}
}

// todo (@stefan) figure out how to optimize the serialization out of the limiter
func (rl *rollupLimiter) AcceptBatch(batch *core.Batch) (bool, error) {
	encodedData, err := rlp.EncodeToBytes(batch.Transactions)
	if err != nil {
		return false, fmt.Errorf("failed to encode data. Cause: %w", err)
	}

	// adjust with a compression factor and add the size of a compressed batch header
	encodedSize := uint64(float64(len(encodedData))*txCompressionFactor) + compressedHeaderSize
	if encodedSize > rl.remainingSize {
		return false, nil
	}

	rl.remainingSize -= encodedSize
	return true, nil
}
