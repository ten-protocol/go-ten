package enclave

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/responses"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/assert"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

const _testEnclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"

// _successfulRollupGasPrice can be deterministically calculated when evaluating the management smart contract.
// It should change only when there are changes to the smart contract or if the gas estimation algorithm is modified.
// Other changes would mean something is broken.
const _successfulRollupGasPrice = 336360

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
	tests := map[string]func(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *viewingkey.ViewingKey){
		"gasEstimateSuccess":             gasEstimateSuccess,
		"gasEstimateNoCallMsgFrom":       gasEstimateNoCallMsgFrom,
		"gasEstimateInvalidBytes":        gasEstimateInvalidBytes,
		"gasEstimateInvalidNumParams":    gasEstimateInvalidNumParams,
		"gasEstimateInvalidParamParsing": gasEstimateInvalidParamParsing,
	}

	idx := 100
	for name, test := range tests {
		// create the enclave
		testEnclave, err := createTestEnclave(nil, idx)
		idx++
		if err != nil {
			t.Fatal(err)
		}

		// create the wallet
		w := datagenerator.RandomWallet(integration.ObscuroChainID)

		// create a VK
		vk, err := viewingkey.GenerateViewingKeyForWallet(w)
		if err != nil {
			t.Fatal(err)
		}

		// execute the tests
		t.Run(name, func(t *testing.T) {
			test(t, w, testEnclave, vk)
		})
	}
}

func gasEstimateSuccess(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	// create the callMsg
	to := datagenerator.RandomAddress()
	callMsg := &ethereum.CallMsg{
		From: w.Address(),
		To:   &to,
		Data: []byte(ManagementContract.ManagementContractMetaData.Bin),
	}

	// create the request payload
	req := []interface{}{
		[]interface{}{
			hexutil.Encode(vk.PublicKey),
			hexutil.Encode(vk.Signature),
		},
		obsclient.ToCallArg(*callMsg),
		nil,
	}
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
	gas, _ := enclave.EstimateGas(encryptedParams)
	if gas.Error() != nil {
		t.Fatal(gas.Error())
	}

	// decrypt with the VK
	decryptedResult, err := vk.PrivateKey.Decrypt(gas.EncUserResponse, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	gasEstimate, err := responses.DecodeResponse[string](decryptedResult)
	if err != nil {
		t.Fatal(err)
	}

	// parse it to Uint64
	decodeUint64, err := hexutil.DecodeUint64(*gasEstimate)
	if err != nil {
		t.Fatal(err)
	}

	if decodeUint64 != _successfulRollupGasPrice {
		t.Fatalf("unexpected gas price %d", decodeUint64)
	}
}

func gasEstimateNoCallMsgFrom(t *testing.T, _ wallet.Wallet, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()

	// create the request
	req := []interface{}{
		[]interface{}{
			hexutil.Encode(vk.PublicKey),
			hexutil.Encode(vk.Signature),
		},
		obsclient.ToCallArg(*callMsg),
		nil,
	}
	delete(req[1].(map[string]interface{}), "from")
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
	resp, _ := enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, resp.Error(), "no from address provided") {
		t.Fatalf("unexpected error - %s", resp.Error())
	}
}

func gasEstimateInvalidBytes(t *testing.T, w wallet.Wallet, enclave common.Enclave, _ *viewingkey.ViewingKey) {
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
	resp, _ := enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, resp.Error(), "invalid character") {
		t.Fatalf("unexpected error - %s", resp.Error())
	}
}

func gasEstimateInvalidNumParams(t *testing.T, _ wallet.Wallet, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	// create the request
	req := []interface{}{
		[]interface{}{
			hexutil.Encode(vk.PublicKey),
			hexutil.Encode(vk.Signature),
		},
	}
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
	resp, _ := enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, resp.Error(), "unexpected number of parameters") {
		t.Fatal("unexpected error")
	}
}

