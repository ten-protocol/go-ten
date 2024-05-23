package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/go/common/compression"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/nodetype"

	"github.com/ten-protocol/go-ten/go/enclave/l2chain"
	"github.com/ten-protocol/go-ten/go/responses"

	"github.com/ten-protocol/go-ten/go/enclave/genesis"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/profiler"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/debugger"
	"github.com/ten-protocol/go-ten/go/enclave/events"

	"github.com/ten-protocol/go-ten/go/enclave/rpc"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	_ "github.com/ten-protocol/go-ten/go/common/tracers/native" // make sure the tracers are loaded

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

var _noHeadBatch = big.NewInt(0)

type enclaveImpl struct {
	config                *config.EnclaveConfig
	storage               storage.Storage
	blockResolver         storage.BlockResolver
	l1BlockProcessor      components.L1BlockProcessor
	rollupConsumer        components.RollupConsumer
	l1Blockchain          *gethcore.BlockChain
	rpcEncryptionManager  *rpc.EncryptionManager
	subscriptionManager   *events.SubscriptionManager
	crossChainProcessors  *crosschain.Processors
	sharedSecretProcessor *components.SharedSecretProcessor

	chain     l2chain.ObscuroChain
	service   nodetype.NodeType
	registry  components.BatchRegistry
	gasOracle gas.Oracle

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	attestationProvider components.AttestationProvider // interface for producing attestation reports and verifying them

	enclaveKey *crypto.EnclaveKey // the enclave's private key (used to identify the enclave and sign messages)

	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	gethEncodingService    gethencoding.EncodingService
	profiler               *profiler.Profiler
	debugger               *debugger.Debugger
	logger                 gethlog.Logger

	stopControl *stopcontrol.StopControl
	mainMutex   sync.Mutex // serialises all data ingestion or creation to avoid weird races
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(
	config *config.EnclaveConfig,
	genesis *genesis.Genesis,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	logger gethlog.Logger,
) common.Enclave {
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Creating enclave service with following config", log.CfgKey, string(jsonConfig))

	// todo (#1053) - add the delay: N hashes

	var prof *profiler.Profiler
	// don't run a profiler on an attested enclave
	if !config.WillAttest && config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultEnclavePort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler", log.ErrKey, err)
		}
	}

	// Initialise the database
	chainConfig := ethchainadapter.ChainParams(big.NewInt(config.ObscuroChainID))
	storage := storage.NewStorageFromConfig(config, chainConfig, logger)

	// Initialise the Ethereum "Blockchain" structure that will allow us to validate incoming blocks
	// todo (#1056) - valid block
	var l1Blockchain *gethcore.BlockChain
	if config.ValidateL1Blocks {
		if config.GenesisJSON == nil {
			logger.Crit("enclave is configured to validate blocks, but genesis JSON is nil")
		}
		l1Blockchain = l2chain.NewL1Blockchain(config.GenesisJSON, logger)
	} else {
		logger.Info("validateBlocks is set to false. L1 blocks will not be validated.")
	}

	// todo (#1474) - make sure the enclave cannot be started in production with WillAttest=false
	var attestationProvider components.AttestationProvider
	if config.WillAttest {
		attestationProvider = &components.EgoAttestationProvider{}
	} else {
		logger.Info("WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &components.DummyAttestationProvider{}
	}

	// attempt to fetch the enclave key from the database
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

	obscuroKey := crypto.GetObscuroKey(logger)

	gethEncodingService := gethencoding.NewGethEncodingService(storage, logger)
	dataEncryptionService := crypto.NewDataEncryptionService(logger)
	dataCompressionService := compression.NewBrotliDataCompressionService()

	crossChainProcessors := crosschain.New(&config.MessageBusAddress, storage, big.NewInt(config.ObscuroChainID), logger)

	gasOracle := gas.NewGasOracle()
	blockProcessor := components.NewBlockProcessor(storage, crossChainProcessors, gasOracle, logger)
	registry := components.NewBatchRegistry(storage, logger)
	batchExecutor := components.NewBatchExecutor(storage, registry, *config, gethEncodingService, crossChainProcessors, genesis, gasOracle, chainConfig, config.GasBatchExecutionLimit, logger)
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
	if config.NodeType == common.Sequencer {
		service = nodetype.NewSequencer(
			blockProcessor,
			batchExecutor,
			registry,
			rProducer,
			rConsumer,
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
			rConsumer,
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
	rpcEncryptionManager := rpc.NewEncryptionManager(ecies.ImportECDSA(obscuroKey), storage, registry, crossChainProcessors, service, config, gasOracle, storage, blockProcessor, chain, logger)
	subscriptionManager := events.NewSubscriptionManager(storage, registry, config.ObscuroChainID, logger)

	// ensure cached chain state data is up-to-date using the persisted batch data
	err = restoreStateDBCache(context.Background(), storage, registry, batchExecutor, genesis, logger)
	if err != nil {
		logger.Crit("failed to resync L2 chain state DB after restart", log.ErrKey, err)
	}

	// TODO ensure debug is allowed/disallowed
	debug := debugger.New(chain, storage, chainConfig)

	logger.Info("Enclave service created successfully.", log.EnclaveIDKey, enclaveKey.EnclaveID())
	return &enclaveImpl{
		config:                 config,
		storage:                storage,
		blockResolver:          storage,
		l1BlockProcessor:       blockProcessor,
		rollupConsumer:         rConsumer,
		l1Blockchain:           l1Blockchain,
		rpcEncryptionManager:   rpcEncryptionManager,
		subscriptionManager:    subscriptionManager,
		crossChainProcessors:   crossChainProcessors,
		mgmtContractLib:        mgmtContractLib,
		attestationProvider:    attestationProvider,
		sharedSecretProcessor:  sharedSecretProcessor,
		enclaveKey:             enclaveKey,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		gethEncodingService:    gethEncodingService,
		profiler:               prof,
		logger:                 logger,
		debugger:               debug,
		stopControl:            stopcontrol.New(),

		chain:     chain,
		registry:  registry,
		service:   service,
		gasOracle: gasOracle,

		mainMutex: sync.Mutex{},
	}
}

func (e *enclaveImpl) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, common.SystemError) {
	return e.service.ExportCrossChainData(ctx, fromSeqNo, toSeqNo)
}

func (e *enclaveImpl) GetBatch(ctx context.Context, hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
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

func (e *enclaveImpl) GetBatchBySeqNo(ctx context.Context, seqNo uint64) (*common.ExtBatch, common.SystemError) {
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

func (e *enclaveImpl) GetRollupData(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, common.SystemError) {
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

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status(ctx context.Context) (common.Status, common.SystemError) {
	if e.stopControl.IsStopping() {
		return common.Status{StatusCode: common.Unavailable}, responses.ToInternalError(fmt.Errorf("requested Status with the enclave stopping"))
	}

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
	return common.Status{StatusCode: common.Running, L1Head: l1HeadHash, L2Head: l2HeadSeqNo}, nil
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() common.SystemError {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveImpl) sendBatch(batch *core.Batch, outChannel chan common.StreamL2UpdatesResponse) {
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
func (e *enclaveImpl) streamEventsForNewHeadBatch(ctx context.Context, batch *core.Batch, receipts types.Receipts, outChannel chan common.StreamL2UpdatesResponse) {
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

func (e *enclaveImpl) StreamL2Updates() (chan common.StreamL2UpdatesResponse, func()) {
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

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitL1Block(ctx context.Context, block *common.L1Block, receipts common.L1Receipts, isLatest bool) (*common.BlockSubmissionResponse, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitL1Block with the enclave stopping"))
	}

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	e.logger.Info("SubmitL1Block", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash())

	// If the block and receipts do not match, reject the block.
	br, err := common.ParseBlockAndReceipts(block, &receipts)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	result, err := e.ingestL1Block(ctx, br)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	if result.IsFork() {
		e.logger.Info(fmt.Sprintf("Detected fork at block %s with height %d", block.Hash(), block.Number()))
	}

	err = e.service.OnL1Block(ctx, block, result)
	if err != nil {
		return nil, e.rejectBlockErr(ctx, fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	bsr := &common.BlockSubmissionResponse{ProducedSecretResponses: e.sharedSecretProcessor.ProcessNetworkSecretMsgs(ctx, br)}
	return bsr, nil
}

func (e *enclaveImpl) ingestL1Block(ctx context.Context, br *common.BlockAndReceipts) (*components.BlockIngestionType, error) {
	e.logger.Info("Start ingesting block", log.BlockHashKey, br.Block.Hash())
	ingestion, err := e.l1BlockProcessor.Process(ctx, br)
	if err != nil {
		// only warn for unexpected errors
		if errors.Is(err, errutil.ErrBlockAncestorNotFound) || errors.Is(err, errutil.ErrBlockAlreadyProcessed) {
			e.logger.Debug("Did not ingest block", log.ErrKey, err, log.BlockHashKey, br.Block.Hash())
		} else {
			e.logger.Warn("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, br.Block.Hash())
		}
		return nil, err
	}

	err = e.rollupConsumer.ProcessRollupsInBlock(ctx, br)
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

func (e *enclaveImpl) SubmitTx(ctx context.Context, encryptedTxParams common.EncryptedTx) (*responses.RawTx, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitTx with the enclave stopping"))
	}
	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedTxParams, rpc.SubmitTxValidate, rpc.SubmitTxExecute)
}

func (e *enclaveImpl) Validator() nodetype.ObsValidator {
	validator, ok := e.service.(nodetype.ObsValidator)
	if !ok {
		panic("enclave service is not a validator but validator was requested!")
	}
	return validator
}

func (e *enclaveImpl) Sequencer() nodetype.Sequencer {
	sequencer, ok := e.service.(nodetype.Sequencer)
	if !ok {
		panic("enclave service is not a sequencer but sequencer was requested!")
	}
	return sequencer
}

func (e *enclaveImpl) SubmitBatch(ctx context.Context, extBatch *common.ExtBatch) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested SubmitBatch with the enclave stopping"))
	}

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "SubmitBatch call completed.", log.BatchHashKey, extBatch.Hash())

	e.logger.Info("Received new p2p batch", log.BatchHeightKey, extBatch.Header.Number, log.BatchHashKey, extBatch.Hash(), "l1", extBatch.Header.L1Proof)
	seqNo := extBatch.Header.SequencerOrderNo.Uint64()
	if seqNo > common.L2GenesisSeqNo+1 {
		_, err := e.storage.FetchBatchBySeqNo(ctx, seqNo-1)
		if err != nil {
			return responses.ToInternalError(fmt.Errorf("could not find previous batch with seq: %d", seqNo-1))
		}
	}

	batch, err := core.ToBatch(extBatch, e.dataEncryptionService, e.dataCompressionService)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not convert batch. Cause: %w", err))
	}

	err = e.Validator().VerifySequencerSignature(batch)
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

	err = e.Validator().ExecuteStoredBatches(ctx)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not execute batches. Cause: %w", err))
	}

	return nil
}

func (e *enclaveImpl) CreateBatch(ctx context.Context, skipBatchIfEmpty bool) common.SystemError {
	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateBatch call ended")
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested CreateBatch with the enclave stopping"))
	}

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	err := e.Sequencer().CreateBatch(ctx, skipBatchIfEmpty)
	if err != nil {
		return responses.ToInternalError(err)
	}

	return nil
}

func (e *enclaveImpl) CreateRollup(ctx context.Context, fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "CreateRollup call ended")
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GenerateRollup with the enclave stopping"))
	}

	// todo - remove once the db operations are more atomic
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	if e.registry.HeadBatchSeq() == nil {
		return nil, responses.ToInternalError(fmt.Errorf("not initialised yet"))
	}

	rollup, err := e.Sequencer().CreateRollup(ctx, fromSeqNo)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return rollup, nil
}

// ObsCall handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_call)
func (e *enclaveImpl) ObsCall(ctx context.Context, encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested ObsCall with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.TenCallValidate, rpc.TenCallExecute)
}

