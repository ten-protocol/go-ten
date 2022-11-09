package rpc

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

const (
	Call                    = "eth_call"
	ChainID                 = "eth_chainId"
	GetBalance              = "eth_getBalance"
	GetCode                 = "eth_getCode"
	GetTransactionByHash    = "eth_getTransactionByHash"
	GetTransactionCount     = "eth_getTransactionCount"
	GetTransactionReceipt   = "eth_getTransactionReceipt"
	SendRawTransaction      = "eth_sendRawTransaction"
	EstimateGas             = "eth_estimateGas"
	GetLogs                 = "eth_getLogs"
	AddViewingKey           = "obscuro_addViewingKey"
	GetBlockHeaderByHash    = "obscuroscan_getBlockHeaderByHash"
	GetCurrentRollupHead    = "obscuroscan_getCurrentRollupHead"
	GetRollup               = "obscuroscan_getRollup"
	GetRollupHeaderByNumber = "obscuroscan_getRollupHeaderByNumber"
	GetRollupForTx          = "obscuroscan_getRollupForTx"
	GetLatestTxs            = "obscuroscan_getLatestTransactions"
	GetTotalTxs             = "obscuroscan_getTotalTransactions"
	Attestation             = "obscuroscan_attestation"
	GetID                   = "test_getID"
	GetHeadBlockHeader      = "test_getHeadBlockHeader"
	GetRollupHeader         = "test_getRollupHeader"
	StopHost                = "test_stopHost"
	Subscribe               = "eth_subscribe"
	SubscribeNamespace      = "eth"
	SubscriptionTypeLogs    = "logs"
)

var ErrNilResponse = errors.New("nil response received from Obscuro node")

// Client is used by client applications to interact with the Obscuro node
type Client interface {
	// Call executes the named method via RPC. (Returns `ErrNilResponse` on nil response from Node, this is used as "not found" for some method calls)
	Call(result interface{}, method string, args ...interface{}) error
	// CallContext If the context is canceled before the call has successfully returned, CallContext returns immediately.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	// Subscribe creates a subscription to the Obscuro host.
	Subscribe(ctx context.Context, result interface{}, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
	// Stop closes the client.
	Stop()
}
