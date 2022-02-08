package simulation

import (
	"fmt"
	"simulation/common"
	ethereum_mock "simulation/ethereum-mock"
	"simulation/obscuro"
)

func printBlock(b common.Block, m ethereum_mock.Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions {
		if tx.TxType == common.RollupTx {
			txs = append(txs, fmt.Sprintf("r_%d", obscuro.DecodeRollup(tx.Rollup).RootHash.ID()))
		} else {
			txs = append(txs, fmt.Sprintf("deposit(%v=%d)", tx.Dest, tx.Amount))
		}
	}
	p, f := b.Parent(m.Resolver)
	if !f {
		panic("wtf")
	}
	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[p=b_%d]. Txs: %v", m.Id, b.RootHash.ID(), b.Height, b.Nonce, p.RootHash.ID(), txs)
}
