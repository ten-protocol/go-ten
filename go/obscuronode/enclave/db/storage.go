package db

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash and true, or (nil, false) if no such Block is stored
	FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block) bool
	// HeightBlock returns the height of the L1 Block
	HeightBlock(block *types.Block) uint64
	// ParentBlock returns the L1 Block's parent and true, or (nil, false)  if no parent Block is stored
	ParentBlock(block *types.Block) (*types.Block, bool)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool
	// FetchHeadBlock - returns the head of the current chain
	FetchHeadBlock() (*types.Block, uint64)

	ProofHeight(rollup *core.Rollup) int64
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
}

type StateStorage interface {
	// FetchBlockState returns the state after ingesting the L1 Block with the given hash
	FetchBlockState(hash obscurocommon.L1RootHash) (*BlockState, bool)
	// FetchHeadState returns the state after ingesting the L1 Block at the head of the chain
	FetchHeadState() *BlockState
	// SetBlockState persists the state after ingesting the L1 Block with the given hash
	SetBlockState(hash obscurocommon.L1RootHash, state *BlockState)
	// FetchRollupState returns the state after adding the rollup with the given hash
	FetchRollupState(hash obscurocommon.L2RootHash) *State
	// SetRollupState persists the state after adding the rollup with the given hash
	SetRollupState(hash obscurocommon.L2RootHash, state *State)
}

type MempoolManager interface {
	// FetchMempoolTxs returns all L2 transactions in the mempool
	FetchMempoolTxs() []nodecommon.L2Tx
	// AddMempoolTx adds an L2 transaction to the mempool
	AddMempoolTx(tx nodecommon.L2Tx)
	// RemoveMempoolTxs removes any L2 transactions whose hash is keyed in the map from the mempool
	RemoveMempoolTxs(toRemove map[common.Hash]common.Hash)
	// FetchRollupTxs returns all transactions in a given rollup keyed by hash and true, or (nil, false) if the rollup is unknown
	FetchRollupTxs(rollup *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool)
	// StoreRollupTxs overwrites the transactions associated with a given rollup
	StoreRollupTxs(rollup *core.Rollup, newTxs map[common.Hash]nodecommon.L2Tx)
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret
	FetchSecret() core.SharedEnclaveSecret
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret core.SharedEnclaveSecret)

	// StoreGenesisRollup stores the rollup genesis
	StoreGenesisRollup(rol *core.Rollup)
	// FetchGenesisRollup returns the rollup genesis
	FetchGenesisRollup() *core.Rollup
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	RollupResolver
	MempoolManager
	SharedSecretStorage
	StateStorage
}

type storageImpl struct {
	db *inMemoryDB
}

func (s *storageImpl) StoreGenesisRollup(rol *core.Rollup) {
	s.db.StoreGenesisRollup(rol)
}

func (s *storageImpl) FetchGenesisRollup() *core.Rollup {
	return s.db.FetchGenesisRollup()
}

func NewStorage() Storage {
	db := NewInMemoryDB()
	return &storageImpl{db: db}
}

func (s *storageImpl) FetchBlockState(hash obscurocommon.L1RootHash) (*BlockState, bool) {
	s.assertSecretAvailable()
	return s.db.FetchBlockState(hash)
}

func (s *storageImpl) SetBlockState(hash obscurocommon.L1RootHash, state *BlockState) {
	s.assertSecretAvailable()
	if state.FoundNewRollup {
		s.db.SetBlockStateNewRollup(hash, state)
	} else {
		s.db.SetBlockState(hash, state)
	}
}

func (s *storageImpl) SetRollupState(hash obscurocommon.L2RootHash, state *State) {
	s.assertSecretAvailable()
	s.db.SetRollupState(hash, state)
}

func (s *storageImpl) FetchHeadState() *BlockState {
	s.assertSecretAvailable()
	val, _ := s.db.FetchBlockState(s.db.FetchHeadBlock())
	return val
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) {
	s.assertSecretAvailable()
	s.db.StoreRollup(rollup)
}

func (s *storageImpl) FetchRollup(hash obscurocommon.L2RootHash) (*core.Rollup, bool) {
	s.assertSecretAvailable()
	return s.db.FetchRollup(hash)
}

func (s *storageImpl) FetchRollups(height uint64) []*core.Rollup {
	s.assertSecretAvailable()
	return s.db.FetchRollups(height)
}

func (s *storageImpl) FetchRollupState(hash obscurocommon.L2RootHash) *State {
	s.assertSecretAvailable()
	return s.db.FetchRollupState(hash)
}

