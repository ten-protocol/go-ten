package obsclient

import (
	"context"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

// AuthObsClient extends the functionality of the ObsClient for all methods that require encryption when communicating with the enclave
// It is created with an EncRPCClient rather than basic RPC client so encryption/decryption is supported
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type AuthObsClient struct {
	ObsClient
}

// NewAuthObsClient constructs an AuthObsClient for sensitive communication with an enclave.
//
// It requires an EncRPCClient specifically even though the AuthObsClient uses a Client interface in its struct because
// the Client interface makes testing easy but an EncRPCClient is required for the actual encrypted communication
func NewAuthObsClient(client *rpcclientlib.EncRPCClient) *AuthObsClient {
	return &AuthObsClient{
		ObsClient: ObsClient{
			RPCClient: client,
		},
	}
}

// DialWithAuth will generate and sign a viewing key for given wallet, then initiate a connection with the RPC node and
//   register the viewing key
func DialWithAuth(rpcurl string, wal wallet.Wallet) (*AuthObsClient, error) {
	viewingKey, err := rpcclientlib.GenerateAndSignViewingKey(wal)
	if err != nil {
		return nil, err
	}
	encClient, err := rpcclientlib.NewEncNetworkClient(rpcurl, viewingKey)
	if err != nil {
		return nil, err
	}
	return NewAuthObsClient(encClient), nil
}

// TransactionByHash returns transaction (if found), isPending (always false currently as we don't search the mempool), error
func (ac *AuthObsClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	var tx types.Transaction
	err := ac.RPCClient.CallContext(ctx, &tx, rpcclientlib.RPCGetTransactionByHash, hash.Hex())
	// todo: revisit isPending result value, included for ethclient equivalence but hardcoded currently
	return &tx, false, err
}

func (ac *AuthObsClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt types.Receipt
	err := ac.RPCClient.CallContext(ctx, &receipt, rpcclientlib.RPCGetTxReceipt, txHash)
	return &receipt, err
}

func (ac *AuthObsClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ac.RPCClient.CallContext(ctx, &result, rpcclientlib.RPCNonce, account, toBlockNumArg(blockNumber))
	return uint64(result), err
}

func (ac *AuthObsClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex string
	err := ac.RPCClient.CallContext(ctx, &hex, rpcclientlib.RPCCall, toCallArg(msg), toBlockNumArg(blockNumber))
	return []byte(hex), err
}

func (ac *AuthObsClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	return ac.RPCClient.CallContext(ctx, nil, rpcclientlib.RPCSendRawTransaction, encodeTx(signedTx))
}
