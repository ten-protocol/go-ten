package rpc

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// ViewingKeySignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const ViewingKeySignedMsgPrefix = "vk"

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

// EncryptionManager manages the decryption and encryption of sensitive RPC requests.
type EncryptionManager struct {
	enclavePrivateKeyECIES *ecies.PrivateKey
	// TODO - Replace with persistent storage.
	// TODO - Handle multiple viewing keys per address.
	viewingKeys map[gethcommon.Address]*ecies.PublicKey // Maps account addresses to viewing public keys.
}

func NewEncryptionManager(enclavePrivateKeyECIES *ecies.PrivateKey) EncryptionManager {
	return EncryptionManager{
		enclavePrivateKeyECIES: enclavePrivateKeyECIES,
		viewingKeys:            make(map[gethcommon.Address]*ecies.PublicKey),
	}
}

// DecryptBytes decrypts the bytes with the enclave's private key if viewing keys are enabled.
func (rpc *EncryptionManager) DecryptBytes(encryptedBytes []byte) ([]byte, error) {
	bytes, err := rpc.enclavePrivateKeyECIES.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return bytes, nil
}

// AddViewingKey - see the description of Enclave.AddViewingKey.
func (rpc *EncryptionManager) AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error {
	// We decrypt the viewing key.
	viewingKeyBytes, err := rpc.enclavePrivateKeyECIES.Decrypt(encryptedViewingKeyBytes, nil, nil)
	if err != nil {
		return fmt.Errorf("could not decrypt viewing key when adding it to enclave. Cause: %w", err)
	}

	// We recalculate the message signed by MetaMask.
	msgToSign := ViewingKeySignedMsgPrefix + hex.EncodeToString(viewingKeyBytes)

	// We recover the key based on the signed message and the signature.
	recoveredAccountPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
	if err != nil {
		return fmt.Errorf("received viewing key but could not validate its signature. Cause: %w", err)
	}
	recoveredAccountAddress := crypto.PubkeyToAddress(*recoveredAccountPublicKey)

	// We decompress the viewing key and create the corresponding ECIES key.
	viewingKey, err := crypto.DecompressPubkey(viewingKeyBytes)
	if err != nil {
		return fmt.Errorf("received viewing key bytes but could not decompress them. Cause: %w", err)
	}
	viewingKeyECIES := ecies.ImportECDSAPublic(viewingKey)

	rpc.viewingKeys[recoveredAccountAddress] = viewingKeyECIES

	return nil
}

// EncryptWithViewingKey encrypts the bytes with a viewing key for the address.
func (rpc *EncryptionManager) EncryptWithViewingKey(address gethcommon.Address, bytes []byte) ([]byte, error) {
	viewingKey := rpc.viewingKeys[address]
	if viewingKey == nil {
		return nil, fmt.Errorf("could not encrypt bytes because it does not have a viewing key for account %s", address.String())
	}

	if len(bytes) == 0 {
		bytes = placeholderResult
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, viewingKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt bytes becauseit could not encrypt the response using a viewing key for account %s", address.String())
	}

	return encryptedBytes, nil
}
