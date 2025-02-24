package l2chain

import (
	"context"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// ObscuroChain - the interface that provides the data access layer to the obscuro l2.
// Operations here should be read only.
type ObscuroChain interface {
	// GetBalanceAtBlock - will return the balance of a specific address at the specific given block number (batch number).
	GetBalanceAtBlock(ctx context.Context, accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error)

	// Call - The interface for executing eth_call RPC commands against obscuro.
	Call(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error, common.SystemError)

	// ObsCallAtBlock - Execute eth_call RPC against obscuro for a specific block (batch) number.
	ObsCallAtBlock(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error, common.SystemError)
}
