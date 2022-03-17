package enclave

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// RollupResolver -database of rollups indexed by the root hash
type RollupResolver interface {
	FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool)
	StoreRollup(rollup *Rollup)
}

// BlockResolver -database of blocks indexed by the root hash
type BlockResolver interface {
	FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool)
	StoreBlock(block *types.Block)
	HeightBlock(block *types.Block) int
}

// This database lives purely in the memory space of an encrypted enclave
type DB interface {
	// Rollup Resolver
	RollupResolver

	// Block resolver
	BlockResolver

	// Gossip
	FetchRollups(height int) []*Rollup

	// State
	FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool)
	SetBlockState(hash obscurocommon.L1RootHash, state BlockState)
	FetchRollupState(hash obscurocommon.L2RootHash) State
	SetRollupState(hash obscurocommon.L2RootHash, state State)
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

type blockAndHeight struct {
	b      *types.Block
	height int
}

type inMemoryDB struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[obscurocommon.L1RootHash]BlockState
	statePerRollup map[obscurocommon.L2RootHash]State
	headBlock      obscurocommon.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*Rollup
	rollups         map[obscurocommon.L2RootHash]*Rollup

	mempool map[common.Hash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[obscurocommon.L1RootHash]*blockAndHeight
	blockM     sync.RWMutex

	transactionsPerBlockCache map[obscurocommon.L2RootHash]map[common.Hash]L2Tx
	txM                       sync.RWMutex

	sharedEnclaveSecret SharedEnclaveSecret
}

func NewInMemoryDB() DB {
	return &inMemoryDB{
		statePerBlock:             make(map[obscurocommon.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[int][]*Rollup),
		rollups:                   make(map[obscurocommon.L2RootHash]*Rollup),
		mempool:                   make(map[common.Hash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[obscurocommon.L2RootHash]State),
		blockCache:                map[obscurocommon.L1RootHash]*blockAndHeight{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[obscurocommon.L2RootHash]map[common.Hash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool) {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetBlockState(hash obscurocommon.L1RootHash, state BlockState) {
	db.assertSecretAvailable()
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerBlock[hash] = state
	// todo - move this logic outside
	if state.foundNewRollup {
		db.rollups[state.Head.Hash()] = state.Head
		db.statePerRollup[state.Head.Hash()] = state.State
	}

	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *inMemoryDB) SetRollupState(hash obscurocommon.L2RootHash, state State) {
	db.assertSecretAvailable()
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerRollup[hash] = state
}

func (db *inMemoryDB) Head() BlockState {
	db.assertSecretAvailable()
	val, _ := db.FetchBlockState(db.headBlock)
	return val
}

func (db *inMemoryDB) StoreRollup(rollup *Rollup) {
	db.assertSecretAvailable()
	height := HeightRollup(db, rollup)

	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.rollups[rollup.Hash()] = rollup
	val, found := db.rollupsByHeight[height]
	if found {
		db.rollupsByHeight[height] = append(val, rollup)
	} else {
		db.rollupsByHeight[height] = []*Rollup{rollup}
	}
}

func (db *inMemoryDB) FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool) {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	r, f := db.rollups[hash]
	return r, f
}

func (db *inMemoryDB) FetchRollups(height int) []*Rollup {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *inMemoryDB) FetchRollupState(hash obscurocommon.L2RootHash) State {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *inMemoryDB) StoreTx(tx L2Tx) {
	db.assertSecretAvailable()
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.Hash()] = tx
}

func (db *inMemoryDB) FetchTxs() []L2Tx {
	db.assertSecretAvailable()
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()
	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *inMemoryDB) PruneTxs(toRemove map[common.Hash]common.Hash) {
	db.assertSecretAvailable()
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	r := make(map[common.Hash]L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	db.mempool = r
}

func (db *inMemoryDB) StoreBlock(b *types.Block) {
	db.assertSecretAvailable()
	db.blockM.Lock()
	defer db.blockM.Unlock()

	if b.ParentHash() == obscurocommon.GenesisHash {
		db.blockCache[b.Hash()] = &blockAndHeight{b, 0}
		return
	}

	p, f := db.blockCache[b.ParentHash()]
	if !f {
		panic("Should not happen")
	}
	db.blockCache[b.Hash()] = &blockAndHeight{b: b, height: p.height + 1}
}

func (db *inMemoryDB) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	db.assertSecretAvailable()
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	val, f := db.blockCache[hash]
	var block *types.Block
	if f {
		block = val.b
	}
	return block, f
}

func (db *inMemoryDB) FetchRollupTxs(r *Rollup) (map[common.Hash]L2Tx, bool) {
	db.assertSecretAvailable()
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[r.Hash()]
	db.txM.RUnlock()
	return val, found
}

func (db *inMemoryDB) AddRollupTxs(r *Rollup, newMap map[common.Hash]L2Tx) {
	db.assertSecretAvailable()
	db.txM.Lock()
	db.transactionsPerBlockCache[r.Hash()] = newMap
	db.txM.Unlock()
}

func (db *inMemoryDB) StoreSecret(secret SharedEnclaveSecret) {
	db.sharedEnclaveSecret = secret
}

func (db *inMemoryDB) FetchSecret() SharedEnclaveSecret {
	return db.sharedEnclaveSecret
}

func (db *inMemoryDB) HeightBlock(block *types.Block) int {
	db.assertSecretAvailable()
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	b, f := db.blockCache[block.Hash()]
	if !f {
		panic("should not happen")
	}
	return b.height
}

func (db *inMemoryDB) assertSecretAvailable() {
	if db.FetchSecret() == nil {
		panic("Enclave not initialized")
	}
}
