package limiters

import (
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

// TestRollupLimiterWithReorgs simulates the Sepolia scenario where a large rollup with many reorged batches
// exceeded the blob size limit (638952 bytes) despite the limiter being set to 131,072 bytes.
// this happened because the limiter only measured transaction sizes but didn't account for:
// 1. full batch headers included for reorged batches (vs compressed deltas for canonical batches)
// 2. RLP structure encoding overhead
// 3. encryption padding
// 4. blob encoding overhead
func TestRollupLimiterWithReorgs(t *testing.T) {
	maxRollupSize := uint64(131072) // 128 KB - matching Sepolia config
	limiter := NewRollupLimiter(maxRollupSize)

	// - ~335 batches (489,081 to 489,555 in the logs, but pos 5915-6390 suggests ~475 total)
	// - each batch has reorg overhead (full header instead of delta)
	// - average ~10 transactions per batch (estimate based on typical load)

	const numBatches = 400
	const txsPerBatch = 10
	const avgTxSize = 200 // bytes per transaction

	acceptedBatches := 0
	for i := 0; i < numBatches; i++ {
		batch := createMockBatch(i, txsPerBatch, avgTxSize)

		accepted, err := limiter.AcceptBatch(batch)
		if err != nil {
			t.Fatalf("AcceptBatch returned error: %v", err)
		}

		if !accepted {
			break
		}
		acceptedBatches++
	}

	t.Logf("Limiter accepted %d out of %d batches (maxRollupSize=%d)", acceptedBatches, numBatches, maxRollupSize)

	// the old limiter (no encodingOverheadFactor) would accept too many batches
	// and result in a 638KB final rollup. With encodingOverheadFactor=2.0, it should reject batches much earlier to
	// stay under the limit.

	expectedBatches := 35
	tolerance := 10

	if acceptedBatches < expectedBatches-tolerance {
		t.Errorf("Limiter is too conservative: accepted %d batches, expected around %d (±%d)",
			acceptedBatches, expectedBatches, tolerance)
	}

	if acceptedBatches > expectedBatches+tolerance {
		t.Errorf("Limiter is too permissive: accepted %d batches, expected around %d (±%d). "+
			"This could result in rollups exceeding the blob size limit.",
			acceptedBatches, expectedBatches, tolerance)
	}
}

// TestRollupLimiterRejectsAfterLimit verifies that the limiter correctly rejects batches
// once the size limit is reached
func TestRollupLimiterRejectsAfterLimit(t *testing.T) {
	smallLimit := uint64(1000)
	limiter := NewRollupLimiter(smallLimit)

	// batch that should be rejected
	largeBatch := createMockBatch(0, 100, 50) // ~5000 bytes of transactions

	accepted, err := limiter.AcceptBatch(largeBatch)
	if err != nil {
		t.Fatalf("AcceptBatch returned error: %v", err)
	}

	if accepted {
		t.Errorf("Limiter should have rejected large batch (estimated size > %d)", smallLimit)
	}
}

// TestRollupLimiterAccumulatesCorrectly verifies that the limiter correctly tracks
// remaining size across multiple batches
func TestRollupLimiterAccumulatesCorrectly(t *testing.T) {
	maxSize := uint64(10000)
	limiter := NewRollupLimiter(maxSize)

	// several small batches
	for i := 0; i < 5; i++ {
		batch := createMockBatch(i, 5, 50) // ~250 bytes per batch, with factors ~550 bytes
		accepted, err := limiter.AcceptBatch(batch)
		if err != nil {
			t.Fatalf("AcceptBatch returned error: %v", err)
		}
		if !accepted {
			t.Fatalf("Batch %d should have been accepted", i)
		}
	}

	// batch that should exceed the limit
	largeBatch := createMockBatch(10, 50, 200) // ~10,000 bytes
	accepted, err := limiter.AcceptBatch(largeBatch)
	if err != nil {
		t.Fatalf("AcceptBatch returned error: %v", err)
	}

	if accepted {
		t.Errorf("Large batch should have been rejected after accumulating previous batches")
	}
}

// createMockBatch creates a mock batch with the specified number of transactions and average size
func createMockBatch(seqNo int, numTxs int, avgTxSize int) *core.Batch {
	txs := make([]*common.L2Tx, numTxs)

	for i := 0; i < numTxs; i++ {
		data := make([]byte, avgTxSize)
		for j := range data {
			data[j] = byte((seqNo + i + j) % 256)
		}

		tx := types.NewTransaction(
			uint64(i),
			gethcommon.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(1000),
			21000,
			big.NewInt(1000000000),
			data,
		)
		txs[i] = tx
	}

	return &core.Batch{
		Transactions: txs,
	}
}
