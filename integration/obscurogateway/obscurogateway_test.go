package faucet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/httputil"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	integrationCommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/tools/walletextension/config"
	"github.com/ten-protocol/go-ten/tools/walletextension/container"
	"github.com/ten-protocol/go-ten/tools/walletextension/lib"
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
	createObscuroNetwork(t, startPort)

	obscuroGatewayConf := config.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: startPort + integration.DefaultObscuroGatewayHTTPPortOffset,
		WalletExtensionPortWS:   startPort + integration.DefaultObscuroGatewayWSPortOffset,
		NodeRPCHTTPAddress:      fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress: fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                 "sys_out",
		VerboseFlag:             false,
		DBType:                  "sqlite",
		TenChainID:              443,
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

	// prefunded wallet
	w := wallet.NewInMemoryWalletFromConfig(genesis.TestnetPrefundedPK, integration.ObscuroChainID, testlog.Logger())

	// run the tests against the exis
	for name, test := range map[string]func(*testing.T, string, string, wallet.Wallet){
		//"testAreTxsMinted":            testAreTxsMinted, this breaks the other tests bc, enable once concurency issues are fixed
		"testErrorHandling":                    testErrorHandling,
		"testMultipleAccountsSubscription":     testMultipleAccountsSubscription,
		"testErrorsRevertedArePassed":          testErrorsRevertedArePassed,
		"testUnsubscribe":                      testUnsubscribe,
		"testClosingConnectionWhileSubscribed": testClosingConnectionWhileSubscribed,
	} {
		t.Run(name, func(t *testing.T) {
			test(t, httpURL, wsURL, w)
		})
	}

	// Gracefully shutdown
	err = obscuroGwContainer.Stop()
	assert.NoError(t, err)
}

