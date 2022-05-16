package txhandler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// ContractHandler defines how a contract might handle packing and unpacking transactions for Obscuro to understand them
type ContractHandler interface {
	// PackTx receives an obscurocommon.L1Transaction object and packs it into a types.TxData object
	// Nonce generation, transaction signature and any other operations are responsibility of the caller
	PackTx(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error)

	// UnPackTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	// Any transaction that is not calling the management contract is purposefully ignored
	UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction

	// Address returns the address where the contract is registered
	Address() *common.Address
}

// TxHandler handles multiple contracts ( including the management contract which has a special treatment )
type TxHandler interface {
	// PackTx receives an obscurocommon.L1Transaction object and packs it into a types.TxData object
	// Nonce generation, transaction signature and any other operations are responsibility of the caller
	PackTx(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error)

	// UnPackTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	// Any transaction that is not calling the management contract is purposefully ignored
	UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction
}
