package simulation

import (
	"simulation/common"
	"simulation/ethereum-mock"
	"time"
)

// L1NetworkCfg - models a full network including artificial random latencies
type L1NetworkCfg struct {
	nodes []*ethereum_mock.Node
	delay common.Latency // the latency
	Stats *Stats
}

// BroadcastBlockL1 broadcast a block to the l1 nodes
func (n L1NetworkCfg) BroadcastBlockL1(b common.Block) {
	for _, m := range n.nodes {
		if m.Id != b.Miner {
			t := m
			common.Schedule(n.delay(), func() { t.P2PReceiveBlock(b) })
		}
	}
	n.Stats.NewBlock(b)
}

// BroadcastL1Tx Broadcasts the L1 tx containing the rollup to the L1 network
func (n L1NetworkCfg) BroadcastL1Tx(tx common.L1Tx) {
	for _, m := range n.nodes {
		t := m
		// the time to broadcast a tx is half that of a L1 block, because it is smaller.
		// todo - find a better way to express this
		d := common.Max(n.delay()/2, 1)
		common.Schedule(d, func() { t.L1P2PGossipTx(tx) })
	}

	// collect Stats
	if tx.TxType == common.RollupTx {
		n.Stats.NewRollup(tx.Rollup)
	}
}

func (n L1NetworkCfg) Start(delay time.Duration) {
	// Start l1 nodes
	for _, m := range n.nodes {
		t := m
		go t.Start()
		// don't start everything at once
		time.Sleep(delay)
	}
}
func (n L1NetworkCfg) Stop() {
	// Start l1 nodes
	for _, m := range n.nodes {
		m.Stop()
	}
}
