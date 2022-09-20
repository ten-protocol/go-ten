package gethenconding

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// CallFieldTo and CallFieldFrom and CallFieldData are the relevant fields in a Call request's params.
	CallFieldTo    = "to"
	CallFieldFrom  = "from"
	CallFieldData  = "data"
	CallFieldValue = "value"
)

// ExtractEthCall extracts the eth_call ethereum.CallMsg from an interface{}
func ExtractEthCall(paramBytes interface{}) (*ethereum.CallMsg, error) {
	// geth lowercases the field name and uses the last seen value
	var toString, fromString, dataString, valueString string
	var to, from gethcommon.Address
	var data []byte
	var value *big.Int
	var ok bool
	var err error
	for field, val := range paramBytes.(map[string]interface{}) {
		switch strings.ToLower(field) {
		case CallFieldTo:
			toString, ok = val.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type supplied in `to` field")
			}
			to = gethcommon.HexToAddress(toString)
		case CallFieldFrom:
			fromString, ok = val.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type supplied in `from` field")
			}
			from = gethcommon.HexToAddress(fromString)
		case CallFieldData:
			dataString, ok = val.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type supplied in `data` field")
			}

			// data can be nil
			if len(dataString) > 0 {
				data, err = hexutil.Decode(dataString)
				if err != nil {
					return nil, fmt.Errorf("could not decode data in CallMsg - %w", err)
				}
			}
		case CallFieldValue:
			valueString, ok = val.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type supplied in `value` field")
			}
			value, err = hexutil.DecodeBig(valueString)
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
