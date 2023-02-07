package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"
	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	errFailedDecrypt = "could not decrypt bytes with viewing key"
	dummyParams      = "dummyParams"
	magicNumber      = 123789
	jsonKeyTopics    = "topics"
	_hostWSPort      = integration.StartPortWalletExtensionUnitTest
	_testOffset      = 100
)

var dummyHash = gethcommon.BigToHash(big.NewInt(magicNumber))

type testHelper struct {
	hostPort       int
	walletHTTPPort int
	walletWSPort   int
}

func TestWalletExtension(t *testing.T) {
	createDummyHost(t, _hostWSPort)
	createWalExt(t, createWalExtCfg(_hostWSPort, _hostWSPort+1, _hostWSPort+2))

	h := &testHelper{
		hostPort:       _hostWSPort,
		walletHTTPPort: _hostWSPort + 1,
		walletWSPort:   _hostWSPort + 2,
	}

	for name, test := range map[string]func(t *testing.T, testHelper *testHelper){
		"canInvokeNonSensitiveMethodsWithoutViewingKey":               canInvokeNonSensitiveMethodsWithoutViewingKey,
		"canInvokeSensitiveMethodsWithViewingKey":                     canInvokeSensitiveMethodsWithViewingKey,
		"cannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount": cannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount,
		"canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys": canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys,
		"cannotSubscribeOverHTTP":                                     cannotSubscribeOverHTTP,
		"canRegisterViewingKeyAndMakeRequestsOverWebsockets":          canRegisterViewingKeyAndMakeRequestsOverWebsockets,
	} {
		t.Run(name, func(t *testing.T) {
			test(t, h)
		})
	}
}

func canInvokeNonSensitiveMethodsWithoutViewingKey(t *testing.T, testHelper *testHelper) {
	respBody, _ := makeWSEthJSONReq(testHelper.hostPort, rpc.ChainID, []interface{}{})
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), l2ChainIDHex) {
		t.Fatalf("expected response containing '%s', got '%s'", l2ChainIDHex, string(respBody))
	}
}

func canInvokeSensitiveMethodsWithViewingKey(t *testing.T, testHelper *testHelper) {
	_, viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.Subscribe {
			continue
		}

		respBody := makeHTTPEthJSONReq(testHelper.walletHTTPPort, method, []interface{}{map[string]interface{}{"params": dummyParams}})
		validateJSONResponse(t, respBody)

		if !strings.Contains(string(respBody), dummyParams) {
			t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
		}
	}
}

func cannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount(t *testing.T, testHelper *testHelper) {
	registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)

	// We set the API to decrypt with a key different to the viewing key we just submitted.
	arbitraryPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(fmt.Sprintf("failed to generate private key. Cause: %s", err))
	}
	arbitraryPublicKeyBytesHex := hex.EncodeToString(crypto.CompressPubkey(&arbitraryPrivateKey.PublicKey))
	dummyAPI.setViewingKey([]byte(arbitraryPublicKeyBytesHex))

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.Subscribe {
			continue
		}

		respBody := makeHTTPEthJSONReq(testHelper.walletHTTPPort, method, []interface{}{map[string]interface{}{}})
		if !strings.Contains(string(respBody), errFailedDecrypt) {
			t.Fatalf("expected response containing '%s', got '%s'", errFailedDecrypt, string(respBody))
		}
	}
}

func canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys(t *testing.T, testHelper *testHelper) {
	// We submit viewing keys for ten arbitrary accounts.
	var viewingKeys [][]byte
	for i := 0; i < 10; i++ {
		_, viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
		viewingKeys = append(viewingKeys, viewingKeyBytes)
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryViewingKey := viewingKeys[len(viewingKeys)/2]
	dummyAPI.setViewingKey(arbitraryViewingKey)

	respBody := makeHTTPEthJSONReq(testHelper.walletHTTPPort, rpc.GetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func cannotSubscribeOverHTTP(t *testing.T, testHelper *testHelper) {
	respBody := makeHTTPEthJSONReq(testHelper.walletHTTPPort, rpc.Subscribe, []interface{}{rpc.SubscriptionTypeLogs})
	fmt.Println(respBody)
	if string(respBody) != walletextension.ErrSubscribeFailHTTP+"\n" {
		t.Fatalf("expected response of '%s', got '%s'", walletextension.ErrSubscribeFailHTTP, string(respBody))
	}
}

func canRegisterViewingKeyAndMakeRequestsOverWebsockets(t *testing.T, testHelper *testHelper) {
	_, viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, true)
	dummyAPI.setViewingKey(viewingKeyBytes)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.Subscribe {
			continue
		}

		respBody, _ := makeWSEthJSONReq(testHelper.walletWSPort, method, []interface{}{map[string]interface{}{"params": dummyParams}})
		validateJSONResponse(t, respBody)

		if !strings.Contains(string(respBody), dummyParams) {
			t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
		}

		return // We only need to test a single sensitive method.
	}
}

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	hostPort := _hostWSPort + _testOffset
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	createDummyHost(t, hostPort)
	createWalExt(t, createWalExtCfg(hostPort, walletHTTPPort, walletWSPort))

	for _, method := range rpc.SensitiveMethods {
		// We use a websocket request because one of the sensitive methods, eth_subscribe, requires it.
		respBody, _ := makeWSEthJSONReq(walletWSPort, method, []interface{}{})
		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey, method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
}

func TestKeysAreReloadedWhenWalletExtensionRestarts(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*2
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	createDummyHost(t, hostPort)
	walExtCfg := createWalExtCfg(hostPort, walletHTTPPort, walletWSPort)
	shutdown := createWalExt(t, walExtCfg)

	_, viewingKeyBytes := registerPrivateKey(t, walletHTTPPort, walletWSPort, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	// We shut down the wallet extension and restart it with the same config, forcing the viewing keys to be reloaded.
	shutdown()
	createWalExt(t, walExtCfg)

	respBody := makeHTTPEthJSONReq(walletHTTPPort, rpc.GetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCanSubscribeForLogsOverWebsockets(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*3
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	createDummyHost(t, hostPort)
	createWalExt(t, createWalExtCfg(hostPort, walletHTTPPort, walletWSPort))

	_, viewingKeyBytes := registerPrivateKey(t, walletHTTPPort, walletWSPort, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	filter := common.FilterCriteriaJSON{Topics: []interface{}{dummyHash}}
	resp, conn := makeWSEthJSONReq(walletWSPort, rpc.Subscribe, []interface{}{rpc.SubscriptionTypeLogs, filter})
	validateSubscriptionResponse(t, resp)

	logsJSON := readMessagesForDuration(t, conn, time.Second)

	// We check we received enough logs.
	if len(logsJSON) < 50 {
		t.Errorf("expected to receive at least 50 logs, only received %d", len(logsJSON))
	}

	// We check that none of the logs were duplicates (i.e. were sent twice).
	assertNoDupeLogs(t, logsJSON)

	// We validate that each log contains the correct topic.
	for _, logJSON := range logsJSON {
		var logResp map[string]interface{}
		err := json.Unmarshal(logJSON, &logResp)
		if err != nil {
			t.Fatalf("could not unmarshal received log from JSON")
		}

		// We extract the topic from the received logs. The API should have set this based on the filter we passed when subscribing.
		logMap := logResp[wecommon.JSONKeyParams].(map[string]interface{})[wecommon.JSONKeyResult].(map[string]interface{})
		firstLogTopic := logMap[jsonKeyTopics].([]interface{})[0].(string)

		if firstLogTopic != dummyHash.Hex() {
			t.Errorf("expected first topic to be '%s', got '%s'", dummyHash.Hex(), firstLogTopic)
		}
	}
}
