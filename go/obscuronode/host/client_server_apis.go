package host

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	p2p P2P
}

func NewObscuroAPI(p2p P2P) *ObscuroAPI {
	return &ObscuroAPI{p2p: p2p}
}

// SendTransactionEncrypted sends the encrypted Obscuro transaction to all peer Obscuro nodes.
func (api *ObscuroAPI) SendTransactionEncrypted(encryptedTx nodecommon.EncryptedTx) {
	api.p2p.BroadcastTx(encryptedTx)
}
