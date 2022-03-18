package enclave

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// TODO - Further generify this interface's methods.

// DB lives purely in the encrypted memory space of an enclave.
// Unlike Storage, methods in this class should have minimal logic, to map them more easily to our chosen datastore.
type DB interface {
	// Blocks
	FetchBlockAndHeight(hash obscurocommon.L1RootHash) (*blockAndHeight, bool)
	StoreBlock(b *types.Block, height int)
	HeadBlock() obscurocommon.L1RootHash

	// Rollups
	FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool)
	StoreRollup(rollup *Rollup, height int)
	FetchRollups(height int) []*Rollup

	// State
	FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool)
	SetBlockState(hash obscurocommon.L1RootHash, state BlockState)
	SetBlockStateNewRollup(hash obscurocommon.L1RootHash, state BlockState)
	FetchState(hash obscurocommon.L2RootHash) State
	SetState(hash obscurocommon.L2RootHash, state State)

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
	blockMutex sync.RWMutex

	transactionsPerBlockCache map[obscurocommon.L2RootHash]map[common.Hash]L2Tx
	txMutex                   sync.RWMutex

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
		blockMutex:                sync.RWMutex{},
		transactionsPerBlockCache: make(map[obscurocommon.L2RootHash]map[common.Hash]L2Tx),
		txMutex:                   sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchBlockState(hash obscurocommon.L1RootHash) (BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetBlockState(hash obscurocommon.L1RootHash, state BlockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *inMemoryDB) SetBlockStateNewRollup(hash obscurocommon.L1RootHash, state BlockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	db.rollups[state.Head.Hash()] = state.Head
	db.statePerRollup[state.Head.Hash()] = state.State
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *inMemoryDB) SetState(hash obscurocommon.L2RootHash, state State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerRollup[hash] = state
}

func (db *inMemoryDB) HeadBlock() obscurocommon.L1RootHash {
	return db.headBlock
}

func (db *inMemoryDB) StoreRollup(rollup *Rollup, height int) {
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
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	r, f := db.rollups[hash]
	return r, f
}

func (db *inMemoryDB) FetchRollups(height int) []*Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	return db.rollupsByHeight[height]
}

func (db *inMemoryDB) FetchState(hash obscurocommon.L2RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	return db.statePerRollup[hash]
}

func (db *inMemoryDB) StoreTx(tx L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	db.mempool[tx.Hash()] = tx
}

func (db *inMemoryDB) FetchTxs() []L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()

	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *inMemoryDB) PruneTxs(toRemove map[common.Hash]common.Hash) {
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

func (db *inMemoryDB) StoreBlock(b *types.Block, height int) {
	db.blockMutex.Lock()
	defer db.blockMutex.Unlock()

	db.blockCache[b.Hash()] = &blockAndHeight{b: b, height: height}
}

func (db *inMemoryDB) FetchBlockAndHeight(hash obscurocommon.L1RootHash) (*blockAndHeight, bool) {
	db.blockMutex.RLock()
	defer db.blockMutex.RUnlock()

	val, f := db.blockCache[hash]
	return val, f
}

func (db *inMemoryDB) FetchRollupTxs(r *Rollup) (map[common.Hash]L2Tx, bool) {
	db.txMutex.RLock()
	defer db.txMutex.RUnlock()

	val, found := db.transactionsPerBlockCache[r.Hash()]
	return val, found
}

func (db *inMemoryDB) AddRollupTxs(r *Rollup, newMap map[common.Hash]L2Tx) {
	db.txMutex.Lock()
	defer db.txMutex.Unlock()

	db.transactionsPerBlockCache[r.Hash()] = newMap
}

func (db *inMemoryDB) StoreSecret(secret SharedEnclaveSecret) {
	db.sharedEnclaveSecret = secret
}

func (db *inMemoryDB) FetchSecret() SharedEnclaveSecret {
	return db.sharedEnclaveSecret
}
