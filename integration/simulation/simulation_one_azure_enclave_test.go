package simulation

import (
	"os"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

const (
	vmIP         = "20.90.208.251" // Todo: replace with the IP of the vm
	azureTestEnv = "AZURE_TEST_ENABLED"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestOneAzureEnclaveNodesMonteCarloSimulation(t *testing.T) {
	if os.Getenv(azureTestEnv) == "" {
		t.Skipf("set %s to run this test", azureTestEnv)
	}
	setupTestLog()

	params := params.SimParams{
		NumberOfNodes:             10,
		NumberOfObscuroWallets:    5,
		AvgBlockDuration:          time.Second,
		SimulationTime:            30 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
	}
	params.AvgNetworkLatency = params.AvgBlockDuration / 15
	params.AvgGossipPeriod = params.AvgBlockDuration / 3

	testSimulation(t, network.NewNetworkWithOneAzureEnclave(vmIP+":11000"), &params)
}
