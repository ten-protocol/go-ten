package enclave

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *core.Rollup, txs []nodecommon.L2Tx, s db.Storage) []nodecommon.L2Tx {
	included := allIncludedTransactions(head, s)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *core.Rollup, s db.Storage) map[common.Hash]nodecommon.L2Tx {
	val, found := s.FetchRollupTxs(b)
	if found {
		return val
	}
	if b.Header.Height == obscurocommon.L2GenesisHeight {
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
func historicTxs(r *core.Rollup, s db.Storage) map[common.Hash]common.Hash {
	i := obscurocommon.HeightCommittedBlocks
	c := r
	for {
		if i == 0 || c.Header.Height == obscurocommon.L2GenesisHeight {
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
	txData := core.TxData(&t)
	switch txData.Type {
	case core.TransferTx:
		txsString = append(txsString, fmt.Sprintf("%d->%d(%d){%d}", obscurocommon.ShortAddress(txData.From), obscurocommon.ShortAddress(txData.To), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	case core.WithdrawalTx:
		txsString = append(txsString, fmt.Sprintf("%d->*(%d){%d}", obscurocommon.ShortAddress(txData.From), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	case core.DepositTx:
		txsString = append(txsString, fmt.Sprintf("*->%d(%d){%d}", obscurocommon.ShortAddress(txData.To), txData.Amount, obscurocommon.ShortHash(t.Hash())))
	}
	return txsString
}

func contains(s []obscurocommon.TxHash, e obscurocommon.TxHash) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
