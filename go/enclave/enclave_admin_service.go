package enclave

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/integration/ethereummock"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	gethcommon "github.com/ethereum/go-ethereum/common"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/compression"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/common/profiler"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/events"
	"github.com/ten-protocol/go-ten/go/enclave/nodetype"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/responses"
)

type enclaveAdminService struct {
	config                 *enclaveconfig.EnclaveConfig
	mainMutex              sync.Mutex   // locks the admin operations
	dataInMutex            sync.RWMutex // controls access to data ingestion
	logger                 gethlog.Logger
	l1BlockProcessor       components.L1BlockProcessor
	validatorService       nodetype.Validator
	sequencerService       nodetype.ActiveSequencer
	service                nodetype.NodeType
	sharedSecretProcessor  *components.SharedSecretProcessor
	rollupConsumer         components.RollupConsumer
	registry               components.BatchRegistry
	daEncryptionService    *crypto.DAEncryptionService
	dataCompressionService compression.DataCompressionService
	storage                storage.Storage
	gethEncodingService    gethencoding.EncodingService
	stopControl            *stopcontrol.StopControl
	profiler               *profiler.Profiler
	subscriptionManager    *events.SubscriptionManager
	enclaveKeyService      *crypto.EnclaveAttestedKeyService
	mempool                *components.TxPool
	sharedSecretService    *crypto.SharedSecretService
	activeSequencer        bool
}

func NewEnclaveAdminAPI(config *enclaveconfig.EnclaveConfig, storage storage.Storage, logger gethlog.Logger, blockProcessor components.L1BlockProcessor, registry components.BatchRegistry, batchExecutor components.BatchExecutor, gethEncodingService gethencoding.EncodingService, stopControl *stopcontrol.StopControl, subscriptionManager *events.SubscriptionManager, enclaveKeyService *crypto.EnclaveAttestedKeyService, mempool *components.TxPool, chainConfig *params.ChainConfig, mgmtContractLib mgmtcontractlib.MgmtContractLib, attestationProvider components.AttestationProvider, sharedSecretService *crypto.SharedSecretService, daEncryptionService *crypto.DAEncryptionService) common.EnclaveAdmin {
	var prof *profiler.Profiler
	// don't run a profiler on an attested enclave
	if !config.WillAttest && config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultEnclavePort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler", log.ErrKey, err)
		}
	}
	sharedSecretProcessor := components.NewSharedSecretProcessor(mgmtContractLib, attestationProvider, enclaveKeyService.EnclaveID(), storage, sharedSecretService, logger)
	sigVerifier, err := getSignatureValidator(config.UseInMemoryDB, storage, logger)
	if err != nil {
		logger.Crit("Could not initialise the signature validator", log.ErrKey, err)
	}
	dataCompressionService := compression.NewBrotliDataCompressionService(int64(config.DecompressionLimit))

	rollupCompression := components.NewRollupCompression(registry, batchExecutor, daEncryptionService, dataCompressionService, storage, gethEncodingService, chainConfig, config, logger)
	rollupProducer := components.NewRollupProducer(enclaveKeyService.EnclaveID(), storage, registry, logger)
	rollupConsumer := components.NewRollupConsumer(mgmtContractLib, registry, rollupCompression, storage, logger, sigVerifier)

	seqSettings := nodetype.SequencerSettings{
		MaxBatchSize:      config.MaxBatchSize,
		MaxRollupSize:     config.MaxRollupSize,
		GasPaymentAddress: config.GasPaymentAddress,
		BatchGasLimit:     config.GasBatchExecutionLimit,
		BaseFee:           config.BaseFee,
	}

	sequencerService := nodetype.NewSequencer(blockProcessor, batchExecutor, registry, rollupProducer, rollupCompression, gethEncodingService, logger, chainConfig, enclaveKeyService, mempool, storage, dataCompressionService, seqSettings)
	validatorService := nodetype.NewValidator(blockProcessor, batchExecutor, registry, chainConfig, storage, sigVerifier, mempool, logger)

	eas := &enclaveAdminService{
		config:                 config,
		mainMutex:              sync.Mutex{},
		dataInMutex:            sync.RWMutex{},
		logger:                 logger,
		l1BlockProcessor:       blockProcessor,
		service:                validatorService,
		sequencerService:       sequencerService,
		validatorService:       validatorService,
		sharedSecretProcessor:  sharedSecretProcessor,
		rollupConsumer:         rollupConsumer,
		registry:               registry,
		daEncryptionService:    daEncryptionService,
		dataCompressionService: dataCompressionService,
		storage:                storage,
		gethEncodingService:    gethEncodingService,
		stopControl:            stopControl,
		profiler:               prof,
		subscriptionManager:    subscriptionManager,
		enclaveKeyService:      enclaveKeyService,
		mempool:                mempool,
		sharedSecretService:    sharedSecretService,
	}

	// if the current enclave was already marked as an active/backup sequencer, it needs to set the right mempool mode
	if eas.isBackupSequencer(context.Background()) || eas.isActiveSequencer(context.Background()) {
		mempool.SetValidateMode(false)
	}

	return eas
}

