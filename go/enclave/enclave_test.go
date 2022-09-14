package enclave

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
)

const enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"

func createRandomEnclave() common.Enclave {
	rndAddr := gethcommon.HexToAddress("hi")
	rndAddr2 := gethcommon.HexToAddress("hi hi")
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

func decryptResponse(viewingKey *rpc.ViewingKey, resultBlob []byte) ([]byte, error) {
	decryptedResult, err := viewingKey.PrivateKey.Decrypt(resultBlob, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt result with viewing key - %w", err)
	}

	return decryptedResult, nil
}

func TestEnclave(t *testing.T) {
	// create the enclave
	randomEnclave := createRandomEnclave()

	// create the wallet
	w := datagenerator.RandomWallet(integration.ObscuroChainID)
	// register the wallets VK
	key, err := rpc.GenerateAndSignViewingKey(w)
	if err != nil {
		t.Fatal(err)
	}

	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		t.Fatal(err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, enclavePublicKey, key.PublicKey, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = randomEnclave.AddViewingKey(encryptedViewingKeyBytes, key.SignedKey)
	if err != nil {
		t.Fatal(err)
	}

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

	// callMsg encrypted with the VL
	encryptedParams, err := ecies.Encrypt(rand.Reader, enclavePublicKey, []byte("[\""+callMsgHex+"\"]"), nil, nil)
	if err != nil {
		t.Errorf("could not encrypt the following request params with enclave public key: %s. Cause: %s", []byte{123}, err)
	}

	// Run gas Estimation
	gas, err := randomEnclave.EstimateGas(encryptedParams)
	if err != nil {
		t.Error(err)
	}

	//
	resp, err := decryptResponse(key, gas)
	if err != nil {
		t.Error(err)
	}

	decodeUint64, err := hexutil.DecodeUint64(string(resp))
	if err != nil {
		t.Fatal(err)
	}

	if decodeUint64 != 10_000_000 {
		t.Fatal("unexpected gas price")
	}
}
