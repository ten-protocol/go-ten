package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/node"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

var (
	localhost      = "127.0.0.1"
	walExtPortHTTP = integration.StartPortWalletExtensionUnitTest
	walExtPortWS   = integration.StartPortWalletExtensionUnitTest + 1
	nodePortWS     = integration.StartPortWalletExtensionUnitTest + 2
	walExtAddr     = fmt.Sprintf("http://%s:%d", localhost, walExtPortHTTP)
	walExtAddrWS   = fmt.Sprintf("ws://%s:%d", localhost, walExtPortWS)
)

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	err := createWalExt(t)
	if err != nil {
		t.Fatalf(fmt.Sprintf("could not create wallet extension. Cause: %s", err.Error()))
	}

	for _, method := range rpc.SensitiveMethods {
		// We use a websocket request because one of the sensitive methods, eth_subscribe, requires it.
		respBody, _ := MakeWSEthJSONReq(walExtAddrWS, method, []interface{}{})

		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey, method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
}

func TestCannotSubscribeOverHTTP(t *testing.T) {
	err := createWalExt(t)
	if err != nil {
		t.Fatalf("could not create wallet extension")
	}

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filters.FilterCriteria{}})
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

func createWalExt(t *testing.T) error {
	err := createDummyHost(t)
	if err != nil {
		return err
	}

	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("could not create persistence file for wallet extension tests")
	}
	cfg := walletextension.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", nodePortWS),
		PersistencePathOverride: testPersistencePath.Name(),
	}

	walExt := walletextension.NewWalletExtension(cfg)
	t.Cleanup(walExt.Shutdown)
	go walExt.Serve(localhost, int(walExtPortHTTP), int(walExtPortWS))

	err = WaitForEndpoint(walExtAddr + walletextension.PathReady)
	if err != nil {
		return err
	}

	return nil
}

// Creates an RPC layer that the wallet extension can connect to. Returns a handle to shut down the host.
func createDummyHost(t *testing.T) error {
	cfg := node.Config{
		WSHost:    localhost,
		WSPort:    int(nodePortWS),
		WSOrigins: []string{"*"},
	}
	rpcServerNode, err := node.New(&cfg)
	if err != nil {
		return fmt.Errorf("could not create new client server. Cause: %w", err)
	}
	t.Cleanup(func() { rpcServerNode.Close() })

	err = rpcServerNode.Start()
	if err != nil {
		return fmt.Errorf("could not create new client server. Cause: %w", err)
	}

	return nil
}
