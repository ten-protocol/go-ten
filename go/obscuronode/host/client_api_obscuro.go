package host

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host *Node
}

func NewObscuroAPI(host *Node) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) {
	api.host.P2p.BroadcastTx(encryptedTx)
}

// GetCurrentBlockHeadHeight returns the height of the current head block.
func (api *ObscuroAPI) GetCurrentBlockHeadHeight() int64 {
	return api.host.nodeDB.GetCurrentBlockHead().Number.Int64()
}

// GetCurrentRollupHead returns the current head rollup's header.
func (api *ObscuroAPI) GetCurrentRollupHead() *nodecommon.Header {
	return api.host.nodeDB.GetCurrentRollupHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash common.Hash) *nodecommon.Header {
	return api.host.nodeDB.GetRollupHeader(hash)
}

// GetTransaction returns the transaction with the given hash.
func (api *ObscuroAPI) GetTransaction(hash common.Hash) *nodecommon.L2Tx {
	return api.host.EnclaveClient.GetTransaction(hash)
}

// Balance returns the balance of the wallet with the given address.
func (api *ObscuroAPI) Balance(address common.Address) uint64 {
	return api.host.EnclaveClient.Balance(address)
}

// Nonce returns the nonce of the wallet with the given address.
func (api *ObscuroAPI) Nonce(address common.Address) uint64 {
	return api.host.EnclaveClient.Nonce(address)
}

// StopHost gracefully stops the host.
func (api *ObscuroAPI) StopHost() {
	api.host.Stop()
}
