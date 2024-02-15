package rpc

import (
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/events"
)

func GetTransactionReceiptValidate(reqParams []any, builder *CallBuilder[gethcommon.Hash, types.Receipt], _ *EncryptionManager) error {
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

func GetTransactionReceiptExecute(builder *CallBuilder[gethcommon.Hash, types.Receipt], rpc *EncryptionManager) error {
	txHash := *builder.Param
	// todo - optimise these calls. This can be done with a single sql
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash)
	// We retrieve the transaction.
	tx, _, _, _, err := rpc.storage.GetTransaction(txHash) //nolint:dogsled
	if err != nil {
		rpc.logger.Trace("error getting tx ", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		return err
	}

	// We retrieve the txSigner's address.
	txSigner, err := core.GetTxSigner(tx)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	if txSigner.Hex() != builder.VK.AccountAddress.Hex() {
		builder.Status = NotAuthorised
		return nil
	}

	// We retrieve the transaction receipt.
	txReceipt, err := rpc.storage.GetTransactionReceipt(txHash)
	if err != nil {
		rpc.logger.Trace("error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		// this is a system error
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	// We filter out irrelevant logs.
	txReceipt.Logs, err = events.FilterLogsForReceipt(txReceipt, &txSigner, rpc.storage)
	if err != nil {
		rpc.logger.Error("error filter logs ", log.TxKey, txHash, log.ErrKey, err)
		// this is a system error
		return err
	}

	rpc.logger.Trace("Successfully retrieved receipt for ", log.TxKey, txHash, "rec", txReceipt)
	builder.ReturnValue = txReceipt
	return nil
}
