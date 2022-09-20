package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/host/node"

	gethnode "github.com/ethereum/go-ethereum/node"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	localhost        = "127.0.0.1"
	errFailedDecrypt = "could not decrypt result with viewing key"
	dummyParams      = "dummyParams"
	magicNumber      = 123789
)

var dummyHash = gethcommon.BigToHash(big.NewInt(magicNumber))

var (
	walExtPortHTTP = integration.StartPortWalletExtensionUnitTest
	walExtPortWS   = integration.StartPortWalletExtensionUnitTest + 1
	nodePortWS     = integration.StartPortWalletExtensionUnitTest + 2
	walExtAddr     = fmt.Sprintf("http://%s:%d", localhost, walExtPortHTTP)
	walExtAddrWS   = fmt.Sprintf("ws://%s:%d", localhost, walExtPortWS)
	dummyAPI       = NewDummyAPI()
)

func TestCanInvokeNonSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	respBody, _ := MakeWSEthJSONReq(walExtAddrWS, rpc.RPCChainID, []interface{}{})

	if !strings.Contains(string(respBody), l2ChainIDHex) {
		t.Fatalf("expected response containing '%s', got '%s'", l2ChainIDHex, string(respBody))
	}
}

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	for _, method := range rpc.SensitiveMethods {
		// We use a websocket request because one of the sensitive methods, eth_subscribe, requires it.
		respBody, _ := MakeWSEthJSONReq(walExtAddrWS, method, []interface{}{})

		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey, method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
}

func TestCanInvokeSensitiveMethodsWithViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	_, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody := MakeHTTPEthJSONReq(walExtAddr, method, []interface{}{map[string]interface{}{"params": dummyParams}})

		if !strings.Contains(string(respBody), dummyParams) {
			t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
		}
	}
}

func TestCannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	RegisterPrivateKey(t, walExtAddr)

	// We set the API to decrypt with a key different to the viewing key we just submitted.
	arbitraryPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(fmt.Sprintf("failed to generate private key. Cause: %s", err))
	}
	arbitraryPublicKeyBytesHex := hex.EncodeToString(crypto.CompressPubkey(&arbitraryPrivateKey.PublicKey))
	dummyAPI.setViewingKey([]byte(arbitraryPublicKeyBytesHex))

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody := MakeHTTPEthJSONReq(walExtAddr, method, []interface{}{map[string]interface{}{}})

		if !strings.Contains(string(respBody), errFailedDecrypt) {
			t.Fatalf("expected response containing '%s', got '%s'", errFailedDecrypt, string(respBody))
		}
	}
}

func TestCanInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	// We submit viewing keys for ten arbitrary accounts.
	var viewingKeys [][]byte
	for i := 0; i < 10; i++ {
		_, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
		viewingKeys = append(viewingKeys, viewingKeyBytes)
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryViewingKey := viewingKeys[len(viewingKeys)/2]
	dummyAPI.setViewingKey(arbitraryViewingKey)

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCGetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCanCallWithoutSettingFromField(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	accountAddr, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	dummyAPI.setViewingKey(viewingKeyBytes)

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCCall, []interface{}{map[string]interface{}{}})

	// We check the automatically-set `from` field is present.
	fromJSON := fmt.Sprintf("\"from\":\"%s\"", strings.ToLower(accountAddr.Hex()))
	if !strings.Contains(string(respBody), fromJSON) {
		t.Fatalf("expected response containing '%s', got '%s'", fromJSON, string(respBody))
	}
}

func TestKeysAreReloadedWhenWalletExtensionRestarts(t *testing.T) {
	createDummyHost(t)
	walExtCfg := createWalExtCfg()
	shutdown := createWalExt(t, walExtCfg)

	_, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	dummyAPI.setViewingKey(viewingKeyBytes)

	// We shut down the wallet extension and restart it with the same config, forcing the viewing keys to be reloaded.
	shutdown()
	createWalExt(t, walExtCfg)

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCGetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCanSubscribeForLogs(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	_, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	dummyAPI.setViewingKey(viewingKeyBytes)

	_, conn := MakeWSEthJSONReq(walExtAddrWS, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filterCriteriaJSON{Topics: []interface{}{dummyHash}}})

	// We set a timeout to kill the test, in case we never receive a log.
	timeout := time.AfterFunc(3*time.Second, func() {
		panic("timed out waiting to receive a log via the subscription")
	})
	defer timeout.Stop()

	// We watch the connection to receive a log...
	_, receivedLogJSON, err := conn.ReadMessage()
	if err != nil {
		panic(fmt.Errorf("could not read log from websocket. Cause: %w", err))
	}

	var receivedLog *types.Log
	err = json.Unmarshal(receivedLogJSON, &receivedLog)
	if err != nil {
		t.Fatalf("could not unmarshall received log from JSON")
	}

	if !strings.Contains(string(receivedLog.Data), dummyHash.Hex()) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyHash.Hex(), string(receivedLog.Data))
	}
}

func TestCannotSubscribeOverHTTP(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs})
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

func createWalExtCfg() *walletextension.Config {
	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		panic("could not create persistence file for wallet extension tests")
	}
	return &walletextension.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", nodePortWS),
		PersistencePathOverride: testPersistencePath.Name(),
	}
}

func createWalExt(t *testing.T, walExtCfg *walletextension.Config) func() {
	walExt := walletextension.NewWalletExtension(*walExtCfg)
	t.Cleanup(walExt.Shutdown)
	go walExt.Serve(localhost, int(walExtPortHTTP), int(walExtPortWS))

	err := WaitForEndpoint(walExtAddr + walletextension.PathReady)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return walExt.Shutdown
}

// Creates an RPC layer that the wallet extension can connect to. Returns a handle to shut down the host.
func createDummyHost(t *testing.T) {
	cfg := gethnode.Config{
		WSHost:    localhost,
		WSPort:    int(nodePortWS),
		WSOrigins: []string{"*"},
	}
	rpcServerNode, err := gethnode.New(&cfg)
	rpcServerNode.RegisterAPIs([]gethrpc.API{
		{
			Namespace: node.APINamespaceObscuro,
			Version:   node.APIVersion1,
			Service:   dummyAPI,
			Public:    true,
		},
		{
			Namespace: node.APINamespaceEth,
			Version:   node.APIVersion1,
			Service:   dummyAPI,
			Public:    true,
		},
	})
	if err != nil {
		t.Fatalf(fmt.Sprintf("could not create new client server. Cause: %s", err))
	}
	t.Cleanup(func() { rpcServerNode.Close() })

	err = rpcServerNode.Start()
	if err != nil {
		t.Fatalf(fmt.Sprintf("could not create new client server. Cause: %s", err))
	}
}

// A structure that JSON-serialises to the expected format for subscription filter criteria.
type filterCriteriaJSON struct {
	BlockHash *gethcommon.Hash     `json:"blockHash"`
	FromBlock *gethrpc.BlockNumber `json:"fromBlock"`
	ToBlock   *gethrpc.BlockNumber `json:"toBlock"`
	Addresses interface{}          `json:"address"`
	Topics    []interface{}        `json:"topics"`
}
