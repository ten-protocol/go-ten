package simulation

import (
	"fmt"
	"net"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const (
	INITIAL_BALANCE = 5000 // nolint:revive,stylecheck
)

// SimParams are the parameters for setting up the simulation.
type SimParams struct {
	NumberOfNodes   int
	NumberOfWallets int

	// A critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	AvgBlockDurationUSecs uint64
	AvgNetworkLatency     uint64 // artificial latency injected between sending and receiving messages on the mock network
	AvgGossipPeriod       uint64 // POBI protocol setting

	SimulationTimeSecs  int // in seconds
	SimulationTimeUSecs int // SimulationTimeSecs converted to Us

	// EfficiencyThresholds represents an acceptable "dead blocks" percentage for this simulation.
	// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
	// We test the results against this threshold to catch eventual protocol errors.
	L1EfficiencyThreshold     float64
	L2EfficiencyThreshold     float64
	L2ToL1EfficiencyThreshold float64
}

// This interface is responsible with knowing how to manage the lifecycle of networks of Ethereum or Obscuro nodes.
// These networks can be composed of in-memory go-routines or of fully fledged existing nodes like Ropsten.
// Implementation notes:
// - This is a work in progress, so there is a lot of code duplication in the implementations
// - Once we implement a few more versions: for example using Ganache, or using enclaves running in azure, etc, we'll revisit and create better abstractions.
type Network interface {
	// Create - returns a group of started Ethereum nodes, a group of started Obscuro nodes, and the Obscuro nodes' P2P addresses.
	// todo - return interfaces to RPC handles to the nodes
	Create(params SimParams, stats *Stats) ([]*ethereum_mock.Node, []*host.Node, []string)
	TearDown()
}

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	MockEthNodes       []*ethereum_mock.Node // the list of mock ethereum nodes - todo - need to be interfaces to rpc handles
	ObscuroNodes       []*host.Node          // the list of Obscuro nodes - todo - need to be interfaces to rpc handles
	ObscuroP2PAddrs    []string              // the P2P addresses of the Obscuro nodes
	AvgBlockDuration   uint64
	TxInjector         *TransactionInjector
	SimulationTimeSecs int
	Stats              *Stats
	Params             *SimParams
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	log.Log(fmt.Sprintf("Genesis block: b_%d.", obscurocommon.ShortHash(obscurocommon.GenesisBlock.Hash())))

	// TODO - Remove this waiting period. The ability for nodes to catch up should be part of the tests.
	waitForP2p(s.ObscuroP2PAddrs)

	timer := time.Now()
	go s.TxInjector.Start()

	// converted to Us
	simulationTimeUSecs := s.SimulationTimeSecs * 1000 * 1000

	// Wait for the simulation time
	time.Sleep(obscurocommon.Duration(uint64(simulationTimeUSecs)))

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), obscurocommon.Duration(uint64(simulationTimeUSecs)))
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
