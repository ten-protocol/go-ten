package tengateway

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ten-protocol/go-ten/tools/walletextension"

	"github.com/go-kit/kit/transport/http/jsonrpc"
	log2 "github.com/ten-protocol/go-ten/go/common/log"

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
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	integrationCommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/tools/walletextension/lib"
	"github.com/valyala/fasthttp"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "tengateway",
		TestSubtype: "test",
		LogLevel:    log.LvlTrace,
	})
}

const (
	testLogs = "../.build/tengateway/"
)

func TestTenGateway(t *testing.T) {
	startPort := integration.TestPorts.TestTenGatewayPort
	createTenNetwork(t, startPort)

	tenGatewayConf := wecommon.Config{
		WalletExtensionHost:            "127.0.0.1",
		WalletExtensionPortHTTP:        startPort + integration.DefaultTenGatewayHTTPPortOffset,
		WalletExtensionPortWS:          startPort + integration.DefaultTenGatewayWSPortOffset,
		NodeRPCHTTPAddress:             fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress:        fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                        "sys_out",
		LogLevel:                       4, // info level
		DBType:                         "sqlite",
		TenChainID:                     5443,
		StoreIncomingTxs:               true,
		RateLimitUserComputeTime:       0,
		RateLimitWindow:                1 * time.Second,
		RateLimitMaxConcurrentRequests: 3,
		SessionKeyExpirationThreshold:  10 * time.Second,
		SessionKeyExpirationInterval:   2 * time.Second,
	}

	tenGwContainer := walletextension.NewContainerFromConfig(tenGatewayConf, testlog.Logger())
	go func() {
		err := tenGwContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	// wait for the msg bus contract to be deployed
	time.Sleep(5 * time.Second)

	// make sure the server is ready to receive requests
	httpURL := fmt.Sprintf("http://%s:%d", tenGatewayConf.WalletExtensionHost, tenGatewayConf.WalletExtensionPortHTTP)
	wsURL := fmt.Sprintf("ws://%s:%d", tenGatewayConf.WalletExtensionHost, tenGatewayConf.WalletExtensionPortWS)

	// make sure the server is ready to receive requests
	err := waitServerIsReady(httpURL)
	require.NoError(t, err)

	// prefunded wallet
	w := wallet.NewInMemoryWalletFromConfig(integrationCommon.TestnetPrefundedPK, integration.TenChainID, testlog.Logger())

	// run the tests against the exis
	for name, test := range map[string]func(*testing.T, int, string, string, wallet.Wallet){
		//"testAreTxsMinted":            testAreTxsMinted, this breaks the other tests bc, enable once concurrency issues are fixed
		"testErrorHandling":                    testErrorHandling,
		"testMultipleAccountsSubscription":     testMultipleAccountsSubscription,
		"testNewHeadsSubscription":             testNewHeadsSubscription,
		"testErrorsRevertedArePassed":          testErrorsRevertedArePassed,
		"testUnsubscribe":                      testUnsubscribe,
		"testClosingConnectionWhileSubscribed": testClosingConnectionWhileSubscribed,
		"testSubscriptionTopics":               testSubscriptionTopics,
		"testDifferentMessagesOnRegister":      testDifferentMessagesOnRegister,
		"testInvokeNonSensitiveMethod":         testInvokeNonSensitiveMethod,
		"testQueryAndRpcTokenModes":            testQueryAndRpcTokenModes,

		"testSessionKeysGetStorageAt":             testSessionKeysGetStorageAt,
		"testSessionKeysSendTransaction":          testSessionKeysSendTransaction,
		"testSessionKeyExpirationAndFundRecovery": testSessionKeyExpirationAndFundRecovery,
		"testSessionKeyFundRecoveryOnDeletion":    testSessionKeyFundRecoveryOnDeletion,
		// "testRateLimiter":                   testRateLimiter,
	} {
		t.Run(name, func(t *testing.T) {
			test(t, startPort, httpURL, wsURL, w)
		})
	}

	// Gracefully shutdown
	// todo remove this sleep when tests stabilize
	time.Sleep(20 * time.Second)
	err = tenGwContainer.Stop()
	assert.NoError(t, err)
}

//func testRateLimiter(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
//	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
//	require.NoError(t, err)
//	testlog.Logger().Info("Created user with encryption token", "t", user0.tgClient.UserID())
//	// register the user so we can call the endpoints that require authentication
//	err = user0.RegisterAccounts()
//	require.NoError(t, err)
//
//	// call BalanceAt - fist call should be successful
//	_, err = user0.HTTPClient.BalanceAt(context.Background(), user0.Wallets[0].Address(), nil)
//	require.NoError(t, err)
//
//	// sleep for a period of time to allow the rate limiter to reset
//	time.Sleep(1 * time.Second)
//
//	// first call after the rate limiter reset should be successful
//	_, err = user0.HTTPClient.BalanceAt(context.Background(), user0.Wallets[0].Address(), nil)
//	require.NoError(t, err)
//
//	address := user0.Wallets[0].Address()
//
//	// make 1000 requests with the same user to "spam" the gateway
//	for i := 0; i < 1000; i++ {
//		msg := ethereum.CallMsg{
//			From: address,
//			To:   &address, // Example: self-call to the user's address
//			Gas:  uint64(i),
//			Data: nil,
//		}
//
//		user0.HTTPClient.EstimateGas(context.Background(), msg)
//	}
//
//	// after 1000 requests, the rate limiter should block the user
//	_, err = user0.HTTPClient.BalanceAt(context.Background(), user0.Wallets[0].Address(), nil)
//	require.Error(t, err)
//	require.Equal(t, "rate limit exceeded", err.Error())
//}

func testSessionKeysGetStorageAt(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user0.tgClient.UserID())

	// Register the user so we can call the endpoints that require authentication
	err = user0.RegisterAccounts()
	require.NoError(t, err)

	// Simple print to verify the test is running
	fmt.Println("testSessionKeysGetStorageAt: Test is running successfully!")
	testlog.Logger().Info("testSessionKeysGetStorageAt: Test is running successfully!")

	// Get the user's balance as a simple operation
	balance, err := user0.HTTPClient.BalanceAt(context.Background(), user0.Wallets[0].Address(), nil)
	require.NoError(t, err)

	// Print the balance to show the test is working
	fmt.Printf("testSessionKeysGetStorageAt: User balance: %s\n", balance.String())
	testlog.Logger().Info("testSessionKeysGetStorageAt: User balance", "balance", balance.String())

	ctx := context.Background()

	// 1) Create session key via eth_getStorageAt (CQ method 0x...0003)
	createSessionKeyAddr := gethcommon.HexToAddress("0x0000000000000000000000000000000000000003")
	skAddrBytes, err := user0.HTTPClient.StorageAt(ctx, createSessionKeyAddr, gethcommon.Hash{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, skAddrBytes)
	skAddress := gethcommon.BytesToAddress(skAddrBytes)
	t.Logf("✓ Session key created: %s", skAddress.Hex())
	fmt.Printf("Session key created: %s\n", skAddress.Hex())

	// 2) Fund the session key from the original wallet
	fundAmount := big.NewInt(0).Mul(big.NewInt(1e15), big.NewInt(1)) // 0.001 TEN
	fromAddr := user0.Wallets[0].Address()
	gasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	gasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: fromAddr, To: &skAddress, Value: fundAmount})
	require.NoError(t, err)
	nonce, err := user0.HTTPClient.PendingNonceAt(ctx, fromAddr)
	require.NoError(t, err)
	legacy := &types.LegacyTx{Nonce: nonce, To: &skAddress, Value: fundAmount, GasPrice: gasPrice, Gas: gasLimit}
	signedFundingTx, err := w.SignTransaction(legacy)
	require.NoError(t, err)
	err = user0.HTTPClient.SendTransaction(ctx, signedFundingTx)
	require.NoError(t, err)

	// wait for receipt
	{
		var rec *types.Receipt
		for i := 0; i < 60; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, signedFundingTx.Hash())
			if err == nil && rec != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
	}
	t.Logf("✓ Session key funded with %s TEN", fundAmount.String())

	// 3) Build an unsigned tx from session key back to original account, send via getStorageAt (CQ 0x...0005)
	returnAmount := big.NewInt(0).Div(fundAmount, big.NewInt(2))
	skGasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	skGasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: skAddress, To: &fromAddr, Value: returnAmount})
	require.NoError(t, err)
	skNonce, err := user0.HTTPClient.PendingNonceAt(ctx, skAddress)
	require.NoError(t, err)
	unsigned := types.NewTx(&types.LegacyTx{Nonce: skNonce, To: &fromAddr, Value: returnAmount, GasPrice: skGasPrice, Gas: skGasLimit})
	blob, err := unsigned.MarshalBinary()
	require.NoError(t, err)
	txB64 := base64.StdEncoding.EncodeToString(blob)

	paramsObj := map[string]string{
		"sessionKeyAddress": skAddress.Hex(),
		"tx":                txB64,
	}
	paramsJSON, err := json.Marshal(paramsObj)
	require.NoError(t, err)

	var txHashBytes hexutil.Bytes
	err = user0.HTTPClient.Client().CallContext(ctx, &txHashBytes, "eth_getStorageAt",
		"0x0000000000000000000000000000000000000005", string(paramsJSON), "latest")
	require.NoError(t, err)
	txHash := gethcommon.BytesToHash(txHashBytes)

	// wait for receipt
	{
		var rec *types.Receipt
		for i := 0; i < 60; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, txHash)
			if err == nil && rec != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
	}
	t.Logf("✓ Return transaction sent: %s TEN", returnAmount.String())

	// 4) Delete the session key via getStorageAt (CQ 0x...0004)
	delParamsObj := map[string]string{
		"sessionKeyAddress": skAddress.Hex(),
	}
	delParamsJSON, err := json.Marshal(delParamsObj)
	require.NoError(t, err)

	var delResult hexutil.Bytes
	err = user0.HTTPClient.Client().CallContext(ctx, &delResult, "eth_getStorageAt",
		"0x0000000000000000000000000000000000000004", string(delParamsJSON), "latest")
	require.NoError(t, err)
	require.Len(t, delResult, 1)
	require.Equal(t, byte(0x01), delResult[0])
	t.Logf("✓ Session key deleted: %s", skAddress.Hex())
}

