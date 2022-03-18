package enclave

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// BlockResolver - database of blocks indexed by the root hash
type BlockResolver interface {
	FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool)
	StoreBlock(block *types.Block)
	HeightBlock(block *types.Block) int
	Parent(b *types.Block) (*types.Block, bool)
	IsAncestor(blockA *types.Block, blockB *types.Block) bool
	IsBlockAncestor(l1BlockHash obscurocommon.L1RootHash, block *types.Block) bool
}

// Storage is the enclave's interface for interacting with the enclave's datastore.
type Storage interface {
	BlockResolver

	// Rollups
	FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool)
	StoreRollup(rollup *Rollup)
	FetchRollups(height int) []*Rollup
	ParentRollup(r *Rollup) *Rollup
	HeightRollup(r *Rollup) int

	// State
	FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool)
	SetBlockState(hash obscurocommon.L1RootHash, state BlockState)
	FetchState(hash obscurocommon.L2RootHash) State
	SetState(hash obscurocommon.L2RootHash, state State)
	Head() BlockState

	// Transactions
	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[common.Hash]common.Hash)
	FetchRollupTxs(r *Rollup) (map[common.Hash]L2Tx, bool)
	AddRollupTxs(*Rollup, map[common.Hash]L2Tx)

	// Shared secret
	StoreSecret(secret SharedEnclaveSecret)
	FetchSecret() SharedEnclaveSecret
}

type storageImpl struct {
	db               DB
	blockDB          db.Database
	rollupDB         *RollupDB
	statePerBlockDB  *BlockStateDB
	secretsDB        db.Database
	statePerRollupDB db.Database
}

func NewStorage() Storage {
	return &storageImpl{
		blockDB:          NewBlockDB(db.NewMemDB()),
		rollupDB:         NewRollupDB(db.NewMemDB()),
		statePerBlockDB:  NewBlockStateDB(db.NewMemDB()),
		statePerRollupDB: db.NewMemDB(),
		secretsDB:        db.NewMemDB(),
	}
}

func (s *storageImpl) FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool) {
	s.assertSecretAvailable()

	blockState, err := s.statePerBlockDB.Get(hash[:])
	if err != nil {
		panic(err)
	}

	return blockState, true
}

func (s *storageImpl) SetBlockState(hash obscurocommon.L1RootHash, state BlockState) {
	s.assertSecretAvailable()
	if state.foundNewRollup {
		err := s.statePerBlockDB.Store(hash[:], state)
		if err != nil {
			panic(err)
		}
		rollupHash := state.Head.Hash()
		err := s.rollupDB.Store(rollupHash[:], *state.Head)
		if err != nil {
			panic(err)
		}

		s.db.SetBlockStateNewRollup(hash, state)
	} else {
		s.db.SetBlockState(hash, state)
	}
}

func (s *storageImpl) SetState(hash obscurocommon.L2RootHash, state State) {
	s.assertSecretAvailable()
	s.db.SetState(hash, state)
}

func (s *storageImpl) Head() BlockState {
	s.assertSecretAvailable()
	val, _ := s.db.FetchBlockState(s.db.HeadBlock())
	return val
}

func (s *storageImpl) StoreRollup(rollup *Rollup) {
	s.assertSecretAvailable()
	height := s.HeightRollup(rollup)
	s.db.StoreRollup(rollup, height)
}

func (s *storageImpl) FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool) {
	s.assertSecretAvailable()
	return s.db.FetchRollup(hash)
}

func (s *storageImpl) FetchRollups(height int) []*Rollup {
	s.assertSecretAvailable()
	return s.db.FetchRollups(height)
}

func (s *storageImpl) FetchState(hash obscurocommon.L2RootHash) State {
	s.assertSecretAvailable()
	return s.db.FetchState(hash)
}

func (s *storageImpl) StoreTx(tx L2Tx) {
	s.assertSecretAvailable()
	s.db.StoreTx(tx)
}

func (s *storageImpl) FetchTxs() []L2Tx {
	s.assertSecretAvailable()
	return s.db.FetchTxs()
}

func (s *storageImpl) PruneTxs(toRemove map[common.Hash]common.Hash) {
	s.assertSecretAvailable()
	s.db.PruneTxs(toRemove)
}

func (s *storageImpl) StoreBlock(b *types.Block) {
	s.assertSecretAvailable()

	var height int
	if b.ParentHash() == obscurocommon.GenesisHash {
		height = 0
	} else {
		bAndHeight, f := s.db.FetchBlockAndHeight(b.ParentHash())
		if !f {
			panic("Should not happen")
		}
		height = bAndHeight.height + 1
	}

	s.db.StoreBlock(b, height)
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

func (s *storageImpl) FetchRollupTxs(r *Rollup) (map[common.Hash]L2Tx, bool) {
	s.assertSecretAvailable()
	return s.db.FetchRollupTxs(r)
}

func (s *storageImpl) AddRollupTxs(r *Rollup, newMap map[common.Hash]L2Tx) {
	s.assertSecretAvailable()
	s.db.AddRollupTxs(r, newMap)
}

func (s *storageImpl) StoreSecret(secret SharedEnclaveSecret) {
	s.db.StoreSecret(secret)
}

func (s *storageImpl) FetchSecret() SharedEnclaveSecret {
	secret, err := s.secretsDB.Get([]byte("currentSecret"))
	if err != nil {
		panic(err)
	}
	return secret
}

func (s *storageImpl) HeightBlock(block *types.Block) int {
	s.assertSecretAvailable()
	val, f := s.db.FetchBlockAndHeight(block.Hash())
	if !f {
		panic("should not happen")
	}
	return val.height
}

func (s *storageImpl) ParentRollup(r *Rollup) *Rollup {
	s.assertSecretAvailable()
	parent, found := s.db.FetchRollup(r.Header.ParentHash)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", r.Hash()))
	}
	return parent
}

func (s *storageImpl) HeightRollup(r *Rollup) int {
	s.assertSecretAvailable()
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(obscurocommon.L2GenesisHeight)
		return obscurocommon.L2GenesisHeight
	}
	v := s.HeightRollup(s.ParentRollup(r)) + 1
	r.Height.Store(v)
	return v
}

func (s *storageImpl) Parent(b *types.Block) (*types.Block, bool) {
	s.assertSecretAvailable()
	return s.FetchBlock(b.Header().ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func (s *storageImpl) IsAncestor(blockA *types.Block, blockB *types.Block) bool {
	s.assertSecretAvailable()
	if blockA.Hash() == blockB.Hash() {
		return true
	}

	if s.HeightBlock(blockA) >= s.HeightBlock(blockB) {
		return false
	}

	p, f := s.Parent(blockB)
	if !f {
		return false
	}

	return s.IsAncestor(blockA, p)
}

// IsBlockAncestor - takes into consideration that the block to verify might be on a branch we haven't received yet
func (s *storageImpl) IsBlockAncestor(l1BlockHash obscurocommon.L1RootHash, block *types.Block) bool {
	s.assertSecretAvailable()
	if l1BlockHash == block.Hash() {
		return true
	}

	if l1BlockHash == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if s.HeightBlock(block) == 0 {
		return false
	}

	resolvedBlock, found := s.FetchBlock(l1BlockHash)
	if found {
		if s.HeightBlock(resolvedBlock) >= s.HeightBlock(block) {
			return false
		}
	}

	p, f := s.Parent(block)
	if !f {
		return false
	}

	return s.IsBlockAncestor(l1BlockHash, p)
}

func (s *storageImpl) assertSecretAvailable() {
	if s.FetchSecret() == nil {
		panic("Enclave not initialized")
	}
}
