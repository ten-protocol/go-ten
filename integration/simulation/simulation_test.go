package simulation

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "avgBlockDurationUSecs" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might result in inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	// define core test parameters
	numberOfNodes := 10
	numberOfWallets := 5

	simulationTimeSecs := 15 // in seconds

	// This is a critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	avgBlockDurationUSecs := uint64(40_000) // in u seconds 1 sec = 1e6 usecs.

	avgNetworkLatency := avgBlockDurationUSecs / 15 // artificial latency injected between sending and receiving messages on the mock network
	avgGossipPeriod := avgBlockDurationUSecs / 3    // POBI protocol setting

	// converted to Us
	simulationTimeUSecs := simulationTimeSecs * 1000 * 1000

	rand.Seed(time.Now().UnixNano())
	uuid.EnableRandPool()

	logFile := setupTestLog("../.build/simulations/")
	defer logFile.Close()

	stats := NewStats(numberOfNodes) // todo - temporary object used to collect metrics. Needs to be replaced with something better

	mockEthNodes, obscuroInMemNodes := CreateBasicNetworkOfInMemoryNodes(numberOfNodes, avgGossipPeriod, avgBlockDurationUSecs, avgNetworkLatency, stats)

	txInjector := NewTransactionInjector(numberOfWallets, avgBlockDurationUSecs, stats, simulationTimeUSecs, mockEthNodes, obscuroInMemNodes)

	simulation := Simulation{
		MockEthNodes:       mockEthNodes,      // the list of mock ethereum nodes
		InMemObscuroNodes:  obscuroInMemNodes, //  the list of in memory obscuro nodes
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
