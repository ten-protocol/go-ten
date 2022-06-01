package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"

	"github.com/obscuronet/obscuro-playground/integration"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes live in the same process, the enclaves run in individual Docker containers, and the Ethereum nodes are mocked out.
// $> docker rm $(docker stop $(docker ps -a -q --filter ancestor=obscuro_enclave --format="{{.ID}}") will stop and remove all images
func TestDockerNodesMonteCarloSimulation(t *testing.T) {
	setupTestLog("docker")

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
		NumberOfNodes:         numberOfNodes,
		AvgBlockDuration:      1 * time.Second,
		SimulationTime:        30 * time.Second,
		L1EfficiencyThreshold: 0.2,
		// Very hard to have precision here as blocks are continually produced and not dependent on the simulation execution thread
		L2EfficiencyThreshold:     0.6, // nodes might produce rollups because they receive a new block
		L2ToL1EfficiencyThreshold: 0.7, // nodes might stop producing rollups but the geth network is still going
		NodeEthWallets:            nodeWallets,
		SimEthWallets:             simWallets,
		StartPort:                 integration.StartPortSimulationDocker,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 20
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 2

	// define the network to use
	prefundedWallets := append(append(nodeWallets, simWallets...), workerWallet) //nolint:makezero
	netw := network.NewBasicNetworkOfNodesWithDockerEnclave(prefundedWallets, workerWallet, contractsBytes)

	testSimulation(t, netw, &simParams)
}
