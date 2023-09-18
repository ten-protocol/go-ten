package gethencoding

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	// The relevant fields in an eth_call request's params.
	callFieldTo                   = "to"
	CallFieldFrom                 = "from"
	callFieldData                 = "data"
	callFieldValue                = "value"
	callFieldGas                  = "gas"
	callFieldNonce                = "nonce"
	callFieldGasPrice             = "gasprice"
	callFieldMaxFeePerGas         = "maxfeepergas"
	callFieldMaxPriorityFeePerGas = "maxpriorityfeepergas"
)

// ExtractEthCallMapString extracts the eth_call gethapi.TransactionArgs from an interface{}
// it ensures that :
// - All types are string
// - All keys are lowercase
// - There is only one key per value
// - From field is set by default
func ExtractEthCallMapString(paramBytes interface{}) (map[string]string, error) {
	// geth lowercase the field name and uses the last seen value
	var valString string
	var ok bool
	callMsg := map[string]string{
		// From field is set by default
		"from": gethcommon.HexToAddress("0x0").Hex(),
	}
	for field, val := range paramBytes.(map[string]interface{}) {
		if val == nil {
			continue
		}
		valString, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type supplied in `%s` field", field)
		}
		if len(strings.TrimSpace(valString)) == 0 {
			continue
		}
		switch strings.ToLower(field) {
		case callFieldTo:
			callMsg[callFieldTo] = valString
		case CallFieldFrom:
			callMsg[CallFieldFrom] = valString
		case callFieldData:
			callMsg[callFieldData] = valString
		case callFieldValue:
			callMsg[callFieldValue] = valString
		case callFieldGas:
			callMsg[callFieldGas] = valString
		case callFieldMaxFeePerGas:
			callMsg[callFieldMaxFeePerGas] = valString
		case callFieldMaxPriorityFeePerGas:
			callMsg[callFieldMaxPriorityFeePerGas] = valString
		default:
			callMsg[field] = valString
		}
	}

	return callMsg, nil
}

// ExtractAddress returns a gethcommon.Address given an interface{}, errors if unexpected values are used
func ExtractAddress(param interface{}) (*gethcommon.Address, error) {
	if param == nil {
		return nil, fmt.Errorf("no address specified")
	}

	paramStr, ok := param.(string)
	if !ok {
		return nil, fmt.Errorf("unexpectd address value")
	}

	if len(strings.TrimSpace(paramStr)) == 0 {
		return nil, fmt.Errorf("no address specified")
	}

	addr := gethcommon.HexToAddress(param.(string))
	return &addr, nil
}

// ExtractOptionalBlockNumber defaults nil or empty block number params to latest block number
func ExtractOptionalBlockNumber(params []interface{}, idx int) (*gethrpc.BlockNumber, error) {
	if len(params) <= idx {
		return ExtractBlockNumber("latest")
	}
	if params[idx] == nil {
		return ExtractBlockNumber("latest")
	}
	if emptyStr, ok := params[idx].(string); ok && len(strings.TrimSpace(emptyStr)) == 0 {
		return ExtractBlockNumber("latest")
	}

	return ExtractBlockNumber(params[idx])
}

// ExtractBlockNumber returns a gethrpc.BlockNumber given an interface{}, errors if unexpected values are used
func ExtractBlockNumber(param interface{}) (*gethrpc.BlockNumber, error) {
	if param == nil {
		return nil, errutil.ErrNotFound
	}

	blockNumber := gethrpc.BlockNumber(0)
	err := blockNumber.UnmarshalJSON([]byte(param.(string)))
	if err != nil {
		return nil, fmt.Errorf("could not parse requested rollup number %s - %w", param.(string), err)
	}

	return &blockNumber, err
}

