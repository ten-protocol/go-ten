package walletextension

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/test"

	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	"github.com/gorilla/websocket"

	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/erc20contract"

	"github.com/obscuronet/go-obscuro/tools/walletextension"

	gethcommon "github.com/ethereum/go-ethereum/common"
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

	reqJSONKeyTo            = "to"
	reqJSONKeyFrom          = "from"
	reqJSONKeyData          = "data"
	respJSONKeyStatus       = "status"
	respJSONKeyContractAddr = "contractAddress"
	latestBlock             = "latest"
	statusSuccess           = "0x1"
	errInsecure             = "enclave could not respond securely to %s request"
	errInvalidRPCMethod     = "rpc request failed: the method %s does not exist/is not available"

	walletExtensionPort   = int(integration.StartPortWalletExtensionTest)
	walletExtensionPortWS = int(integration.StartPortWalletExtensionTest + 1)
	networkStartPort      = integration.StartPortWalletExtensionTest + 2
	nodeRPCHTTPPort       = networkStartPort + network.DefaultHostRPCHTTPOffset
	nodeRPCWSPort         = networkStartPort + network.DefaultHostRPCWSOffset

	// Returned by the EVM to indicate a zero result.
	zeroResult  = "0x0000000000000000000000000000000000000000000000000000000000000000"
	zeroBalance = "0x0"

	faucetAlloc = 750000000000000 // The amount the faucet allocates to each Obscuro wallet.
)

var (
	walletExtensionAddrHTTP = fmt.Sprintf("http://%s:%d", network.Localhost, walletExtensionPort)
	walletExtensionAddrWS   = fmt.Sprintf("ws://%s:%d", network.Localhost, walletExtensionPortWS)
	walletExtensionConfig   = createWalletExtensionConfig()

	dummyAccountAddress = gethcommon.HexToAddress("0x8D97689C9818892B700e27F316cc3E41e17fBeb9")
	deployERC20Tx       = types.LegacyTx{
		Gas:      1025_000_000,
		GasPrice: gethcommon.Big1,
		Data:     erc20contract.L2BytecodeWithDefaultSupply("TST"),
	}
)

