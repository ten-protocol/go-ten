package test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/obscuronet/go-obscuro/integration"
	"io"
	"net/http"
	"testing"
	"time"

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
	walExtPortHTTP = integration.StartPortWalletExtensionUnitTest
	walExtPortWS   = integration.StartPortWalletExtensionUnitTest + 1
	walExtAddr     = fmt.Sprintf("http://%s:%d", localhost, walExtPortHTTP)
	walExtAddrWS   = fmt.Sprintf("ws://%s:%d", localhost, walExtPortWS)
)

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

	resp, err := http.Post(walExtAddr, "text/html", reqBody) //nolint:noctx,gosec
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(fmt.Errorf("received error response from wallet extension: %w", err))
	}
	if resp == nil {
		panic("did not receive a response from the wallet extension")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return respBody
}

// Makes an Ethereum JSON RPC request over websockets and returns the response body.
func makeWSEthJSONReq(method string, params interface{}) ([]byte, *websocket.Conn) {
	reqBody := prepareRequestBody(method, params)
	return makeRequestWS(walExtAddrWS, reqBody.Bytes())
}

// Formats a method and its parameters as a Ethereum JSON RPC request.
func prepareRequestBody(method string, params interface{}) *bytes.Buffer {
	reqBodyBytes, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "1",
	})
	if err != nil {
		panic(fmt.Errorf("failed to prepare request body. Cause: %w", err))
	}
	return bytes.NewBuffer(reqBodyBytes)
}

// Generates a new account and registers it with the node.
func registerPrivateKey(t *testing.T, useWS bool) (gethcommon.Address, []byte) {
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	accountAddr := crypto.PubkeyToAddress(accountPrivateKey.PublicKey)

	var viewingKeyBytes []byte
	if useWS {
		viewingKeyBytes = generateViewingKeyWS(accountAddr.String())
	} else {
		viewingKeyBytes = generateViewingKey(accountAddr.String())
	}

	signature := signViewingKey(accountPrivateKey, viewingKeyBytes)

	if useWS {
		submitViewingKeyWS(accountAddr.String(), signature)
	} else {
		submitViewingKey(accountAddr.String(), signature)
	}

	return accountAddr, viewingKeyBytes
}

// Generates a viewing key.
func generateViewingKey(accountAddress string) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}
	generateViewingKeyBody := bytes.NewBuffer(generateViewingKeyBodyBytes)
	resp, err := http.Post(walExtAddr+walletextension.PathGenerateViewingKey, "application/json", generateViewingKeyBody) //nolint:noctx
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

func generateViewingKeyWS(accountAddress string) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}

	viewingKeyBytes, _ := makeRequestWS(walExtAddrWS+walletextension.PathGenerateViewingKey, generateViewingKeyBodyBytes)
	return viewingKeyBytes
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
func submitViewingKey(accountAddr string, signature []byte) {
	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeySignature: hex.EncodeToString(signature),
		walletextension.ReqJSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}

	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	resp, err := http.Post(walExtAddr+walletextension.PathSubmitViewingKey, "application/json", submitViewingKeyBody) //nolint:noctx
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
}

func submitViewingKeyWS(accountAddr string, signature []byte) {
	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeySignature: hex.EncodeToString(signature),
		walletextension.ReqJSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}

	makeRequestWS(walExtAddrWS+walletextension.PathSubmitViewingKey, submitViewingKeyBodyBytes)
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

	err = conn.WriteMessage(websocket.TextMessage, body) //nolint:noctx
	if err != nil {
		panic(err)
	}

	_, reqResp, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	return reqResp, conn
}
