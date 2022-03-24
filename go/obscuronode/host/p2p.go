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
}

type P2PImpl struct {
	TxAddress      string
	RollupAddress  string
	txListener     net.Listener
	rollupListener net.Listener
}

func (p *P2PImpl) listenForTxs(txP2PCh chan nodecommon.EncryptedTx) {
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

func (p *P2PImpl) listenForRollups(rollupsP2PCh chan obscurocommon.EncodedRollup) {
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

func (p *P2PImpl) stopListeningForTxs() {
	// todo - joel - implement
}

func (p *P2PImpl) stopListeningForRollups() {
	// todo - joel - implement
}

func (p *P2PImpl) handleTx(txP2PCh chan nodecommon.EncryptedTx, listener net.Listener) {
	encryptedTx := readBytes(listener)

	t := nodecommon.L2Tx{}
	// We only post the transaction if it decodes correctly.
	if err := rlp.DecodeBytes(encryptedTx, &t); err == nil {
		txP2PCh <- encryptedTx
	}
}

func (p *P2PImpl) handleRollup(rollupsP2PCh chan obscurocommon.EncodedRollup, listener net.Listener) {
	encodedRollup := readBytes(listener)

	r := nodecommon.Rollup{}
	// We only post the rollup if it decodes correctly.
	if err := rlp.DecodeBytes(encodedRollup, &r); err == nil {
		rollupsP2PCh <- encodedRollup
	}
}

func readBytes(listener net.Listener) []byte {
	conn, err := listener.Accept()
	if err != nil {
		panic("Could not accept any further connections.")
	}
	defer func(conn net.Conn) {
		if err = conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	bytes, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	return bytes
}
