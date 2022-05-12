package clientserver

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	p2p host.P2P
}

func NewObscuroAPI(p2p host.P2P) *ObscuroAPI {
	return &ObscuroAPI{p2p: p2p}
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) { //nolint
	api.p2p.BroadcastTx(encryptedTx)
}
