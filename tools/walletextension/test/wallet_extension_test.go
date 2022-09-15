package test

import (
	"context"
	"fmt"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/gorilla/websocket"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

var (
	upgrader       = websocket.Upgrader{}
	localhost      = "127.0.0.1"
	walExtPortHTTP = integration.StartPortWalletExtensionUnitTest
	walExtPortWS   = integration.StartPortWalletExtensionUnitTest + 1
	nodePortWS     = integration.StartPortWalletExtensionUnitTest + 2
	walExtAddr     = fmt.Sprintf("http://%s:%d", localhost, walExtPortHTTP)
)

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	shutdown, err := createWalExt()
	defer shutdown()
	if err != nil {
		t.Fatalf("could not create wallet extension")
	}

	for _, method := range rpc.SensitiveMethods {
		respBody := MakeHTTPEthJSONReq(walExtAddr, method, []interface{}{})
		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey+"\n", method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
}

func TestCannotSubscribeOverHTTP(t *testing.T) {
	shutdown, err := createWalExt()
	defer shutdown()
	if err != nil {
		t.Fatalf("could not create wallet extension")
	}

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filters.FilterCriteria{}})
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

func createWalExt() (func(), error) {
	server, err := createDummyHost()
	if err != nil {
		return nil, err
	}

	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		server.Shutdown(context.Background()) //nolint:errcheck
		return nil, fmt.Errorf("could not create persistence file for wallet extension tests")
	}
	cfg := walletextension.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", nodePortWS),
		PersistencePathOverride: testPersistencePath.Name(),
	}

	walExt := walletextension.NewWalletExtension(cfg)
	go walExt.Serve(localhost, int(walExtPortHTTP), int(walExtPortWS))

	err = WaitForEndpoint(walExtAddr + walletextension.PathReady)
	if err != nil {
		walExt.Shutdown()
		server.Shutdown(context.Background()) //nolint:errcheck
		return nil, err
	}

	return func() {
		server.Shutdown(context.Background()) //nolint:errcheck
		walExt.Shutdown()
	}, nil
}

// Creates a dummy host that the wallet extension can connect to.
func createDummyHost() (*http.Server, error) {
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", localhost, nodePortWS)} //nolint:gosec
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic("could not upgrade websocket connection in request")
		}
	})

	go func() {
		server.ListenAndServe() //nolint:errcheck
	}()

	err := WaitForEndpoint(fmt.Sprintf("http://%s:%d/ready", localhost, nodePortWS))
	if err != nil {
		server.Shutdown(context.Background()) //nolint:errcheck
		return nil, fmt.Errorf("could not retrieve host endpoint after waiting")
	}

	return server, nil
}
