package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

// This test creates a network of standalone L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate via sockets, and with standalone enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestStandaloneNodesMonteCarloSimulation(t *testing.T) {
	// todo - joel - pull some of this setup into a common method
	// define core test parameters
	numberOfNodes := 10
	numberOfWallets := 5

	simulationTimeSecs := 15 // in seconds

	// This is a critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	// todo - joel - too high?
	avgBlockDurationUSecs := uint64(250_000) // in u seconds 1 sec = 1e6 usecs.

	avgNetworkLatency := avgBlockDurationUSecs / 15 // artificial latency injected between sending and receiving messages on the mock network
	avgGossipPeriod := avgBlockDurationUSecs / 3    // POBI protocol setting

	// converted to Us
	simulationTimeUSecs := simulationTimeSecs * 1000 * 1000

	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	logFile := setupTestLog("../.build/simulations/")
	defer logFile.Close()

	stats := NewStats(numberOfNodes)

	mockEthNodes, obscuroInMemNodes := CreateBasicNetworkOfStandaloneNodes(numberOfNodes, avgGossipPeriod, avgBlockDurationUSecs, avgNetworkLatency, stats)

	txInjector := NewTransactionInjector(numberOfWallets, avgBlockDurationUSecs, stats, simulationTimeUSecs, mockEthNodes, obscuroInMemNodes)

	// todo - joel - change `ObscuroNodes` to just `ObscuroNodes`
	simulation := Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		ObscuroNodes:       obscuroInMemNodes, // the list of in memory obscuro nodes
		AvgBlockDuration:   avgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: simulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	// todo - joel - approx. correct values?
	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.38}
	checkBlockchainValidity(t, &simulation, efficiencies)

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))

	simulation.Stop()
}
