package host

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

// An in-memory implementation of `clientserver.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	obscuroAPI ObscuroAPI
}

func NewInMemObscuroClient(p2p P2P) obscuroclient.Client {
	NewObscuroAPI(p2p)

	return inMemObscuroClient{obscuroAPI: *NewObscuroAPI(p2p)}
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c inMemObscuroClient) Call(_ interface{}, method string, args ...interface{}) error {
	if method == obscuroclient.RPCSendTransactionEncrypted {
		// todo - joel - be more resilient to errors here
		c.obscuroAPI.SendTransactionEncrypted(args[0].(nodecommon.EncryptedTx))
	}
	return nil
}

func (c inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
