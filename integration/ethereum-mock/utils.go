package ethereum_mock

import (
	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

// LCA - returns the least common ancestor of the 2 blocks
func LCA(a *common2.Block, b *common2.Block, r common2.BlockResolver) *common2.Block {
	if a.Height(r) == common2.L1GenesisHeight || b.Height(r) == common2.L1GenesisHeight {
		return a
	}
	if a.Hash() == b.Hash() {
		return a
	}
	if a.Height(r) > b.Height(r) {
		p, f := a.Parent(r)
		if !f {
			panic("wtf")
		}
		return LCA(p, b, r)
	}
	if b.Height(r) > a.Height(r) {
		p, f := b.Parent(r)
		if !f {
			panic("wtf")
		}

		return LCA(a, p, r)
	}
	p1, f := a.Parent(r)
	if !f {
		panic("wtf")
	}
	p2, f := b.Parent(r)
	if !f {
		panic("wtf")
	}

	return LCA(p1, p2, r)
}

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findNotIncludedTxs(head *common2.Block, txs []*common2.L1Tx, r common2.BlockResolver, db TxDb) []*common2.L1Tx {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *common2.Block, r common2.BlockResolver, db TxDb) map[common2.TxHash]*common2.L1Tx {
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
		newMap[tx.Id] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []*common2.L1Tx, toRemove map[common2.TxHash]*common2.L1Tx) (r []*common2.L1Tx) {
	for _, t := range base {
		_, f := toRemove[t.Id]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs []*common2.L1Tx) map[common2.TxHash]*common2.L1Tx {
	m := make(map[common2.TxHash]*common2.L1Tx)
	for _, tx := range txs {
		m[tx.Id] = tx
	}
	return m
}

func BlocksBetween(a *common2.Block, b *common2.Block, r common2.BlockResolver) []*common2.Block {
	if a.Hash() == b.Hash() {
		return []*common2.Block{a}
	}
	blocks := make([]*common2.Block, 0)
	c := b
	f := false
	for {
		blocks = append(blocks, c)
		if c.Hash() == a.Hash() {
			break
		}
		c, f = c.Parent(r)
		if !f {
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
