package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// obscuroPrivateKeyHex is the private key used for sensitive communication with the enclave.
	// TODO - Replace this fixed key with a key derived from the master seed.
	obscuroPrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"
	sharedSecretLen      = 32
)

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [sharedSecretLen]byte

func GetObscuroKey() *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(obscuroPrivateKeyHex)
	if err != nil {
		log.Panic("failed to create enclave private key. Cause: %s", err)
	}
	return key
}

func GenerateEntropy() SharedEnclaveSecret {
	secret := make([]byte, sharedSecretLen)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		log.Panic("could not generate secret. Cause: %s", err)
	}
	var temp [sharedSecretLen]byte
	copy(temp[:], secret)
	return temp
}

func DecryptSecret(secret common.EncryptedSharedEnclaveSecret, privateKey *ecdsa.PrivateKey) (*SharedEnclaveSecret, error) {
	if privateKey == nil {
		return nil, errors.New("private key not found - shouldn't happen")
	}
	value, err := decryptWithPrivateKey(secret, privateKey)
	if err != nil {
		return nil, err
	}
	var temp SharedEnclaveSecret
	copy(temp[:], value)
	return &temp, nil
}

func EncryptSecret(pubKeyEncoded []byte, secret SharedEnclaveSecret, nodeShortID uint64) (common.EncryptedSharedEnclaveSecret, error) {
	common.LogWithID(nodeShortID, "Encrypting secret with public key %s", gethcommon.Bytes2Hex(pubKeyEncoded))
	key, err := crypto.DecompressPubkey(pubKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := encryptWithPublicKey(secret[:], key)
	if err != nil {
		common.LogWithID(nodeShortID, "Failed to encrypt key, err: %s\nsecret: %v\npubkey: %v\nencKey:%v", err, secret, pubKeyEncoded, encKey)
	}
	return encKey, err
}

// Encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *ecdsa.PublicKey) ([]byte, error) {
	ciphertext, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(pub), msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with public key. %w", err)
	}
	return ciphertext, nil
}

// Decrypts data with private key
func decryptWithPrivateKey(ciphertext []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	plaintext, err := ecies.ImportECDSA(priv).Decrypt(ciphertext, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with private key. %w", err)
	}
	return plaintext, nil
}

// GeneratePublicRandomness - generate 32 bytes of randomness, which will be exposed in the rollup header.
func GeneratePublicRandomness() []byte {
	return randomBytes(gethcommon.HashLength)
}

// PrivateRollupRnd - combine public randomness with private randomness in a way that protects the secret.
func PrivateRollupRnd(publicRnd []byte, secret []byte) []byte {
	return crypto.Keccak256Hash(publicRnd, secret).Bytes()
}

func randomBytes(length int) []byte {
	byteArr := make([]byte, length)
	if _, err := rand.Read(byteArr); err != nil {
		// todo - what should happen?
		panic(err)
	}
	return byteArr
}

// PerTransactionRnd - calculates a per tx random value
func PerTransactionRnd(privateRnd []byte, tCount int) []byte {
	return crypto.Keccak256Hash(privateRnd, big.NewInt(int64(tCount)).Bytes()).Bytes()
}
