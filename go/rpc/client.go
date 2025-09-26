package rpc

import (
	"context"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// these are public RPC methods exposed by a TEN node
const (
	BatchNumber        = "ten_batchNumber"
	ChainID            = "ten_chainId"
	GetBatchByHash     = "ten_getBatchByHash"
	GetBatchByNumber   = "ten_getBatchByNumber"
	GetCode            = "ten_getCode"
	GasPrice           = "ten_gasPrice"
	GetCrossChainProof = "ten_getCrossChainProof"

	Health = "ten_health"
	Config = "ten_config"
	RPCKey = "ten_rpcKey"

	StopHost                 = "test_stopHost"
	SubscribeNamespace       = "ten"
	SubscriptionTypeLogs     = "logs"
	SubscriptionTypeNewHeads = "newHeads"

	GetBatchByTx               = "scan_getBatchByTx"
	GetLatestRollupHeader      = "scan_getLatestRollupHeader"
	GetTotalTxCount            = "scan_getTotalTransactionCount"
	GetTotalTxsQuery           = "scan_getTotalTransactionsQuery"
	GetHistoricalContractCount = "scan_getHistoricalContractCount"
	GetHistoricalTxCount       = "scan_getHistoricalTransactionCount"
	GetTotalContractCount      = "scan_getTotalContractCount"
	GetPublicTransactionData   = "scan_getPublicTransactionData"
	GetBatchListing            = "scan_getBatchListing"
	GetBlockListing            = "scan_getBlockListing"
	GetBatch                   = "scan_getBatch"
	GetLatestBatch             = "scan_getLatestBatch"
	GetBatchByHeight           = "scan_getBatchByHeight"
	GetBatchBySeqNo            = "scan_getBatchBySeq"
	GetTransaction             = "scan_getTransaction"

	GetRollupListing        = "scan_getRollupListing"
	GetRollupByHash         = "scan_getRollupByHash"
	GetRollupBatches        = "scan_getRollupBatches"
	GetRollupBySeqNo        = "scan_getRollupBySeqNo"
	GetBatchTransactions    = "scan_getBatchTransactions"
	GetPersonalTransactions = "scan_getPersonalTransactions"
	Search                  = "scan_search"
)

// Client is used by client applications to interact with the TEN node
type Client interface {
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// CallContext If the context is canceled before the call has successfully returned, CallContext returns immediately.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	// Subscribe creates a subscription to the Obscuro host.
	Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
	// Stop closes the client.
	Stop()
}
