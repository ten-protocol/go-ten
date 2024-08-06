package components

import (
	"context"
	"errors"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

var ErrDuplicateRollup = errors.New("duplicate rollup received")

type BlockIngestionType struct {
	// PreGenesis is true if there is no stored L1 head block.
	// (L1 head is only stored when there is an L2 state to associate it with. Soon we will start consuming from the
	// genesis block and then, we should only see one block ingested in a 'PreGenesis' state)
	PreGenesis bool

	// ChainFork contains information about the status of the new block in the chain
	ChainFork *common.ChainFork

	// Block that is already on the canonical chain
	OldCanonicalBlock bool
}

func (bit *BlockIngestionType) IsFork() bool {
	if bit.ChainFork == nil {
		return false
	}
	return bit.ChainFork.IsFork()
}

type L1BlockProcessor interface {
	Process(ctx context.Context, br *common.BlockAndReceipts) (*BlockIngestionType, error)
	GetHead(context.Context) (*common.L1Block, error)
	GetCrossChainContractAddress() *gethcommon.Address
	HealthCheck() (bool, error)
}

// BatchExecutionContext - Contains all of the data that each batch depends on
type BatchExecutionContext struct {
	BlockPtr     common.L1BlockHash // Block is needed for the cross chain messages
	ParentPtr    common.L2BatchHash
	Transactions common.L2Transactions
	AtTime       uint64
	Creator      gethcommon.Address
	ChainConfig  *params.ChainConfig
	SequencerNo  *big.Int
	BaseFee      *big.Int
}

// ComputedBatch - a structure representing the result of a batch
// computation where `Batch` is the newly computed batch and `Receipts`
// are the receipts for the executed transactions inside this batch.
// The `Commit` function allows for committing the stateDB resulting from
// the computation of the batch. One might not want to commit in case the
// resulting batch differs than what is being validated for example.
type ComputedBatch struct {
	Batch    *core.Batch
	Receipts types.Receipts
	// while executing the batch, we collect the newly created contracts mapped by the transaction that created them
	CreatedContracts map[gethcommon.Hash][]*gethcommon.Address
	Commit           func(bool) (gethcommon.Hash, error)
}

type BatchExecutor interface {
	// ComputeBatch - a more primitive ExecuteBatch
	// Call with same BatchContext should always produce identical extBatch - idempotent
	// Should be safe to call in parallel
	// failForEmptyBatch bool is used to skip batch production
	ComputeBatch(ctx context.Context, batchContext *BatchExecutionContext, failForEmptyBatch bool) (*ComputedBatch, error)

	// ExecuteBatch - executes the transactions and xchain messages, returns the receipts, and updates the stateDB
	ExecuteBatch(context.Context, *core.Batch) (types.Receipts, map[gethcommon.Hash][]*gethcommon.Address, error)

	// CreateGenesisState - will create and commit the genesis state in the stateDB for the given block hash,
	// and uint64 timestamp representing the time now. In this genesis state is where one can
	// find preallocated funds for faucet. TODO - make this an option
	CreateGenesisState(context.Context, common.L1BlockHash, uint64, gethcommon.Address, *big.Int) (*core.Batch, *types.Transaction, error)
}

type BatchRegistry interface {
	// BatchesAfter - Given a hash, will return batches following it until the head batch and the l1 blocks referenced by those batches
	BatchesAfter(ctx context.Context, batchSeqNo uint64, upToL1Height uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, []*types.Block, error)

	// GetBatchStateAtHeight - creates a stateDB for the block number
	GetBatchStateAtHeight(ctx context.Context, blockNumber *gethrpc.BlockNumber) (*state.StateDB, error)

	// GetBatchState - creates a stateDB for the block hash
	GetBatchState(ctx context.Context, hash *common.L2BatchHash) (*state.StateDB, error)

	// GetBatchAtHeight - same as `GetBatchStateAtHeight`, but instead returns the full batch
	// rather than its stateDB only.
	GetBatchAtHeight(ctx context.Context, height gethrpc.BlockNumber) (*core.Batch, error)

	// SubscribeForExecutedBatches - register a callback for new batches
	SubscribeForExecutedBatches(func(*core.Batch, types.Receipts))
	UnsubscribeFromBatches()

	OnBatchExecuted(batch *common.BatchHeader, receipts types.Receipts)
	OnL1Reorg(*BlockIngestionType)

	// HasGenesisBatch - returns if genesis batch is available yet or not, or error in case
	// the function is unable to determine.
	HasGenesisBatch() (bool, error)

	HeadBatchSeq() *big.Int

	HealthCheck() (bool, error)
}

type RollupProducer interface {
	// CreateInternalRollup - creates a rollup starting from the end of the last rollup that has been stored on the L1
	CreateInternalRollup(ctx context.Context, fromBatchNo uint64, upToL1Height uint64, limiter limiters.RollupLimiter) (*core.Rollup, error)
}

type RollupConsumer interface {
	// ProcessRollupsInBlock - extracts the rollup from the block's transactions
	// and verifies its integrity, saving and processing any batches that have
	// not been seen previously.
	ProcessRollupsInBlock(ctx context.Context, b *common.BlockAndReceipts) error
}

type BlobResolver interface {
	// FetchBlobs Fetches the blob data using beacon chain APIs
	FetchBlobs(ctx context.Context, b *types.Header, hashes []ethadapter.IndexedBlobHash) ([]*ethadapter.Blob, error)
}
