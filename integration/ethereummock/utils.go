package ethereummock

import (
	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

// LCA - returns the least common ancestor of the 2 blocks
func LCA(blockA *common2.Block, blockB *common2.Block, resolver common2.BlockResolver) *common2.Block {
	if blockA.Height(resolver) == common2.L1GenesisHeight || blockB.Height(resolver) == common2.L1GenesisHeight {
		return blockA
	}
	if blockA.Hash() == blockB.Hash() {
		return blockA
	}
	if blockA.Height(resolver) > blockB.Height(resolver) {
		p, f := blockA.Parent(resolver)
		if !f {
			panic("wtf")
		}
		return LCA(p, blockB, resolver)
	}
	if blockB.Height(resolver) > blockA.Height(resolver) {
		p, f := blockB.Parent(resolver)
		if !f {
			panic("wtf")
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, f := blockA.Parent(resolver)
	if !f {
		panic("wtf")
	}
	parentBlockB, f := blockB.Parent(resolver)
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
	if b.Height(r) == common2.L1GenesisHeight {
		return makeMap(b.Transactions)
	}
	newMap := make(map[common2.TxHash]*common2.L1Tx)
	p, f := b.Parent(r)
	if !f {
		panic("wtf")
	}
	for k, v := range allIncludedTransactions(p, r, db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions {
		newMap[tx.ID] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []*common2.L1Tx, toRemove map[common2.TxHash]*common2.L1Tx) (r []*common2.L1Tx) {
	for _, t := range base {
		_, f := toRemove[t.ID]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs []*common2.L1Tx) map[common2.TxHash]*common2.L1Tx {
	m := make(map[common2.TxHash]*common2.L1Tx)
	for _, tx := range txs {
		m[tx.ID] = tx
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
		tempBlock, found = tempBlock.Parent(resolver)
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
