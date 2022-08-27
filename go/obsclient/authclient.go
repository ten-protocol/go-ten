package obsclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

// AuthObsClient extends the functionality of the ObsClient for all methods that require encryption when communicating
//
//	with the enclave
type AuthObsClient struct {
	ObsClient
	c rpcclientlib.Client
}

// NewAuthObsClient constructs an AuthObsClient for sensitive communication with an enclave.
//
//	It requires an EncRPCClient specifically even though the AuthObsClient uses a Client interface in its struct because
//	the Client interface makes testing easy but an EncRPCClient is required for the actual encrypted communication
func NewAuthObsClient(client *rpcclientlib.EncRPCClient) *AuthObsClient {
	return &AuthObsClient{
		ObsClient: ObsClient{
			c: client,
		},
		c: client,
	}
}

func (ac *AuthObsClient) Close() {
	ac.c.Stop()
	ac.ObsClient.Close()
}

// TransactionByHash returns transaction (if found), isPending (always false currently as we don't search the mempool), error
func (ac *AuthObsClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	var tx types.Transaction
	err := ac.c.CallContext(ctx, &tx, rpcclientlib.RPCGetTransactionByHash, hash.Hex())
	// todo: revisit isPending result value, included for ethclient equivalence but hardcoded currently
	return &tx, false, err
}

func (ac *AuthObsClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt types.Receipt
	err := ac.c.CallContext(ctx, &receipt, rpcclientlib.RPCGetTxReceipt, txHash)
	return &receipt, err
}

func (ac *AuthObsClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ac.c.CallContext(ctx, &result, rpcclientlib.RPCNonce, account, toBlockNumArg(blockNumber))
	return uint64(result), err
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ac *AuthObsClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex string
	err := ac.c.CallContext(ctx, &hex, rpcclientlib.RPCCall, toCallArg(msg), toBlockNumArg(blockNumber))
	return []byte(hex), err
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ac *AuthObsClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	return ac.c.CallContext(ctx, nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
}
