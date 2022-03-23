package host

import (
	"io/ioutil"
	"net"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type P2P interface {
	listenForTxs(txP2PCh chan nodecommon.EncryptedTx)
	stopListeningForTxs()
}

type P2PImpl struct {
	TxAddress  string
	txListener net.Listener
}

func (p *P2PImpl) listenForTxs(txP2PCh chan nodecommon.EncryptedTx) {
	listener, err := net.Listen("tcp", p.TxAddress)
	if err != nil {
		panic(err)
	}
	p.txListener = listener

	go func() {
		for {
			p.listen(txP2PCh, listener)
		}
	}()
}

func (p *P2PImpl) stopListeningForTxs() {
	// todo - joel - implement
}

func (p *P2PImpl) listen(txP2PCh chan nodecommon.EncryptedTx, listener net.Listener) {
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
