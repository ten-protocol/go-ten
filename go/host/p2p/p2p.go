package p2p

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/host"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	tcp = "tcp"
)

// A P2P message's type.
type msgType uint8

const (
	msgTypeTx msgType = iota
	msgTypeBatch
)

// Message associates an encoded message to its type.
type Message struct {
	Type     msgType
	Contents []byte
}

// NewSocketP2PLayer - returns the Socket implementation of the P2P
func NewSocketP2PLayer(config config.HostConfig, logger gethlog.Logger) host.P2P {
	return &p2pImpl{
		ourAddress:    config.P2PBindAddress,
		peerAddresses: []string{},
		nodeID:        common.ShortAddress(config.ID),
		p2pTimeout:    config.P2PConnectionTimeout,
		logger:        logger,
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

func (p *p2pImpl) BroadcastTx(tx common.EncryptedTx) error {
	msg := Message{Type: msgTypeTx, Contents: tx}
	return p.broadcast(msg)
}

func (p *p2pImpl) BroadcastBatch(batch *common.ExtBatch) error {
	encodedBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	msg := Message{Type: msgTypeBatch, Contents: encodedBatch}
	return p.broadcast(msg)
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
		return
	}

	msg := Message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		p.logger.Warn("failed to decode message received from peer: ", log.ErrKey, err)
		return
	}

	switch msg.Type {
	case msgTypeTx:
		// The transaction is encrypted, so we cannot check that it's correctly formed.
		callback.ReceiveTx(msg.Contents)
	case msgTypeBatch:
		callback.ReceiveBatch(msg.Contents)
	}
}

// Broadcasts a message to all peers.
func (p *p2pImpl) broadcast(msg Message) error {
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

// sendBytes Sends the bytes over P2P to the given address.
func (p *p2pImpl) sendBytes(wg *sync.WaitGroup, address string, tx []byte) {
	defer wg.Done()

	conn, err := net.DialTimeout(tcp, address, p.p2pTimeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		p.logger.Warn(fmt.Sprintf("could not connect to peer on address %s", address), log.ErrKey, err)
		return
	}

	_, err = conn.Write(tx)
	if err != nil {
		p.logger.Warn(fmt.Sprintf("could not send message to peer on address %s", address), log.ErrKey, err)
	}
}
