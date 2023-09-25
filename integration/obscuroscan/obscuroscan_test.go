package faucet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"

	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/config"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/container"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/stretchr/testify/assert"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "obscuroscan",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/obscuroscan/"
)

func TestObscuroscan(t *testing.T) {
	//t.Skip("Commented it out until more testing is driven from this test")
	startPort := integration.StartPortObscuroscanUnitTest
	createObscuroNetwork(t, startPort)

	obsScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultObscuroscanHTTPPortOffset),
		LogPath:         "sys_out",
	}
	serverAddress := fmt.Sprintf("http://%s", obsScanConfig.ServerAddress)

	obsScanContainer, err := container.NewObscuroScanContainer(obsScanConfig)
	require.NoError(t, err)

	err = obsScanContainer.Start()
	require.NoError(t, err)

	// wait for the msg bus contract to be deployed
	time.Sleep(10 * time.Second)

	// make sure the server is ready to receive requests
	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	// Issue tests
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/count/contracts/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":1}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/count/transactions/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":1}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type itemRes struct {
		Item common.BatchHeader `json:"item"`
	}

	itemObj := itemRes{}
	err = json.Unmarshal(body, &itemObj)
	assert.NoError(t, err)
	batchHead := itemObj.Item

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/latest/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batchHeader/%s", serverAddress, batchHead.Hash().String()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/transactions/?offset=0&size=99", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type publicTxsRes struct {
		Result common.TransactionListingResponse `json:"result"`
	}

	publicTxsObj := publicTxsRes{}
	err = json.Unmarshal(body, &publicTxsObj)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(publicTxsObj.Result.TransactionsData))
	assert.Equal(t, uint64(1), publicTxsObj.Result.Total)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batches/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchlisting struct {
		Result common.BatchListingResponse `json:"result"`
	}

	batchlistingObj := batchlisting{}
	err = json.Unmarshal(body, &batchlistingObj)
	assert.NoError(t, err)
	assert.LessOrEqual(t, 9, len(batchlistingObj.Result.BatchesData))
	assert.LessOrEqual(t, uint64(9), batchlistingObj.Result.Total)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/blocks/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type blockListing struct {
		Result common.BlockListingResponse `json:"result"`
	}

	blocklistingObj := blockListing{}
	err = json.Unmarshal(body, &blocklistingObj)
	assert.NoError(t, err)
	// assert.LessOrEqual(t, 9, len(blocklistingObj.Result.BlocksData))
	// assert.LessOrEqual(t, uint64(9), blocklistingObj.Result.Total)

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s", serverAddress, batchlistingObj.Result.BatchesData[0].Hash()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchFetch struct {
		Item *common.ExtBatch `json:"item"`
	}

	batchObj := batchFetch{}
	err = json.Unmarshal(body, &batchObj)
	assert.NoError(t, err)
	assert.Equal(t, batchlistingObj.Result.BatchesData[0].Hash(), batchObj.Item.Hash())

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/info/obscuro/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type configFetch struct {
		Item common.ObscuroNetworkInfo `json:"item"`
	}

	configFetchObj := configFetch{}
	err = json.Unmarshal(body, &configFetchObj)
	assert.NoError(t, err)
	assert.NotEqual(t, configFetchObj.Item.SequencerID, gethcommon.Address{})

	issueTransactions(
		t,
		fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		wallet.NewInMemoryWalletFromConfig("8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b", integration.ObscuroChainID, testlog.Logger()),
		100,
	)

	fmt.Println("Running for 1 hour...")
	time.Sleep(time.Hour)
	// Gracefully shutdown
	err = obsScanContainer.Stop()
	assert.NoError(t, err)
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

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) {
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
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}

func issueTransactions(t *testing.T, hostWSAddr string, issuerWallet wallet.Wallet, numbTxs int) {
	ctx := context.Background()

	vk, err := viewingkey.GenerateViewingKeyForWallet(issuerWallet)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(hostWSAddr, vk, testlog.Logger())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(ctx, nil)
	assert.Nil(t, err)

	if balance.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("not enough balance: has %s has %s obx", issuerWallet.Address().Hex(), balance.String())
	}

	var receipts []gethcommon.Hash
	for i := 0; i < numbTxs; i++ {
		toAddr := datagenerator.RandomAddress()
		nonce, err := authClient.NonceAt(ctx, nil)
		assert.Nil(t, err)

		issuerWallet.SetNonce(nonce)
		estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
			Nonce:    issuerWallet.GetNonceAndIncrement(),
			To:       &toAddr,
			Value:    big.NewInt(100),
			Gas:      uint64(1_000_000),
			GasPrice: gethcommon.Big1,
		})
		assert.Nil(t, err)

		signedTx, err := issuerWallet.SignTransaction(estimatedTx)
		assert.Nil(t, err)

		err = authClient.SendTransaction(ctx, signedTx)
		assert.Nil(t, err)

		fmt.Printf("Issued Tx: %s \n", signedTx.Hash().Hex())
		receipts = append(receipts, signedTx.Hash())
		time.Sleep(1500 * time.Millisecond)
	}

	for _, txHash := range receipts {
		fmt.Printf("Checking for tx receipt for %s \n", txHash)
		var receipt *types.Receipt
		for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
			receipt, err = authClient.TransactionReceipt(ctx, txHash)
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
	}
}
