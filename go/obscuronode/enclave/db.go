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
	// FetchBlockAndHeight returns the L1 block with the given hash, its height and true, or (nil, false) if no such block is stored
	FetchBlockAndHeight(hash obscurocommon.L1RootHash) (*blockAndHeight, bool)
	// StoreBlock persists the L1 block and its height in the chain
	StoreBlock(b *types.Block, height int)
	// FetchHeadBlock returns the L1 block at the head of the chain
	FetchHeadBlock() obscurocommon.L1RootHash

	// FetchRollup returns the rollup with the given hash and true, or (nil, false) if no such rollup is stored
	FetchRollup(hash obscurocommon.L2RootHash) (*Rollup, bool)
	// FetchRollups returns all the rollup with the given height
	FetchRollups(height int) []*Rollup
	// StoreRollup persists the rollup
	StoreRollup(rollup *Rollup, height int)

	// FetchBlockState returns the state after ingesting the L1 block with the given hash
	FetchBlockState(hash obscurocommon.L1RootHash) (blockState, bool)
	// SetBlockState persists the state after ingesting the L1 block with the given hash
	SetBlockState(hash obscurocommon.L1RootHash, state blockState)
	// SetBlockStateNewRollup persists the state after ingesting the L1 block with the given hash that contains a new rollup
	SetBlockStateNewRollup(hash obscurocommon.L1RootHash, state blockState)
	// FetchRollupState returns the state after adding the rollup with the given hash
	FetchRollupState(hash obscurocommon.L2RootHash) State
	// SetRollupState persists the state after adding the rollup with the given hash
	SetRollupState(hash obscurocommon.L2RootHash, state State)

	// FetchMempoolTxs returns all L2 transactions in the mempool
	FetchMempoolTxs() []L2Tx
	// AddMempoolTx adds an L2 transaction to the mempool
	AddMempoolTx(tx L2Tx)
	// RemoveMempoolTxs removes any L2 transactions whose hash is keyed in the map from the mempool
	RemoveMempoolTxs(remove map[common.Hash]common.Hash)
	// FetchRollupTxs returns all transactions in a given rollup keyed by hash and true, or (nil, false) if the rollup is unknown
	FetchRollupTxs(r *Rollup) (map[common.Hash]L2Tx, bool)
	// StoreRollupTxs overwrites the transactions associated with a given rollup
	StoreRollupTxs(*Rollup, map[common.Hash]L2Tx)

	// FetchSecret returns the enclave's secret
	FetchSecret() SharedEnclaveSecret
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret SharedEnclaveSecret)
}

type blockAndHeight struct {
	b      *types.Block
	height int
}

type inMemoryDB struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[obscurocommon.L1RootHash]blockState
	statePerRollup map[obscurocommon.L2RootHash]State
	headBlock      obscurocommon.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*Rollup
	rollups         map[obscurocommon.L2RootHash]*Rollup

	mempool map[common.Hash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[obscurocommon.L1RootHash]*blockAndHeight
	blockMutex sync.RWMutex

	txsPerBlockCache map[obscurocommon.L2RootHash]map[common.Hash]L2Tx
	txMutex          sync.RWMutex

	sharedEnclaveSecret SharedEnclaveSecret
}

func NewInMemoryDB() DB {
	return &inMemoryDB{
		statePerBlock:    make(map[obscurocommon.L1RootHash]blockState),
		stateMutex:       sync.RWMutex{},
		rollupsByHeight:  make(map[int][]*Rollup),
		rollups:          make(map[obscurocommon.L2RootHash]*Rollup),
		mempool:          make(map[common.Hash]L2Tx),
		mpMutex:          sync.RWMutex{},
		statePerRollup:   make(map[obscurocommon.L2RootHash]State),
		blockCache:       map[obscurocommon.L1RootHash]*blockAndHeight{},
		blockMutex:       sync.RWMutex{},
		txsPerBlockCache: make(map[obscurocommon.L2RootHash]map[common.Hash]L2Tx),
		txMutex:          sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchBlockState(hash obscurocommon.L1RootHash) (blockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetBlockState(hash obscurocommon.L1RootHash, state blockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *inMemoryDB) SetBlockStateNewRollup(hash obscurocommon.L1RootHash, state blockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerBlock[hash] = state
	db.rollups[state.head.Hash()] = state.head
	db.statePerRollup[state.head.Hash()] = state.state
	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *inMemoryDB) SetRollupState(hash obscurocommon.L2RootHash, state State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()

	db.statePerRollup[hash] = state
}

func (db *inMemoryDB) FetchHeadBlock() obscurocommon.L1RootHash {
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

func (db *inMemoryDB) FetchRollupState(hash obscurocommon.L2RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()

	return db.statePerRollup[hash]
}

func (db *inMemoryDB) AddMempoolTx(tx L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	db.mempool[tx.Hash()] = tx
}

func (db *inMemoryDB) FetchMempoolTxs() []L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()

	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *inMemoryDB) RemoveMempoolTxs(toRemove map[common.Hash]common.Hash) {
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

	val, found := db.txsPerBlockCache[r.Hash()]
	return val, found
}

func (db *inMemoryDB) StoreRollupTxs(r *Rollup, newMap map[common.Hash]L2Tx) {
	db.txMutex.Lock()
	defer db.txMutex.Unlock()

	db.txsPerBlockCache[r.Hash()] = newMap
}

func (db *inMemoryDB) StoreSecret(secret SharedEnclaveSecret) {
	db.sharedEnclaveSecret = secret
}

func (db *inMemoryDB) FetchSecret() SharedEnclaveSecret {
	return db.sharedEnclaveSecret
}
