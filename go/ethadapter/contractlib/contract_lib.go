package contractlib

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

// ContractLib - common functions between contract libs
type ContractLib interface {
	DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error)
	GetContractAddr() *gethcommon.Address
	IsMock() bool
}
