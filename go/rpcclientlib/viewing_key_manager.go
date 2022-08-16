package rpcclientlib

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// ViewingKeyManager interface allows us to provide different versions of the viewing key client. The initial motivation for this was:
// - we want the go-obscuro AuthClient to always be bound to a single account for code that it easy to reason about
// - we want the WalletExtension to
type ViewingKeyManager interface {
	// IsReady returns true if ViewingKeyManager has been initialised with a viewing key
	IsReady() bool
	// DecryptBytes expects response bytes and it will attempt to decrypt them using available Viewing Key(s)
	// todo: this method should probably take the request as well, so multi-acc implementations can use it for context
	DecryptBytes(respBytes []byte) ([]byte, error)
}

// SingleAccountVKManager is bound to one account and VK for its lifetime, this is the default way to use ViewingKeyClient
type SingleAccountVKManager struct {
	Address    common.Address
	ViewingKey *ecies.PrivateKey
}

func (s SingleAccountVKManager) IsReady() bool {
	// ViewingKey was signed and set at initialization time
	return true
}

func (s SingleAccountVKManager) DecryptBytes(respBytes []byte) ([]byte, error) {
	return s.ViewingKey.Decrypt(respBytes, nil, nil)
}
