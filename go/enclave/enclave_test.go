package enclave

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"

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
	tests := map[string]func(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *rpc.ViewingKey){
		"gasEstimateSuccess":        gasEstimateSuccess,
		"gasEstimateNoVKRegistered": gasEstimateNoVKRegistered,
		"gasEstimateNoCallMsgFrom":  gasEstimateNoCallMsgFrom,
		"gasEstimateInvalidCallMsg": gasEstimateInvalidCallMsg,
	}

	for name, test := range tests {
		// create the enclave
		randomEnclave := createRandomEnclave()

		// create the wallet
		w := datagenerator.RandomWallet(integration.ObscuroChainID)

		// register the VK with the enclave
		vk, err := registerWalletViewingKey(t, randomEnclave, w)
		if err != nil {
			t.Fatalf("unable to register wallets VK - %s", err)
		}

		// execute the tests
		t.Run(name, func(t *testing.T) {
			test(t, w, randomEnclave, vk)
		})
	}
}

func gasEstimateSuccess(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()
	callMsgBytes, err := json.Marshal(callMsg)
	if err != nil {
		t.Error(err)
	}

	// callMsg serialized as a Hex
	callMsgHex := "0x" + gethcommon.Bytes2Hex(callMsgBytes)
	if err != nil {
		panic(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, []byte("[\""+callMsgHex+"\"]"), nil, nil)
	if err != nil {
		t.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %s", []byte{123}, err)
	}

	// Run gas Estimation
	gas, err := enclave.EstimateGas(encryptedParams)
	if err != nil {
		t.Error(err)
	}

	// decrypt with the VK
	decryptedResult, err := vk.PrivateKey.Decrypt(gas, nil, nil)
	if err != nil {
		t.Error(err)
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
	callMsgBytes, err := json.Marshal(callMsg)
	if err != nil {
		t.Error(err)
	}

	// callMsg serialized as a Hex
	callMsgHex := "0x" + gethcommon.Bytes2Hex(callMsgBytes)
	if err != nil {
		panic(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, []byte("[\""+callMsgHex+"\"]"), nil, nil)
	if err != nil {
		t.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %s", []byte{123}, err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if err == nil {
		t.Fatal("enclave does not have the viewing key to successfully encrypt")
	}
}

func gasEstimateNoCallMsgFrom(t *testing.T, _ wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsgBytes, err := json.Marshal(callMsg)
	if err != nil {
		t.Error(err)
	}

	// callMsg serialized as a Hex
	callMsgHex := "0x" + gethcommon.Bytes2Hex(callMsgBytes)
	if err != nil {
		panic(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, []byte("[\""+callMsgHex+"\"]"), nil, nil)
	if err != nil {
		t.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %s", []byte{123}, err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if err == nil {
		t.Fatal("enclave does not have the viewing key to successfully encrypt")
	}
}

func gasEstimateInvalidCallMsg(t *testing.T, _ wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// create a L2Tx instead
	callMsg := datagenerator.CreateL2Tx()
	callMsgBytes, err := json.Marshal(callMsg)
	if err != nil {
		t.Error(err)
	}

	// callMsg serialized as a Hex
	callMsgHex := "0x" + gethcommon.Bytes2Hex(callMsgBytes)
	if err != nil {
		panic(err)
	}

	// callMsg encrypted with the VK
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, []byte("[\""+callMsgHex+"\"]"), nil, nil)
	if err != nil {
		t.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %s", []byte{123}, err)
	}

	// Run gas Estimation
	_, err = enclave.EstimateGas(encryptedParams)
	if err == nil {
		t.Fatal("enclave should not parse invalid params")
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

// createRandomEnclave returns a new instance of the enclave with random values
func createRandomEnclave() common.Enclave {
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
	return NewEnclave(enclaveConfig, nil, nil, nil)
}
