package obsclient

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"math/big"
)

// AuthObsClient extends the functionality of the ObsClient for all methods that require encryption when communicating with
//	the enclave
type AuthObsClient struct {
	ObsClient
	c rpcclientlib.Client
}

// NewAuthObsClient constructs an AuthObsClient for sensitive communication with an enclave.
// 	It requires an EncRPCClient specifically even though the AuthObsClient uses a Client interface in its struct because
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

//func Dial(rawurl string) (*Client, error) {
//func DialContext(ctx context.Context, rawurl string) (*Client, error) {
//func NewClient(c *rpc.Client) *Client {
//func (ec *Client) Close() {
//func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
//func (ec *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
//func (ec *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
//func (ec *Client) BlockNumber(ctx context.Context) (uint64, error) {
//func (ec *Client) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
//func (ec *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
//func (ec *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
//func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
//func (ec *Client) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
//func (ec *Client) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
//func (ec *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
//func (ec *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
//func (ec *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
//func (ec *Client) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
//func (ec *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
//func (ec *Client) NetworkID(ctx context.Context) (*big.Int, error) {
//func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
//func (ec *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
//func (ec *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
//func (ec *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
//func (ec *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
//func (ec *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
//func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
//func (ec *Client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
//func (ec *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
//func (ec *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
//func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
//func (ec *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
//func (ec *Client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
//func (ec *Client) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
//func (ec *Client) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
//func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
//func (ec *Client) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
//func (ec *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
//func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
//func toBlockNumArg(number *big.Int) string {
//func toCallArg(msg ethereum.CallMsg) interface{} {
//func (p *rpcProgress) toSyncProgress() *ethereum.SyncProgress {
