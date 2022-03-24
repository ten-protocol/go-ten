package simulation

import (
	"fmt"
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
	p2p              p2p.P2P
}

// NewL2Network returns an instance of a configured L2 Network (no nodes)
func NewL2Network(
	nrNodes int,
	avgBlockDuration uint64,
	avgLatency uint64,
	newP2P func(ourAddress string, allAddresses []string) p2p.P2P,
) *L2NetworkCfg {
	// We generate the P2P addresses for each node on the network.
	var nodeAddresses []string
	for i := 1; i <= nrNodes; i++ {
		nodeAddresses = append(nodeAddresses, fmt.Sprintf("localhost:%d", P2P_START_PORT+i))
	}

	p2p := newP2P("localhost:11000", nodeAddresses)
	p2p.Listen(nil, nil)

	return &L2NetworkCfg{
		avgLatency:       avgLatency,
		avgBlockDuration: avgBlockDuration,
		p2p:              p2p,
		nodeAddresses:    nodeAddresses,
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
		// we want to control the delay, so we use the send function
		obscurocommon.Schedule(cfg.delay()/2, func() {
			cfg.p2p.SendBytes(address, msgEncoded)
		})
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
