//nolint:contextcheck
package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	// todo: this is a convenience for testnet testing and will eventually be retrieved from the L1
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// An in-memory implementation of `rpc.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI       *clientapi.ObscuroAPI
	ethAPI           *clientapi.EthereumAPI
	filterAPI        *clientapi.FilterAPI
	obscuroScanAPI   *clientapi.ObscuroScanAPI
	testAPI          *clientapi.TestAPI
	enclavePublicKey *ecies.PublicKey
}

func NewInMemObscuroClient(nodeHost hostcommon.Host) rpc.Client {
	logger := testlog.Logger().New(log.CmpKey, log.RPCClientCmp)
	// todo: this is a convenience for testnet but needs to replaced by a parameter and/or retrieved from the target host
	enclPubECDSA, err := crypto.DecompressPubkey(gethcommon.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}
	enclPubKey := ecies.ImportECDSAPublic(enclPubECDSA)

	return &inMemObscuroClient{
		obscuroAPI:       clientapi.NewObscuroAPI(nodeHost),
		ethAPI:           clientapi.NewEthereumAPI(nodeHost, logger),
		filterAPI:        clientapi.NewFilterAPI(nodeHost, logger),
		obscuroScanAPI:   clientapi.NewObscuroScanAPI(nodeHost),
		testAPI:          clientapi.NewTestAPI(nodeHost, nil),
		enclavePublicKey: enclPubKey,
	}
}

func NewInMemoryEncRPCClient(host hostcommon.Host, viewingKey *rpc.ViewingKey, logger gethlog.Logger) *rpc.EncRPCClient {
	inMemClient := NewInMemObscuroClient(host)
	encClient, err := rpc.NewEncRPCClient(inMemClient, viewingKey, logger)
	if err != nil {
		panic(err)
	}
	return encClient
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case rpc.SendRawTransaction:
		return c.sendRawTransaction(args)

	case rpc.GetTransactionByHash:
		return c.getTransactionByHash(result, args)

	case rpc.Call:
		return c.rpcCall(result, args)

	case rpc.GetTransactionCount:
		return c.getTransactionCount(result, args)

	case rpc.GetTransactionReceipt:
		return c.getTransactionReceipt(result, args)

	case rpc.RollupNumber:
		*result.(*hexutil.Uint64) = c.ethAPI.BlockNumber()
		return nil

	case rpc.StopHost:
		c.testAPI.StopHost()
		return nil

	case rpc.AddViewingKey:
		return c.addViewingKey(args)

	case rpc.GetLogs:
		return c.getLogs(result, args)

	case rpc.GetRollupByNumber:
		return c.getRollupByNumber(result, args)

	case rpc.GetRollupByHash:
		return c.getRollupByHash(result, args)

	case rpc.Health:
		return c.health(result)

	case rpc.GetTotalTxs:
		return c.getTotalTransactions(result)

	case rpc.GetLatestTxs:
		return c.getLatestTransactions(result, args)

	case rpc.GetBatchForTx:
		return c.getBatchForTx(result, args)

	case rpc.GetBatch:
		return c.getBatch(result, args)

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}
}

// CallContext not currently supported by in-memory obscuro client, the context will be ignored.
func (c *inMemObscuroClient) CallContext(_ context.Context, result interface{}, method string, args ...interface{}) error {
	return c.Call(result, method, args...)
}

func (c *inMemObscuroClient) Subscribe(context.Context, interface{}, string, interface{}, ...interface{}) (*gethrpc.ClientSubscription, error) {
	panic("not implemented")
}

func (c *inMemObscuroClient) sendRawTransaction(args []interface{}) error {
	encBytes, err := getEncryptedBytes(args, rpc.SendRawTransaction)
	if err != nil {
		return err
	}

	_, err = c.ethAPI.SendRawTransaction(context.Background(), encBytes)
	return err
}

func (c *inMemObscuroClient) getTransactionByHash(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionByHash)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionByHash(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionByHash, err)
	}

	// GetTransactionByHash returns string pointer, we want string
	if encryptedResponse != nil {
		*result.(*interface{}) = *encryptedResponse
	}
	return nil
}

func (c *inMemObscuroClient) rpcCall(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.Call)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.Call(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.Call, err)
	}
	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getTransactionReceipt(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionReceipt)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionReceipt(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionReceipt, err)
	}

	// GetTransactionReceipt returns string pointer, we want string
	if encryptedResponse != nil {
		*result.(*interface{}) = *encryptedResponse
	}
	return nil
}

