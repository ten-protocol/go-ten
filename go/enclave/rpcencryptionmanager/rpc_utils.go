package rpcencryptionmanager

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// ExtractTxHash - Returns the transaction hash from the params of a eth_getTransactionReceipt request.
func ExtractTxHash(getTxReceiptParams []byte) (common.Hash, error) {
	var paramsJSONList []string
	err := json.Unmarshal(getTxReceiptParams, &paramsJSONList)
	if err != nil {
		return common.Hash{}, fmt.Errorf("could not parse JSON params in eth_getTransactionReceipt "+
			"request. JSON params are: %s. Cause: %w", string(getTxReceiptParams), err)
	}
	txHash := common.HexToHash(paramsJSONList[0]) // The only argument is the transaction hash.
	return txHash, err
}
