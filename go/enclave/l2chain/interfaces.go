package l2chain

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// ObscuroChain - the interface that provides the data access layer to the obscuro l2.
// Operations here should be read only.
type ObscuroChain interface {
	// GetBalance - Returns the balance of an address given the standard provided BlockNumber which in our case is batch number.
	// Note this method also returns the address we should encrypt the balance for, unlike `GetBalanceAtBlock`
	GetBalance(accountAddress gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, *hexutil.Big, error)

	// GetBalanceAtBlock - will return the balance of a specific address at the specific given block number (batch number).
	GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error)

	// ObsCall - The interface for executing eth_call RPC commands against obscuro.
	ObsCall(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error)

	// ObsCallAtBlock - Execute eth_call RPC against obscuro for a specific block (batch) number.
	ObsCallAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error)

	// GetChainStateAtTransaction - returns the stateDB after applying all the transactions in the batch leading to the desired transaction.
	GetChainStateAtTransaction(batch *core.Batch, txIndex int, reexec uint64) (gethcore.Message, vm.BlockContext, *state.StateDB, error)
}
