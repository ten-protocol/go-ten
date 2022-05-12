package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash and true, or (nil, false) if no such Block is stored
	FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block) bool
	// ParentBlock returns the L1 Block's parent and true, or (nil, false)  if no parent Block is stored
	ParentBlock(block *types.Block) (*types.Block, bool)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool
	// FetchHeadBlock - returns the head of the current chain
	FetchHeadBlock() *types.Block
	// ProofHeight - return the height of the L1 proof, or GenesisHeight - if the block is not known
	ProofHeight(rollup *core.Rollup) int64
	// Proof - returns the block used as proof for the rollup
	Proof(rollup *core.Rollup) *types.Block
}

type RollupResolver interface {
	// FetchRollup returns the rollup with the given hash and true, or (nil, false) if no such rollup is stored
	FetchRollup(hash obscurocommon.L2RootHash) (*core.Rollup, bool)
	// FetchRollups returns all the proposed rollups with the given height
	FetchRollups(height uint64) []*core.Rollup
	// StoreRollup persists the rollup
	StoreRollup(rollup *core.Rollup)
	// ParentRollup returns the rollup's parent rollup
	ParentRollup(rollup *core.Rollup) *core.Rollup
	// FetchRollupTxs returns all transactions in a given rollup keyed by hash and true, or (nil, false) if the rollup is unknown
	FetchRollupTxs(rollup *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool)
	// StoreRollupTxs overwrites the transactions associated with a given rollup
	StoreRollupTxs(rollup *core.Rollup, newTxs map[common.Hash]nodecommon.L2Tx)
	// StoreGenesisRollup stores the rollup genesis
	StoreGenesisRollup(rol *core.Rollup)
	// FetchGenesisRollup returns the rollup genesis
	FetchGenesisRollup() *core.Rollup
}

type BlockStateStorage interface {
	// FetchBlockState returns the head rollup found in the block
	FetchBlockState(blockHash obscurocommon.L1RootHash) (*core.BlockState, bool)
	// FetchHeadState returns the head rollup
	FetchHeadState() *core.BlockState
	// SetBlockState save the rollup-block mapping
	SetBlockState(blockHash obscurocommon.L1RootHash, state *core.BlockState, rollup *core.Rollup)
	// CreateStateDB create a database that can be used to execute transactions
	CreateStateDB(hash obscurocommon.L2RootHash) StateDB
	// GenesisStateDB create the original empty StateDB
	GenesisStateDB() StateDB
}

// StateDB - is the conceptual equivalent of the geth vm.StateDB
type StateDB interface {
	GetBalance(address common.Address) uint64
	SetBalance(address common.Address, balance uint64)
	AddWithdrawal(txHash obscurocommon.TxHash)
	Copy() StateDB
	StateRoot() common.Hash
	Withdrawals() []obscurocommon.TxHash

	// Commit saves the changes made during transaction execution to a persistent db
	Commit(currentRoot obscurocommon.L2RootHash)
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret
	FetchSecret() core.SharedEnclaveSecret
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret core.SharedEnclaveSecret)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	RollupResolver
	SharedSecretStorage
	BlockStateStorage
}
