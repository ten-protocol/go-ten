package ethereummock

import (
	"bytes"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// LCA - returns the least common ancestor of the 2 blocks
func LCA(blockA *types.Block, blockB *types.Block, resolver db.BlockResolver) *types.Block {
	if blockA.NumberU64() == common.L1GenesisHeight || blockB.NumberU64() == common.L1GenesisHeight {
		return blockA
	}
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return blockA
	}
	if blockA.NumberU64() > blockB.NumberU64() {
		p, f := resolver.ParentBlock(blockA)
		if !f {
			log.Panic("Should not happen. Parent not found")
		}
		return LCA(p, blockB, resolver)
	}
	if blockB.NumberU64() > blockA.NumberU64() {
		p, f := resolver.ParentBlock(blockB)
		if !f {
			log.Panic("Should not happen. Parent not found")
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, f := resolver.ParentBlock(blockA)
	if !f {
		log.Panic("Should not happen. Parent not found")
	}
	parentBlockB, f := resolver.ParentBlock(blockB)
	if !f {
		log.Panic("Should not happen. Parent not found")
	}

	return LCA(parentBlockA, parentBlockB, resolver)
}

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findNotIncludedTxs(head *types.Block, txs []*types.Transaction, r db.BlockResolver, db TxDB) []*types.Transaction {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *types.Block, r db.BlockResolver, db TxDB) map[common.TxHash]*types.Transaction {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if b.NumberU64() == common.L1GenesisHeight {
		return makeMap(b.Transactions())
	}
	newMap := make(map[common.TxHash]*types.Transaction)
	p, f := r.ParentBlock(b)
	if !f {
		log.Panic("Should not happen. Parent not found")
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

func removeExisting(base []*types.Transaction, toRemove map[common.TxHash]*types.Transaction) (r []*types.Transaction) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs types.Transactions) map[common.TxHash]*types.Transaction {
	m := make(map[common.TxHash]*types.Transaction)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}
