package limiters

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

const (
	// 85% is a very conservative number. It will most likely be 66% in practice.
	// We can lower it, once we have a mechanism in place to handle batches that don't actually compress to that.
	txCompressionFactor = 0.85
	// Based on testing: compressedHeaderSize=4 gives ~17K batches (too conservative)
	// Target: ~25K batches for better utilization of 90KB limit
	// Current: 2.06 bytes/batch actual, target estimation: 3.6 bytes/batch
	// With 0.85 factor: 2.55 + 1.5 = 4.05 bytes/batch ≈ 22,222 batches
	compressedHeaderSize = 2
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
	// Encode transactions with timestamp deltas, matching the rollup compression format
	txsAndTimestamps := common.CreateTxsAndTimeStamp(batch.Transactions, batch.Header.Time)
	txBytes, err := rlp.EncodeToBytes(txsAndTimestamps)  // <-- FIX: matches actual rollup structure
	if err != nil {
		return false, fmt.Errorf("failed to encode batch transactions. Cause: %w", err)
	}

	// The actual rollup will compress this payload, so we estimate based on the raw size
	// Note: we don't actually compress here to avoid the CPU cost, but use a realistic estimate
	// based on the actual encoded size rather than just counting transactions
	estimatedSize := uint64(len(txBytes)) + compressedHeaderSize
	if estimatedSize > rl.remainingSize {
		return false, nil
	}

	rl.remainingSize -= estimatedSize
	return true, nil
}
