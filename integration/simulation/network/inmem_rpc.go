package network

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"
	"github.com/obscuronet/go-obscuro/go/responses"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

// the MockServer has the same 'Call' interface as an RPC client, but it
// calls methods directly on the host (no network calls)
type MockServer struct {
	obscuroAPI     *clientapi.ObscuroAPI
	ethAPI         *clientapi.EthereumAPI
	filterAPI      *clientapi.FilterAPI
	obscuroScanAPI *clientapi.ObscuroScanAPI
	testAPI        *clientapi.TestAPI
}

func SetupMockServer(hostStop func() error) (host.ServiceFactory[host.RPCServerService], *MockServer) {
	server := &MockServer{}
	factory := func(config *config.HostConfig, serviceLocator host.ServiceLocator, logger gethlog.Logger) (host.RPCServerService, error) {
		server.obscuroAPI = clientapi.NewObscuroAPI(serviceLocator)
		server.ethAPI = clientapi.NewEthereumAPI(config.ObscuroChainID, serviceLocator, logger.New(log.CmpKey, "eth-rpc"))
		server.obscuroScanAPI = clientapi.NewObscuroScanAPI(serviceLocator)
		server.testAPI = clientapi.NewTestAPI(hostStop)
		server.filterAPI = clientapi.NewFilterAPI(serviceLocator, logger.New(log.CmpKey, "filter-rpc"))

		return server, nil
	}
	return factory, server
}

// Call bypasses RPC, and invokes methods on the node directly.
func (m *MockServer) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case rpc.SendRawTransaction:
		return m.sendRawTransaction(result, args)

	case rpc.GetTransactionByHash:
		return m.getTransactionByHash(result, args)

	case rpc.Call:
		return m.rpcCall(result, args)

	case rpc.GetTransactionCount:
		return m.getTransactionCount(result, args)

	case rpc.GetTransactionReceipt:
		return m.getTransactionReceipt(result, args)

	case rpc.BatchNumber:
		*result.(*hexutil.Uint64) = m.ethAPI.BlockNumber()
		return nil

	case rpc.StopHost:
		return m.testAPI.StopHost()

	case rpc.GetLogs:
		return m.getLogs(result, args)

	case rpc.GetBatchByNumber:
		return m.getBatchByNumber(result, args)

	case rpc.GetBatchByHash:
		return m.getBatchByHash(result, args)

	case rpc.Health:
		return m.health(result)

	case rpc.GetTotalTxs:
		return m.getTotalTransactions(result)

	case rpc.GetLatestTxs:
		return m.getLatestTransactions(result, args)

	case rpc.GetBatchForTx:
		return m.getBatchForTx(result, args)

	case rpc.GetBatch:
		return m.getBatch(result, args)

	default:
		return fmt.Errorf("RPC method %s is unknown", method)
	}
}

// CallContext not currently supported by in-memory obscuro client, the context will be ignored.
func (m *MockServer) CallContext(_ context.Context, result interface{}, method string, args ...interface{}) error {
	return m.Call(result, method, args...) //nolint: contextcheck
}

func (m *MockServer) Subscribe(context.Context, interface{}, string, interface{}, ...interface{}) (*gethrpc.ClientSubscription, error) {
	panic("not implemented")
}

func (m *MockServer) sendRawTransaction(result interface{}, args []interface{}) error {
	encBytes, err := getEncryptedBytes(args, rpc.SendRawTransaction)
	if err != nil {
		return err
	}

	encryptedResponse, err := m.ethAPI.SendRawTransaction(context.Background(), encBytes)
	if err == nil {
		*result.(*responses.EnclaveResponse) = encryptedResponse
	}

	return err
}

func (m *MockServer) getTransactionByHash(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionByHash)
	if err != nil {
		return err
	}
	encryptedResponse, err := m.ethAPI.GetTransactionByHash(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionByHash, err)
	}

	// GetTransactionByHash returns EnclaveResponse
	*result.(*responses.EnclaveResponse) = encryptedResponse
	return nil
}

func (m *MockServer) rpcCall(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.Call)
	if err != nil {
		return err
	}
	encryptedResponse, err := m.ethAPI.Call(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.Call, err)
	}
	*result.(*responses.EnclaveResponse) = encryptedResponse
	return nil
}

