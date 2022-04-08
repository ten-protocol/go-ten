package simulation

import (
	"fmt"
	"net"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const (
	INITIAL_BALANCE = 5000 // nolint:revive,stylecheck
)

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	EthClients       []ethclient.Client // the list of mock ethereum clients
	ObscuroNodes     []*host.Node       // the list of Obscuro nodes - todo - need to be interfaces to rpc handles
	ObscuroP2PAddrs  []string           // the P2P addresses of the Obscuro nodes
	AvgBlockDuration uint64
	TxInjector       *TransactionInjector
	SimulationTime   time.Duration
	Stats            *stats.Stats
	Params           *params.SimParams
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	// TODO - Remove this waiting period. The ability for nodes to catch up should be part of the tests.
	waitForP2p(s.ObscuroP2PAddrs)

	timer := time.Now()
	go s.TxInjector.Start()

	// Wait for the simulation time
	time.Sleep(s.SimulationTime)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
	time.Sleep(time.Second)
}

func (s *Simulation) Stop() {
	// nothing to do for now
}

// Waits for the L2 nodes to be ready to process P2P messages.
func waitForP2p(obscuroP2PAddrs []string) {
	for _, addr := range obscuroP2PAddrs {
		for {
			conn, _ := net.Dial("tcp", addr)
			if conn != nil {
				if closeErr := conn.Close(); closeErr != nil {
					panic(closeErr)
				}
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