func (c *inMemObscuroClient) getTransactionCount(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionCount)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.ethAPI.GetTransactionCount(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionCount, err)
	}

	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getLogs(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetLogs)
	if err != nil {
		return err
	}
	encryptedResponse, err := c.filterAPI.GetLogs(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetLogs, err)
	}
	*result.(*interface{}) = encryptedResponse
	return nil
}

func (c *inMemObscuroClient) getRollupByNumber(result interface{}, args []interface{}) error {
	blockNumberHex, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected int64", rpc.GetRollupByNumber, args[0])
	}

	blockNumber, err := hexutil.DecodeUint64(blockNumberHex)
	if err != nil {
		return fmt.Errorf("arg to %s could not be decoded from hex. Cause: %w", rpc.GetRollupByNumber, err)
	}

	headerMap, err := c.ethAPI.GetBlockByNumber(nil, gethrpc.BlockNumber(blockNumber), false) //nolint:staticcheck
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetRollupByNumber, err)
	}

	headerJSON, err := json.Marshal(headerMap)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to JSON. Cause: %w", rpc.GetRollupByNumber, err)
	}
	var header common.BatchHeader
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to rollup header. Cause: %w", rpc.GetRollupByNumber, err)
	}

	*result.(**common.BatchHeader) = &header
	return nil
}

func (c *inMemObscuroClient) getRollupByHash(result interface{}, args []interface{}) error {
	blockHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected common.Hash", rpc.GetRollupByHash, args[0])
	}

	headerMap, err := c.ethAPI.GetBlockByHash(nil, blockHash, false) //nolint:staticcheck
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetRollupByHash, err)
	}

	headerJSON, err := json.Marshal(headerMap)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to JSON. Cause: %w", rpc.GetRollupByHash, err)
	}
	var header common.BatchHeader
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return fmt.Errorf("could not marshal %s response to rollup header. Cause: %w", rpc.GetRollupByHash, err)
	}

	*result.(**common.BatchHeader) = &header
	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}

func (c *inMemObscuroClient) SetViewingKey(_ *ecies.PrivateKey, _ []byte) {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) RegisterViewingKey(_ gethcommon.Address, _ []byte) error {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (c *inMemObscuroClient) addViewingKey(args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 args to %s, got %d", rpc.AddViewingKey, len(args))
	}

	vk, ok := args[0].([]byte)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type []byte", rpc.AddViewingKey, args[0])
	}

	sig, ok := args[1].([]byte)
	if !ok {
		return fmt.Errorf("second arg to %s is of type %T, expected type []byte", rpc.AddViewingKey, args[1])
	}
	return c.obscuroAPI.AddViewingKey(vk, sig)
}

func (c *inMemObscuroClient) health(result interface{}) error {
	*result.(**hostcommon.HealthCheck) = &hostcommon.HealthCheck{OverallHealth: true}
	return nil
}

func (c *inMemObscuroClient) getTotalTransactions(result interface{}) error {
	totalTxs, err := c.obscuroScanAPI.GetTotalTransactions()
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTotalTxs, err)
	}

	*result.(**big.Int) = totalTxs
	return nil
}

func (c *inMemObscuroClient) getLatestTransactions(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetLatestTxs, len(args))
	}
	numTxs, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetLatestTxs, args[0])
	}

	latestTxs, err := c.obscuroScanAPI.GetLatestTransactions(numTxs)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetLatestTxs, err)
	}

	*result.(*[]gethcommon.Hash) = latestTxs
	return nil
}

func (c *inMemObscuroClient) getBatchForTx(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatchForTx, len(args))
	}
	txHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatchForTx, args[0])
	}

	batch, err := c.obscuroScanAPI.GetBatchForTx(txHash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchForTx, err)
	}

	*result.(**common.ExtBatch) = batch
	return nil
}

func (c *inMemObscuroClient) getBatch(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatch, len(args))
	}
	batchHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatch, args[0])
	}

	batch, err := c.obscuroScanAPI.GetBatch(batchHash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatch, err)
	}

	*result.(**common.ExtBatch) = batch
	return nil
}

// getEncryptedBytes expects args to have a single element and it to be of type bytes (client doesn't know anything about what's getting passed through on sensitive methods)
func getEncryptedBytes(args []interface{}, methodName string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arg to %s, got %d", methodName, len(args))
	}
	encBytes, ok := args[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("first arg to %s is of type %T, expected []byte", methodName, args[0])
	}
	return encBytes, nil
}
