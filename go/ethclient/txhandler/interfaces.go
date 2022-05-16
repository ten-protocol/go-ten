package txhandler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// ContractHandler defines how a contract might handle packing and unpacking transactions for Obscuro to understand them
type ContractHandler interface {
	// Pack receives an obscurocommon.L1Transaction object and packs it into a types.TxData object
	Pack(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error)

	// UnPack receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	UnPack(tx *types.Transaction) obscurocommon.L1Transaction

	// Address returns the address where the contract is deployed
	Address() *common.Address
}

// TxHandler handles multiple contracts routing the transaction to the correct contract if it's registered
// ( including the management contract which has a special treatment )
type TxHandler interface {
	// PackTx receives an obscurocommon.L1Transaction object and packs it into a types.TxData object
	PackTx(tx obscurocommon.L1Transaction, from common.Address, nonce uint64) (types.TxData, error)

	// UnPackTx receives a *types.Transaction and converts it to an obscurocommon.L1Transaction
	// Any transaction NOT calling a registered contract will be ignored
	UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction
}
