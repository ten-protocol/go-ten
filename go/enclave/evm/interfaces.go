package evm

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

type EVMFacade interface {
	ExecuteTx(tx *common.L2PricedTransaction, s *state.StateDB, header *types.Header, gp *gethcore.GasPool, usedGas *uint64, tCount int, noBaseFee bool) *core.TxExecResult
	ExecuteCall(ctx context.Context, msg *gethcore.Message, s *state.StateDB, header *common.BatchHeader, isEstimateGas bool) (*gethcore.ExecutionResult, error, common.SystemError)
	DumpStateDB(label string, db *state.StateDB, from gethcommon.Address, to gethcommon.Address)
}

type ContractVisibilityReader interface {
	ReadVisibilityConfig(ctx context.Context, evm *vm.EVM, contractAddress gethcommon.Address, gasCap uint64) (*core.ContractVisibilityConfig, uint64, error)
}