func testMultipleAccountsSubscription(t *testing.T, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewUser([]wallet.Wallet{w}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user0.ogClient.UserID())

	user1, err := NewUser([]wallet.Wallet{datagenerator.RandomWallet(integration.ObscuroChainID), datagenerator.RandomWallet(integration.ObscuroChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user0.ogClient.UserID())

	user2, err := NewUser([]wallet.Wallet{datagenerator.RandomWallet(integration.ObscuroChainID), datagenerator.RandomWallet(integration.ObscuroChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user0.ogClient.UserID())

	// register all the accounts for that user
	err = user0.RegisterAccounts()
	require.NoError(t, err)
	err = user1.RegisterAccounts()
	require.NoError(t, err)
	err = user2.RegisterAccounts()
	require.NoError(t, err)

	var amountToTransfer int64 = 1_000_000_000_000_000_000
	// Transfer some funds to user1 and user2 wallets, because they need it to make transactions
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)

	// Print balances of all registered accounts to check if all accounts have funds
	err = user0.PrintUserAccountsBalances()
	require.NoError(t, err)
	err = user1.PrintUserAccountsBalances()
	require.NoError(t, err)
	err = user2.PrintUserAccountsBalances()
	require.NoError(t, err)

	// deploy events contract
	deployTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user0.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user0.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	// check if value was changed in the smart contract with the interactions above
	pack, _ := eventsContractABI.Pack("message2")
	result, err := user1.HTTPClient.CallContract(context.Background(), ethereum.CallMsg{
		From: user1.Wallets[0].Address(),
		To:   &contractReceipt.ContractAddress,
		Data: pack,
	}, nil)
	require.NoError(t, err)

	resultMessage := string(bytes.TrimRight(result[64:], "\x00"))
	require.NoError(t, err)

	// check if the value is the same as hardcoded in smart contract
	hardcodedMessageValue := "foo"
	assert.Equal(t, hardcodedMessageValue, resultMessage)

	// subscribe with all three users for all events in deployed contract
	var user0logs []types.Log
	var user1logs []types.Log
	var user2logs []types.Log
	subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user0.WSClient, &user0logs)
	subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user1.WSClient, &user1logs)
	subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user2.WSClient, &user2logs)

	// user1 calls setMessage and setMessage2 on deployed smart contract with the account
	// that was registered as the first in OG
	user1MessageValue := "user1PublicEvent"
	// interact with smart contract and cause events to be emitted
	_, err = integrationCommon.InteractWithSmartContract(user1.HTTPClient, user1.Wallets[0], eventsContractABI, "setMessage", "user1PrivateEvent", contractReceipt.ContractAddress)
	require.NoError(t, err)
	_, err = integrationCommon.InteractWithSmartContract(user1.HTTPClient, user1.Wallets[0], eventsContractABI, "setMessage2", "user1PublicEvent", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// check if value was changed in the smart contract with the interactions above
	pack, _ = eventsContractABI.Pack("message2")
	result, err = user1.HTTPClient.CallContract(context.Background(), ethereum.CallMsg{
		From: user1.Wallets[0].Address(),
		To:   &contractReceipt.ContractAddress,
		Data: pack,
	}, nil)
	require.NoError(t, err)

	resultMessage = string(bytes.TrimRight(result[64:], "\x00"))
	assert.Equal(t, user1MessageValue, resultMessage)

	// user2 calls setMessage and setMessage2 on deployed smart contract with the account
	// that was registered as the second in OG
	user2MessageValue := "user2PublicEvent"
	// interact with smart contract and cause events to be emitted
	_, err = integrationCommon.InteractWithSmartContract(user2.HTTPClient, user2.Wallets[1], eventsContractABI, "setMessage", "user2PrivateEvent", contractReceipt.ContractAddress)
	require.NoError(t, err)
	_, err = integrationCommon.InteractWithSmartContract(user2.HTTPClient, user2.Wallets[1], eventsContractABI, "setMessage2", "user2PublicEvent", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// check if value was changed in the smart contract with the interactions above
	pack, _ = eventsContractABI.Pack("message2")
	result, err = user1.HTTPClient.CallContract(context.Background(), ethereum.CallMsg{
		From: user1.Wallets[0].Address(),
		To:   &contractReceipt.ContractAddress,
		Data: pack,
	}, nil)
	require.NoError(t, err)
	resultMessage = string(bytes.TrimRight(result[64:], "\x00"))
	assert.Equal(t, user2MessageValue, resultMessage)

	// wait a few seconds to be completely sure all events arrived
	time.Sleep(time.Second * 3)

	// Assert the number of logs received by each client
	// user0 should see two lifecycle events (1 for each interaction with setMessage2)
	assert.Equal(t, 2, len(user0logs))
	// user1 should see three events (two lifecycle events - same as user0) and event with his interaction with setMessage
	assert.Equal(t, 3, len(user1logs))
	// user2 should see three events (two lifecycle events - same as user0) and event with his interaction with setMessage
	assert.Equal(t, 3, len(user2logs))
}

func testAreTxsMinted(t *testing.T, httpURL, wsURL string, w wallet.Wallet) { //nolint: unused
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

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
}

func testErrorHandling(t *testing.T, httpURL, wsURL string, w wallet.Wallet) {
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	// register an account
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// make requests to geth for comparison

	for _, req := range []string{
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":[],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getgetget","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1,"extra":"extra_field"}`,
		`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "0x1234"]],"id":1}`,
	} {
		// ensure the geth request is issued correctly (should return 200 ok with jsonRPCError)
		_, response, err := httputil.PostDataJSON(ogClient.HTTP(), []byte(req))
		require.NoError(t, err)

		// unmarshall the response to JSONRPCMessage
		jsonRPCError := wecommon.JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err)

		// repeat the process for the gateway
		_, response, err = httputil.PostDataJSON(fmt.Sprintf("http://localhost:%d", integration.StartPortObscuroGatewayUnitTest), []byte(req))
		require.NoError(t, err)

		// we only care about format
		jsonRPCError = wecommon.JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err)
	}
}

func testErrorsRevertedArePassed(t *testing.T, httpURL, wsURL string, w wallet.Wallet) {
	// set up the ogClient
	ogClient := lib.NewObscuroGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	// use a standard eth client via the og
	ethStdClient, err := ethclient.Dial(ogClient.HTTP())
	require.NoError(t, err)

	// check the balance
	balance, err := ethStdClient.BalanceAt(context.Background(), w.Address(), nil)
	require.NoError(t, err)
	require.True(t, big.NewInt(0).Cmp(balance) == -1)

	// deploy errors contract
	deployTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(errorsContractBytecode),
	}

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = ethStdClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	receipt, err := integrationCommon.AwaitReceiptEth(context.Background(), ethStdClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	pack, _ := errorsContractABI.Pack("force_require")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted: Forced require")

	// convert error to WE error
	errBytes, err := json.Marshal(err)
	require.NoError(t, err)
	weError := wecommon.JSONError{}
	err = json.Unmarshal(errBytes, &weError)
	require.NoError(t, err)
	require.Equal(t, weError.Message, "execution reverted: Forced require")
	require.Equal(t, weError.Data, "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e466f726365642072657175697265000000000000000000000000000000000000")
	require.Equal(t, weError.Code, 3)

	pack, _ = errorsContractABI.Pack("force_revert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted: Forced revert")

	pack, _ = errorsContractABI.Pack("force_assert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, err.Error(), "execution reverted")
}

func testUnsubscribe(t *testing.T, httpURL, wsURL string, w wallet.Wallet) {
	// create a user with multiple accounts
	user, err := NewUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.ObscuroChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user.ogClient.UserID())

	// register all the accounts for the user
	err = user.RegisterAccounts()
	require.NoError(t, err)

	// deploy events contract
	deployTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	fmt.Println("Deployed contract address: ", contractReceipt.ContractAddress)

	// subscribe to an event
	var userLogs []types.Log
	subscription := subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user.WSClient, &userLogs)

	// make an action that will trigger events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage", "foo", contractReceipt.ContractAddress)
	require.NoError(t, err)

	assert.Equal(t, 1, len(userLogs))

	// Unsubscribe from events
	subscription.Unsubscribe()

	// make another action that will trigger events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage", "bar", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// check that we are not receiving events after unsubscribing
	assert.Equal(t, 1, len(userLogs))
}

