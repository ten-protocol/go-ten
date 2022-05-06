package db

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type storageImpl struct {
	db *InMemoryDB
}

func (s *storageImpl) StoreGenesisRollup(rol *core.Rollup) {
	s.db.StoreGenesisRollup(rol)
}

func (s *storageImpl) FetchGenesisRollup() *core.Rollup {
	return s.db.FetchGenesisRollup()
}

func NewStorage(db *InMemoryDB) Storage {
	return &storageImpl{db: db}
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

func (s *storageImpl) StoreBlock(b *types.Block) bool {
	s.assertSecretAvailable()
	s.db.StoreBlock(b)
	return true
}

func (s *storageImpl) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	s.assertSecretAvailable()
	b, f := s.db.FetchBlock(hash)
	return b, f
}

func (s *storageImpl) FetchHeadBlock() *types.Block {
	s.assertSecretAvailable()
	b, _ := s.db.FetchBlock(s.db.FetchHeadBlock())
	return b
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

	if maybeAncestor.NumberU64() >= block.NumberU64() {
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

	if block.NumberU64() == obscurocommon.L1GenesisHeight {
		return false
	}

	resolvedBlock, found := s.FetchBlock(maybeAncestor)
	if found {
		if resolvedBlock.NumberU64() >= block.NumberU64() {
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
	return int64(v.NumberU64())
}

func (s *storageImpl) Proof(r *core.Rollup) *types.Block {
	v, f := s.FetchBlock(r.Header.L1Proof)
	if !f {
		panic("Could not find proof for this rollup")
	}
	return v
}

func (s *storageImpl) FetchBlockState(hash obscurocommon.L1RootHash) (*BlockState, bool) {
	return s.db.FetchBlockState(hash)
}

func (s *storageImpl) SetBlockState(hash obscurocommon.L1RootHash, state *BlockState) {
	if state.FoundNewRollup {
		s.db.SetBlockStateNewRollup(hash, state)
	} else {
		s.db.SetBlockState(hash, state)
	}
}

func (s *storageImpl) CreateStateDB(hash obscurocommon.L2RootHash) StateDB {
	parent := s.db.FetchRollupState(hash)
	newState := CopyStateNoWithdrawals(parent)
	return NewStateDB(s.db, hash, newState)
}

func (s *storageImpl) GenesisStateDB() StateDB {
	return NewStateDB(s.db, obscurocommon.GenesisHash, EmptyState())
}

func (s *storageImpl) FetchHeadState() *BlockState {
	val, _ := s.db.FetchBlockState(s.db.FetchHeadBlock())
	return val
}
