package limiters

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

const (
	// Based on testing: compressedHeaderSize=4 gives ~17K batches (too conservative)
	// Target: ~25K batches for better utilization of 90KB limit
	// Current: 2.06 bytes/batch actual, target estimation: 3.6 bytes/batch
	compressedHeaderSize = 2
)

type rollupLimiter struct {
	remainingSize uint64
	logger        gethlog.Logger
}

func NewRollupLimiter(size uint64, logger gethlog.Logger) RollupLimiter {
	return &rollupLimiter{
		remainingSize: size,
		logger:        logger,
	}
}

// AcceptBatch estimates the size of a batch's transaction payload to determine if it fits in the rollup
func (rl *rollupLimiter) AcceptBatch(batch *core.Batch) (bool, error) {
	// Encode transactions with timestamp deltas, matching the actual rollup compression format
	// This gives a realistic size estimate instead of overestimating
	txsAndTimestamps := common.CreateTxsAndTimeStamp(batch.Transactions, batch.Header.Time)
	txBytes, err := rlp.EncodeToBytes(txsAndTimestamps)
	if err != nil {
		return false, fmt.Errorf("failed to encode batch transactions. Cause: %w", err)
	}

	// Estimate the compressed size - the actual rollup will compress this
	// The rollup uses CompressRollup with brotli.BestCompression on the combined batch payloads
	// Based on real data, this achieves ~10-15% of original size for transaction-heavy batches
	// Use 15% to be conservative and avoid overshooting the blob limit
	compressionFactor := 0.15
	estimatedCompressedSize := uint64(float64(len(txBytes))*compressionFactor) + compressedHeaderSize

	accepted := estimatedCompressedSize <= rl.remainingSize
	
	rl.logger.Debug("RollupLimiter AcceptBatch",
		"batch_seq", batch.SeqNo().Uint64(),
		"tx_count", len(batch.Transactions),
		"raw_bytes", len(txBytes),
		"estimated_compressed", estimatedCompressedSize,
		"remaining_space", rl.remainingSize,
		"accepted", accepted)

	if !accepted {
		return false, nil
	}

	rl.remainingSize -= estimatedCompressedSize
	return true, nil
}
