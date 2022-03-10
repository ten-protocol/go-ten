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
	for _, tx := range b.Transactions() {
		t := common3.TxData(tx)
		if t.TxType == common3.RollupTx {
			r := common2.DecodeRollup(t.Rollup)
			txs = append(txs, fmt.Sprintf("r_%d", common3.ShortHash(r.Hash())))
		} else if t.TxType == common3.DepositTx {
			txs = append(txs, fmt.Sprintf("deposit(%d=%d)", common3.ShortAddress(t.Dest), t.Amount))
		}
	}
	p, f := m.Resolver.ParentBlock(b)
	if !f {
		panic("wtf")
	}

	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[parent=b_%d]. Txs: %v",
		common3.ShortNodeID(m.ID), common3.ShortHash(b.Hash()), m.Resolver.HeightBlock(b), common3.ShortNonce(b.Header().Nonce), common3.ShortHash(p.Hash()), txs)
}
