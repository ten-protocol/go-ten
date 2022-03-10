package simulation

import (
	"fmt"

	common3 "github.com/obscuronet/obscuro-playground/go/common"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func printBlock(b *common3.Block, m ethereum_mock.Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions {
		if tx.TxType == common3.RollupTx {
			r := common2.DecodeRollup(tx.Rollup)
			txs = append(txs, fmt.Sprintf("r_%d", common3.ShortHash(r.Hash())))
		} else {
			txs = append(txs, fmt.Sprintf("deposit(%v=%d)", tx.Dest, tx.Amount))
		}
	}
	p, f := b.Parent(m.Resolver)
	if !f {
		panic("wtf")
	}

	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[parent=b_%d]. Txs: %v", m.ID, common3.ShortHash(b.Hash()), b.Height(m.Resolver), b.Header.Nonce, common3.ShortHash(p.Hash()), txs)
}
