//nolint:dupl
package simulation

import (
	"github.com/obscuronet/obscuro-playground/go/log"
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"

	"github.com/obscuronet/obscuro-playground/integration"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/erc20contract"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// TestGethSimulation runs the simulation against a private geth network using Clique (PoA)
func TestGethSimulation(t *testing.T) {
	setupTestLog("geth-in-mem")

	log.Info("Starting Geth in-memory simulation.")

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
	netw := network.NewNetworkInMemoryGeth(prefundedWallets, workerWallet, contractsBytes)

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
		StartPort:                 integration.StartPortSimulationGethInMem,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, netw, simParams)
}
