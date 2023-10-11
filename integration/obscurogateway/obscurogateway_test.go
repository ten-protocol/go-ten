package faucet

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
	"github.com/obscuronet/go-obscuro/tools/walletextension/container"
	"github.com/obscuronet/go-obscuro/tools/walletextension/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "obscurogateway",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/obscurogateway/"
)

func TestObscuroGateway(t *testing.T) {
	startPort := integration.StartPortObscuroGatewayUnitTest
	wallets := createObscuroNetwork(t, startPort)

	obscuroGatewayConf := config.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: startPort + integration.DefaultObscuroGatewayHTTPPortOffset,
		WalletExtensionPortWS:   startPort + integration.DefaultObscuroGatewayWSPortOffset,
		NodeRPCHTTPAddress:      fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress: fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                 "sys_out",
		VerboseFlag:             false,
		DBType:                  "sqlite",
	}

	obscuroGwContainer := container.NewWalletExtensionContainerFromConfig(obscuroGatewayConf, testlog.Logger())
	go func() {
		err := obscuroGwContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	// wait for the msg bus contract to be deployed
	time.Sleep(5 * time.Second)

	// make sure the server is ready to receive requests
	httpURL := fmt.Sprintf("http://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortHTTP)
	wsURL := fmt.Sprintf("ws://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortWS)

	// make sure the server is ready to receive requests
	err := waitServerIsReady(httpURL)
	require.NoError(t, err)

	// join + register against the og
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)
	err = ogClient.Join()
	require.NoError(t, err)

	w := wallets.L2FaucetWallet
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// use a standard eth client via the og
	ethStdClient, err := ethclient.Dial(ogClient.HTTP())
	require.NoError(t, err)

	// check the balance
	balance, err := ethStdClient.BalanceAt(context.Background(), w.Address(), nil)
	require.NoError(t, err)
	require.True(t, big.NewInt(0).Cmp(balance) == -1)

	// issue a tx and check it was successfully minted
	txHash := transferRandomAddr(t, ethStdClient, w)
	receipt, err := ethStdClient.TransactionReceipt(context.Background(), txHash)
	assert.NoError(t, err)
	require.True(t, receipt.Status == 1)

	// Gracefully shutdown
	err = obscuroGwContainer.Stop()
	assert.NoError(t, err)
}

func transferRandomAddr(t *testing.T, client *ethclient.Client, w wallet.Wallet) common.TxHash {
	ctx := context.Background()
	toAddr := datagenerator.RandomAddress()
	nonce, err := client.NonceAt(ctx, w.Address(), nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	}
	assert.Nil(t, err)

	fmt.Println("Transferring from:", w.Address(), " to:", toAddr)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = client.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = client.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}
	fmt.Println("Successfully minted the transaction - ", signedTx.Hash())
	return signedTx.Hash()
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) *params.SimWallets {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
		WithPrefunding:   true,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
	return wallets
}

func waitServerIsReady(serverAddr string) error {
	for now := time.Now(); time.Since(now) < 30*time.Second; time.Sleep(500 * time.Millisecond) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/health/", serverAddr))
		if err != nil {
			// give it time to boot up
			if strings.Contains(err.Error(), "connection") {
				continue
			}
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("timed out before server was ready")
}
