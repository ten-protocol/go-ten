package p2p

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/common/subscription"
	"github.com/obscuronet/go-obscuro/go/config"

	gethlog "github.com/ethereum/go-ethereum/log"
)

type NoInboundP2P struct {
	batchSubscribers *subscription.Manager[host.P2PBatchHandler]
	txSubscribers    *subscription.Manager[host.P2PTxHandler]
	batchReqHandlers *subscription.Manager[host.P2PBatchRequestHandler]

	isSequencer      bool
	ourPublicAddress string
	stopControl      *stopcontrol.StopControl
	logger           gethlog.Logger
	sl               p2pServiceLocator
	sequencerAddress string
	p2pTimeout       time.Duration
	ourBindAddress   string

	listener net.Listener
	running  atomic.Bool // new connections won't be accepted if this is false
}

func NewNoInboundP2P(config *config.HostConfig, serviceLocator p2pServiceLocator, logger gethlog.Logger) *NoInboundP2P {
	return &NoInboundP2P{
		batchSubscribers: subscription.NewManager[host.P2PBatchHandler](),
		txSubscribers:    subscription.NewManager[host.P2PTxHandler](),
		batchReqHandlers: subscription.NewManager[host.P2PBatchRequestHandler](),

		isSequencer:      config.NodeType == common.Sequencer,
		ourPublicAddress: config.P2PPublicAddress,
		ourBindAddress:   config.P2PBindAddress,
		stopControl:      stopcontrol.New(),
		logger:           logger,
		sl:               serviceLocator,
		p2pTimeout:       config.P2PConnectionTimeout,
	}
}

func (n *NoInboundP2P) Start() error {
	// Only the sequencer accepts data in
	if !n.isSequencer {
		return nil
	}
	listener, err := net.Listen("tcp", n.ourBindAddress)
	if err != nil {
		return fmt.Errorf("could not listen for P2P connections on %s: %w", n.ourBindAddress, err)
	}

	n.logger.Info("P2P server started listening", "bindAddress", n.ourBindAddress, "publicAddress", n.ourPublicAddress)
	n.running.Store(true)
	n.listener = listener

	go n.handleConnections()

	// ensure we have re-synced the peer list from management contract after startup
	go n.RefreshPeerList()

	return nil
}

func (n *NoInboundP2P) Stop() error {
	return nil
}

func (n *NoInboundP2P) HealthStatus() host.HealthStatus {
	return &host.BasicErrHealthStatus{
		ErrMsg: "",
	}
}

func (n *NoInboundP2P) SubscribeForBatches(handler host.P2PBatchHandler) func() {
	return func() {}
}

func (n *NoInboundP2P) SubscribeForTx(handler host.P2PTxHandler) func() {
	return n.txSubscribers.Subscribe(handler)
}

func (n *NoInboundP2P) SubscribeForBatchRequests(handler host.P2PBatchRequestHandler) func() {
	return func() {}
}

func (n *NoInboundP2P) RefreshPeerList() {
	// we only care about the sequencer
	if n.sequencerAddress != "" {
		return
	}

	var newPeers []string
	err := retry.Do(func() error {
		if n.stopControl.IsStopping() {
			return retry.FailFast(fmt.Errorf("p2p service is stopped - abandoning peer list refresh"))
		}

		var retryErr error
		newPeers, retryErr = n.sl.L1Publisher().FetchLatestPeersList()
		if retryErr != nil {
			n.logger.Error("failed to fetch latest peer list from L1", log.ErrKey, retryErr)
			return retryErr
		}
		return nil
	}, retry.NewTimeoutStrategy(1*time.Minute, 5*time.Second))
	if err != nil {
		n.logger.Error("unable to fetch latest peer list from L1", log.ErrKey, err)
		return
	}

	if len(newPeers) == 0 {
		n.logger.Error("unable to fetch sequencer address from L1", log.ErrKey, err)
		return
	}

	n.sequencerAddress = newPeers[0]
}

