package ethereummock

import (
	"github.com/ethereum/go-ethereum/core/types"
	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

// LCA - returns the least common ancestor of the 2 blocks
func LCA(blockA *common2.Block, blockB *common2.Block, resolver common2.BlockResolver) *common2.Block {
	if resolver.Height(blockA) == common2.L1GenesisHeight || resolver.Height(blockB) == common2.L1GenesisHeight {
		return blockA
	}
	if blockA.Hash() == blockB.Hash() {
		return blockA
	}
	if resolver.Height(blockA) > resolver.Height(blockB) {
		p, f := resolver.Parent(blockA)
		if !f {
			panic("wtf")
		}
		return LCA(p, blockB, resolver)
	}
	if resolver.Height(blockB) > resolver.Height(blockA) {
		p, f := resolver.Parent(blockB)
		if !f {
			panic("wtf")
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, f := resolver.Parent(blockA)
	if !f {
		panic("wtf")
	}
	parentBlockB, f := resolver.Parent(blockB)
	if !f {
		panic("wtf")
	}

	return LCA(parentBlockA, parentBlockB, resolver)
}

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findNotIncludedTxs(head *common2.Block, txs []*common2.L1Tx, r common2.BlockResolver, db TxDB) []*common2.L1Tx {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *common2.Block, r common2.BlockResolver, db TxDB) map[common2.TxHash]*common2.L1Tx {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if r.Height(b) == common2.L1GenesisHeight {
		return makeMap(b.Transactions())
	}
	newMap := make(map[common2.TxHash]*common2.L1Tx)
	p, f := r.Parent(b)
	if !f {
		panic("wtf")
	}
	for k, v := range allIncludedTransactions(p, r, db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions() {
		newMap[tx.Hash()] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []*common2.L1Tx, toRemove map[common2.TxHash]*common2.L1Tx) (r []*common2.L1Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs types.Transactions) map[common2.TxHash]*common2.L1Tx {
	m := make(map[common2.TxHash]*common2.L1Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func BlocksBetween(blockA *common2.Block, blockB *common2.Block, resolver common2.BlockResolver) []*common2.Block {
	if blockA.Hash() == blockB.Hash() {
		return []*common2.Block{blockA}
	}
	blocks := make([]*common2.Block, 0)
	tempBlock := blockB
	var found bool
	for {
		blocks = append(blocks, tempBlock)
		if tempBlock.Hash() == blockA.Hash() {
			break
		}
		tempBlock, found = resolver.Parent(tempBlock)
		if !found {
			panic("should not happen")
		}
	}
	n := len(blocks)
	result := make([]*common2.Block, n)
	for i, block := range blocks {
		result[n-i-1] = block
	}
	return result
}
