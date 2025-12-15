package limiters

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ethereum/go-ethereum/rlp"
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
	// encodingOverheadFactor accounts for RLP structure encoding, encryption padding, blob encoding overhead,
	// and rollup header/structure overhead that are not captured by just measuring transaction sizes.
	// This factor is deliberately conservative (2.0 = 100% overhead) because large rollups with many batches
	// can have significant overhead from batch headers, deltas, and structural encoding.
	// Without this, large rollups can exceed the blob size limit (e.g., 543KB of transactions -> 638KB final)
	encodingOverheadFactor = 2.0
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

	// Calculate estimated final size accounting for:
	// 1. Compression (txCompressionFactor)
	// 2. RLP structure encoding, encryption, and blob encoding overhead (encodingOverheadFactor)
	// 3. Compressed header size per batch (compressedHeaderSize)
	estimatedSize := uint64(float64(len(encodedData))*txCompressionFactor*encodingOverheadFactor) + compressedHeaderSize
	if estimatedSize > rl.remainingSize {
		return false, nil
	}

	rl.remainingSize -= estimatedSize
	return true, nil
}
