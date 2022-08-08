package p2p

import (
	"bytes"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/host"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/obscuronet/go-obscuro/go/common"
)

// MockP2P - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type MockP2P struct {
	CurrentNode host.Host
	Nodes       []host.MockHost

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

func (netw *MockP2P) StartListening(host.Host) {
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
func (netw *MockP2P) BroadcastRollup(r common.EncodedRollup) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	for _, a := range netw.Nodes {
		if !bytes.Equal(a.Config().ID.Bytes(), netw.CurrentNode.Config().ID.Bytes()) {
			t := a
			common.Schedule(netw.delay(), func() { t.ReceiveRollup(r) })
		}
	}

	return nil
}

func (netw *MockP2P) BroadcastTx(tx common.EncryptedTx) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	for _, a := range netw.Nodes {
		if !bytes.Equal(a.Config().ID.Bytes(), netw.CurrentNode.Config().ID.Bytes()) {
			t := a
			common.Schedule(netw.delay()/2, func() { t.ReceiveTx(tx) })
		}
	}

	return nil
}

// delay returns an expected delay on the l2
func (netw *MockP2P) delay() time.Duration {
	return testcommon.RndBtwTime(netw.avgLatency/10, 2*netw.avgLatency)
}
