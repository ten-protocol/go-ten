package simulation

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/obscuronode/obscuroclient"

	"github.com/obscuronet/go-obscuro/go/ethclient"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/obscuronet/go-obscuro/go/log"
	"github.com/obscuronet/go-obscuro/go/obscurocommon"
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

	s.WaitForObscuroGenesis()

	// arbitrary sleep to wait for RPC clients to get up and running
	time.Sleep(5 * time.Second)

	timer := time.Now()
	log.Info("Starting injection")
	go s.TxInjector.Start()

	stoppingDelay := s.Params.AvgBlockDuration * 7

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - stoppingDelay)
	log.Info("Stopping injection")

	s.TxInjector.Stop()

	// allow for some time after tx injection was stopped so that the network can process all transactions
	time.Sleep(stoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
}

func (s *Simulation) Stop() {
	// nothing to do for now
}

func (s *Simulation) WaitForObscuroGenesis() {
	// grab an L1 client
	client := s.EthClients[0]

	for {
		// spin through the L1 blocks periodically to see if the genesis rollup has arrived
		head := client.FetchHeadBlock()
		for _, b := range client.BlocksBetween(obscurocommon.GenesisBlock, head) {
			for _, tx := range b.Transactions() {
				t := s.Params.MgmtContractLib.DecodeTx(tx)
				if t == nil {
					continue
				}
				if _, ok := t.(*obscurocommon.L1RollupTx); ok {
					// exit at the first obscuro rollup we see
					return
				}
			}
		}
		time.Sleep(s.Params.AvgBlockDuration)
		log.Trace("Waiting for the Obscuro genesis rollup...")
	}
}
