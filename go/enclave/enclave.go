package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common/compression"

	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/system"

	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/responses"

	"github.com/ten-protocol/go-ten/go/enclave/genesis"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/events"

	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	_ "github.com/ten-protocol/go-ten/go/common/tracers/native" // make sure the tracers are loaded

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type enclaveImpl struct {
	initAPI  common.EnclaveInit
	adminAPI common.EnclaveAdmin
	rpcAPI   common.EnclaveClientRPC

	stopControl *stopcontrol.StopControl
}

// NewEnclave creates and initializes all the services of the enclave.
//
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(config *enclaveconfig.EnclaveConfig, genesis *genesis.Genesis, mgmtContractLib mgmtcontractlib.MgmtContractLib, logger gethlog.Logger) common.Enclave {
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Creating enclave service with following config", log.CfgKey, string(jsonConfig))

	chainConfig := ethchainadapter.ChainParams(big.NewInt(config.TenChainID))

	// Initialise the database
	cachingService := storage.NewCacheService(logger, config.UseInMemoryDB)
	storage := storage.NewStorageFromConfig(config, cachingService, chainConfig, logger)

	// attempt to fetch the enclave key from the database
	// the enclave key is part of the attestation and identifies the current enclave
	// if this is the first time the enclave starts, it has to generate a new key
	enclaveKeyService := crypto.NewEnclaveAttestedKeyService(logger)
	err := loadOrCreateEnclaveKey(storage, enclaveKeyService, logger)
	if err != nil {
		logger.Crit("Failed to load or create enclave key", log.ErrKey, err)
	}

	sharedSecretService := crypto.NewSharedSecretService(logger)
	err = loadSharedSecret(storage, sharedSecretService, logger)
	if err != nil {
		logger.Crit("Failed to load shared secret", log.ErrKey, err)
	}

	daEncryptionService := crypto.NewDAEncryptionService(sharedSecretService, logger)
	rpcKeyService := crypto.NewRPCKeyService(sharedSecretService, logger)

	crossChainProcessors := crosschain.New(&config.MessageBusAddress, storage, logger)

	// initialise system contracts
	scb := system.NewSystemContractCallbacks(storage, &config.SystemContractOwner, logger)
	err = scb.Load(crossChainProcessors.Local)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		logger.Crit("failed to load system contracts", log.ErrKey, err)
	}

	// start the mempool in validate only. Based on the config, it might become sequencer
	evmEntropyService := crypto.NewEvmEntropyService(sharedSecretService, logger)
	gethEncodingService := gethencoding.NewGethEncodingService(storage, cachingService, evmEntropyService, logger)
	batchRegistry := components.NewBatchRegistry(storage, config, gethEncodingService, logger)
	mempool, err := components.NewTxPool(batchRegistry.EthChain(), config.MinGasPrice, true, logger)
	if err != nil {
		logger.Crit("unable to init eth tx pool", log.ErrKey, err)
	}

	gasOracle := gas.NewGasOracle()
	blockProcessor := components.NewBlockProcessor(storage, crossChainProcessors, gasOracle, logger)
	dataCompressionService := compression.NewBrotliDataCompressionService()
	batchExecutor := components.NewBatchExecutor(storage, batchRegistry, *config, gethEncodingService, crossChainProcessors, genesis, gasOracle, chainConfig, scb, evmEntropyService, mempool, dataCompressionService, logger)

	// ensure cached chain state data is up-to-date using the persisted batch data
	err = restoreStateDBCache(context.Background(), storage, batchRegistry, batchExecutor, genesis, logger)
	if err != nil {
		logger.Crit("failed to resync L2 chain state DB after restart", log.ErrKey, err)
	}

	subscriptionManager := events.NewSubscriptionManager(storage, batchRegistry, config.TenChainID, logger)

	// todo (#1474) - make sure the enclave cannot be started in production with WillAttest=false
	attestationProvider := components.NewAttestationProvider(enclaveKeyService, config.WillAttest, logger)

	// signal to stop the enclave
	stopControl := stopcontrol.New()

	// these services are directly exposed as the API of the Enclave
	initAPI := NewEnclaveInitAPI(config, storage, logger, blockProcessor, enclaveKeyService, attestationProvider, sharedSecretService, daEncryptionService, rpcKeyService)
	adminAPI := NewEnclaveAdminAPI(config, storage, logger, blockProcessor, batchRegistry, batchExecutor, gethEncodingService, stopControl, subscriptionManager, enclaveKeyService, mempool, chainConfig, mgmtContractLib, attestationProvider, sharedSecretService, daEncryptionService)
	rpcAPI := NewEnclaveRPCAPI(config, storage, logger, blockProcessor, batchRegistry, gethEncodingService, cachingService, mempool, chainConfig, crossChainProcessors, scb, subscriptionManager, genesis, gasOracle, sharedSecretService, rpcKeyService)

	logger.Info("Enclave service created successfully.", log.EnclaveIDKey, enclaveKeyService.EnclaveID())
	return &enclaveImpl{
		initAPI:     initAPI,
		adminAPI:    adminAPI,
		rpcAPI:      rpcAPI,
		stopControl: stopControl,
	}
}

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status(ctx context.Context) (common.Status, common.SystemError) {
	return e.initAPI.Status(ctx)
}

