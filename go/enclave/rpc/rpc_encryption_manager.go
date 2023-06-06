package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// EncryptionManager manages the decryption and encryption of sensitive RPC requests.
type EncryptionManager struct {
	enclavePrivateKeyECIES *ecies.PrivateKey
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
