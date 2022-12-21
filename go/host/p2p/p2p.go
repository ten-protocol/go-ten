package p2p

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

const (
	tcp = "tcp"

	msgTypeTx msgType = iota
	msgTypeBatches
	msgTypeBatchRequest

	_thresholdErrorFailure = 100

	_failedMessageRead        = "msg/inbound/failed_read"
	_failedMessageDecode      = "msg/inbound/failed_decode"
	_failedConnectSendMessage = "msg/outbound/failed_peer_connect"
	_failedWriteSendMessage   = "msg/outbound/failed_write"
	_receivedMessage          = "msg/inbound/success_received"
)

var _alertPeriod = 5 * time.Minute

// A P2P message's type.
type msgType uint8

// Associates an encoded message to its type.
type message struct {
	Sender   string // TODO this needs to be authed in the future
	Type     msgType
	Contents []byte
}

// NewSocketP2PLayer - returns the Socket implementation of the P2P
func NewSocketP2PLayer(config *config.HostConfig, logger gethlog.Logger, metricReg gethmetrics.Registry) host.P2P {
	return &p2pImpl{
		ourAddress:      config.P2PBindAddress,
		peerAddresses:   []string{},
		nodeID:          common.ShortAddress(config.ID),
		p2pTimeout:      config.P2PConnectionTimeout,
		logger:          logger,
		peerTracker:     newPeerTracker(),
		hostGauges:      map[string]map[string]gethmetrics.Gauge{},
		metricsRegistry: metricReg,
	}
}

type p2pImpl struct {
	ourAddress        string
	peerAddresses     []string
	listener          net.Listener
	listenerInterrupt *int32 // A value of 1 indicates that new connections should not be accepted
	nodeID            uint64
	p2pTimeout        time.Duration
	logger            gethlog.Logger
	peerTracker       *peerTracker
	// hostGauges holds a map of gauges per host per event to track p2p metrics and health status
	hostGauges      map[string]map[string]gethmetrics.Gauge
	metricsRegistry gethmetrics.Registry
}

func (p *p2pImpl) StartListening(callback host.Host) {
	// We listen for P2P connections.
	listener, err := net.Listen("tcp", p.ourAddress)
	if err != nil {
		p.logger.Crit(fmt.Sprintf("could not listen for P2P connections on %s.", p.ourAddress), log.ErrKey, err)
	}

	p.logger.Info(fmt.Sprintf("Started listening on port: %s", p.ourAddress))
	i := int32(0)
	p.listenerInterrupt = &i
	p.listener = listener

	go p.handleConnections(callback)
}

func (p *p2pImpl) StopListening() error {
	atomic.StoreInt32(p.listenerInterrupt, 1)

	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}

func (p *p2pImpl) UpdatePeerList(newPeers []string) {
	p.logger.Info(fmt.Sprintf("Updated peer list - old: %s new: %s", p.peerAddresses, newPeers))
	p.peerAddresses = newPeers
}

func (p *p2pImpl) SendTxToSequencer(tx common.EncryptedTx) error {
	msg := message{Sender: p.ourAddress, Type: msgTypeTx, Contents: tx}
	return p.send(msg, p.getSequencer())
}

func (p *p2pImpl) BroadcastBatch(batchMsg *host.BatchMsg) error {
	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourAddress, Type: msgTypeBatches, Contents: encodedBatchMsg}
	return p.broadcast(msg)
}

