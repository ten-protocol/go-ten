package walletextension

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	walletExtensionAddr = "http://localhost:3000"
	chainID             = "0x539"                // Chain ID in hex.
	alloc               = "0x3635c9adc5dea00000" // Default account allocation in hex.
)

func TestCanMakeNonSensitiveRequestWithoutSubmittingViewingKey(t *testing.T) {
	runConfig := RunConfig{LocalNetwork: true}
	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	respBody := makeEthJSONReq(t, "eth_chainId", []string{})

	var respJSON map[string]string
	err := json.Unmarshal(respBody, &respJSON)
	if err != nil {
		t.Fatal(err)
	}

	if respJSON["result"] != chainID {
		t.Fatalf("Expected chainId of %s, got %s", "1337", respJSON["result"])
	}
}

func TestCannotRetrieveBalanceWithoutSubmittingViewingKey(t *testing.T) {
	account := "0x8D97689C9818892B700e27F316cc3E41e17fBeb9"

	runConfig := RunConfig{LocalNetwork: true}
	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	respBody := makeEthJSONReq(t, "eth_getBalance", []string{account, "latest"})

	trimmedRespBody := strings.TrimSpace(string(respBody))
	errPrefix := "enclave could not respond securely to eth_getBalance request because there is no viewing key for the account"
	if trimmedRespBody != errPrefix {
		t.Fatalf("Expected error message with prefix \"%s\", got \"%s\"", errPrefix, trimmedRespBody)
	}
}

func TestCanRetrieveBalanceAfterSubmittingViewingKey(t *testing.T) {
	account := "0x8D97689C9818892B700e27F316cc3E41e17fBeb9"

	runConfig := RunConfig{LocalNetwork: true, PrefundedAccounts: []string{account}}
	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	resp, err := http.Get(walletExtensionAddr + pathGenerateViewingKey) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	submitViewingKeyBodyBytes, _ := json.Marshal(map[string]string{"signedBytes": "dummySignedBytes"})
	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	resp, err = http.Post(walletExtensionAddr+pathSubmitViewingKey, "application/json", submitViewingKeyBody) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	getBalanceRespBody := makeEthJSONReq(t, "eth_getBalance", []string{account, "latest"})

	var getBalanceRespJSON map[string]interface{}
	err = json.Unmarshal(getBalanceRespBody, &getBalanceRespJSON)
	if err != nil {
		t.Fatal(err)
	}

	if getBalanceRespJSON["result"] != alloc {
		t.Fatalf("Expected balance of %s, got %s", alloc, getBalanceRespJSON["result"])
	}
}

// Makes an Ethereum JSON RPC request and returns the response body.
func makeEthJSONReq(t *testing.T, method string, params []string) []byte {
	reqBodyBytes, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "1",
	})
	reqBody := bytes.NewBuffer(reqBodyBytes)
	resp, err := http.Post(walletExtensionAddr, "text/html", reqBody) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return respBody
}
