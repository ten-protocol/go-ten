package rpc

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// EncryptionManager manages the decryption and encryption of sensitive RPC requests.
type EncryptionManager struct {
	enclavePrivateKeyECIES *ecies.PrivateKey
	// todo (#1445) - replace with persistent storage
	// todo - handle multiple viewing keys per address
	viewingKeys map[gethcommon.Address]*ecies.PublicKey // Maps account addresses to viewing public keys.
}

func NewEncryptionManager(enclavePrivateKeyECIES *ecies.PrivateKey) EncryptionManager {
	return EncryptionManager{
		enclavePrivateKeyECIES: enclavePrivateKeyECIES,
	}
}

// DecryptBytes decrypts the bytes with the enclave's private key.
func (rpc *EncryptionManager) DecryptBytes(encryptedBytes []byte) ([]byte, error) {
	bytes, err := rpc.enclavePrivateKeyECIES.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return bytes, nil
}