func (e *enclaveImpl) GetTransactionCount(ctx context.Context, encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransactionCount with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetTransactionCountValidate, rpc.GetTransactionCountExecute)
}

func (e *enclaveImpl) GetTransaction(ctx context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (*responses.TxByHash, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransaction with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetTransactionValidate, rpc.GetTransactionExecute)
}

func (e *enclaveImpl) GetTransactionReceipt(ctx context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.TxReceipt, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransactionReceipt with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetTransactionReceiptValidate, rpc.GetTransactionReceiptExecute)
}

func (e *enclaveImpl) Attestation(ctx context.Context) (*common.AttestationReport, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested ObsCall with the enclave stopping"))
	}

	if e.enclaveKey == nil {
		return nil, responses.ToInternalError(fmt.Errorf("public key not initialized, we can't produce the attestation report"))
	}
	report, err := e.attestationProvider.GetReport(ctx, e.enclaveKey.PublicKeyBytes(), e.enclaveKey.EnclaveID(), e.config.HostAddress)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not produce remote report. Cause %w", err))
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret(ctx context.Context) (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GenerateSecret with the enclave stopping"))
	}

	secret := crypto.GenerateEntropy(e.logger)
	err := e.storage.StoreSecret(ctx, secret)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	encSec, err := crypto.EncryptSecret(e.enclaveKey.PublicKeyBytes(), secret, e.logger)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed to encrypt secret. Cause: %w", err))
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(ctx context.Context, s common.EncryptedSharedEnclaveSecret) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested InitEnclave with the enclave stopping"))
	}

	secret, err := crypto.DecryptSecret(s, e.enclaveKey.PrivateKey())
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

