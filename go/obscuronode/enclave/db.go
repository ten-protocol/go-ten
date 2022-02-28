package enclave

import (
	"fmt"
	"sync"

	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

// RollupResolver -database of rollups indexed by the root hash
type RollupResolver interface {
	FetchRollup(hash common2.L2RootHash) *Rollup
	Parent(r Rollup) *Rollup
	Height(r Rollup) int
}

// This database lives purely in the memory space of an encrypted enclave
type DB interface {
	FetchState(hash common2.L1RootHash) (BlockState, bool)
	SetState(hash common2.L1RootHash, state BlockState)
	FetchRollup(hash common2.L2RootHash) *Rollup
	ExistRollup(hash common2.L2RootHash) bool
	FetchRollupState(hash common2.L2RootHash) State
	SetRollupState(hash common2.L2RootHash, state State)
	Head() BlockState
	Balance(address common2.Address) uint64
	FetchRollups(height int) []*Rollup
	StoreRollup(rollup *Rollup)
	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[common2.TxHash]common2.TxHash)
	Resolve(hash common2.L1RootHash) (*common2.Block, bool)
	Store(node *common2.Block)
	Txs(r *Rollup) (map[common2.TxHash]L2Tx, bool)
	AddTxs(*Rollup, map[common2.TxHash]L2Tx)
	Height(*Rollup) int
	Parent(*Rollup) *Rollup
}

type inMemoryDB struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[common2.L1RootHash]BlockState
	statePerRollup map[common2.L2RootHash]State
	headBlock      common2.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*Rollup
	rollups         map[common2.L2RootHash]*Rollup

	mempool map[common2.TxHash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[common2.L1RootHash]*common2.Block
	blockM     sync.RWMutex

	transactionsPerBlockCache map[common2.L2RootHash]map[common2.TxHash]L2Tx
	txM                       sync.RWMutex
}

func NewInMemoryDB() DB {
	return &inMemoryDB{
		statePerBlock:             make(map[common2.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[int][]*Rollup),
		rollups:                   make(map[common2.L2RootHash]*Rollup),
		mempool:                   make(map[common2.TxHash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[common2.L2RootHash]State),
		blockCache:                map[common2.L1RootHash]*common2.Block{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[common2.L2RootHash]map[common2.TxHash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchState(hash common2.L1RootHash) (BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetState(hash common2.L1RootHash, state BlockState) {
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

func (db *inMemoryDB) SetRollupState(hash common2.L2RootHash, state State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerRollup[hash] = state
}

func (db *inMemoryDB) Head() BlockState {
	val, _ := db.FetchState(db.headBlock)
	return val
}

func (db *inMemoryDB) Balance(address common2.Address) uint64 {
	return db.Head().State[address]
}

func (db *inMemoryDB) StoreRollup(rollup *Rollup) {
	height := db.Height(rollup)

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

func (db *inMemoryDB) FetchRollup(hash common2.L2RootHash) *Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	r, f := db.rollups[hash]
	if !f {
		panic(fmt.Sprintf("Could not find rollup: r_%s", hash))
	}
	return r
}

func (db *inMemoryDB) ExistRollup(hash common2.L2RootHash) bool {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	_, f := db.rollups[hash]
	return f
}

func (db *inMemoryDB) FetchRollups(height int) []*Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *inMemoryDB) FetchRollupState(hash common2.L2RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *inMemoryDB) StoreTx(tx L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.ID] = tx
}

func (db *inMemoryDB) FetchTxs() []L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()
	// txStr := make([]common.TxHash, 0)
	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
		// txStr = append(txStr, tx.ID)
	}
	// common.Log(fmt.Sprintf(">>> %v <<<", txStr))
	return mpCopy
}

func (db *inMemoryDB) PruneTxs(toRemove map[common2.TxHash]common2.TxHash) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	r := make(map[common2.TxHash]L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	// fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}

func (db *inMemoryDB) Store(node *common2.Block) {
	db.blockM.Lock()
	db.blockCache[node.Hash()] = node
	db.blockM.Unlock()
}

func (db *inMemoryDB) Resolve(hash common2.L1RootHash) (*common2.Block, bool) {
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	v, f := db.blockCache[hash]
	return v, f
}

func (db *inMemoryDB) Txs(r *Rollup) (map[common2.TxHash]L2Tx, bool) {
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[r.Hash()]
	db.txM.RUnlock()
	return val, found
}

func (db *inMemoryDB) AddTxs(r *Rollup, newMap map[common2.TxHash]L2Tx) {
	db.txM.Lock()
	db.transactionsPerBlockCache[r.Hash()] = newMap
	db.txM.Unlock()
}

func (db *inMemoryDB) Parent(r *Rollup) *Rollup {
	return db.FetchRollup(r.Header.ParentHash)
}

func (db *inMemoryDB) Height(r *Rollup) int {
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(common2.L2GenesisHeight)
		return common2.L2GenesisHeight
	}
	v := db.Height(db.Parent(r)) + 1
	r.Height.Store(v)
	return v
}
