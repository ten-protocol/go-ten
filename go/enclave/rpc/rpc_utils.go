package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

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

// ExtractCallParamTo extracts and parses the `to` field of an eth_call request.
func ExtractCallParamTo(callParams []byte) (gethcommon.Address, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(callParams, &paramsJSONMap)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not parse JSON params in eth_call request. JSON params are: %s. Cause: %w", string(callParams), err)
	}

	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	contractAddressString, ok := txArgs.(map[string]interface{})[CallFieldTo].(string)
	if !ok {
		return gethcommon.Address{}, fmt.Errorf("`to` field in request params was missing or not of expected type string")
	}

	return gethcommon.HexToAddress(contractAddressString), nil
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
