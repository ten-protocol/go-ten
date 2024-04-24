package p2p

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/config"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
)

const (
	tcp = "tcp"

	msgTypeTx msgType = iota
	msgTypeBatches
	msgTypeBatchRequest
)

var (
	_alertPeriod        = 5 * time.Minute
	errUnknownSequencer = errors.New("sequencer address not known")
)

// A P2P message's type.
type msgType uint8

// Associates an encoded message to its type.
type message struct {
	Sender   string // todo (#1619) - this needs to be authed in the future
	Type     msgType
	Contents []byte
}

type p2pServiceLocator interface {
	L1Publisher() host.L1Publisher
	L2Repo() host.L2BatchRepository
}

// NewSocketP2PLayer - returns the Socket implementation of the P2P
func NewSocketP2PLayer(config *config.HostConfig, serviceLocator p2pServiceLocator, logger gethlog.Logger, metricReg gethmetrics.Registry) *Service {
	return &Service{
		batchSubscribers: subscription.NewManager[host.P2PBatchHandler](),
		txSubscribers:    subscription.NewManager[host.P2PTxHandler](),
		batchReqHandlers: subscription.NewManager[host.P2PBatchRequestHandler](),

		sl: serviceLocator,

		isSequencer:      config.NodeType == common.Sequencer,
		ourBindAddress:   config.P2PBindAddress,
		ourPublicAddress: config.P2PPublicAddress,
		peerAddresses:    []string{},
		p2pTimeout:       config.P2PConnectionTimeout,

		peerAddressesMutex: sync.RWMutex{},

		// monitoring
		peerTracker:     newPeerTracker(),
		metricsRegistry: metricReg,
		logger:          logger,

		isIncomingP2PDisabled: config.IsInboundP2PDisabled,
	}
}

type Service struct {
	batchSubscribers *subscription.Manager[host.P2PBatchHandler]
	txSubscribers    *subscription.Manager[host.P2PTxHandler]
	batchReqHandlers *subscription.Manager[host.P2PBatchRequestHandler]

	listener net.Listener
	running  atomic.Bool // new connections won't be accepted if this is false

	sl p2pServiceLocator

	isSequencer      bool
	ourBindAddress   string
	ourPublicAddress string
	peerAddresses    []string
	p2pTimeout       time.Duration

	peerTracker           *peerTracker
	metricsRegistry       gethmetrics.Registry
	logger                gethlog.Logger
	peerAddressesMutex    sync.RWMutex
	isIncomingP2PDisabled bool
}

func (p *Service) Start() error {
	p.running.Store(true)

	if p.isIncomingP2PDisabled {
		go p.RefreshPeerList()
		return nil
	}

	// We listen for P2P connections.
	listener, err := net.Listen("tcp", p.ourBindAddress)
	if err != nil {
		return fmt.Errorf("could not listen for P2P connections on %s: %w", p.ourBindAddress, err)
	}

	p.logger.Info("P2P server started listening", "bindAddress", p.ourBindAddress, "publicAddress", p.ourPublicAddress)

	p.listener = listener

	go p.handleConnections()

	// ensure we have re-synced the peer list from management contract after startup
	go p.RefreshPeerList()

	return nil
}

func (p *Service) Stop() error {
	p.logger.Info("Shutting down P2P.")
	p.running.Store(false)
	if p.listener != nil {
		// todo immediately shutting down the listener seems to impact other hosts shutdown process
		time.Sleep(time.Second)
		return p.listener.Close()
	}
	return nil
}

func (p *Service) HealthStatus(context.Context) host.HealthStatus {
	msg := ""
	if err := p.verifyHealth(); err != nil {
		msg = err.Error()
	}
	return &host.BasicErrHealthStatus{
		ErrMsg: msg,
	}
}

func (p *Service) SubscribeForBatches(handler host.P2PBatchHandler) func() {
	if p.isIncomingP2PDisabled {
		return nil
	}
	return p.batchSubscribers.Subscribe(handler)
}

func (p *Service) SubscribeForTx(handler host.P2PTxHandler) func() {
	return p.txSubscribers.Subscribe(handler)
}

func (p *Service) SubscribeForBatchRequests(handler host.P2PBatchRequestHandler) func() {
	if p.isIncomingP2PDisabled {
		return nil
	}
	return p.batchReqHandlers.Subscribe(handler)
}

