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

type BlockStatus uint16

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

const (
	Latest = BlockStatus(0)
	Fork   = BlockStatus(1)
)

type BlockConsumer interface {
	ConsumeBlock(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, error)
	GetHead() (*common.L1Block, error)
}

// Contains all of the data that each batch depends on
type BatchContext struct {
	BlockPtr     common.L1BlockHash // Block is needed for the cross chain messages
	ParentPtr    common.L2BatchHash
	Transactions common.L2Transactions
	AtTime       uint64
	Randomness   gethcommon.Hash
	Creator      gethcommon.Address
	ChainConfig  *params.ChainConfig
}

type ComputedBatch struct {
	Batch    *core.Batch
	Receipts types.Receipts
	Commit   func(bool) (gethcommon.Hash, error)
}
type BatchProducer interface {
	// Call with same BatchContext should always produce identical extBatch - idempotent
	// Should be safe to call in parallel
	ComputeBatch(*BatchContext) (*ComputedBatch, error)
	CreateGenesisState(common.L1BlockHash, common.L2Address, uint64) (*core.Batch, *types.Transaction, error)
}

type BatchRegistry interface {
	StoreBatch(*core.Batch, types.Receipts) error
	GetHeadBatch() (*core.Batch, error)
	GetHeadBatchFor(common.L1BlockHash) (*core.Batch, error)
	GetBatch(common.L2BatchHash) (*core.Batch, error)
	Subscribe() chan *core.Batch
	Unsubscribe()
	FindAncestralBatchFor(*common.L1Block) (*core.Batch, error)
	HasGenesisBatch() (bool, error)
	GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error)
	BatchesAfter(batchHash gethcommon.Hash) ([]*core.Batch, error)
	GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error)
}

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	return append(txReceipts, depositReceipts...)
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }

type RollupProducer interface {
	CreateRollup() (*core.Rollup, error)
}

type RollupConsumer interface {
	// ProcessL1Block - extracts the rollups from the block's transactions
	// and verifies their integrity, saving and processing any batches that have
	// not been seenp previously.
	ProcessL1Block(b *common.BlockAndReceipts) ([]*core.Rollup, error)
}
