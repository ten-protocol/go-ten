package p2p

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
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

// P2P manages P2P communication between L2 nodes.
type P2P interface {
	// Listen starts listening for transaction and rollup P2P connections.
	Listen(chan nodecommon.EncryptedTx, chan obscurocommon.EncodedRollup)
	// StopListening stops listening for transaction and rollup P2P connections.
	StopListening()

	// BroadcastTx broadcasts a transaction to all network peers over P2P.
	BroadcastTx([]byte)
	// BroadcastRollup broadcasts a rollup to all network peers over P2P.
	BroadcastRollup([]byte)
}

// NewP2P returns a new P2P object.
// allAddresses is a list of all the transaction P2P addresses on the network, possibly including ourAddress.
func NewP2P(ourAddress string, allAddresses []string) P2P {
	// We filter out our P2P address if it's contained in the list of all P2P addresses.
	var peerAddresses []string
	for _, a := range allAddresses {
		if a != ourAddress {
			peerAddresses = append(peerAddresses, a)
		}
	}

	return &p2pImpl{
		OurAddress:    ourAddress,
		PeerAddresses: peerAddresses,
	}
}

type p2pImpl struct {
	OurAddress     string
	PeerAddresses  []string
	txListener     net.Listener
	rollupListener net.Listener
}

func (p *p2pImpl) Listen(txP2PCh chan nodecommon.EncryptedTx, rollupsP2PCh chan obscurocommon.EncodedRollup) {
	// We listen for P2P connections.
	txListener, err := net.Listen("tcp", p.OurAddress)
	if err != nil {
		panic(err)
	}
	p.txListener = txListener
	go p.handleConnections(txP2PCh, rollupsP2PCh, txListener)
}

func (p *p2pImpl) StopListening() {
	if p.txListener != nil {
		if err := p.txListener.Close(); err != nil {
			log.Log(fmt.Sprintf("failed to close transaction P2P listener cleanly: %v", err))
		}
		p.txListener = nil
	}

	if p.rollupListener != nil {
		if err := p.rollupListener.Close(); err != nil {
			log.Log(fmt.Sprintf("failed to close rollup P2P listener cleanly: %v", err))
		}
		p.rollupListener = nil
	}
}

func (p *p2pImpl) BroadcastTx(bytes []byte) {
	p.broadcast(Tx, bytes)
}

func (p *p2pImpl) BroadcastRollup(bytes []byte) {
	p.broadcast(Rollup, bytes)
}

// Listens for connections and handles them in a separate goroutine.
func (p *p2pImpl) handleConnections(txP2PCh chan nodecommon.EncryptedTx, rollupsP2PCh chan obscurocommon.EncodedRollup, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Could not accept any further connections.")
		}
		go handle(conn, txP2PCh, rollupsP2PCh)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func handle(conn net.Conn, txP2PCh chan nodecommon.EncryptedTx, rollupsP2PCh chan obscurocommon.EncodedRollup) {
	if conn != nil {
		defer func(conn net.Conn) {
			if closeErr := conn.Close(); closeErr != nil {
				panic(closeErr)
			}
		}(conn)
	}

	encodedMsg, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	msg := Message{}
	err = rlp.DecodeBytes(encodedMsg, &msg)
	if err != nil {
		panic(err)
	}

	switch msg.Type {
	case Tx:
		tx := nodecommon.L2Tx{}
		err := rlp.DecodeBytes(msg.MsgContents, &tx)

		// We only post the transaction if it decodes correctly.
		if err == nil {
			txP2PCh <- msg.MsgContents
		} else {
			log.Log(fmt.Sprintf("failed to decode transaction received from peer: %v", err))
		}
	case Rollup:
		rollup := nodecommon.Rollup{}
		err := rlp.DecodeBytes(msg.MsgContents, &rollup)

		// We only post the rollup if it decodes correctly.
		if err == nil {
			rollupsP2PCh <- msg.MsgContents
		} else {
			log.Log(fmt.Sprintf("failed to decode rollup received from peer: %v", err))
		}
	}
}

// Creates a P2P message and broadcasts it to all peers.
func (p *p2pImpl) broadcast(msgType Type, bytes []byte) {
	msg := Message{Type: msgType, MsgContents: bytes}
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}

	for _, address := range p.PeerAddresses {
		sendBytes(address, msgEncoded)
	}
}

// Sends the bytes over P2P to the given address.
func sendBytes(address string, tx []byte) {
	conn, err := net.Dial("tcp", address)
	if conn != nil {
		defer func(conn net.Conn) {
			if closeErr := conn.Close(); closeErr != nil {
				panic(closeErr)
			}
		}(conn)
	}
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(tx)
	if err != nil {
		panic(err)
	}
}
