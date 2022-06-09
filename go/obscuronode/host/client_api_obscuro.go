package host

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// GetID returns the ID of the host.
func (api *ObscuroAPI) GetID() common.Address {
	return api.host.ID
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) {
	api.host.P2p.BroadcastTx(encryptedTx)
}

// GetCurrentBlockHead returns the current head block's header.
func (api *ObscuroAPI) GetCurrentBlockHead() *types.Header {
	return api.host.nodeDB.GetCurrentBlockHead()
}

// GetCurrentRollupHead returns the current head rollup's header.
func (api *ObscuroAPI) GetCurrentRollupHead() *nodecommon.Header {
	return api.host.nodeDB.GetCurrentRollupHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash common.Hash) *nodecommon.Header {
	return api.host.nodeDB.GetRollupHeader(hash)
}

// GetRollup returns the rollup with the given hash.
func (api *ObscuroAPI) GetRollup(hash common.Hash) *nodecommon.ExtRollup {
	return api.host.EnclaveClient.GetRollup(hash)
}

// GetTransaction returns the transaction with the given hash.
func (api *ObscuroAPI) GetTransaction(hash common.Hash) *nodecommon.L2Tx {
	return api.host.EnclaveClient.GetTransaction(hash)
}

// ExecContract returns the result of executing the smart contract as a user.
// `data` is generally generated from the ABI of a smart contract.
func (api *ObscuroAPI) ExecContract(from common.Address, contractAddress common.Address, data []byte) OffChainResponse {
	r, err := api.host.EnclaveClient.ExecuteOffChainTransaction(from, contractAddress, data)
	return OffChainResponse{
		Response: r,
		Error:    err,
	}
}

// Nonce returns the nonce of the wallet with the given address.
func (api *ObscuroAPI) Nonce(address common.Address) uint64 {
	return api.host.EnclaveClient.Nonce(address)
}

// AddViewingKey stores the viewing key on the enclave.
func (api *ObscuroAPI) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	return api.host.EnclaveClient.AddViewingKey(viewingKeyBytes, signature)
}

// StopHost gracefully stops the host.
func (api *ObscuroAPI) StopHost() {
	go api.host.Stop()
}

type OffChainResponse struct {
	Response nodecommon.EncryptedResponse
	Error    error
}
