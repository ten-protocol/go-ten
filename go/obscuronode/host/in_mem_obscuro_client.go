package host

import (
	"fmt"

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
		// TODO - Extract this checking logic as the set of RPC operations grows.
		if len(args) != 0 {
			return fmt.Errorf("expected 0 args to %s, got %d", obscuroclient.RPCSendTransactionEncrypted, len(args))
		}
		tx, ok := args[0].(nodecommon.EncryptedTx)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type EncryptedTx", obscuroclient.RPCSendTransactionEncrypted)
		}

		c.obscuroAPI.SendTransactionEncrypted(tx)
	}
	return nil
}

func (c inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
