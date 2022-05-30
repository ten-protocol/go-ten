package simulation

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	ethereum_mock "github.com/obscuronet/obscuro-playground/integration/ethereummock"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

const azureTestEnv = "AZURE_TEST_ENABLED"

// TODO: we really need tests to demonstrate the unhappy-cases in the attestation scenario:
//	 - if someone puts a dodgy public key on a request with a genuine attestation report they shouldn't get secret
//	 - if owner doesn't match - they shouldn't get secret

// Todo: replace with the IPs of the VMs you are testing, see the azuredeployer README for more info.
var vmIPs = []string{"20.254.65.172", "20.254.67.124"}

// This test creates a network of L2 nodes consisting of just the Azure nodes configured above.
//
// It then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestAzureEnclaveNodesMonteCarloSimulation(t *testing.T) {
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

	if len(vmIPs) > simParams.NumberOfNodes {
		panic(fmt.Sprintf("have %d VMs but only %d nodes", len(vmIPs), simParams.NumberOfNodes))
	}

	for i := 0; i < simParams.NumberOfNodes+1; i++ {
		simParams.NodeEthWallets = append(simParams.NodeEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
		simParams.SimEthWallets = append(simParams.SimEthWallets, datagenerator.RandomWallet(integration.EthereumChainID))
	}

	testSimulation(t, network.NewNetworkWithAzureEnclaves(vmIPs), &simParams)
}
