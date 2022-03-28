package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/integration/ethereummock"
)

type createNetworkFunc = func(nrNodes int, avgGossipPeriod uint64, avgBlockDurationUSecs uint64, avgLatency uint64, stats *Stats) ([]*ethereummock.Node, []*host.Node)

// testSimulation encapsulates the shared logic for simulating and testing various types of nodes.
func testSimulation(t *testing.T, avgBlockDurationUSecs uint64, createNetwork createNetworkFunc, efficiencies EfficiencyThresholds) {
	// define core test parameters
	numberOfNodes := 10
	numberOfWallets := 5

	simulationTimeSecs := 15 // in seconds

	avgNetworkLatency := avgBlockDurationUSecs / 15 // artificial latency injected between sending and receiving messages on the mock network
	avgGossipPeriod := avgBlockDurationUSecs / 3    // POBI protocol setting

	// converted to Us
	simulationTimeUSecs := simulationTimeSecs * 1000 * 1000

	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	logFile := setupTestLog("../.build/simulations/")
	defer logFile.Close()

	stats := NewStats(numberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	mockEthNodes, obscuroInMemNodes := createNetwork(numberOfNodes, avgGossipPeriod, avgBlockDurationUSecs, avgNetworkLatency, stats)

	txInjector := NewTransactionInjector(numberOfWallets, avgBlockDurationUSecs, stats, simulationTimeUSecs, mockEthNodes, obscuroInMemNodes)

	simulation := Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		ObscuroNodes:       obscuroInMemNodes, //  the list of in memory obscuro nodes
		AvgBlockDuration:   avgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: simulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkBlockchainValidity(t, &simulation, efficiencies)

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))

	simulation.Stop()
}
