package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/ethereummock"
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

// An alias for a function that creates a group of Ethereum and Obscuro nodes.
type createNetworkFunc = func(params SimParams, stats *Stats) ([]*ethereummock.Node, []*host.Node)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, createNetwork createNetworkFunc, params SimParams, efficiencies EfficiencyThresholds) {
	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	logFile := setupTestLog("../.build/simulations/")
	defer logFile.Close()

	stats := NewStats(params.NumberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	mockEthNodes, obscuroNodes := createNetwork(params, stats)

	txInjector := NewTransactionInjector(params.NumberOfWallets, params.AvgBlockDurationUSecs, stats, params.SimulationTimeUSecs, mockEthNodes, obscuroNodes)

	simulation := Simulation{
		MockEthNodes:       mockEthNodes, // the list of mock ethereum nodes
		ObscuroNodes:       obscuroNodes, //  the list of obscuro nodes
		AvgBlockDuration:   params.AvgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: params.SimulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	simulation.Start()
	simulation.Stop()

	// run tests
	checkNetworkValidity(t, &simulation, &params, efficiencies)

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))
}
