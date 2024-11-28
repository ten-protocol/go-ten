package enclave

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

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
	mainMutex              sync.Mutex // serialises all data ingestion or creation to avoid weird races
	logger                 gethlog.Logger
	l1BlockProcessor       components.L1BlockProcessor
	service                nodetype.NodeType
	sharedSecretProcessor  *components.SharedSecretProcessor
	rollupConsumer         components.RollupConsumer
	registry               components.BatchRegistry
	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	storage                storage.Storage
	gethEncodingService    gethencoding.EncodingService
	stopControl            *stopcontrol.StopControl
	profiler               *profiler.Profiler
	subscriptionManager    *events.SubscriptionManager
	enclaveKeyService      *components.EnclaveKeyService
}

func NewEnclaveAdminService(config *enclaveconfig.EnclaveConfig, logger gethlog.Logger, l1BlockProcessor components.L1BlockProcessor, service nodetype.NodeType, sharedSecretProcessor *components.SharedSecretProcessor, rollupConsumer components.RollupConsumer, registry components.BatchRegistry, dataEncryptionService crypto.DataEncryptionService, dataCompressionService compression.DataCompressionService, storage storage.Storage, gethEncodingService gethencoding.EncodingService, stopControl *stopcontrol.StopControl, subscriptionManager *events.SubscriptionManager, enclaveKeyService *components.EnclaveKeyService) common.EnclaveAdmin {
	var prof *profiler.Profiler
	// don't run a profiler on an attested enclave
	if !config.WillAttest && config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultEnclavePort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler", log.ErrKey, err)
		}
	}

	return &enclaveAdminService{
		config:                 config,
		mainMutex:              sync.Mutex{},
		logger:                 logger,
		l1BlockProcessor:       l1BlockProcessor,
		service:                service,
		sharedSecretProcessor:  sharedSecretProcessor,
		rollupConsumer:         rollupConsumer,
		registry:               registry,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		storage:                storage,
		gethEncodingService:    gethEncodingService,
		stopControl:            stopControl,
		profiler:               prof,
		subscriptionManager:    subscriptionManager,
		enclaveKeyService:      enclaveKeyService,
	}
}

func (e *enclaveAdminService) AddSequencer(id common.EnclaveID, proof types.Receipt) common.SystemError {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	// by default all enclaves start their life as a validator

	// store in the database the enclave id
	err := e.storage.StoreNodeType(context.Background(), id, common.BackupSequencer)
	if err != nil {
		return responses.ToInternalError(err)
	}

	// compare the id with the current enclaveId and if they match - do something so that the current enclave behaves as a "backup sequencer"
	// the host will specifically mark the active enclave
	//currentEnclaveId, err := e.initService.EnclaveID(context.Background())
	//if err != nil {
	//	return err
	//}

	//if currentEnclaveId == id {
	//	todo
	//}

	// todo - use the proof
	return nil
}

func (e *enclaveAdminService) MakeActive() common.SystemError {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	if !e.isBackupSequencer(context.Background()) {
		return fmt.Errorf("only backup sequencer can become active")
	}
	// todo
	// change the node type service
	// do something with the mempool
	// make some other checks?
	// Once we've got the sequencer Enclave IDs permission list monitoring we should include that check here probably.
	// We could even make it so that all sequencer enclaves start as backup and it can't be activated until the permissioning is done?

	return nil
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveAdminService) SubmitL1Block(ctx context.Context, blockHeader *types.Header, processed *common.ProcessedL1Data) (*common.BlockSubmissionResponse, common.SystemError) {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	e.logger.Info("SubmitL1Block", log.BlockHeightKey, blockHeader.Number, log.BlockHashKey, blockHeader.Hash())

	// Verify the block header matches the one in processedData
	if blockHeader.Hash() != processed.BlockHeader.Hash() {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("block header mismatch"))
	}
	result, err := e.ingestL1Block(ctx, processed)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	if result.IsFork() {
		e.logger.Info(fmt.Sprintf("Detected fork at block %s with height %d", blockHeader.Hash(), blockHeader.Number))
	}

	err = e.service.OnL1Block(ctx, blockHeader, result)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	bsr := &common.BlockSubmissionResponse{ProducedSecretResponses: e.sharedSecretProcessor.ProcessNetworkSecretMsgs(ctx, processed)}
	return bsr, nil
}

