package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/ten-protocol/go-ten/go/common/gethapi"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
)

// ExtractTx returns the common.L2Tx from the params of an eth_sendRawTransaction request.
func ExtractTx(txBinary string) (*common.L2Tx, error) {
	// We need to extract the transaction hex from the JSON list encoding. We remove the leading `0x`.
	txBytes := gethcommon.Hex2Bytes(txBinary[2:])

	tx := &common.L2Tx{}
	err := tx.UnmarshalBinary(txBytes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal transaction from bytes. Cause: %w", err)
	}

	return tx, nil
}

type CallParamsWithBlock struct {
	callParams *gethapi.TransactionArgs
	block      *gethrpc.BlockNumber
}

func storeTxEnabled[P any, R any](rpc *EncryptionManager, builder *CallBuilder[P, R]) bool {
	if !rpc.config.StoreExecutedTransactions {
		builder.Err = fmt.Errorf("the current TEN enclave is not configured to respond to this query")
		return false
	}
	return true
}

// newRevertError creates a revertError instance with the provided revert data.
func newRevertError(revert []byte) *errutil.DataError {
	err := vm.ErrExecutionReverted

	reason, errUnpack := abi.UnpackRevert(revert)
	if errUnpack == nil {
		err = fmt.Errorf("%w: %v", vm.ErrExecutionReverted, reason)
	}
	return &errutil.DataError{
		Err:    err.Error(),
		Code:   3, // See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
		Reason: hexutil.Encode(revert),
	}
}
