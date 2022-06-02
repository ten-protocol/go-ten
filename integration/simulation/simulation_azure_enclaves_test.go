package simulation

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/obscuronode/wallet"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"
)

const azureTestEnv = "AZURE_TEST_ENABLED"

// TODO: we really need tests to demonstrate the unhappy-cases in the attestation scenario:
//	 - if someone puts a dodgy public key on a request with a genuine attestation report they shouldn't get secret
//	 - if owner doesn't match - they shouldn't get secret

// Todo: replace with the IPs of the VMs you are testing, see the azuredeployer README for more info.
var vmIPs = []string{"20.90.162.69"}

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

	numberOfNodes := 5
	numberOfSimWallets := 5

	// create the ethereum wallets to be used by the nodes and prefund them
	nodeWallets := make([]wallet.Wallet, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		nodeWallets[i] = datagenerator.RandomWallet(integration.EthereumChainID)
	}
	// create the ethereum wallets to be used by the simulation and prefund them
	simWallets := make([]wallet.Wallet, numberOfSimWallets)
	for i := 0; i < numberOfSimWallets; i++ {
		simWallets[i] = datagenerator.RandomWallet(integration.EthereumChainID)
	}
	// create one extra wallet as the worker wallet ( to deploy contracts )
	workerWallet := datagenerator.RandomWallet(integration.EthereumChainID)

	// define contracts to be deployed
	contractsBytes := []string{
		mgmtcontractlib.MgmtContractByteCode,
		erc20contract.ContractByteCode,
	}

	simParams := params.SimParams{
		NumberOfNodes:             numberOfNodes,
		AvgBlockDuration:          time.Second,
		SimulationTime:            60 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
		StartPort:                 integration.StartPortSimulationAzureEnclave,

		NodeEthWallets: nodeWallets,
		SimEthWallets:  simWallets,
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	if len(vmIPs) > simParams.NumberOfNodes {
		panic(fmt.Sprintf("have %d VMs but only %d nodes", len(vmIPs), simParams.NumberOfNodes))
	}

	// define the network to use
	prefundedWallets := append(append(nodeWallets, simWallets...), workerWallet) //nolint:makezero
	netw := network.NewNetworkWithAzureEnclaves(vmIPs, prefundedWallets, workerWallet, contractsBytes)

	testSimulation(t, netw, &simParams)
}
