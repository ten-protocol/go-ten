package simulation

import (
	"os"
	"testing"
	"time"

	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

const (
	vmIP         = "20.90.208.251" // Todo: replace with the IP of the vm
	azureTestEnv = "AZURE_TEST_ENABLED"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// The genesis node is connected to a remote enclave service running in Azure, while all other enclave services are local.
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestOneAzureEnclaveNodesMonteCarloSimulation(t *testing.T) {
	if os.Getenv(azureTestEnv) == "" {
		t.Skipf("set the variable to run this test: `%s=true`", azureTestEnv)
	}
	setupTestLog("azure-enclave")

	simParams := params.SimParams{
		NumberOfNodes:             10,
		AvgBlockDuration:          time.Second,
		SimulationTime:            30 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
		StartPort:                 integration.StartPortSimulationAzureEnclave,

		MgmtContractLib:  ethereum_mock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereum_mock.NewERC20ContractLibMock(),
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	for i := 0; i < simParams.NumberOfNodes+1; i++ {
		simParams.NodeEthWallets = append(simParams.NodeEthWallets, datagenerator.RandomWallet(integration.ChainID))
		simParams.SimEthWallets = append(simParams.SimEthWallets, datagenerator.RandomWallet(integration.ChainID))
	}

	testSimulation(t, network.NewNetworkWithOneAzureEnclave(vmIP+":11000"), &simParams)
}
