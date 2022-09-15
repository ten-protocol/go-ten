package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

// WaitForWalletExtension waits for wallet extension to be ready. Times out after three seconds.
func WaitForWalletExtension(walExtAddr string) error {
	retries := 30
	for i := 0; i < retries; i++ {
		resp, err := http.Get(walExtAddr + walletextension.PathReady) //nolint:noctx
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

// MakeHTTPEthJSONReq makes an Ethereum JSON RPC request over HTTP and returns the response body.
func MakeHTTPEthJSONReq(address string, method string, params interface{}) []byte {
	reqBody := PrepareRequestBody(method, params)

	resp, err := http.Post(address, "text/html", reqBody) //nolint:noctx,gosec
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

// PrepareRequestBody formats a method and its parameters as a Ethereum JSON RPC request.
func PrepareRequestBody(method string, params interface{}) *bytes.Buffer {
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