func testSessionKeysSendTransaction(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user0.tgClient.UserID())

	// Register the user so we can call the endpoints that require authentication
	err = user0.RegisterAccounts()
	require.NoError(t, err)

	// Simple print to verify the test is running
	fmt.Println("testSessionKeysSendTransaction: Test is running successfully!")
	testlog.Logger().Info("testSessionKeysSendTransaction: Test is running successfully!")

	// Get the user's balance as a simple operation
	balance, err := user0.HTTPClient.BalanceAt(context.Background(), user0.Wallets[0].Address(), nil)
	require.NoError(t, err)

	// Print the balance to show the test is working
	fmt.Printf("testSessionKeysSendTransaction: User balance: %s\n", balance.String())
	testlog.Logger().Info("testSessionKeysSendTransaction: User balance", "balance", balance.String())

	ctx := context.Background()

	// 1) Create session key via eth_getStorageAt (CQ method 0x...0003)
	createSessionKeyAddr := gethcommon.HexToAddress("0x0000000000000000000000000000000000000003")
	skAddrBytes, err := user0.HTTPClient.StorageAt(ctx, createSessionKeyAddr, gethcommon.Hash{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, skAddrBytes)
	skAddress := gethcommon.BytesToAddress(skAddrBytes)
	t.Logf("✓ Session key created: %s", skAddress.Hex())
	fmt.Printf("Session key created: %s\n", skAddress.Hex())

	// 2) Fund the session key from the original wallet
	fundAmount := big.NewInt(0).Mul(big.NewInt(1e15), big.NewInt(1)) // 0.001 TEN
	fromAddr := user0.Wallets[0].Address()
	gasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	gasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: fromAddr, To: &skAddress, Value: fundAmount})
	require.NoError(t, err)
	nonce, err := user0.HTTPClient.PendingNonceAt(ctx, fromAddr)
	require.NoError(t, err)
	legacy := &types.LegacyTx{Nonce: nonce, To: &skAddress, Value: fundAmount, GasPrice: gasPrice, Gas: gasLimit}
	signedFundingTx, err := w.SignTransaction(legacy)
	require.NoError(t, err)
	err = user0.HTTPClient.SendTransaction(ctx, signedFundingTx)
	require.NoError(t, err)

	// wait for receipt
	{
		var rec *types.Receipt
		for i := 0; i < 60; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, signedFundingTx.Hash())
			if err == nil && rec != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
	}
	t.Logf("✓ Session key funded with %s TEN", fundAmount.String())

	// 3) Send transaction using eth_sendTransaction with session key in From field
	returnAmount := big.NewInt(0).Div(fundAmount, big.NewInt(2))
	skGasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	skGasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: skAddress, To: &fromAddr, Value: returnAmount})
	require.NoError(t, err)
	skNonce, err := user0.HTTPClient.PendingNonceAt(ctx, skAddress)
	require.NoError(t, err)

	// Send transaction using eth_sendTransaction (this will use our new SendTransaction method)
	var txHash gethcommon.Hash
	err = user0.HTTPClient.Client().CallContext(ctx, &txHash, "eth_sendTransaction", map[string]interface{}{
		"from":     skAddress.Hex(),
		"to":       fromAddr.Hex(),
		"value":    fmt.Sprintf("0x%x", returnAmount),
		"gas":      fmt.Sprintf("0x%x", skGasLimit),
		"gasPrice": fmt.Sprintf("0x%x", skGasPrice),
		"nonce":    fmt.Sprintf("0x%x", skNonce),
	})
	require.NoError(t, err)
	require.NotEqual(t, gethcommon.Hash{}, txHash)
	t.Logf("✓ Transaction sent via eth_sendTransaction: %s", txHash.Hex())

	// wait for receipt
	{
		var rec *types.Receipt
		for i := 0; i < 60; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, txHash)
			if err == nil && rec != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
	}
	t.Logf("✓ Return transaction confirmed: %s TEN", returnAmount.String())

	// 4) Test that eth_sendTransaction fails with non-session key address
	// Use the user's own address (which is not a session key) to test the failure case
	nonSessionKeyAddr := user0.Wallets[0].Address()
	t.Logf("Testing eth_sendTransaction with non-session key address: %s", nonSessionKeyAddr.Hex())

	// Get the current nonce for this address to avoid "nonce too low" errors
	nonceForNonSessionKey, err := user0.HTTPClient.PendingNonceAt(ctx, nonSessionKeyAddr)
	require.NoError(t, err)
	t.Logf("Using nonce %d for non-session key address", nonceForNonSessionKey)

	var failTxHash gethcommon.Hash
	err = user0.HTTPClient.Client().CallContext(ctx, &failTxHash, "eth_sendTransaction", map[string]interface{}{
		"from":     nonSessionKeyAddr.Hex(),
		"to":       fromAddr.Hex(),
		"value":    fmt.Sprintf("0x%x", big.NewInt(1000)),
		"gas":      fmt.Sprintf("0x%x", uint64(21000)),
		"gasPrice": fmt.Sprintf("0x%x", skGasPrice),
		"nonce":    fmt.Sprintf("0x%x", nonceForNonSessionKey),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "session key address")
	require.Contains(t, err.Error(), "not found")
	t.Logf("✓ eth_sendTransaction correctly rejected with non-session key address")

	// 5) Delete the session key via getStorageAt (CQ 0x...0004)
	delParamsObj := map[string]string{
		"sessionKeyAddress": skAddress.Hex(),
	}
	delParamsJSON, err := json.Marshal(delParamsObj)
	require.NoError(t, err)

	var delResult hexutil.Bytes
	err = user0.HTTPClient.Client().CallContext(ctx, &delResult, "eth_getStorageAt",
		"0x0000000000000000000000000000000000000004", string(delParamsJSON), "latest")
	require.NoError(t, err)
	require.Len(t, delResult, 1)
	require.Equal(t, byte(0x01), delResult[0])
	t.Logf("✓ Session key deleted: %s", skAddress.Hex())
}

