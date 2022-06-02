package db

import (
	"sync"

	"github.com/obscuronet/go-obscuro/go/obscuronode/enclave/core"

	"github.com/obscuronet/go-obscuro/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/obscurocommon"
)

// InMemoryDB lives purely in the encrypted memory space of an enclave.
// Unlike Storage, methods in this class should have minimal logic, to map them more easily to our chosen datastore.
type InMemoryDB struct {
	txsPerRollupCache map[obscurocommon.L2RootHash]map[common.Hash]nodecommon.L2Tx
	txMutex           sync.RWMutex // Controls access to `txsPerRollupCache`
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		txsPerRollupCache: make(map[obscurocommon.L2RootHash]map[common.Hash]nodecommon.L2Tx),
		txMutex:           sync.RWMutex{},
	}
}

func (db *InMemoryDB) FetchRollupTxs(r *core.Rollup) (map[common.Hash]nodecommon.L2Tx, bool) {
	db.txMutex.RLock()
	defer db.txMutex.RUnlock()

	val, found := db.txsPerRollupCache[r.Hash()]
	return val, found
}

func (db *InMemoryDB) StoreRollupTxs(r *core.Rollup, newMap map[common.Hash]nodecommon.L2Tx) {
	db.txMutex.Lock()
	defer db.txMutex.Unlock()

	db.txsPerRollupCache[r.Hash()] = newMap
}
