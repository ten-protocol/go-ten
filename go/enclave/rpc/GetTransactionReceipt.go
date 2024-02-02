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

func ExtractGetTransactionReceiptRequest(reqParams []any, rpc *EncryptionManager) (*UserRPCRequest1[types.Transaction], error) {
	// Parameters are [Hash]
	if len(reqParams) < 1 {
		return nil, fmt.Errorf("unexpected number of parameters")
	}
	txHashStr, ok := reqParams[0].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected address parameter")
	}

	txHash := gethcommon.HexToHash(txHashStr)

	// todo - optimise these calls. This can be done with a single sql
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash)
	// We retrieve the transaction.
	tx, _, _, _, err := rpc.storage.GetTransaction(txHash)
	if err != nil {
		rpc.logger.Trace("error getting tx ", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return nil, nil
		}
		return nil, err
	}

	// We retrieve the sender's address.
	sender, err := core.GetSender(tx)
	if err != nil {
		rpc.logger.Trace("error getting sender tx ", log.TxKey, txHash, log.ErrKey, err)
		return nil, fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionReceipt response. Cause: %w", err)
	}
	return &UserRPCRequest1[types.Transaction]{&sender, tx}, nil
}

func ExecuteGetTransactionReceipt(decodedParams *UserRPCRequest1[types.Transaction], rpc *EncryptionManager) (*UserResponse[types.Receipt], error) {
	if decodedParams == nil {
		return nil, nil
	}
	tx := decodedParams.Param1
	sender := decodedParams.Sender

	txHash := tx.Hash()
	// We retrieve the transaction receipt.
	txReceipt, err := rpc.storage.GetTransactionReceipt(txHash)
	if err != nil {
		rpc.logger.Trace("error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return nil, nil
		}
		// this is a system error
		err = fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
		return nil, err
	}

	// We filter out irrelevant logs.
	txReceipt.Logs, err = events.FilterLogsForReceipt(txReceipt, sender, rpc.storage)
	if err != nil {
		rpc.logger.Error("error filter logs ", log.TxKey, txHash, log.ErrKey, err)
		// this is a system error
		return nil, err
	}

	rpc.logger.Trace("Successfully retrieved receipt for ", log.TxKey, txHash, "rec", txReceipt)

	return &UserResponse[types.Receipt]{txReceipt, nil}, nil
}
