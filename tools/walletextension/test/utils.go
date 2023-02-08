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
	"regexp"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/gorilla/websocket"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethnode "github.com/ethereum/go-ethereum/node"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	enclaverpc "github.com/obscuronet/go-obscuro/go/enclave/rpc"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
)

const jsonID = "1"

func createWalExtCfg(connectPort, wallHTTPPort, wallWSPort int) *walletextension.Config {
	testPersistencePath, err := os.CreateTemp("", "")
	if err != nil {
		panic("could not create persistence file for wallet extension tests")
	}
	return &walletextension.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", connectPort),
		PersistencePathOverride: testPersistencePath.Name(),
		WalletExtensionPort:     wallHTTPPort,
		WalletExtensionPortWS:   wallWSPort,
	}
}

func createWalExt(t *testing.T, walExtCfg *walletextension.Config) func() {
	// todo - log somewhere else?
	logger := log.New(log.WalletExtCmp, int(gethlog.LvlInfo), log.SysOut)

	walExt := walletextension.NewWalletExtension(*walExtCfg, logger)
	go walExt.Serve(common.Localhost, walExtCfg.WalletExtensionPort, walExtCfg.WalletExtensionPortWS)

	err := waitForEndpoint(fmt.Sprintf("http://%s:%d%s", common.Localhost, walExtCfg.WalletExtensionPort, walletextension.PathReady))
	if err != nil {
		t.Fatalf(err.Error())
	}

	return walExt.Shutdown
}

// Creates an RPC layer that the wallet extension can connect to. Returns a handle to shut down the host.
func createDummyHost(t *testing.T, wsRPCPort int) (*DummyAPI, func() error) {
	dummyAPI := NewDummyAPI()
	cfg := gethnode.Config{
		WSHost:    common.Localhost,
		WSPort:    wsRPCPort,
		WSOrigins: []string{"*"},
	}
	rpcServerNode, err := gethnode.New(&cfg)
	rpcServerNode.RegisterAPIs([]gethrpc.API{
		{
			Namespace: hostcontainer.APINamespaceObscuro,
			Version:   hostcontainer.APIVersion1,
			Service:   dummyAPI,
			Public:    true,
		},
		{
			Namespace: hostcontainer.APINamespaceEth,
			Version:   hostcontainer.APIVersion1,
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
	return dummyAPI, rpcServerNode.Close
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
func makeHTTPEthJSONReq(port int, method string, params interface{}) []byte {
	reqBody := prepareRequestBody(method, params)
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d", common.Localhost, port), reqBody)
}

// Makes an Ethereum JSON RPC request over websockets and returns the response body.
func makeWSEthJSONReq(port int, method string, params interface{}) ([]byte, *websocket.Conn) {
	reqBody := prepareRequestBody(method, params)
	return makeRequestWS(fmt.Sprintf("ws://%s:%d", common.Localhost, port), reqBody)
}

func makeWSEthJSONReqWithConn(conn *websocket.Conn, method string, params interface{}) []byte {
	reqBody := prepareRequestBody(method, params)
	return issueRequestWS(conn, reqBody)
}

func openWSConn(port int) (*websocket.Conn, error) {
	conn, dialResp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d", common.Localhost, port), nil)
	if dialResp != nil && dialResp.Body != nil {
		defer dialResp.Body.Close()
	}
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		panic(fmt.Errorf("received error response from wallet extension: %w", err))
	}
	return conn, err
}

// Formats a method and its parameters as a Ethereum JSON RPC request.
func prepareRequestBody(method string, params interface{}) []byte {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		common.JSONKeyRPCVersion: jsonrpc.Version,
		common.JSONKeyMethod:     method,
		common.JSONKeyParams:     params,
		common.JSONKeyID:         jsonID,
	})
	if err != nil {
		panic(fmt.Errorf("failed to prepare request body. Cause: %w", err))
	}
	return reqBodyBytes
}

// Generates a new account and registers it with the node.
func registerPrivateKey(t *testing.T, walletHTTPPort, walletWSPort int, useWS bool) []byte {
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	accountAddr := crypto.PubkeyToAddress(accountPrivateKey.PublicKey)

	viewingKeyBytes := generateViewingKey(walletHTTPPort, walletWSPort, accountAddr.String(), useWS)
	signature := signViewingKey(accountPrivateKey, viewingKeyBytes)
	submitViewingKey(accountAddr.String(), walletHTTPPort, walletWSPort, signature, useWS)

	return viewingKeyBytes
}

// Generates a viewing key.
func generateViewingKey(wallHTTPPort, wallWSPort int, accountAddress string, useWS bool) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		common.JSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}

	if useWS {
		viewingKeyBytes, _ := makeRequestWS(fmt.Sprintf("ws://%s:%d%s", common.Localhost, wallWSPort, walletextension.PathGenerateViewingKey), generateViewingKeyBodyBytes)
		return viewingKeyBytes
	}
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d%s", common.Localhost, wallHTTPPort, walletextension.PathGenerateViewingKey), generateViewingKeyBodyBytes)
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
func submitViewingKey(accountAddr string, wallHTTPPort, wallWSPort int, signature []byte, useWS bool) {
	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		common.JSONKeySignature: hex.EncodeToString(signature),
		common.JSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}

	if useWS {
		makeRequestWS(fmt.Sprintf("ws://%s:%d%s", common.Localhost, wallWSPort, walletextension.PathSubmitViewingKey), submitViewingKeyBodyBytes)
	} else {
		makeRequestHTTP(fmt.Sprintf("http://%s:%d%s", common.Localhost, wallHTTPPort, walletextension.PathSubmitViewingKey), submitViewingKeyBodyBytes)
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

	return issueRequestWS(conn, body), conn
}

