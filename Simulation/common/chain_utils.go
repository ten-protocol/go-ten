package common

import (
	"github.com/google/uuid"
	"sync"
)

type RootHash = uuid.UUID
type TxHash = uuid.UUID

type ChainNode interface {
	Parent() ChainNode
	Height() int
	RootHash() RootHash
	Txs() []Tx
}

type Tx interface {
	Hash() TxHash
}

//LCA - returns the least common ancestor of the 2 blocks
func LCA(a ChainNode, b ChainNode) ChainNode {
	if a.Height() == -1 || b.Height() == -1 {
		return a
	}
	if a.RootHash() == b.RootHash() {
		return a
	}
	if a.Height() > b.Height() {
		return LCA(a.Parent(), b)
	}
	if b.Height() > a.Height() {
		return LCA(a, b.Parent())
	}
	return LCA(a.Parent(), b.Parent())
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(a ChainNode, b ChainNode) bool {
	if a.RootHash() == b.RootHash() {
		return true
	}
	if a.Height() >= b.Height() {
		return false
	}
	return IsAncestor(a, b.Parent())
}

var transactionsPerBlockCache = make(map[RootHash]map[TxHash]Tx)
var rpbcM = &sync.RWMutex{}

//FindNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
//todo - inefficient
func FindNotIncludedTxs(head ChainNode, txs []Tx) []Tx {
	included := allIncludedTransactions(head)
	return removeExisting(txs, included)
}

func makeMap(txs []Tx) map[TxHash]Tx {
	m := make(map[TxHash]Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func allIncludedTransactions(b ChainNode) map[TxHash]Tx {
	rpbcM.RLock()
	val, found := transactionsPerBlockCache[b.RootHash()]
	rpbcM.RUnlock()
	if found {
		return val
	}
	if b.Height() == -1 {
		return makeMap(b.Txs())
	}
	var newMap = make(map[TxHash]Tx)
	for k, v := range allIncludedTransactions(b.Parent()) {
		newMap[k] = v
	}
	for _, tx := range b.Txs() {
		newMap[tx.Hash()] = tx
	}
	rpbcM.Lock()
	transactionsPerBlockCache[b.RootHash()] = newMap
	rpbcM.Unlock()
	return newMap
}

func removeExisting(base []Tx, toRemove map[TxHash]Tx) (r []Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}
