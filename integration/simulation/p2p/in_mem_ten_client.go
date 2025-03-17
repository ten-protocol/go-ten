package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/host/container"
	"github.com/ten-protocol/go-ten/go/host/rpc/clientapi"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration/common/testlog"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// todo - move this from the P2P folder
// An in-memory implementation of `rpc.Client` that speaks directly to the node.
type inMemTenClient struct {
	tenAPI     *clientapi.TenAPI
	ethAPI     *clientapi.ChainAPI
	filterAPI  *clientapi.FilterAPI
	tenScanAPI *clientapi.ScanAPI
	testAPI    *clientapi.TestAPI
}

func NewInMemTenClient(hostContainer *container.HostContainer) rpc.Client {
	logger := testlog.Logger().New(log.CmpKey, log.RPCClientCmp)
	return &inMemTenClient{
		tenAPI:     clientapi.NewTenAPI(hostContainer.Host(), logger),
		ethAPI:     clientapi.NewChainAPI(hostContainer.Host(), logger),
		filterAPI:  clientapi.NewFilterAPI(hostContainer.Host(), logger),
		tenScanAPI: clientapi.NewScanAPI(hostContainer.Host(), logger),
		testAPI:    clientapi.NewTestAPI(hostContainer),
	}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemTenClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case tenrpc.EncRPC:
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg for encrypted method, got %d", len(args))
		}
		encryptedRPCRequest, ok := args[0].(common.EncryptedRPCRequest)
		if !ok {
			return fmt.Errorf("first arg to encrypted method is of type %T, expected EncryptedRPCRequest", args[0])
		}

		encryptedResponse, err := c.tenAPI.EncryptedRPC(context.Background(), encryptedRPCRequest)
		if err == nil {
			*result.(*responses.EnclaveResponse) = encryptedResponse
		}
		return err

	case rpc.BatchNumber:
		*result.(*hexutil.Uint64) = c.ethAPI.BatchNumber()
		return nil

	case rpc.StopHost:
		return c.testAPI.StopHost()

	case rpc.GetBatchByNumber:
		return c.getBatchByNumber(result, args)

	case rpc.GetBatchByHash:
		return c.getBatchByHash(result, args)

	case rpc.Health:
		return c.health(result)

	case rpc.GetTotalTxCount:
		return c.getTotalTransactions(result)

	case rpc.GetBatchByTx:
		return c.getBatchByTx(result, args)

	case rpc.GetBatch:
		return c.getBatch(result, args)

	case rpc.GetBatchListing:
		return c.getBatchListingDeprecated(result, args)

	case rpc.GetRollupListing:
		return c.getRollupListing(result, args)

	case rpc.GetPublicTransactionData:
		return c.getPublicTransactionData(result, args)

	case rpc.GasPrice:
		return c.getGasPrice(result)

	case rpc.Config:
		return c.tenConfig(result)

	case rpc.RPCKey:
		key, err := c.tenAPI.RpcKey()
		*result.(*[]byte) = key
		return err

	case rpc.GetCode:
		return c.getCode(result, args)

	case rpc.GetCrossChainProof:
		return c.getCrossChainProof(result, args)

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}
}

func (c *inMemTenClient) getCrossChainProof(result interface{}, args []interface{}) error {
	messageType, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("invalid argument type: expected string")
	}
	crossChainMessage, ok := args[1].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("invalid argument type: expected gethcommon.Hash")
	}
	proof, err := c.tenAPI.GetCrossChainProof(context.Background(), messageType, crossChainMessage)
	if err != nil {
		return fmt.Errorf("failed to get cross chain proof: %w", err)
	}
	*result.(*clientapi.CrossChainProof) = proof
	return nil
}

func (c *inMemTenClient) getCode(result interface{}, args []interface{}) error {
	address, ok := args[0].(gethcommon.Address)
	if !ok {
		return fmt.Errorf("invalid argument type: expected gethcommon.Address")
	}

	blockNumber, err := gethencoding.ExtractOptionalBlockNumber(args, 1)
	if err != nil {
		return fmt.Errorf("arg to %s is of type %T, expected int64", rpc.GetCode, args[0])
	}

	code, err := c.ethAPI.GetCode(context.Background(), address, *blockNumber)
	if err != nil {
		return err
	}

	*result.(*hexutil.Bytes) = code
	return nil
}

// CallContext not currently supported by in-memory TEN client, the context will be ignored.
func (c *inMemTenClient) CallContext(_ context.Context, result interface{}, method string, args ...interface{}) error {
	return c.Call(result, method, args...) //nolint: contextcheck
}

func (c *inMemTenClient) tenConfig(result interface{}) error {
	cfg, err := c.tenAPI.Config()
	if err != nil {
		return err
	}

	importantContracts := cfg.ImportantContracts

	publicSystemContracts := make(map[string]gethcommon.Address)
	for key, value := range cfg.PublicSystemContracts {
		publicSystemContracts[key] = gethcommon.Address(value)
	}

	tenNetworkInfo := &common.TenNetworkInfo{
		NetworkConfigAddress:            gethcommon.Address(cfg.NetworkConfigAddress),
		L2MessageBusAddress:             gethcommon.Address(cfg.L2MessageBusAddress),
		TransactionPostProcessorAddress: gethcommon.Address(cfg.TransactionPostProcessorAddress),
		L1StartHash:                     cfg.L1StartHash,
		ImportantContracts:              &importantContracts,
		PublicSystemContracts:           publicSystemContracts,
	}

	*result.(*common.TenNetworkInfo) = *tenNetworkInfo
	return nil
}