// addSequencer is used internally to add a sequencer enclaveID to the pool of attested enclaves.
// If it is the current enclave it will change the behaviour of this enclave to be a backup sequencer (ready to become active).
func (e *enclaveAdminService) addSequencer(id common.EnclaveID, _ types.Receipt) common.SystemError {
	e.logger.Info("Storing new sequencer enclaveID", log.EnclaveIDKey, id)
	err := e.storage.StoreNodeType(context.Background(), id, common.Sequencer)
	if err != nil {
		return responses.ToInternalError(err)
	}

	if e.enclaveKeyService.EnclaveID() == id {
		e.mempool.SetValidateMode(false)
	}

	return nil
}

func (e *enclaveAdminService) MakeActive() common.SystemError {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	if !e.isBackupSequencer(context.Background()) {
		// host may see this if it tries to promote its enclave before its ID has been added to the permission pool
		return fmt.Errorf("only backup sequencer can become active")
	}

	e.activeSequencer = true
	e.service = e.sequencerService
	e.logger.Info("Enclave is now active sequencer.")

	return nil
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
// Security considerations:
// - In phase 1 we assume the Sequencer node is benevolent.
// - Each batch (created and signed by the sequencer) refers to an L1 Block, which must be available for the batch to be processed
// - If the (validator) node operator is malicious and feeds the enclave with invalid L1 blocks, then batch processing will stop
func (e *enclaveAdminService) SubmitL1Block(ctx context.Context, blockData *common.ProcessedL1Data) (*common.BlockSubmissionResponse, common.SystemError) {
	e.dataInMutex.Lock()
	defer e.dataInMutex.Unlock()
	blockHeader := blockData.BlockHeader

	e.logger.Info("SubmitL1Block", log.BlockHeightKey, blockHeader.Number, log.BlockHashKey, blockHeader.Hash())

	// TODO verify proof provided with block blockData.Proof

	result, rollupMetadata, err := e.ingestL1Block(ctx, blockData)
	if err != nil {
		// only critical errors ie duplicate block or signed rollup error are returned so we can continue processing if non-critical
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	if result.IsFork() {
		e.logger.Info(fmt.Sprintf("Detected fork at block %s with height %d", blockHeader.Hash(), blockHeader.Number))
	}

	err = e.service.OnL1Block(ctx, blockHeader, result)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	err = e.storage.UpdateProcessed(ctx, blockHeader.Hash())
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	// in phase 1, only if the enclave is a sequencer, it can respond to shared secret requests
	canShareSecret := e.isBackupSequencer(ctx) || e.isActiveSequencer(ctx) || e.sharedSecretService.IsGenesis()

	bsr := &common.BlockSubmissionResponse{
		RollupMetadata:          rollupMetadata,
		ProducedSecretResponses: e.sharedSecretProcessor.ProcessNetworkSecretMsgs(ctx, blockData, canShareSecret),
	}

	// doing this after the network secret msgs to make sure we have stored the attestation before promotion.
	e.processSequencerPromotions(blockData)

	return bsr, nil
}

func (e *enclaveAdminService) processSequencerPromotions(blockData *common.ProcessedL1Data) {
	// todo handle sequencer revoked - could move all of this into a separate processor
	sequencerAddedTxs := blockData.GetEvents(common.SequencerAddedTx)
	for _, tx := range sequencerAddedTxs {
		if tx.HasSequencerEnclaveID() {
			err := e.addSequencer(tx.SequencerEnclaveID, *tx.Receipt)
			if err != nil {
				e.logger.Crit("Encountered error while adding sequencer enclaveID", log.ErrKey, err)
			}
		}
	}
}

func (e *enclaveAdminService) SubmitBatch(ctx context.Context, extBatch *common.ExtBatch) common.SystemError {
	if e.isActiveSequencer(ctx) {
		e.logger.Crit("Can't submit a batch to the active sequencer")
	}

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "SubmitBatch call completed.", &core.RelaxedThresholds, log.BatchHashKey, extBatch.Hash())

	e.logger.Info("Received new p2p batch", log.BatchHeightKey, extBatch.Header.Number, log.BatchHashKey, extBatch.Hash(), "l1", extBatch.Header.L1Proof)
	seqNo := extBatch.Header.SequencerOrderNo.Uint64()
	if seqNo > common.L2GenesisSeqNo+1 {
		_, err := e.storage.FetchBatchHeaderBySeqNo(ctx, seqNo-1)
		if err != nil {
			return responses.ToInternalError(fmt.Errorf("could not find previous batch with seq: %d", seqNo-1))
		}
	}

	batch, err := core.ToBatch(extBatch, e.daEncryptionService, e.dataCompressionService)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not convert batch. Cause: %w", err))
	}

	err = e.validator().VerifySequencerSignature(batch)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("invalid batch received. Could not verify signature. Cause: %w", err))
	}

	// calculate the converted hash, and store it in the db for chaining of the converted chain
	convertedHeader, err := e.gethEncodingService.CreateEthHeaderForBatch(ctx, extBatch.Header)
	if err != nil {
		return err
	}

	// todo - review whether we need to lock here.
	e.dataInMutex.Lock()
	// if the signature is valid, then store the batch together with the converted hash
	err = e.storage.StoreBatch(ctx, batch, convertedHeader.Hash())
	if err != nil {
		e.dataInMutex.Unlock()
		return responses.ToInternalError(fmt.Errorf("could not store batch. Cause: %w", err))
	}
	e.dataInMutex.Unlock()

	err = e.validator().ExecuteStoredBatches(ctx)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not execute batches. Cause: %w", err))
	}

	return nil
}

