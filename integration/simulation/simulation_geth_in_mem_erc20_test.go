package simulation

import (
	"testing"
	"time"

	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
	"github.com/obscuronet/obscuro-playground/integration/datagenerator"
	"github.com/obscuronet/obscuro-playground/integration/simulation/network"
	"github.com/obscuronet/obscuro-playground/integration/simulation/params"
)

// TestGethSimulation runs the simulation against a private geth network using Clique (PoA)
func TestGethSimulation(t *testing.T) {
	setupTestLog()

	numberOfNodes := 5

	// randomly create the ethereum wallets to be used and prefund them
	// create one extra wallet as the worker wallet
	wallets := make([]wallet.Wallet, numberOfNodes+1)
	for i := 0; i < numberOfNodes+1; i++ {
		wallets[i] = datagenerator.RandomWallet()
	}

	// The last wallet as the worker wallet ( to deposit and inject transactions )
	workerWallet := wallets[numberOfNodes]

	// define contracts to be deployed
	contractsBytes := []string{
		contracts.MgmtContractByteCode,
		contracts.StableTokenERC20ContractByteCode,
	}

	// define the network to use
	netw := network.NewNetworkInMemoryGeth(wallets, workerWallet, contractsBytes)

	simParams := &params.SimParams{
		NumberOfNodes:             numberOfNodes,
		NumberOfObscuroWallets:    numberOfNodes,
		AvgBlockDuration:          1 * time.Second,
		SimulationTime:            60 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.9,
		L2ToL1EfficiencyThreshold: 0.9, // one rollup every 2 blocks
		EthWallets:                wallets,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	testSimulation(t, netw, simParams)
}
