package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/exec"
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
}

// An alias for a function that returns a group of Ethereum nodes, a group of Obscuro nodes, and the Obscuro nodes' P2P addresses.
type createNetworkFunc = func(params SimParams, stats *Stats) ([]exec.EthNode, []*host.Node, []string)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, createNetwork createNetworkFunc, params SimParams, efficiencies EfficiencyThresholds) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	stats := NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	ethNodes, obscuroNodes, p2pAddrs := createNetwork(params, stats)

	txInjector := NewTransactionInjector(
		params.NumberOfWallets,
		params.AvgBlockDurationUSecs,
		stats,
		params.SimulationTimeUSecs,
		ethNodes,
		obscuroNodes,
	)

	simulation := Simulation{
		EthNodes:           ethNodes, // the list of ethereum nodes
		ObscuroNodes:       obscuroNodes,
		ObscuroP2PAddrs:    p2pAddrs,
		AvgBlockDuration:   params.AvgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: params.SimulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkNetworkValidity(t, &simulation, &params, efficiencies)

	simulation.Stop()

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
}
