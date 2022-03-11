package simulation

import (
	"time"

	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// L2NetworkCfg - models a full network including artificial random latencies
type L2NetworkCfg struct {
	nodes []*obscuro_node.Node
	delay common.Latency // the latency
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (cfg *L2NetworkCfg) BroadcastRollup(r common.EncodedRollup) {
	for _, a := range cfg.nodes {
		rol := nodecommon.DecodeRollup(r)
		if a.ID != rol.Header.Agg {
			t := a
			common.Schedule(cfg.delay(), func() { t.P2PGossipRollup(r) })
		}
	}
}

func (cfg *L2NetworkCfg) BroadcastTx(tx nodecommon.EncryptedTx) {
	for _, a := range cfg.nodes {
		t := a
		common.Schedule(cfg.delay()/2, func() { t.P2PReceiveTx(tx) })
	}
}

func (cfg *L2NetworkCfg) Start(delay time.Duration) {
	// Start l1 nodes
	for _, m := range cfg.nodes {
		t := m
		go t.Start()
		// don't start everything at once
		time.Sleep(delay)
	}
}

func (cfg *L2NetworkCfg) Stop() {
	for _, m := range cfg.nodes {
		m.Stop()
		// fmt.Printf("Stopped L2 node: %d\n", m.ID)
	}
}
