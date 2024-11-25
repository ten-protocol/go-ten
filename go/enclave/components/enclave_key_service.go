package components

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/signature"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

type EnclaveKeyService struct {
	storage    storage.Storage
	logger     gethlog.Logger
	enclaveKey *crypto.EnclaveKey
}

func NewEnclaveKeyService(storage storage.Storage, logger gethlog.Logger) *EnclaveKeyService {
	return &EnclaveKeyService{storage: storage, logger: logger}
}

func (eks *EnclaveKeyService) LoadOrCreateEnclaveKey() error {
	enclaveKey, err := eks.storage.GetEnclaveKey(context.Background())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			eks.logger.Crit("Failed to fetch enclave key", log.ErrKey, err)
		}
		// enclave key not found - new key should be generated
		// todo (#1053) - revisit the crypto for this key generation/lifecycle before production
		eks.logger.Info("Generating new enclave key")
		enclaveKey, err = crypto.GenerateEnclaveKey()
		if err != nil {
			eks.logger.Crit("Failed to generate enclave key.", log.ErrKey, err)
		}
		err = eks.storage.StoreEnclaveKey(context.Background(), enclaveKey)
		if err != nil {
			eks.logger.Crit("Failed to store enclave key.", log.ErrKey, err)
		}
	}
	eks.logger.Info(fmt.Sprintf("Enclave key available. EnclaveID=%s, publicKey=%s", enclaveKey.EnclaveID(), gethcommon.Bytes2Hex(enclaveKey.PublicKeyBytes())))
	eks.enclaveKey = enclaveKey
	return err
}

func (eks *EnclaveKeyService) Sign(payload gethcommon.Hash) ([]byte, error) {
	return signature.Sign(payload.Bytes(), eks.enclaveKey.PrivateKey())
}

func (eks *EnclaveKeyService) EnclaveID() common.EnclaveID {
	return eks.enclaveKey.EnclaveID()
}

func (eks *EnclaveKeyService) PublicKey() *ecdsa.PublicKey {
	return eks.enclaveKey.PublicKey()
}

func (eks *EnclaveKeyService) PublicKeyBytes() []byte {
	return eks.enclaveKey.PublicKeyBytes()
}

func (eks *EnclaveKeyService) Decrypt(s common.EncryptedSharedEnclaveSecret) (*crypto.SharedEnclaveSecret, error) {
	return crypto.DecryptSecret(s, eks.enclaveKey.PrivateKey())
}
