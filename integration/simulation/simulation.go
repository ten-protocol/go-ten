package simulation

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const (
	INITIAL_BALANCE = 5000 // nolint:revive,stylecheck
)

// Simulation represents the data which to set up and run a simulated network
type Simulation struct {
	MockEthNodes       []*ethereum_mock.Node // the list of mock ethereum nodes
	InMemObscuroNodes  []*host.Node          //  the list of in memory obscuro nodes
	AvgBlockDuration   uint64
	TxInjector         *TransactionInjector
	SimulationTimeSecs int
	Stats              *Stats
}

// Start executes the simulation given all the params. Starts all nodes, and injects transactions.
func (s *Simulation) Start() {
	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	// The sequence of starting the nodes is important to catch various edge cases.
	// Here we first start the mock layer 1 nodes, with a pause between them of a fraction of a block duration.
	// The reason is to make sure that they catch up correctly.
	// Then we pause for a while, to give the L1 network enough time to create a number of blocks, which will have to be ingested by the Obscuro nodes
	// Then, we begin the starting sequence of the Obscuro nodes, again with a delay between them, to test that they are able to cach up correctly.
	// Note: Other simulations might test variations of this pattern.
	for _, m := range s.MockEthNodes {
		t := m
		go t.Start()
		time.Sleep(time.Duration(s.AvgBlockDuration / 8))
	}

	time.Sleep(time.Duration(s.AvgBlockDuration * 20))
	for _, m := range s.InMemObscuroNodes {
		t := m
		go t.Start()
		time.Sleep(time.Duration(s.AvgBlockDuration / 3))
	}

	timer := time.Now()
	go s.TxInjector.Start()

	// converted to Us
	simulationTimeUSecs := s.SimulationTimeSecs * 1000 * 1000

	// Wait for the simulation time
	time.Sleep(obscurocommon.Duration(uint64(simulationTimeUSecs)))

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), obscurocommon.Duration(uint64(simulationTimeUSecs)))
	time.Sleep(time.Second)
}

// Stop closes down the L2 and L1 networks.
func (s *Simulation) Stop() {
	// stop L2 first and then L1
	go func() {
		for _, n := range s.InMemObscuroNodes {
			n.Stop()
		}
	}()
	go func() {
		for _, m := range s.MockEthNodes {
			t := m
			go t.Stop()
			// fmt.Printf("Stopped L1 node: %d.\n", m.ID)
		}
	}()
}
