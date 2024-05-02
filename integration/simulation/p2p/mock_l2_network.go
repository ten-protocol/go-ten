package p2p

import (
	"context"
	"math/big"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/go/common/async"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"

	testcommon "github.com/ten-protocol/go-ten/integration/common"
)

const _sequencerID = "0"

type MockP2PNetwork struct {
	nodes map[string]*MockP2P

	avgLatency                  time.Duration
	avgBlockDuration            time.Duration
	nodeWithIncomingP2PDisabled int
}

type MockP2PNetworkIntf interface {
	NewNode(id int) host.P2PHostService
}

func NewMockP2PNetwork(avgBlockDuration time.Duration, avgLatency time.Duration, nodeWithIncomingP2PDisabled int) MockP2PNetworkIntf {
	return &MockP2PNetwork{
		nodes:                       make(map[string]*MockP2P),
		avgBlockDuration:            avgBlockDuration,
		avgLatency:                  avgLatency,
		nodeWithIncomingP2PDisabled: nodeWithIncomingP2PDisabled,
	}
}

func (m *MockP2PNetwork) NewNode(id int) host.P2PHostService {
	idStr := strconv.Itoa(id)
	isIncomingP2PDisabled := m.nodeWithIncomingP2PDisabled != 0 && m.nodeWithIncomingP2PDisabled == id
	node := NewMockP2P(m, idStr, isIncomingP2PDisabled)
	m.nodes[idStr] = node
	return node
}

func (m *MockP2PNetwork) RequestBatchesFromSequencer(id string, fromSeqNo *big.Int) {
	seqNode := m.nodes[_sequencerID]
	async.Schedule(m.delay()/2, func() { seqNode.ReceiveBatchRequest(id, fromSeqNo) })
}

func (m *MockP2PNetwork) SendTransactionToSequencer(tx common.EncryptedTx) {
	seqNode := m.nodes[_sequencerID]
	async.Schedule(m.delay()/2, func() { seqNode.ReceiveTransaction(tx) })
}

func (m *MockP2PNetwork) BroadcastBatch(fromNodeID string, batches []*common.ExtBatch) {
	for _, node := range m.nodes {
		if node.id != fromNodeID {
			tempNode := node
			async.Schedule(m.delay()/2, func() { tempNode.ReceiveBatches(batches, true) })
		}
	}
}

func (m *MockP2PNetwork) RespondToBatchRequest(requesterID string, batches []*common.ExtBatch) {
	async.Schedule(m.delay()/2, func() {
		requester, ok := m.nodes[requesterID]
		if !ok {
			panic("requester not found in mock p2p service")
		}
		requester.ReceiveBatches(batches, false)
	})
}

// delay returns an expected delay on the l2
func (m *MockP2PNetwork) delay() time.Duration {
	return testcommon.RndBtwTime(m.avgLatency/10, 2*m.avgLatency)
}

// MockP2P - models the p2p service of a host, but instead of sending messages over tcp it uses the `MockP2PNetwork` to distribute messages
type MockP2P struct {
	id      string
	network *MockP2PNetwork // reference to the mock network

	batchSubscribers *subscription.Manager[host.P2PBatchHandler]
	txSubscribers    *subscription.Manager[host.P2PTxHandler]
	batchReqHandlers *subscription.Manager[host.P2PBatchRequestHandler]

	listenerInterrupt     *int32
	isIncomingP2PDisabled bool
}

// NewMockP2P returns an instance of a configured L2 Network (no nodes)
func NewMockP2P(network *MockP2PNetwork, id string, isIncomingP2PDisabled bool) *MockP2P {
	i := int32(0)
	return &MockP2P{
		id:      id,
		network: network,

		batchSubscribers:      subscription.NewManager[host.P2PBatchHandler](),
		txSubscribers:         subscription.NewManager[host.P2PTxHandler](),
		batchReqHandlers:      subscription.NewManager[host.P2PBatchRequestHandler](),
		listenerInterrupt:     &i,
		isIncomingP2PDisabled: isIncomingP2PDisabled,
	}
}

func (n *MockP2P) Start() error {
	// nothing to do here, since communication is direct through the in memory objects
	return nil
}

func (n *MockP2P) Stop() error {
	atomic.StoreInt32(n.listenerInterrupt, 1)
	return nil
}

func (n *MockP2P) HealthStatus(context.Context) host.HealthStatus {
	return &host.BasicErrHealthStatus{ErrMsg: ""}
}

func (n *MockP2P) UpdatePeerList([]string) {
	// Do nothing.
}

func (n *MockP2P) SendTxToSequencer(tx common.EncryptedTx) error {
	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}
	n.network.SendTransactionToSequencer(tx)
	return nil
}

func (n *MockP2P) BroadcastBatches(batches []*common.ExtBatch) error {
	if n.isIncomingP2PDisabled {
		return nil
	}

	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}

	n.network.BroadcastBatch(n.id, batches)

	return nil
}

func (n *MockP2P) SubscribeForBatches(handler host.P2PBatchHandler) func() {
	if n.isIncomingP2PDisabled {
		return func() {}
	}
	return n.batchSubscribers.Subscribe(handler)
}

func (n *MockP2P) SubscribeForTx(handler host.P2PTxHandler) func() {
	if n.isIncomingP2PDisabled {
		return func() {}
	}
	return n.txSubscribers.Subscribe(handler)
}

func (n *MockP2P) SubscribeForBatchRequests(handler host.P2PBatchRequestHandler) func() {
	if n.isIncomingP2PDisabled {
		return func() {}
	}

	return n.batchReqHandlers.Subscribe(handler)
}

func (n *MockP2P) RequestBatchesFromSequencer(fromSeqNo *big.Int) error {
	if n.isIncomingP2PDisabled {
		return nil
	}

	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}
	n.network.RequestBatchesFromSequencer(n.id, fromSeqNo)
	return nil
}

func (n *MockP2P) RespondToBatchRequest(requesterID string, batches []*common.ExtBatch) error {
	if n.isIncomingP2PDisabled {
		return nil
	}

	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}
	n.network.RespondToBatchRequest(requesterID, batches)
	return nil
}

// ReceiveTransaction is a mock method that simulates receiving a batch from a peer and then forwarding to all subscribers
func (n *MockP2P) ReceiveTransaction(tx common.EncryptedTx) {
	for _, sub := range n.txSubscribers.Subscribers() {
		sub.HandleTransaction(tx)
	}
}

// ReceiveBatches is a mock method that simulates receiving a batch from a peer and then forwarding to all subscribers
func (n *MockP2P) ReceiveBatches(batches []*common.ExtBatch, isLive bool) {
	if n.isIncomingP2PDisabled {
		return
	}

	for _, sub := range n.batchSubscribers.Subscribers() {
		sub.HandleBatches(batches, isLive)
	}
}

// ReceiveBatchRequest is a mock method that simulates receiving a batch request from a peer and then forwarding to all subscribers
func (n *MockP2P) ReceiveBatchRequest(requestID string, fromSeqNo *big.Int) {
	if n.isIncomingP2PDisabled {
		return
	}

	for _, sub := range n.batchReqHandlers.Subscribers() {
		sub.HandleBatchRequest(requestID, fromSeqNo)
	}
}

func (n *MockP2P) RefreshPeerList() {
	// no-op
}
