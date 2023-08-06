package p2p

import (
	"math/big"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/common/subscription"

	"github.com/obscuronet/go-obscuro/go/common/async"

	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"
)

const _sequencerID = "0"

type MockP2PNetwork struct {
	nodes map[string]*MockP2P

	avgLatency       time.Duration
	avgBlockDuration time.Duration
}

func NewMockP2PNetwork(avgBlockDuration time.Duration, avgLatency time.Duration) *MockP2PNetwork {
	return &MockP2PNetwork{
		nodes:            make(map[string]*MockP2P),
		avgBlockDuration: avgBlockDuration,
		avgLatency:       avgLatency,
	}
}

func (m *MockP2PNetwork) P2PServiceFactory(id int) host.ServiceFactory[host.P2PService] {
	idStr := strconv.Itoa(id)
	node := NewMockP2P(idStr, m)
	m.nodes[idStr] = node
	// return a factory function that can be used by the test host to create a p2p service
	return func(_ *config.HostConfig, _ host.ServiceLocator, _ log.Logger) (host.P2PService, error) {
		return node, nil
	}
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

	batchSubscribers *subscription.Manager[hostcommon.P2PBatchHandler]
	txSubscribers    *subscription.Manager[hostcommon.P2PTxHandler]
	batchReqHandlers *subscription.Manager[hostcommon.P2PBatchRequestHandler]

	listenerInterrupt *int32
}

// NewMockP2P returns an instance of a configured L2 Network (no nodes)
func NewMockP2P(id string, network *MockP2PNetwork) *MockP2P {
	i := int32(0)
	return &MockP2P{
		id:      id,
		network: network,

		batchSubscribers:  subscription.NewManager[hostcommon.P2PBatchHandler](),
		txSubscribers:     subscription.NewManager[hostcommon.P2PTxHandler](),
		batchReqHandlers:  subscription.NewManager[hostcommon.P2PBatchRequestHandler](),
		listenerInterrupt: &i,
	}
}

func (n *MockP2P) Start() error {
	// nothing to do here, since communication is direct through the in memory objects
	return nil
}

func (n *MockP2P) Stop() {
	atomic.StoreInt32(n.listenerInterrupt, 1)
}

func (n *MockP2P) HealthStatus() hostcommon.HealthStatus {
	return &hostcommon.BasicErrHealthStatus{ErrMsg: ""}
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
	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}

	n.network.BroadcastBatch(n.id, batches)

	return nil
}

func (n *MockP2P) SubscribeForBatches(handler hostcommon.P2PBatchHandler) func() {
	return n.batchSubscribers.Subscribe(handler)
}

func (n *MockP2P) SubscribeForTx(handler hostcommon.P2PTxHandler) func() {
	return n.txSubscribers.Subscribe(handler)
}

func (n *MockP2P) SubscribeForBatchRequests(handler hostcommon.P2PBatchRequestHandler) func() {
	return n.batchReqHandlers.Subscribe(handler)
}

func (n *MockP2P) RequestBatchesFromSequencer(fromSeqNo *big.Int) error {
	if atomic.LoadInt32(n.listenerInterrupt) == 1 {
		return nil
	}
	n.network.RequestBatchesFromSequencer(n.id, fromSeqNo)
	return nil
}

func (n *MockP2P) RespondToBatchRequest(requesterID string, batches []*common.ExtBatch) error {
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
	for _, sub := range n.batchSubscribers.Subscribers() {
		sub.HandleBatches(batches, isLive)
	}
}

// ReceiveBatchRequest is a mock method that simulates receiving a batch request from a peer and then forwarding to all subscribers
func (n *MockP2P) ReceiveBatchRequest(requestID string, fromSeqNo *big.Int) {
	for _, sub := range n.batchReqHandlers.Subscribers() {
		sub.HandleBatchRequest(requestID, fromSeqNo)
	}
}

func (n *MockP2P) RefreshPeerList() {
	// no-op
}
