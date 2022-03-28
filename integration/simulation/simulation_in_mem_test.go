package simulation

import (
	"testing"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "avgBlockDurationUSecs" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might result in inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	// This is a critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	avgBlockDurationUSecs := uint64(40_000) // in u seconds 1 sec = 1e6 usecs.

	createNetwork := CreateBasicNetworkOfInMemoryNodes
	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.32}

	testSimulation(t, avgBlockDurationUSecs, createNetwork, efficiencies)
}
