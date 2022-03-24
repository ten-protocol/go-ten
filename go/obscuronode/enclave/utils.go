package enclave

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *Rollup, txs []nodecommon.L2Tx, s Storage) []nodecommon.L2Tx {
	included := allIncludedTransactions(head, s)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *Rollup, s Storage) map[common.Hash]nodecommon.L2Tx {
	val, found := s.FetchRollupTxs(b)
	if found {
		return val
	}
	if s.HeightRollup(b) == obscurocommon.L2GenesisHeight {
		return makeMap(b.Transactions)
	}
	newMap := make(map[common.Hash]nodecommon.L2Tx)
	for k, v := range allIncludedTransactions(s.ParentRollup(b), s) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions {
		newMap[tx.Hash()] = tx
	}
	s.StoreRollupTxs(b, newMap)
	return newMap
}

func removeExisting(base []nodecommon.L2Tx, toRemove map[common.Hash]nodecommon.L2Tx) (r []nodecommon.L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

// Returns all transactions found 20 levels below
func historicTxs(r *Rollup, s Storage, genesisRollup *Rollup) map[common.Hash]common.Hash {
	i := obscurocommon.HeightCommittedBlocks
	c := r
	for {
		if i == 0 || c.Hash() == genesisRollup.Hash() {
			return toMap(c.Transactions)
		}
		i--
		c = s.ParentRollup(c)
	}
}

func makeMap(txs []nodecommon.L2Tx) map[common.Hash]nodecommon.L2Tx {
	m := make(map[common.Hash]nodecommon.L2Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func toMap(txs []nodecommon.L2Tx) map[common.Hash]common.Hash {
	m := make(map[common.Hash]common.Hash)
	for _, tx := range txs {
		m[tx.Hash()] = tx.Hash()
	}
	return m
}

func printTxs(txs []nodecommon.L2Tx) (txsString []string) {
	for _, t := range txs {
		txsString = printTx(t, txsString)
	}
	return txsString
}

func printTx(t nodecommon.L2Tx, txsString []string) []string {
	txData := TxData(&t)
	switch txData.Type {
	case TransferTx:
		txsString = append(txsString, fmt.Sprintf("%d->%d(%d){%d}", obscurocommon.ShortAddress(txData.From), obscurocommon.ShortAddress(txData.To), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	case WithdrawalTx:
		txsString = append(txsString, fmt.Sprintf("%d->*(%d){%d}", obscurocommon.ShortAddress(txData.From), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	}
	return txsString
}