func TestMain(m *testing.M) {
	log.OutputToFile(testlog.Setup(&testlog.Cfg{LogDir: testLogs, TestType: "wal-ext", TestSubtype: "test"}))

	// We share a single Obscuro network across tests. Otherwise, every test takes 20 seconds at a minimum.
	teardown, err := createObscuroNetwork()
	defer teardown()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestCannotGetAnothersBalanceAfterSubmittingViewingKey(t *testing.T) {
	createWalletExtension(t)
	test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	respBody := test.MakeHTTPEthJSONReq(walletExtensionAddrHTTP, rpc.RPCGetBalance, []string{dummyAccountAddress.Hex(), latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpc.RPCGetBalance)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message to contain \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanCallWithoutSettingFromField(t *testing.T) {
	createWalletExtension(t)
	accountAddress, _, _ := test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	balanceData := erc20contractlib.CreateBalanceOfData(accountAddress)
	convertedData := (hexutil.Bytes)(balanceData)
	reqParams := map[string]interface{}{
		reqJSONKeyTo:   bridge.HOCContract,
		reqJSONKeyData: convertedData,
	}

	callJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCCall, []interface{}{reqParams, latestBlock})

	if callJSON[walletextension.RespJSONKeyResult] != zeroResult {
		t.Fatalf("Expected call result of %s, got %s", zeroResult, callJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCannotCallForAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	createWalletExtension(t)
	test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	// We submit a transaction to the Obscuro ERC20 contract. By transferring an amount of zero, we avoid the need to
	// deposit any funds in the ERC20 contract.
	balanceData := erc20contractlib.CreateBalanceOfData(dummyAccountAddress)
	convertedData := (hexutil.Bytes)(balanceData)
	reqParams := map[string]interface{}{
		reqJSONKeyTo: bridge.HOCContract,
		// We send the request from a different address than the one we created a viewing key for.
		reqJSONKeyFrom: dummyAccountAddress.Hex(),
		reqJSONKeyData: convertedData,
	}

	respBody := test.MakeHTTPEthJSONReq(walletExtensionAddrHTTP, rpc.RPCCall, []interface{}{reqParams, latestBlock})
	expectedErr := fmt.Sprintf(errInsecure, rpc.RPCCall)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCannotSubmitTxFromAnotherAddressAfterSubmittingViewingKey(t *testing.T) {
	createWalletExtension(t)
	test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	// We submit a transaction using another account.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	txBinaryHex := signAndSerialiseTransaction(txWallet, &deployERC20Tx)

	respBody := test.MakeHTTPEthJSONReq(walletExtensionAddrHTTP, rpc.RPCSendRawTransaction, []interface{}{txBinaryHex})
	expectedErr := fmt.Sprintf(errInsecure, rpc.RPCSendRawTransaction)

	if !strings.Contains(string(respBody), expectedErr) {
		t.Fatalf("Expected error message \"%s\", got \"%s\"", expectedErr, respBody)
	}
}

func TestCanDecryptSuccessfullyAfterSubmittingMultipleViewingKeys(t *testing.T) {
	createWalletExtension(t)

	// We submit a viewing key for a random account.
	var accountAddrs []string
	for i := 0; i < 10; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		accountAddr := crypto.PubkeyToAddress(privateKey.PublicKey).String()
		_, err = test.GenerateAndSubmitViewingKey(walletExtensionAddrHTTP, accountAddr, privateKey)
		if err != nil {
			t.Fatalf(err.Error())
		}
		accountAddrs = append(accountAddrs, accountAddr)
	}

	// We request the balance of a random account about halfway through the list.
	randAccountAddr := accountAddrs[len(accountAddrs)/2]
	getBalanceJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCGetBalance, []string{randAccountAddr, latestBlock})

	if getBalanceJSON[walletextension.RespJSONKeyResult] != zeroBalance {
		t.Fatalf("Expected balance of %s, got %s", zeroBalance, getBalanceJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCanDecryptSuccessfullyAfterRestartingWalletExtension(t *testing.T) {
	walletExtension := createWalletExtension(t)
	accountAddr, _, _ := test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	// We shut down the wallet extension and restart it, forcing the viewing keys to be reloaded.
	walletExtension.Shutdown()
	createWalletExtension(t)

	getBalanceJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCGetBalance, []string{accountAddr.String(), latestBlock})

	if getBalanceJSON[walletextension.RespJSONKeyResult] != zeroBalance {
		t.Fatalf("Expected balance of %s, got %s", zeroBalance, getBalanceJSON[walletextension.RespJSONKeyResult])
	}
}

func TestCanGetErrorOverWS(t *testing.T) {
	createWalletExtension(t)

	invalidMethod := "invalidRPCMethod"
	respJSON, _ := makeWSEthJSONReqAsJSON(invalidMethod, []string{})

	expectedErr := fmt.Sprintf(errInvalidRPCMethod, invalidMethod)
	if respJSON[userconn.RespJSONKeyErr] != expectedErr {
		t.Fatalf("Expected error '%s', got '%s'", expectedErr, respJSON[userconn.RespJSONKeyErr])
	}
}

func TestCanSubscribeForLogs(t *testing.T) {
	createWalletExtension(t)
	test.RegisterPrivateKey(t, walletExtensionAddrHTTP)

	_, conn := makeWSEthJSONReqAsJSON(rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filterCriteriaJSON{}})

	// We watch the connection for events...
	var receivedLogJSON []byte
	go func() {
		var err error
		_, receivedLogJSON, err = conn.ReadMessage()
		if err != nil {
			panic(fmt.Errorf("could not read log from websocket. Cause: %w", err))
		}
	}()

	// ... then trigger an event...
	txReceiptJSON := triggerEvent(t)

	// ... and wait up to thirty seconds for the event to be received.
	for i := 0; i < 30; i++ {
		if receivedLogJSON != nil {
			break
		}
		time.Sleep(time.Second)
	}
	if receivedLogJSON == nil {
		t.Fatalf("waited for 30 seconds without receiving a log")
	}

	// We convert the received JSON to a log object.
	var receivedLog *types.Log
	err := json.Unmarshal(receivedLogJSON, &receivedLog)
	if err != nil {
		t.Fatalf("could not unmarshall received log from JSON")
	}

	// We check the event we received was emitted by the expected contract.
	contractAddr := txReceiptJSON[walletextension.RespJSONKeyResult].(map[string]interface{})[respJSONKeyContractAddr].(string)
	logAddrLowercase := strings.ToLower(contractAddr)
	if logAddrLowercase != contractAddr {
		t.Fatalf("Expected event with contract address '%s', got '%s'", logAddrLowercase, contractAddr)
	}
}

func createWalletExtensionConfig() *walletextension.Config {
	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		panic("could not create persistence file for wallet extension tests")
	}

	return &walletextension.Config{
		WalletExtensionPort:     walletExtensionPort,
		WalletExtensionPortWS:   walletExtensionPortWS,
		NodeRPCHTTPAddress:      fmt.Sprintf("%s:%d", network.Localhost, nodeRPCHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("%s:%d", network.Localhost, nodeRPCWSPort),
		PersistencePathOverride: testPersistencePath.Name(),
	}
}

// Creates and serves a wallet extension.
func createWalletExtension(t *testing.T) *walletextension.WalletExtension {
	return createWalletExtensionWithConfig(t, walletExtensionConfig)
}

// Creates and serves a wallet extension with custom configuration.
func createWalletExtensionWithConfig(t *testing.T, config *walletextension.Config) *walletextension.WalletExtension {
	walletExtension := walletextension.NewWalletExtension(*config)
	t.Cleanup(walletExtension.Shutdown)

	go walletExtension.Serve(network.Localhost, walletExtensionPort, walletExtensionPortWS)
	err := test.WaitForEndpoint(walletExtensionAddrHTTP)
	if err != nil {
		t.Fatal(err)
	}

	return walletExtension
}

// Makes an Ethereum JSON RPC request over HTTP and returns the response body as JSON.
func makeHTTPEthJSONReqAsJSON(method string, params interface{}) map[string]interface{} {
	respBody := test.MakeHTTPEthJSONReq(walletExtensionAddrHTTP, method, params)
	return convertRespBodyToJSON(respBody)
}

// Makes an Ethereum JSON RPC request over websockets and returns the response body as JSON.
func makeWSEthJSONReqAsJSON(method string, params interface{}) (map[string]interface{}, *websocket.Conn) {
	respBody, conn := test.MakeWSEthJSONReq(walletExtensionAddrWS, method, params)
	return convertRespBodyToJSON(respBody), conn
}

// Converts the response body bytes to JSON.
func convertRespBodyToJSON(respBody []byte) map[string]interface{} {
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

// Creates a single-node Obscuro network for testing, and deploys an ERC20 contract to it.
func createObscuroNetwork() (func(), error) {
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
	_, err := obscuroNetwork.Create(&simParams, simStats)
	if err != nil {
		return obscuroNetwork.TearDown, fmt.Errorf("failed to create test Obscuro network. Cause: %w", err)
	}

	// Create a wallet extension to allow the creation of the ERC20 contracts.
	walletExtension := walletextension.NewWalletExtension(*walletExtensionConfig)
	defer walletExtension.Shutdown()
	go walletExtension.Serve(network.Localhost, walletExtensionPort, walletExtensionPortWS)
	err = test.WaitForEndpoint(walletExtensionAddrHTTP)
	if err != nil {
		return obscuroNetwork.TearDown, fmt.Errorf("failed to create test Obscuro network. Cause: %w", err)
	}

	// Set up the ERC20 wallet.
	erc20Wallet := wallets.Tokens[bridge.HOC].L2Owner
	_, err = test.GenerateAndSubmitViewingKey(walletExtensionAddrHTTP, erc20Wallet.Address().Hex(), erc20Wallet.PrivateKey())
	if err != nil {
		return obscuroNetwork.TearDown, fmt.Errorf("failed to create test Obscuro network. Cause: %w", err)
	}
	err = fundAccount(erc20Wallet.Address())
	if err != nil {
		return obscuroNetwork.TearDown, fmt.Errorf("failed to create test Obscuro network. Cause: %w", err)
	}

	_, err = sendTransactionAndAwaitConfirmation(erc20Wallet, deployERC20Tx)
	return obscuroNetwork.TearDown, err
}

// Submits a transaction and awaits the transaction receipt.
func sendTransactionAndAwaitConfirmation(txWallet wallet.Wallet, tx types.LegacyTx) (map[string]interface{}, error) {
	// Set the transaction's nonce.
	nonceJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCGetTransactionCount, []interface{}{txWallet.Address().Hex(), latestBlock})
	nonceString, ok := nonceJSON[walletextension.RespJSONKeyResult].(string)
	if !ok {
		respJSON, err := json.Marshal(nonceJSON)
		if err != nil {
			respJSON = []byte(fmt.Sprintf("can't read response as json, cause: %s response: %v", err, nonceJSON))
		}
		return nil, fmt.Errorf("retrieved nonce was not of type string, resp: %s", respJSON)
	}
	nonce, err := hexutil.DecodeUint64(nonceString)
	if err != nil {
		return nil, fmt.Errorf("could not parse nonce from string. Cause: %w", err)
	}
	tx.Nonce = nonce

	// Send the transaction.
	txBinaryHex := signAndSerialiseTransaction(txWallet, &tx)
	sendTxJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCSendRawTransaction, []interface{}{txBinaryHex})

	// Verify the transaction was successful.
	txHash, ok := sendTxJSON[walletextension.RespJSONKeyResult].(string)
	if !ok {
		return nil, fmt.Errorf("could not retrieve transaction hash from JSON result, failed to deploy ERC20")
	}

	counter := 0
	for {
		if counter > 10 {
			return nil, fmt.Errorf("could not get ERC20 receipt after 10 seconds")
		}
		getReceiptJSON := makeHTTPEthJSONReqAsJSON(rpc.RPCGetTxReceipt, []interface{}{txHash})
		getReceiptJSONResult, ok := getReceiptJSON[walletextension.RespJSONKeyResult].(map[string]interface{})
		if ok && getReceiptJSONResult[respJSONKeyStatus] == statusSuccess {
			return getReceiptJSON, nil
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
	txBinaryHex := "0x" + gethcommon.Bytes2Hex(txBinary)
	if err != nil {
		panic(err)
	}

	return txBinaryHex
}

// Funds the account from the faucet account.
func fundAccount(dest gethcommon.Address) error {
	// We create the faucet wallet.
	faucetPrivKey, err := crypto.HexToECDSA(rollupchain.FaucetPrivateKeyHex)
	if err != nil {
		return fmt.Errorf("could not initialise faucet private key")
	}
	faucetWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), faucetPrivKey)

	// We generate a viewing key for the faucet.
	_, err = test.GenerateAndSubmitViewingKey(walletExtensionAddrHTTP, faucetWallet.Address().Hex(), faucetPrivKey)
	if err != nil {
		return err
	}

	// We submit the transaction and await confirmation.
	tx := types.LegacyTx{
		Value:    big.NewInt(faucetAlloc),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		To:       &dest,
	}
	_, err = sendTransactionAndAwaitConfirmation(faucetWallet, tx)
	return err
}

// Causes an event, to allow us to test subscriptions.
// TODO - #453 - Introduce a simpler way to cause an event.
func triggerEvent(t *testing.T) map[string]interface{} {
	// We cause an event by deploying an ERC20 contract.
	_, privateKey, _ := test.RegisterPrivateKey(t, walletExtensionAddrHTTP)
	txWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.ObscuroChainID), privateKey)
	err := fundAccount(txWallet.Address())
	if err != nil {
		t.Fatal(err)
	}
	_, err = txWallet.SignTransaction(&deployERC20Tx)
	if err != nil {
		panic(fmt.Errorf("could not sign transaction. Cause: %w", err))
	}
	receipt, err := sendTransactionAndAwaitConfirmation(txWallet, deployERC20Tx)
	if err != nil {
		t.Fatal(err)
	}
	return receipt
}

// A structure that JSON-serialises to the expected format for subscription filter criteria.
type filterCriteriaJSON struct {
	BlockHash *gethcommon.Hash     `json:"blockHash"`
	FromBlock *gethrpc.BlockNumber `json:"fromBlock"`
	ToBlock   *gethrpc.BlockNumber `json:"toBlock"`
	Addresses interface{}          `json:"address"`
	Topics    []interface{}        `json:"topics"`
}