func (e *enclaveAdminService) CreateBatch(ctx context.Context, skipBatchIfEmpty bool) common.SystemError {
	if !e.isActiveSequencer(ctx) {
		e.logger.Crit("Only the active sequencer can create batches")
	}

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateBatch call ended", &core.RelaxedThresholds)

	e.dataInMutex.RLock()
	defer e.dataInMutex.RUnlock()

	err := e.sequencer().CreateBatch(ctx, skipBatchIfEmpty)
	if err != nil {
		return responses.ToInternalError(err)
	}

	return nil
}

func (e *enclaveAdminService) CreateRollup(ctx context.Context, fromSeqNo uint64) (*common.CreateRollupResult, common.SystemError) {
	if !e.isActiveSequencer(ctx) {
		e.logger.Crit("Only the active sequencer can create rollups")
	}
	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateRollup call ended", &core.RelaxedThresholds)

	// allow the simultaneous production of rollups and batches
	e.dataInMutex.RLock()
	defer e.dataInMutex.RUnlock()

	if e.registry.HeadBatchSeq() == nil {
		return nil, responses.ToInternalError(fmt.Errorf("not initialised yet"))
	}

	result, err := e.sequencer().CreateRollup(ctx, fromSeqNo)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return result, nil
}

func (e *enclaveAdminService) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, common.SystemError) {
	bundle, err := exportCrossChainData(ctx, e.storage, fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	bytes, err := bundle.HashPacked()
	if err != nil {
		return nil, err
	}

	sig, err := e.enclaveKeyService.Sign(bytes)
	if err != nil {
		return nil, err
	}

	bundle.Signature = sig
	return bundle, nil
}

func (e *enclaveAdminService) GetBatch(ctx context.Context, hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
	batch, err := e.storage.FetchBatch(ctx, hash)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed getting batch. Cause: %w", err))
	}

	b, err := batch.ToExtBatch(e.daEncryptionService, e.dataCompressionService)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return b, nil
}

func (e *enclaveAdminService) GetBatchBySeqNo(ctx context.Context, seqNo uint64) (*common.ExtBatch, common.SystemError) {
	batch, err := e.storage.FetchBatchBySeqNo(ctx, seqNo)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed getting batch. Cause: %w", err))
	}

	b, err := batch.ToExtBatch(e.daEncryptionService, e.dataCompressionService)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return b, nil
}

