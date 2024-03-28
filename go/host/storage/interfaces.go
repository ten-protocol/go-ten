package storage

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"io"
)

type Storage interface {
	BatchResolver
	BlockResolver
	DatabaseResolver
	io.Closer
}

type BatchResolver interface {
	// AddBatch stores the batch
	AddBatch(batch *common.ExtBatch) error
	// FetchBatchBySeqNo returns the batch with the given seq number.
	FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error)
	//// FetchBatch returns the batch with the given hash.
	//FetchBatch(hash common.L2BatchHash) (*core.Batch, error)
	//// FetchBatchHeader returns the batch header with the given hash.
	//FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error)
	//// FetchBatchByHeight returns the batch on the canonical chain with the given height.
	//FetchBatchByHeight(height uint64) (*core.Batch, error)
	//// FetchBatchBySeqNo returns the batch with the given seq number.
	//FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error)
	//// FetchHeadBatch returns the current head batch of the canonical chain.
	//FetchHeadBatch() (*core.Batch, error)
	//// FetchCurrentSequencerNo returns the sequencer number
	//FetchCurrentSequencerNo() (*big.Int, error)
	//// FetchBatchesByBlock returns all batches with the block hash as the L1 proof
	//FetchBatchesByBlock(common.L1BlockHash) ([]*core.Batch, error)
	//// BatchWasExecuted - return true if the batch was executed
	//BatchWasExecuted(hash common.L2BatchHash) (bool, error)
	//// FetchHeadBatchForBlock returns the hash of the head batch at a given L1 block.
	//FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error)
}

type BlockResolver interface {
	AddBlock(b *types.Header, rollupHash common.L2RollupHash) error
	AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error
}

type DatabaseResolver interface {
	GetDB() *sql.DB
}