// testSessionKeyFundRecoveryOnDeletion verifies that deleting a session key triggers
// a refund of (balance - gas) back to the user's primary account.
func testSessionKeyFundRecoveryOnDeletion(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	fmt.Println("=== Starting testSessionKeyFundRecoveryOnDeletion ===")

	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	fmt.Println("✓ Created user with encryption token:", user0.tgClient.UserID())

	// Register the user so we can call the endpoints that require authentication
	err = user0.RegisterAccounts()
	require.NoError(t, err)
	fmt.Println("✓ User accounts registered")

	ctx := context.Background()

	// 1) Create session key via eth_getStorageAt (CQ method 0x...0003)
	fmt.Println("Step 1: Creating session key...")
	createSessionKeyAddr := gethcommon.HexToAddress("0x0000000000000000000000000000000000000003")
	skAddrBytes, err := user0.HTTPClient.StorageAt(ctx, createSessionKeyAddr, gethcommon.Hash{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, skAddrBytes)
	skAddress := gethcommon.BytesToAddress(skAddrBytes)
	fmt.Printf("✓ Session key created: %s\n", skAddress.Hex())

	// 2) Fund the session key from the original wallet (user's first account)
	fmt.Println("Step 2: Funding session key...")
	fundAmount := big.NewInt(0).Mul(big.NewInt(1e15), big.NewInt(1)) // 0.001 ETH (reduced amount)
	fromAddr := user0.Wallets[0].Address()
	fmt.Printf("Funding from: %s\n", fromAddr.Hex())
	fmt.Printf("Funding amount: %s\n", fundAmount.String())

	gasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	fmt.Printf("Gas price: %s\n", gasPrice.String())

	gasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: fromAddr, To: &skAddress, Value: fundAmount})
	require.NoError(t, err)
	fmt.Printf("Gas limit: %d\n", gasLimit)

	nonce, err := user0.HTTPClient.PendingNonceAt(ctx, fromAddr)
	require.NoError(t, err)
	fmt.Printf("Nonce: %d\n", nonce)

	legacy := &types.LegacyTx{Nonce: nonce, To: &skAddress, Value: fundAmount, GasPrice: gasPrice, Gas: gasLimit}
	signedFundingTx, err := w.SignTransaction(legacy)
	require.NoError(t, err)
	err = user0.HTTPClient.SendTransaction(ctx, signedFundingTx)
	require.NoError(t, err)
	fmt.Printf("✓ Funding transaction sent: %s\n", signedFundingTx.Hash().Hex())

	// wait for receipt of funding tx
	fmt.Println("Waiting for funding transaction receipt...")
	{
		var rec *types.Receipt
		for i := 0; i < 10; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, signedFundingTx.Hash())
			if err == nil && rec != nil {
				break
			}
			fmt.Printf("Waiting for receipt... attempt %d\n", i+1)
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
		fmt.Println("✓ Funding transaction confirmed")
	}

	// 3) Record pre-deletion balances
	fmt.Println("Step 3: Recording pre-deletion balances...")
	skInitialBalance, err := user0.HTTPClient.BalanceAt(ctx, skAddress, nil)
	require.NoError(t, err)
	fmt.Printf("Session key initial balance: %s\n", skInitialBalance.String())
	require.True(t, skInitialBalance.Cmp(big.NewInt(0)) > 0)

	userInitialBalance, err := user0.HTTPClient.BalanceAt(ctx, fromAddr, nil)
	require.NoError(t, err)
	fmt.Printf("User initial balance: %s\n", userInitialBalance.String())

	// 4) Delete the session key via getStorageAt (CQ 0x...0004)
	fmt.Println("Step 4: Deleting session key...")
	delParamsObj := map[string]string{
		"sessionKeyAddress": skAddress.Hex(),
	}
	delParamsJSON, err := json.Marshal(delParamsObj)
	require.NoError(t, err)

	var delResult hexutil.Bytes
	err = user0.HTTPClient.Client().CallContext(ctx, &delResult, "eth_getStorageAt",
		"0x0000000000000000000000000000000000000004", string(delParamsJSON), "latest")
	require.NoError(t, err)
	require.Len(t, delResult, 1)
	require.Equal(t, byte(0x01), delResult[0])
	fmt.Printf("✓ Session key deleted: %s\n", skAddress.Hex())

	// 5) Wait briefly to allow refund tx to be mined (it is sent internally)
	fmt.Println("Step 5: Waiting for fund recovery...")
	var skFinalBalance *big.Int
	for i := 0; i < 10; i++ {
		skFinalBalance, err = user0.HTTPClient.BalanceAt(ctx, skAddress, nil)
		require.NoError(t, err)
		fmt.Printf("Session key balance check %d: %s\n", i+1, skFinalBalance.String())
		if skFinalBalance.Cmp(big.NewInt(0)) == 0 { // fully drained
			fmt.Println("✓ Session key balance is zero - funds recovered")
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	require.NotNil(t, skFinalBalance)
	fmt.Printf("Final session key balance: %s\n", skFinalBalance.String())

	// It may not be exactly zero if dust threshold prevented sending, but should be <= initial.
	require.True(t, skFinalBalance.Cmp(skInitialBalance) <= 0)

	// 6) Verify user's primary account increased by at least some positive amount
	fmt.Println("Step 6: Verifying user balance increase...")
	userFinalBalance, err := user0.HTTPClient.BalanceAt(ctx, fromAddr, nil)
	require.NoError(t, err)
	fmt.Printf("User final balance: %s\n", userFinalBalance.String())
	fmt.Printf("Balance increase: %s\n", big.NewInt(0).Sub(userFinalBalance, userInitialBalance).String())

	require.Truef(t, userFinalBalance.Cmp(userInitialBalance) > 0,
		"expected user balance to increase: before=%s after=%s",
		userInitialBalance.String(), userFinalBalance.String())

	fmt.Println("=== testSessionKeyFundRecoveryOnDeletion completed successfully ===")
}

func testSessionKeyExpirationAndFundRecovery(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user0.tgClient.UserID())

	// Register the user so we can call the endpoints that require authentication
	err = user0.RegisterAccounts()
	require.NoError(t, err)

	// Sanity log to mark test start
	testlog.Logger().Info("testSessionKeyExpirationAndFundRecovery: started")

	ctx := context.Background()

	// 1) Create session key via eth_getStorageAt (CQ method 0x...0003)
	createSessionKeyAddr := gethcommon.HexToAddress("0x0000000000000000000000000000000000000003")
	skAddrBytes, err := user0.HTTPClient.StorageAt(ctx, createSessionKeyAddr, gethcommon.Hash{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, skAddrBytes)
	skAddress := gethcommon.BytesToAddress(skAddrBytes)

	// 2) Fund the session key from the original wallet
	fundAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1)) // 1 TEN
	fromAddr := user0.Wallets[0].Address()
	gasPrice, err := user0.HTTPClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	gasLimit, err := user0.HTTPClient.EstimateGas(ctx, ethereum.CallMsg{From: fromAddr, To: &skAddress, Value: fundAmount})
	require.NoError(t, err)
	nonce, err := user0.HTTPClient.PendingNonceAt(ctx, fromAddr)
	require.NoError(t, err)
	legacy := &types.LegacyTx{Nonce: nonce, To: &skAddress, Value: fundAmount, GasPrice: gasPrice, Gas: gasLimit}
	signedFundingTx, err := w.SignTransaction(legacy)
	require.NoError(t, err)
	err = user0.HTTPClient.SendTransaction(ctx, signedFundingTx)
	require.NoError(t, err)

	// wait for receipt
	{
		var rec *types.Receipt
		for i := 0; i < 30; i++ {
			rec, err = user0.HTTPClient.TransactionReceipt(ctx, signedFundingTx.Hash())
			if err == nil && rec != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.NotNil(t, rec)
		require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)
	}
	// Session key funded from the user's primary account
	t.Logf("✓ Session key funded with %s TEN", fundAmount.String())

	// 3) Record initial balances for assertions
	initialBalance, err := user0.HTTPClient.BalanceAt(ctx, skAddress, nil)
	require.NoError(t, err)
	require.Equal(t, fundAmount, initialBalance)
	t.Logf("✓ Initial session key balance: %s TEN", initialBalance.String())

	// 4) Check initial balance of user's first account
	initialUserBalance, err := user0.HTTPClient.BalanceAt(ctx, fromAddr, nil)
	require.NoError(t, err)
	t.Logf("✓ Initial user balance: %s TEN", initialUserBalance.String())

	// 5) Wait for session key expiration (default is 10 seconds, wait 12 seconds to be safe)
	t.Logf("⏳ Waiting for session key expiration (12 seconds)...")
	time.Sleep(12 * time.Second)
	t.Logf("✓ Session key should now be expired")

	// 6) After expiration, the service should initiate fund recovery
	finalBalance, err := user0.HTTPClient.BalanceAt(ctx, skAddress, nil)
	require.NoError(t, err)
	t.Logf("✓ Final session key balance: %s TEN", finalBalance.String())

	// 7) Poll user's pending balance until it increases (async refund tx inclusion)
	var finalUserBalance *big.Int
	{
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			bal, err := user0.HTTPClient.PendingBalanceAt(ctx, fromAddr)
			require.NoError(t, err)
			if bal.Cmp(initialUserBalance) > 0 {
				finalUserBalance = bal
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		if finalUserBalance == nil {
			// take one last reading for logging and fail with a clear message
			bal, _ := user0.HTTPClient.PendingBalanceAt(ctx, fromAddr)
			t.Logf("Pending balance did not increase within the timeout. last=%s, initial=%s", bal.String(), initialUserBalance.String())
			require.FailNow(t, "Timed out waiting for recovered funds to reflect in user's pending balance")
		}
	}
	t.Logf("✓ Final user balance: %s TEN", finalUserBalance.String())

	// 9) Verify that the user's balance increased by the recovered amount (minus gas costs)
	// The user should have received the funds back, minus some gas costs
	balanceIncrease := big.NewInt(0).Sub(finalUserBalance, initialUserBalance)
	t.Logf("✓ Balance increase: %s TEN", balanceIncrease.String())

	// The balance increase should be positive (user received funds back)
	// and should be close to the original fund amount (minus gas costs)
	require.True(t, balanceIncrease.Cmp(big.NewInt(0)) > 0, "User balance should have increased due to fund recovery")

	// The recovered amount should be at least 90% of the original fund amount
	// (allowing for gas costs)
	expectedMinRecovery := big.NewInt(0).Div(big.NewInt(0).Mul(fundAmount, big.NewInt(90)), big.NewInt(100))
	require.True(t, balanceIncrease.Cmp(expectedMinRecovery) >= 0,
		"Recovered amount should be at least 90%% of original fund amount")

	t.Logf("✓ Session key expiration and fund recovery test completed successfully!")
}

