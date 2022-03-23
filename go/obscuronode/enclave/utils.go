package enclave

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *Rollup, txs []L2Tx, s Storage) []L2Tx {
	included := allIncludedTransactions(head, s)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *Rollup, s Storage) map[common.Hash]L2Tx {
	val, found := s.FetchRollupTxs(b)
	if found {
		return val
	}
	if s.HeightRollup(b) == obscurocommon.L2GenesisHeight {
		return makeMap(b.Transactions)
	}
	newMap := make(map[common.Hash]L2Tx)
	for k, v := range allIncludedTransactions(s.ParentRollup(b), s) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions {
		newMap[tx.Hash()] = tx
	}
	s.StoreRollupTxs(b, newMap)
	return newMap
}

func removeExisting(base []L2Tx, toRemove map[common.Hash]L2Tx) (r []L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

// Returns all transactions found 20 levels below
func historicTxs(r *Rollup, s Storage) map[common.Hash]common.Hash {
	i := obscurocommon.HeightCommittedBlocks
	c := r
	for {
		if i == 0 || s.HeightRollup(c) == obscurocommon.L2GenesisHeight {
			return toMap(c.Transactions)
		}
		i--
		c = s.ParentRollup(c)
	}
}

func makeMap(txs []L2Tx) map[common.Hash]L2Tx {
	m := make(map[common.Hash]L2Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func toMap(txs []L2Tx) map[common.Hash]common.Hash {
	m := make(map[common.Hash]common.Hash)
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
		txsString = append(txsString, fmt.Sprintf("%d->%d(%d){%d}", obscurocommon.ShortAddress(txData.From), obscurocommon.ShortAddress(txData.To), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	case WithdrawalTx:
		txsString = append(txsString, fmt.Sprintf("%d->*(%d){%d}", obscurocommon.ShortAddress(txData.From), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	}
	return txsString
}
