package obsclient

import (
	"context"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

// AuthObsClient extends the functionality of the ObsClient for all methods that require encryption when communicating with the enclave
// It is created with an EncRPCClient rather than basic RPC client so encryption/decryption is supported
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type AuthObsClient struct {
	ObsClient
	Account *common.Address
}

// NewAuthObsClient constructs an AuthObsClient for sensitive communication with an enclave.
//
// It requires an EncRPCClient specifically even though the AuthObsClient uses a Client interface in its struct because
// the Client interface makes testing easy but an EncRPCClient is required for the actual encrypted communication
func NewAuthObsClient(client *rpc.EncRPCClient) *AuthObsClient {
	return &AuthObsClient{
		ObsClient: ObsClient{
			RPCClient: client,
		},
		Account: client.Account(),
	}
}

// DialWithAuth will generate and sign a viewing key for given wallet, then initiate a connection with the RPC node and
//
//	register the viewing key
func DialWithAuth(rpcurl string, wal wallet.Wallet) (*AuthObsClient, error) {
	viewingKey, err := rpc.GenerateAndSignViewingKey(wal)
	if err != nil {
		return nil, err
	}
	encClient, err := rpc.NewEncNetworkClient(rpcurl, viewingKey)
	if err != nil {
		return nil, err
	}
	return NewAuthObsClient(encClient), nil
}

// TransactionByHash returns transaction (if found), isPending (always false currently as we don't search the mempool), error
func (ac *AuthObsClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	var tx types.Transaction
	err := ac.RPCClient.CallContext(ctx, &tx, rpc.RPCGetTransactionByHash, hash.Hex())
	// todo: revisit isPending result value, included for ethclient equivalence but hardcoded currently
	return &tx, false, err
}

func (ac *AuthObsClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt types.Receipt
	err := ac.RPCClient.CallContext(ctx, &receipt, rpc.RPCGetTxReceipt, txHash)
	return &receipt, err
}

// NonceAt retrieves the nonce for the account registered on this client (due to obscuro privacy restrictions,
// nonce cannot be requested for other accounts)
func (ac *AuthObsClient) NonceAt(ctx context.Context, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ac.RPCClient.CallContext(ctx, &result, rpc.RPCNonce, ac.Account.Hex(), toBlockNumArg(blockNumber))
	return uint64(result), err
}

func (ac *AuthObsClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex string
	err := ac.RPCClient.CallContext(ctx, &hex, rpc.RPCCall, toCallArg(msg), toBlockNumArg(blockNumber))
	return []byte(hex), err
}

func (ac *AuthObsClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	return ac.RPCClient.CallContext(ctx, nil, rpc.RPCSendRawTransaction, encodeTx(signedTx))
}

// BalanceAt retrieves the native balance for the account registered on this client (due to obscuro privacy restrictions,
// balance cannot be requested for other accounts)
func (ac *AuthObsClient) BalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	var result string
	err := ac.RPCClient.CallContext(ctx, &result, rpc.RPCGetBalance, ac.Account.Hex(), toBlockNumArg(blockNumber))
	if err != nil {
		return big.NewInt(0), err
	}
	balance, err := hexutil.DecodeBig(result)
	return balance, err
}
