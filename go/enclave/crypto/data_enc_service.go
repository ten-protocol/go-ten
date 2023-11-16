package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	// RollupEncryptionKeyHex is the AES key used to encrypt and decrypt the transaction blob in rollups.
	// todo (#1053) - replace this fixed key with derived, rotating keys.
	RollupEncryptionKeyHex = "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c"
	// NonceLength is the nonce's length in bytes for encrypting and decrypting transactions.
	NonceLength = 12
)

// DataEncryptionService handles the encryption and decryption of the transaction blobs stored inside a rollup.
type DataEncryptionService interface {
	Encrypt(blob []byte) ([]byte, error)
	Decrypt(blob []byte) ([]byte, error)
}

type dataEncryptionServiceImpl struct {
	cipher cipher.AEAD
	logger gethlog.Logger
}

func NewDataEncryptionService(logger gethlog.Logger) DataEncryptionService {
	key := gethcommon.Hex2Bytes(RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Crit("could not initialise AES cipher for enclave rollup key.", log.ErrKey, err)
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		logger.Crit("could not initialise wrapper for AES cipher for enclave rollup key. ", log.ErrKey, err)
	}
	return dataEncryptionServiceImpl{
		cipher: cipher,
		logger: logger,
	}
}

// todo (#1053) - modify this logic so that transactions with different reveal periods are in different blobs, as per the whitepaper.
func (t dataEncryptionServiceImpl) Encrypt(blob []byte) ([]byte, error) {
	nonce := make([]byte, NonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		t.logger.Error("could not generate nonce to encrypt transactions.", log.ErrKey, err)
		return nil, err
	}

	// todo - ensure this nonce is not used too many times (2^32?) with the same key, to avoid risk of repeat.
	ciphertext := t.cipher.Seal(nil, nonce, blob, nil)
	// We prepend the nonce to the ciphertext, so that it can be retrieved when decrypting.
	return append(nonce, ciphertext...), nil //nolint:makezero
}

func (t dataEncryptionServiceImpl) Decrypt(blob []byte) ([]byte, error) {
	// The nonce is prepended to the ciphertext.
	nonce := blob[0:NonceLength]
	ciphertext := blob[NonceLength:]

	plaintext, err := t.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		t.logger.Error("could not decrypt blob.", log.ErrKey, err)
		return nil, err
	}

	return plaintext, nil
}
