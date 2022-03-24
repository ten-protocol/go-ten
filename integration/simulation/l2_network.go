package simulation

import (
	"net"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/p2p"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// L2NetworkCfg - models a full network including artificial random latencies
type L2NetworkCfg struct {
	nodes            []*host.Node
	nodeAddresses    []string
	avgLatency       uint64
	avgBlockDuration uint64
}

// NewL2Network returns an instance of a configured L2 Network (no nodes)
func NewL2Network(avgBlockDuration uint64, avgLatency uint64) *L2NetworkCfg {
	return &L2NetworkCfg{
		avgLatency:       avgLatency,
		avgBlockDuration: avgBlockDuration,
	}
}

func (cfg *L2NetworkCfg) BroadcastTx(tx nodecommon.EncryptedTx) {
	msg := p2p.Message{Type: p2p.Tx, MsgContents: tx}
	msgEncoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}

	for _, a := range cfg.nodeAddresses {
		address := a
		obscurocommon.Schedule(cfg.delay()/2, func() {
			broadcastBytes(address, msgEncoded)
		})
	}
}

func broadcastBytes(address string, tx []byte) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	_, err = conn.Write(tx)
	if err != nil {
		panic(err)
	}
}

// Start kicks off the l2 nodes waiting a delay between each node
func (cfg *L2NetworkCfg) Start(delay time.Duration) {
	// Start l1 nodes
	for _, m := range cfg.nodes {
		t := m
		go t.Start()
		time.Sleep(delay)
	}
}

func (cfg *L2NetworkCfg) Stop() {
	for _, n := range cfg.nodes {
		n.Stop()
	}
}

// delay returns an expected delay on the l2
func (cfg *L2NetworkCfg) delay() uint64 {
	return obscurocommon.RndBtw((cfg.avgBlockDuration/25)/10, (cfg.avgBlockDuration/25)*2)
}
