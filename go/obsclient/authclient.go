package obsclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
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
//
//	register the viewing key
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

// TransactionByHash returns transaction (if found), isPending (always false currently as we don't search the mempool), error
func (ac *AuthObsClient) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (*types.Transaction, bool, error) {
	var tx responses.TxType
	err := ac.rpcClient.CallContext(ctx, &tx, rpc.GetTransactionByHash, hash.Hex())
	if err != nil {
		return nil, false, err
	}
	// todo (#1491) - revisit isPending result value, included for ethclient equivalence but hardcoded currently
	return &tx, false, nil
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
	var result responses.ReceiptType
	var emptyHash gethcommon.Hash
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetTransactionReceipt, txHash)
	if err != nil {
		return nil, err
	}
	if result.TxHash != emptyHash {
		return &result, nil
	}

	return nil, nil
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
	var result responses.CallType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.Call, ToCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func (ac *AuthObsClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	var result responses.RawTxType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.SendRawTransaction, encodeTx(signedTx))
	if err != nil {
		return err
	}
	println(result.Hex())
	return nil
}

// BalanceAt retrieves the native balance for the account registered on this client (due to obscuro privacy restrictions,
// balance cannot be requested for other accounts)
func (ac *AuthObsClient) BalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	var result responses.BalanceType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetBalance, ac.account, toBlockNumArg(blockNumber))
	if err != nil {
		return big.NewInt(0), err
	}

	return result.ToInt(), nil
}

func (ac *AuthObsClient) SubscribeFilterLogs(ctx context.Context, filterCriteria filters.FilterCriteria, ch chan common.IDAndLog) (ethereum.Subscription, error) {
	return ac.rpcClient.Subscribe(ctx, nil, rpc.SubscribeNamespace, ch, rpc.SubscriptionTypeLogs, filterCriteria)
}

func (ac *AuthObsClient) GetLogs(ctx context.Context, filterCriteria common.FilterCriteriaJSON) ([]*types.Log, error) {
	var result responses.LogsType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetLogs, filterCriteria, ac.account)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ac *AuthObsClient) Address() gethcommon.Address {
	return ac.account
}

func (ac *AuthObsClient) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (uint64, error) {
	var result responses.GasType
	err := ac.rpcClient.CallContext(ctx, &result, rpc.EstimateGas, ToCallArg(*msg))
	if err != nil {
		return 0, err
	}

	return hexutil.DecodeUint64(result.String())
}

func (ac *AuthObsClient) EstimateGasAndGasPrice(txData types.TxData) types.TxData {
	unEstimatedTx := types.NewTx(txData)

	gasLimit, err := ac.EstimateGas(context.Background(), &ethereum.CallMsg{
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

// GetReceiptsByAddress retrieves the receipts for the account registered on this client (due to obscuro privacy restrictions,
// balance cannot be requested for other accounts)
func (ac *AuthObsClient) GetReceiptsByAddress(ctx context.Context, address *gethcommon.Address) (types.Receipts, error) {
	var result types.Receipts
	err := ac.rpcClient.CallContext(ctx, &result, rpc.GetStorageAt, address, nil, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
