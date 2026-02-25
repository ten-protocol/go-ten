package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// Encryptor provides AES-GCM encryption/decryption with the following characteristics:
//   - Uses AES-256-GCM (Galois/Counter Mode) with a 32-byte key
//   - Generates a random 12-byte nonce for each encryption operation using crypto/rand
//   - The nonce is prepended to the ciphertext output from Encrypt() and is generated
//     using crypto/rand.Reader for cryptographically secure random values
//
// Additionally provides HMAC-SHA256 hashing functionality:
// - Uses the same 32-byte key as the encryption operations
// - Generates a 32-byte (256-bit) message authentication code
// - Suitable for creating secure message digests and verifying data integrity
type Encryptor struct {
	gcm cipher.AEAD
	key []byte
}

func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes long")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Encryptor{gcm: gcm, key: key}, nil
}

func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return e.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < e.gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:e.gcm.NonceSize()], ciphertext[e.gcm.NonceSize():]
	return e.gcm.Open(nil, nonce, ciphertext, nil)
}

func (e *Encryptor) HashWithHMAC(data []byte) []byte {
	h := hmac.New(sha256.New, e.key)
	h.Write(data)
	return h.Sum(nil)
}

// GetKey returns the encryption key
func (e *Encryptor) GetKey() []byte {
	return e.key
}
