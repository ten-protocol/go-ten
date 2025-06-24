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

	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend/config"
	"github.com/ten-protocol/go-ten/tools/tenscan/backend/container"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/valyala/fasthttp"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "tenscan",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/tenscan/"
)

func TestTenscan(t *testing.T) {
	startPort := integration.TestPorts.TestTenscanPort
	createTenNetwork(t, integration.TestPorts.TestTenscanPort)

	tenScanConfig := &config.Config{
		NodeHostAddress: fmt.Sprintf("http://127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		ServerAddress:   fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultTenscanHTTPPortOffset),
		LogPath:         "sys_out",
	}
	serverAddress := fmt.Sprintf("http://%s", tenScanConfig.ServerAddress)

	tenScanContainer, err := container.NewTenScanContainer(tenScanConfig)
	require.NoError(t, err)

	err = tenScanContainer.Start()
	require.NoError(t, err)

	// wait for the msg bus contract to be deployed
	time.Sleep(30 * time.Second)

	// make sure the server is ready to receive requests
	err = waitServerIsReady(serverAddress)
	require.NoError(t, err)

	issueTransactions(
		t,
		fmt.Sprintf("ws://127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		wallet.NewInMemoryWalletFromConfig(testcommon.TestnetPrefundedPK, integration.TenChainID, testlog.Logger()),
		5,
	)

	err = waitForFirstRollup(serverAddress)
	require.NoError(t, err)

	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/count/contracts/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":22}", string(body))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/count/transactions/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "{\"count\":5}", string(body))

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

	statusCode, _, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s", serverAddress, batchHead.Hash().String()))
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
	assert.Equal(t, 5, len(publicTxsObj.Result.TransactionsData))
	assert.Equal(t, uint64(5), publicTxsObj.Result.Total)

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
	// check results are descending order (latest first)
	assert.LessOrEqual(t, batchlistingObj.Result.BatchesData[1].Height.Cmp(batchlistingObj.Result.BatchesData[0].Height), 0)
	// check "hash" field is included in json response
	assert.Contains(t, string(body), "\"hash\"")

	// fetch block listing
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/blocks/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type blockListing struct {
		Result common.BlockListingResponse `json:"result"`
	}

	blocklistingObj := blockListing{}
	err = json.Unmarshal(body, &blocklistingObj)
	assert.NoError(t, err)

	// fetch batch by hash
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s", serverAddress, batchlistingObj.Result.BatchesData[0].Header.Hash()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchFetch struct {
		Item *common.ExtBatch `json:"item"`
	}

	batchObj := batchFetch{}
	err = json.Unmarshal(body, &batchObj)
	assert.NoError(t, err)
	assert.Equal(t, batchlistingObj.Result.BatchesData[0].Header.Hash(), batchObj.Item.Header.Hash())

	// fetch batch by sequence num
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/seq/%s", serverAddress, batchlistingObj.Result.BatchesData[0].SequencerOrderNo))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type pubBatchFetch struct {
		Item *common.PublicBatch `json:"item"`
	}

	pubBatchObj := pubBatchFetch{}
	err = json.Unmarshal(body, &pubBatchObj)
	assert.NoError(t, err)

	assert.Equal(t, batchlistingObj.Result.BatchesData[0].Header.SequencerOrderNo, pubBatchObj.Item.SequencerOrderNo)

	// fetch rollup listing
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollups/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type rollupListing struct {
		Result common.RollupListingResponse `json:"result"`
	}

	rollupListingObj := rollupListing{}
	err = json.Unmarshal(body, &rollupListingObj)
	assert.NoError(t, err)
	assert.LessOrEqual(t, 1, len(rollupListingObj.Result.RollupsData))
	assert.LessOrEqual(t, uint64(1), rollupListingObj.Result.Total)
	assert.Contains(t, string(body), "\"hash\"")

	// test pagination for rollup batches
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/%s/batches?offset=0&size=1", serverAddress, rollupListingObj.Result.RollupsData[0].Header.Hash()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchListing struct {
		Result common.BatchListingResponse `json:"result"`
	}

	firstPageObj := batchListing{}
	err = json.Unmarshal(body, &firstPageObj)
	assert.NoError(t, err)
	firstPageBatches := firstPageObj.Result.BatchesData

	// second page
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/%s/batches?offset=1&size=2", serverAddress, rollupListingObj.Result.RollupsData[0].Header.Hash()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	secondPageObj := batchListing{}
	err = json.Unmarshal(body, &secondPageObj)
	assert.NoError(t, err)
	secondPageBatches := secondPageObj.Result.BatchesData

	// verify different pages have different batches
	assert.NotEqual(t, firstPageBatches[0].SequencerOrderNo, secondPageBatches[0].SequencerOrderNo)

	// test pagination for batch transactions
	batchHash := firstPageBatches[0].Header.Hash()

	// first page of transactions
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s/transactions?offset=0&size=2", serverAddress, batchHash))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type batchTxListing struct {
		Result common.TransactionListingResponse `json:"result"`
	}

	firstPageTxObj := batchTxListing{}
	err = json.Unmarshal(body, &firstPageTxObj)
	assert.NoError(t, err)
	firstPageTxs := firstPageTxObj.Result.TransactionsData

	// second page of transactions
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/%s/transactions?offset=2&size=2", serverAddress, batchHash))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	secondPageTxObj := batchTxListing{}
	err = json.Unmarshal(body, &secondPageTxObj)
	assert.NoError(t, err)
	secondPageTxs := secondPageTxObj.Result.TransactionsData

	// verify different pages have different transactions
	if len(firstPageTxs) > 0 && len(secondPageTxs) > 0 {
		assert.NotEqual(t, firstPageTxs[0].TransactionHash, secondPageTxs[0].TransactionHash)
	}

	// verify transactions are ordered by timestamp descending
	if len(firstPageTxs) > 1 {
		assert.GreaterOrEqual(t, firstPageTxs[0].BatchTimestamp, firstPageTxs[1].BatchTimestamp)
	}

	// fetch transaction listing
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/transactions/?offset=0&size=10", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type txListing struct {
		Result common.TransactionListingResponse `json:"result"`
	}

	txListingObj := txListing{}
	err = json.Unmarshal(body, &txListingObj)
	assert.NoError(t, err)
	assert.LessOrEqual(t, 5, len(txListingObj.Result.TransactionsData))
	assert.LessOrEqual(t, uint64(5), txListingObj.Result.Total)

	// fetch batch by height from tx
	batchHeight := txListingObj.Result.TransactionsData[0].BatchHeight
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/batch/height/%s", serverAddress, batchHeight))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type publicBatchFetch struct {
		Item *common.PublicBatch `json:"item"`
	}

	publicBatchObj := publicBatchFetch{}
	err = json.Unmarshal(body, &publicBatchObj)
	assert.NoError(t, err)
	assert.True(t, publicBatchObj.Item.Height.Cmp(batchHeight) == 0)

	// fetch tx by hash
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/transaction/%s", serverAddress, txListingObj.Result.TransactionsData[0].TransactionHash))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type txFetch struct {
		Item *common.PublicTransaction `json:"item"`
	}

	txObj := txFetch{}
	err = json.Unmarshal(body, &txObj)
	assert.NoError(t, err)
	assert.True(t, txObj.Item.Finality == common.BatchFinal)

	// search API tests
	rollupHash := rollupListingObj.Result.RollupsData[0].Header.Hash()
	txHash := txListingObj.Result.TransactionsData[0].TransactionHash
	batchSequence := batchlistingObj.Result.BatchesData[0].SequencerOrderNo.String()

	// search by rollup hash
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/search/?query=%s", serverAddress, rollupHash.Hex()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type searchResponse struct {
		Result common.SearchResponse `json:"result"`
	}

	searchObj := searchResponse{}
	err = json.Unmarshal(body, &searchObj)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), searchObj.Result.Total)
	assert.Equal(t, 1, len(searchObj.Result.ResultsData))
	assert.Equal(t, "rollup", searchObj.Result.ResultsData[0].Type)
	assert.Equal(t, strings.TrimPrefix(rollupHash.Hex(), "0x"), searchObj.Result.ResultsData[0].Hash)

	// search by batch hash
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/search/?query=%s", serverAddress, batchHash.Hex()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	err = json.Unmarshal(body, &searchObj)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), searchObj.Result.Total)
	assert.Equal(t, 1, len(searchObj.Result.ResultsData))
	assert.Equal(t, "batch", searchObj.Result.ResultsData[0].Type)
	assert.Equal(t, strings.TrimPrefix(batchHash.Hex(), "0x"), searchObj.Result.ResultsData[0].Hash)

	// search by tx hash
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/search/?query=%s", serverAddress, txHash.Hex()))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	err = json.Unmarshal(body, &searchObj)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), searchObj.Result.Total)
	assert.Equal(t, 1, len(searchObj.Result.ResultsData))
	assert.Equal(t, "transaction", searchObj.Result.ResultsData[0].Type)
	assert.Equal(t, strings.TrimPrefix(txHash.Hex(), "0x"), searchObj.Result.ResultsData[0].Hash)

	// search by sequence number
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/search/?query=%s", serverAddress, batchSequence))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	err = json.Unmarshal(body, &searchObj)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, uint64(1), searchObj.Result.Total)
	assert.GreaterOrEqual(t, 1, len(searchObj.Result.ResultsData))

	// since there is a diff of 1 between sequence number and height we will find two results we just need to check the seq num is present
	foundExpectedSequence := false
	for _, result := range searchObj.Result.ResultsData {
		if result.Type == "batch" && result.Sequence.String() == batchSequence {
			foundExpectedSequence = true
			break
		}
	}
	assert.True(t, foundExpectedSequence, "Expected to find batch with sequence %s", batchSequence)

	// non-existent hash should return nothing
	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/items/search/?query=0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	err = json.Unmarshal(body, &searchObj)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), searchObj.Result.Total)
	assert.Equal(t, 0, len(searchObj.Result.ResultsData))

	statusCode, body, err = fasthttp.Get(nil, fmt.Sprintf("%s/info/obscuro/", serverAddress))
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)

	type configFetch struct {
		Item common.TenNetworkInfo `json:"item"`
	}

	configFetchObj := configFetch{}
	err = json.Unmarshal(body, &configFetchObj)
	assert.NoError(t, err)

	err = tenScanContainer.Stop()
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

