package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	// obscuroPrivateKeyHex is the private key used for sensitive communication with the enclave.
	// todo (#1053) - replace this fixed key with a key derived from the master seed.
	obscuroPrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"
	sharedSecretLen      = 32
)

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [sharedSecretLen]byte

func GetObscuroKey(logger gethlog.Logger) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(obscuroPrivateKeyHex)
	if err != nil {
		logger.Crit("failed to create enclave private key", log.ErrKey, err)
	}
	return key
}

func GenerateEntropy(logger gethlog.Logger) SharedEnclaveSecret {
	secret := make([]byte, sharedSecretLen)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		logger.Crit("could not generate secret", log.ErrKey, err)
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

func EncryptSecret(pubKeyEncoded []byte, secret SharedEnclaveSecret, logger gethlog.Logger) (common.EncryptedSharedEnclaveSecret, error) {
	logger.Info(fmt.Sprintf("Encrypting secret with public key %s", gethcommon.Bytes2Hex(pubKeyEncoded)))
	key, err := crypto.DecompressPubkey(pubKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := encryptWithPublicKey(secret[:], key)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to encrypt key, err: %s\nsecret: %v\npubkey: %v\nencKey:%v", err, secret, pubKeyEncoded, encKey))
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
		return nil, fmt.Errorf("could not decrypt with private key. %w", err)
	}
	return plaintext, nil
}

// CalculateRootBatchEntropy - calculates entropy per batch
// In Obscuro, we use a root entropy per batch, which is then used to calculate randomness exposed to individual transactions
// The RootBatchEntropy is calculated based on the shared secret and the batch height
// This ensures that sibling batches will naturally use the same root entropy so that transactions will have the same results
// Note that this formula is vulnerable to the unlikely event of a secret leak.
// todo (crypto) - find a way to hash in timestamp or something else then it would make it harder for attacker, such that sibling batches naturally have the same entropy.
func CalculateRootBatchEntropy(rootEntropy []byte, batchHeight *big.Int) gethcommon.Hash {
	return crypto.Keccak256Hash(rootEntropy, batchHeight.Bytes())
}

// CalculateTxRnd - calculates the randomness exposed to individual transactions
// In Obscuro, each tx must have its own randomness, independent from the others, because otherwise a malicious transaction
// could reveal information.
func CalculateTxRnd(rootBatchEntropy []byte, tCount int) gethcommon.Hash {
	return crypto.Keccak256Hash(rootBatchEntropy, intToBytes(tCount))
}

func intToBytes(val int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int64(val))
	if err != nil {
		panic(fmt.Sprintf("Could not convert int to bytes. Cause: %s", err))
	}
	return buf.Bytes()
}
