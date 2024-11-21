package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common/compression"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/system"
	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/nodetype"

	"github.com/ten-protocol/go-ten/go/enclave/l2chain"
	"github.com/ten-protocol/go-ten/go/responses"

	"github.com/ten-protocol/go-ten/go/enclave/genesis"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/debugger"
	"github.com/ten-protocol/go-ten/go/enclave/events"

	"github.com/ten-protocol/go-ten/go/enclave/rpc"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	_ "github.com/ten-protocol/go-ten/go/common/tracers/native" // make sure the tracers are loaded

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type enclaveImpl struct {
	initService  common.EnclaveInit
	adminService common.EnclaveAdmin
	rpcService   common.EnclaveClientRPC
	stopControl  *stopcontrol.StopControl
}

// NewEnclave creates and initializes all the services of the enclave.
//
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(config *enclaveconfig.EnclaveConfig, genesis *genesis.Genesis, mgmtContractLib mgmtcontractlib.MgmtContractLib, logger gethlog.Logger) common.Enclave {
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Creating enclave service with following config", log.CfgKey, string(jsonConfig))

	// todo (#1053) - add the delay: N hashes

	// Initialise the database
	cachingService := storage.NewCacheService(logger, config.UseInMemoryDB)
	chainConfig := ethchainadapter.ChainParams(big.NewInt(config.ObscuroChainID))
	storage := storage.NewStorageFromConfig(config, cachingService, chainConfig, logger)

	// todo (#1474) - make sure the enclave cannot be started in production with WillAttest=false
	var attestationProvider components.AttestationProvider
	if config.WillAttest {
		attestationProvider = &components.EgoAttestationProvider{}
	} else {
		logger.Info("WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &components.DummyAttestationProvider{}
	}

	// attempt to fetch the enclave key from the database
	// the enclave key is part of the attestation and identifies the current enclave
	enclaveKey, err := loadOrCreateEnclaveKey(storage, logger)
	if err != nil {
		logger.Crit("Failed to load or create enclave key", "err", err)
	}

	gethEncodingService := gethencoding.NewGethEncodingService(storage, cachingService, logger)
	dataEncryptionService := crypto.NewDataEncryptionService(logger)
	dataCompressionService := compression.NewBrotliDataCompressionService()

	crossChainProcessors := crosschain.New(&config.MessageBusAddress, storage, big.NewInt(config.ObscuroChainID), logger)

	systemContractsWallet := system.GetPlaceholderWallet(chainConfig.ChainID, logger)
	scb := system.NewSystemContractCallbacks(systemContractsWallet, storage, logger)

	gasOracle := gas.NewGasOracle()
	blockProcessor := components.NewBlockProcessor(storage, crossChainProcessors, gasOracle, logger)
	registry := components.NewBatchRegistry(storage, logger)
	batchExecutor := components.NewBatchExecutor(storage, registry, *config, gethEncodingService, crossChainProcessors, genesis, gasOracle, chainConfig, config.GasBatchExecutionLimit, scb, logger)
	sigVerifier, err := components.NewSignatureValidator(storage)
	rProducer := components.NewRollupProducer(enclaveKey.EnclaveID(), storage, registry, logger)
	if err != nil {
		logger.Crit("Could not initialise the signature validator", log.ErrKey, err)
	}
	rollupCompression := components.NewRollupCompression(registry, batchExecutor, dataEncryptionService, dataCompressionService, storage, gethEncodingService, chainConfig, logger)
	rConsumer := components.NewRollupConsumer(mgmtContractLib, registry, rollupCompression, storage, logger, sigVerifier)
	sharedSecretProcessor := components.NewSharedSecretProcessor(mgmtContractLib, attestationProvider, enclaveKey.EnclaveID(), storage, logger)

	blockchain := ethchainadapter.NewEthChainAdapter(big.NewInt(config.ObscuroChainID), registry, storage, gethEncodingService, *config, logger)
	mempool, err := txpool.NewTxPool(blockchain, config.MinGasPrice, logger)
	if err != nil {
		logger.Crit("unable to init eth tx pool", log.ErrKey, err)
	}

	var service nodetype.NodeType
	if config.NodeType == common.ActiveSequencer {
		service = nodetype.NewSequencer(
			blockProcessor,
			batchExecutor,
			registry,
			rProducer,
			rollupCompression,
			gethEncodingService,
			logger,
			chainConfig,
			enclaveKey,
			mempool,
			storage,
			dataEncryptionService,
			dataCompressionService,
			nodetype.SequencerSettings{
				MaxBatchSize:      config.MaxBatchSize,
				MaxRollupSize:     config.MaxRollupSize,
				GasPaymentAddress: config.GasPaymentAddress,
				BatchGasLimit:     config.GasBatchExecutionLimit,
				BaseFee:           config.BaseFee,
			},
			blockchain,
		)
	} else {
		service = nodetype.NewValidator(
			blockProcessor,
			batchExecutor,
			registry,
			chainConfig,
			storage,
			sigVerifier,
			mempool,
			enclaveKey,
			logger,
		)
	}

	chain := l2chain.NewChain(
		storage,
		*config,
		gethEncodingService,
		chainConfig,
		genesis,
		logger,
		registry,
		config.GasLocalExecutionCapFlag,
	)
	obscuroKey := crypto.GetObscuroKey(logger)
	rpcEncryptionManager := rpc.NewEncryptionManager(ecies.ImportECDSA(obscuroKey), storage, cachingService, registry, crossChainProcessors, service, config, gasOracle, storage, blockProcessor, chain, logger)
	subscriptionManager := events.NewSubscriptionManager(storage, registry, config.ObscuroChainID, logger)

	// ensure cached chain state data is up-to-date using the persisted batch data
	err = restoreStateDBCache(context.Background(), storage, registry, batchExecutor, genesis, logger)
	if err != nil {
		logger.Crit("failed to resync L2 chain state DB after restart", log.ErrKey, err)
	}

	err = scb.Load()
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		logger.Crit("failed to load system contracts", log.ErrKey, err)
	}

	// TODO ensure debug is allowed/disallowed
	debug := debugger.New(chain, storage, chainConfig)
	stopControl := stopcontrol.New()
	initService := NewEnclaveInitService(config, storage, blockProcessor, logger, enclaveKey, attestationProvider)
	adminService := NewEnclaveAdminService(config, logger, blockProcessor, service, sharedSecretProcessor, rConsumer, registry, dataEncryptionService, dataCompressionService, storage, gethEncodingService, stopControl, subscriptionManager)
	rpcService := NewEnclaveRPCService(rpcEncryptionManager, registry, subscriptionManager, config, debug, storage, crossChainProcessors, scb)
	logger.Info("Enclave service created successfully.", log.EnclaveIDKey, enclaveKey.EnclaveID())
	return &enclaveImpl{
		initService:  initService,
		adminService: adminService,
		rpcService:   rpcService,
		stopControl:  stopControl,
	}
}

