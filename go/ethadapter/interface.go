package ethadapter

import (
	"errors"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrSubscriptionNotSupported return from BlockListener subscription if client doesn't support streaming (in-mem simulation)
var ErrSubscriptionNotSupported = errors.New("block subscription not supported")

// EthClient defines the interface for RPC communications with the ethereum nodes
// TODO Some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type EthClient interface {
	BlockNumber() (uint64, error)                                                 // retrieves the number of the head block
	BlockByHash(id gethcommon.Hash) (*types.Block, error)                         // retrieves a block given a hash
	BlockByNumber(n *big.Int) (*types.Block, error)                               // retrieves a block given a number - returns head block if n is nil
	SendTransaction(signedTx *types.Transaction) error                            // issues an ethereum transaction (expects signed tx)
	TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error)              // fetches the ethereum transaction receipt
	Nonce(address gethcommon.Address) (uint64, error)                             // fetches the account nonce to use in the next transaction
	BalanceAt(account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) // fetches the balance of the account

	Info() Info                                                         // retrieves the node Info
	FetchHeadBlock() (*types.Block, error)                              // retrieves the block at head height
	BlocksBetween(block *types.Block, head *types.Block) []*types.Block // returns the blocks between two blocks
	IsBlockAncestor(block *types.Block, proof common.L1RootHash) bool   // returns if the node considers a block the ancestor
	BlockListener() (chan *types.Header, ethereum.Subscription)         // subscribes to new blocks and returns a listener with the blocks heads and the subscription handler

	CallContract(msg ethereum.CallMsg) ([]byte, error) // Runs the provided call message on the latest block.

	Stop() // tries to cleanly stop the client and release any resources

	EthClient() *ethclient.Client // returns the underlying eth client
}

// Info forces the RPC EthClient to return the data in the same format (independently of its implementation)
type Info struct {
	L2ID gethcommon.Address // the address of the Obscuro node this client is dedicated to
}
