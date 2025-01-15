package ethadapter

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthClient defines the interface for RPC communications with the ethereum nodes
// todo (#1617) - some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type EthClient interface {
	BlockNumber() (uint64, error)                                                 // retrieves the number of the head block
	FetchHeadBlock() (*types.Header, error)                                       // retrieves the block at head height
	HeaderByHash(id gethcommon.Hash) (*types.Header, error)                       // retrieves a block header given a hash
	BlockByHash(id gethcommon.Hash) (*types.Block, error)                         // retrieves a block given a hash
	HeaderByNumber(n *big.Int) (*types.Header, error)                             // retrieves a block given a number - returns head block if n is nil
	SendTransaction(signedTx *types.Transaction) error                            // issues an ethereum transaction (expects signed tx)
	TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error)              // fetches the ethereum transaction receipt
	TransactionByHash(hash gethcommon.Hash) (*types.Transaction, bool, error)     // fetches the ethereum tx
	Nonce(address gethcommon.Address) (uint64, error)                             // fetches the account nonce to use in the next transaction
	BalanceAt(account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) // fetches the balance of the account
	GetLogs(q ethereum.FilterQuery) ([]types.Log, error)                          // fetches the logs for a given query
	CallContract(msg ethereum.CallMsg) ([]byte, error)                            // Runs the provided call message on the latest block.
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)                       // Suggests the gas tip cap
	EthClient() *ethclient.Client                                                 // returns the underlying eth client
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)       // estimates the gas for a given call
	BlockListener() (chan *types.Header, ethereum.Subscription)                   // subscribes to new blocks and returns a listener with the blocks heads and the subscription handler
	ReconnectIfClosed() error                                                     // closes and creates a new connection
	Alive() bool                                                                  // returns whether the connection is live or not
	Info() Info                                                                   // retrieves the node Info
	Stop()                                                                        // tries to cleanly stop the client and release any resources

	// todo - all the below should be removed from the interface
	BlocksBetween(block *types.Header, head *types.Header) []*types.Header      // returns the blocks between two blocks
	IsBlockAncestor(block *types.Header, maybeAncestor common.L1BlockHash) bool // returns if the node considers a block the ancestor
	FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error)
}

// Info forces the RPC EthClient to return the data in the same format (independently of its implementation)
type Info struct {
	L2ID gethcommon.Address // the address of the Obscuro node this client is dedicated to
}
