package rpcclientlib

import (
	"errors"
)

const (
	RPCGetID                   = "obscuro_getID"
	RPCGetCurrentBlockHead     = "obscuro_getCurrentBlockHead"
	RPCGetCurrentRollupHead    = "obscuro_getCurrentRollupHead"
	RPCGetRollupHeader         = "obscuro_getRollupHeader"
	RPCGetRollup               = "obscuro_getRollup"
	RPCAddViewingKey           = "obscuro_addViewingKey"
	RPCStopHost                = "obscuro_stopHost"
	RPCCall                    = "eth_call"
	RPCChainID                 = "eth_chainId"
	RPCGetBalance              = "eth_getBalance"
	RPCGetTransactionByHash    = "eth_getTransactionByHash"
	RPCNonce                   = "eth_getTransactionCount"
	RPCGetTxReceipt            = "eth_getTransactionReceipt"
	RPCSendRawTransaction      = "eth_sendRawTransaction"
	RPCGetBlockHeaderByHash    = "obscuroscan_getBlockHeaderByHash"
	RPCGetRollupHeaderByNumber = "obscuroscan_getRollupHeaderByNumber"
	RPCGetRollupForTx          = "obscuroscan_getRollupForTx"
	RPCGetLatestTxs            = "obscuroscan_getLatestTransactions"
	RPCGetTotalTxs             = "obscuroscan_getTotalTransactions"
	RPCAttestation             = "obscuroscan_attestation"
)

var ErrNilResponse = errors.New("nil response received from Obscuro node")

// Client is used by client applications to interact with the Obscuro node
type Client interface {
	// Call executes the named method via RPC. (Returns `ErrNilResponse` on nil response from Node, this is used as "not found" for some method calls)
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}