// RefreshPeerList - fetches the latest peer list from L1 and updates the peerAddresses.
// Note: this is designed to be run in a separate goroutine, it will retry a few times before giving up.
func (p *Service) RefreshPeerList() {
	var newPeers []string
	err := retry.Do(func() error {
		if !p.running.Load() {
			return retry.FailFast(fmt.Errorf("p2p service is stopped - abandoning peer list refresh"))
		}

		var retryErr error
		newPeers, retryErr = p.sl.L1Publisher().FetchLatestPeersList()
		if retryErr != nil {
			p.logger.Error("failed to fetch latest peer list from L1", log.ErrKey, retryErr)
			return retryErr
		}
		return nil
	}, retry.NewTimeoutStrategy(1*time.Minute, 5*time.Second))
	if err != nil {
		p.logger.Error("unable to fetch latest peer list from L1", log.ErrKey, err)
		return
	}

	p.peerAddressesMutex.Lock()
	defer p.peerAddressesMutex.Unlock()
	p.logger.Info(fmt.Sprintf("Updated peer list - old: %s new: %s", p.peerAddresses, newPeers))
	p.peerAddresses = newPeers
}

func (p *Service) SendTxToSequencer(tx common.EncryptedTx) error {
	if p.isSequencer {
		return errors.New("sequencer cannot send tx to itself")
	}
	msg := message{Sender: p.ourPublicAddress, Type: msgTypeTx, Contents: tx}
	sequencer, err := p.getSequencer()
	if err != nil {
		return fmt.Errorf("failed to find sequencer - %w", err)
	}
	return p.send(msg, sequencer)
}

func (p *Service) BroadcastBatches(batches []*common.ExtBatch) error {
	if p.isIncomingP2PDisabled {
		return nil
	}
	if !p.isSequencer {
		return errors.New("only sequencer can broadcast batches")
	}
	batchMsg := host.BatchMsg{
		Batches: batches,
		IsLive:  true,
	}

	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourPublicAddress, Type: msgTypeBatches, Contents: encodedBatchMsg}
	return p.broadcast(msg)
}