// ExtractEthCall extracts the eth_call gethapi.TransactionArgs from an interface{}
func ExtractEthCall(param interface{}) (*gethapi.TransactionArgs, error) {
	// geth lowercases the field name and uses the last seen value
	var valString string
	var to, from *gethcommon.Address
	var data *hexutil.Bytes
	var value, gasPrice, maxFeePerGas, maxPriorityFeePerGas *hexutil.Big
	var ok bool
	zeroUint := hexutil.Uint64(0)
	nonce := &zeroUint
	// if gas is not set it should be null
	gas := (*hexutil.Uint64)(nil)

	for field, val := range param.(map[string]interface{}) {
		if val == nil {
			continue
		}
		valString, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type supplied in `%s` field", field)
		}
		if len(strings.TrimSpace(valString)) == 0 {
			continue
		}
		switch strings.ToLower(field) {
		case callFieldTo:
			toVal := gethcommon.HexToAddress(valString)
			to = &toVal
		case CallFieldFrom:
			fromVal := gethcommon.HexToAddress(valString)
			from = &fromVal
		case callFieldData:
			dataVal, err := hexutil.Decode(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode data in CallMsg - %w", err)
			}
			data = (*hexutil.Bytes)(&dataVal)
		case callFieldValue:
			valueVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			value = (*hexutil.Big)(valueVal)
		case callFieldNonce:
			nonceVal, err := hexutil.DecodeUint64(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			nonce = (*hexutil.Uint64)(&nonceVal)
		case callFieldGas:
			gasVal, err := hexutil.DecodeUint64(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			gas = (*hexutil.Uint64)(&gasVal)

		case callFieldGasPrice:
			valueVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			value = (*hexutil.Big)(valueVal)

		case callFieldMaxFeePerGas:
			maxFeePerGasVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			maxFeePerGas = (*hexutil.Big)(maxFeePerGasVal)

		case callFieldMaxPriorityFeePerGas:
			maxPriorityFeePerGasVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			maxPriorityFeePerGas = (*hexutil.Big)(maxPriorityFeePerGasVal)
		}
	}

	// convert the params[0] into an ethereum.CallMsg
	callMsg := &gethapi.TransactionArgs{
		From:                 from,
		To:                   to,
		Gas:                  gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		Value:                value,
		Data:                 data,
		Nonce:                nonce,
		AccessList:           nil,
	}

	return callMsg, nil
}

// CreateEthHeaderForBatch - the EVM requires an Ethereum "block" header.
// In this function we are creating one from the Batch Header
func CreateEthHeaderForBatch(h *common.BatchHeader, secret []byte) (*types.Header, error) {
	// deterministically calculate private randomness that will be exposed to the evm
	randomness := crypto.CalculateRootBatchEntropy(secret, h.Number)

	baseFee := uint64(0)
	if h.BaseFee != nil {
		baseFee = h.BaseFee.Uint64()
	}

	return &types.Header{
		ParentHash:  h.ParentHash,
		Root:        h.Root,
		TxHash:      h.TxHash,
		ReceiptHash: h.ReceiptHash,
		Difficulty:  big.NewInt(0),
		Number:      h.Number,
		GasLimit:    h.GasLimit,
		GasUsed:     0,
		BaseFee:     big.NewInt(0).SetUint64(baseFee),
		Coinbase:    h.Coinbase,
		Time:        h.Time,
		MixDigest:   randomness,
		Nonce:       types.BlockNonce{},
	}, nil
}

// DecodeParamBytes decodes the parameters byte array into a slice of interfaces
// Helps each calling method to manage the positional data
func DecodeParamBytes(paramBytes []byte) ([]interface{}, error) {
	var paramList []interface{}

	if err := json.Unmarshal(paramBytes, &paramList); err != nil {
		return nil, fmt.Errorf("unable to unmarshal params - %w", err)
	}
	return paramList, nil
}

// ExtractViewingKey returns the viewingkey pubkey and the signature from the request
func ExtractViewingKey(vkBytesIntf interface{}) ([]byte, []byte, error) {
	vkBytesList, ok := vkBytesIntf.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("unable to cast the vk to []interface")
	}

	if len(vkBytesList) != 2 {
		return nil, nil, fmt.Errorf("wrong size of viewing key params")
	}

	vkPubkeyHexBytes, err := hexutil.Decode(vkBytesList[0].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode data in vk pub key - %w", err)
	}

	accountSignatureHexBytes, err := hexutil.Decode(vkBytesList[1].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode data in vk signature - %w", err)
	}

	return vkPubkeyHexBytes, accountSignatureHexBytes, nil
}

func ExtractPrivateCustomQuery(_ interface{}, query interface{}) (*common.PrivateCustomQueryListTransactions, error) {
	// Convert the map to a JSON string
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	var result common.PrivateCustomQueryListTransactions
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
