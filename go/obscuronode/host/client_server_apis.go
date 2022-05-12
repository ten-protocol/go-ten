package host

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	p2p *P2P
	db  *DB
}

func NewObscuroAPI(p2p *P2P, db *DB) *ObscuroAPI {
	return &ObscuroAPI{
		p2p: p2p,
		db:  db,
	}
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) {
	(*api.p2p).BroadcastTx(encryptedTx)
}

// GetCurrentBlockHeadHeight returns the current block head from the node database.
func (api *ObscuroAPI) GetCurrentBlockHeadHeight() int64 {
	return api.db.GetCurrentBlockHead().Number.Int64()
}

// GetCurrentRollupHead returns the current head rollup's header.
func (api *ObscuroAPI) GetCurrentRollupHead() *nodecommon.Header {
	return api.db.GetCurrentRollupHead()
}
