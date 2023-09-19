package test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/vkhandler"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/stretchr/testify/assert"

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
		"canInvokeSensitiveMethodsWithViewingKey":                     canInvokeSensitiveMethodsWithViewingKey,
		"canInvokeNonSensitiveMethodsWithoutViewingKey":               canInvokeNonSensitiveMethodsWithoutViewingKey,
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

			assert.NoError(t, shutdownWallet())
			assert.NoError(t, shutDownHost())
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
	address, vkPubKeyBytes, signature := simulateViewingKeyRegister(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
	testHelper.hostAPI.setViewingKey(address, vkPubKeyBytes, signature)

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
	addr1, _, _ := simulateViewingKeyRegister(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)

	_, hexVKPubKeyBytes2, signature2 := simulateViewingKeyRegister(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)

	// We set the API to decrypt with a key different to the viewing key we just submitted.
	testHelper.hostAPI.setViewingKey(addr1, hexVKPubKeyBytes2, signature2)

	for _, method := range rpc.SensitiveMethods {
		// Subscriptions have to be tested separately, as they return results differently.
		if method == rpc.Subscribe {
			continue
		}

		respBody := makeHTTPEthJSONReq(testHelper.walletHTTPPort, method, []interface{}{map[string]interface{}{}})
		if !strings.Contains(string(respBody), vkhandler.ErrInvalidAddressSignature.Error()) {
			t.Fatalf("expected response containing '%s', got '%s'", errFailedDecrypt, string(respBody))
		}
	}
}

func canInvokeSensitiveMethodsAfterSubmittingMultipleViewingKeys(t *testing.T, testHelper *testHelper) {
	type tempVKHolder struct {
		address     *gethcommon.Address
		hexVKPubKey []byte
		signature   []byte
	}
	// We submit viewing keys for ten arbitrary accounts.
	var viewingKeys []tempVKHolder

	for i := 0; i < 10; i++ {
		address, hexVKPubKeyBytes, signature := simulateViewingKeyRegister(t, testHelper.walletHTTPPort, testHelper.walletWSPort, false)
		viewingKeys = append(viewingKeys, tempVKHolder{
			address:     address,
			hexVKPubKey: hexVKPubKeyBytes,
			signature:   signature,
		})
	}

	// We set the API to decrypt with an arbitrary key from the list we just generated.
	arbitraryViewingKey := viewingKeys[len(viewingKeys)/2]
	testHelper.hostAPI.setViewingKey(arbitraryViewingKey.address, arbitraryViewingKey.hexVKPubKey, arbitraryViewingKey.signature)

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
	address, hexVKPubKeyBytes, signature := simulateViewingKeyRegister(t, testHelper.walletHTTPPort, testHelper.walletWSPort, true)
	testHelper.hostAPI.setViewingKey(address, hexVKPubKeyBytes, signature)

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
	defer shutdownWallet() //nolint: errcheck

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

	addr, viewingKeyBytes, signature := simulateViewingKeyRegister(t, walletHTTPPort, walletWSPort, false)
	dummyAPI.setViewingKey(addr, viewingKeyBytes, signature)

	// We shut down the wallet extension and restart it with the same config, forcing the viewing keys to be reloaded.
	err := shutdownWallet()
	assert.NoError(t, err)

	shutdownWallet = createWalExt(t, walExtCfg)
	defer shutdownWallet() //nolint: errcheck

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
	defer shutdownWallet() //nolint: errcheck

	dummyAPI.setViewingKey(simulateViewingKeyRegister(t, walletHTTPPort, walletWSPort, false))

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

func TestGetStorageAtForReturningUserID(t *testing.T) {
	hostPort := _hostWSPort + _testOffset*8
	walletHTTPPort := hostPort + 1
	walletWSPort := hostPort + 2

	createDummyHost(t, hostPort)
	walExtCfg := createWalExtCfg(hostPort, walletHTTPPort, walletWSPort)
	createWalExtCfg(hostPort, walletHTTPPort, walletWSPort)
	createWalExt(t, walExtCfg)

	// create userID
	respJoin := makeHTTPEthJSONReqWithPath(walletHTTPPort, "v1/join")
	userID := string(respJoin)

	// make a request to GetStorageAt with correct parameters to get userID that exists in the database
	respBody := makeHTTPEthJSONReqWithUserID(walletHTTPPort, rpc.GetStorageAt, []interface{}{"getUserID", "0", nil}, userID)
	validateJSONResponse(t, respBody)

	if !strings.Contains(string(respBody), userID) {
		t.Fatalf("expected response containing '%s', got '%s'", userID, string(respBody))
	}

	// make a request to GetStorageAt with correct parameters, but userID that is not present in the database
	invalidUserID := "abc123"
	respBody2 := makeHTTPEthJSONReqWithUserID(walletHTTPPort, rpc.GetStorageAt, []interface{}{"getUserID", "0", nil}, invalidUserID)

	if !strings.Contains(string(respBody2), "method eth_getStorageAt cannot be called with an unauthorised client - no signed viewing keys found") {
		t.Fatalf("expected method eth_getStorageAt cannot be called with an unauthorised client - no signed viewing keys found, got '%s'", string(respBody2))
	}

	// make a request to GetStorageAt with userID that is in the database, but wrong parameters
	respBody3 := makeHTTPEthJSONReqWithUserID(walletHTTPPort, rpc.GetStorageAt, []interface{}{"abc", "0", nil}, userID)
	if strings.Contains(string(respBody3), userID) {
		t.Fatalf("expected response not containing userID as the parameters are wrong ")
	}

	// make a request with wrong rpcMethod
	respBody4 := makeHTTPEthJSONReqWithUserID(walletHTTPPort, rpc.GetBalance, []interface{}{"getUserID", "0", nil}, userID)
	if strings.Contains(string(respBody4), userID) {
		t.Fatalf("expected response not containing userID as the parameters are wrong ")
	}
}