func (e *enclaveAdminService) GetRollupData(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, common.SystemError) {
	rollupMetadata, err := e.storage.FetchRollupMetadata(ctx, hash)
	if err != nil {
		return nil, err
	}
	metadata := &common.PublicRollupMetadata{
		FirstBatchSequence: rollupMetadata.FirstBatchSequence,
		StartTime:          rollupMetadata.StartTime,
	}
	return metadata, nil
}

func (e *enclaveAdminService) StreamL2Updates() (chan common.StreamL2UpdatesResponse, func()) {
	l2UpdatesChannel := make(chan common.StreamL2UpdatesResponse, 100)

	if e.stopControl.IsStopping() {
		close(l2UpdatesChannel)
		return l2UpdatesChannel, func() {}
	}

	e.registry.SubscribeForExecutedBatches(func(batch *core.Batch, receipts types.Receipts) {
		e.sendBatch(batch, l2UpdatesChannel)
		if receipts != nil {
			e.streamEventsForNewHeadBatch(context.Background(), batch, receipts, l2UpdatesChannel)
		}
	})

	return l2UpdatesChannel, func() {
		e.registry.UnsubscribeFromBatches()
	}
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveAdminService) HealthCheck(ctx context.Context) (bool, common.SystemError) {
	if e.stopControl.IsStopping() {
		return false, responses.ToInternalError(fmt.Errorf("requested HealthCheck with the enclave stopping"))
	}

	// if we have seen no sequencer permissioned on the L1 yet then we are in an unusual bootstrapping network state
	// and can return healthy
	enclaveIDs, err := e.storage.GetSequencerEnclaveIDs(ctx)
	if err != nil {
		return false, fmt.Errorf("could not get sequencer enclaveIDs. Cause: %w", err)
	}
	if len(enclaveIDs) == 0 {
		e.logger.Debug("No sequencer enclaveIDs found permissioned from L1, network is bootstrapping")
		return true, nil
	}

	// check the storage health
	storageHealthy, err := e.storage.HealthCheck(ctx)
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Info("HealthCheck failed for the enclave storage", log.ErrKey, err)
		return false, nil
	}

	// todo (#1148) - enclave healthcheck operations
	l1blockHealthy, err := e.l1BlockProcessor.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Info("HealthCheck failed for the l1 block processor", log.ErrKey, err)
		return false, nil
	}

	l2batchHealthy, err := e.registry.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Info("HealthCheck failed for the l2 batch registry", log.ErrKey, err)
		return false, nil
	}

	return storageHealthy && l1blockHealthy && l2batchHealthy, nil
}

