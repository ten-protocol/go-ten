package simulation

import (
	"testing"
)

// This test creates a network of standalone L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate via sockets, and with standalone enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestStandaloneNodesMonteCarloSimulation(t *testing.T) {
	// This is a critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	avgBlockDurationUSecs := uint64(250_000) // in u seconds 1 sec = 1e6 usecs.

	createNetwork := CreateBasicNetworkOfStandaloneNodes
	efficiencies := EfficiencyThresholds{0.2, 0.3, 0.4}

	testSimulation(t, avgBlockDurationUSecs, createNetwork, efficiencies)
}
