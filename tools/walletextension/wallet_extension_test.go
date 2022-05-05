package walletextension

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Chain ID in hex.
var chainId = "0x539"

// Default account allocation in hex.
var alloc = "0x3635c9adc5dea00000"

func TestCanMakeNonSensitiveRequestWithoutSubmittingViewingKey(t *testing.T) {
	runConfig := RunConfig{LocalNetwork: true}
	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	reqBodyBytes, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_chainId",
		"params":  []string{},
		"id":      "1",
	})
	reqBody := bytes.NewBuffer(reqBodyBytes)

	resp, err := http.Post("http://localhost:3000/", "text/html", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var respJSON map[string]string
	err = json.Unmarshal(respBody, &respJSON)
	if err != nil {
		t.Fatal(err)
	}

	if respJSON["result"] != chainId {
		t.Fatalf("Expected chainId of %s, got %s", "1337", respJSON["result"])
	}
}

func TestCannotRetrieveBalanceWithoutSubmittingViewingKey(t *testing.T) {
	runConfig := RunConfig{LocalNetwork: true}
	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	reqBodyBytes, _ := json.Marshal(map[string]string{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  "",
		"id":      "1",
	})
	reqBody := bytes.NewBuffer(reqBodyBytes)

	resp, err := http.Post("http://localhost:3000/", "text/html", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

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

	_, err := http.Get("http://localhost:3000" + pathGenerateViewingKey)
	if err != nil {
		t.Fatal(err)
	}

	submitViewingKeyBodyBytes, _ := json.Marshal(map[string]string{"signedBytes": "dummySignedBytes"})
	submitViewingKeyBody := bytes.NewBuffer(submitViewingKeyBodyBytes)
	_, err = http.Post("http://localhost:3000"+pathSubmitViewingKey, "application/json", submitViewingKeyBody)
	if err != nil {
		t.Fatal(err)
	}

	getBalanceBodyBytes, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []string{account, "latest"},
		"id":      "1",
	})
	getBalanceBody := bytes.NewBuffer(getBalanceBodyBytes)
	getBalanceResp, err := http.Post("http://localhost:3000/", "text/html", getBalanceBody)
	if err != nil {
		t.Fatal(err)
	}
	getBalanceRespBody, err := ioutil.ReadAll(getBalanceResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var getBalanceRespJSON map[string]interface{}
	err = json.Unmarshal(getBalanceRespBody, &getBalanceRespJSON)
	if err != nil {
		t.Fatal(err)
	}

	if getBalanceRespJSON["result"] != alloc {
		t.Fatalf("Expected balance of %s, got %s", alloc, getBalanceRespJSON["result"])
	}
}
