package simulation

import (
	"fmt"
	"time"

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
func NewL2Network(nrNodes int, avgBlockDuration uint64, avgLatency uint64, p2pFactory *p2p.Factory) *L2NetworkCfg {
	// We generate the P2P addresses for each node on the network.
	var nodeAddresses []string
	for i := 1; i <= nrNodes; i++ {
		nodeAddresses = append(nodeAddresses, fmt.Sprintf("localhost:%d", P2P_START_PORT+i))
	}

	p2pNetwork := p2p.NewDelayP2P(
		p2pFactory.NewP2P(fmt.Sprintf("localhost:%d", P2P_START_PORT+100), nodeAddresses),
		func() uint64 {
			return obscurocommon.RndBtw((avgBlockDuration/25)/10, (avgBlockDuration/25)*2)
		})
	p2pNetwork.Listen(nil, nil)

	return &L2NetworkCfg{
		avgLatency:       avgLatency,
		avgBlockDuration: avgBlockDuration,
		p2p:              p2pNetwork,
		nodeAddresses:    nodeAddresses,
	}
}

func (cfg *L2NetworkCfg) BroadcastTx(tx nodecommon.EncryptedTx) {
	cfg.p2p.BroadcastTx(tx)
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
