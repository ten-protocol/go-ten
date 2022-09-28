package test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
func makeHTTPEthJSONReq(walExtAddr string, method string, params interface{}) []byte {
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
func makeWSEthJSONReq(walExtAddr string, method string, params interface{}) ([]byte, *websocket.Conn) {
	conn, resp, err := websocket.DefaultDialer.Dial(walExtAddr, nil)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		panic(fmt.Errorf("received error response from wallet extension: %w", err))
	}

	reqBody := prepareRequestBody(method, params)
	err = conn.WriteMessage(websocket.TextMessage, reqBody.Bytes())
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		panic(fmt.Errorf("received error response when writing to wallet extension websocket: %w", err))
	}

	_, respBody, err := conn.ReadMessage()
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		panic(fmt.Errorf("received error response when reading from wallet extension websocket: %w", err))
	}

	return respBody, conn
}

// Generates a new account and registers it with the node.
func registerPrivateKey(t *testing.T, walExtAddr string) (gethcommon.Address, []byte) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	accountAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	viewingKeyBytes := generateAndSubmitViewingKey(walExtAddr, accountAddr.String(), privateKey)
	return accountAddr, viewingKeyBytes
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

// Generates a signed viewing key and submits it to the wallet extension.
func generateAndSubmitViewingKey(walExtAddr string, accountAddr string, accountPrivateKey *ecdsa.PrivateKey) []byte {
	viewingKeyBytes := generateViewingKey(walExtAddr, accountAddr)
	signature := signViewingKey(accountPrivateKey, viewingKeyBytes)
	return submitViewingKey(walExtAddr, accountAddr, signature, viewingKeyBytes)
}

// Generates a viewing key.
func generateViewingKey(walExtAddr string, accountAddress string) []byte {
	generateViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeyAddress: accountAddress,
	})
	if err != nil {
		panic(err)
	}
	generateViewingKeyBody := bytes.NewBuffer(generateViewingKeyBodyBytes)
	resp, err := http.Post(walExtAddr+walletextension.PathGenerateViewingKey, "application/json", generateViewingKeyBody) //nolint:noctx
	if err != nil {
		panic(err)
	}
	viewingKey, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return viewingKey
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
func submitViewingKey(walExtAddr string, accountAddr string, signature []byte, viewingKeyBytes []byte) []byte {
	submitViewingKeyBodyBytes, err := json.Marshal(map[string]interface{}{
		walletextension.ReqJSONKeySignature: hex.EncodeToString(signature),
		walletextension.ReqJSONKeyAddress:   accountAddr,
	})
	if err != nil {
		panic(err)
	}
	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	resp, err := http.Post(walExtAddr+walletextension.PathSubmitViewingKey, "application/json", submitViewingKeyBody) //nolint:noctx
	if err != nil {
		panic(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, err := io.ReadAll(resp.Body)
		if err == nil {
			panic(fmt.Errorf("request to add viewing key failed with status %s: %s", resp.Status, respBody))
		}
		panic(fmt.Errorf("request to add viewing key failed with status %s", resp.Status))
	}
	if err != nil {
		panic(err)
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return viewingKeyBytes
}
