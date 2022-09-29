package test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	gethnode "github.com/ethereum/go-ethereum/node"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/host/node"
	"github.com/obscuronet/go-obscuro/integration"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	enclaverpc "github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"

	"github.com/gorilla/websocket"
)

const (
	localhost = "127.0.0.1"
)

var (
	walExtPort   = integration.StartPortWalletExtensionUnitTest
	walExtPortWS = integration.StartPortWalletExtensionUnitTest + 1
	walExtAddr   = fmt.Sprintf("http://%s:%d", localhost, walExtPort)
	walExtAddrWS = fmt.Sprintf("ws://%s:%d", localhost, walExtPortWS)
	nodePortWS   = integration.StartPortWalletExtensionUnitTest + 2
	dummyAPI     = NewDummyAPI()
)

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
	go walExt.Serve(localhost, walExtPort, walExtPortWS)

	err := waitForEndpoint(walExtAddr + walletextension.PathReady)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return walExt.Shutdown
}

// Creates an RPC layer that the wallet extension can connect to. Returns a handle to shut down the host.
func createDummyHost(t *testing.T) {
	cfg := gethnode.Config{
		WSHost:    localhost,
		WSPort:    nodePortWS,
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

// Waits for the endpoint to be available. Times out after three seconds.
func waitForEndpoint(addr string) error {
	retries := 30
	for i := 0; i < retries; i++ {
		resp, err := http.Get(addr) //nolint:noctx,gosec
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if err == nil {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("could not establish connection to wallet extension")
}

// Makes an Ethereum JSON RPC request over HTTP and returns the response body.
func makeHTTPEthJSONReq(method string, params interface{}) []byte {
	reqBody := prepareRequestBody(method, params)
	return makeRequestHTTP(walExtAddr, reqBody)
}

// Makes an Ethereum JSON RPC request over websockets and returns the response body.
func makeWSEthJSONReq(method string, params interface{}) ([]byte, *websocket.Conn) {
	reqBody := prepareRequestBody(method, params)
	return makeRequestWS(walExtAddrWS, reqBody)
}

// Formats a method and its parameters as a Ethereum JSON RPC request.
func prepareRequestBody(method string, params interface{}) []byte {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		accountmanager.JSONKeyRPCVersion: jsonrpc.Version,
		accountmanager.JSONKeyMethod:     method,
		accountmanager.JSONKeyParams:     params,
		walletextension.JSONKeyID:        "1",
	})
	if err != nil {
		panic(fmt.Errorf("failed to prepare request body. Cause: %w", err))
	}
	return reqBodyBytes
}

// Generates a new account and registers it with the node.
func registerPrivateKey(t *testing.T, useWS bool) (gethcommon.Address, []byte) {
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	accountAddr := crypto.PubkeyToAddress(accountPrivateKey.PublicKey)

	viewingKeyBytes := generateViewingKey(accountAddr.String(), useWS)
	signature := signViewingKey(accountPrivateKey, viewingKeyBytes)
	submitViewingKey(accountAddr.String(), signature, useWS)

	return accountAddr, viewingKeyBytes
}

// Generates a viewing key.
func generateViewingKey(accountAddress string, useWS bool) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.JSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}

	if useWS {
		viewingKeyBytes, _ := makeRequestWS(walExtAddrWS+walletextension.PathGenerateViewingKey, generateViewingKeyBodyBytes)
		return viewingKeyBytes
	}
	return makeRequestHTTP(walExtAddr+walletextension.PathGenerateViewingKey, generateViewingKeyBodyBytes)
}

// Signs a viewing key.
func signViewingKey(privateKey *ecdsa.PrivateKey, viewingKey []byte) []byte {
	msgToSign := enclaverpc.ViewingKeySignedMsgPrefix + string(viewingKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), privateKey)
	if err != nil {
		panic(err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	return signatureWithLeadBytes
}

// Submits a viewing key.
func submitViewingKey(accountAddr string, signature []byte, useWS bool) {
	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.JSONKeySignature: hex.EncodeToString(signature),
		walletextension.JSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}

	if useWS {
		makeRequestWS(walExtAddrWS+walletextension.PathSubmitViewingKey, submitViewingKeyBodyBytes)
	} else {
		makeRequestHTTP(walExtAddr+walletextension.PathSubmitViewingKey, submitViewingKeyBodyBytes)
	}
}

// Sends the body to the URL over HTTP, and returns the result.
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

// Sends the body to the URL over a websocket connection, and returns the result.
func makeRequestWS(url string, body []byte) ([]byte, *websocket.Conn) {
	conn, dialResp, err := websocket.DefaultDialer.Dial(url, nil)
	if dialResp != nil && dialResp.Body != nil {
		defer dialResp.Body.Close()
	}
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		panic(fmt.Errorf("received error response from wallet extension: %w", err))
	}

	err = conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		panic(err)
	}

	_, reqResp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	return reqResp, conn
}
