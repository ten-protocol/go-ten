package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// CallFieldTo and CallFieldFrom and CallFieldData are the relevant fields in a Call request's params.
	CallFieldTo   = "to"
	CallFieldFrom = "from"
	CallFieldData = "data"
)

// ExtractTxHash returns the transaction hash from the params of an eth_getTransactionReceipt request.
func ExtractTxHash(getTxReceiptParams []byte) (gethcommon.Hash, error) {
	var paramsJSONList []string
	err := json.Unmarshal(getTxReceiptParams, &paramsJSONList)
	if err != nil {
		return gethcommon.Hash{}, fmt.Errorf("could not parse JSON params in eth_getTransactionReceipt "+
			"request. JSON params are: %s. Cause: %w", string(getTxReceiptParams), err)
	}
	txHash := gethcommon.HexToHash(paramsJSONList[0]) // The only argument is the transaction hash.
	return txHash, err
}

// ExtractTx returns the common.L2Tx from the params of an eth_sendRawTransaction request.
func ExtractTx(sendRawTxParams []byte) (*common.L2Tx, error) {
	// We need to extract the transaction hex from the JSON list encoding. We remove the leading `"[0x`, and the trailing `]"`.
	txBinary := sendRawTxParams[4 : len(sendRawTxParams)-2]
	txBytes := gethcommon.Hex2Bytes(string(txBinary))

	tx := &common.L2Tx{}
	err := tx.UnmarshalBinary(txBytes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal transaction from bytes. Cause: %w", err)
	}

	return tx, nil
}

// ExtractAddress - Returns the address from a common.EncryptedParamsGetTransactionCount blob
func ExtractAddress(getTransactionCountParams []byte) (gethcommon.Address, error) {
	var paramsJSONList []string
	err := json.Unmarshal(getTransactionCountParams, &paramsJSONList)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not parse JSON params in eth_getTransactionCount request. Cause: %w", err)
	}
	txHash := gethcommon.HexToAddress(paramsJSONList[0]) // The only argument is the transaction hash.
	return txHash, err
}

// GetViewingKeyAddressForTransaction returns the address whose viewing key should be used to encrypt the response,
// given a transaction.
func GetViewingKeyAddressForTransaction(tx *common.L2Tx) (gethcommon.Address, error) {
	// TODO - Once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`.
	signer := types.NewLondonSigner(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not recover sender for transaction. Cause: %w", err)
	}
	return sender, nil
}

// ExtractCallParamTo extracts and parses the `to` field of an eth_call request.
func ExtractCallParamTo(callParams []byte) (*gethcommon.Address, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(callParams, &paramsJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON params in eth_call request. JSON params are: %s. Cause: %w", string(callParams), err)
	}

	// to field is null on contract creation
	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	if to := txArgs.(map[string]interface{})[CallFieldTo]; to == nil {
		return nil, nil //nolint:nilnil
	}

	contractAddressString, ok := txArgs.(map[string]interface{})[CallFieldTo].(string)
	if !ok {
		return nil, fmt.Errorf("`to` field in request params was missing or not of expected type string")
	}
	contractAddress := gethcommon.HexToAddress(contractAddressString)
	return &contractAddress, nil
}

// ExtractCallParamFrom extracts and parses the `from` field of an eth_call request.
// This is also the address whose viewing key should be used to encrypt the response.
func ExtractCallParamFrom(callParams []byte) (gethcommon.Address, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(callParams, &paramsJSONMap)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not parse JSON params in eth_call request. JSON "+
			"params are: %s. Cause: %w", string(callParams), err)
	}

	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	fromString, ok := txArgs.(map[string]interface{})[CallFieldFrom].(string)
	if !ok {
		return gethcommon.Address{}, fmt.Errorf("`from` field in request params is missing or was not of " +
			"expected type string. The `from` field is required to encrypt the response")
	}

	from := gethcommon.HexToAddress(fromString)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not decode data in Call request. Cause: %w", err)
	}
	return from, nil
}

// ExtractCallParamData extracts and parses the `data` field of an eth_call request.
func ExtractCallParamData(callParams []byte) ([]byte, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(callParams, &paramsJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON params in eth_call request. JSON params are: %s. Cause: %w", string(callParams), err)
	}

	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	dataString, ok := txArgs.(map[string]interface{})[CallFieldData].(string)
	if !ok {
		return nil, fmt.Errorf("`data` field in request params is missing or was not of expected type string")
	}

	data, err := hexutil.Decode(dataString)
	if err != nil {
		return nil, fmt.Errorf("could not decode data in Call request. Cause: %w", err)
	}
	return data, nil
}
