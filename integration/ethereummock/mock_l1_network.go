package ethereummock

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/log"
)

// MockEthNetwork - models a full network including artificial random latencies
// This is the gateway through which the mock L1 nodes communicate with each other
type MockEthNetwork struct {
	CurrentNode *Node

	AllNodes []*Node

	// config
	avgLatency       time.Duration
	avgBlockDuration time.Duration

	Stats host.StatsCollector
}

// NewMockEthNetwork returns an instance of a configured L1 Network (no nodes)
func NewMockEthNetwork(avgBlockDuration time.Duration, avgLatency time.Duration, stats host.StatsCollector) *MockEthNetwork {
	return &MockEthNetwork{
		Stats:            stats,
		avgLatency:       avgLatency,
		avgBlockDuration: avgBlockDuration,
	}
}

// BroadcastBlock broadcast a block to the l1 nodes
func (n *MockEthNetwork) BroadcastBlock(b obscurocommon.EncodedBlock, p obscurocommon.EncodedBlock) {
	bl, _ := b.Decode()
	for _, m := range n.AllNodes {
		if m.ID != n.CurrentNode.ID {
			t := m
			obscurocommon.Schedule(n.delay(), func() { t.P2PReceiveBlock(b, p) })
		} else {
			log.Info(printBlock(bl, *m))
		}
	}

	n.Stats.NewBlock(bl)
}

// BroadcastTx Broadcasts the L1 tx containing the rollup to the L1 network
func (n *MockEthNetwork) BroadcastTx(tx *types.Transaction) {
	for _, m := range n.AllNodes {
		if m.ID != n.CurrentNode.ID {
			t := m
			// the time to broadcast a tx is half that of a L1 block, because it is smaller.
			// todo - find a better way to express this
			d := n.delay() / 2
			obscurocommon.Schedule(d, func() { t.P2PGossipTx(tx) })
		}
	}
}

// delay returns an expected delay on the l1 network
func (n *MockEthNetwork) delay() time.Duration {
	return obscurocommon.RndBtwTime(n.avgLatency/10, 2*n.avgLatency)
}

func printBlock(b *types.Block, m Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions() {
		t := m.txHandler.UnPackTx(tx)
		if t == nil {
			continue
		}
		switch l1Tx := t.(type) {
		case *obscurocommon.L1RollupTx:
			r := nodecommon.DecodeRollupOrPanic(l1Tx.Rollup)
			txs = append(txs, fmt.Sprintf("r_%d", obscurocommon.ShortHash(r.Hash())))
		case *obscurocommon.L1DepositTx:
			txs = append(txs, fmt.Sprintf("deposit(%d=%d)", obscurocommon.ShortAddress(l1Tx.To), l1Tx.Amount))
		}
	}
	p, f := m.Resolver.ParentBlock(b)
	if !f {
		panic("wtf")
	}

	return fmt.Sprintf("> M%d: create b_%d(Height=%d, L1Nonce=%d)[parent=b_%d]. Txs: %v",
		obscurocommon.ShortAddress(m.ID), obscurocommon.ShortHash(b.Hash()), b.NumberU64(), obscurocommon.ShortNonce(b.Header().Nonce), obscurocommon.ShortHash(p.Hash()), txs)
}
