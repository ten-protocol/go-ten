package enclave

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/components"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/responses"

	"github.com/ten-protocol/go-ten/go/common"
	_ "github.com/ten-protocol/go-ten/go/common/tracers/native" // make sure the tracers are loaded
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	gethlog "github.com/ethereum/go-ethereum/log"
)

var _noHeadBatch = big.NewInt(0)

type enclaveInitService struct {
	config              *enclaveconfig.EnclaveConfig
	storage             storage.Storage
	l1BlockProcessor    components.L1BlockProcessor
	logger              gethlog.Logger
	sharedSecretService *crypto.SharedSecretService
	enclaveKeyService   *crypto.EnclaveAttestedKeyService // the enclave's private key (used to identify the enclave and sign messages)
	attestationProvider components.AttestationProvider    // interface for producing attestation reports and verifying them
	daEncryptionService *crypto.DAEncryptionService
	rpcKeyService       *crypto.RPCKeyService
}

func NewEnclaveInitAPI(config *enclaveconfig.EnclaveConfig, storage storage.Storage, logger gethlog.Logger, l1BlockProcessor components.L1BlockProcessor, enclaveKeyService *crypto.EnclaveAttestedKeyService, attestationProvider components.AttestationProvider, sharedSecretService *crypto.SharedSecretService, daEncryptionService *crypto.DAEncryptionService, rpcKeyService *crypto.RPCKeyService) common.EnclaveInit {
	return &enclaveInitService{
		config:              config,
		storage:             storage,
		l1BlockProcessor:    l1BlockProcessor,
		logger:              logger,
		enclaveKeyService:   enclaveKeyService,
		attestationProvider: attestationProvider,
		sharedSecretService: sharedSecretService,
		daEncryptionService: daEncryptionService,
		rpcKeyService:       rpcKeyService,
	}
}

func (e *enclaveInitService) Attestation(ctx context.Context) (*common.AttestationReport, common.SystemError) {
	if e.enclaveKeyService.PublicKey() == nil {
		return nil, responses.ToInternalError(fmt.Errorf("public key not initialized, we can't produce the attestation report"))
	}
	report, err := e.attestationProvider.CreateAttestationReport(ctx, e.config.HostAddress)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not produce remote report. Cause %w", err))
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
// it returns it encrypted with the enclave key
func (e *enclaveInitService) GenerateSecret(ctx context.Context) (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	e.sharedSecretService.GenerateSharedSecret()
	secret := e.sharedSecretService.Secret()
	if secret == nil {
		return nil, responses.ToInternalError(errors.New("failed to generate secret"))
	}
	err := e.storage.StoreSecret(ctx, *secret)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}

	err = e.notifyCryptoServices(*secret)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}

	encSec, err := e.enclaveKeyService.Encrypt(secret[:])
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed to encrypt secret. Cause: %w", err))
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a shared secret received from another enclave
func (e *enclaveInitService) InitEnclave(ctx context.Context, s common.EncryptedSharedEnclaveSecret) common.SystemError {
	secret, err := e.enclaveKeyService.Decrypt(s)
	if err != nil {
		return responses.ToInternalError(err)
	}
	err = e.storage.StoreSecret(ctx, crypto.SharedEnclaveSecret(secret))
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	var fixedSizeSecret crypto.SharedEnclaveSecret
	copy(fixedSizeSecret[:], secret)

	// notify the encryption services that depend on the shared secret
	err = e.notifyCryptoServices(fixedSizeSecret)
	if err != nil {
		return responses.ToInternalError(err)
	}
	return nil
}

func (e *enclaveInitService) notifyCryptoServices(sharedSecret crypto.SharedEnclaveSecret) error {
	e.sharedSecretService.SetSharedSecret(&sharedSecret)
	err := e.rpcKeyService.Initialise()
	if err != nil {
		return err
	}
	return e.daEncryptionService.Initialise()
}

func (e *enclaveInitService) EnclaveID(context.Context) (common.EnclaveID, common.SystemError) {
	return e.enclaveKeyService.EnclaveID(), nil
}

func (e *enclaveInitService) RPCEncryptionKey(ctx context.Context) ([]byte, common.SystemError) {
	k, err := e.rpcKeyService.PublicKey()
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return k, nil
}