func (e *enclaveImpl) EnclaveID(context.Context) (common.EnclaveID, common.SystemError) {
	return e.enclaveKey.EnclaveID(), nil
}

// GetBalance handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_getBalance)
func (e *enclaveImpl) GetBalance(ctx context.Context, encryptedParams common.EncryptedParamsGetBalance) (*responses.Balance, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetBalance with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetBalanceValidate, rpc.GetBalanceExecute)
}

func (e *enclaveImpl) GetCode(ctx context.Context, address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetCode with the enclave stopping"))
	}

	stateDB, err := e.registry.GetBatchState(ctx, batchHash)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not create stateDB. Cause: %w", err))
	}
	return stateDB.GetCode(address), nil
}

func (e *enclaveImpl) Subscribe(ctx context.Context, id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested SubscribeForExecutedBatches with the enclave stopping"))
	}

	encodedSubscription, err := e.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	return e.subscriptionManager.AddSubscription(id, encodedSubscription)
}

func (e *enclaveImpl) Unsubscribe(id gethrpc.ID) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Unsubscribe with the enclave stopping"))
	}

	e.subscriptionManager.RemoveSubscription(id)
	return nil
}

func (e *enclaveImpl) Stop() common.SystemError {
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

// EstimateGas decrypts CallMsg data, runs the gas estimation for the data.
// Using the callMsg.From Viewing Key, returns the encrypted gas estimation
func (e *enclaveImpl) EstimateGas(ctx context.Context, encryptedParams common.EncryptedParamsEstimateGas) (*responses.Gas, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested EstimateGas with the enclave stopping"))
	}

	defer core.LogMethodDuration(e.logger, measure.NewStopwatch(), "enclave.go:EstimateGas()")
	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.EstimateGasValidate, rpc.EstimateGasExecute)
}

