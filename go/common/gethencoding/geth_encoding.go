package gethencoding

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// The relevant fields in an eth_call request's params.
	callFieldTo                   = "to"
	CallFieldFrom                 = "from"
	callFieldData                 = "data"
	callFieldValue                = "value"
	callFieldGas                  = "gas"
	callFieldGasPrice             = "gasprice"
	callFieldMaxFeePerGas         = "maxfeepergas"
	callFieldMaxPriorityFeePerGas = "maxpriorityfeepergas"
)

// ExtractEthCallMapString extracts the eth_call ethereum.CallMsg from an interface{}
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

// ExtractEthCall extracts the eth_call ethereum.CallMsg from an interface{}
func ExtractEthCall(paramBytes interface{}) (*gethapi.TransactionArgs, error) {
	// geth lowercases the field name and uses the last seen value
	var valString string
	var to, from gethcommon.Address
	var data *hexutil.Bytes
	var value, gasPrice, maxFeePerGas, maxPriorityFeePerGas *hexutil.Big
	var ok bool
	var gas *hexutil.Uint64

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
			to = gethcommon.HexToAddress(valString)
		case CallFieldFrom:
			from = gethcommon.HexToAddress(valString)
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
		From:                 &from,
		To:                   &to,
		Gas:                  gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		Value:                value,
		Data:                 data,
		AccessList:           nil,
	}

	return callMsg, nil
}
