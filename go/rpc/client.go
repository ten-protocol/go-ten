package rpc

import (
	"context"
	"errors"
)

const (
	RPCAddViewingKey           = "obscuro_addViewingKey"
	RPCCall                    = "eth_call"
	RPCChainID                 = "eth_chainId"
	RPCGetBalance              = "eth_getBalance"
	RPCGetCode                 = "eth_getCode"
	RPCGetTransactionByHash    = "eth_getTransactionByHash"
	RPCNonce                   = "eth_getTransactionCount"
	RPCGetTxReceipt            = "eth_getTransactionReceipt"
	RPCSendRawTransaction      = "eth_sendRawTransaction"
	RPCGetBlockHeaderByHash    = "obscuroscan_getBlockHeaderByHash"
	RPCGetCurrentRollupHead    = "obscuroscan_getCurrentRollupHead"
	RPCGetRollup               = "obscuroscan_getRollup"
	RPCGetRollupHeaderByNumber = "obscuroscan_getRollupHeaderByNumber"
	RPCGetRollupForTx          = "obscuroscan_getRollupForTx"
	RPCGetLatestTxs            = "obscuroscan_getLatestTransactions"
	RPCGetTotalTxs             = "obscuroscan_getTotalTransactions"
	RPCAttestation             = "obscuroscan_attestation"
	RPCGetID                   = "test_getID"
	RPCGetCurrentBlockHead     = "test_getCurrentBlockHead"
	RPCGetRollupHeader         = "test_getRollupHeader"
	RPCStopHost                = "test_stopHost"
)

var ErrNilResponse = errors.New("nil response received from Obscuro node")

// Client is used by client applications to interact with the Obscuro node
type Client interface {
	// Call executes the named method via RPC. (Returns `ErrNilResponse` on nil response from Node, this is used as "not found" for some method calls)
	Call(result interface{}, method string, args ...interface{}) error
	// CallContext If the context is canceled before the call has successfully returned, CallContext returns immediately.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}
