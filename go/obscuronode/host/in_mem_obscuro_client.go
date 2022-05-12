package host

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

// An in-memory implementation of `clientserver.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	nodeID     common.Address
	obscuroAPI ObscuroAPI
}

func NewInMemObscuroClient(nodeID int64, p2p P2P) obscuroclient.Client {
	return &inMemObscuroClient{
		obscuroAPI: *NewObscuroAPI(p2p),
		nodeID:     common.BigToAddress(big.NewInt(nodeID)),
	}
}

func (c *inMemObscuroClient) ID() common.Address {
	return c.nodeID
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(_ interface{}, method string, args ...interface{}) error {
	if method == obscuroclient.RPCSendTransactionEncrypted {
		// TODO - Extract this checking logic as the set of RPC operations grows.
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCSendTransactionEncrypted, len(args))
		}
		tx, ok := args[0].(nodecommon.EncryptedTx)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type EncryptedTx", obscuroclient.RPCSendTransactionEncrypted)
		}

		c.obscuroAPI.SendTransactionEncrypted(tx)
	}
	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
