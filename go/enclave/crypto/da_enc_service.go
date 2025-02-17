package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"sync"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/log"
)

const (
	// GCMNonceLength is the nonce's length in bytes for encrypting and decrypting transactions.
	GCMNonceLength = 12
	// daSuffix is used for generating the encryption key from the shared secret
	daSuffix = 0
)

// DAEncryptionService - handles encryption/decryption of the data stored in the DA layer
// using AES-GCM with a shared secret. It prepends the nonce to encrypted data.
//
// Thread-safe for concurrent usage.
type DAEncryptionService struct {
	sharedSecretService *SharedSecretService
	cipher              cipher.AEAD
	logger              gethlog.Logger
	mu                  sync.RWMutex
}

func NewDAEncryptionService(sharedSecretService *SharedSecretService, logger gethlog.Logger) *DAEncryptionService {
	da := &DAEncryptionService{
		sharedSecretService: sharedSecretService,
		logger:              logger,
	}

	// ignore the error because the service will be initialised later
	_ = da.Initialise()

	return da
}

func (t *DAEncryptionService) Initialise() error {
	if !t.sharedSecretService.IsInitialised() {
		return errors.New("shared secret service is not initialised")
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	var err error
	t.cipher, err = createCipher(t.sharedSecretService)
	if err != nil {
		return fmt.Errorf("error creating cypher: %w", err)
	}
	return nil
}

func createCipher(sharedSecretService *SharedSecretService) (cipher.AEAD, error) {
	key := sharedSecretService.ExtendEntropy([]byte{byte(daSuffix)})
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave DA key. cause %w", err)
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise GCM cipher for enclave DA key. cause %w", err)
	}
	return cipher, nil
}

func (t *DAEncryptionService) Encrypt(blob []byte) ([]byte, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.cipher == nil {
		return nil, errors.New("not initialised")
	}

	nonce, err := generateSecureEntropy(GCMNonceLength)
	if err != nil {
		t.logger.Error("could not generate nonce to encrypt transactions.", log.ErrKey, err)
		return nil, fmt.Errorf("nonce generation failed: %w", err)
	}

	result := make([]byte, GCMNonceLength+len(blob)+t.cipher.Overhead())
	copy(result[:GCMNonceLength], nonce)

	t.cipher.Seal(result[GCMNonceLength:GCMNonceLength], nonce, blob, nil)
	return result, nil
}

func (t *DAEncryptionService) Decrypt(blob []byte) ([]byte, error) {
	if len(blob) <= GCMNonceLength {
		return nil, errors.New("invalid encrypted blob size")
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.cipher == nil {
		return nil, errors.New("not initialised")
	}

	// The nonce is prepended to the ciphertext.
	nonce := blob[0:GCMNonceLength]
	ciphertext := blob[GCMNonceLength:]

	plaintext, err := t.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		t.logger.Error("could not decrypt blob.", log.ErrKey, err)
		return nil, err
	}

	return plaintext, nil
}
