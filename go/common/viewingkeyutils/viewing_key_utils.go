package viewingkeyutils

import (
	"encoding/json"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Viewing-key-related functionality shared across the viewing-key client and the enclave.

const CallFieldFrom = "from"

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

// GetViewingKeyAddressForBalanceRequest returns the address whose viewing key should be used to encrypt the response,
// given the params of an eth_getBalance request.
func GetViewingKeyAddressForBalanceRequest(balanceParams []byte) (gethcommon.Address, error) {
	var paramsJSONMap []string
	err := json.Unmarshal(balanceParams, &paramsJSONMap)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not parse JSON params in eth_getBalance request. JSON "+
			"params are: %s. Cause: %w", string(balanceParams), err)
	}
	// The first argument is the address, the second the block.
	return gethcommon.HexToAddress(paramsJSONMap[0]), nil
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
	return from, nil
}

// ExtractTx returns the common.L2Tx from the params of an eth_sendRawTransaction request.
func ExtractTx(sendRawTxParams []byte) (*common.L2Tx, error) {
	// We need to extract the transaction hex from the JSON list encoding. We remove the leading `"[0x`, and the trailing `]"`.
	txBinary := sendRawTxParams[4 : len(sendRawTxParams)-2]
	txBytes := gethcommon.Hex2Bytes(string(txBinary))

	tx := &common.L2Tx{}
	err := tx.UnmarshalBinary(txBytes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall transaction from bytes. Cause: %w", err)
	}

	return tx, nil
}
