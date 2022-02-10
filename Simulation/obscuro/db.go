package obscuro

import (
	"simulation/common"
	"sync"
)

// This database lives purely in the memory space of an encrypted enclave
type Db interface {
	FetchState(hash common.L1RootHash) (BlockState, bool)
	SetState(hash common.L1RootHash, state BlockState)

	FetchRollup(hash common.L2RootHash) Rollup
	FetchRollupState(hash common.L2RootHash) State
	SetRollupState(hash common.L2RootHash, state State)

	Head() BlockState
	Balance(address common.Address) uint64

	FetchRollups(height uint32) []Rollup
	StoreRollup(height uint32, rollup Rollup)

	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[common.TxHash]common.TxHash)

	Resolve(hash common.L1RootHash) (common.Block, bool)
	Store(node common.Block)

	Txs(block Rollup) (map[common.TxHash]L2Tx, bool)
	AddTxs(Rollup, map[common.TxHash]L2Tx)
}

type InMemoryDb struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[common.L1RootHash]BlockState
	statePerRollup map[common.L2RootHash]State
	headBlock      common.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[uint32][]Rollup // todo encoded
	rollups         map[common.L2RootHash]Rollup

	mempool map[common.TxHash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[common.L1RootHash]common.Block
	blockM     sync.RWMutex

	transactionsPerBlockCache map[common.L2RootHash]map[common.TxHash]L2Tx
	txM                       sync.RWMutex
}

func NewInMemoryDb() Db {
	return &InMemoryDb{
		statePerBlock:             make(map[common.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[uint32][]Rollup),
		rollups:                   make(map[common.TxHash]Rollup),
		mempool:                   make(map[common.TxHash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[common.L2RootHash]State),
		blockCache:                map[common.L1RootHash]common.Block{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[common.L2RootHash]map[common.TxHash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *InMemoryDb) FetchState(hash common.L1RootHash) (BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *InMemoryDb) SetState(hash common.L1RootHash, state BlockState) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerBlock[hash] = state
	// todo - move this logic outside
	if state.foundNewRollup {
		db.rollups[state.Head.RootHash] = state.Head
		db.statePerRollup[state.Head.RootHash] = state.State
	}

	// todo - is there any other logic here?
	db.headBlock = hash
}

func (db *InMemoryDb) SetRollupState(hash common.L2RootHash, state State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerRollup[hash] = state
}

func (db *InMemoryDb) Head() BlockState {
	val, _ := db.FetchState(db.headBlock)
	return val
}

func (db *InMemoryDb) Balance(address common.Address) uint64 {
	return db.Head().State[address]
}

func (db *InMemoryDb) StoreRollup(height uint32, rollup Rollup) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.rollups[rollup.RootHash] = rollup
	val, found := db.rollupsByHeight[height]
	if found {
		db.rollupsByHeight[height] = append(val, rollup)
	} else {
		db.rollupsByHeight[height] = []Rollup{rollup}
	}
}
func (db *InMemoryDb) FetchRollup(hash common.L2RootHash) Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollups[hash]
}

func (db *InMemoryDb) FetchRollups(height uint32) []Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *InMemoryDb) FetchRollupState(hash common.L2RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *InMemoryDb) StoreTx(tx L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.Id] = tx
}

func (db *InMemoryDb) FetchTxs() []L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()
	//txStr := make([]common.TxHash, 0)
	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
		//txStr = append(txStr, tx.Id)
	}
	//common.Log(fmt.Sprintf(">>> %v <<<", txStr))
	return mpCopy
}

func (db *InMemoryDb) PruneTxs(toRemove map[common.TxHash]common.TxHash) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	r := make(map[common.TxHash]L2Tx, 0)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	//fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}

func (db *InMemoryDb) Store(node common.Block) {
	db.blockM.Lock()
	db.blockCache[node.Hash()] = node
	db.blockM.Unlock()
}

func (db *InMemoryDb) Resolve(hash common.L1RootHash) (common.Block, bool) {
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	v, f := db.blockCache[hash]
	return v, f
}

func (db *InMemoryDb) Txs(b Rollup) (map[common.TxHash]L2Tx, bool) {
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[b.RootHash]
	db.txM.RUnlock()
	return val, found
}

func (db *InMemoryDb) AddTxs(b Rollup, newMap map[common.TxHash]L2Tx) {
	db.txM.Lock()
	db.transactionsPerBlockCache[b.RootHash] = newMap
	db.txM.Unlock()
}
