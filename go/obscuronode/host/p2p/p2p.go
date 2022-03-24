package p2p

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// TODO - Provide configurable timeouts on P2P connections.

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
// allTxAddresses is a list of all the transaction P2P addresses on the network, possibly including out own.
// allRollupAddresses is a list of all the rollup P2P addresses on the network, possibly including out own.
// TODO - Consolidate `ourTxAddress` and `ourRollupAddress` into a single address.
// todo - joel - consolidate
func NewP2P(ourTxAddress string, ourRollupAddress string, allTxAddresses []string, allRollupAddresses []string) P2P {
	// We filter out our transaction P2P address if it's contained in the list of all transaction P2P addresses.
	var peerTxAddresses []string
	for _, a := range allTxAddresses {
		if a != ourTxAddress {
			peerTxAddresses = append(peerTxAddresses, a)
		}
	}

	// We filter out our rollup P2P address if it's contained in the list of all rollup P2P addresses.
	var peerRollupAddresses []string
	for _, a := range allRollupAddresses {
		if a != ourRollupAddress {
			peerRollupAddresses = append(peerRollupAddresses, a)
		}
	}

	return &p2pImpl{
		TxAddress:           ourTxAddress,
		RollupAddress:       ourRollupAddress,
		PeerTxAddresses:     peerTxAddresses,
		PeerRollupAddresses: peerRollupAddresses,
	}
}

type p2pImpl struct {
	TxAddress           string
	RollupAddress       string
	PeerTxAddresses     []string
	PeerRollupAddresses []string
	txListener          net.Listener
	rollupListener      net.Listener
}

func (p *p2pImpl) Listen(txP2PCh chan nodecommon.EncryptedTx, rollupsP2PCh chan obscurocommon.EncodedRollup) {
	// We listen for transaction P2P connections.
	txListener, err := net.Listen("tcp", p.TxAddress)
	if err != nil {
		panic(err)
	}
	p.txListener = txListener
	go p.handleTxs(txP2PCh, txListener)

	// We listen for rollup P2P connections.
	rollupListener, err := net.Listen("tcp", p.RollupAddress)
	if err != nil {
		panic(err)
	}
	p.rollupListener = rollupListener
	go p.handleRollups(rollupsP2PCh, rollupListener)
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
	for _, address := range p.PeerTxAddresses {
		sendBytes(address, bytes)
	}
}

func (p *p2pImpl) BroadcastRollup(bytes []byte) {
	for _, address := range p.PeerRollupAddresses {
		sendBytes(address, bytes)
	}
}

func (p *p2pImpl) handleTxs(txP2PCh chan nodecommon.EncryptedTx, listener net.Listener) {
	for {
		encryptedTx := readAllBytes(listener)
		tx := nodecommon.L2Tx{}
		err := rlp.DecodeBytes(encryptedTx, &tx)

		// We only post the transaction if it decodes correctly.
		if err == nil {
			txP2PCh <- encryptedTx
		} else {
			log.Log(fmt.Sprintf("failed to decode transaction received from peer: %v", err))
		}
	}
}

func (p *p2pImpl) handleRollups(rollupsP2PCh chan obscurocommon.EncodedRollup, listener net.Listener) {
	for {
		encodedRollup := readAllBytes(listener)
		r := nodecommon.Rollup{}
		err := rlp.DecodeBytes(readAllBytes(listener), &r)

		// We only post the rollup if it decodes correctly.
		if err == nil {
			rollupsP2PCh <- encodedRollup
		} else {
			log.Log(fmt.Sprintf("failed to decode rollup received from peer: %v", err))
		}
	}
}

// Accepts the next connection, and reads and returns all bytes.
func readAllBytes(listener net.Listener) []byte {
	conn, err := listener.Accept()
	if conn != nil {
		defer func(conn net.Conn) {
			if closeErr := conn.Close(); closeErr != nil {
				panic(closeErr)
			}
		}(conn)
	}
	if err != nil {
		panic("Could not accept any further connections.")
	}

	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	return bytes
}

// Sends the bytes over P2P to the given address.
func sendBytes(address string, tx []byte) {
	// todo - joel - use connection pool
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
