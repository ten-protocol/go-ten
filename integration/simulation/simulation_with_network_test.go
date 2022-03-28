package simulation

import (
	"github.com/google/uuid"
	"math/rand"
	"testing"
	"time"
)

// This test creates a network of standalone L2 nodes, then injects transactions, and finally checks the resulting output blockchain
func TestStandaloneNodesMonteCarloSimulation(t *testing.T) {
	// todo - joel - pull some of this setup into a common method
	// define core test parameters
	numberOfNodes := 10
	numberOfWallets := 5

	simulationTimeSecs := 15 // in seconds

	// This is a critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	avgBlockDurationUSecs := uint64(100_000) // in u seconds 1 sec = 1e6 usecs.

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

	// todo - joel - change `InMemObscuroNodes` to just `ObscuroNodes`
	simulation := Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		InMemObscuroNodes:  obscuroInMemNodes, // the list of in memory obscuro nodes
		AvgBlockDuration:   avgBlockDurationUSecs,
		TxInjector:         txInjector,
		SimulationTimeSecs: simulationTimeSecs,
		Stats:              stats,
	}

	// execute the simulation
	simulation.Start()

	// run tests
	checkBlockchainValidity(t, &simulation)

	// generate and print the final stats
	t.Logf("Simulation results:%+v", NewOutputStats(&simulation))

	simulation.Stop()
}