func testNewHeadsSubscription(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)

	receivedHeads := make([]*types.Header, 0)
	newHeadChan := make(chan *types.Header)
	subscription, err := user0.WSClient.SubscribeNewHead(context.Background(), newHeadChan)
	require.NoError(t, err)

	// Listen for new heads in a goroutine
	go func() {
		for {
			select {
			case err := <-subscription.Err():
				// if err != nil {
				testlog.Logger().Info("Error from new head subscription", log2.ErrKey, err)
				return
				//}
			case newHead := <-newHeadChan:
				// append logs to be visible from the main thread
				receivedHeads = append(receivedHeads, newHead)
			}
		}
	}()

	// sleep for 5 seconds and there should be at least 2 heads received in this interval
	time.Sleep(5 * time.Second)
	subscription.Unsubscribe()
	require.True(t, len(receivedHeads) > 1)
}

func testMultipleAccountsSubscription(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user0.tgClient.UserID())

	_, err = user0.HTTPClient.ChainID(context.Background())
	require.NoError(t, err)

	user1, err := NewGatewayUser([]wallet.Wallet{datagenerator.RandomWallet(integration.TenChainID), datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user1.tgClient.UserID())

	user2, err := NewGatewayUser([]wallet.Wallet{datagenerator.RandomWallet(integration.TenChainID), datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user2.tgClient.UserID())

	// register all the accounts for that user
	err = user0.RegisterAccountsPersonalSign()
	require.NoError(t, err)
	err = user1.RegisterAccountsPersonalSign()
	require.NoError(t, err)
	err = user2.RegisterAccountsPersonalSign()
	require.NoError(t, err)

	var amountToTransfer int64 = 1_000_000_000_000_000_000
	// Transfer some funds to user1 and user2 wallets, because they need it to make transactions
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user2.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)

	// Print balances of all registered accounts to check if all accounts have funds
	balances, err := user1.GetUserAccountsBalances()
	require.NoError(t, err)
	for _, balance := range balances {
		require.NotZero(t, balance.Uint64())
	}
	balances, err = user2.GetUserAccountsBalances()
	require.NoError(t, err)
	for _, balance := range balances {
		require.NotZero(t, balance.Uint64())
	}

	// deploy events contract
	// Get current nonce to avoid conflicts with previous tests
	currentNonce, err := user0.HTTPClient.PendingNonceAt(context.Background(), w.Address())
	require.NoError(t, err)

	deployTx := &types.LegacyTx{
		Nonce:    currentNonce,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	err = getFeeAndGas(user0.HTTPClient, w, deployTx)
	require.NoError(t, err)

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user0.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user0.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	_, err = user0.HTTPClient.CodeAt(context.Background(), contractReceipt.ContractAddress, big.NewInt(int64(rpc.LatestBlockNumber)))
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
	_, err = subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user0.WSClient, &user0logs)
	require.NoError(t, err)
	_, err = subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user1.WSClient, &user1logs)
	require.NoError(t, err)
	_, err = subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user2.WSClient, &user2logs)
	require.NoError(t, err)

	// user1 calls setMessage and setMessage2 on deployed smart contract with the account
	// that was registered as the first in TG
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
	// that was registered as the second in TG
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

	_, err = user0.HTTPClient.FilterLogs(context.TODO(), ethereum.FilterQuery{
		Addresses: []gethcommon.Address{contractReceipt.ContractAddress},
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(10000),
		Topics:    nil,
	})
	require.NoError(t, err)
}

