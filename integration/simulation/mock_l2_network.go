package simulation

import (
	"sync/atomic"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroInMemNetwork - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type ObscuroInMemNetwork struct {
	currentNode *host.Node
	Nodes       []*host.Node

	avgLatency       uint64
	avgBlockDuration uint64

	listenerInterrupt *int32
}

// NewObscuroInMemNetwork returns an instance of a configured L2 Network (no nodes)
func NewObscuroInMemNetwork(avgBlockDuration uint64, avgLatency uint64) *ObscuroInMemNetwork {
	i := int32(0)
	return &ObscuroInMemNetwork{
		avgLatency:        avgLatency,
		avgBlockDuration:  avgBlockDuration,
		listenerInterrupt: &i,
	}
}

func (netw *ObscuroInMemNetwork) StartListening(callback host.P2PCallback) {
	// nothing to do here, since communication is direct through the in memory objects
}

func (netw *ObscuroInMemNetwork) StopListening() {
	atomic.StoreInt32(netw.listenerInterrupt, 1)
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (netw *ObscuroInMemNetwork) BroadcastRollup(r obscurocommon.EncodedRollup) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if a.ID != netw.currentNode.ID {
			t := a
			obscurocommon.Schedule(netw.delay(), func() { t.P2PGossipRollup(r) })
		}
	}
}

func (netw *ObscuroInMemNetwork) BroadcastTx(tx nodecommon.EncryptedTx) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if a.ID != netw.currentNode.ID {
			t := a
			obscurocommon.Schedule(netw.delay()/2, func() { t.P2PReceiveTx(tx) })
		}
	}
}

// delay returns an expected delay on the l2
func (netw *ObscuroInMemNetwork) delay() uint64 {
	return obscurocommon.RndBtw(netw.avgLatency/10, 2*netw.avgLatency)
}
