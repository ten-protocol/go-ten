package rpcclientlib

const (
	RPCGetID                   = "obscuro_getID"
	RPCGetCurrentBlockHead     = "obscuro_getCurrentBlockHead"
	RPCGetBlockHeaderByHash    = "obscuro_getBlockHeaderByHash"
	RPCGetCurrentRollupHead    = "obscuro_getCurrentRollupHead"
	RPCGetRollupHeader         = "obscuro_getRollupHeader"
	RPCGetRollupHeaderByNumber = "obscuro_getRollupHeaderByNumber"
	RPCGetRollup               = "obscuro_getRollup"
	RPCNonce                   = "obscuro_nonce"
	RPCAddViewingKey           = "obscuro_addViewingKey"
	RPCGetRollupForTx          = "obscuro_getRollupForTx"
	RPCGetLatestTxs            = "obscuro_getLatestTransactions"
	RPCGetTotalTxs             = "obscuro_getTotalTransactions"
	RPCAttestation             = "obscuro_attestation"
	RPCStopHost                = "obscuro_stopHost"
	RPCCall                    = "eth_call"
	RPCChainID                 = "eth_chainId"
	RPCGetBalance              = "eth_getBalance"
	RPCGetTransactionByHash    = "eth_getTransactionByHash"
	RPCGetTxReceipt            = "eth_getTransactionReceipt"
	RPCSendRawTransaction      = "eth_sendRawTransaction"
)

// Client is used by client applications to interact with the Obscuro node
type Client interface {
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}
