package enclave

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"

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
	Balance(address common.Address) *big.Int
	FetchRollups(height int) []*Rollup
	StoreRollup(rollup *Rollup)
	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[common.Hash]common.Hash)
	Resolve(hash common2.L1RootHash) (*common2.Block, bool)
	Store(node *common2.Block)
	Txs(r *Rollup) (map[common.Hash]L2Tx, bool)
	AddTxs(*Rollup, map[common.Hash]L2Tx)
	Height(*Rollup) int
	Parent(*Rollup) *Rollup
	StoreSecret(secret SharedEnclaveSecret)
	FetchSecret() SharedEnclaveSecret
}

type inMemoryDB struct {
	// the State is dependent on the L1 block alone
	statePerBlock  map[common2.L1RootHash]BlockState
	statePerRollup map[common2.L2RootHash]State
	headBlock      common2.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*Rollup
	rollups         map[common2.L2RootHash]*Rollup

	mempool map[common.Hash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[common2.L1RootHash]*common2.Block
	blockM     sync.RWMutex

	transactionsPerBlockCache map[common2.L2RootHash]map[common.Hash]L2Tx
	txM                       sync.RWMutex

	sharedEnclaveSecret SharedEnclaveSecret
}

func NewInMemoryDB() DB {
	return &inMemoryDB{
		statePerBlock:             make(map[common2.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[int][]*Rollup),
		rollups:                   make(map[common2.L2RootHash]*Rollup),
		mempool:                   make(map[common.Hash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[common2.L2RootHash]State),
		blockCache:                map[common2.L1RootHash]*common2.Block{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[common2.L2RootHash]map[common.Hash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchState(hash common2.L1RootHash) (BlockState, bool) {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetState(hash common2.L1RootHash, state BlockState) {
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

func (db *inMemoryDB) SetRollupState(hash common2.L2RootHash, state State) {
	db.assertSecretAvailable()
	db.stateMutex.Lock()
	defer db.stateMutex.Unlock()
	db.statePerRollup[hash] = state
}

func (db *inMemoryDB) Head() BlockState {
	db.assertSecretAvailable()
	val, _ := db.FetchState(db.headBlock)
	return val
}

func (db *inMemoryDB) Balance(address common.Address) *big.Int {
	db.assertSecretAvailable()
	return db.Head().State[address]
}

func (db *inMemoryDB) StoreRollup(rollup *Rollup) {
	db.assertSecretAvailable()
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
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	r, f := db.rollups[hash]
	if !f {
		panic(fmt.Sprintf("Could not find rollup: r_%s", hash))
	}
	return r
}

func (db *inMemoryDB) ExistRollup(hash common2.L2RootHash) bool {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	_, f := db.rollups[hash]
	return f
}

func (db *inMemoryDB) FetchRollups(height int) []*Rollup {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *inMemoryDB) FetchRollupState(hash common2.L2RootHash) State {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.statePerRollup[hash]
}

func (db *inMemoryDB) StoreTx(tx L2Tx) {
	db.assertSecretAvailable()
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	db.mempool[tx.Tx.Hash()] = tx
}

func (db *inMemoryDB) FetchTxs() []L2Tx {
	db.assertSecretAvailable()
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
	// fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}

func (db *inMemoryDB) Store(node *common2.Block) {
	db.assertSecretAvailable()
	db.blockM.Lock()
	db.blockCache[node.Hash()] = node
	db.blockM.Unlock()
}

func (db *inMemoryDB) Resolve(hash common2.L1RootHash) (*common2.Block, bool) {
	db.assertSecretAvailable()
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	v, f := db.blockCache[hash]
	return v, f
}

func (db *inMemoryDB) Txs(r *Rollup) (map[common.Hash]L2Tx, bool) {
	db.assertSecretAvailable()
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[r.Hash()]
	db.txM.RUnlock()
	return val, found
}

func (db *inMemoryDB) AddTxs(r *Rollup, newMap map[common.Hash]L2Tx) {
	db.assertSecretAvailable()
	db.txM.Lock()
	db.transactionsPerBlockCache[r.Hash()] = newMap
	db.txM.Unlock()
}

func (db *inMemoryDB) Parent(r *Rollup) *Rollup {
	db.assertSecretAvailable()
	return db.FetchRollup(r.Header.ParentHash)
}

func (db *inMemoryDB) Height(r *Rollup) int {
	db.assertSecretAvailable()
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

func (db *inMemoryDB) StoreSecret(secret SharedEnclaveSecret) {
	db.sharedEnclaveSecret = secret
}

func (db *inMemoryDB) FetchSecret() SharedEnclaveSecret {
	return db.sharedEnclaveSecret
}

func (db *inMemoryDB) assertSecretAvailable() {
	if db.sharedEnclaveSecret == nil {
		panic("Enclave not initialized")
	}
}
