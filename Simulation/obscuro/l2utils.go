package obscuro

import (
	"sync"
)

//FindNotIncludedTxs - given a list of txs, it keeps only the ones that were not included in the rollup
//todo - inefficient
func FindNotIncludedTxs(head *Rollup, txs []*L2Tx) []*L2Tx {
	included := allIncludedTxs(head)
	return removeExistingTxs(txs, included)
}

var txsPerRollupCache = make(map[L2RootHash]map[L2TxId]*L2Tx)
var trcm = &sync.RWMutex{}

func allIncludedTxs(r *Rollup) map[L2TxId]*L2Tx {
	trcm.RLock()
	val, found := txsPerRollupCache[r.rootHash]
	trcm.RUnlock()
	if found {
		return val
	}
	if r.height == -1 {
		return makeTxMap(r.txs)
	}
	var newMap = make(map[L2TxId]*L2Tx)
	for k, v := range allIncludedTxs(r.parent) {
		newMap[k] = v
	}
	for _, tx := range r.txs {
		newMap[tx.id] = tx
	}
	trcm.Lock()
	txsPerRollupCache[r.rootHash] = newMap
	trcm.Unlock()
	return newMap
}

func makeTxMap(txs []*L2Tx) map[L2TxId]*L2Tx {
	m := make(map[L2TxId]*L2Tx)
	for _, tx := range txs {
		m[tx.id] = tx
	}
	return m
}

func removeExistingTxs(base []*L2Tx, toRemove map[L2TxId]*L2Tx) (r []*L2Tx) {
	for _, tx := range base {
		_, f := toRemove[tx.id]
		if !f {
			r = append(r, tx)
		}
	}
	return
}
