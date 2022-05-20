package simulation

import (
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

	numberOfNodes := 5

	// randomly create the ethereum wallets to be used and prefund them
	// create one extra wallet as the worker wallet
	wallets := make([]wallet.Wallet, numberOfNodes+1)
	for i := 0; i < numberOfNodes+1; i++ {
		wallets[i] = datagenerator.RandomWallet(integration.ChainID)
	}

	// The last wallet as the worker wallet ( to deposit and inject transactions )
	workerWallet := wallets[numberOfNodes]

	// define contracts to be deployed
	contractsBytes := []string{
		mgmtcontractlib.MgmtContractByteCode,
		erc20contract.ContractByteCode,
	}

	// define the network to use
	netw := network.NewNetworkInMemoryGeth(wallets, workerWallet, contractsBytes)

	simParams := &params.SimParams{
		NumberOfNodes:          numberOfNodes,
		NumberOfObscuroWallets: numberOfNodes,
		AvgBlockDuration:       1 * time.Second,
		SimulationTime:         30 * time.Second,
		L1EfficiencyThreshold:  0.2,
		// Very hard to have precision here as blocks are continually produced and not dependent on the simulation execution thread
		L2EfficiencyThreshold:     0.6, // nodes might produce rollups because they receive a new block
		L2ToL1EfficiencyThreshold: 0.7, // nodes might stop producing rollups but the geth network is still going
		EthWallets:                wallets,
		StartPort:                 integration.StartPortSimulationGethInMem,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, netw, simParams)
}