func loadOrCreateEnclaveKey(storage storage.Storage, logger gethlog.Logger) (*crypto.EnclaveKey, error) {
	enclaveKey, err := storage.GetEnclaveKey(context.Background())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			logger.Crit("Failed to fetch enclave key", log.ErrKey, err)
		}
		// enclave key not found - new key should be generated
		// todo (#1053) - revisit the crypto for this key generation/lifecycle before production
		logger.Info("Generating new enclave key")
		enclaveKey, err = crypto.GenerateEnclaveKey()
		if err != nil {
			logger.Crit("Failed to generate enclave key.", log.ErrKey, err)
		}
		err = storage.StoreEnclaveKey(context.Background(), enclaveKey)
		if err != nil {
			logger.Crit("Failed to store enclave key.", log.ErrKey, err)
		}
	}
	logger.Info(fmt.Sprintf("Enclave key available. EnclaveID=%s, publicKey=%s", enclaveKey.EnclaveID(), gethcommon.Bytes2Hex(enclaveKey.PublicKeyBytes())))
	return enclaveKey, err
}

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status(ctx context.Context) (common.Status, common.SystemError) {
	return e.initService.Status(ctx)
}

func (e *enclaveImpl) Attestation(ctx context.Context) (*common.AttestationReport, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.initService.Attestation(ctx)
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret(ctx context.Context) (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.initService.GenerateSecret(ctx)
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(ctx context.Context, s common.EncryptedSharedEnclaveSecret) common.SystemError {
	return e.initService.InitEnclave(ctx, s)
}

func (e *enclaveImpl) EnclaveID(ctx context.Context) (common.EnclaveID, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return common.EnclaveID{}, systemError
	}
	return e.initService.EnclaveID(ctx)
}

func (e *enclaveImpl) DebugTraceTransaction(ctx context.Context, txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	return e.rpcService.DebugTraceTransaction(ctx, txHash, config)
}

func (e *enclaveImpl) GetTotalContractCount(ctx context.Context) (*big.Int, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcService.GetTotalContractCount(ctx)
}

func (e *enclaveImpl) EnclavePublicConfig(ctx context.Context) (*common.EnclavePublicConfig, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcService.EnclavePublicConfig(ctx)
}

func (e *enclaveImpl) EncryptedRPC(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.EnclaveResponse, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcService.EncryptedRPC(ctx, encryptedParams)
}

func (e *enclaveImpl) GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash gethrpc.BlockNumberOrHash) ([]byte, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.rpcService.GetCode(ctx, address, blockNrOrHash)
}

