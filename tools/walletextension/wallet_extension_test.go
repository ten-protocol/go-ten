package walletextension

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCannotRetrieveBalanceWithoutSubmittingViewingKey(t *testing.T) {
	runConfig := RunConfig{
		LocalNetwork:      true,
		PrefundedAccounts: []string{"0x41F534DB02c6953FB6d9Bd9Eff8B55C364819700"},
	}

	stopNodesFunc := StartWalletExtension(runConfig)
	defer stopNodesFunc()

	reqBodyBytes, _ := json.Marshal(map[string]string{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  "",
		"id":      "1",
	})
	reqBody := bytes.NewBuffer(reqBodyBytes)

	resp, err := http.Post("http://localhost:3000/", "application/json", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// TODO - Investigate why HTTP error message is being truncated.
	errPrefix := "enclave could not respond securely to eth_getBalance request because there is no viewing key for the"
	if !strings.HasPrefix(string(body), errPrefix) {
		t.Fatalf("Expected error message with prefix \"%s\", got \"%s\"", errPrefix, string(body))
	}
}