func (p *p2pImpl) RequestBatchesFromSequencer(batchRequest *common.BatchRequest) error {
	if len(p.peerAddresses) == 0 {
		return errors.New("no peers available to request batches")
	}
	encodedBatchRequest, err := rlp.EncodeToBytes(batchRequest)
	if err != nil {
		return fmt.Errorf("could not encode batch request using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourAddress, Type: msgTypeBatchRequest, Contents: encodedBatchRequest}
	// TODO - #718 - Allow missing batches to be requested from peers other than sequencer?
	return p.send(msg, p.getSequencer())
}

func (p *p2pImpl) SendBatches(batchMsg *host.BatchMsg, to string) error {
	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batches using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourAddress, Type: msgTypeBatches, Contents: encodedBatchMsg}
	return p.send(msg, to)
}

// Status returns the current status of the p2p layer
func (p *p2pImpl) Status() *host.P2PStatus {
	return p.status()
}

// HealthCheck returns whether the p2p is considered healthy
// Currently it considers itself unhealthy
// if there's more than 100 failures on a given fail type
// if there's a known peer for which a message hasn't been received
func (p *p2pImpl) HealthCheck() bool {
	currentStatus := p.status()

	if currentStatus.FailedReceivedMessages >= _thresholdErrorFailure ||
		currentStatus.FailedSendMessage >= _thresholdErrorFailure {
		return false
	}

	var noMsgReceivedPeers []string
	for peer, lastMsgTimestamp := range p.peerTracker.receivedMessagesByPeer() {
		if time.Now().After(lastMsgTimestamp.Add(_alertPeriod)) {
			noMsgReceivedPeers = append(noMsgReceivedPeers, peer)
			p.logger.Warn("no message from peer in the alert period",
				log.HostCmp, p.ourAddress,
				"peer", peer,
				"alertPeriod", _alertPeriod,
			)
		}
	}
	if len(noMsgReceivedPeers) > 0 { //nolint: gosimple
		return false
	}

	return true
}

// Listens for connections and handles them in a separate goroutine.
func (p *p2pImpl) handleConnections(callback host.Host) {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			if atomic.LoadInt32(p.listenerInterrupt) != 1 {
				p.logger.Warn("host could not form P2P connection", log.ErrKey, err)
			}
			return
		}
		go p.handle(conn, callback)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func (p *p2pImpl) handle(conn net.Conn, callback host.Host) {
	if conn != nil {
		defer conn.Close()
	}

	encodedMsg, err := io.ReadAll(conn)
	if err != nil {
		p.logger.Warn("failed to read message from peer", log.ErrKey, err)
		p.incHostGaugeMetric(conn.RemoteAddr().String(), _failedMessageRead)
		return
	}

	msg := message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		p.logger.Warn("failed to decode message received from peer: ", log.ErrKey, err)
		p.incHostGaugeMetric(conn.RemoteAddr().String(), _failedMessageDecode)
		return
	}

	switch msg.Type {
	case msgTypeTx:
		// The transaction is encrypted, so we cannot check that it's correctly formed.
		callback.ReceiveTx(msg.Contents)
	case msgTypeBatches:
		callback.ReceiveBatches(msg.Contents)
	case msgTypeBatchRequest:
		callback.ReceiveBatchRequest(msg.Contents)
	}
	p.incHostGaugeMetric(msg.Sender, _receivedMessage)
	p.peerTracker.receivedPeerMsg(msg.Sender)
}

// Broadcasts a message to all peers.
func (p *p2pImpl) broadcast(msg message) error {
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to peers. Cause: %w", err)
	}

	var wg sync.WaitGroup
	for _, address := range p.peerAddresses {
		wg.Add(1)
		go p.sendBytes(&wg, address, msgEncoded)
	}
	wg.Wait()

	return nil
}

// Sends a message to the provided address.
func (p *p2pImpl) send(msg message, to string) error {
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to sequencer. Cause: %w", err)
	}
	p.sendBytes(nil, to, msgEncoded)
	return nil
}

// Sends the bytes to the provided address.
func (p *p2pImpl) sendBytes(wg *sync.WaitGroup, address string, tx []byte) {
	if wg != nil {
		defer wg.Done()
	}

	conn, err := net.DialTimeout(tcp, address, p.p2pTimeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		p.logger.Warn(fmt.Sprintf("could not connect to peer on address %s", address), log.ErrKey, err)
		p.incHostGaugeMetric(address, _failedConnectSendMessage)
		return
	}

	_, err = conn.Write(tx)
	if err != nil {
		p.logger.Warn(fmt.Sprintf("could not send message to peer on address %s", address), log.ErrKey, err)
		p.incHostGaugeMetric(address, _failedWriteSendMessage)
	}
}

// Retrieves the sequencer's address.
// TODO - #718 - Use better method to identify the sequencer?
func (p *p2pImpl) getSequencer() string {
	return p.peerAddresses[0]
}

// status returns the current status of the p2p layer
func (p *p2pImpl) status() *host.P2PStatus {
	status := &host.P2PStatus{
		FailedReceivedMessages: int64(0),
		FailedSendMessage:      int64(0),
		ReceivedMessages:       int64(0),
	}

	for _, hostGauge := range p.hostGauges {
		for gaugeName, gauge := range hostGauge {
			switch gaugeName {
			case _receivedMessage:
				status.ReceivedMessages = gauge.Value()
			case _failedMessageRead:
			case _failedMessageDecode:
				status.FailedReceivedMessages += gauge.Value()
			case _failedWriteSendMessage:
			case _failedConnectSendMessage:
				status.FailedSendMessage += gauge.Value()
			}
		}
	}
	return status
}

func (p *p2pImpl) incHostGaugeMetric(host string, gaugeName string) {
	if _, ok := p.hostGauges[host]; !ok {
		p.hostGauges[host] = map[string]gethmetrics.Gauge{}
	}
	if _, ok := p.hostGauges[host][gaugeName]; !ok {
		p.hostGauges[host][gaugeName] = gethmetrics.NewRegisteredGauge(gaugeName, p.metricsRegistry)
	}
	p.hostGauges[host][gaugeName].Inc(1)
}