// Creates a single-node TEN network for testing.
func createTenNetwork(t *testing.T, startPort int) {
	// Create the TEN network.
	wallets := params.NewSimWallets(1, 1, integration.EthereumChainID, integration.TenChainID)
	simParams := params.SimParams{
		NumberOfNodes:       1,
		AvgBlockDuration:    2 * time.Second,
		ContractRegistryLib: ethereummock.NewContractRegistryLibMock(),
		ERC20ContractLib:    ethereummock.NewERC20ContractLibMock(),
		Wallets:             wallets,
		StartPort:           startPort,
		WithPrefunding:      true,
		L1BeaconPort:        integration.TestPorts.TestTenscanPort + integration.DefaultPrysmGatewayPortOffset,
	}

	tenNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(tenNetwork.TearDown)
	_, err := tenNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test TEN network. Cause: %s", err))
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
		t.Errorf("not enough balance: has %s has %s ten", issuerWallet.Address().Hex(), balance.String())
	}

	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)
	issuerWallet.SetNonce(nonce)

	var receipts []gethcommon.Hash
	for i := 0; i < numbTxs; i++ {
		toAddr := datagenerator.RandomAddress()

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
		for start := time.Now(); time.Since(start) < 2*time.Minute; time.Sleep(time.Second) {
			receipt, err = authClient.TransactionReceipt(ctx, txHash)
			if err == nil {
				break
			}
			if receipt != nil && receipt.Status == 1 {
				break
			}
			fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
		}

		if receipt == nil {
			t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 60*time.Second, receipt)
		}
		if receipt.Status == 0 {
			t.Fatalf("Tx Failed")
		}
	}
}

func waitForFirstRollup(serverAddress string) error {
	for now := time.Now(); time.Since(now) < 4*time.Minute; time.Sleep(5 * time.Second) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/items/rollup/latest/", serverAddress))
		if err != nil {
			if strings.Contains(err.Error(), "connection") {
				continue
			}
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("timed out before rollup was found")
}
