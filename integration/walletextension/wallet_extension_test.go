package walletextension

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/obscuronet/go-obscuro/tools/walletextension"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	testLogs     = "../.build/wallet_extension/"
	l2ChainIDHex = "0x309"

	reqJSONKeyTo      = "to"
	reqJSONKeyFrom    = "from"
	reqJSONKeyData    = "data"
	respJSONKeyStatus = "status"
	latestBlock       = "latest"
	statusSuccess     = "0x1"
	errInsecure       = "enclave could not respond securely to %s request"

	networkStartPort = integration.StartPortWalletExtensionTest + 1
	nodeRPCHTTPPort  = integration.StartPortWalletExtensionTest + 1 + network.DefaultHostRPCHTTPOffset
	nodeRPCWSPort    = integration.StartPortWalletExtensionTest + 1 + network.DefaultHostRPCWSOffset
	httpProtocol     = "http://"

	// Returned by the EVM to indicate a zero result.
	zeroResult  = "0x0000000000000000000000000000000000000000000000000000000000000000"
	zeroBalance = "0x0"
)

var (
	// The log file used across all the wallet extension tests.
	logFile = testlog.Setup(
		&testlog.Cfg{LogDir: testLogs, TestType: "wal-ext", TestSubtype: "test"},
	)

	walletExtensionAddr   = fmt.Sprintf("%s:%d", network.Localhost, integration.StartPortWalletExtensionTest)
	walletExtensionConfig = walletextension.Config{
		WalletExtensionPort:     int(integration.StartPortWalletExtensionTest),
		NodeRPCHTTPAddress:      fmt.Sprintf("%s:%d", network.Localhost, nodeRPCHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("%s:%d", network.Localhost, nodeRPCWSPort),
	}
	
	dummyAccountAddress = common.HexToAddress("0x8D97689C9818892B700e27F316cc3E41e17fBeb9")
	deployERC20Tx       = types.LegacyTx{
		Gas:      1025_000_000,
		GasPrice: common.Big0,
		Data:     erc20contract.L2BytecodeWithDefaultSupply("TST"),
	}
)

func TestCanMakeNonSensitiveRequestWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("req-no-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)

	respJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCChainID, []string{})

	if respJSON[walletextension.RespJSONKeyResult] != l2ChainIDHex {
		t.Fatalf("Expected chainId of %s, got %s", l2ChainIDHex, respJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotGetBalanceWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("bal-no-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCGetBalance, []string{dummyAccountAddress.Hex(), latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCGetBalance)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message to contain \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanGetOwnBalanceAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("bal-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	accountAddr, _ := registerPrivateKey(t)

	getBalanceJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCGetBalance, []string{accountAddr.String(), latestBlock})

	if getBalanceJSON[walletextension.RespJSONKeyResult] != zeroBalance {
		t.Fatalf("Expected balance of %s, got %s", zeroBalance, getBalanceJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotGetAnothersBalanceAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-bal-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	registerPrivateKey(t)

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCGetBalance, []string{dummyAccountAddress.Hex(), latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCGetBalance)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message to contain \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCannotCallWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("call-no-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)

	// We generate an account, but do not register it with the node.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	transferTxBytes := erc20contractlib.CreateTransferTxData(accountAddress, 0)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WOBXContract,
		reqJSONKeyFrom: accountAddress.String(),
		reqJSONKeyData: "0x" + common.Bytes2Hex(transferTxBytes),
	}

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCCall, []interface{}{reqParams, latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCCall)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanCallAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("call-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	accountAddress, _ := registerPrivateKey(t)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	balanceData := erc20contractlib.CreateBalanceOfData(accountAddress)
	convertedData := (hexutil.Bytes)(balanceData)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WOBXContract.Hex(),
		reqJSONKeyFrom: accountAddress.String(),
		reqJSONKeyData: convertedData,
	}

	callJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCCall, []interface{}{reqParams, latestBlock})

	if callJSON[walletextension.RespJSONKeyResult] != zeroResult {
		t.Fatalf("Expected call result of %s, got %s", zeroResult, callJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCanCallWithoutSettingFromField(t *testing.T) {
	setupWalletTestLog("call-no-from-field")

	createWalletExtension(t)
	createObscuroNetwork(t)
	accountAddress, _ := registerPrivateKey(t)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	balanceData := erc20contractlib.CreateBalanceOfData(accountAddress)
	convertedData := (hexutil.Bytes)(balanceData)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.WOBXContract,
		reqJSONKeyData: convertedData,
	}

	callJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCCall, []interface{}{reqParams, latestBlock})

	if callJSON[walletextension.RespJSONKeyResult] != zeroResult {
		t.Fatalf("Expected call result of %s, got %s", zeroResult, callJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotCallForAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-call-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	registerPrivateKey(t)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	balanceData := erc20contractlib.CreateBalanceOfData(dummyAccountAddress)
	convertedData := (hexutil.Bytes)(balanceData)
	reqParams := map[string]interface{}{
		reqJSONKeyTo: bridge.WOBXContract,
		// We send the request from a different address than the one we created a viewing key for.
		reqJSONKeyFrom: dummyAccountAddress.Hex(),
		reqJSONKeyData: convertedData,
	}

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCCall, []interface{}{reqParams, latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCCall)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCannotSubmitTxWithoutSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("submit-tx-no-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	txBinaryHex := signAndSerialiseTransaction(txWallet, &deployERC20Tx)

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCSendRawTransaction, []interface{}{txBinaryHex})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCSendRawTransaction)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanSubmitTxAndGetTxReceiptAndTxAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("submit-tx-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	_, privateKey := registerPrivateKey(t)

	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	signedTx, err := txWallet.SignTransaction(&deployERC20Tx)
	if err != nil {
		panic(fmt.Errorf("could not sign transaction. Cause: %w", err))
	}

	// We check the transaction receipt contains the correct transaction hash.
	txReceiptJSON := sendTransactionAndAwaitConfirmation(txWallet, deployERC20Tx)
	txReceiptResult := fmt.Sprintf("%s", txReceiptJSON[walletextension.RespJSONKeyResult])
	expectedTxReceiptJSON := fmt.Sprintf("transactionHash:%s", signedTx.Hash())
	if !strings.Contains(txReceiptResult, expectedTxReceiptJSON) {
		t.Fatalf("Expected transaction receipt containing %s, got %s", expectedTxReceiptJSON, txReceiptResult)
	}

	// We check we can retrieve the transaction by hash.
	getTxJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCGetTransactionByHash, []string{signedTx.Hash().Hex()})
	getTxJSONResult := fmt.Sprintf("%s", getTxJSON[walletextension.RespJSONKeyResult])
	expectedGetTxJSON := fmt.Sprintf("hash:%s", signedTx.Hash())
	if !strings.Contains(getTxJSONResult, expectedGetTxJSON) {
		t.Fatalf("Expected transaction containing %s, got %s", expectedGetTxJSON, getTxJSONResult)
	}
}

func TestCannotSubmitTxFromAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	setupWalletTestLog("others-submit-tx-with-viewing-key")

	createWalletExtension(t)
	createObscuroNetwork(t)
	registerPrivateKey(t)

	// We submit a transaction using another account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	txBinaryHex := signAndSerialiseTransaction(txWallet, &deployERC20Tx)

	respBody := makeEthJSONReq(walletExtensionAddr, rpcclientlib.RPCSendRawTransaction, []interface{}{txBinaryHex})
	expectedErr := fmt.Sprintf(errInsecure, rpcclientlib.RPCSendRawTransaction)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanDecryptSuccessfullyAfterSubmittingMultipleViewingKeys(t *testing.T) {
	setupWalletTestLog("bal-with-mult-viewing-keys")

	createWalletExtension(t)
	createObscuroNetwork(t)

	// We submit a viewing key for a random account.
	var accountAddrs []string
	for i := 0; i < 10; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		accountAddr := crypto.PubkeyToAddress(privateKey.PublicKey).String()
		generateAndSubmitViewingKey(accountAddr, privateKey)
		accountAddrs = append(accountAddrs, accountAddr)
	}

	// We request the balance of a random account about halfway through the list.
	randAccountAddr := accountAddrs[len(accountAddrs)/2]
	getBalanceJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCGetBalance, []string{randAccountAddr, latestBlock})

	if getBalanceJSON[walletextension.RespJSONKeyResult] != zeroBalance {
		t.Fatalf("Expected balance of %s, got %s", zeroBalance, getBalanceJSON[walletextension.RespJSONKeyResult])
	}
}

// Creates and serves a wallet extension.
func createWalletExtension(t *testing.T) {
	walletExtension := walletextension.NewWalletExtension(walletExtensionConfig)
	t.Cleanup(walletExtension.Shutdown)
	go walletExtension.Serve(walletExtensionAddr)
	waitForWalletExtension(t, walletExtensionAddr)
}

// Waits for wallet extension to be ready. Times out after three seconds.
func waitForWalletExtension(t *testing.T, walletExtensionAddr string) {
	retries := 30
	for i := 0; i < retries; i++ {
		resp, err := http.Get(httpProtocol + walletExtensionAddr + walletextension.PathReady) //nolint:noctx
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if err == nil {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatal("could not establish connection to wallet extension")
}

// Makes an Ethereum JSON RPC request and returns the response body.
func makeEthJSONReq(walletExtensionAddr string, method string, params interface{}) []byte {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "1",
	})
	if err != nil {
		panic(err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)

	var resp *http.Response
	// We retry for three seconds to handle node start-up time.
	timeout := time.Now().Add(3 * time.Second)
	for i := time.Now(); i.Before(timeout); i = time.Now() {
		resp, err = http.Post(httpProtocol+walletExtensionAddr, "text/html", reqBody) //nolint:noctx
		if err == nil {
			break
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

	if err != nil {
		panic(fmt.Errorf("received error response from wallet extension: %w", err))
	}
	if resp == nil {
		panic("did not receive a response from the wallet extension")
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return respBody
}

// Makes an Ethereum JSON RPC request and returns the response body as JSON.
func makeEthJSONReqAsJSON(method string, params interface{}) map[string]interface{} {
	respBody := makeEthJSONReq(walletExtensionAddr, method, params)

	if respBody[0] != '{' {
		panic(fmt.Errorf("expected JSON response but received: %s", respBody))
	}

	var respBodyJSON map[string]interface{}
	err := json.Unmarshal(respBody, &respBodyJSON)
	if err != nil {
		panic(err)
	}

	return respBodyJSON
}

// Generates a signed viewing key and submits it to the wallet extension.
func generateAndSubmitViewingKey(accountAddr string, accountPrivateKey *ecdsa.PrivateKey) {
	viewingKey := generateViewingKey(accountAddr, walletExtensionAddr)
	signature := signViewingKey(accountPrivateKey, viewingKey)

	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeySignature: hex.EncodeToString(signature),
		walletextension.ReqJSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}
	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	resp, err := http.Post(httpProtocol+walletExtensionAddr+walletextension.PathSubmitViewingKey, "application/json", submitViewingKeyBody) //nolint:noctx
	if err != nil {
		panic(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(fmt.Errorf("request to add viewing key failed with following status: %s", resp.Status))
	}
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

// Generates a viewing key.
func generateViewingKey(accountAddress string, walletExtensionAddr string) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}
	generateViewingKeyBody := bytes.NewBuffer(generateViewingKeyBodyBytes)
	resp, err := http.Post(httpProtocol+walletExtensionAddr+walletextension.PathGenerateViewingKey, "application/json", generateViewingKeyBody) //nolint:noctx
	if err != nil {
		panic(err)
	}
	viewingKey, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return viewingKey
}

// Signs a viewing key.
func signViewingKey(privateKey *ecdsa.PrivateKey, viewingKey []byte) []byte {
	msgToSign := rpc.ViewingKeySignedMsgPrefix + string(viewingKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), privateKey)
	if err != nil {
		panic(err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	return signatureWithLeadBytes
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
		StartPort:        int(networkStartPort),
	}
	simStats := stats.NewStats(simParams.NumberOfNodes)
	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, simStats)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}

	// Set up the ERC20 wallet.
	erc20Wallet := wallets.Tokens[bridge.OBX].L2Owner
	generateAndSubmitViewingKey(walletExtensionAddr, erc20Wallet.PrivateKey())

	sendTransactionAndAwaitConfirmation(erc20Wallet, deployERC20Tx)
}

// Generates a new account and registers it with the node.
func registerPrivateKey(t *testing.T) (common.Address, *ecdsa.PrivateKey) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	accountAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	generateAndSubmitViewingKey(accountAddr.String(), privateKey)
	return accountAddr, privateKey
}

// Submits a transaction and awaits the transaction receipt.
func sendTransactionAndAwaitConfirmation(txWallet wallet.Wallet, tx types.LegacyTx) map[string]interface{} {
	// Set the transaction's nonce.
	nonceJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCNonce, []interface{}{txWallet.Address().Hex(), latestBlock})
	nonceString, ok := nonceJSON[walletextension.RespJSONKeyResult].(string)
	if !ok {
		panic(fmt.Errorf("retrieved nonce was not of type string"))
	}
	nonce, err := hexutil.DecodeUint64(nonceString)
	if err != nil {
		panic(fmt.Errorf("could not parse nonce from string. Cause: %w", err))
	}
	tx.Nonce = nonce

	// Send the transaction.
	txBinaryHex := signAndSerialiseTransaction(txWallet, &tx)
	sendTxJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCSendRawTransaction, []interface{}{txBinaryHex})

	// Verify the transaction was successful.
	txHash, ok := sendTxJSON[walletextension.RespJSONKeyResult].(string)
	if !ok {
		panic("could not retrieve transaction hash from JSON result, failed to deploy ERC20")
	}

	counter := 0
	for {
		if counter > 10 {
			panic("could not get ERC20 receipt after 10 seconds")
		}
		getReceiptJSON := makeEthJSONReqAsJSON(rpcclientlib.RPCGetTxReceipt, []interface{}{txHash})
		getReceiptJSONResult, ok := getReceiptJSON[walletextension.RespJSONKeyResult].(map[string]interface{})
		if ok && getReceiptJSONResult[respJSONKeyStatus] == statusSuccess {
			return getReceiptJSON
		}
		time.Sleep(1 * time.Second)
		counter++
	}
}

// Signs and serialises a transaction for submission to the node.
func signAndSerialiseTransaction(wallet wallet.Wallet, tx types.TxData) string {
	signedTx, err := wallet.SignTransaction(tx)
	if err != nil {
		panic(err)
	}
	// We convert the transaction to the form expected for sending transactions via RPC.
	txBinary, err := signedTx.MarshalBinary()
	if err != nil {
		panic(err)
	}
	txBinaryHex := "0x" + common.Bytes2Hex(txBinary)
	if err != nil {
		panic(err)
	}

	return txBinaryHex
}

func setupWalletTestLog(testName string) {
	// We reuse the same file for every test.
	log.OutputToFile(logFile)

	log.Info("-----------")
	log.Info("Starting test: %s", testName)
	log.Info("-----------")
}
