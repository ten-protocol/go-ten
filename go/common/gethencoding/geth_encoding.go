package gethencoding

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// The relevant fields in an eth_call request's params.
	callFieldTo    = "to"
	CallFieldFrom  = "from"
	callFieldData  = "data"
	callFieldValue = "value"
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
		default:
			callMsg[field] = valString
		}
	}

	return callMsg, nil
}

// ExtractEthCall extracts the eth_call ethereum.CallMsg from an interface{}
func ExtractEthCall(paramBytes interface{}) (*ethereum.CallMsg, error) {
	// geth lowercases the field name and uses the last seen value
	var valString string
	var to, from gethcommon.Address
	var data []byte
	var value *big.Int
	var ok bool
	var err error
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
			data, err = hexutil.Decode(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode data in CallMsg - %w", err)
			}
		case callFieldValue:
			value, err = hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
		}
	}

	// convert the params[0] into an ethereum.CallMsg
	callMsg := &ethereum.CallMsg{
		From:       from,
		To:         &to,
		Gas:        0,
		GasPrice:   nil,
		GasFeeCap:  nil,
		GasTipCap:  nil,
		Value:      value,
		Data:       data,
		AccessList: nil,
	}

	return callMsg, nil
}
