package responses

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type ViewingKeyEncryptor func([]byte) ([]byte, error)

// UserResponse - The response struct that contains either data or result
// which will be decoded only on the client side.
type UserResponse[T any] struct {
	Result *T
	ErrStr *string
}

// Error - converts the encoded string in the response into a normal error and returns it.
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
	BalanceType = hexutil.Big
	CallType    = string
	ReceiptType = types.Receipt
	RawTxType   = common.Hash
	TxType      = types.Transaction
	NonceType   = string
	GasType     = hexutil.Uint64
	LogsType    = []*types.Log
)
