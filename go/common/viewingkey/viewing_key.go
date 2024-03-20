package viewingkey

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave.
// It is the client-side perspective of the viewing key used for decrypting incoming traffic.
type ViewingKey struct {
	Account                 *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey              *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey               []byte              // ViewingKey public key in decrypt data from the enclave
	SignatureWithAccountKey []byte              // ViewingKey public key signed by the Accounts Private key - Allows to retrieve the Account address
	SignatureType           SignatureType       // Type of signature used to sign the public key
}

// RPCSignedViewingKey - used for transporting a minimalist viewing key via
// every RPC request to a sensitive method, including Log subscriptions.
// only the public key and the signature are required
// the account address is sent as well to aid validation
type RPCSignedViewingKey struct {
	PublicKey               []byte
	SignatureWithAccountKey []byte
	SignatureType           SignatureType
}