func (e *enclaveImpl) Subscribe(ctx context.Context, id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) common.SystemError {
	return e.rpcService.Subscribe(ctx, id, encryptedSubscription)
}

func (e *enclaveImpl) Unsubscribe(id gethrpc.ID) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.rpcService.Unsubscribe(id)
}

func (e *enclaveImpl) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.ExportCrossChainData(ctx, fromSeqNo, toSeqNo)
}

func (e *enclaveImpl) GetBatch(ctx context.Context, hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.GetBatch(ctx, hash)
}

func (e *enclaveImpl) GetBatchBySeqNo(ctx context.Context, seqNo uint64) (*common.ExtBatch, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.GetBatchBySeqNo(ctx, seqNo)
}

func (e *enclaveImpl) GetRollupData(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.GetRollupData(ctx, hash)
}

func (e *enclaveImpl) StreamL2Updates() (chan common.StreamL2UpdatesResponse, func()) {
	return e.adminService.StreamL2Updates()
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitL1Block(ctx context.Context, blockHeader *types.Header, receipts []*common.TxAndReceiptAndBlobs) (*common.BlockSubmissionResponse, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.SubmitL1Block(ctx, blockHeader, receipts)
}

func (e *enclaveImpl) SubmitBatch(ctx context.Context, extBatch *common.ExtBatch) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.adminService.SubmitBatch(ctx, extBatch)
}

func (e *enclaveImpl) CreateBatch(ctx context.Context, skipBatchIfEmpty bool) common.SystemError {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return systemError
	}
	return e.adminService.CreateBatch(ctx, skipBatchIfEmpty)
}

func (e *enclaveImpl) CreateRollup(ctx context.Context, fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return nil, systemError
	}
	return e.adminService.CreateRollup(ctx, fromSeqNo)
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveImpl) HealthCheck(ctx context.Context) (bool, common.SystemError) {
	if systemError := checkStopping(e.stopControl); systemError != nil {
		return false, systemError
	}
	return e.adminService.HealthCheck(ctx)
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() common.SystemError {
	return e.adminService.StopClient()
}

func (e *enclaveImpl) Stop() common.SystemError {
	return e.adminService.Stop()
}

func checkStopping(s *stopcontrol.StopControl) common.SystemError {
	if s.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("enclave is stopping"))
	}
	return nil
}
