package enclave

import (
	"fmt"

	common2 "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/common"
)

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *Rollup, txs []L2Tx, db DB) []L2Tx {
	included := allIncludedTransactions(head, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *Rollup, db DB) map[common2.Hash]L2Tx {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if db.Height(b) == common.L2GenesisHeight {
		return makeMap(b.Transactions)
	}
	newMap := make(map[common2.Hash]L2Tx)
	for k, v := range allIncludedTransactions(db.Parent(b), db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions {
		newMap[tx.Hash()] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []L2Tx, toRemove map[common2.Hash]L2Tx) (r []L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

// Returns all transactions found 20 levels below
func historicTxs(r *Rollup, db DB) map[common2.Hash]common2.Hash {
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

func makeMap(txs []L2Tx) map[common2.Hash]L2Tx {
	m := make(map[common2.Hash]L2Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func toMap(txs []L2Tx) map[common2.Hash]common2.Hash {
	m := make(map[common2.Hash]common2.Hash)
	for _, tx := range txs {
		m[tx.Hash()] = tx.Hash()
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
	txData := TxData(&t)
	switch txData.Type {
	case TransferTx:
		txsString = append(txsString, fmt.Sprintf("%v->%v(%d){%d}", txData.From, txData.To, txData.Amount, t.Hash()))
	case WithdrawalTx:
		txsString = append(txsString, fmt.Sprintf("%v->*(%d){%d}", txData.From, txData.Amount, t.Hash()))
	}
	return txsString
}
