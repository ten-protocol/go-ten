package simulation

import (
	common2 "github.com/ethereum/go-ethereum/common"
	"time"

	common3 "github.com/obscuronet/obscuro-playground/go/common"
	obscuro_node "github.com/obscuronet/obscuro-playground/go/obscuronode"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/common"
)

// L2NetworkCfg - models a full network including artificial random latencies
type L2NetworkCfg struct {
	nodes []*obscuro_node.Node
	delay common3.Latency // the latency
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (cfg *L2NetworkCfg) BroadcastRollup(r common3.EncodedRollup) {
	for _, a := range cfg.nodes {
		rol := common.DecodeRollup(r)
		if common2.Address(a.ID) != rol.Header.Agg {
			t := a
			common3.Schedule(cfg.delay(), func() { t.P2PGossipRollup(r) })
		}
	}
}

func (cfg *L2NetworkCfg) BroadcastTx(tx common.EncryptedTx) {
	for _, a := range cfg.nodes {
		t := a
		common3.Schedule(cfg.delay()/2, func() { t.P2PReceiveTx(tx) })
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