func (e *enclaveImpl) Attestation(ctx context.Context) (*common.AttestationReport, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.initAPI.Attestation(ctx)
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret(ctx context.Context) (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.initAPI.GenerateSecret(ctx)
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(ctx context.Context, s common.EncryptedSharedEnclaveSecret) common.SystemError {
	return e.initAPI.InitEnclave(ctx, s)
}

func (e *enclaveImpl) EnclaveID(ctx context.Context) (common.EnclaveID, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return common.EnclaveID{}, systemError
	}
	return e.initAPI.EnclaveID(ctx)
}

func (e *enclaveImpl) RPCEncryptionKey(ctx context.Context) ([]byte, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.initAPI.RPCEncryptionKey(ctx)
}

func (e *enclaveImpl) DebugTraceTransaction(ctx context.Context, txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	return e.rpcAPI.DebugTraceTransaction(ctx, txHash, config)
}

func (e *enclaveImpl) GetTotalContractCount(ctx context.Context) (*big.Int, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcAPI.GetTotalContractCount(ctx)
}

func (e *enclaveImpl) EnclavePublicConfig(ctx context.Context) (*common.EnclavePublicConfig, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcAPI.EnclavePublicConfig(ctx)
}

func (e *enclaveImpl) EncryptedRPC(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.EnclaveResponse, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcAPI.EncryptedRPC(ctx, encryptedParams)
}

func (e *enclaveImpl) GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash gethrpc.BlockNumberOrHash) ([]byte, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcAPI.GetCode(ctx, address, blockNrOrHash)
}

func (e *enclaveImpl) Subscribe(ctx context.Context, id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) common.SystemError {
	return e.rpcAPI.Subscribe(ctx, id, encryptedSubscription)
}

func (e *enclaveImpl) Unsubscribe(id gethrpc.ID) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.rpcAPI.Unsubscribe(id)
}

func (e *enclaveImpl) MakeActive() common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.adminAPI.MakeActive()
}

func (e *enclaveImpl) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.ExportCrossChainData(ctx, fromSeqNo, toSeqNo)
}

func (e *enclaveImpl) GetBatch(ctx context.Context, hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.GetBatch(ctx, hash)
}

func (e *enclaveImpl) GetBatchBySeqNo(ctx context.Context, seqNo uint64) (*common.ExtBatch, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.GetBatchBySeqNo(ctx, seqNo)
}

func (e *enclaveImpl) GetRollupData(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.GetRollupData(ctx, hash)
}

func (e *enclaveImpl) StreamL2Updates() (chan common.StreamL2UpdatesResponse, func()) {
	return e.adminAPI.StreamL2Updates()
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitL1Block(ctx context.Context, processed *common.ProcessedL1Data) (*common.BlockSubmissionResponse, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.SubmitL1Block(ctx, processed)
}

func (e *enclaveImpl) SubmitBatch(ctx context.Context, extBatch *common.ExtBatch) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.adminAPI.SubmitBatch(ctx, extBatch)
}

func (e *enclaveImpl) CreateBatch(ctx context.Context, skipBatchIfEmpty bool) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.adminAPI.CreateBatch(ctx, skipBatchIfEmpty)
}

func (e *enclaveImpl) CreateRollup(ctx context.Context, fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminAPI.CreateRollup(ctx, fromSeqNo)
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveImpl) HealthCheck(ctx context.Context) (bool, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return false, systemError
	}
	return e.adminAPI.HealthCheck(ctx)
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() common.SystemError {
	return e.adminAPI.StopClient()
}

func (e *enclaveImpl) Stop() common.SystemError {
	return e.adminAPI.Stop()
}

func checkStopping(s *stopcontrol.StopControl) common.SystemError {
	if s.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("enclave is stopping"))
	}
	return nil
}

func loadSharedSecret(storage storage.Storage, sharedSecretService *crypto.SharedSecretService, logger gethlog.Logger) error {
	sharedSecret, err := storage.FetchSecret(context.Background())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		logger.Crit("Failed to fetch secret", "err", err)
	}
	if sharedSecret != nil {
		sharedSecretService.SetSharedSecret(sharedSecret)
	}
	return nil
}

func loadOrCreateEnclaveKey(storage storage.Storage, enclaveKeyService *crypto.EnclaveAttestedKeyService, logger gethlog.Logger) error {
	enclaveKey, err := storage.GetEnclaveKey(context.Background())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("failed to load enclave key: %w", err)
	}
	if enclaveKey != nil {
		enclaveKeyService.SetEnclaveKey(enclaveKey)
		return nil
	}

	// enclave key not found - new key should be generated
	logger.Info("Generating new enclave key")
	enclaveKey, err = enclaveKeyService.GenerateEnclaveKey()
	if err != nil {
		return fmt.Errorf("failed to generate enclave key: %w", err)
	}
	err = storage.StoreEnclaveKey(context.Background(), enclaveKey)
	if err != nil {
		return fmt.Errorf("failed to store enclave key: %w", err)
	}

	enclaveKeyService.SetEnclaveKey(enclaveKey)
	return nil
}
