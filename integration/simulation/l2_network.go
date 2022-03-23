package simulation

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// L2NetworkCfg - models a full network including artificial random latencies
type L2NetworkCfg struct {
	nodes            []*host.Node
	nodeP2PAddresses []string
	enclaveServers   []*grpc.Server
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

// todo - joel - send to each node address instead
// BroadcastRollup Broadcasts the rollup to all L2 peers
func (cfg *L2NetworkCfg) BroadcastRollup(r obscurocommon.EncodedRollup) {
	for _, a := range cfg.nodes {
		rol := nodecommon.DecodeRollupOrPanic(r)
		if a.ID != rol.Header.Agg {
			t := a
			obscurocommon.Schedule(cfg.delay(), func() { t.P2PGossipRollup(r) })
		}
	}
}

func (cfg *L2NetworkCfg) BroadcastTx(tx nodecommon.EncryptedTx) {
	time.Sleep(1 * time.Second) // todo - joel - get rid of this wait somehow

	for _, a := range cfg.nodeP2PAddresses {
		address := a
		obscurocommon.Schedule(cfg.delay()/2, func() {
			broadcastTxToNode(address, tx)
		})
	}
}

func broadcastTxToNode(address string, tx nodecommon.EncryptedTx) {
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

	for _, es := range cfg.enclaveServers {
		es.GracefulStop()
	}
}

// delay returns an expected delay on the l2
func (cfg *L2NetworkCfg) delay() uint64 {
	return obscurocommon.RndBtw(cfg.avgLatency/10, 2*cfg.avgLatency)
}