func (e *enclaveImpl) GetLogs(ctx context.Context, encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetLogs with the enclave stopping"))
	}
	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetLogsValidate, rpc.GetLogsExecute)
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveImpl) HealthCheck(ctx context.Context) (bool, common.SystemError) {
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

func (e *enclaveImpl) DebugTraceTransaction(ctx context.Context, txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested DebugTraceTransaction with the enclave stopping"))
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, responses.ToInternalError(fmt.Errorf("debug namespace not enabled"))
	}

	jsonMsg, err := e.debugger.DebugTraceTransaction(ctx, txHash, config)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		// TODO *Pedro* MOVE THIS TO Enclave Response
		return json.RawMessage(err.Error()), nil
	}

	return jsonMsg, nil
}

func (e *enclaveImpl) DebugEventLogRelevancy(ctx context.Context, txHash gethcommon.Hash) (json.RawMessage, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested DebugEventLogRelevancy with the enclave stopping"))
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, responses.ToInternalError(fmt.Errorf("debug namespace not enabled"))
	}

	jsonMsg, err := e.debugger.DebugEventLogRelevancy(ctx, txHash)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		// TODO *Pedro* MOVE THIS TO Enclave Response
		return json.RawMessage(err.Error()), nil
	}

	return jsonMsg, nil
}

func (e *enclaveImpl) GetTotalContractCount(ctx context.Context) (*big.Int, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTotalContractCount with the enclave stopping"))
	}

	return e.storage.GetContractCount(ctx)
}

