package contractdeployer

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"
	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/tools/contractdeployer"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	contractDeployerPrivateKeyHex = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
	guessingGameParamOne          = "100"
	guessingGameParamTwo          = "0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
	latestBlock                   = "latest"
	emptyCode                     = "0x"
)

var (
	config = &contractdeployer.Config{
		NodeHost:          network.Localhost,
		NodePort:          integration.StartPortContractDeployerTest + network.DefaultHostRPCHTTPOffset,
		IsL1Deployment:    false,
		PrivateKey:        contractDeployerPrivateKeyHex,
		ChainID:           big.NewInt(integration.ObscuroChainID),
		ContractName:      contractdeployer.GuessingGameContract,
		ConstructorParams: []string{guessingGameParamOne, guessingGameParamTwo},
	}
	nodeAddress = fmt.Sprintf("%s:%d", config.NodeHost, config.NodePort)
)

func TestCanDeployGuessingGameContract(t *testing.T) {
	createObscuroNetwork(t)
	// This sleep is required to ensure the initial rollup exists, and thus contract deployer can check its balance.
	time.Sleep(2 * time.Second)
	contractAddr, err := contractdeployer.Deploy(config)
	if err != nil {
		panic(err)
	}

	contractDeployerWallet := getWallet(contractDeployerPrivateKeyHex)
	contractDeployerClient := getClient(contractDeployerWallet)

	var deployedCode string
	err = contractDeployerClient.Call(&deployedCode, rpcclientlib.RPCGetCode, contractAddr, latestBlock)
	if err != nil {
		panic(err)
	}

	if deployedCode == emptyCode {
		t.Fatal("contract was deployed but could not get code")
	}
}

func TestFaucetSendsFundsOnlyIfNeeded(t *testing.T) {
	createObscuroNetwork(t)

	faucetWallet := getWallet(rollupchain.FaucetPrivateKeyHex)
	faucetClient := getClient(faucetWallet)

	// We send more than enough to the contract deployer, to make sure prefunding won't be needed.
	prefundContractDeployer(faucetWallet, faucetClient, contractdeployer.Prealloc*3)

	// We check the faucet's balance before and after the deployment. Since the contract deployer has already been sent
	// sufficient funds, the faucet should have been to dispense any more, leaving its balance unchanged.
	var faucetInitialBalance string
	err := faucetClient.Call(&faucetInitialBalance, rpcclientlib.RPCGetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	_, err = contractdeployer.Deploy(config)
	if err != nil {
		panic(err)
	}

	var faucetBalanceAfterDeploy string
	// We create a new faucet client because deploying the contract will have overwritten the faucet's viewing key on the node.
	faucetClient = getClient(faucetWallet)
	err = faucetClient.Call(&faucetBalanceAfterDeploy, rpcclientlib.RPCGetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	if faucetInitialBalance != faucetBalanceAfterDeploy {
		t.Fatal("contract deployment allocated extra funds to contract deployer, despite sufficient funds")
	}
}

func getWallet(privateKeyHex string) wallet.Wallet {
	faucetPrivKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic("could not initialise faucet private key")
	}
	faucetWallet := wallet.NewInMemoryWalletFromPK(config.ChainID, faucetPrivKey)
	return faucetWallet
}

// Creates a single-node Obscuro network for testing.
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

// Returns a viewing-key client with a registered viewing key.
func getClient(wallet wallet.Wallet) *rpcclientlib.EncRPCClient {
	viewingKey, err := rpcclientlib.GenerateAndSignViewingKey(wallet)
	if err != nil {
		panic(err)
	}
	client, err := rpcclientlib.NewEncNetworkClient(nodeAddress, viewingKey)
	if err != nil {
		panic(err)
	}
	return client
}

// Prefunds the contract deployer with the given amount.
func prefundContractDeployer(faucetWallet wallet.Wallet, faucetClient *rpcclientlib.EncRPCClient, prealloc int64) {
	contractDeployerPrivKey, err := crypto.HexToECDSA(contractDeployerPrivateKeyHex)
	if err != nil {
		panic("could not initialise contract deployer private key")
	}
	contractDeployerWallet := wallet.NewInMemoryWalletFromPK(config.ChainID, contractDeployerPrivKey)

	contractDeployerAddr := contractDeployerWallet.Address()
	tx := &types.LegacyTx{
		Nonce:    0, // We can assume this is the first transaction for the new wallet.
		Value:    big.NewInt(prealloc),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		To:       &contractDeployerAddr,
	}
	signedTx, err := faucetWallet.SignTransaction(tx)
	if err != nil {
		panic(err)
	}

	err = faucetClient.Call(nil, rpcclientlib.RPCSendRawTransaction, testcommon.EncodeTx(signedTx))
	if err != nil {
		panic(err)
	}

	err = testcommon.AwaitReceipt(faucetClient, signedTx.Hash())
	if err != nil {
		panic(err)
	}
}
