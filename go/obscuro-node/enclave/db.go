package enclave

import (
	common2 "github.com/obscuronet/obscuro-playground/go/common"
	"sync"
)

// RollupResolver -database of rollups indexed by the root hash
type RollupResolver interface {
	FetchRollup(hash common2.L2RootHash) *EnclaveRollup
	Parent(r EnclaveRollup) *EnclaveRollup
	Height(r EnclaveRollup) int
}

// This database lives purely in the memory space of an encrypted enclave
type Db interface {
	FetchState(hash common2.L1RootHash) (BlockState, bool)
	SetState(hash common2.L1RootHash, state BlockState)
	FetchRollup(hash common2.L2RootHash) *EnclaveRollup
	FetchRollupState(hash common2.L2RootHash) State
	SetRollupState(hash common2.L2RootHash, state State)
	Head() BlockState
	Balance(address common2.Address) uint64
	FetchRollups(height int) []*EnclaveRollup
	StoreRollup(rollup *EnclaveRollup)
	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[common2.TxHash]common2.TxHash)
	Resolve(hash common2.L1RootHash) (*common2.Block, bool)
	Store(node *common2.Block)
	Txs(r *EnclaveRollup) (map[common2.TxHash]L2Tx, bool)
	AddTxs(*EnclaveRollup, map[common2.TxHash]L2Tx)
	Height(*EnclaveRollup) int
	Parent(*EnclaveRollup) *EnclaveRollup
}

type inMemoryDb struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[common2.L1RootHash]BlockState
	statePerRollup map[common2.L2RootHash]State
	headBlock      common2.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*EnclaveRollup // todo encoded
	rollups         map[common2.L2RootHash]*EnclaveRollup

	mempool map[common2.TxHash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[common2.L1RootHash]*common2.Block
	blockM     sync.RWMutex

	transactionsPerBlockCache map[common2.L2RootHash]map[common2.TxHash]L2Tx
	txM                       sync.RWMutex
}

func NewInMemoryDb() Db {
	return &inMemoryDb{
		statePerBlock:             make(map[common2.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[int][]*EnclaveRollup),
		rollups:                   make(map[common2.L2RootHash]*EnclaveRollup),
		mempool:                   make(map[common2.TxHash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[common2.L2RootHash]State),
		blockCache:                map[common2.L1RootHash]*common2.Block{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[common2.L2RootHash]map[common2.TxHash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *inMemoryDb) FetchState(hash common2.L1RootHash) (BlockState, bool) {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDb) SetState(hash common2.L1RootHash, state BlockState) {
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

func (db *inMemoryDb) SetRollupState(hash common2.L2RootHash, state State) {
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerRollup[hash] = state
}

func (db *inMemoryDb) Head() BlockState {
	val, _ := db.FetchState(db.headBlock)
	return val
}

func (db *inMemoryDb) Balance(address common2.Address) uint64 {
	return db.Head().State[address]
}

func (db *inMemoryDb) StoreRollup(rollup *EnclaveRollup) {
	height := db.Height(rollup)

	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.rollups[rollup.Hash()] = rollup
	val, found := db.rollupsByHeight[height]
	if found {
		db.rollupsByHeight[height] = append(val, rollup)
	} else {
		db.rollupsByHeight[height] = []*EnclaveRollup{rollup}
	}
}

func (db *inMemoryDb) FetchRollup(hash common2.L2RootHash) *EnclaveRollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	r, f := db.rollups[hash]
	if !f {
		panic("wtf")
	}
	return r
}

func (db *inMemoryDb) FetchRollups(height int) []*EnclaveRollup {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *inMemoryDb) FetchRollupState(hash common2.L2RootHash) State {
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *inMemoryDb) StoreTx(tx L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.Id] = tx
}

func (db *inMemoryDb) FetchTxs() []L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()
	// txStr := make([]common.TxHash, 0)
	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
		// txStr = append(txStr, tx.Id)
	}
	// common.Log(fmt.Sprintf(">>> %v <<<", txStr))
	return mpCopy
}

func (db *inMemoryDb) PruneTxs(toRemove map[common2.TxHash]common2.TxHash) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	r := make(map[common2.TxHash]L2Tx, 0)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	// fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}

func (db *inMemoryDb) Store(node *common2.Block) {
	db.blockM.Lock()
	db.blockCache[node.Hash()] = node
	db.blockM.Unlock()
}

func (db *inMemoryDb) Resolve(hash common2.L1RootHash) (*common2.Block, bool) {
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	v, f := db.blockCache[hash]
	return v, f
}

func (db *inMemoryDb) Txs(r *EnclaveRollup) (map[common2.TxHash]L2Tx, bool) {
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[r.Hash()]
	db.txM.RUnlock()
	return val, found
}

func (db *inMemoryDb) AddTxs(r *EnclaveRollup, newMap map[common2.TxHash]L2Tx) {
	db.txM.Lock()
	db.transactionsPerBlockCache[r.Hash()] = newMap
	db.txM.Unlock()
}

func (db *inMemoryDb) Parent(r *EnclaveRollup) *EnclaveRollup {
	return db.FetchRollup(r.Header.ParentHash)
}

func (db *inMemoryDb) Height(r *EnclaveRollup) int {
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