func (n *NoInboundP2P) SendTxToSequencer(tx common.EncryptedTx) error {
	if n.isSequencer {
		return errors.New("sequencer cannot send tx to itself")
	}
	msg := message{Sender: n.ourPublicAddress, Type: msgTypeTx, Contents: tx}
	if n.sequencerAddress == "" {
		return fmt.Errorf("failed to find sequencer - no sequencerAddress")
	}
	return n.send(msg, n.sequencerAddress)
}

func (n *NoInboundP2P) BroadcastBatches(_ []*common.ExtBatch) error {
	return nil
}

func (n *NoInboundP2P) RequestBatchesFromSequencer(*big.Int) error {
	return nil
}

func (n *NoInboundP2P) RespondToBatchRequest(_ string, _ []*common.ExtBatch) error {
	return nil
}

// Sends a message to the provided address.
func (n *NoInboundP2P) send(msg message, to string) error {
	// sanity check the message to discover bugs
	if !(msg.Type >= 1 && msg.Type <= 3) {
		n.logger.Error(fmt.Sprintf("Sending message with wrong message type: %v", msg))
	}
	if len(msg.Sender) == 0 {
		n.logger.Error(fmt.Sprintf("Sending message with wrong sender type: %v", msg))
	}
	if len(msg.Contents) == 0 {
		n.logger.Error(fmt.Sprintf("Sending message with empty contents: %v", msg))
	}

	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to sequencer. Cause: %w", err)
	}
	err = n.sendBytesWithRetry(nil, to, msgEncoded)
	if err != nil {
		return err
	}
	return nil
}

// Sends the bytes to the provided address.
// Until introducing libp2p (or equivalent), we have a simple retry
func (n *NoInboundP2P) sendBytesWithRetry(wg *sync.WaitGroup, address string, msgEncoded []byte) error {
	if wg != nil {
		defer wg.Done()
	}
	// retry for about 2 seconds
	err := retry.Do(func() error {
		return n.sendBytes(address, msgEncoded)
	}, retry.NewDoublingBackoffStrategy(100*time.Millisecond, 5))
	return err
}

// Sends the bytes to the provided address.
func (n *NoInboundP2P) sendBytes(address string, tx []byte) error {
	conn, err := net.DialTimeout(tcp, address, n.p2pTimeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		n.logger.Warn(fmt.Sprintf("could not connect to peer on address %s", address), log.ErrKey, err)
		return err
	}

	_, err = conn.Write(tx)
	if err != nil {
		n.logger.Warn(fmt.Sprintf("could not send message to peer on address %s", address), log.ErrKey, err)
		return err
	}
	return nil
}

// Listens for connections and handles them in a separate goroutine.
func (n *NoInboundP2P) handleConnections() {
	for n.running.Load() {
		conn, err := n.listener.Accept()
		if err != nil {
			if n.running.Load() {
				n.logger.Warn("host could not form P2P connection", log.ErrKey, err)
			}
			return
		}
		go n.handle(conn)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func (n *NoInboundP2P) handle(conn net.Conn) {
	if conn != nil {
		defer conn.Close()
	}

	encodedMsg, err := io.ReadAll(conn)
	if err != nil {
		n.logger.Warn("failed to read message from peer", log.ErrKey, err)
		return
	}

	msg := message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		n.logger.Warn("failed to decode message received from peer: ", log.ErrKey, err)
		return
	}

	switch msg.Type {
	case msgTypeTx:
		if !n.isSequencer {
			n.logger.Error("received transaction from peer, but not a sequencer node")
			return
		}
		// The transaction is encrypted, so we cannot check that it's correctly formed.
		for _, txSubs := range n.txSubscribers.Subscribers() {
			txSubs.HandleTransaction(msg.Contents)
		}
	case msgTypeBatches:
		if n.isSequencer {
			n.logger.Error("received batch from peer, but this is a sequencer node")
			return
		}
		var batchMsg *host.BatchMsg
		err := rlp.DecodeBytes(msg.Contents, &batchMsg)
		if err != nil {
			n.logger.Warn("unable to decode batch received from peer", log.ErrKey, err)
			// nothing to send to subscribers
			break
		}

	case msgTypeBatchRequest:
		if !n.isSequencer {
			n.logger.Error("received batch request from peer, but not a sequencer node")
			return
		}
	}
}