func (c *inMemTenClient) Subscribe(context.Context, string, interface{}, ...interface{}) (*gethrpc.ClientSubscription, error) {
	panic("not implemented")
}

func (c *inMemTenClient) getBatchByNumber(result interface{}, args []interface{}) error {
	blockNumberHex, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected int64", rpc.GetBatchByNumber, args[0])
	}

	blockNumber, err := hexutil.DecodeUint64(blockNumberHex)
	if err != nil {
		return fmt.Errorf("arg to %s could not be decoded from hex. Cause: %w", rpc.GetBatchByNumber, err)
	}

	headerMap, err := c.ethAPI.GetBatchByNumber(nil, gethrpc.BlockNumber(blockNumber), false) //nolint:staticcheck
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchByNumber, err)
	}

	headerJSON, err := json.Marshal(headerMap)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to JSON. Cause: %w", rpc.GetBatchByNumber, err)
	}
	var header common.BatchHeader
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to rollup header. Cause: %w", rpc.GetBatchByNumber, err)
	}

	*result.(**common.BatchHeader) = &header
	return nil
}

func (c *inMemTenClient) getBatchByHash(result interface{}, args []interface{}) error {
	blockHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected common.Hash", rpc.GetBatchByHash, args[0])
	}

	headerMap, err := c.ethAPI.GetBatchByHash(nil, blockHash, false) //nolint:staticcheck
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchByHash, err)
	}

	headerJSON, err := json.Marshal(headerMap)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to JSON. Cause: %w", rpc.GetBatchByHash, err)
	}
	var header common.BatchHeader
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to rollup header. Cause: %w", rpc.GetBatchByHash, err)
	}

	*result.(**common.BatchHeader) = &header
	return nil
}

func (c *inMemTenClient) Stop() {
	// There is no RPC connection to close.
}

func (c *inMemTenClient) SetViewingKey(_ *ecies.PrivateKey, _ []byte) {
	panic("viewing key encryption/decryption is not currently supported by in-memory ten-client")
}

func (c *inMemTenClient) RegisterViewingKey(_ gethcommon.Address, _ []byte) error {
	panic("viewing key encryption/decryption is not currently supported by in-memory ten-client")
}

func (c *inMemTenClient) health(result interface{}) error {
	resPtr, ok := result.(*hostcommon.HealthCheck)
	if !ok {
		return fmt.Errorf("invalid type for result: expected *hostcommon.HealthCheck")
	}

	healthPtr, err := c.tenAPI.Health(context.Background())
	if err != nil {
		return err
	}

	*resPtr = *healthPtr // Dereference the pointer to assign the value
	return nil
}

func (c *inMemTenClient) getTotalTransactions(result interface{}) error {
	totalTxs, err := c.tenScanAPI.GetTotalTransactionCount()
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTotalTxCount, err)
	}

	*result.(**big.Int) = totalTxs
	return nil
}

func (c *inMemTenClient) getBatchByTx(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatchByTx, len(args))
	}
	txHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatchByTx, args[0])
	}

	batch, err := c.tenScanAPI.GetBatchByTx(txHash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchByTx, err)
	}

	*result.(**common.ExtBatch) = batch
	return nil
}

func (c *inMemTenClient) getBatch(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatch, len(args))
	}
	batchHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatch, args[0])
	}

	batch, err := c.tenScanAPI.GetBatch(batchHash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatch, err)
	}

	*result.(**common.ExtBatch) = batch
	return nil
}

func (c *inMemTenClient) getBatchListingDeprecated(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatchListing, len(args))
	}
	pagination, ok := args[0].(*common.QueryPagination)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatchListing, args[0])
	}

	batches, err := c.tenScanAPI.GetBatchListing(pagination)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchListing, err)
	}

	res, ok := result.(*common.BatchListingResponseDeprecated)
	if !ok {
		return fmt.Errorf("result is of type %T, expected *common.BatchListingResponseDeprecated", result)
	}
	*res = *batches
	return nil
}

func (c *inMemTenClient) getRollupListing(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetRollupListing, len(args))
	}
	pagination, ok := args[0].(*common.QueryPagination)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetRollupListing, args[0])
	}

	rollups, err := c.tenScanAPI.GetRollupListing(pagination)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetRollupListing, err)
	}

	res, ok := result.(*common.RollupListingResponse)
	if !ok {
		return fmt.Errorf("result is of type %T, expected *common.BatchListingResponseDeprecated", result)
	}
	*res = *rollups
	return nil
}

func (c *inMemTenClient) getGasPrice(result interface{}) error {
	gasPrice, err := c.ethAPI.GasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GasPrice, err)
	}

	*result.(*hexutil.Big) = *gasPrice
	return nil
}

func (c *inMemTenClient) getPublicTransactionData(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetPublicTransactionData, len(args))
	}
	pagination, ok := args[0].(*common.QueryPagination)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetPublicTransactionData, args[0])
	}

	txs, err := c.tenScanAPI.GetPublicTransactionData(pagination)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetPublicTransactionData, err)
	}

	res, ok := result.(*common.TransactionListingResponse)
	if !ok {
		return fmt.Errorf("result is of type %T, expected *common.BatchListingResponseDeprecated", result)
	}
	*res = *txs
	return nil
}
