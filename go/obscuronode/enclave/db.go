package enclave

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/common"
)

// RollupResolver -database of rollups indexed by the root hash
type RollupResolver interface {
	FetchRollup(hash common.L2RootHash) *Rollup
	ParentRollup(r Rollup) *Rollup
	HeightRollup(r Rollup) int
	StoreRollup(rollup *Rollup)
	ExistRollup(hash common.L2RootHash) bool
}

// This database lives purely in the memory space of an encrypted enclave
type DB interface {
	// Rollup Resolver
	FetchRollup(hash common.L2RootHash) *Rollup
	StoreRollup(rollup *Rollup)
	ParentRollup(*Rollup) *Rollup
	HeightRollup(*Rollup) int
	ExistRollup(hash common.L2RootHash) bool

	// Gossip
	FetchGossipedRollups(height int) []*Rollup

	// Block resolver
	HeightBlock(block *types.Block) int
	ParentBlock(block *types.Block) (*types.Block, bool)
	ResolveBlock(hash common.L1RootHash) (*types.Block, bool)
	StoreBlock(node *types.Block)

	// State
	FetchState(hash common.L1RootHash) (BlockState, bool)
	SetState(hash common.L1RootHash, state BlockState)
	FetchRollupState(hash common.L2RootHash) State
	SetRollupState(hash common.L2RootHash, state State)
	Head() BlockState
	Balance(address gethcommon.Address) uint64

	// Transactions
	FetchTxs() []L2Tx
	StoreTx(tx L2Tx)
	PruneTxs(remove map[gethcommon.Hash]gethcommon.Hash)
	Txs(r *Rollup) (map[gethcommon.Hash]L2Tx, bool)
	AddTxs(*Rollup, map[gethcommon.Hash]L2Tx)

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
	statePerBlock  map[common.L1RootHash]BlockState
	statePerRollup map[common.L2RootHash]State
	headBlock      common.L1RootHash
	stateMutex     sync.RWMutex

	rollupsByHeight map[int][]*Rollup
	rollups         map[common.L2RootHash]*Rollup

	mempool map[gethcommon.Hash]L2Tx
	mpMutex sync.RWMutex

	blockCache map[common.L1RootHash]*blockAndHeight
	blockM     sync.RWMutex

	transactionsPerBlockCache map[common.L2RootHash]map[gethcommon.Hash]L2Tx
	txM                       sync.RWMutex

	sharedEnclaveSecret SharedEnclaveSecret
}

func NewInMemoryDB() DB {
	return &inMemoryDB{
		statePerBlock:             make(map[common.L1RootHash]BlockState),
		stateMutex:                sync.RWMutex{},
		rollupsByHeight:           make(map[int][]*Rollup),
		rollups:                   make(map[common.L2RootHash]*Rollup),
		mempool:                   make(map[gethcommon.Hash]L2Tx),
		mpMutex:                   sync.RWMutex{},
		statePerRollup:            make(map[common.L2RootHash]State),
		blockCache:                map[common.L1RootHash]*blockAndHeight{},
		blockM:                    sync.RWMutex{},
		transactionsPerBlockCache: make(map[common.L2RootHash]map[gethcommon.Hash]L2Tx),
		txM:                       sync.RWMutex{},
	}
}

func (db *inMemoryDB) FetchState(hash common.L1RootHash) (BlockState, bool) {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	val, found := db.statePerBlock[hash]
	return val, found
}

func (db *inMemoryDB) SetState(hash common.L1RootHash, state BlockState) {
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

func (db *inMemoryDB) SetRollupState(hash common.L2RootHash, state State) {
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

func (db *inMemoryDB) Balance(address gethcommon.Address) uint64 {
	db.assertSecretAvailable()
	return db.Head().State[address]
}

func (db *inMemoryDB) StoreRollup(rollup *Rollup) {
	db.assertSecretAvailable()
	height := db.HeightRollup(rollup)

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

func (db *inMemoryDB) FetchRollup(hash common.L2RootHash) *Rollup {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	r, f := db.rollups[hash]
	if !f {
		panic(fmt.Sprintf("Could not find rollup: r_%s", hash))
	}
	return r
}

func (db *inMemoryDB) ExistRollup(hash common.L2RootHash) bool {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	_, f := db.rollups[hash]
	return f
}

func (db *inMemoryDB) FetchGossipedRollups(height int) []*Rollup {
	db.assertSecretAvailable()
	db.stateMutex.RLock()
	defer db.stateMutex.RUnlock()
	return db.rollupsByHeight[height]
}

func (db *inMemoryDB) FetchRollupState(hash common.L2RootHash) State {
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
	// txStr := make([]nodegethcommon.TxHash, 0)
	mpCopy := make([]L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
		// txStr = append(txStr, tx.ID)
	}
	// nodegethcommon.Log(fmt.Sprintf(">>> %v <<<", txStr))
	return mpCopy
}

func (db *inMemoryDB) PruneTxs(toRemove map[gethcommon.Hash]gethcommon.Hash) {
	db.assertSecretAvailable()
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	r := make(map[gethcommon.Hash]L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	// fmt.Printf("len(mempool) := %d\n", len(r))
	db.mempool = r
}

func (db *inMemoryDB) StoreBlock(b *types.Block) {
	db.assertSecretAvailable()
	db.blockM.Lock()
	defer db.blockM.Unlock()

	if b.ParentHash() == common.GenesisHash {
		db.blockCache[b.Hash()] = &blockAndHeight{b, 0}
		return
	}

	p, f := db.blockCache[b.ParentHash()]
	if !f {
		panic("Should not happen")
	}
	db.blockCache[b.Hash()] = &blockAndHeight{b: b, height: p.height + 1}
}

func (db *inMemoryDB) ResolveBlock(hash common.L1RootHash) (*types.Block, bool) {
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

func (db *inMemoryDB) Txs(r *Rollup) (map[gethcommon.Hash]L2Tx, bool) {
	db.assertSecretAvailable()
	db.txM.RLock()
	val, found := db.transactionsPerBlockCache[r.Hash()]
	db.txM.RUnlock()
	return val, found
}

func (db *inMemoryDB) AddTxs(r *Rollup, newMap map[gethcommon.Hash]L2Tx) {
	db.assertSecretAvailable()
	db.txM.Lock()
	db.transactionsPerBlockCache[r.Hash()] = newMap
	db.txM.Unlock()
}

func (db *inMemoryDB) ParentRollup(r *Rollup) *Rollup {
	db.assertSecretAvailable()
	return db.FetchRollup(r.Header.ParentHash)
}

func (db *inMemoryDB) HeightRollup(r *Rollup) int {
	db.assertSecretAvailable()
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(common.L2GenesisHeight)
		return common.L2GenesisHeight
	}
	v := db.HeightRollup(db.ParentRollup(r)) + 1
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

func (db *inMemoryDB) HeightBlock(block *types.Block) int {
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	b, f := db.blockCache[block.Hash()]
	if !f {
		panic("should not happen")
	}
	return b.height
}

func (db *inMemoryDB) ParentBlock(block *types.Block) (*types.Block, bool) {
	db.blockM.RLock()
	defer db.blockM.RUnlock()
	b, f := db.blockCache[block.ParentHash()]
	return b.b, f
}
