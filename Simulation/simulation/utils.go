package simulation

import (
	"fmt"
	"simulation/common"
	ethereum_mock "simulation/ethereum-mock"
	"simulation/obscuro"
)

func printBlock(b *common.Block, m ethereum_mock.Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions {
		if tx.TxType == common.RollupTx {
			r := obscuro.DecodeRollup(tx.Rollup)
			txs = append(txs, fmt.Sprintf("r_%s", common.Str(r.Hash())))
		} else {
			txs = append(txs, fmt.Sprintf("deposit(%v=%d)", tx.Dest, tx.Amount))
		}
	}
	p, f := b.Parent(m.Resolver)
	if !f {
		panic("wtf")
	}
	return fmt.Sprintf("> M%d: create b_%s(Height=%d, Nonce=%d)[p=b_%s]. Txs: %v", m.Id, common.Str(b.Hash()), b.Height(m.Resolver), b.Header.Nonce, common.Str(p.Hash()), txs)
}
