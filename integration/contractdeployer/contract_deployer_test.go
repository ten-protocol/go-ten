package contractdeployer

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/tools/contractdeployer"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	contractDeployerPrivateKey = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
	guessingGameParamOne       = "100"
	guessingGameParamTwo       = "0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
)

func TestCanDeployGuessingGameContract(t *testing.T) {
	createObscuroNetwork(t)

	config := contractdeployer.DefaultConfig()
	config.NodeHost = network.Localhost
	config.NodePort = integration.StartPortContractDeployerTest + network.DefaultHostRPCHTTPOffset
	config.PrivateKey = contractDeployerPrivateKey
	config.ContractName = contractdeployer.GuessingGameContract
	config.ConstructorParams = []string{guessingGameParamOne, guessingGameParamTwo}

	err := contractdeployer.Deploy(config)
	if err != nil {
		panic(err)
	}
}

// Creates a single-node Obscuro network for testing, and deploys an ERC20 contract to it.
func createObscuroNetwork(t *testing.T) {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		AvgGossipPeriod:  1 * time.Second / 3,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        integration.StartPortContractDeployerTest,
	}
	simStats := stats.NewStats(simParams.NumberOfNodes)
	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, simStats)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}
