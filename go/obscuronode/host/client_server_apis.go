package host

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	p2p     *P2P
	db      *DB
	enclave *nodecommon.Enclave
}

func NewObscuroAPI(p2p *P2P, db *DB, enclave *nodecommon.Enclave) *ObscuroAPI {
	return &ObscuroAPI{
		p2p:     p2p,
		db:      db,
		enclave: enclave,
	}
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) {
	(*api.p2p).BroadcastTx(encryptedTx)
}

// GetCurrentBlockHeadHeight returns the height of the current head block.
func (api *ObscuroAPI) GetCurrentBlockHeadHeight() int64 {
	return api.db.GetCurrentBlockHead().Number.Int64()
}

// GetCurrentRollupHead returns the current head rollup's header.
func (api *ObscuroAPI) GetCurrentRollupHead() *nodecommon.Header {
	return api.db.GetCurrentRollupHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash common.Hash) *nodecommon.Header {
	return api.db.GetRollupHeader(hash)
}

// GetTransaction returns the transaction with the given hash.
func (api *ObscuroAPI) GetTransaction(hash common.Hash) *nodecommon.L2Tx {
	return (*api.enclave).GetTransaction(hash)
}

// Balance returns the balance of the wallet with the given address.
func (api *ObscuroAPI) Balance(address common.Address) uint64 {
	return (*api.enclave).Balance(address)
}