func gasEstimateInvalidParamParsing(t *testing.T, w wallet.Wallet, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	// create the callMsg
	callMsg := datagenerator.CreateCallMsg()
	callMsg.From = w.Address()

	// create the request
	req := []interface{}{
		[]interface{}{
			hexutil.Encode(vk.PublicKey),
			hexutil.Encode(vk.Signature),
		},
		callMsg,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// callMsg encrypted with the Enclave Key
	encryptedParams, err := ecies.Encrypt(rand.Reader, _enclavePubKey, reqBytes, nil, nil)
	if err != nil {
		t.Fatalf("could not encrypt the following request params with enclave public key - %s", err)
	}

	// Run gas Estimation
	resp, _ := enclave.EstimateGas(encryptedParams)
	if !assert.ErrorContains(t, resp.Error(), "unexpected type supplied in") {
		t.Fatal("unexpected error")
	}
}

// TestGetBalance runs the GetBalance tests
func TestGetBalance(t *testing.T) {
	tests := map[string]func(t *testing.T, prefund []genesis.Account, enclave common.Enclave, vk *viewingkey.ViewingKey){
		"getBalanceSuccess":             getBalanceSuccess,
		"getBalanceRequestUnsuccessful": getBalanceRequestUnsuccessful,
	}

	idx := 0
	for name, test := range tests {
		// create the wallet
		w := datagenerator.RandomWallet(integration.ObscuroChainID)

		// prefund the wallet
		prefundedAddresses := []genesis.Account{
			{
				Address: w.Address(),
				Amount:  big.NewInt(123_000_000),
			},
		}

		// create the enclave
		testEnclave, err := createTestEnclave(prefundedAddresses, idx)
		idx++
		if err != nil {
			t.Fatal(err)
		}

		// create a VK
		vk, err := viewingkey.GenerateViewingKeyForWallet(w)
		if err != nil {
			t.Fatal(err)
		}

		// execute the tests
		t.Run(name, func(t *testing.T) {
			test(t, prefundedAddresses, testEnclave, vk)
		})
	}
}

func getBalanceSuccess(t *testing.T, prefund []genesis.Account, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	// create the request payload
	req := []interface{}{
		[]interface{}{
			hexutil.Encode(vk.PublicKey),
			hexutil.Encode(vk.Signature),
		},
		prefund[0].Address.Hex(),
		"latest",
	}
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
	gas, _ := enclave.GetBalance(encryptedParams)
	if err != nil {
		t.Fatal(err)
	}

	// decrypt with the VK
	decryptedResult, err := vk.PrivateKey.Decrypt(gas.EncUserResponse, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// parse it
	balance, err := responses.DecodeResponse[hexutil.Big](decryptedResult)
	if err != nil {
		t.Fatal(err)
	}

	// make sure its de expected value
	if prefund[0].Amount.Cmp(balance.ToInt()) != 0 {
		t.Errorf("unexpected balance, expected %d, got %d", prefund[0].Amount, balance.ToInt())
	}
}

func getBalanceRequestUnsuccessful(t *testing.T, prefund []genesis.Account, enclave common.Enclave, vk *viewingkey.ViewingKey) {
	type errorTest struct {
		request  []interface{}
		errorStr string
	}
	vkSerialized := []interface{}{
		hexutil.Encode(vk.PublicKey),
		hexutil.Encode(vk.Signature),
	}

	for subtestName, test := range map[string]errorTest{
		"No1stArg": {
			request:  []interface{}{vkSerialized, nil, "latest"},
			errorStr: "no address specified",
		},
		"No2ndArg": {
			request:  []interface{}{vkSerialized, prefund[0].Address.Hex()},
			errorStr: "unexpected number of parameters",
		},
		"Nil2ndArg": {
			request:  []interface{}{vkSerialized, prefund[0].Address.Hex(), nil},
			errorStr: "unable to extract requested block number - not found",
		},
		"Rubbish2ndArg": {
			request:  []interface{}{vkSerialized, prefund[0].Address.Hex(), "Rubbish"},
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
			enclaveResp, _ := enclave.GetBalance(encryptedParams)
			err = enclaveResp.Error()

			// If there is no enclave error we must get
			// the internal user error
			if err == nil {
				encodedResp, encErr := vk.PrivateKey.Decrypt(enclaveResp.EncUserResponse, nil, nil)
				if encErr != nil {
					t.Fatal(encErr)
				}
				_, err = responses.DecodeResponse[hexutil.Big](encodedResp)
			}

			if !assert.ErrorContains(t, err, test.errorStr) {
				t.Fatal("unexpected error")
			}
		})
	}
}

// createTestEnclave returns a test instance of the enclave
func createTestEnclave(prefundedAddresses []genesis.Account, idx int) (common.Enclave, error) {
	enclaveConfig := &config.EnclaveConfig{
		HostID:         gethcommon.BigToAddress(big.NewInt(int64(idx))),
		L1ChainID:      integration.EthereumChainID,
		ObscuroChainID: integration.ObscuroChainID,
		WillAttest:     false,
		UseInMemoryDB:  true,
		MinGasPrice:    big.NewInt(1),
	}
	logger := log.New(log.TestLogCmp, int(gethlog.LvlError), log.SysOut)

	obsGenesis := &genesis.TestnetGenesis
	if len(prefundedAddresses) > 0 {
		obsGenesis = &genesis.Genesis{Accounts: prefundedAddresses}
	}
	enclave := NewEnclave(enclaveConfig, obsGenesis, nil, logger)

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

func createFakeGenesis(enclave common.Enclave, addresses []genesis.Account) error {
	// Random Layer 1 block where the genesis rollup is set
	blk := types.NewBlock(&types.Header{}, nil, nil, nil, &trie.StackTrie{})
	_, err := enclave.SubmitL1Block(*blk, make(types.Receipts, 0), true)
	if err != nil {
		return err
	}

	// make sure the state is updated otherwise balances will not be available
	genesisPreallocStateDB, err := enclave.(*enclaveImpl).storage.EmptyStateDB()
	if err != nil {
		return fmt.Errorf("could not initialise empty state DB. Cause: %w", err)
	}

	for _, prefundedAddr := range addresses {
		genesisPreallocStateDB.SetBalance(prefundedAddr.Address, prefundedAddr.Amount)
	}

	_, err = genesisPreallocStateDB.Commit(0, false)
	if err != nil {
		return err
	}

	genesisBatch := dummyBatch(blk.Hash(), common.L2GenesisHeight, genesisPreallocStateDB)

	// We update the database
	err = enclave.(*enclaveImpl).storage.StoreBatch(genesisBatch)
	if err != nil {
		return err
	}
	return enclave.(*enclaveImpl).storage.StoreExecutedBatch(genesisBatch, nil)
}

func dummyBatch(blkHash gethcommon.Hash, height uint64, state *state.StateDB) *core.Batch {
	h := common.BatchHeader{
		ParentHash:       common.L1BlockHash{},
		L1Proof:          blkHash,
		Root:             state.IntermediateRoot(true),
		Number:           big.NewInt(int64(height)),
		SequencerOrderNo: big.NewInt(int64(height + 1)), // seq number starts at 1 so need to offset
		ReceiptHash:      types.EmptyRootHash,
		Time:             uint64(time.Now().Unix()),
	}
	return &core.Batch{
		Header:       &h,
		Transactions: []*common.L2Tx{},
	}
}
