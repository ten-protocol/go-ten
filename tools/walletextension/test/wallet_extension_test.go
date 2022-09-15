package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

func TestCannotSubscribeOverHTTP(t *testing.T) {
	server := createDummyHost()
	defer server.Shutdown(context.Background()) //nolint:errcheck

	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		panic("could not create persistence file for wallet extension tests")
	}
	cfg := walletextension.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", nodePortWS),
		PersistencePathOverride: testPersistencePath.Name(),
	}

	walExt := walletextension.NewWalletExtension(cfg)
	defer walExt.Shutdown()

	err = WaitForWalletExtension(walExtAddr)
	if err != nil {
		t.Fatal(err)
	}

	go walExt.Serve(localhost, int(walExtPortHTTP), int(walExtPortWS))

	respBody := MakeHTTPEthJSONReq(walExtAddr, rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filters.FilterCriteria{}})
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

// Creates a dummy host that the wallet extension can connect to.
func createDummyHost() *http.Server {
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", localhost, nodePortWS)} //nolint:gosec
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic("could not upgrade websocket connection in request")
		}
	})

	go func() {
		server.ListenAndServe() //nolint:errcheck
	}()

	return server
}
