package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// CallFieldTo and CallFieldFrom and CallFieldData are the relevant fields in a Call request's params.
	CallFieldTo    = "to"
	CallFieldFrom  = "from"
	CallFieldData  = "data"
	CallFieldValue = "value"
)

// ExtractTxHash returns the transaction hash from the params of an eth_getTransactionReceipt request.
func ExtractTxHash(getTxReceiptParams []byte) (gethcommon.Hash, error) {
	var paramsJSONList []string
	err := json.Unmarshal(getTxReceiptParams, &paramsJSONList)
	if err != nil {
		return gethcommon.Hash{}, fmt.Errorf("could not parse JSON params in eth_getTransactionReceipt "+
			"request. JSON params are: %s. Cause: %w", string(getTxReceiptParams), err)
	}
	if len(paramsJSONList) != 1 {
		return gethcommon.Hash{}, fmt.Errorf("expected a single param (the tx hash) but received %d params", len(paramsJSONList))
	}
	txHash := paramsJSONList[0]

	return gethcommon.HexToHash(txHash), nil
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

// ExtractEthCall extracts the eth_call [ethereum.CallMsg, gethrpc.BlockNumberOrHash] from a byte slice
func ExtractEthCall(paramBytes []byte) (*ethereum.CallMsg, *gethrpc.BlockNumberOrHash, error) {
	// extract params from byte slice to array of strings
	var paramList []interface{}
	err := json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode EthCall params - %w", err)
	}

	// params are [callMsg, block number (optional) ]
	if len(paramList) < 1 {
		return nil, nil, fmt.Errorf("required at least 1 params, but received %d", len(paramList))
	}

	// geth lowercases the field name and uses the last seen value
	var toString, fromString, dataString, valueString string
	var to, from gethcommon.Address
	var data []byte
	var value *big.Int
	var ok bool
	for field, val := range paramList[0].(map[string]interface{}) {
		switch strings.ToLower(field) {
		case CallFieldTo:
			toString, ok = val.(string)
			if !ok {
				return nil, nil, fmt.Errorf("unexpected type supplied in `to` field")
			}
			to = gethcommon.HexToAddress(toString)
		case CallFieldFrom:
			fromString, ok = val.(string)
			if !ok {
				return nil, nil, fmt.Errorf("unexpected type supplied in `from` field")
			}
			from = gethcommon.HexToAddress(fromString)
		case CallFieldData:
			dataString, ok = val.(string)
			if !ok {
				return nil, nil, fmt.Errorf("unexpected type supplied in `data` field")
			}

			// data can be nil
			if len(dataString) > 0 {
				data, err = hexutil.Decode(dataString)
				if err != nil {
					return nil, nil, fmt.Errorf("could not decode data in CallMsg - %w", err)
				}
			}
		case CallFieldValue:
			valueString, ok = val.(string)
			if !ok {
				return nil, nil, fmt.Errorf("unexpected type supplied in `value` field")
			}
			value, err = hexutil.DecodeBig(valueString)
			if err != nil {
				return nil, nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
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

	// todo actually hook the block number
	return callMsg, nil, nil
}