func testSubscriptionTopics(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user0, err := NewGatewayUser([]wallet.Wallet{w}, httpURL, wsURL)
	require.NoError(t, err)

	user1, err := NewGatewayUser([]wallet.Wallet{datagenerator.RandomWallet(integration.TenChainID), datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)

	// register all the accounts for that user
	err = user0.RegisterAccounts()
	require.NoError(t, err)
	err = user1.RegisterAccounts()
	require.NoError(t, err)

	var amountToTransfer int64 = 1_000_000_000_000_000_000
	// Transfer some funds to user1 to be able to make transactions
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[0].Address(), amountToTransfer)
	require.NoError(t, err)
	_, err = transferETHToAddress(user0.HTTPClient, user0.Wallets[0], user1.Wallets[1].Address(), amountToTransfer)
	require.NoError(t, err)

	// Print balances of all registered accounts to check if all accounts have funds
	balances, err := user0.GetUserAccountsBalances()
	require.NoError(t, err)
	for _, balance := range balances {
		require.NotZero(t, balance.Uint64())
	}
	balances, err = user1.GetUserAccountsBalances()
	require.NoError(t, err)
	for _, balance := range balances {
		require.NotZero(t, balance.Uint64())
	}

	// deploy events contract
	// Get current nonce to avoid conflicts with previous tests
	currentNonce, err := user0.HTTPClient.PendingNonceAt(context.Background(), w.Address())
	require.NoError(t, err)

	deployTx := &types.LegacyTx{
		Nonce:    currentNonce,
		Gas:      uint64(10_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	err = getFeeAndGas(user0.HTTPClient, w, deployTx)
	require.NoError(t, err)

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user0.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user0.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	tx, _, err := user0.HTTPClient.TransactionByHash(context.Background(), signedTx.Hash())
	if err != nil {
		return
	}
	require.Equal(t, signedTx.Hash(), tx.Hash())

	// user0 subscribes to all events from that smart contract, user1 only an event with a topic of his first account
	var user0logs []types.Log
	var user1logs []types.Log
	var topics [][]gethcommon.Hash
	t1 := gethcommon.BytesToHash(user1.Wallets[1].Address().Bytes())
	topics = append(topics, nil)
	topics = append(topics, []gethcommon.Hash{t1})
	_, err = subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user0.WSClient, &user0logs)
	require.NoError(t, err)
	_, err = subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, topics, user1.WSClient, &user1logs)
	require.NoError(t, err)

	// user0 calls setMessage on deployed smart contract with the account twice and expects two events
	_, err = integrationCommon.InteractWithSmartContract(user0.HTTPClient, user0.Wallets[0], eventsContractABI, "setMessage", "user0Event1", contractReceipt.ContractAddress)
	require.NoError(t, err)
	_, err = integrationCommon.InteractWithSmartContract(user0.HTTPClient, user0.Wallets[0], eventsContractABI, "setMessage", "user0Event2", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// user1 calls setMessage on deployed smart contract with two different accounts and expects only one event,
	// because only the first address is in the topic filter of the subscription
	_, err = integrationCommon.InteractWithSmartContract(user1.HTTPClient, user1.Wallets[0], eventsContractABI, "setMessage", "user1Event1", contractReceipt.ContractAddress)
	require.NoError(t, err)
	_, err = integrationCommon.InteractWithSmartContract(user1.HTTPClient, user1.Wallets[1], eventsContractABI, "setMessage", "user1Event2", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// wait a few seconds to be completely sure all events arrived
	time.Sleep(time.Second * 3)

	// Assert the number of logs received by each client
	// user0 should see two lifecycle events (1 for each interaction with the smart contract)
	assert.Equal(t, 2, len(user0logs))
	// user1 should see only one event (the other is filtered out because of the topic filter)
	assert.Equal(t, 1, len(user1logs))
}

func testAreTxsMinted(t *testing.T, httpURL, wsURL string, w wallet.Wallet) { //nolint: unused
	// set up the tgClient
	ogClient := lib.NewTenGatewayLibrary(httpURL, wsURL)

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

func testErrorHandling(t *testing.T, startPort int, httpURL, wsURL string, w wallet.Wallet) {
	// set up the tgClient
	ogClient := lib.NewTenGatewayLibrary(httpURL, wsURL)

	// join + register against the og
	err := ogClient.Join()
	require.NoError(t, err)

	// register an account
	err = ogClient.RegisterAccount(w.PrivateKey(), w.Address())
	require.NoError(t, err)

	privateTxsBytes, _ := json.Marshal(common.ListPrivateTransactionsQueryParams{
		Address:          gethcommon.HexToAddress("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
		Pagination:       common.QueryPagination{Size: 10},
		ShowSyntheticTxs: false,
		ShowAllPublicTxs: false,
	})

	privateTxs := strings.ReplaceAll(string(privateTxsBytes), `"`, `\"`)

	// make requests to geth for comparison
	for _, req := range []string{
		`{"jsonrpc":"2.0","method":"eth_getStorageAt","params":["` + common.ListPrivateTransactionsCQMethod + `", "` + privateTxs + `","latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getLogs","params":[[]],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":[]}],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"fromBlock":"0x387","topics":["0xc6d8c0af6d21f291e7c359603aa97e0ed500f04db6e983b9fce75a91c6b8da6b"]}],"id":1}`,
		`{"jsonrpc":"2.0","method":"debug_eventLogRelevancy","params":[{"fromBlock":"0x387","topics":["0xc6d8c0af6d21f291e7c359603aa97e0ed500f04db6e983b9fce75a91c6b8da6b"]}],"id":1}`,
		//`{"jsonrpc":"2.0","method":"eth_subscribe","params":["logs"],"id":1}`,
		//`{"jsonrpc":"2.0","method":"eth_subscribe","params":["logs",{"topics":[]}],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`, // test caching
		`{"jsonrpc":"2.0","method":"eth_gasPrice","params": [],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_gasPrice","params": [],"id":1}`, // test caching
		`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["latest", false],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_feeHistory","params":[1, "latest", [50]],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":[],"id":1}`,
		//`{"jsonrpc":"2.0","method":"eth_getgetget","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "latest"],"id":1,"extra":"extra_field"}`,
		`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[["0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", "0x1234"]],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x0000000000000000000000000000000000000000000000000000000000000000"],"id":1}`,
		`{"jsonrpc":"2.0","method":"eth_maxPriorityFeePerGas","params":[],"id":1}`,
	} {
		// ensure the gateway request is issued correctly (should return 200 ok with jsonRPCError)
		_, response, err := httputil.PostDataJSON(ogClient.HTTP(), []byte(req))
		require.NoError(t, err)
		fmt.Printf("Resp: %s\n", response)

		// unmarshall the response to JSONRPCMessage
		jsonRPCError := JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err, req, response)

		// repeat the process for geth
		_, response, err = httputil.PostDataJSON(fmt.Sprintf("http://localhost:%d", startPort+integration.DefaultGethHTTPPortOffset), []byte(req))
		require.NoError(t, err)

		// we only care about format
		jsonRPCError = JSONRPCMessage{}
		err = json.Unmarshal(response, &jsonRPCError)
		require.NoError(t, err)
	}
}

func testErrorsRevertedArePassed(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	// set up the tgClient
	ogClient := lib.NewTenGatewayLibrary(httpURL, wsURL)

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
	// Get current nonce to avoid conflicts with previous tests
	currentNonce, err := ethStdClient.PendingNonceAt(context.Background(), w.Address())
	require.NoError(t, err)

	deployTx := &types.LegacyTx{
		Nonce:    currentNonce,
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(errorsContractBytecode),
	}

	err = getFeeAndGas(ethStdClient, w, deployTx)
	require.NoError(t, err)

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
	require.Equal(t, "execution reverted: Forced require", err.Error())

	// convert error to WE error
	errBytes, err := json.Marshal(err)
	require.NoError(t, err)
	weError := JSONError{}
	err = json.Unmarshal(errBytes, &weError)
	require.NoError(t, err)
	require.Equal(t, "execution reverted: Forced require", weError.Message)
	expectedData := "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e466f726365642072657175697265000000000000000000000000000000000000"
	require.Equal(t, expectedData, weError.Data)
	require.Equal(t, 3, weError.Code)

	pack, _ = errorsContractABI.Pack("force_revert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, "execution reverted: Forced revert", err.Error())

	pack, _ = errorsContractABI.Pack("force_assert")
	_, err = ethStdClient.CallContract(context.Background(), ethereum.CallMsg{
		From: w.Address(),
		To:   &receipt.ContractAddress,
		Data: pack,
	}, nil)
	require.Error(t, err)
	require.Equal(t, "execution reverted: assert(false)", err.Error())
}

func testUnsubscribe(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	// create a user with multiple accounts
	user, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user.tgClient.UserID())

	_, err = user.HTTPClient.ChainID(context.Background())
	require.NoError(t, err)

	// register all the accounts for the user
	err = user.RegisterAccounts()
	require.NoError(t, err)

	// deploy events contract
	// Get current nonce to avoid conflicts with previous tests
	currentNonce, err := user.HTTPClient.PendingNonceAt(context.Background(), w.Address())
	require.NoError(t, err)

	deployTx := &types.LegacyTx{
		Nonce:    currentNonce,
		Gas:      uint64(10_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	require.NoError(t, getFeeAndGas(user.HTTPClient, w, deployTx))

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	testlog.Logger().Info("Deployed contract address: ", "addr", contractReceipt.ContractAddress)

	// subscribe to an event
	var userLogs []types.Log
	subscription, err := subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user.WSClient, &userLogs)
	require.NoError(t, err)

	// make an action that will trigger events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage", "foo", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// wait for the first log to arrive (subscription consumes logs asynchronously)
	{
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if len(userLogs) >= 1 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		assert.Equal(t, 1, len(userLogs))
	}

	// Unsubscribe from events
	subscription.Unsubscribe()

	// make another action that will trigger events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage", "bar", contractReceipt.ContractAddress)
	require.NoError(t, err)

	// give a short time window to ensure no new logs are received after unsubscribing
	time.Sleep(1 * time.Second)
	// check that we are not receiving events after unsubscribing
	assert.Equal(t, 1, len(userLogs))
}

func testClosingConnectionWhileSubscribed(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	// create a user with multiple accounts
	user, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user.tgClient.UserID())

	_, err = user.HTTPClient.ChainID(context.Background())
	require.NoError(t, err)

	// register all the accounts for the user
	err = user.RegisterAccounts()
	require.NoError(t, err)

	// deploy events contract
	// Get current nonce to avoid conflicts with previous tests
	currentNonce, err := user.HTTPClient.PendingNonceAt(context.Background(), w.Address())
	require.NoError(t, err)

	deployTx := &types.LegacyTx{
		Nonce:    currentNonce,
		Gas:      uint64(10_000_000),
		GasPrice: gethcommon.Big1,
		Data:     gethcommon.FromHex(eventsContractBytecode),
	}

	require.NoError(t, getFeeAndGas(user.HTTPClient, w, deployTx))

	signedTx, err := w.SignTransaction(deployTx)
	require.NoError(t, err)

	err = user.HTTPClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	contractReceipt, err := integrationCommon.AwaitReceiptEth(context.Background(), user.HTTPClient, signedTx.Hash(), time.Minute)
	require.NoError(t, err)

	testlog.Logger().Info("Deployed contract address: ", "addr", contractReceipt.ContractAddress)

	// subscribe to an event
	var userLogs []types.Log
	subscription, err := subscribeToEvents([]gethcommon.Address{contractReceipt.ContractAddress}, nil, user.WSClient, &userLogs)
	require.NoError(t, err)

	// Close the websocket connection and make sure nothing breaks, but user does not receive events
	user.WSClient.Close()

	// make an action that will emmit events
	_, err = integrationCommon.InteractWithSmartContract(user.HTTPClient, user.Wallets[0], eventsContractABI, "setMessage2", "foo", contractReceipt.ContractAddress)
	require.NoError(t, err)
	// but with closed connection we don't receive any logs
	assert.Equal(t, 0, len(userLogs))

	// re-establish connection
	wsClient, err := ethclient.Dial(wsURL + "/v1/" + "?token=" + user.tgClient.UserID())
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

func testDifferentMessagesOnRegister(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user, err := NewGatewayUser([]wallet.Wallet{w, datagenerator.RandomWallet(integration.TenChainID)}, httpURL, wsURL)
	require.NoError(t, err)
	testlog.Logger().Info("Created user with encryption token", "t", user.tgClient.UserID())

	// register all the accounts for the user with EIP-712 message format
	err = user.RegisterAccounts()
	require.NoError(t, err)

	// register all the accounts for the user with personal sign message format
	err = user.RegisterAccountsPersonalSign()
	require.NoError(t, err)
}

func testInvokeNonSensitiveMethod(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	user, err := NewGatewayUser([]wallet.Wallet{w}, httpURL, wsURL)
	require.NoError(t, err)

	// call one of the non-sensitive methods with unauthenticated user
	// and make sure gateway is not complaining about not having viewing keys
	respBody := makeHTTPEthJSONReq(httpURL, "eth_chainId", user.tgClient.UserID(), nil)
	if strings.Contains(string(respBody), fmt.Sprintf("method %s cannot be called with an unauthorised client - no signed viewing keys found", "eth_chainId")) {
		t.Errorf("sensitive method called without authenticating viewingkeys and did fail because of it:  %s", "eth_chainId")
	}
}

func testQueryAndRpcTokenModes(t *testing.T, _ int, httpURL, wsURL string, w wallet.Wallet) {
	// 1) Create a user and authenticate (register address)
	user, err := NewGatewayUser([]wallet.Wallet{w}, httpURL, wsURL)
	require.NoError(t, err)

	// Register the user so REST /query can validate
	require.NoError(t, user.RegisterAccounts())

	// 2) Call REST /v1/query/?token=...&a=<address>
	addrHex := user.Wallets[0].Address().Hex()
	queryURL := fmt.Sprintf("%s/v1/query/?token=%s&a=%s", httpURL, user.tgClient.UserID(), addrHex)
	status, body, err := fasthttp.Get(nil, queryURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	// Response format: {"status": true|false}
	type queryResp struct {
		Status bool `json:"status"`
	}
	var qres queryResp
	require.NoError(t, json.Unmarshal(body, &qres))
	require.True(t, qres.Status, "expected registered address to be found")

	// 3) Call JSON-RPC eth_getBalance in two ways and compare results
	ethClientWithQuery, err := ethclient.Dial(fmt.Sprintf("%s/v1/?token=%s", httpURL, user.tgClient.UserID()))
	require.NoError(t, err)
	defer ethClientWithQuery.Close()

	ethClientWithPath, err := ethclient.Dial(fmt.Sprintf("%s/v1/%s", httpURL, user.tgClient.UserID()))
	require.NoError(t, err)
	defer ethClientWithPath.Close()

	balanceQuery, err := ethClientWithQuery.BalanceAt(context.Background(), user.Wallets[0].Address(), nil)
	require.NoError(t, err)

	balancePath, err := ethClientWithPath.BalanceAt(context.Background(), user.Wallets[0].Address(), nil)
	require.NoError(t, err)

	require.Equal(t, 0, balanceQuery.Cmp(balancePath), "balances via query vs path token should match")
}

func makeRequestHTTP(url string, body []byte) []byte {
	generateViewingKeyBody := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", generateViewingKeyBody) //nolint:noctx,gosec
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	viewingKey, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return viewingKey
}

func makeHTTPEthJSONReq(url string, method string, userID string, params interface{}) []byte {
	reqBody := prepareRequestBody(method, params)
	return makeRequestHTTP(fmt.Sprintf("%s/v1/?token=%s", url, userID), reqBody)
}

func prepareRequestBody(method string, params interface{}) []byte {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		wecommon.JSONKeyRPCVersion: jsonrpc.Version,
		wecommon.JSONKeyMethod:     method,
		wecommon.JSONKeyParams:     params,
		wecommon.JSONKeyID:         "1",
	})
	if err != nil {
		panic(fmt.Errorf("failed to prepare request body. Cause: %w", err))
	}
	return reqBodyBytes
}

func transferRandomAddr(t *testing.T, client *ethclient.Client, w wallet.Wallet) common.TxHash { //nolint: unused
	ctx := context.Background()
	toAddr := datagenerator.RandomAddress()
	nonce, err := client.NonceAt(ctx, w.Address(), nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := &types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	}
	assert.Nil(t, err)

	testlog.Logger().Info("Transferring from:", "addr", w.Address(), " to:", toAddr)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = client.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

	_, err = integrationCommon.AwaitReceiptEth(context.Background(), client, signedTx.Hash(), time.Minute)
	assert.NoError(t, err)

	testlog.Logger().Info("Successfully minted the transaction - ", "tx", signedTx.Hash())
	return signedTx.Hash()
}

// Creates a single-node TEN network for testing.
func createTenNetwork(t *testing.T, startPort int) {
	// Create the TEN network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.TenChainID)
	simParams := params.SimParams{
		NumberOfNodes:       numberOfNodes,
		AvgBlockDuration:    2 * time.Second,
		ContractRegistryLib: ethereummock.NewContractRegistryLibMock(),
		ERC20ContractLib:    ethereummock.NewERC20ContractLibMock(),
		Wallets:             wallets,
		StartPort:           startPort,
		WithPrefunding:      true,
		L1BeaconPort:        integration.TestPorts.TestTenGatewayPort + integration.DefaultPrysmGatewayPortOffset,
	}

	tenNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(tenNetwork.TearDown)
	_, err := tenNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test TEN network. Cause: %s", err))
	}
}

func waitServerIsReady(serverAddr string) error {
	for now := time.Now(); time.Since(now) < 30*time.Second; time.Sleep(500 * time.Millisecond) {
		statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/v1/health/", serverAddr))
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

func getFeeAndGas(client *ethclient.Client, wallet wallet.Wallet, legacyTx *types.LegacyTx) error {
	tx := types.NewTx(legacyTx)

	history, err := client.FeeHistory(context.Background(), 1, nil, nil)
	if err != nil || len(history.BaseFee) == 0 {
		return err
	}

	estimate, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		// From:  wallet.Address(),
		To:    tx.To(),
		Value: tx.Value(),
		Data:  tx.Data(),
	})
	if err != nil {
		return err
	}

	legacyTx.Gas = estimate
	legacyTx.GasPrice = history.BaseFee[0] // big.NewInt(gethparams.InitialBaseFee)

	return nil
}

func transferETHToAddress(client *ethclient.Client, wallet wallet.Wallet, toAddress gethcommon.Address, amount int64) (*types.Receipt, error) { //nolint:unparam
	// Get current nonce to avoid conflicts
	currentNonce, err := client.PendingNonceAt(context.Background(), wallet.Address())
	if err != nil {
		return nil, err
	}

	transferTx1 := types.LegacyTx{
		Nonce:    currentNonce,
		To:       &toAddress,
		Value:    big.NewInt(amount),
		Gas:      uint64(10_000_000),
		GasPrice: gethcommon.Big1,
		Data:     nil,
	}

	err = getFeeAndGas(client, wallet, &transferTx1)
	if err != nil {
		return nil, err
	}

	signedTx, err := wallet.SignTransaction(&transferTx1)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return integrationCommon.AwaitReceiptEth(context.Background(), client, signedTx.Hash(), 30*time.Second)
}

func subscribeToEvents(addresses []gethcommon.Address, topics [][]gethcommon.Hash, client *ethclient.Client, logs *[]types.Log) (ethereum.Subscription, error) {
	// Make a subscription
	filterQuery := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: big.NewInt(2),
		// ToBlock:   big.NewInt(10000),
		Topics: topics,
	}
	logsCh := make(chan types.Log)

	subscription, err := client.SubscribeFilterLogs(context.Background(), filterQuery, logsCh)
	if err != nil {
		testlog.Logger().Info("Failed to subscribe to filter logs", log2.ErrKey, err)
		return nil, err
	}

	// Listen for logs in a goroutine
	go func() {
		for {
			select {
			case err := <-subscription.Err():
				testlog.Logger().Info("Error from logs subscription", log2.ErrKey, err)
				return
			case log := <-logsCh:
				// append logs to be visible from the main thread
				*logs = append(*logs, log)
			}
		}
	}()

	return subscription, nil
}
