package rpcencryptionmanager

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ViewingKeySignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const ViewingKeySignedMsgPrefix = "vk"

// PlaceholderResult is used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil
// using ECIES throws an exception.
var PlaceholderResult = []byte("<nil result>")

// RPCEncryptionManager manages the decryption and encryption of sensitive RPC requests.
type RPCEncryptionManager struct {
	viewingKeysEnabled     bool
	enclavePrivateKeyECIES *ecies.PrivateKey
	// TODO - Replace with persistent storage.
	// TODO - Handle multiple viewing keys per address.
	viewingKeys map[common.Address]*ecies.PublicKey
}

func NewRPCEncryptionManager(viewingKeysEnabled bool, enclavePrivateKeyECIES *ecies.PrivateKey) RPCEncryptionManager {
	return RPCEncryptionManager{
		viewingKeysEnabled:     viewingKeysEnabled,
		enclavePrivateKeyECIES: enclavePrivateKeyECIES,
		viewingKeys:            make(map[common.Address]*ecies.PublicKey),
	}
}

// DecryptWithEnclaveKey decrypts the bytes with the enclave's public key.
func (e *RPCEncryptionManager) DecryptWithEnclaveKey(encryptedBytes []byte) ([]byte, error) {
	if !e.viewingKeysEnabled {
		return encryptedBytes, nil
	}

	paramBytes, err := e.enclavePrivateKeyECIES.Decrypt(encryptedBytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt bytes with enclave private key. Cause: %w", err)
	}

	return paramBytes, nil
}

// AddViewingKey - see the description of Enclave.AddViewingKey.
func (e *RPCEncryptionManager) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	// We recalculate the message signed by MetaMask.
	msgToSign := ViewingKeySignedMsgPrefix + hex.EncodeToString(viewingKeyBytes)

	// We recover the key based on the signed message and the signature.
	recoveredPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
	if err != nil {
		return fmt.Errorf("received viewing key but could not validate its signature. Cause: %w", err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPublicKey)

	// We decompress the viewing key and create the corresponding ECIES key.
	viewingKey, err := crypto.DecompressPubkey(viewingKeyBytes)
	if err != nil {
		return fmt.Errorf("received viewing key bytes but could not decompress them. Cause: %w", err)
	}
	eciesPublicKey := ecies.ImportECDSAPublic(viewingKey)

	e.viewingKeys[recoveredAddress] = eciesPublicKey

	return nil
}

// EncryptWithViewingKey encrypts the bytes with a viewing key for the address.
func (e *RPCEncryptionManager) EncryptWithViewingKey(address common.Address, bytes []byte) (nodecommon.EncryptedResponse, error) {
	if !e.viewingKeysEnabled {
		return bytes, nil
	}

	viewingKey := e.viewingKeys[address]
	if viewingKey == nil {
		return nil, fmt.Errorf("could not encrypt bytes because it does not have a viewing key for account %s", address.String())
	}

	if bytes == nil {
		bytes = PlaceholderResult
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, viewingKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt bytes becauseit could not encrypt the response using a viewing key for account %s", address.String())
	}

	return encryptedBytes, nil
}
