package enclave

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/common"
)

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *Rollup, txs []L2Tx, db DB) []L2Tx {
	included := allIncludedTransactions(head, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *Rollup, db DB) map[common.TxHash]L2Tx {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if db.Height(b) == common.L2GenesisHeight {
		return makeMap(b.Transactions)
	}
	newMap := make(map[common.TxHash]L2Tx)
	for k, v := range allIncludedTransactions(db.Parent(b), db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions {
		newMap[tx.ID] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []L2Tx, toRemove map[common.TxHash]L2Tx) (r []L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.ID]
		if !f {
			r = append(r, t)
		}
	}
	return
}

// Returns all transactions found 20 levels below
func historicTxs(r *Rollup, db DB) map[common.TxHash]common.TxHash {
	i := common.HeightCommittedBlocks
	c := r
	for {
		if i == 0 || db.Height(c) == common.L2GenesisHeight {
			return toMap(c.Transactions)
		}
		i--
		c = db.Parent(c)
	}
}

func makeMap(txs []L2Tx) map[common.TxHash]L2Tx {
	m := make(map[common.TxHash]L2Tx)
	for _, tx := range txs {
		m[tx.ID] = tx
	}
	return m
}

func toMap(txs []L2Tx) map[common.TxHash]common.TxHash {
	m := make(map[common.TxHash]common.TxHash)
	for _, tx := range txs {
		m[tx.ID] = tx.ID
	}
	return m
}

func printTxs(txs []L2Tx) (txsString []string) {
	for _, t := range txs {
		txsString = printTx(t, txsString)
	}
	return txsString
}

func printTx(t L2Tx, txsString []string) []string {
	switch t.TxType {
	case TransferTx:
		txsString = append(txsString, fmt.Sprintf("%v->%v(%d){%d}", t.From, t.To, t.Amount, t.ID.ID()))
	case WithdrawalTx:
		txsString = append(txsString, fmt.Sprintf("%v->*(%d){%d}", t.From, t.Amount, t.ID.ID()))
	}
	return txsString
}
