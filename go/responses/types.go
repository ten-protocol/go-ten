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
	GetBalance   = EnclaveResponse // The response for an RPC getBalance request, as a JSON object encrypted with the viewing key of the user.
	Call         = EnclaveResponse // As above, but for an RPC call request.
	GetTxReceipt = EnclaveResponse // As above, but for an RPC getTransactionReceipt request.
	SendRawTx    = EnclaveResponse // As above, but for an RPC sendRawTransaction request.
	GetTxByHash  = EnclaveResponse // As above, but for an RPC getTransactionByHash request.
	GetTxCount   = EnclaveResponse // As above, but for an RPC getTransactionCount request.
	Logs         = EnclaveResponse // As above, but for a log subscription response.
	EstimateGas  = EnclaveResponse // As above, but for an RPC estimateGas response.
	GetLogs      = EnclaveResponse
)

// Data Types

type (
	BalanceType     = UserResponse[hexutil.Big]
	CallType        = UserResponse[string]
	ReceiptType     = UserResponse[types.Receipt]
	SendRawTxType   = UserResponse[common.Hash]
	TxType          = UserResponse[types.Transaction]
	NonceType       = UserResponse[string]
	EstimateGasType = UserResponse[hexutil.Uint64]
	GetLogsType     = UserResponse[[]*types.Log]
)
