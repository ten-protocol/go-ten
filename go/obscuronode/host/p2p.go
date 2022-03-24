package host

import (
	"io/ioutil"
	"net"

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
	conn, err := listener.Accept()
	if err != nil {
		println("Could not accept any further connections.")
	}
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	encryptedTx, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	txP2PCh <- encryptedTx
}

func (p *P2PImpl) handleRollup(rollupsP2PCh chan obscurocommon.EncodedRollup, listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		println("Could not accept any further connections.")
	}
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	encryptedTx, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	rollupsP2PCh <- encryptedTx
}