func (e *enclaveAdminService) SubmitBatch(ctx context.Context, extBatch *common.ExtBatch) common.SystemError {
	if e.isActiveSequencer(ctx) {
		e.logger.Crit("Can't submit a batch to the active sequencer")
	}

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "SubmitBatch call completed.", log.BatchHashKey, extBatch.Hash())

	e.logger.Info("Received new p2p batch", log.BatchHeightKey, extBatch.Header.Number, log.BatchHashKey, extBatch.Hash(), "l1", extBatch.Header.L1Proof)
	seqNo := extBatch.Header.SequencerOrderNo.Uint64()
	if seqNo > common.L2GenesisSeqNo+1 {
		_, err := e.storage.FetchBatchHeaderBySeqNo(ctx, seqNo-1)
		if err != nil {
			return responses.ToInternalError(fmt.Errorf("could not find previous batch with seq: %d", seqNo-1))
		}
	}

	batch, err := core.ToBatch(extBatch, e.dataEncryptionService, e.dataCompressionService)
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

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	// if the signature is valid, then store the batch together with the converted hash
	err = e.storage.StoreBatch(ctx, batch, convertedHeader.Hash())
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not store batch. Cause: %w", err))
	}

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

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateBatch call ended")

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	err := e.sequencer().CreateBatch(ctx, skipBatchIfEmpty)
	if err != nil {
		return responses.ToInternalError(err)
	}

	return nil
}

func (e *enclaveAdminService) CreateRollup(ctx context.Context, fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	if !e.isActiveSequencer(ctx) {
		e.logger.Crit("Only the active sequencer can create rollups")
	}
	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateRollup call ended")

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	if e.registry.HeadBatchSeq() == nil {
		return nil, responses.ToInternalError(fmt.Errorf("not initialised yet"))
	}

	rollup, err := e.sequencer().CreateRollup(ctx, fromSeqNo)
	// TODO do we need to store the blob hashes here so we can check them against our records?
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return rollup, nil
}

func (e *enclaveAdminService) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, common.SystemError) {
	bundle, err := exportCrossChainData(ctx, e.storage, fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	sig, err := e.enclaveKeyService.Sign(bundle.HashPacked())
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

	b, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
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

	b, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
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
	extBatch, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
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

func (e *enclaveAdminService) ingestL1Block(ctx context.Context, processed *common.ProcessedL1Data) (*components.BlockIngestionType, error) {
	e.logger.Info("Start ingesting block", log.BlockHashKey, processed.BlockHeader.Hash())
	ingestion, err := e.l1BlockProcessor.Process(ctx, processed)
	if err != nil {
		// only warn for unexpected errors
		if errors.Is(err, errutil.ErrBlockAncestorNotFound) || errors.Is(err, errutil.ErrBlockAlreadyProcessed) {
			e.logger.Debug("Did not ingest block", log.ErrKey, err, log.BlockHashKey, processed.BlockHeader.Hash())
		} else {
			e.logger.Warn("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, processed.BlockHeader.Hash())
		}
		return nil, err
	}

	err = e.rollupConsumer.ProcessBlobsInBlock(ctx, processed)
	if err != nil && !errors.Is(err, components.ErrDuplicateRollup) {
		e.logger.Error("Encountered error while processing l1 block", log.ErrKey, err)
		// Unsure what to do here; block has been stored
	}

	if ingestion.IsFork() {
		e.registry.OnL1Reorg(ingestion)
		err := e.service.OnL1Fork(ctx, ingestion.ChainFork)
		if err != nil {
			return nil, err
		}
	}
	return ingestion, nil
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
	return e.getNodeType(ctx) == common.ActiveSequencer
}

func (e *enclaveAdminService) isBackupSequencer(ctx context.Context) bool {
	return e.getNodeType(ctx) == common.BackupSequencer
}

func (e *enclaveAdminService) isValidator(ctx context.Context) bool { //nolint:unused
	return e.getNodeType(ctx) == common.Validator
}

func (e *enclaveAdminService) getNodeType(ctx context.Context) common.NodeType {
	id := e.enclaveKeyService.EnclaveID()
	_, nodeType, err := e.storage.GetEnclavePubKey(ctx, id)
	if err != nil {
		e.logger.Crit("could not read enclave pub key", log.ErrKey, err)
		return 0
	}
	return nodeType
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
