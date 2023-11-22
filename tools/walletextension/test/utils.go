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

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/gorilla/websocket"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/config"
	"github.com/ten-protocol/go-ten/tools/walletextension/container"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethnode "github.com/ethereum/go-ethereum/node"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

const jsonID = "1"

func createWalExtCfg(connectPort, wallHTTPPort, wallWSPort int) *config.Config { //nolint: unparam
	testDBPath, err := os.CreateTemp("", "")
	if err != nil {
		panic("could not create persistence file for wallet extension tests")
	}
	return &config.Config{
		NodeRPCWebsocketAddress: fmt.Sprintf("localhost:%d", connectPort),
		DBPathOverride:          testDBPath.Name(),
		WalletExtensionPortHTTP: wallHTTPPort,
		WalletExtensionPortWS:   wallWSPort,
		DBType:                  "sqlite",
	}
}

func createWalExt(t *testing.T, walExtCfg *config.Config) func() error {
	// todo (@ziga) - log somewhere else?
	logger := log.New(log.WalletExtCmp, int(gethlog.LvlInfo), log.SysOut)

	wallExtContainer := container.NewWalletExtensionContainerFromConfig(*walExtCfg, logger)
	go wallExtContainer.Start() //nolint: errcheck

	err := waitForEndpoint(fmt.Sprintf("http://%s:%d%s", walExtCfg.WalletExtensionHost, walExtCfg.WalletExtensionPortHTTP, common.PathReady))
	if err != nil {
		t.Fatalf(err.Error())
	}

	return wallExtContainer.Stop
}

// Creates an RPC layer that the wallet extension can connect to. Returns a handle to shut down the host.
func createDummyHost(t *testing.T, wsRPCPort int) (*DummyAPI, func() error) { //nolint: unparam
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
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d/v1/", common.Localhost, port), reqBody)
}

// Makes an Ethereum JSON RPC request over HTTP to specific endpoint and returns the response body.
func makeHTTPEthJSONReqWithPath(port int, path string) []byte {
	reqBody := prepareRequestBody("", "")
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d/%s", common.Localhost, port, path), reqBody)
}

// Makes an Ethereum JSON RPC request over HTTP and returns the response body with userID query paremeter.
func makeHTTPEthJSONReqWithUserID(port int, method string, params interface{}, userID string) []byte { //nolint: unparam
	reqBody := prepareRequestBody(method, params)
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d/v1/?u=%s", common.Localhost, port, userID), reqBody)
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
func simulateViewingKeyRegister(t *testing.T, walletHTTPPort, walletWSPort int, useWS bool) (*gethcommon.Address, []byte, []byte) {
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	accountAddr := crypto.PubkeyToAddress(accountPrivateKey.PublicKey)

	compressedHexVKBytes := generateViewingKey(walletHTTPPort, walletWSPort, accountAddr.String(), useWS)
	mmSignature := signViewingKey(accountPrivateKey, compressedHexVKBytes)
	submitViewingKey(accountAddr.String(), walletHTTPPort, walletWSPort, mmSignature, useWS)

	// transform the metamask signature to the geth compatible one
	sigStr := hex.EncodeToString(mmSignature)
	// and then we extract the signature bytes in the same way as the wallet extension
	outputSig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		panic(fmt.Errorf("failed to decode signature string: %w", err))
	}
	// This same change is made in geth internals, for legacy reasons to be able to recover the address:
	//	https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	outputSig[64] -= 27

	// keys are expected to be a []byte of hex string
	vkPubKeyBytes, err := hex.DecodeString(string(compressedHexVKBytes))
	if err != nil {
		panic(fmt.Errorf("unexpected hex string"))
	}

	return &accountAddr, vkPubKeyBytes, outputSig
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
		viewingKeyBytes, _ := makeRequestWS(fmt.Sprintf("ws://%s:%d%s", common.Localhost, wallWSPort, common.PathGenerateViewingKey), generateViewingKeyBodyBytes)
		return viewingKeyBytes
	}
	return makeRequestHTTP(fmt.Sprintf("http://%s:%d%s", common.Localhost, wallHTTPPort, common.PathGenerateViewingKey), generateViewingKeyBodyBytes)
}

// Signs a viewing key like metamask
func signViewingKey(privateKey *ecdsa.PrivateKey, compressedHexVKBytes []byte) []byte {
	// compressedHexVKBytes already has the key in the hex format
	// it should be decoded back into raw bytes
	viewingKey, err := hex.DecodeString(string(compressedHexVKBytes))
	if err != nil {
		panic(err)
	}
	msgToSign := viewingkey.GenerateSignMessage(viewingKey)
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
		makeRequestWS(fmt.Sprintf("ws://%s:%d%s", common.Localhost, wallWSPort, common.PathSubmitViewingKey), submitViewingKeyBodyBytes)
	} else {
		makeRequestHTTP(fmt.Sprintf("http://%s:%d%s", common.Localhost, wallHTTPPort, common.PathSubmitViewingKey), submitViewingKeyBodyBytes)
	}
}

// Sends the body to the URL over HTTP, and returns the result.
func makeRequestHTTP(url string, body []byte) []byte {
	generateViewingKeyBody := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", generateViewingKeyBody) //nolint:noctx,gosec
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp == nil || resp.Body == nil {
		return nil
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
//func readMessagesForDuration(t *testing.T, conn *websocket.Conn, duration time.Duration) [][]byte {
//	// We set a timeout to kill the test, in case we never receive a log.
//	timeout := time.AfterFunc(duration*3, func() {
//		t.Fatalf("timed out waiting to receive a log via the subscription")
//	})
//	defer timeout.Stop()
//
//	var msgs [][]byte
//	endTime := time.Now().Add(duration)
//	for {
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			t.Fatalf("could not read message from websocket. Cause: %s", err)
//		}
//		msgs = append(msgs, msg)
//		if time.Now().After(endTime) {
//			return msgs
//		}
//	}
//}

// Asserts that there are no duplicate logs in the provided list.
//func assertNoDupeLogs(t *testing.T, logsJSON [][]byte) {
//	logCount := make(map[string]int)
//
//	for _, logJSON := range logsJSON {
//		// Check if the log is already in the logCount map.
//		_, exist := logCount[string(logJSON)]
//		if exist {
//			logCount[string(logJSON)]++ // If it is, increase the count for that log by one.
//		} else {
//			logCount[string(logJSON)] = 1 // Otherwise, start a count for that log starting at one.
//		}
//	}
//
//	for logJSON, count := range logCount {
//		if count > 1 {
//			t.Errorf("received duplicate log with body %s", logJSON)
//		}
//	}
//}

// Checks that the response to a request is correctly formatted, and returns the result field.
func validateJSONResponse(t *testing.T, resp []byte) {
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
}

// Checks that the response to a subscription request is correctly formatted.
//func validateSubscriptionResponse(t *testing.T, resp []byte) {
//	result := validateJSONResponse(t, resp)
//	pattern := "0x.*"
//	resultString, ok := result.(string)
//	if !ok || !regexp.MustCompile(pattern).MatchString(resultString) {
//		t.Fatalf("subscription response did not contain expected result. Expected pattern matching %s, got %s", pattern, resultString)
//	}
//}
