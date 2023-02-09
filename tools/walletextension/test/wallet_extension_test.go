package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"

	gethcommon "github.com/ethereum/go-ethereum/common"
	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"
)

const (
	errFailedDecrypt = "could not decrypt bytes with viewing key"
	dummyParams      = "dummyParams"
	jsonKeyTopics    = "topics"
	_hostWSPort      = integration.StartPortWalletExtensionUnitTest
	_testOffset      = 100 // offset each test by a multiplier of the offset to avoid port colision. ie: 	hostPort := _hostWSPort + _testOffset*2
)

type testHelper struct {
	hostPort       int
	walletHTTPPort int
	walletWSPort   int
	hostAPI        *DummyAPI
}

func TestWalletExtension(t *testing.T) {
	i := 0
	for name, test := range map[string]func(t *testing.T, testHelper *testHelper){
		"canInvokeNonSensitiveMethodsWithoutViewingKey":               canInvokeNonSensitiveMethodsWithoutViewingKey,
		"canInvokeSensitiveMethodsWithViewingKey":                     canInvokeSensitiveMethodsWithViewingKey,
		"cannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount": cannotInvokeSensitiveMethodsWithViewingKeyForAnotherAccount,
		"canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys": canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys,
		"cannotSubscribeOverHTTP":                                     cannotSubscribeOverHTTP,
		"canRegisterViewingKeyAndMakeRequestsOverWebsockets":          canRegisterViewingKeyAndMakeRequestsOverWebsockets,
	} {
		t.Run(name, func(t *testing.T) {
			hostPort := _hostWSPort + i*_testOffset
			dummyAPI, shutDownHost := createDummyHost(t, hostPort)
			shutdownWallet := createWalExt(t, createWalExtCfg(hostPort, hostPort+1, hostPort+2))

			h := &testHelper{
				hostPort:       hostPort,
				walletHTTPPort: hostPort + 1,
				walletWSPort:   hostPort + 2,
				hostAPI:        dummyAPI,
			}

			test(t, h)

			shutdownWallet()
			err := shutDownHost()
			if err != nil {
				t.Fatal(err)
			}
		})
		i++
	}
}

func canInvokeNonSensitiveMethodsWithoutViewingKey(t *testing.T, testHelper *testHelper) {
	respBody, wsConnWE := makeWSEthJSONReq(testHelper.hostPort, rpc.ChainID, []interface{}{})
	defer wsConnWE.Close()

	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), l2ChainIDHex) {
		t.Fatalf("expected response containing '%s', got '%s'", l2ChainIDHex, string(respBody))
	}
}

func canInvokeSensitiveMethodsWithViewingKey(t *testing.T, testHelper *testHelper) {
	viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
	testHelper.hostAPI.setViewingKey(viewingKeyBytes)

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
	testHelper.hostAPI.setViewingKey([]byte(arbitraryPublicKeyBytesHex))

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
		viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
		viewingKeys = append(viewingKeys, viewingKeyBytes)
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryViewingKey := viewingKeys[len(viewingKeys)/2]
	testHelper.hostAPI.setViewingKey(arbitraryViewingKey)

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
	viewingKeyBytes := registerPrivateKey(t, testHelper.walletHTTPPort, testHelper.walletWSPort, true)
	testHelper.hostAPI.setViewingKey(viewingKeyBytes)

	conn, err := openWSConn(testHelper.walletWSPort)
	if err != nil {
		t.Fatal(err)
	}

	respBody := makeWSEthJSONReqWithConn(conn, rpc.GetTransactionReceipt, []interface{}{map[string]interface{}{"params": dummyParams}})
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCannotInvokeSensitiveMethodsWithoutViewingKey(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*7
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	_, shutdownHost := createDummyHost(t, hostPort)
	defer shutdownHost() //nolint: errcheck

	shutdownWallet := createWalExt(t, createWalExtCfg(hostPort, walletHTTPPort, walletWSPort))
	defer shutdownWallet()

	conn, err := openWSConn(walletWSPort)
	if err != nil {
		t.Fatal(err)
	}

	for _, method := range rpc.SensitiveMethods {
		// We use a websocket request because one of the sensitive methods, eth_subscribe, requires it.
		respBody := makeWSEthJSONReqWithConn(conn, method, []interface{}{})
		if !strings.Contains(string(respBody), fmt.Sprintf(accountmanager.ErrNoViewingKey, method)) {
			t.Fatalf("expected response containing '%s', got '%s'", fmt.Sprintf(accountmanager.ErrNoViewingKey, method), string(respBody))
		}
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeysAreReloadedWhenWalletExtensionRestarts(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*8
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	dummyAPI, shutdownHost := createDummyHost(t, hostPort)
	defer shutdownHost() //nolint: errcheck
	walExtCfg := createWalExtCfg(hostPort, walletHTTPPort, walletWSPort)
	shutdownWallet := createWalExt(t, walExtCfg)

	viewingKeyBytes := registerPrivateKey(t, walletHTTPPort, walletWSPort, false)
	dummyAPI.setViewingKey(viewingKeyBytes)

	// We shut down the wallet extension and restart it with the same config, forcing the viewing keys to be reloaded.
	shutdownWallet()
	shutdownWallet = createWalExt(t, walExtCfg)
	defer shutdownWallet()

	respBody := makeHTTPEthJSONReq(walletHTTPPort, rpc.GetBalance, []interface{}{map[string]interface{}{"params": dummyParams}})
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), dummyParams) {
		t.Fatalf("expected response containing '%s', got '%s'", dummyParams, string(respBody))
	}
}

func TestCanSubscribeForLogsOverWebsockets(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*9
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	dummyHash := gethcommon.BigToHash(big.NewInt(1234))

	dummyAPI, shutdownHost := createDummyHost(t, hostPort)
	defer shutdownHost() //nolint: errcheck
	shutdownWallet := createWalExt(t, createWalExtCfg(hostPort, walletHTTPPort, walletWSPort))
	defer shutdownWallet()

	viewingKeyBytes := registerPrivateKey(t, walletHTTPPort, walletWSPort, false)
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
