package host

import (
	"io/ioutil"
	"net"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type P2P interface {
	listenForTxs(txP2PCh chan nodecommon.EncryptedTx)
	listenForRollups(rollupsP2PCh chan obscurocommon.EncodedRollup)
	stopListeningForTxs()
	stopListeningForRollups()

	broadcastRollup(r obscurocommon.EncodedRollup)
}

func NewP2P(idx int, txAddresses []string, rollupAddresses []string) P2P {
	txAddrsCopy := make([]string, len(txAddresses))
	copy(txAddrsCopy, txAddresses)
	peerTxAddresses := append(txAddrsCopy[0:idx], txAddrsCopy[idx+1:]...)

	rollupAddrsCopy := make([]string, len(rollupAddresses))
	copy(rollupAddrsCopy, rollupAddresses)
	peerRollupAddresses := append(rollupAddrsCopy[0:idx], rollupAddrsCopy[idx+1:]...)

	return &p2pImpl{TxAddress: txAddresses[idx], RollupAddress: rollupAddresses[idx], PeerTxAddresses: peerTxAddresses, PeerRollupAddresses: peerRollupAddresses}
}

type p2pImpl struct {
	TxAddress           string
	RollupAddress       string
	PeerTxAddresses     []string
	PeerRollupAddresses []string
	txListener          net.Listener
	rollupListener      net.Listener
}

func (p *p2pImpl) listenForTxs(txP2PCh chan nodecommon.EncryptedTx) {
	listener, err := net.Listen("tcp", p.TxAddress)
	if err != nil {
		panic(err)
	}
	p.txListener = listener

	go func() {
		for {
			p.handleTx(txP2PCh, listener)
		}
	}()
}

func (p *p2pImpl) listenForRollups(rollupsP2PCh chan obscurocommon.EncodedRollup) {
	listener, err := net.Listen("tcp", p.RollupAddress)
	if err != nil {
		panic(err)
	}
	p.rollupListener = listener

	go func() {
		for {
			p.handleRollup(rollupsP2PCh, listener)
		}
	}()
}

func (p *p2pImpl) stopListeningForTxs() {
	// todo - joel - implement
}

func (p *p2pImpl) stopListeningForRollups() {
	// todo - joel - implement
}

func (p *p2pImpl) handleTx(txP2PCh chan nodecommon.EncryptedTx, listener net.Listener) {
	encryptedTx := readBytes(listener)

	t := nodecommon.L2Tx{}
	// We only post the transaction if it decodes correctly.
	if err := rlp.DecodeBytes(encryptedTx, &t); err == nil {
		txP2PCh <- encryptedTx
	}
}

func (p *p2pImpl) handleRollup(rollupsP2PCh chan obscurocommon.EncodedRollup, listener net.Listener) {
	encodedRollup := readBytes(listener)

	r := nodecommon.Rollup{}
	// We only post the rollup if it decodes correctly.
	if err := rlp.DecodeBytes(encodedRollup, &r); err == nil {
		rollupsP2PCh <- encodedRollup
	}
}

func readBytes(listener net.Listener) []byte {
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

// BroadcastRollup Broadcasts the rollup to all L2 peers
func (p *p2pImpl) broadcastRollup(r obscurocommon.EncodedRollup) {
	for _, a := range p.PeerRollupAddresses {
		address := a
		obscurocommon.Schedule(delay(), func() { broadcastBytes(address, r) })
	}
}

// todo - joel - pulled this over from l2 network. need to decide on proper latency
func delay() uint64 {
	avgLatency := uint64(40_000)
	return obscurocommon.RndBtw(avgLatency/10, 2*avgLatency)
}

func broadcastBytes(address string, tx []byte) {
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
