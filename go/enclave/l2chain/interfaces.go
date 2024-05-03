package l2chain

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// ObscuroChain - the interface that provides the data access layer to the obscuro l2.
// Operations here should be read only.
type ObscuroChain interface {
	// AccountOwner - returns the account that owns the address.
	// For EOA - the actual address.
	// For Contracts - the address of the deployer.
	// Note - this might be subject to change if we implement a more flexible mechanism
	// todo - support BlockNumberOrHash
	AccountOwner(ctx context.Context, address gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, error)

	// GetBalanceAtBlock - will return the balance of a specific address at the specific given block number (batch number).
	GetBalanceAtBlock(ctx context.Context, accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error)

	// ObsCall - The interface for executing eth_call RPC commands against obscuro.
	ObsCall(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error)

	// ObsCallAtBlock - Execute eth_call RPC against obscuro for a specific block (batch) number.
	ObsCallAtBlock(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error)

	// GetChainStateAtTransaction - returns the stateDB after applying all the transactions in the batch leading to the desired transaction.
	GetChainStateAtTransaction(ctx context.Context, batch *core.Batch, txIndex int, reexec uint64) (*gethcore.Message, vm.BlockContext, *state.StateDB, error)
}