func (m *MockServer) getTransactionReceipt(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionReceipt)
	if err != nil {
		return err
	}
	encryptedResponse, err := m.ethAPI.GetTransactionReceipt(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionReceipt, err)
	}

	// GetTransactionReceipt returns EnclaveResponse
	*result.(*responses.EnclaveResponse) = encryptedResponse
	return nil
}

func (m *MockServer) getTransactionCount(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetTransactionCount)
	if err != nil {
		return err
	}
	encryptedResponse, err := m.ethAPI.GetTransactionCount(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTransactionCount, err)
	}

	*result.(*responses.EnclaveResponse) = encryptedResponse
	return nil
}

func (m *MockServer) getLogs(result interface{}, args []interface{}) error {
	enc, err := getEncryptedBytes(args, rpc.GetLogs)
	if err != nil {
		return err
	}
	encryptedResponse, err := m.filterAPI.GetLogs(context.Background(), enc)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetLogs, err)
	}
	*result.(*responses.EnclaveResponse) = encryptedResponse
	return nil
}

func (m *MockServer) getBatchByNumber(result interface{}, args []interface{}) error {
	blockNumberHex, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected int64", rpc.GetBatchByNumber, args[0])
	}

	blockNumber, err := hexutil.DecodeUint64(blockNumberHex)
	if err != nil {
		return fmt.Errorf("arg to %s could not be decoded from hex. Cause: %w", rpc.GetBatchByNumber, err)
	}

	headerMap, err := m.ethAPI.GetBlockByNumber(nil, gethrpc.BlockNumber(blockNumber), false) //nolint:staticcheck
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

func (m *MockServer) getBatchByHash(result interface{}, args []interface{}) error {
	blockHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("arg to %s is of type %T, expected common.Hash", rpc.GetBatchByHash, args[0])
	}

	headerMap, err := m.ethAPI.GetBlockByHash(nil, blockHash, false) //nolint:staticcheck
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

func (m *MockServer) Start() error {
	// There is no RPC connection to open.
	return nil
}

func (m *MockServer) Stop() {
	// There is no RPC connection to close.
}

func (m *MockServer) HealthStatus() hostcommon.HealthStatus {
	// always healthy, this is to satisfy the Service interface
	return &hostcommon.BasicErrHealthStatus{ErrMsg: ""}
}

func (m *MockServer) SetViewingKey(_ *ecies.PrivateKey, _ []byte) {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (m *MockServer) RegisterViewingKey(_ gethcommon.Address, _ []byte) error {
	panic("viewing key encryption/decryption is not currently supported by in-memory obscuro-client")
}

func (m *MockServer) health(result interface{}) error {
	*result.(**hostcommon.HealthCheck) = &hostcommon.HealthCheck{OverallHealth: true}
	return nil
}

func (m *MockServer) getTotalTransactions(result interface{}) error {
	totalTxs, err := m.obscuroScanAPI.GetTotalTransactions()
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetTotalTxs, err)
	}

	*result.(**big.Int) = totalTxs
	return nil
}

func (m *MockServer) getLatestTransactions(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetLatestTxs, len(args))
	}
	numTxs, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetLatestTxs, args[0])
	}

	latestTxs, err := m.obscuroScanAPI.GetLatestTransactions(numTxs)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetLatestTxs, err)
	}

	*result.(*[]gethcommon.Hash) = latestTxs
	return nil
}

func (m *MockServer) getBatchForTx(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatchForTx, len(args))
	}
	txHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatchForTx, args[0])
	}

	batch, err := m.obscuroScanAPI.GetBatchForTx(txHash)
	if err != nil {
		return fmt.Errorf("`%s` call failed. Cause: %w", rpc.GetBatchForTx, err)
	}

	*result.(**common.ExtBatch) = batch
	return nil
}

func (m *MockServer) getBatch(result interface{}, args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg to %s, got %d", rpc.GetBatch, len(args))
	}
	batchHash, ok := args[0].(gethcommon.Hash)
	if !ok {
		return fmt.Errorf("first arg to %s is of type %T, expected type int", rpc.GetBatch, args[0])
	}

	batch, err := m.obscuroScanAPI.GetBatch(batchHash)
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
