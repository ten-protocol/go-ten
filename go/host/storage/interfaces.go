package storage

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"math/big"
)

type BatchResolver interface {
	// FetchBatch returns the batch with the given hash.
	FetchBatch(hash common.L2BatchHash) (*core.Batch, error)
	// FetchBatchHeader returns the batch header with the given hash.
	FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error)
	// FetchBatchByHeight returns the batch on the canonical chain with the given height.
	FetchBatchByHeight(height uint64) (*core.Batch, error)
	// FetchBatchBySeqNo returns the batch with the given seq number.
	FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error)
	// FetchHeadBatch returns the current head batch of the canonical chain.
	FetchHeadBatch() (*core.Batch, error)
	// FetchCurrentSequencerNo returns the sequencer number
	FetchCurrentSequencerNo() (*big.Int, error)
	// FetchBatchesByBlock returns all batches with the block hash as the L1 proof
	FetchBatchesByBlock(common.L1BlockHash) ([]*core.Batch, error)
	//// FetchNonCanonicalBatchesBetween - returns all reorged batches between the sequences
	//FetchNonCanonicalBatchesBetween(startSeq uint64, endSeq uint64) ([]*core.Batch, error)
	//// FetchCanonicalUnexecutedBatches - return the list of the unexecuted batches that are canonical
	//FetchCanonicalUnexecutedBatches(*big.Int) ([]*core.Batch, error)

	//FetchConvertedHash(hash common.L2BatchHash) (gethcommon.Hash, error)

	// BatchWasExecuted - return true if the batch was executed
	BatchWasExecuted(hash common.L2BatchHash) (bool, error)

	// FetchHeadBatchForBlock returns the hash of the head batch at a given L1 block.
	FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error)

	// StoreBatch stores an un-executed batch.
	StoreBatch(batch *core.Batch, convertedHash gethcommon.Hash) error
	// StoreExecutedBatch - store the batch after it was executed
	StoreExecutedBatch(batch *core.Batch, receipts []*types.Receipt) error

	// StoreRollup
	StoreRollup(rollup *common.ExtRollup, header *common.CalldataRollupHeader) error
	FetchReorgedRollup(reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error)
}
