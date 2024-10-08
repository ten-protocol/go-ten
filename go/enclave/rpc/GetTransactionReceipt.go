package rpc

import (
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"
)

func GetTransactionReceiptValidate(reqParams []any, builder *CallBuilder[gethcommon.Hash, map[string]interface{}], _ *EncryptionManager) error {
	// Parameters are [Hash]
	if len(reqParams) < 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	txHashStr, ok := reqParams[0].(string)
	if !ok {
		builder.Err = fmt.Errorf("invalid transaction hash")
		return nil
	}

	txHash := gethcommon.HexToHash(txHashStr)
	builder.Param = &txHash
	return nil
}

func GetTransactionReceiptExecute(builder *CallBuilder[gethcommon.Hash, map[string]interface{}], rpc *EncryptionManager) error {
	txHash := *builder.Param
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash)
	requester := builder.VK.AccountAddress

	// We retrieve the transaction receipt.
	bareReceipt, logs, err := rpc.storage.GetTransactionReceipt(builder.ctx, txHash, requester)
	if err != nil {
		rpc.logger.Trace("error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		// this is a system error
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	rpc.logger.Trace("Successfully retrieved receipt for ", log.TxKey, txHash, "rec", bareReceipt)
	r := marshalReceipt(bareReceipt, logs)
	builder.ReturnValue = &r
	return nil
}

// marshalReceipt marshals a transaction receipt into a JSON object.
// taken from geth
func marshalReceipt(receipt *enclavedb.BareReceipt, logs []*types.Log) map[string]interface{} {
	var effGasPrice *hexutil.Big
	if receipt.EffectiveGasPrice != nil {
		effGasPrice = (*hexutil.Big)(big.NewInt(int64(*receipt.EffectiveGasPrice)))
	}

	fields := map[string]interface{}{
		"blockHash":         receipt.BlockHash,
		"blockNumber":       hexutil.Uint64(receipt.BlockNumber.Uint64()),
		"transactionHash":   receipt.TxHash,
		"transactionIndex":  hexutil.Uint64(receipt.TransactionIndex),
		"from":              receipt.From,
		"to":                receipt.To,
		"gasUsed":           hexutil.Uint64(receipt.CumulativeGasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   receipt.CreatedContract,
		"logs":              logs,
		"logsBloom":         types.Bloom{},
		"type":              hexutil.Uint(receipt.TxType),
		"effectiveGasPrice": effGasPrice,
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if logs == nil {
		fields["logs"] = []*types.Log{}
	}

	return fields
}
