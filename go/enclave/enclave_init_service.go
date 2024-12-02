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

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	_ "github.com/ten-protocol/go-ten/go/common/tracers/native" // make sure the tracers are loaded
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

var _noHeadBatch = big.NewInt(0)

type enclaveInitService struct {
	config              *enclaveconfig.EnclaveConfig
	storage             storage.Storage
	l1BlockProcessor    components.L1BlockProcessor
	logger              gethlog.Logger
	enclaveKeyService   *components.EnclaveKeyService  // the enclave's private key (used to identify the enclave and sign messages)
	attestationProvider components.AttestationProvider // interface for producing attestation reports and verifying them
}

func NewEnclaveInitService(config *enclaveconfig.EnclaveConfig, storage storage.Storage, logger gethlog.Logger, l1BlockProcessor components.L1BlockProcessor, enclaveKeyService *components.EnclaveKeyService, attestationProvider components.AttestationProvider) common.EnclaveInit {
	return &enclaveInitService{
		config:              config,
		storage:             storage,
		l1BlockProcessor:    l1BlockProcessor,
		logger:              logger,
		enclaveKeyService:   enclaveKeyService,
		attestationProvider: attestationProvider,
	}
}

// Status is only implemented by the RPC wrapper
func (e *enclaveInitService) Status(ctx context.Context) (common.Status, common.SystemError) {
	_, err := e.storage.FetchSecret(ctx)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return common.Status{StatusCode: common.AwaitingSecret, L2Head: _noHeadBatch}, nil
		}
		return common.Status{StatusCode: common.Unavailable}, responses.ToInternalError(err)
	}
	var l1HeadHash gethcommon.Hash
	l1Head, err := e.l1BlockProcessor.GetHead(ctx)
	if err != nil {
		// this might be normal while enclave is starting up, just send empty hash
		e.logger.Debug("failed to fetch L1 head block for status response", log.ErrKey, err)
	} else {
		l1HeadHash = l1Head.Hash()
	}
	// we use zero when there's no head batch yet, the first seq number is 1
	l2HeadSeqNo := _noHeadBatch
	// this is the highest seq number that has been received and stored on the enclave (it may not have been executed)
	currSeqNo, err := e.storage.FetchCurrentSequencerNo(ctx)
	if err != nil {
		// this might be normal while enclave is starting up, just send empty hash
		e.logger.Debug("failed to fetch L2 head batch for status response", log.ErrKey, err)
	} else {
		l2HeadSeqNo = currSeqNo
	}
	enclaveID := e.enclaveKeyService.EnclaveID()
	return common.Status{StatusCode: common.Running, L1Head: l1HeadHash, L2Head: l2HeadSeqNo, EnclaveID: enclaveID}, nil
}

func (e *enclaveInitService) Attestation(ctx context.Context) (*common.AttestationReport, common.SystemError) {
	if e.enclaveKeyService.PublicKey() == nil {
		return nil, responses.ToInternalError(fmt.Errorf("public key not initialized, we can't produce the attestation report"))
	}
	report, err := e.attestationProvider.GetReport(ctx, e.enclaveKeyService.PublicKeyBytes(), e.enclaveKeyService.EnclaveID(), e.config.HostAddress)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not produce remote report. Cause %w", err))
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveInitService) GenerateSecret(ctx context.Context) (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	secret := crypto.GenerateEntropy(e.logger)
	err := e.storage.StoreSecret(ctx, secret)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	encSec, err := crypto.EncryptSecret(e.enclaveKeyService.PublicKeyBytes(), secret, e.logger)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed to encrypt secret. Cause: %w", err))
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveInitService) InitEnclave(ctx context.Context, s common.EncryptedSharedEnclaveSecret) common.SystemError {
	secret, err := e.enclaveKeyService.Decrypt(s)
	if err != nil {
		return responses.ToInternalError(err)
	}
	err = e.storage.StoreSecret(ctx, *secret)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	e.logger.Trace(fmt.Sprintf("Secret decrypted and stored. Secret: %v", secret))
	return nil
}

func (e *enclaveInitService) EnclaveID(context.Context) (common.EnclaveID, common.SystemError) {
	return e.enclaveKeyService.EnclaveID(), nil
}
