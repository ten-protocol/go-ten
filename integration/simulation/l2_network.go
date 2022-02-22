package simulation

import (
	common3 "github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/obscuro-node"
	"github.com/obscuronet/obscuro-playground/go/obscuro-node/common"
	"time"
)

// L2NetworkCfg - models a full network including artificial random latencies
type L2NetworkCfg struct {
	nodes []*obscuro_node.Node
	delay common3.Latency // the latency
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (c *L2NetworkCfg) BroadcastRollup(r common3.EncodedRollup) {
	for _, a := range c.nodes {
		rol := common.DecodeRollup(r)
		if a.Id != rol.Header.Agg {
			t := a
			common3.Schedule(c.delay(), func() { t.P2PGossipRollup(r) })
		}
	}
}

func (c *L2NetworkCfg) BroadcastTx(tx common.EncryptedTx) {
	for _, a := range c.nodes {
		t := a
		common3.Schedule(c.delay()/2, func() { t.P2PReceiveTx(tx) })
	}
}

func (n *L2NetworkCfg) Start(delay time.Duration) {
	// Start l1 nodes
	for _, m := range n.nodes {
		t := m
		go t.Start()
		// don't start everything at once
		time.Sleep(delay)
	}
}

func (n *L2NetworkCfg) Stop() {
	for _, m := range n.nodes {
		m.Stop()
		// fmt.Printf("Stopped L2 node: %d\n", m.Id)
	}
}
