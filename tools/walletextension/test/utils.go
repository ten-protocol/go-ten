package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
