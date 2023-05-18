package components

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

const (
	SubscriptionChannelBuffer = 10
)

type BlockIngestionType struct {
	// IsLatest is true if this block was the canonical head of the L1 chain at the time it was submitted to enclave
	// (if false then we are behind and catching up, expect to be fed another block immediately afterwards)
	IsLatest bool

	// Fork is true if the ingested block is on a different branch to previously known head
	// (resulting in rewinding of one or more blocks that we had previously considered canonical)
	Fork bool

	// PreGenesis is true if there is no stored L1 head block.
	// (L1 head is only stored when there is an L2 state to associate it with. Soon we will start consuming from the
	// genesis block and then, we should only see one block ingested in a 'PreGenesis' state)
	PreGenesis bool
}

type L1BlockProcessor interface {
	Process(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, error)
	GetHead() (*common.L1Block, error)
}

// Contains all of the data that each batch depends on
type BatchExecutionContext struct {
	BlockPtr     common.L1BlockHash // Block is needed for the cross chain messages
	ParentPtr    common.L2BatchHash
	Transactions common.L2Transactions
	AtTime       uint64
	Randomness   gethcommon.Hash
	Creator      gethcommon.Address
	ChainConfig  *params.ChainConfig
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
	Commit   func(bool) (gethcommon.Hash, error)
}
type BatchProducer interface {
	// ComputeBatch - will formulate a batch and execute the transactions according to the
	// state provided in the BatchContext.
	// Call with same BatchContext should always produce identical extBatch - idempotent
	// Should be safe to call in parallel
	ComputeBatch(*BatchExecutionContext) (*ComputedBatch, error)

	// CreateGenesisState - will create and commit the genesis state in the stateDB for the given block hash,
	// sequencer address and uint64 timestamp representing the time now. In this genesis state is where one can
	// find preallocated funds for faucet. TODO - make this an option
	CreateGenesisState(common.L1BlockHash, common.L2Address, uint64) (*core.Batch, *types.Transaction, error)
}

type BatchRegistry interface {
	// StoreBatch - will store the batch and receipts in storage.
	// Furthermore any heads and pointers would be updated here and
	// after all is done the batch will be pushed to the subscribers
	// in order to update them.
	StoreBatch(*core.Batch, types.Receipts) error

	// GetHeadBatch - Returns the batch considered to be the L2 head.
	GetHeadBatch() (*core.Batch, error)

	// GetHeadBatchFor - Returns the head batch for given L1 block.
	// Note that this does not seek backwards so if a given block does not have
	// any pointers to batches then `ErrNotFound` will be returned.
	GetHeadBatchFor(common.L1BlockHash) (*core.Batch, error)

	// GetBatch - Given a batch hash returns the matching batch taken from the
	// storage. If no such batch is found then nil and `ErrNotFound` should
	// be returned.
	GetBatch(common.L2BatchHash) (*core.Batch, error)

	// FindAncestralBatchFor - this function is akin to `GetHeadBatchFor`, but
	// it does not check the block passed as a parameter, but rather starts from its
	// parent block and then seeks backwards until a head batch is found or genesis is reached.
	FindAncestralBatchFor(*common.L1Block) (*core.Batch, error)

	// BatchesAfter - Given a hash, will return batches following it until the head batch
	BatchesAfter(batchHash gethcommon.Hash) ([]*core.Batch, error)

	// GetBatchStateAtHeight - creates a stateDB that represents the state committed when
	// the batch with height matching the blockNumber was created and stored.
	GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error)

	// GetBatchStateAtHeight - same as `GetBatchStateAtHeight`, but instead returns the full batch
	// rather than its stateDB only.
	GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error)

	// Subscribe - creates and returns a channel that will be used to push any newly created batches
	// to the subscriber.
	Subscribe(lastKnownHead *common.L2BatchHash) (chan *core.Batch, error)
	// Unsubscribe - informs the registry that the subscriber is no longer listening, allowing it to
	// gracefully terminate any streaming and stop queueing new batches.
	Unsubscribe()

	SubscribeForEvents() chan uint64
	UnsubscribeFromEvents()

	// HasGenesisBatch - returns if genesis batch is available yet or not, or error in case
	// the function is unable to determine.
	HasGenesisBatch() (bool, error)

	// CommitBatch - Commits any state changes from a batch computation to the database.
	CommitBatch(cb *ComputedBatch) error
}

type RollupProducer interface {
	// CreateRollup - creates a rollup starting from the end of the last rollup
	// that has been stored and continues it towards what we consider the current L2 head.
	CreateRollup() (*core.Rollup, error)
}

type RollupConsumer interface {
	// ProcessL1Block - extracts the rollup from the block's transactions
	// and verifies its integrity, saving and processing any batches that have
	// not been seen previously.
	ProcessL1Block(b *common.BlockAndReceipts) (*core.Rollup, error)
}
