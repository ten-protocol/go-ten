package rpc

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

const (
	BatchNumber           = "eth_blockNumber"
	Call                  = "eth_call"
	ChainID               = "eth_chainId"
	GetBalance            = "eth_getBalance"
	GetBatchByHash        = "eth_getBlockByHash"
	GetBatchByNumber      = "eth_getBlockByNumber"
	GetCode               = "eth_getCode"
	GetTransactionByHash  = "eth_getTransactionByHash"
	GetTransactionCount   = "eth_getTransactionCount"
	GetTransactionReceipt = "eth_getTransactionReceipt"
	SendRawTransaction    = "eth_sendRawTransaction"
	EstimateGas           = "eth_estimateGas"
	GetLogs               = "eth_getLogs"
	GetStorageAt          = "eth_getStorageAt"
	GasPrice              = "eth_gasPrice"

	Health = "obscuro_health"
	Config = "obscuro_config"

	GetBlockHeaderByHash = "tenscan_getBlockHeaderByHash"
	GetBatch             = "tenscan_getBatch"
	GetBatchForTx        = "tenscan_getBatchForTx"
	GetLatestTxs         = "tenscan_getLatestTransactions"
	GetTotalTxs          = "tenscan_getTotalTransactions"
	Attestation          = "tenscan_attestation"
	StopHost             = "test_stopHost"
	Subscribe            = "eth_subscribe"
	Unsubscribe          = "eth_unsubscribe"
	SubscribeNamespace   = "eth"
	SubscriptionTypeLogs = "logs"

	// GetL1RollupHeaderByHash  = "scan_getL1RollupHeaderByHash"
	// GetActiveNodeCount       = "scan_getActiveNodeCount"

	GetLatestRollupHeader    = "scan_getLatestRollupHeader"
	GetTotalTransactionCount = "scan_getTotalTransactionCount"
	GetTotalContractCount    = "scan_getTotalContractCount"
	GetPublicTransactionData = "scan_getPublicTransactionData"
	GetBatchListing          = "scan_getBatchListing"
	GetBlockListing          = "scan_getBlockListing"
	GetFullBatchByHash       = "scan_getBatchByHash"
)

var ErrNilResponse = errors.New("nil response received from TEN node")

// Client is used by client applications to interact with the TEN node
type Client interface {
	// Call executes the named method via RPC. (Returns `ErrNilResponse` on nil response from Node, this is used as "not found" for some method calls)
	Call(result interface{}, method string, args ...interface{}) error
	// CallContext If the context is canceled before the call has successfully returned, CallContext returns immediately.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	// Subscribe creates a subscription to the TEN host.
	Subscribe(ctx context.Context, result interface{}, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
	// Stop closes the client.
	Stop()
}
