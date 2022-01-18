package obscuro

import (
	"sync"
)

//FindNotIncludedL1Txs - given a list of rollups, it keeps only the ones that were not included in the block
//todo - inefficient
func FindNotIncludedL1Txs(head *Block, l1Txs []*L1Tx) []*L1Tx {
	included := allIncludedTransactions(head)
	return removeExisting(l1Txs, included)
}

var transactionsPerBlockCache = make(map[L1RootHash]map[L1TxId]*L1Tx)
var rpbcM = &sync.RWMutex{}

func makeMap(txs []*L1Tx) map[L1TxId]*L1Tx {
	m := make(map[L1TxId]*L1Tx)
	for _, tx := range txs {
		m[tx.id] = tx
	}
	return m
}

func allIncludedTransactions(b *Block) map[L1TxId]*L1Tx {
	rpbcM.RLock()
	val, found := transactionsPerBlockCache[b.rootHash]
	rpbcM.RUnlock()
	if found {
		return val
	}
	if b.height == -1 {
		return makeMap(b.txs)
	}
	var newMap = make(map[L1TxId]*L1Tx)
	for k, v := range allIncludedTransactions(b.parent) {
		newMap[k] = v
	}
	for _, tx := range b.txs {
		newMap[tx.id] = tx
	}
	rpbcM.Lock()
	transactionsPerBlockCache[b.rootHash] = newMap
	rpbcM.Unlock()
	return newMap
}

func removeExisting(base []*L1Tx, toRemove map[L1TxId]*L1Tx) (r []*L1Tx) {
	for _, t := range base {
		_, f := toRemove[t.id]
		if !f {
			r = append(r, t)
		}
	}
	return
}
