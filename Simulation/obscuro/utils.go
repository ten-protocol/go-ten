package obscuro

import "simulation/common"

// Returns all transactions found 20 levels below
func historicTxs(r common.Rollup) map[common.TxHash]common.TxHash {
	i := common.HeightCommittedBlocks
	c := &r
	for {
		if i == 0 || c.H == common.GenesisHeight {
			return common.ToMap(c.Transactions)
		}
		i--
		c = c.ParentRollup()
	}
}

func makeMap(txs []common.L2Tx) map[common.TxHash]common.L2Tx {
	m := make(map[common.TxHash]common.L2Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}
