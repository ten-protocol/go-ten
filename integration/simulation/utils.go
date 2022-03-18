package simulation

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

func printBlock(b *types.Block, m ethereum_mock.Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions() {
		t := obscurocommon.TxData(tx)
		if t.TxType == obscurocommon.RollupTx {
			r := nodecommon.DecodeRollup(t.Rollup)
			txs = append(txs, fmt.Sprintf("r_%d", obscurocommon.ShortHash(r.Hash())))
		} else if t.TxType == obscurocommon.DepositTx {
			txs = append(txs, fmt.Sprintf("deposit(%d=%d)", obscurocommon.ShortAddress(t.Dest), t.Amount))
		}
	}
	p, f := m.Resolver.ParentBlock(b)
	if !f {
		panic("wtf")
	}

	return fmt.Sprintf("> M%d: create b_%d(Height=%d, Nonce=%d)[parent=b_%d]. Txs: %v",
		obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Hash()), m.Resolver.HeightBlock(b), obscurocommon.ShortNonce(b.Header().Nonce), obscurocommon.ShortHash(p.Hash()), txs)
}