func (e *enclaveAdminService) Status(ctx context.Context) (common.Status, common.SystemError) {
	initialised := e.sharedSecretService.IsInitialised()
	if !initialised {
		return common.Status{StatusCode: common.AwaitingSecret, L2Head: _noHeadBatch}, nil
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
	return common.Status{StatusCode: common.Running, L1Head: l1HeadHash, L2Head: l2HeadSeqNo, EnclaveID: enclaveID, IsActiveSequencer: e.activeSequencer}, nil
}

func (e *enclaveAdminService) Stop() common.SystemError {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	// block all requests
	e.stopControl.Stop()

	if e.profiler != nil {
		if err := e.profiler.Stop(); err != nil {
			e.logger.Error("Could not stop profiler", log.ErrKey, err)
			return err
		}
	}

	if e.registry != nil {
		e.registry.UnsubscribeFromBatches()
	}

	err := e.service.Close()
	if err != nil {
		e.logger.Error("Could not stop node service", log.ErrKey, err)
	}

	time.Sleep(time.Second)
	err = e.storage.Close()
	if err != nil {
		e.logger.Error("Could not stop db", log.ErrKey, err)
		return err
	}

	return nil
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveAdminService) StopClient() common.SystemError {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveAdminService) sendBatch(batch *core.Batch, outChannel chan common.StreamL2UpdatesResponse) {
	if batch.SeqNo().Uint64()%10 == 0 {
		e.logger.Info("Streaming batch to host", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SeqNo())
	} else {
		e.logger.Debug("Streaming batch to host", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SeqNo())
	}
	extBatch, err := batch.ToExtBatch(e.daEncryptionService, e.dataCompressionService)
	if err != nil {
		// this error is unrecoverable
		e.logger.Crit("failed to convert batch", log.ErrKey, err)
	}
	resp := common.StreamL2UpdatesResponse{
		Batch: extBatch,
	}
	outChannel <- resp
}

// this function is only called when the executed batch is the new head
func (e *enclaveAdminService) streamEventsForNewHeadBatch(ctx context.Context, batch *core.Batch, receipts types.Receipts, outChannel chan common.StreamL2UpdatesResponse) {
	logs, err := e.subscriptionManager.GetSubscribedLogsForBatch(ctx, batch, receipts)
	e.logger.Debug("Stream Events for", log.BatchHashKey, batch.Hash(), "nr_events", len(logs))
	if err != nil {
		e.logger.Error("Error while getting subscription logs", log.ErrKey, err)
		return
	}
	if logs != nil {
		outChannel <- common.StreamL2UpdatesResponse{
			Logs: logs,
		}
	}
}

func (e *enclaveAdminService) ingestL1Block(ctx context.Context, processed *common.ProcessedL1Data) (*components.BlockIngestionType, []common.ExtRollupMetadata, error) {
	e.logger.Info("Start ingesting block", log.BlockHashKey, processed.BlockHeader.Hash())
	ingestion, err := e.l1BlockProcessor.Process(ctx, processed)
	if err != nil {
		if e.isCriticalError(err) {
			e.logger.Debug("Did not ingest block", log.ErrKey, err, log.BlockHashKey, processed.BlockHeader.Hash())
		} else {
			e.logger.Warn("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, processed.BlockHeader.Hash())
		}
		return nil, nil, err
	}

	var rollupMetadataList []common.ExtRollupMetadata
	if processed.HasEvents(common.RollupTx) {
		rollupMetadataList, err = e.processRollups(ctx, processed, rollupMetadataList)
		if err != nil && e.isCriticalError(err) {
			// only propagate the error if we encounter a critical error on a sequencer signed rollup
			return nil, nil, err
		}
	}

	// Handle any L1 fork events
	if ingestion.IsFork() {
		e.registry.OnL1Reorg(ingestion)
		if err := e.service.OnL1Fork(ctx, ingestion.ChainFork); err != nil {
			return nil, nil, err
		}
	}
	return ingestion, rollupMetadataList, nil
}

func (e *enclaveAdminService) processRollups(ctx context.Context, processed *common.ProcessedL1Data, rollupMetadataList []common.ExtRollupMetadata) ([]common.ExtRollupMetadata, error) {
	rollupTxs := processed.GetEvents(common.RollupTx)
	txsSeen := make(map[gethcommon.Hash]bool)

	// verify and process each rollup one by one
	for _, rollupTx := range rollupTxs {
		extRollup, err := e.rollupConsumer.ExtractAndVerifyRollupData(rollupTx)
		if err != nil {
			// this will only be returned once we've verified the sequencer signature
			if errors.Is(err, errutil.ErrCriticalRollupProcessing) {
				return nil, err
			}
			// anything non-critical we skip the processing
			e.logger.Error("Error processing rollups from L1 data", log.ErrKey, err, log.BlockHashKey, processed.BlockHeader.Hash())
			continue
		}
		rHash := rollupTx.Transaction.Hash()
		// prevent the case where someone pushes a blob to the same slot. multiple rollups can be found in a block,
		// but they must come from unique transactions
		if txsSeen[rHash] {
			return nil, fmt.Errorf("multiple rollups from same transaction: %s. Err: %w", rHash, errutil.ErrCriticalRollupProcessing)
		}

		rollupMetadata, err := e.rollupConsumer.ProcessRollup(ctx, extRollup)
		if err != nil {
			// critical error as the sequencer has signed this rollup
			return nil, fmt.Errorf("failed to process rollup: %w", errutil.ErrCriticalRollupProcessing)
		}

		if rollupMetadata != nil {
			rollupMetadataList = append(rollupMetadataList, *rollupMetadata)
		}
		txsSeen[rHash] = true
	}
	if len(rollupMetadataList) == 0 {
		e.logger.Warn("No rollups found in block when rollupTxs present", log.BlockHashKey, processed.BlockHeader.Hash())
		return nil, nil
	}
	if len(rollupMetadataList) > 1 {
		// this is allowed as long as they come from unique transactions
		e.logger.Trace(fmt.Sprintf("Multiple rollups %d in block %s", len(rollupMetadataList), processed.BlockHeader.Hash()))
	}
	return rollupMetadataList, nil
}

func (e *enclaveAdminService) rejectBlockErr(ctx context.Context, cause error) *errutil.BlockRejectError {
	var hash common.L1BlockHash
	l1Head, err := e.l1BlockProcessor.GetHead(ctx)
	// todo - handle error
	if err == nil {
		hash = l1Head.Hash()
	}
	return &errutil.BlockRejectError{
		L1Head:  hash,
		Wrapped: cause,
	}
}

func (e *enclaveAdminService) validator() nodetype.Validator {
	validator, ok := e.service.(nodetype.Validator)
	if !ok {
		panic("enclave service is not a validator but validator was requested!")
	}
	return validator
}

func (e *enclaveAdminService) sequencer() nodetype.ActiveSequencer {
	sequencer, ok := e.service.(nodetype.ActiveSequencer)
	if !ok {
		panic("enclave service is not a sequencer but sequencer was requested!")
	}
	return sequencer
}

func (e *enclaveAdminService) isActiveSequencer(ctx context.Context) bool {
	return e.activeSequencer && e.getNodeType(ctx) == common.Sequencer
}

func (e *enclaveAdminService) isBackupSequencer(ctx context.Context) bool {
	return e.getNodeType(ctx) == common.Sequencer && !e.activeSequencer
}

func (e *enclaveAdminService) isValidator(ctx context.Context) bool { //nolint:unused
	return e.getNodeType(ctx) == common.Validator
}

func (e *enclaveAdminService) getNodeType(ctx context.Context) common.NodeType {
	id := e.enclaveKeyService.EnclaveID()
	attestedEnclave, err := e.storage.GetEnclavePubKey(ctx, id)
	if err != nil {
		// this log message doesn't need to be info level, we can assume enclave is in this state until we see a msg like
		// "Store attestation. Owner: <EnclaveID>" in the logs (i.e. after initial L1 catchup)
		e.logger.Trace("could not read enclave pub key. Defaulting to validator type", log.ErrKey, err)
		return common.Validator
	}
	return attestedEnclave.Type
}

func exportCrossChainData(ctx context.Context, storage storage.Storage, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, error) {
	canonicalBatches, err := storage.FetchCanonicalBatchesBetween((ctx), fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	if len(canonicalBatches) == 0 {
		return nil, errutil.ErrCrossChainBundleNoBatches
	}

	// todo - siliev - all those fetches need to be atomic
	header, err := storage.FetchHeadBatchHeader(ctx)
	if err != nil {
		return nil, err
	}

	blockHash := header.L1Proof
	batchHash := canonicalBatches[len(canonicalBatches)-1].Hash()

	block, err := storage.FetchBlock(ctx, blockHash)
	if err != nil {
		return nil, err
	}

	crossChainHashes := make([][]byte, 0)
	for _, batch := range canonicalBatches {
		if batch.CrossChainRoot != gethcommon.BigToHash(gethcommon.Big0) {
			crossChainHashes = append(crossChainHashes, batch.CrossChainRoot.Bytes())
		}
	}

	bundle := &common.ExtCrossChainBundle{
		LastBatchHash:        batchHash, // unused for now.
		L1BlockHash:          block.Hash(),
		L1BlockNum:           big.NewInt(0).Set(block.Number),
		CrossChainRootHashes: crossChainHashes,
	} // todo: check fromSeqNo
	return bundle, nil
}

func getSignatureValidator(useInMemDB bool, storage storage.Storage, logger gethlog.Logger) (components.SequencerSignatureVerifier, error) {
	if useInMemDB {
		return ethereummock.NewMockSignatureValidator(), nil
	}
	return components.NewSignatureValidator(storage, logger)
}

// isCriticalError returns true if the error should cause block processing to stop
func (e *enclaveAdminService) isCriticalError(err error) bool {
	return errors.Is(err, errutil.ErrCriticalRollupProcessing) ||
		errors.Is(err, errutil.ErrBlockAncestorNotFound) ||
		errors.Is(err, errutil.ErrBlockAlreadyProcessed)
}
