package simulation

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const initialBalance = 5000

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	EthClients       []ethclient.EthClient   // the list of mock ethereum clients
	ObscuroClients   []*obscuroclient.Client // the list of Obscuro host clients
	ObscuroP2PAddrs  []string                // the P2P addresses of the Obscuro nodes
	AvgBlockDuration uint64
	TxInjector       *TransactionInjector
	SimulationTime   time.Duration
	Stats            *stats.Stats
	Params           *params.SimParams
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	log.Info(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	timer := time.Now()
	fmt.Printf("Starting injection\n")
	go s.TxInjector.Start()

	stoppingDelay := s.Params.AvgBlockDuration * 4

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - stoppingDelay)
	fmt.Printf("Stopping injection\n")

	s.TxInjector.Stop()

	// allow for some time after tx injection was stopped so that the network can process all transactions
	time.Sleep(stoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
}

func (s *Simulation) Stop() {
	// nothing to do for now
}
