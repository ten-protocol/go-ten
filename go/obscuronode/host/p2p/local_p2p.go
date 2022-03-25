package p2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// networkLayer is a fake multicast network
// it's shared across the p2p package and allows for each instance to receive comms
var networkLayer = make(map[string]chan []byte)

// localP2PImpl implements the P2P in a channel based approach
type localP2PImpl struct {
	address       string
	peerAddresses []string
	txChan        chan nodecommon.EncryptedTx
	rollupChan    chan obscurocommon.EncodedRollup
}

// NewLocalP2P returns a new P2P object that runs locally using channels.
func NewLocalP2P(ourAddress string, allAddresses []string) P2P {
	// We filter out our P2P address if it's contained in the list of all P2P addresses.
	var peerAddresses []string
	for _, a := range allAddresses {
		if a != ourAddress {
			peerAddresses = append(peerAddresses, a)
		}
	}

	// ensure the current client network layer has been created
	if _, found := networkLayer[ourAddress]; !found {
		networkLayer[ourAddress] = make(chan []byte, 1000)
	}

	return &localP2PImpl{
		address:       ourAddress,
		peerAddresses: peerAddresses,
	}
}

func (l *localP2PImpl) Listen(txChan chan nodecommon.EncryptedTx, rollupChan chan obscurocommon.EncodedRollup) {
	l.txChan = txChan
	l.rollupChan = rollupChan
	listener := networkLayer[l.address]
	go func() {
		for {
			data, ok := <-listener
			if !ok {
				return
			}
			l.handle(data)
		}
	}()
}

// StopListening closes the channels
func (l *localP2PImpl) StopListening() {
	close(l.txChan)
	close(l.rollupChan)
	close(networkLayer[l.address])
}

func (l *localP2PImpl) BroadcastTx(bytes []byte) {
	l.broadcast(Tx, bytes)
}

func (l *localP2PImpl) BroadcastRollup(bytes []byte) {
	l.broadcast(Rollup, bytes)
}

// sendBytes issues bytes across the network
func (l *localP2PImpl) sendBytes(address string, data []byte) {
	networkLayer[address] <- data
}

// Creates a P2P message and broadcasts it to all peers.
func (l *localP2PImpl) broadcast(msgType Type, bytes []byte) {
	msg := Message{Type: msgType, MsgContents: bytes}
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}

	for _, a := range l.peerAddresses {
		address := a
		l.sendBytes(address, msgEncoded)
	}
}

// Receives and decodes a P2P message, and pushes it to the correct channel.
func (l *localP2PImpl) handle(data []byte) {
	msg := Message{}
	err := rlp.DecodeBytes(data, &msg)
	if err != nil {
		panic(err)
	}

	switch msg.Type {
	case Tx:
		tx := nodecommon.L2Tx{}
		err := rlp.DecodeBytes(msg.MsgContents, &tx)
		// We only post the transaction if it decodes correctly.
		if err != nil {
			log.Log(fmt.Sprintf("failed to decode transaction received from peer: %v", err))
			return
		}
		l.txChan <- msg.MsgContents

	case Rollup:
		rollup := nodecommon.Rollup{}
		err := rlp.DecodeBytes(msg.MsgContents, &rollup)
		// We only post the rollup if it decodes correctly.
		if err != nil {
			log.Log(fmt.Sprintf("failed to decode rollup received from peer: %v", err))
			return
		}
		l.rollupChan <- msg.MsgContents
	default:
		panic(msg)
	}
}