func (p *Service) RequestBatchesFromSequencer(fromSeqNo *big.Int) error {
	if p.isIncomingP2PDisabled {
		return nil
	}
	if p.isSequencer {
		return errors.New("sequencer cannot request batches from itself")
	}
	batchRequest := &common.BatchRequest{
		Requester: p.ourPublicAddress,
		FromSeqNo: fromSeqNo,
	}
	defer core.LogMethodDuration(p.logger, measure.NewStopwatch(), "Requested batches from sequencer", "fromSeqNo", batchRequest.FromSeqNo)

	encodedBatchRequest, err := rlp.EncodeToBytes(batchRequest)
	if err != nil {
		return fmt.Errorf("could not encode batch request using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourPublicAddress, Type: msgTypeBatchRequest, Contents: encodedBatchRequest}
	// todo (#718) - allow missing batches to be requested from peers other than sequencer?
	sequencer, err := p.getSequencer()
	if err != nil {
		return fmt.Errorf("failed to find sequencer - %w", err)
	}
	return p.send(msg, sequencer)
}

func (p *Service) RespondToBatchRequest(requestID string, batches []*common.ExtBatch) error {
	if p.isIncomingP2PDisabled {
		return nil
	}
	if !p.isSequencer {
		return errors.New("only sequencer can respond to batch requests")
	}
	batchMsg := &host.BatchMsg{
		Batches: batches,
		IsLive:  false,
	}

	encodedBatchMsg, err := rlp.EncodeToBytes(batchMsg)
	if err != nil {
		return fmt.Errorf("could not encode batches using RLP. Cause: %w", err)
	}

	msg := message{Sender: p.ourPublicAddress, Type: msgTypeBatches, Contents: encodedBatchMsg}
	return p.send(msg, requestID)
}

// HealthCheck returns whether the p2p is considered healthy
// Currently it considers itself unhealthy
// if there's more than 100 failures on a given fail type
// if there's a known peer for which a message hasn't been received
func (p *Service) verifyHealth() error {
	if p.isIncomingP2PDisabled {
		return nil
	}
	var noMsgReceivedPeers []string
	for peer, lastMsgTimestamp := range p.peerTracker.receivedMessagesByPeer() {
		if time.Now().After(lastMsgTimestamp.Add(_alertPeriod)) {
			noMsgReceivedPeers = append(noMsgReceivedPeers, peer)
			p.logger.Warn("no message from peer in the alert period",
				"ourAddress", p.ourPublicAddress,
				"peer", peer,
				"alertPeriod", _alertPeriod,
			)
		}
	}
	if len(noMsgReceivedPeers) > 0 {
		return errors.New("no message received from peers")
	}

	return nil
}

// Listens for connections and handles them in a separate goroutine.
func (p *Service) handleConnections() {
	for p.running.Load() {
		conn, err := p.listener.Accept()
		if err != nil {
			if p.running.Load() {
				p.logger.Debug("Could not form P2P connection", log.ErrKey, err)
			}
			return
		}
		go p.handle(conn)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func (p *Service) handle(conn net.Conn) {
	if conn != nil {
		defer conn.Close()
	}

	encodedMsg, err := io.ReadAll(conn)
	if err != nil {
		p.logger.Debug("Failed to read message from peer", log.ErrKey, err)
		return
	}

	msg := message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		p.logger.Debug("Failed to decode message received from peer: ", log.ErrKey, err)
		return
	}

	switch msg.Type {
	case msgTypeTx:
		if !p.isSequencer {
			p.logger.Error("Received transaction from peer, but not a sequencer node")
			return
		}
		// The transaction is encrypted, so we cannot check that it's correctly formed.
		for _, txSubs := range p.txSubscribers.Subscribers() {
			txSubs.HandleTransaction(msg.Contents)
		}
	case msgTypeBatches:
		if p.isSequencer {
			p.logger.Error("received batch from peer, but this is a sequencer node")
			return
		}
		var batchMsg *host.BatchMsg
		err := rlp.DecodeBytes(msg.Contents, &batchMsg)
		if err != nil {
			p.logger.Warn("unable to decode batch received from peer", log.ErrKey, err)
			// nothing to send to subscribers
			break
		}
		// todo - check the batch signature
		for _, batchSubs := range p.batchSubscribers.Subscribers() {
			go batchSubs.HandleBatches(batchMsg.Batches, batchMsg.IsLive)
		}
	case msgTypeBatchRequest:
		if !p.isSequencer {
			p.logger.Error("received batch request from peer, but not a sequencer node")
			return
		}
		// this is an incoming request, p2p service is responsible for finding the response and returning it
		go p.handleBatchRequest(msg.Contents)
	}
	p.peerTracker.receivedPeerMsg(msg.Sender)
}

// Broadcasts a message to all peers.
func (p *Service) broadcast(msg message) error {
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to peers. Cause: %w", err)
	}

	// clone current known addresses
	p.peerAddressesMutex.RLock()
	currentAddresses := make([]string, len(p.peerAddresses))
	copy(currentAddresses, p.peerAddresses)
	p.peerAddressesMutex.RUnlock()

	for _, address := range currentAddresses {
		closureAddr := address
		go func() {
			err := p.sendBytesWithRetry(closureAddr, msgEncoded)
			if err != nil {
				p.logger.Debug("Could not send message to peer", "peer", closureAddr, log.ErrKey, err)
			}
		}()
	}

	return nil
}

// Sends a message to the provided address.
func (p *Service) send(msg message, to string) error {
	// sanity check the message to discover bugs
	if !(msg.Type >= 1 && msg.Type <= 3) {
		p.logger.Error(fmt.Sprintf("Sending message with wrong message type: %v", msg))
	}
	if len(msg.Sender) == 0 {
		p.logger.Error(fmt.Sprintf("Sending message with wrong sender type: %v", msg))
	}
	if len(msg.Contents) == 0 {
		p.logger.Error(fmt.Sprintf("Sending message with empty contents: %v", msg))
	}

	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to sequencer. Cause: %w", err)
	}
	err = p.sendBytesWithRetry(to, msgEncoded)
	if err != nil {
		return err
	}
	return nil
}

// Sends the bytes to the provided address.
// Until introducing libp2p (or equivalent), we have a simple retry
func (p *Service) sendBytesWithRetry(address string, msgEncoded []byte) error {
	// retry for about 2 seconds
	err := retry.Do(func() error {
		return p.sendBytes(address, msgEncoded)
	}, retry.NewDoublingBackoffStrategy(100*time.Millisecond, 5))
	return err
}

// Sends the bytes to the provided address.
func (p *Service) sendBytes(address string, tx []byte) error {
	conn, err := net.DialTimeout(tcp, address, p.p2pTimeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		p.logger.Debug(fmt.Sprintf("could not connect to peer on address %s", address), log.ErrKey, err)
		return err
	}

	_, err = conn.Write(tx)
	if err != nil {
		p.logger.Debug(fmt.Sprintf("could not send message to peer on address %s", address), log.ErrKey, err)
		return err
	}
	return nil
}

// Retrieves the sequencer's address.
// todo (#718) - use better method to identify the sequencer?
func (p *Service) getSequencer() (string, error) {
	p.peerAddressesMutex.RLock()
	defer p.peerAddressesMutex.RUnlock()

	if len(p.peerAddresses) == 0 {
		return "", errUnknownSequencer
	}
	return p.peerAddresses[0], nil
}

func (p *Service) handleBatchRequest(encodedBatchRequest common.EncodedBatchRequest) {
	var batchRequest *common.BatchRequest
	err := rlp.DecodeBytes(encodedBatchRequest, &batchRequest)
	if err != nil {
		p.logger.Warn("unable to decode batch request received from peer using RLP", log.ErrKey, err)
		return
	}

	// todo (@matt) should this response be synchronous?
	for _, requestHandler := range p.batchReqHandlers.Subscribers() {
		go requestHandler.HandleBatchRequest(batchRequest.Requester, batchRequest.FromSeqNo)
	}
}
