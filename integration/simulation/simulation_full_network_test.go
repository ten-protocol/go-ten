//nolint:dupl
package simulation

import (
	"github.com/obscuronet/obscuro-playground/go/log"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"

	"github.com/obscuronet/obscuro-playground/integration/datagenerator"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process. The L1 network is a private geth network using Clique (PoA).
func TestFullNetworkMonteCarloSimulation(t *testing.T) {
	setupTestLog("socket")

	log.Info("Starting full network simulation.")

	numberOfNodes := 5
	numberOfSimWallets := 5

	// create the ethereum obsWallets to be used by the nodes and prefund them
	nodeWallets := make([]wallet.Wallet, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		nodeWallets[i] = datagenerator.RandomWallet(integration.EthereumChainID)
	}
	// create the ethereum obsWallets to be used by the simulation and prefund them
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

	// define the network to use
	prefundedWallets := append(append(nodeWallets, simWallets...), workerWallet) //nolint:makezero
	netw := network.NewNetworkOfSocketNodes(prefundedWallets, workerWallet, contractsBytes)

	simParams := &params.SimParams{
		NumberOfNodes:         numberOfNodes,
		AvgBlockDuration:      1 * time.Second,
		SimulationTime:        30 * time.Second,
		L1EfficiencyThreshold: 0.2,
		// Very hard to have precision here as blocks are continually produced and not dependent on the simulation execution thread
		L2EfficiencyThreshold:     0.6, // nodes might produce rollups because they receive a new block
		L2ToL1EfficiencyThreshold: 0.7, // nodes might stop producing rollups but the geth network is still going
		NodeEthWallets:            nodeWallets,
		SimEthWallets:             simWallets,
		StartPort:                 integration.StartPortSimulationFullNetwork,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, netw, simParams)
}