func (s *storageImpl) AddMempoolTx(tx nodecommon.L2Tx) {
	s.assertSecretAvailable()
	s.db.AddMempoolTx(tx)
}

func (s *storageImpl) FetchMempoolTxs() []nodecommon.L2Tx {
	s.assertSecretAvailable()
	return s.db.FetchMempoolTxs()
}

func (s *storageImpl) RemoveMempoolTxs(toRemove map[common.Hash]common.Hash) {
	s.assertSecretAvailable()
	s.db.RemoveMempoolTxs(toRemove)
}

func (s *storageImpl) StoreBlock(b *types.Block) bool {
	s.assertSecretAvailable()

	var height uint64
	if b.ParentHash() == obscurocommon.GenesisHash {
		height = obscurocommon.L1GenesisHeight
	} else {
		bAndHeight, f := s.db.FetchBlockAndHeight(b.ParentHash())
		if !f {
			log.Log(fmt.Sprintf("unable to store Block: %s without its parent: %s", b.Hash(), b.ParentHash()))
			return false
		}
		height = bAndHeight.height + 1
	}

	s.db.StoreBlock(b, height)
	return true
}

func (s *storageImpl) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	s.assertSecretAvailable()
	val, f := s.db.FetchBlockAndHeight(hash)
	var block *types.Block
	if f {
		block = val.b
	}
	return block, f
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, uint64) {
	s.assertSecretAvailable()
	bh, _ := s.db.FetchBlockAndHeight(s.db.FetchHeadBlock())
	return bh.b, bh.height
}

func (s *storageImpl) FetchRollupTxs(r *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool) {
	s.assertSecretAvailable()
	return s.db.FetchRollupTxs(r)
}

func (s *storageImpl) StoreRollupTxs(r *core.Rollup, newTxs map[common.Hash]nodecommon.L2Tx) {
	s.assertSecretAvailable()
	s.db.StoreRollupTxs(r, newTxs)
}

func (s *storageImpl) StoreSecret(secret core.SharedEnclaveSecret) {
	s.db.StoreSecret(secret)
}

func (s *storageImpl) FetchSecret() core.SharedEnclaveSecret {
	return s.db.FetchSecret()
}

func (s *storageImpl) HeightBlock(block *types.Block) uint64 {
	s.assertSecretAvailable()
	val, f := s.db.FetchBlockAndHeight(block.Hash())
	if !f {
		panic("should not happen")
	}
	return val.height
}

func (s *storageImpl) ParentRollup(r *core.Rollup) *core.Rollup {
	s.assertSecretAvailable()
	parent, found := s.db.FetchRollup(r.Header.ParentHash)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", r.Hash()))
	}
	return parent
}

func (s *storageImpl) ParentBlock(b *types.Block) (*types.Block, bool) {
	s.assertSecretAvailable()
	return s.FetchBlock(b.Header().ParentHash)
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	s.assertSecretAvailable()
	if maybeAncestor.Hash() == block.Hash() {
		return true
	}

	if s.HeightBlock(maybeAncestor) >= s.HeightBlock(block) {
		return false
	}

	p, f := s.ParentBlock(block)
	if !f {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	s.assertSecretAvailable()
	if maybeAncestor == block.Hash() {
		return true
	}

	if maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if s.HeightBlock(block) == obscurocommon.L1GenesisHeight {
		return false
	}

	resolvedBlock, found := s.FetchBlock(maybeAncestor)
	if found {
		if s.HeightBlock(resolvedBlock) >= s.HeightBlock(block) {
			return false
		}
	}

	p, f := s.ParentBlock(block)
	if !f {
		return false
	}

	return s.IsBlockAncestor(p, maybeAncestor)
}

func (s *storageImpl) assertSecretAvailable() {
	// TODO uncomment this
	//if s.FetchSecret() == nil {
	//	panic("Enclave not initialized")
	//}
}

// ProofHeight - return the height of the L1 proof, or GenesisHeight - if the block is not known
// todo - find a better way. This is a workaround to handle rollups created with proofs that haven't propagated yet
func (s *storageImpl) ProofHeight(r *core.Rollup) int64 {
	v, f := s.FetchBlock(r.Header.L1Proof)
	if !f {
		return -1
	}
	return int64(s.HeightBlock(v))
}

func (s *storageImpl) Proof(r *core.Rollup) *types.Block {
	v, f := s.FetchBlock(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}
