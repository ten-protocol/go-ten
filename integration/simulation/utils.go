package simulation

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/common"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func printBlock(b *types.Block, m ethereum_mock.Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions() {
		t := common.TxData(tx)
		if t.TxType == common.RollupTx {
			r := nodecommon.DecodeRollup(t.Rollup)
			txs = append(txs, fmt.Sprintf("r_%d", common.ShortHash(r.Hash())))
		} else if t.TxType == common.DepositTx {
			txs = append(txs, fmt.Sprintf("deposit(%d=%d)", common.ShortAddress(t.Dest), t.Amount))
		}
	}
	p, f := m.Resolver.ParentBlock(b)
	if !f {
		panic("wtf")
	}

	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[parent=b_%d]. Txs: %v",
		common.ShortAddress(m.ID), common.ShortHash(b.Hash()), m.Resolver.HeightBlock(b), common.ShortNonce(b.Header().Nonce), common.ShortHash(p.Hash()), txs)
}
