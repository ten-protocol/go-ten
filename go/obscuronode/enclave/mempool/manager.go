package mempool

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type mempoolManager struct {
	mpMutex sync.RWMutex // Controls access to `mempool`
	mempool map[common.Hash]nodecommon.L2Tx
}

func New() Manager {
	return &mempoolManager{
		mempool: make(map[common.Hash]nodecommon.L2Tx),
		mpMutex: sync.RWMutex{},
	}
}

func (db *mempoolManager) AddMempoolTx(tx nodecommon.L2Tx) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	db.mempool[tx.Hash()] = tx
}

func (db *mempoolManager) FetchMempoolTxs() []nodecommon.L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()

	mpCopy := make([]nodecommon.L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *mempoolManager) RemoveMempoolTxs(toRemove map[common.Hash]common.Hash) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	r := make(map[common.Hash]nodecommon.L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	db.mempool = r
}
