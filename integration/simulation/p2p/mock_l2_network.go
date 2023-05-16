package p2p

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/async"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"
)

// MockP2P - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type MockP2P struct {
	CurrentNode host.Host
	Nodes       []host.Host

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

func (netw *MockP2P) SendTxToSequencer(tx common.EncryptedTx) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}
	async.Schedule(netw.delay()/2, func() { netw.Nodes[0].ReceiveTx(tx) })
	return nil
}

func (netw *MockP2P) BroadcastBatch(batchMsg *host.BatchMsg) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	for _, node := range netw.Nodes {
		if node.Config().ID.Hex() != netw.CurrentNode.Config().ID.Hex() {
			tempNode := node
			async.Schedule(netw.delay()/2, func() { tempNode.ReceiveBatches(encodedBatchMsg) })
		}
	}

	return nil
}

func (netw *MockP2P) RequestBatchesFromSequencer(batchRequest *common.BatchRequest) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	encodedBatchRequest, err := rlp.EncodeToBytes(batchRequest)
	if err != nil {
		return fmt.Errorf("could not encode batch request using RLP. Cause: %w", err)
	}
	async.Schedule(netw.delay()/2, func() { netw.Nodes[0].ReceiveBatchRequest(encodedBatchRequest) })
	return nil
}

func (netw *MockP2P) SendBatches(batchMsg *host.BatchMsg, requesterAddress string) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	var requester host.Host
	for _, node := range netw.Nodes {
		if node.Config().P2PPublicAddress == requesterAddress {
			requester = node
		}
	}

	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	async.Schedule(netw.delay()/2, func() { requester.ReceiveBatches(encodedBatchMsg) })
	return nil
}

func (netw *MockP2P) Status() *host.P2PStatus {
	return &host.P2PStatus{}
}

func (netw *MockP2P) HealthCheck() bool {
	return true
}

// delay returns an expected delay on the l2
func (netw *MockP2P) delay() time.Duration {
	return testcommon.RndBtwTime(netw.avgLatency/10, 2*netw.avgLatency)
}
