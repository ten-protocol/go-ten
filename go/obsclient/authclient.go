package obsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// AuthObsClient extends the functionality of the ObsClient for all methods that require encryption when communicating with the enclave
// It is created with an EncRPCClient rather than basic RPC client so encryption/decryption is supported
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type AuthObsClient struct {
	ObsClient
	account gethcommon.Address
}

// NewAuthObsClient constructs an AuthObsClient for sensitive communication with an enclave.
//
// It requires an EncRPCClient specifically even though the AuthObsClient uses a Client interface in its struct because
// the Client interface makes testing easy but an EncRPCClient is required for the actual encrypted communication
func NewAuthObsClient(client *rpc.EncRPCClient) *AuthObsClient {
	return &AuthObsClient{
		ObsClient: ObsClient{
			rpcClient: client,
		},
		account: *client.Account(),
	}
}

// DialWithAuth will generate and sign a viewing key for given wallet, then initiate a connection with the RPC node and
// register the viewing key
func DialWithAuth(rpcurl string, wal wallet.Wallet, logger gethlog.Logger) (*AuthObsClient, error) {
	viewingKey, err := viewingkey.GenerateViewingKeyForWallet(wal)
	if err != nil {
		return nil, err
	}
	encClient, err := rpc.NewEncNetworkClient(rpcurl, viewingKey, logger)
	if err != nil {
		return nil, err
	}

	authClient := NewAuthObsClient(encClient)
	return authClient, nil
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string             `json:"blockNumber,omitempty"`
	BlockHash   *gethcommon.Hash    `json:"blockHash,omitempty"`
	From        *gethcommon.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

// TransactionByHash returns transaction (if found), isPending (always false currently as we don't search the mempool), error
func (ac *AuthObsClient) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (tx *types.Transaction, isPending bool, err error) {
	var result *rpcTransaction
	err = ac.rpcClient.CallContext(ctx, &result, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if result == nil {
		return nil, false, ethereum.NotFound
	} else if _, r, _ := result.tx.RawSignatureValues(); r == nil {
		return nil, false, errors.New("server returned transaction without signature")
	}
	if result.From != nil && result.BlockHash != nil {
		setSenderFromServer(result.tx, *result.From, *result.BlockHash)
	}
	return result.tx, result.BlockNumber == nil, nil
}

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer struct {
	addr      gethcommon.Address
	blockhash gethcommon.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer(tx *types.Transaction, addr gethcommon.Address, block gethcommon.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	_, _ = types.Sender(&senderFromServer{addr, block}, tx)
}

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(_ *types.Transaction) (gethcommon.Address, error) {
	if s.addr == (gethcommon.Address{}) {
		return gethcommon.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) ChainID() *big.Int {
	panic("can't sign with senderFromServer")
}

func (s *senderFromServer) Hash(_ *types.Transaction) gethcommon.Hash {
	panic("can't sign with senderFromServer")
}

func (s *senderFromServer) SignatureValues(_ *types.Transaction, _ []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}

func (ac *AuthObsClient) GasPrice(ctx context.Context) (*big.Int, error) {
	var result responses.GasPriceType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GasPrice)
	if err != nil {
		return nil, err
	}

	return result.ToInt(), nil
}

func (ac *AuthObsClient) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	var r *types.Receipt
	err := ac.rpcClient.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

// NonceAt retrieves the nonce for the account registered on this client (due to obscuro privacy restrictions,
// nonce cannot be requested for other accounts)
func (ac *AuthObsClient) NonceAt(ctx context.Context, blockNumber *big.Int) (uint64, error) {
	var result responses.NonceType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetTransactionCount, ac.account, toBlockNumArg(blockNumber))
	if err != nil {
		return 0, err
	}

	return hexutil.DecodeUint64(result)
}

func (ac *AuthObsClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ac.rpcClient.CallContext(ctx, &hex, "eth_call", ToCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (ac *AuthObsClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	var result responses.RawTxType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.SendRawTransaction, encodeTx(signedTx))
	if err != nil {
		return err
	}
	return nil
}

// BalanceAt retrieves the native balance for the account registered on this client (due to obscuro privacy restrictions,
// balance cannot be requested for other accounts)
func (ac *AuthObsClient) BalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := ac.rpcClient.CallContext(ctx, &result, "eth_getBalance", ac.account, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}

func (ac *AuthObsClient) SubscribeFilterLogs(ctx context.Context, filterCriteria common.FilterCriteria, ch chan types.Log) (ethereum.Subscription, error) {
	return ac.rpcClient.Subscribe(ctx, rpc.SubscribeNamespace, ch, rpc.SubscriptionTypeLogs, filterCriteria)
}

func (ac *AuthObsClient) GetLogs(ctx context.Context, filterCriteria common.FilterCriteria) ([]*types.Log, error) {
	var result responses.LogsType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetLogs, filterCriteria)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ac *AuthObsClient) Address() gethcommon.Address {
	return ac.account
}

func (ac *AuthObsClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ac.rpcClient.CallContext(ctx, &hex, "eth_estimateGas", ToCallArg(msg))
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

func (ac *AuthObsClient) EstimateGasAndGasPrice(txData types.TxData) types.TxData {
	unEstimatedTx := types.NewTx(txData)

	gasLimit, err := ac.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  ac.Address(),
		To:    unEstimatedTx.To(),
		Value: unEstimatedTx.Value(),
		Data:  unEstimatedTx.Data(),
	})
	if err != nil {
		gasLimit = unEstimatedTx.Gas()
	}

	gasPrice, err := ac.GasPrice(context.Background())
	if err != nil {
		// params.InitialBaseFee should be the new standard gas price.
		// If the gas price is too low, then the gas required to be put in a transaction
		// becomes astronomical.
		gasPrice = big.NewInt(params.InitialBaseFee)
	}

	return &types.LegacyTx{
		Nonce:    unEstimatedTx.Nonce(),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       unEstimatedTx.To(),
		Value:    unEstimatedTx.Value(),
		Data:     unEstimatedTx.Data(),
	}
}

// GetPrivateTransactions retrieves the receipts for the specified account (must be registered on this client), returns requested range of receipts and the total number of receipts for that acc
func (ac *AuthObsClient) GetPrivateTransactions(ctx context.Context, address *gethcommon.Address, pagination common.QueryPagination) (types.Receipts, uint64, error) {
	queryParam := &common.ListPrivateTransactionsQueryParams{
		Address:    *address,
		Pagination: pagination,
	}
	queryParamStr, err := json.Marshal(queryParam)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to marshal query params - %w", err)
	}
	var result common.PrivateTransactionsQueryResponse
	// todo @matt don't use the string method name here and avoid stringifying the params
	err = ac.rpcClient.CallContext(ctx, &result, rpc.GetPersonalTransactions, common.ListPrivateTransactionsCQMethod, string(queryParamStr), nil)
	if err != nil {
		return nil, 0, err
	}

	return result.Receipts, result.Total, nil
}
