//nolint:unused
package enclave

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"

	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const _testEnclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"

var _enclavePubKey *ecies.PublicKey

func init() { //nolint:gochecknoinits
	// fetch the usable enclave pub key
	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(_testEnclavePublicKeyHex))
	if err != nil {
		panic(err)
	}

	_enclavePubKey = ecies.ImportECDSAPublic(enclPubECDSA)
}

// TestGasEstimation runs the GasEstimation tests
func TestGasEstimation(t *testing.T) {
	// TODO create a Headstate in hs := rc.storage.FetchHeadState()
	t.Skip("Skipping the gas estimation tests..")
	tests := map[string]func(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *rpc.ViewingKey){
		"gasEstimateSuccess":             gasEstimateSuccess,
		"gasEstimateNoVKRegistered":      gasEstimateNoVKRegistered,
		"gasEstimateNoCallMsgFrom":       gasEstimateNoCallMsgFrom,
		"gasEstimateInvalidBytes":        gasEstimateInvalidBytes,
		"gasEstimateInvalidNumParams":    gasEstimateInvalidNumParams,
		"gasEstimateInvalidParamParsing": gasEstimateInvalidParamParsing,
	}

	for name, test := range tests {
		// create the enclave
		testEnclave := createTestEnclave()

		// create the wallet
		w := datagenerator.RandomWallet(integration.ObscuroChainID)

		// register the VK with the enclave
		vk, err := registerWalletViewingKey(t, testEnclave, w)
		if err != nil {
			t.Fatalf("unable to register wallets VK - %s", err)
		}

		// execute the tests
		t.Run(name, func(t *testing.T) {
			test(t, w, testEnclave, vk)
		})
	}
}

func gasEstimateSuccess(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request payload
	req := []interface{}{obsclient.ToCallArg(*callMsg), nil}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	gas, err := enclave.EstimateGas(encryptedParams)
	if err != nil {
		t.Fatal(err)
	}

	// decrypt with the VK
	decryptedResult, err := vk.PrivateKey.Decrypt(gas, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// parse it to Uint64
	decodeUint64, err := hexutil.DecodeUint64(string(decryptedResult))
	if err != nil {
		t.Fatal(err)
	}

	if decodeUint64 != 5_000_000_000 {
		t.Fatal("unexpected gas price")
	}
}

func gasEstimateNoVKRegistered(t *testing.T, _ wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// use a non-registered wallet
	w := datagenerator.RandomWallet(integration.ObscuroChainID)

	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request
	req := []interface{}{obsclient.ToCallArg(*callMsg), nil}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, err, "could not encrypt bytes because it does not have a viewing key for account") {
		t.Fatalf("unexpected error - %s", err)
	}
}

func gasEstimateNoCallMsgFrom(t *testing.T, _ wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()

	// create the request
	req := []interface{}{obsclient.ToCallArg(*callMsg), nil}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, err, "could not encrypt bytes because it does not have a viewing key for account") {
		t.Fatalf("unexpected error - %s", err)
	}
}

func gasEstimateInvalidBytes(t *testing.T, w wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request
	req := []interface{}{obsclient.ToCallArg(*callMsg), nil}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	reqBytes = append(reqBytes, []byte("this should break stuff")...)

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, err, "invalid character") {
		t.Fatalf("unexpected error - %s", err)
	}
}

func gasEstimateInvalidNumParams(t *testing.T, w wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request
	var req []interface{}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, err, "required at least 1 params, but received 0") {
		t.Fatal("unexpected error")
	}
}

func gasEstimateInvalidParamParsing(t *testing.T, w wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request
	req := []interface{}{callMsg}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, err, "unexpected type supplied in") {
		t.Fatal("unexpected error")
	}
}

// registerWalletViewingKey takes a wallet and registers a VK with the enclave
func registerWalletViewingKey(t *testing.T, enclave common.Enclave, w wallet.Wallet) (*rpc.ViewingKey, error) {
	// generate the VK from the wallet
	key, err := rpc.GenerateAndSignViewingKey(w)
	if err != nil {
		t.Fatal(err)
	}

	// encrypt the VK
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, _enclavePubKey, key.PublicKey, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// add the VK to the enclave
	return key, enclave.AddViewingKey(encryptedViewingKeyBytes, key.SignedKey)
}

// createTestEnclave returns a test instance of the enclave
func createTestEnclave() common.Enclave {
	rndAddr := gethcommon.HexToAddress("contract1")
	rndAddr2 := gethcommon.HexToAddress("contract2")
	enclaveConfig := config.EnclaveConfig{
		L1ChainID:              integration.EthereumChainID,
		ObscuroChainID:         integration.ObscuroChainID,
		WillAttest:             false,
		UseInMemoryDB:          true,
		ERC20ContractAddresses: []*gethcommon.Address{&rndAddr, &rndAddr2},
		MinGasPrice:            big.NewInt(1),
	}
	logger := log.New(log.TestLogCmp, int(gethlog.LvlError), log.SysOut)
	return NewEnclave(enclaveConfig, nil, nil, logger)
}