func (e *enclaveImpl) GetCustomQuery(ctx context.Context, encryptedParams common.EncryptedParamsGetStorageAt) (*responses.PrivateQueryResponse, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetReceiptsByAddress with the enclave stopping"))
	}

	return rpc.WithVKEncryption(ctx, e.rpcEncryptionManager, encryptedParams, rpc.GetCustomQueryValidate, rpc.GetCustomQueryExecute)
}

func (e *enclaveImpl) EnclavePublicConfig(context.Context) (*common.EnclavePublicConfig, common.SystemError) {
	address, systemError := e.crossChainProcessors.GetL2MessageBusAddress()
	if systemError != nil {
		return nil, systemError
	}
	return &common.EnclavePublicConfig{L2MessageBusAddress: address}, nil
}

func (e *enclaveImpl) rejectBlockErr(ctx context.Context, cause error) *errutil.BlockRejectError {
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

// this function looks at the batch chain and makes sure the resulting stateDB snapshots are available, replaying them if needed
// (if there had been a clean shutdown and all stateDB data was persisted this should do nothing)
func restoreStateDBCache(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, producer components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
	if registry.HeadBatchSeq() == nil {
		// not initialised yet
		return nil
	}
	batch, err := storage.FetchBatchBySeqNo(ctx, registry.HeadBatchSeq().Uint64())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// there is no head batch, this is probably a new node - there is no state to rebuild
			logger.Info("no head batch found in DB after restart", log.ErrKey, err)
			return nil
		}
		return fmt.Errorf("unexpected error fetching head batch to resync- %w", err)
	}
	if !stateDBAvailableForBatch(ctx, registry, batch.Hash()) {
		logger.Info("state not available for latest batch after restart - rebuilding stateDB cache from batches")
		err = replayBatchesToValidState(ctx, storage, registry, producer, gen, logger)
		if err != nil {
			return fmt.Errorf("unable to replay batches to restore valid state - %w", err)
		}
	}
	return nil
}

// The enclave caches a stateDB instance against each batch hash, this is the input state when producing the following
// batch in the chain and is used to query state at a certain height.
//
// This method checks if the stateDB data is available for a given batch hash (so it can be restored if not)
func stateDBAvailableForBatch(ctx context.Context, registry components.BatchRegistry, hash common.L2BatchHash) bool {
	_, err := registry.GetBatchState(ctx, &hash)
	return err == nil
}

// replayBatchesToValidState is used to repopulate the stateDB cache with data from persisted batches. Two step process:
// 1. step backwards from head batch until we find a batch that is already in stateDB cache, builds list of batches to replay
// 2. iterate that list of batches from the earliest, process the transactions to calculate and cache the stateDB
// todo (#1416) - get unit test coverage around this (and L2 Chain code more widely, see ticket #1416 )
func replayBatchesToValidState(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, batchExecutor components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
	// this slice will be a stack of batches to replay as we walk backwards in search of latest valid state
	// todo - consider capping the size of this batch list using FIFO to avoid memory issues, and then repeating as necessary
	var batchesToReplay []*core.Batch
	// `batchToReplayFrom` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	batchToReplayFrom, err := storage.FetchBatchBySeqNo(ctx, registry.HeadBatchSeq().Uint64())
	if err != nil {
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(ctx, registry, batchToReplayFrom.Hash()) {
		batchesToReplay = append(batchesToReplay, batchToReplayFrom)
		if batchToReplayFrom.NumberU64() == 0 {
			// no more parents to check, replaying from genesis
			break
		}
		batchToReplayFrom, err = storage.FetchBatch(ctx, batchToReplayFrom.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
	}
	logger.Info("replaying batch data into stateDB cache", "fromBatch", batchesToReplay[len(batchesToReplay)-1].NumberU64(),
		"toBatch", batchesToReplay[0].NumberU64())
	// loop through the slice of batches without stateDB data to cache the state (loop in reverse because slice is newest to oldest)
	for i := len(batchesToReplay) - 1; i >= 0; i-- {
		batch := batchesToReplay[i]

		// if genesis batch then create the genesis state before continuing on with remaining batches
		if batch.NumberU64() == 0 {
			err := gen.CommitGenesisState(storage)
			if err != nil {
				return err
			}
			continue
		}

		// calculate the stateDB after this batch and store it in the cache
		_, err := batchExecutor.ExecuteBatch(ctx, batch)
		if err != nil {
			return err
		}
	}

	return nil
}
