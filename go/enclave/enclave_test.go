package enclave

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/contracts/managementcontract/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
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
		"gasEstimateSuccess":             gasEstimateSuccess,
		"gasEstimateNoVKRegistered":      gasEstimateNoVKRegistered,
		"gasEstimateNoCallMsgFrom":       gasEstimateNoCallMsgFrom,
		"gasEstimateInvalidBytes":        gasEstimateInvalidBytes,
		"gasEstimateInvalidNumParams":    gasEstimateInvalidNumParams,
		"gasEstimateInvalidParamParsing": gasEstimateInvalidParamParsing,
	}

	for name, test := range tests {
		// create the enclave
		testEnclave, err := createTestEnclave(nil)
		if err != nil {
			t.Fatal(err)
		}

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
	to := datagenerator.RandomAddress()
	callMsg := &ethereum.CallMsg{
		From: w.Address(),
		To:   &to,
		Data: []byte(ManagementContract.ManagementContractMetaData.Bin),
	}

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

	if decodeUint64 != 393608 {
		t.Fatal("unexpected gas price")
	}
}

func gasEstimateNoVKRegistered(t *testing.T, _ wallet.Wallet, enclave common.Enclave, _ *rpc.ViewingKey) {
	// use a non-registered wallet
	w := datagenerator.RandomWallet(integration.ObscuroChainID)

	// create the callMsg
	to := datagenerator.RandomAddress()
	callMsg := &ethereum.CallMsg{
		From: w.Address(),
		To:   &to,
		Data: []byte(ManagementContract.ManagementContractMetaData.Bin),
	}

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
	delete(req[0].(map[string]interface{}), "from")
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
	if !assert.ErrorContains(t, err, "no from address provided") {
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

// TestGetBalance runs the GetBalance tests
func TestGetBalance(t *testing.T) {
	tests := map[string]func(t *testing.T, w wallet.Wallet, prefund []prefundedAddress, enclave common.Enclave, vk *rpc.ViewingKey){
		"getBalanceSuccess":     getBalanceSuccess,
		"getBalanceRequestFail": getBalanceRequestFail,
	}

	for name, test := range tests {
		// create the wallet
		w := datagenerator.RandomWallet(integration.ObscuroChainID)

		// prefund the wallet
		prefundedAddresses := []prefundedAddress{
			{
				address: w.Address(),
				amount:  big.NewInt(123_000_000),
			},
		}

		// create the enclave
		testEnclave, err := createTestEnclave(prefundedAddresses)
		if err != nil {
			t.Fatal(err)
		}

		// register the VK with the enclave
		vk, err := registerWalletViewingKey(t, testEnclave, w)
		if err != nil {
			t.Fatalf("unable to register wallets VK - %s", err)
		}

		// execute the tests
		t.Run(name, func(t *testing.T) {
			test(t, w, prefundedAddresses, testEnclave, vk)
		})
	}
}

func getBalanceSuccess(t *testing.T, w wallet.Wallet, prefund []prefundedAddress, enclave common.Enclave, vk *rpc.ViewingKey) {
	// create the request payload
	req := []interface{}{prefund[0].address.Hex(), "latest"}
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
	gas, err := enclave.GetBalance(encryptedParams)
	if err != nil {
		t.Fatal(err)
	}

	// decrypt with the VK
	decryptedResult, err := vk.PrivateKey.Decrypt(gas, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// parse it
	balance, err := hexutil.DecodeBig(string(decryptedResult))
	if err != nil {
		t.Fatal(err)
	}

	// make sure its de expected value
	if prefund[0].amount.Cmp(balance) != 0 {
		t.Errorf("unexpected balance, expected %d, got %d", prefund[0].amount, balance)
	}
}

func getBalanceRequestFail(t *testing.T, w wallet.Wallet, prefund []prefundedAddress, enclave common.Enclave, vk *rpc.ViewingKey) {
	type errorTest struct {
		request  []interface{}
		errorStr string
	}
	for subtestName, test := range map[string]errorTest{
		"No1stArg": {
			request:  []interface{}{nil, "latest"},
			errorStr: "no address specified",
		},
		"No2ndArg": {
			request:  []interface{}{prefund[0].address.Hex()},
			errorStr: "required exactly two params, but received 1",
		},
		"Nil2ndArg": {
			request:  []interface{}{prefund[0].address.Hex(), nil},
			errorStr: "empty hex string",
		},
		"Rubbish2ndArg": {
			request:  []interface{}{prefund[0].address.Hex(), "Rubbish"},
			errorStr: "hex string without 0x prefix",
		},
	} {
		t.Run(subtestName, func(t *testing.T) {
			reqBytes, err := json.Marshal(test.request)
			if err != nil {
				t.Fatal(err)
			}

			// callMsg encrypted with the VK
			encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
			if err != nil {
				t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
			}

			// Run gas Estimation
			_, err = enclave.GetBalance(encryptedParams)
			if err == nil {
				t.Fatal(err)
			}

			if !assert.ErrorContains(t, err, test.errorStr) {
				t.Fatal("unexpected error")
			}
		})
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
func createTestEnclave(prefundedAddresses []prefundedAddress) (common.Enclave, error) {
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
	enclave := NewEnclave(enclaveConfig, nil, nil, logger)

	_, err := enclave.GenerateSecret()
	if err != nil {
		return nil, err
	}

	err = createFakeGenesis(enclave, prefundedAddresses)
	if err != nil {
		return nil, err
	}

	return enclave, nil
}

func createFakeGenesis(enclave common.Enclave, addresses []prefundedAddress) error {
	// Random Layer 1 block where the genesis rollup is set
	blk := types.NewBlock(&types.Header{}, nil, nil, nil, &trie.StackTrie{})
	_, err := enclave.SubmitL1Block(*blk, true)
	if err != nil {
		return err
	}

	// make sure the state is updated otherwise balances will not be available
	genesisPreallocStateDB, err := enclave.(*enclaveImpl).storage.EmptyStateDB()
	if err != nil {
		return fmt.Errorf("could not initialise empty state DB. Cause: %w", err)
	}

	for _, prefundedAddr := range addresses {
		genesisPreallocStateDB.SetBalance(prefundedAddr.address, prefundedAddr.amount)
	}

	_, err = genesisPreallocStateDB.Commit(false)
	if err != nil {
		return err
	}

	// make sure the genesis is stored the rollup storage
	genRollup := core.NewRollup(
		blk.Hash(),
		nil,
		common.L2GenesisHeight,
		gethcommon.HexToAddress("0x0"),
		[]*common.L2Tx{},
		[]common.Withdrawal{},
		common.GenerateNonce(),
		genesisPreallocStateDB.IntermediateRoot(true),
	)

	err = enclave.(*enclaveImpl).storage.StoreGenesisRollup(genRollup)
	if err != nil {
		return err
	}

	// make sure the genesis is stored as the new Head of the rollup chain
	bs := &core.BlockState{
		Block:          blk.Hash(),
		HeadRollup:     genRollup.Hash(),
		FoundNewRollup: true,
	}

	return enclave.(*enclaveImpl).storage.StoreNewHead(bs, genRollup, nil, nil)
}

type prefundedAddress struct {
	address gethcommon.Address
	amount  *big.Int
}
