package obscuro

import (
	"simulation/common"
	"sync"
)

// This database lives purely in the memory space of an encrypted enclave
type Db interface {
	FetchState(hash common.RootHash) (BlockState, bool)
	SetState(hash common.RootHash, state BlockState)

	FetchRollup(hash common.RootHash) common.Rollup
	FetchRollupState(hash common.RootHash) State
	SetRollupState(hash common.RootHash, state State)

	Head() BlockState
	Balance(address common.Address) uint64

	FetchRollups(height uint32) []common.Rollup
	StoreRollup(height uint32, rollup common.Rollup)

	FetchTxs() []common.L2Tx
	StoreTx(tx common.L2Tx)
	PruneTxs(remove map[common.TxHash]common.TxHash)
}

type InMemoryDb struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[common.RootHash]BlockState
	statePerRollup map[common.RootHash]State
	headBlock      common.RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[uint32][]common.Rollup // todo encoded
	rollups         map[common.RootHash]common.Rollup

	mempool map[common.TxHash]common.L2Tx
	mpMutex sync.RWMutex

	speculativeRollup currentWork
}

func NewInMemoryDb() Db {
	return &InMemoryDb{
		statePerBlock:   make(map[common.RootHash]BlockState),
		stateMutex:      sync.RWMutex{},
		rollupsByHeight: make(map[uint32][]common.Rollup),
		rollups:         make(map[common.TxHash]common.Rollup),
		mpMutex:         sync.RWMutex{},
		statePerRollup:  make(map[common.RootHash]State),
	}
}

func (db *InMemoryDb) FetchState(hash common.RootHash) (BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *InMemoryDb) SetState(hash common.RootHash, state BlockState) {
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

func (db *InMemoryDb) SetRollupState(hash common.RootHash, state State) {
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

func (db *InMemoryDb) StoreRollup(height uint32, rollup common.Rollup) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.rollups[rollup.RootHash] = rollup
	val, found := db.rollupsByHeight[height]
	if found {
		db.rollupsByHeight[height] = append(val, rollup)
	} else {
		db.rollupsByHeight[height] = []common.Rollup{rollup}
	}
}
func (db *InMemoryDb) FetchRollup(hash common.RootHash) common.Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollups[hash]
}

func (db *InMemoryDb) FetchRollups(height uint32) []common.Rollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *InMemoryDb) FetchRollupState(hash common.RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *InMemoryDb) StoreTx(tx common.L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.Id] = tx
}

func (db *InMemoryDb) FetchTxs() []common.L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()
	//txStr := make([]common.TxHash, 0)
	mpCopy := make([]common.L2Tx, 0)
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
	r := make(map[common.TxHash]common.L2Tx, 0)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	//fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}
