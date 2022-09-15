package test

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/host/node"

	gethnode "github.com/ethereum/go-ethereum/node"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	localhost        = "127.0.0.1"
	errFailedDecrypt = "failed to decrypt result with viewing key"
)

var (
	walExtPortHTTP = integration.StartPortWalletExtensionUnitTest
	walExtPortWS   = integration.StartPortWalletExtensionUnitTest + 1
	nodePortWS     = integration.StartPortWalletExtensionUnitTest + 2
	walExtAddr     = fmt.Sprintf("http://%s:%d", localhost, walExtPortHTTP)
	walExtAddrWS   = fmt.Sprintf("ws://%s:%d", localhost, walExtPortWS)
	walExtCfg      = createWalExtCfg()
	dummyEthAPI    = &DummyEthAPI{}
)

func TestCanInvokeNonSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t)

	respBody, _ := MakeWSEthJSONReq(walExtAddrWS, rpc.RPCChainID, []interface{}{})

	if !strings.Contains(string(respBody), l2ChainIDHex) {
		t.Fatalf("expected response containing '%s', got '%s'", l2ChainIDHex, string(respBody))
	}
}

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t)

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
	createWalExt(t)

	// We register a viewing key and pass it to the API, so that the RPC layer can properly encrypt responses.
	_, _, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	err := dummyEthAPI.setViewingKey(viewingKeyBytes)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody := MakeHTTPEthJSONReq(walExtAddr, method, []interface{}{map[string]interface{}{}})

		if !strings.Contains(string(respBody), successMsg) {
			t.Fatalf("expected response containing '%s', got '%s'", successMsg, string(respBody))
		}
	}
}

func TestCannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount(t *testing.T) {
	createDummyHost(t)
	createWalExt(t)

	RegisterPrivateKey(t, walExtAddr)

	// We set the API to decrypt with a key different to the viewing key we just submitted.
	arbitraryPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(fmt.Sprintf("failed to generate private key. Cause: %s", err))
	}
	arbitraryPublicKeyBytesHex := hex.EncodeToString(crypto.CompressPubkey(&arbitraryPrivateKey.PublicKey))
	err = dummyEthAPI.setViewingKey([]byte(arbitraryPublicKeyBytesHex))
	if err != nil {
		t.Fatalf(err.Error())
	}

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
	createWalExt(t)

	// todo - joel - update naming below
	// todo - joel - get rid of equivalent integration test

	// We submit viewing keys for ten arbitrary accounts.
	var viewingKeys [][]byte
	for i := 0; i < 10; i++ {
		_, _, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
		viewingKeys = append(viewingKeys, viewingKeyBytes)
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryPrivateKey := viewingKeys[len(viewingKeys)/2]
	err := dummyEthAPI.setViewingKey(arbitraryPrivateKey)
	if err != nil {
		t.Fatalf(err.Error())
	}

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCGetBalance, []interface{}{map[string]interface{}{}})

	if !strings.Contains(string(respBody), successMsg) {
		t.Fatalf("expected response containing '%s', got '%s'", successMsg, string(respBody))
	}
}

func TestKeysAreReloadedWhenWalletExtensionRestarts(t *testing.T) {
	createDummyHost(t)
	shutdown := createWalExt(t)

	// We register a viewing key and pass it to the API, so that the RPC layer can properly encrypt responses.
	_, _, viewingKeyBytes := RegisterPrivateKey(t, walExtAddr)
	err := dummyEthAPI.setViewingKey(viewingKeyBytes)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// We shut down the wallet extension and restart it, forcing the viewing keys to be reloaded.
	shutdown()
	createWalExt(t)

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCGetBalance, []interface{}{map[string]interface{}{}})

	if !strings.Contains(string(respBody), successMsg) {
		t.Fatalf("expected response containing '%s', got '%s'", successMsg, string(respBody))
	}
}

func TestCannotSubscribeOverHTTP(t *testing.T) {
	createDummyHost(t)
	createWalExt(t)

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filters.FilterCriteria{}})
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

func createWalExt(t *testing.T) func() {
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
			Service:   &DummyObscuroAPI{},
			Public:    true,
		},
		{
			Namespace: node.APINamespaceEth,
			Version:   node.APIVersion1,
			Service:   dummyEthAPI,
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
