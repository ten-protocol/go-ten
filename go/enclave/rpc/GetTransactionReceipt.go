package rpc

import (
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
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
	bareReceipt, err := rpc.storage.GetTransactionReceipt(builder.ctx, txHash, requester)
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
	r := bareReceipt.MarshalToJson()
	builder.ReturnValue = &r
	return nil
}
