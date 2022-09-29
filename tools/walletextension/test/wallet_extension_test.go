package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	errFailedDecrypt = "could not decrypt bytes with viewing key"
	dummyParams      = "dummyParams"
	magicNumber      = 123789
	jsonKeyTopics    = "topics"
)

var dummyHash = gethcommon.BigToHash(big.NewInt(magicNumber))

func TestCanInvokeNonSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	respBody, _ := makeWSEthJSONReq(rpc.RPCChainID, []interface{}{})

	if !strings.Contains(string(respBody), l2ChainIDHex) {
		t.Fatalf("expected response containing '%s', got '%s'", l2ChainIDHex, string(respBody))
	}
}

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	for _, method := range rpc.SensitiveMethods {
		// We use a websocket request because one of the sensitive methods, eth_subscribe, requires it.
		respBody, _ := makeWSEthJSONReq(method, []interface{}{})

		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey, method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
}

func TestCanInvokeSensitiveMethodsWithViewingKey(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	_, viewingKeyBytes := registerPrivateKey(t, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody := makeHTTPEthJSONReq(method, []interface{}{map[string]interface{}{"params": dummyParams}})

		if !strings.Contains(string(respBody), dummyParams) {
			t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
		}
	}
}

func TestCannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	registerPrivateKey(t, false)

	// We set the API to decrypt with a key different to the viewing key we just submitted.
	arbitraryPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(fmt.Sprintf("failed to generate private key. Cause: %s", err))
	}
	arbitraryPublicKeyBytesHex := hex.EncodeToString(crypto.CompressPubkey(&arbitraryPrivateKey.PublicKey))
	dummyAPI.setViewingKey([]byte(arbitraryPublicKeyBytesHex))

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody := makeHTTPEthJSONReq(method, []interface{}{map[string]interface{}{}})

		if !strings.Contains(string(respBody), errFailedDecrypt) {
			t.Fatalf("expected response containing '%s', got '%s'", errFailedDecrypt, string(respBody))
		}
	}
}

func TestCanInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	// We submit viewing keys for ten arbitrary accounts.
	var viewingKeys [][]byte
	for i := 0; i < 10; i++ {
		_, viewingKeyBytes := registerPrivateKey(t, false)
		viewingKeys = append(viewingKeys, viewingKeyBytes)
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryViewingKey := viewingKeys[len(viewingKeys)/2]
	dummyAPI.setViewingKey(arbitraryViewingKey)

	respBody := makeHTTPEthJSONReq(rpc.RPCGetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCanCallWithoutSettingFromField(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	vkAddress, viewingKeyBytes := registerPrivateKey(t, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range []string{rpc.RPCCall, rpc.RPCEstimateGas} {
		respBody := makeHTTPEthJSONReq(method, []interface{}{map[string]interface{}{
			"To":    "0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC",
			"data":  "0x70a0823100000000000000000000000013e23ca74de0206c56ebae8d51b5622eff1e9944",
			"value": nil,
			"Value": "",
		}})

		// RPCCall and RPCEstimateGas payload might be manipulated ( added the From field information )
		if !strings.Contains(strings.ToLower(string(respBody)), strings.ToLower(vkAddress.Hex())) {
			t.Fatalf("expected response containing '%s', got '%s'", strings.ToLower(vkAddress.Hex()), string(respBody))
		}
	}
}

func TestKeysAreReloadedWhenWalletExtensionRestarts(t *testing.T) {
	createDummyHost(t)
	walExtCfg := createWalExtCfg()
	shutdown := createWalExt(t, walExtCfg)

	_, viewingKeyBytes := registerPrivateKey(t, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	// We shut down the wallet extension and restart it with the same config, forcing the viewing keys to be reloaded.
	shutdown()
	createWalExt(t, walExtCfg)

	respBody := makeHTTPEthJSONReq(rpc.RPCGetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCannotSubscribeOverHTTP(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	respBody := makeHTTPEthJSONReq(rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs})
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

func TestCanRegisterViewingKeyAndMakeRequestsOverWebsockets(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	_, viewingKeyBytes := registerPrivateKey(t, true)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.RPCSubscribe {
			continue
		}

		respBody, _ := makeWSEthJSONReq(method, []interface{}{map[string]interface{}{"params": dummyParams}})

		if !strings.Contains(string(respBody), dummyParams) {
			t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
		}

		return // We only need to test a single sensitive method.
	}
}

func TestCanSubscribeForLogsOverWebsockets(t *testing.T) {
	createDummyHost(t)
	createWalExt(t, createWalExtCfg())

	_, viewingKeyBytes := registerPrivateKey(t, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	_, conn := makeWSEthJSONReq(rpc.RPCSubscribe, []interface{}{rpc.RPCSubscriptionTypeLogs, filterCriteriaJSON{Topics: []interface{}{dummyHash}}})

	// We set a timeout to kill the test, in case we never receive a log.
	timeout := time.AfterFunc(3*time.Second, func() {
		t.Fatalf("timed out waiting to receive a log via the subscription")
	})
	defer timeout.Stop()

	// We watch the connection to receive a log...
	_, respJSON, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("could not read log from websocket. Cause: %s", err)
	}

	var resp map[string]interface{}
	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		t.Fatalf("could not unmarshal received log from JSON")
	}

	// We extract the topic from the received logs. The API should have set this based on the filter we passed when subscribing.
	logMap := resp[accountmanager.JSONKeyParams].(map[string]interface{})[accountmanager.JSONKeyResult].(map[string]interface{})
	logTopic := logMap[jsonKeyTopics].([]interface{})[0].(string)

	if !strings.Contains(logTopic, dummyHash.Hex()) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyHash.Hex(), logTopic)
	}
}

// A structure that JSON-serialises to the expected format for subscription filter criteria.
type filterCriteriaJSON struct {
	BlockHash *gethcommon.Hash     `json:"blockHash"`
	FromBlock *gethrpc.BlockNumber `json:"fromBlock"`
	ToBlock   *gethrpc.BlockNumber `json:"toBlock"`
	Addresses interface{}          `json:"address"`
	Topics    []interface{}        `json:"topics"`
}
