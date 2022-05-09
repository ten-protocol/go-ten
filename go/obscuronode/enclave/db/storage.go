package db

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	obscurorawdb "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rawdb"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type storageImpl struct {
	tempDB *InMemoryDB // todo - has to be replaced completely by the ethdb.Database
	db     ethdb.Database
}

func NewStorage(db *InMemoryDB) Storage {
	return &storageImpl{tempDB: db, db: rawdb.NewMemoryDatabase()}
}

func (s *storageImpl) StoreGenesisRollup(rol *core.Rollup) {
	obscurorawdb.WriteGenesisHash(s.db, rol.Hash())
	s.StoreRollup(rol)
}

func (s *storageImpl) FetchGenesisRollup() *core.Rollup {
	hash := obscurorawdb.ReadGenesisHash(s.db)
	if hash == nil {
		return nil
	}
	r, _ := s.FetchRollup(*hash)
	return r
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) {
	s.assertSecretAvailable()

	batch := s.db.NewBatch()
	obscurorawdb.WriteRollup(batch, rollup)
	if err := batch.Write(); err != nil {
		panic(err)
	}
}

func (s *storageImpl) FetchRollup(hash obscurocommon.L2RootHash) (*core.Rollup, bool) {
	s.assertSecretAvailable()
	r := obscurorawdb.ReadRollup(s.db, hash)
	if r != nil {
		return r, true
	}
	return nil, false
}

func (s *storageImpl) FetchRollups(height uint64) []*core.Rollup {
	s.assertSecretAvailable()
	return obscurorawdb.ReadRollupsForHeight(s.db, height)
}

func (s *storageImpl) StoreBlock(b *types.Block) bool {
	s.assertSecretAvailable()
	rawdb.WriteBlock(s.db, b)
	return true
}

func (s *storageImpl) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	s.assertSecretAvailable()
	height := rawdb.ReadHeaderNumber(s.db, hash)
	if height == nil {
		return nil, false
	}
	return rawdb.ReadBlock(s.db, hash, *height), true
}

func (s *storageImpl) FetchHeadBlock() *types.Block {
	s.assertSecretAvailable()
	b, _ := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	return b
}

func (s *storageImpl) FetchRollupTxs(r *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool) {
	s.assertSecretAvailable()
	return s.tempDB.FetchRollupTxs(r)
}

func (s *storageImpl) StoreRollupTxs(r *core.Rollup, newTxs map[common.Hash]nodecommon.L2Tx) {
	s.assertSecretAvailable()
	s.tempDB.StoreRollupTxs(r, newTxs)
}

func (s *storageImpl) StoreSecret(secret core.SharedEnclaveSecret) {
	obscurorawdb.WriteSharedSecret(s.db, secret)
}

func (s *storageImpl) FetchSecret() core.SharedEnclaveSecret {
	ss := obscurorawdb.ReadSharedSecret(s.db)
	if ss != nil {
		return *ss
	}
	// todo - I guess this is fixed by Matt
	return core.SharedEnclaveSecret{}
}

func (s *storageImpl) ParentRollup(r *core.Rollup) *core.Rollup {
	s.assertSecretAvailable()
	parent, found := s.FetchRollup(r.Header.ParentHash)
	if !found {
		log.Log(fmt.Sprintf("Could not find rollup: r_%d", obscurocommon.ShortHash(r.Hash())))
		return nil
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
	return s.tempDB.FetchBlockState(hash)
}

func (s *storageImpl) SetBlockState(hash obscurocommon.L1RootHash, state *BlockState, rollup *core.Rollup) {
	if state.FoundNewRollup {
		s.StoreRollup(rollup)
	}
	s.tempDB.SetBlockState(hash, state)
	rawdb.WriteHeadHeaderHash(s.db, state.Block)
}

func (s *storageImpl) CreateStateDB(hash obscurocommon.L2RootHash) StateDB {
	parent := s.tempDB.FetchRollupState(hash)
	newState := CopyStateNoWithdrawals(parent)
	return NewStateDB(s.tempDB, hash, newState)
}

func (s *storageImpl) GenesisStateDB() StateDB {
	return NewStateDB(s.tempDB, obscurocommon.GenesisHash, EmptyState())
}

func (s *storageImpl) FetchHeadState() *BlockState {
	val, _ := s.tempDB.FetchBlockState(rawdb.ReadHeadHeaderHash(s.db))
	return val
}
