package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

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
type DAEncryptionService struct {
	sharedSecretService *SharedSecretService
	cipher              *cipher.AEAD
	logger              gethlog.Logger
}

func NewDAEncryptionService(sharedSecretService *SharedSecretService, logger gethlog.Logger) *DAEncryptionService {
	da := &DAEncryptionService{
		sharedSecretService: sharedSecretService,
		logger:              logger,
	}

	_ = da.Initialise()

	return da
}

func (t *DAEncryptionService) Initialise() error {
	if !t.sharedSecretService.IsInitialised() {
		return errors.New("shared secret service is not initialised")
	}
	var err error
	t.cipher, err = createCypher(t.sharedSecretService)
	if err != nil {
		return fmt.Errorf("error creating cypher: %w", err)
	}
	return nil
}

func createCypher(sharedSecretService *SharedSecretService) (*cipher.AEAD, error) {
	key := sharedSecretService.ExtendEntropy([]byte{byte(daSuffix)})
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave DA key. cause %w", err)
	}
	cipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise GCM cipher for enclave DA key. cause %w", err)
	}
	return &cipher, nil
}

func (t *DAEncryptionService) Encrypt(blob []byte) ([]byte, error) {
	if t.cipher == nil {
		return nil, errors.New("not initialised")
	}

	nonce, err := generateSecureEntropy(GCMNonceLength)
	if err != nil {
		t.logger.Error("could not generate nonce to encrypt transactions.", log.ErrKey, err)
		return nil, err
	}

	ciphertext := (*t.cipher).Seal(nil, nonce, blob, nil)
	// We prepend the nonce to the ciphertext, so that it can be retrieved when decrypting.
	return append(nonce, ciphertext...), nil //nolint:makezero
}

func (t *DAEncryptionService) Decrypt(blob []byte) ([]byte, error) {
	if t.cipher == nil {
		return nil, errors.New("not initialised")
	}

	// The nonce is prepended to the ciphertext.
	nonce := blob[0:GCMNonceLength]
	ciphertext := blob[GCMNonceLength:]

	plaintext, err := (*t.cipher).Open(nil, nonce, ciphertext, nil)
	if err != nil {
		t.logger.Error("could not decrypt blob.", log.ErrKey, err)
		return nil, err
	}

	return plaintext, nil
}
