package crypto

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
)

const (
	sharedSecretLenInBytes = 32
)

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [sharedSecretLenInBytes]byte

// SharedSecretService provides functionality to encapsulate, generate, extend, and encrypt the shared secret of the TEN network.
type SharedSecretService struct {
	secret    *SharedEnclaveSecret
	isGenesis bool
	logger    gethlog.Logger
}

func NewSharedSecretService(logger gethlog.Logger) *SharedSecretService {
	return &SharedSecretService{logger: logger}
}

// GenerateSharedSecret - called only by the genesis
func (sss *SharedSecretService) GenerateSharedSecret() {
	secret, err := generateSecureEntropy(sharedSecretLenInBytes)
	if err != nil {
		sss.logger.Crit("could not generate secret", log.ErrKey, err)
	}
	var tempSecret SharedEnclaveSecret
	copy(tempSecret[:], secret)
	sss.secret = &tempSecret
	sss.isGenesis = true
}

// Secret - should only be used before storing it
func (sss *SharedSecretService) Secret() *SharedEnclaveSecret {
	return sss.secret
}

func (sss *SharedSecretService) SetSharedSecret(ss *SharedEnclaveSecret) {
	sss.secret = ss
}

// ExtendEntropy derives more entropy from the shared secret
func (sss *SharedSecretService) ExtendEntropy(extra []byte) []byte {
	secretHash := crypto.Keccak256(sss.secret[:])
	return crypto.Keccak256(secretHash, extra)
}

func (sss *SharedSecretService) EncryptSecretWithKey(pubKey []byte) (common.EncryptedSharedEnclaveSecret, error) {
	sss.logger.Info(fmt.Sprintf("Encrypting secret with public key %s", gethcommon.Bytes2Hex(pubKey)))
	key, err := crypto.DecompressPubkey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := encryptWithPublicKey(sss.secret[:], key)
	if err != nil {
		sss.logger.Info("Failed to encrypt key", log.ErrKey, err)
	}
	return encKey, err
}

func (sss *SharedSecretService) IsInitialised() bool {
	return sss.secret != nil
}

func (sss *SharedSecretService) IsGenesis() bool {
	return sss.isGenesis
}