func testClosingConnectionWhileSubscribed(t *testing.T, httpURL, wsURL string, w wallet.Wallet) {
	// create a user with multiple accounts
	user, err := NewUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.ObscuroChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Printf("Created user with UserID: %s\n", user.ogClient.UserID())

	// register all the accounts for the user
	err = user.RegisterAccounts()
	require.NoError(t, err)

	// deploy events contract
	deployTx := &types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	fmt.Println("Deployed contract address: ", contractReceipt.ContractAddress)

	// subscribe to an event
	var userLogs []types.Log
	subscription := subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user.WSClient, &userLogs)

	// Close the websocket connection and make sure nothing breaks, but user does not receive events
	user.WSClient.Close()

	// make an action that will emmit events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage2", "foo", contractReceipt.ContractAddress)
	require.NoError(t, err)
	// but with closed connection we don't receive any logs
	assert.Equal(t, 0, len(userLogs))

	// re-establish connection
	wsClient, err := ethclient.Dial(wsURL + "/v1/" + "?u=" + user.ogClient.UserID())
	require.NoError(t, err)
	user.WSClient = wsClient

	// make an action that will emmit events again
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage2", "foo", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// closing connection (above) unsubscribes, and we still should see no logs
	assert.Equal(t, 0, len(userLogs))

	// Call unsubscribe (should handle it without issues even if it is already unsubscribed by closing the channel)
	subscription.Unsubscribe()
}

func transferRandomAddr(t *testing.T, client *ethclient.Client, w wallet.Wallet) common.TxHash { //nolint: unused
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

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client, signedTx.Hash(), time.Minute)
	assert.NoError(t, err)

	fmt.Println("Successfully minted the transaction - ", signedTx.Hash())
	return signedTx.Hash()
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
		WithPrefunding:   true,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
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

func transferETHToAddress(client *ethclient.Client, wallet wallet.Wallet, toAddress gethcommon.Address, amount int64) (*types.Receipt, error) { //nolint:unparam
	transferTx1 := types.LegacyTx{
		Nonce:    wallet.GetNonceAndIncrement(),
		To:       &toAddress,
		Value:    big.NewInt(amount),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     nil,
	}
	signedTx, err := wallet.SignTransaction(&transferTx1)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return integrationCommon.AwaitReceiptEth(context.Background(), client, signedTx.Hash(), 2*time.Second)
}

func subscribeToEvents(addresses []gethcommon.Address, topics [][]gethcommon.Hash, client *ethclient.Client, logs *[]types.Log) ethereum.Subscription { //nolint:unparam
	// Make a subscription
	filterQuery := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: big.NewInt(0), // todo (@ziga) - without those we get errors - fix that and make them configurable
		ToBlock:   big.NewInt(10000),
		Topics:    topics,
	}
	logsCh := make(chan types.Log)

	subscription, err := client.SubscribeFilterLogs(context.Background(), filterQuery, logsCh)
	if err != nil {
		fmt.Printf("Failed to subscribe to filter logs: %v\n", err)
	}

	// Listen for logs in a goroutine
	go func() {
		for {
			select {
			case err := <-subscription.Err():
				fmt.Printf("Error from logs subscription: %v\n", err)
				return
			case log := <-logsCh:
				// append logs to be visible from the main thread
				*logs = append(*logs, log)
			}
		}
	}()

	return subscription
}
