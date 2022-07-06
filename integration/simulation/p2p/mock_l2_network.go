package p2p

import (
	"bytes"
	"sync/atomic"
	"time"

	"github.com/obscuronet/obscuro-playground/go/host"

	"github.com/obscuronet/obscuro-playground/go/common"
)

// MockP2P - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type MockP2P struct {
	CurrentNode *host.Node
	Nodes       []*host.Node

	avgLatency       time.Duration
	avgBlockDuration time.Duration

	listenerInterrupt *int32
}

// NewMockP2P returns an instance of a configured L2 Network (no nodes)
func NewMockP2P(avgBlockDuration time.Duration, avgLatency time.Duration) *MockP2P {
	i := int32(0)
	return &MockP2P{
		avgLatency:        avgLatency,
		avgBlockDuration:  avgBlockDuration,
		listenerInterrupt: &i,
	}
}

func (netw *MockP2P) StartListening(host.P2PCallback) {
	// nothing to do here, since communication is direct through the in memory objects
}

func (netw *MockP2P) StopListening() error {
	atomic.StoreInt32(netw.listenerInterrupt, 1)
	return nil
}

func (netw *MockP2P) UpdatePeerList([]string) {
	// Do nothing.
}

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (netw *MockP2P) BroadcastRollup(r common.EncodedRollup) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if !bytes.Equal(a.ID.Bytes(), netw.CurrentNode.ID.Bytes()) {
			t := a
			common.Schedule(netw.delay(), func() { t.ReceiveRollup(r) })
		}
	}
}

func (netw *MockP2P) BroadcastTx(tx common.EncryptedTx) {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return
	}

	for _, a := range netw.Nodes {
		if !bytes.Equal(a.ID.Bytes(), netw.CurrentNode.ID.Bytes()) {
			t := a
			common.Schedule(netw.delay()/2, func() { t.ReceiveTx(tx) })
		}
	}
}

// delay returns an expected delay on the l2
func (netw *MockP2P) delay() time.Duration {
	return common.RndBtwTime(netw.avgLatency/10, 2*netw.avgLatency)
}