// issues request on an existing ws connection
func issueRequestWS(conn *websocket.Conn, body []byte) []byte {
	err := conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		panic(err)
	}

	_, reqResp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	return reqResp
}

// Reads messages from the connection for the provided duration, and returns the read messages.
func readMessagesForDuration(t *testing.T, conn *websocket.Conn, duration time.Duration) [][]byte {
	// We set a timeout to kill the test, in case we never receive a log.
	timeout := time.AfterFunc(duration*3, func() {
		t.Fatalf("timed out waiting to receive a log via the subscription")
	})
	defer timeout.Stop()

	var msgs [][]byte
	endTime := time.Now().Add(duration)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("could not read message from websocket. Cause: %s", err)
		}
		msgs = append(msgs, msg)
		if time.Now().After(endTime) {
			return msgs
		}
	}
}

// Asserts that there are no duplicate logs in the provided list.
func assertNoDupeLogs(t *testing.T, logsJSON [][]byte) {
	logCount := make(map[string]int)

	for _, logJSON := range logsJSON {
		// Check if the log is already in the logCount map.
		_, exist := logCount[string(logJSON)]
		if exist {
			logCount[string(logJSON)]++ // If it is, increase the count for that log by one.
		} else {
			logCount[string(logJSON)] = 1 // Otherwise, start a count for that log starting at one.
		}
	}

	for logJSON, count := range logCount {
		if count > 1 {
			t.Errorf("received duplicate log with body %s", logJSON)
		}
	}
}

// Checks that the response to a request is correctly formatted, and returns the result field.
func validateJSONResponse(t *testing.T, resp []byte) interface{} {
	var respJSON map[string]interface{}
	err := json.Unmarshal(resp, &respJSON)
	if err != nil {
		t.Fatalf("could not unmarshal response to JSON")
	}

	id := respJSON[common.JSONKeyID]
	jsonRPCVersion := respJSON[common.JSONKeyRPCVersion]
	result := respJSON[common.JSONKeyResult]

	if id != jsonID {
		t.Fatalf("response did not contain expected ID. Expected 1, got %s", id)
	}
	if jsonRPCVersion != jsonrpc.Version {
		t.Fatalf("response did not contain expected RPC version. Expected 2.0, got %s", jsonRPCVersion)
	}
	if result == nil {
		t.Fatalf("response did not contain `result` field")
	}

	return result
}

// Checks that the response to a subscription request is correctly formatted.
func validateSubscriptionResponse(t *testing.T, resp []byte) {
	result := validateJSONResponse(t, resp)
	pattern := "0x.*"
	resultString, ok := result.(string)
	if !ok || !regexp.MustCompile(pattern).MatchString(resultString) {
		t.Fatalf("subscription response did not contain expected result. Expected pattern matching %s, got %s", pattern, resultString)
	}
}
