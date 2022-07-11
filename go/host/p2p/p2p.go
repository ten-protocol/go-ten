package p2p

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync/atomic"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/obscuronet/obscuro-playground/go/config"

	"github.com/obscuronet/obscuro-playground/go/host"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
)

// TODO - Provide configurable timeouts on P2P connections.

// Type indicates the type of a P2P message.
type Type uint8

const (
	Tx Type = iota
	Rollup
)

// Message associates an encoded message to its type.
type Message struct {
	Type        Type
	MsgContents []byte
}

// NewSocketP2PLayer - returns the Socket implementation of the P2P
func NewSocketP2PLayer(config config.HostConfig) host.P2P {
	return &p2pImpl{
		ourAddress:    config.P2PAddress,
		peerAddresses: []string{},
		nodeID:        common.ShortAddress(config.ID),
	}
}

type p2pImpl struct {
	ourAddress        string
	peerAddresses     []string
	listener          net.Listener
	listenerInterrupt *int32 // A value of 1 indicates that new connections should not be accepted
	nodeID            uint64
}

func (p *p2pImpl) StartListening(callback host.P2PCallback) {
	// We listen for P2P connections.
	listener, err := net.Listen("tcp", p.ourAddress)
	if err != nil {
		log.Panic("could not listen for P2P connections on %s. Cause: %s", p.ourAddress, err)
	}

	common.LogWithID(p.nodeID, "Start listening on port: %s", p.ourAddress)
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
	p.peerAddresses = newPeers
}

func (p *p2pImpl) BroadcastTx(tx common.EncryptedTx) error {
	return p.broadcast(Tx, tx, p.peerAddresses)
}

func (p *p2pImpl) BroadcastRollup(r common.EncodedRollup) error {
	return p.broadcast(Rollup, r, p.peerAddresses)
}

// Listens for connections and handles them in a separate goroutine.
func (p *p2pImpl) handleConnections(callback host.P2PCallback) {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			if atomic.LoadInt32(p.listenerInterrupt) != 1 {
				common.WarnWithID(p.nodeID, "host could not from P2P connection: %s", err)
			}
			return
		}
		go p.handle(conn, callback)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func (p *p2pImpl) handle(conn net.Conn, callback host.P2PCallback) {
	if conn != nil {
		defer conn.Close()
	}

	encodedMsg, err := ioutil.ReadAll(conn)
	if err != nil {
		common.WarnWithID(p.nodeID, "failed to read message from peer: %v", err)
		return
	}

	msg := Message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		common.WarnWithID(p.nodeID, "failed to decode message received from peer: %v", err)
		return
	}

	switch msg.Type {
	case Tx:
		// The transaction is encrypted, so we cannot check that it's correctly formed.
		callback.ReceiveTx(msg.MsgContents)
	case Rollup:
		// We check that the rollup decodes correctly.
		rollup := common.EncryptedRollup{}
		if err = rlp.DecodeBytes(msg.MsgContents, &rollup); err != nil {
			common.WarnWithID(p.nodeID, "failed to decode rollup received from peer: %v", err)
			return
		}

		callback.ReceiveRollup(msg.MsgContents)
	}
}

// Creates a P2P message and broadcasts it to all peers.
func (p *p2pImpl) broadcast(msgType Type, bytes []byte, toAddresses []string) error {
	msg := Message{Type: msgType, MsgContents: bytes}
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return fmt.Errorf("could not encode message to send to peers. Cause: %w", err)
	}

	for _, address := range toAddresses {
		p.sendBytes(address, msgEncoded)
	}

	return nil
}

// sendBytes Sends the bytes over P2P to the given address.
func (p *p2pImpl) sendBytes(address string, tx []byte) {
	conn, err := net.Dial("tcp", address)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		common.WarnWithID(p.nodeID, "could not send message to peer on address %s: %v", address, err)
		return
	}

	_, err = conn.Write(tx)
	if err != nil {
		common.WarnWithID(p.nodeID, "could not send message to peer on address %s: %v", address, err)
	}
}
