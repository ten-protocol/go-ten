package responses

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type Encryptor func([]byte) ([]byte, error)

type UserResponse[T any] struct {
	Result *T
	ErrStr *string
}

func (ur *UserResponse[T]) Error() error {
	if ur.ErrStr != nil {
		return fmt.Errorf(*ur.ErrStr)
	}
	return nil
}

// Responses

type (
	Balance   = EnclaveResponse // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	Call      = EnclaveResponse // As above, but for an RPC call request.
	TxReceipt = EnclaveResponse // As above, but for an RPC getTransactionReceipt request.
	RawTx     = EnclaveResponse // As above, but for an RPC sendRawTransaction request.
	TxByHash  = EnclaveResponse // As above, but for an RPC getTransactionByHash request.
	TxCount   = EnclaveResponse // As above, but for an RPC getTransactionCount request.
	Gas       = EnclaveResponse // As above, but for an RPC estimateGas response.
	Logs      = EnclaveResponse
)

// Data Types

type (
	BalanceType = UserResponse[hexutil.Big]
	CallType    = UserResponse[string]
	ReceiptType = UserResponse[types.Receipt]
	RawTxType   = UserResponse[common.Hash]
	TxType      = UserResponse[types.Transaction]
	NonceType   = UserResponse[string]
	GasType     = UserResponse[hexutil.Uint64]
	LogsType    = UserResponse[[]*types.Log]
)
