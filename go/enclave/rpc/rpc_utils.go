package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ExtractTxHash returns the transaction hash from the params of an eth_getTransactionReceipt request.
func ExtractTxHash(getTxReceiptParams []byte) (gethcommon.Hash, error) {
	var txHex string
	err := json.Unmarshal(getTxReceiptParams, &txHex)
	if err != nil {
		return gethcommon.Hash{}, fmt.Errorf("could not parse JSON params in eth_getTransactionReceipt "+
			"request. JSON params are: %s. Cause: %w", string(getTxReceiptParams), err)
	}
	return gethcommon.HexToHash(txHex), nil
}

// ExtractTx returns the common.L2Tx from the params of an eth_sendRawTransaction request.
func ExtractTx(sendRawTxParams []byte) (*common.L2Tx, error) {
	// We remove the leading `"0x` and the trailing `"` from the transaction hex.
	txBinary := sendRawTxParams[3 : len(sendRawTxParams)-1]
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

// ConvertToCallMsg converts the interface to a *ethereum.CallMsg
func ConvertToCallMsg(callMsgInterface interface{}) (*ethereum.CallMsg, error) {
	callMsgBytes, err := json.Marshal(callMsgInterface)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall callMsg - %w", err)
	}
	var callMsg ethereum.CallMsg
	err = json.Unmarshal(callMsgBytes, &callMsg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse callMsg - %w", err)
	}
	return &callMsg, err
}
