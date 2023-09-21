package enclave

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/gas"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/obscuronet/go-obscuro/go/enclave/vkhandler"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/nodetype"

	"github.com/obscuronet/go-obscuro/go/enclave/l2chain"
	"github.com/obscuronet/go-obscuro/go/responses"

	"github.com/obscuronet/go-obscuro/go/enclave/genesis"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/common/syserr"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/debugger"
	"github.com/obscuronet/go-obscuro/go/enclave/events"

	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	_ "github.com/obscuronet/go-obscuro/go/common/tracers/native" // make sure the tracers are loaded

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

var _noHeadBatch = big.NewInt(0)

type enclaveImpl struct {
	config                *config.EnclaveConfig
	storage               storage.Storage
	blockResolver         storage.BlockResolver
	l1BlockProcessor      components.L1BlockProcessor
	rollupConsumer        components.RollupConsumer
	l1Blockchain          *gethcore.BlockChain
	rpcEncryptionManager  rpc.EncryptionManager
	subscriptionManager   *events.SubscriptionManager
	crossChainProcessors  *crosschain.Processors
	sharedSecretProcessor *components.SharedSecretProcessor

	chain    l2chain.ObscuroChain
	service  nodetype.NodeType
	registry components.BatchRegistry

	// todo (#627) - use the ethconfig.Config instead
	GlobalGasCap uint64   //         5_000_000_000, // todo (#627) - make config
	BaseFee      *big.Int //              gethcommon.Big0,

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	attestationProvider components.AttestationProvider // interface for producing attestation reports and verifying them

	enclaveKey    *ecdsa.PrivateKey // this is a key specific to this enclave, which is included in the Attestation. Used for signing rollups and for encryption of the shared secret.
	enclavePubKey []byte            // the public key of the above

	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
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

	zeroTimestamp := uint64(0)
	// Initialise the database
	chainConfig := params.ChainConfig{
		ChainID:             big.NewInt(config.ObscuroChainID),
		HomesteadBlock:      gethcommon.Big0,
		DAOForkBlock:        gethcommon.Big0,
		EIP150Block:         gethcommon.Big0,
		EIP155Block:         gethcommon.Big0,
		EIP158Block:         gethcommon.Big0,
		ByzantiumBlock:      gethcommon.Big0,
		ConstantinopleBlock: gethcommon.Big0,
		PetersburgBlock:     gethcommon.Big0,
		IstanbulBlock:       gethcommon.Big0,
		MuirGlacierBlock:    gethcommon.Big0,
		BerlinBlock:         gethcommon.Big0,
		LondonBlock:         gethcommon.Big0,

		CancunTime:   &zeroTimestamp,
		ShanghaiTime: &zeroTimestamp,
		PragueTime:   &zeroTimestamp,
		VerkleTime:   &zeroTimestamp,
	}
	storage := storage.NewStorageFromConfig(config, &chainConfig, logger)

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
	enclaveKey, err := storage.GetEnclaveKey()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			logger.Crit("Failed to fetch enclave key", log.ErrKey, err)
		}
		// enclave key not found - new key should be generated
		// todo (#1053) - revisit the crypto for this key generation/lifecycle before production
		logger.Info("Generating the Obscuro key")
		enclaveKey, err = gethcrypto.GenerateKey()
		if err != nil {
			logger.Crit("Failed to generate enclave key.", log.ErrKey, err)
		}
		err = storage.StoreEnclaveKey(enclaveKey)
		if err != nil {
			logger.Crit("Failed to store enclave key.", log.ErrKey, err)
		}
	}

	serializedEnclavePubKey := gethcrypto.CompressPubkey(&enclaveKey.PublicKey)
	logger.Info(fmt.Sprintf("Generated public key %s", gethcommon.Bytes2Hex(serializedEnclavePubKey)))

	obscuroKey := crypto.GetObscuroKey(logger)
	rpcEncryptionManager := rpc.NewEncryptionManager(ecies.ImportECDSA(obscuroKey))

	dataEncryptionService := crypto.NewDataEncryptionService(logger)
	dataCompressionService := compression.NewBrotliDataCompressionService()

	memp := mempool.New(config.ObscuroChainID, logger)

	crossChainProcessors := crosschain.New(&config.MessageBusAddress, storage, big.NewInt(config.ObscuroChainID), logger)

	subscriptionManager := events.NewSubscriptionManager(&rpcEncryptionManager, storage, logger)

	gasOracle := gas.NewGasOracle()
	blockProcessor := components.NewBlockProcessor(storage, crossChainProcessors, gasOracle, logger)
	batchExecutor := components.NewBatchExecutor(storage, crossChainProcessors, genesis, gasOracle, &chainConfig, logger)
	sigVerifier, err := components.NewSignatureValidator(config.SequencerID, storage)
	registry := components.NewBatchRegistry(storage, logger)
	rProducer := components.NewRollupProducer(config.SequencerID, dataEncryptionService, config.ObscuroChainID, config.L1ChainID, storage, registry, blockProcessor, logger)
	if err != nil {
		logger.Crit("Could not initialise the signature validator", log.ErrKey, err)
	}
	rollupCompression := components.NewRollupCompression(registry, batchExecutor, dataEncryptionService, dataCompressionService, storage, &chainConfig, logger)
	rConsumer := components.NewRollupConsumer(mgmtContractLib, registry, rollupCompression, storage, logger, sigVerifier)
	sharedSecretProcessor := components.NewSharedSecretProcessor(mgmtContractLib, attestationProvider, storage, logger)

	var service nodetype.NodeType
	if config.NodeType == common.Sequencer {
		service = nodetype.NewSequencer(
			blockProcessor,
			batchExecutor,
			registry,
			rProducer,
			rConsumer,
			rollupCompression,
			logger,
			config.HostID,
			&chainConfig,
			enclaveKey,
			memp,
			storage,
			dataEncryptionService,
			dataCompressionService,
			nodetype.SequencerSettings{
				MaxBatchSize:      config.MaxBatchSize,
				MaxRollupSize:     config.MaxRollupSize,
				GasPaymentAddress: config.GasPaymentAddress,
				BatchGasLimit:     config.GasLimit,
				BaseFee:           config.BaseFee,
			},
		)
	} else {
		service = nodetype.NewValidator(blockProcessor, batchExecutor, registry, rConsumer, &chainConfig, config.SequencerID, storage, sigVerifier, logger)
	}

	chain := l2chain.NewChain(
		storage,
		&chainConfig,
		genesis,
		logger,
		registry,
	)

	// ensure cached chain state data is up-to-date using the persisted batch data
	err = restoreStateDBCache(storage, batchExecutor, genesis, logger)
	if err != nil {
		logger.Crit("failed to resync L2 chain state DB after restart", log.ErrKey, err)
	}

	// TODO ensure debug is allowed/disallowed
	debug := debugger.New(chain, storage, &chainConfig)

	logger.Info("Enclave service created with following config", log.CfgKey, config.HostID)
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
		enclavePubKey:          serializedEnclavePubKey,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		profiler:               prof,
		logger:                 logger,
		debugger:               debug,
		stopControl:            stopcontrol.New(),

		chain:    chain,
		registry: registry,
		service:  service,

		GlobalGasCap: 5_000_000_000, // todo (#627) - make config
		BaseFee:      gethcommon.Big0,

		mainMutex: sync.Mutex{},
	}
}

func (e *enclaveImpl) GetBatch(hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
	batch, err := e.storage.FetchBatch(hash)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed getting batch. Cause: %w", err))
	}

	b, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return b, nil
}

func (e *enclaveImpl) GetBatchBySeqNo(seqNo uint64) (*common.ExtBatch, common.SystemError) {
	batch, err := e.storage.FetchBatchBySeqNo(seqNo)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed getting batch. Cause: %w", err))
	}

	b, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return b, nil
}

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status() (common.Status, common.SystemError) {
	if e.stopControl.IsStopping() {
		return common.Status{StatusCode: common.Unavailable}, responses.ToInternalError(fmt.Errorf("requested Status with the enclave stopping"))
	}

	_, err := e.storage.FetchSecret()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return common.Status{StatusCode: common.AwaitingSecret, L2Head: _noHeadBatch}, nil
		}
		return common.Status{StatusCode: common.Unavailable}, responses.ToInternalError(err)
	}
	var l1HeadHash gethcommon.Hash
	l1Head, err := e.storage.FetchHeadBlock()
	if err != nil {
		// this might be normal while enclave is starting up, just send empty hash
		e.logger.Debug("failed to fetch L1 head block for status response", log.ErrKey, err)
	} else {
		l1HeadHash = l1Head.Hash()
	}
	// we use zero when there's no head batch yet, the first seq number is 1
	l2HeadSeqNo := _noHeadBatch
	// this is the highest seq number that has been received and stored on the enclave (it may not have been executed)
	currSeqNo, err := e.storage.FetchCurrentSequencerNo()
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
	e.logger.Info("Streaming batch to client", log.BatchHashKey, batch.Hash())
	extBatch, err := batch.ToExtBatch(e.dataEncryptionService, e.dataCompressionService)
	if err != nil {
		e.logger.Crit("failed to convert batch", log.ErrKey, err)
	}
	resp := common.StreamL2UpdatesResponse{
		Batch: extBatch,
	}
	outChannel <- resp
}

// this function is only called when the executed batch is the new head
func (e *enclaveImpl) streamEventsForNewHeadBatch(batch *core.Batch, receipts types.Receipts, outChannel chan common.StreamL2UpdatesResponse) {
	logs, err := e.subscriptionManager.GetSubscribedLogsForBatch(batch, receipts)
	e.logger.Info("Stream Events for", log.BatchHashKey, batch.Hash(), "nr_events", len(logs))
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
			e.streamEventsForNewHeadBatch(batch, receipts, l2UpdatesChannel)
		}
	})

	return l2UpdatesChannel, func() {
		e.registry.UnsubscribeFromBatches()
	}
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitL1Block(block types.Block, receipts types.Receipts, _ bool) (*common.BlockSubmissionResponse, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitL1Block with the enclave stopping"))
	}

	e.logger.Info("SubmitL1Block", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash())

	// If the block and receipts do not match, reject the block.
	br, err := common.ParseBlockAndReceipts(&block, &receipts)
	if err != nil {
		return nil, e.rejectBlockErr(fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	result, err := e.ingestL1Block(br)
	if err != nil {
		return nil, e.rejectBlockErr(fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	if result.IsFork() {
		e.logger.Info(fmt.Sprintf("Detected fork at block %s with height %d", block.Hash(), block.Number()))
	}

	err = e.service.OnL1Block(block, result)
	if err != nil {
		return nil, e.rejectBlockErr(fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	bsr := &common.BlockSubmissionResponse{ProducedSecretResponses: e.sharedSecretProcessor.ProcessNetworkSecretMsgs(br)}
	return bsr, nil
}

func (e *enclaveImpl) ingestL1Block(br *common.BlockAndReceipts) (*components.BlockIngestionType, error) {
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	ingestion, err := e.l1BlockProcessor.Process(br)
	if err != nil {
		// only warn for unexpected errors
		if errors.Is(err, errutil.ErrBlockAncestorNotFound) || errors.Is(err, errutil.ErrBlockAlreadyProcessed) {
			e.logger.Debug("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, br.Block.Hash())
		} else {
			e.logger.Warn("Failed ingesting block", log.ErrKey, err, log.BlockHashKey, br.Block.Hash())
		}
		return nil, err
	}

	err = e.rollupConsumer.ProcessRollupsInBlock(br)
	if err != nil && !errors.Is(err, components.ErrDuplicateRollup) {
		e.logger.Error("Encountered error while processing l1 block", log.ErrKey, err)
		// Unsure what to do here; block has been stored
	}

	if ingestion.IsFork() {
		err := e.service.OnL1Fork(ingestion.ChainFork)
		if err != nil {
			return nil, err
		}
	}
	return ingestion, nil
}

func (e *enclaveImpl) SubmitTx(tx common.EncryptedTx) (*responses.RawTx, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitTx with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(tx)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_call params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Transaction]
	if len(paramList) != 2 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}

	decryptedTx, err := rpc.ExtractTx(paramList[1].(string))
	if err != nil {
		e.logger.Info("could not decrypt transaction. ", log.ErrKey, err)
		return responses.AsPlaintextError(fmt.Errorf("could not decrypt transaction. Cause: %w", err)), nil
	}

	e.logger.Debug("Submitted transaction", log.TxKey, decryptedTx.Hash())

	viewingKeyAddress, err := rpc.GetSender(decryptedTx)
	if err != nil {
		if errors.Is(err, types.ErrInvalidSig) {
			return responses.AsPlaintextError(fmt.Errorf("transaction contains invalid signature")), nil
		}
		return responses.AsPlaintextError(fmt.Errorf("could not recover from address. Cause: %w", err)), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(&viewingKeyAddress, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	if e.crossChainProcessors.Local.IsSyntheticTransaction(*decryptedTx) {
		return responses.AsPlaintextError(responses.ToInternalError(fmt.Errorf("synthetic transaction coming from external rpc"))), nil
	}
	if err = e.checkGas(decryptedTx); err != nil {
		e.logger.Info("gas check failed", log.ErrKey, err.Error())
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	if err = e.service.SubmitTransaction(decryptedTx); err != nil {
		e.logger.Debug("Could not submit transaction", log.TxKey, decryptedTx.Hash(), log.ErrKey, err)
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	hash := decryptedTx.Hash().Hex()
	return responses.AsEncryptedResponse(&hash, vkHandler), nil
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

func (e *enclaveImpl) SubmitBatch(extBatch *common.ExtBatch) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested SubmitBatch with the enclave stopping"))
	}

	callStart := time.Now()
	defer func() {
		e.logger.Info("SubmitBatch call completed.", "start", callStart, log.DurationKey, time.Since(callStart), log.BatchHashKey, extBatch.Hash())
	}()

	e.logger.Info("SubmitBatch", log.BatchHeightKey, extBatch.Header.Number, log.BatchHashKey, extBatch.Hash(), "l1", extBatch.Header.L1Proof)
	batch, err := core.ToBatch(extBatch, e.dataEncryptionService, e.dataCompressionService)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not convert batch. Cause: %w", err))
	}

	err = e.Validator().VerifySequencerSignature(batch)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("invalid batch received. Could not verify signature. Cause: %w", err))
	}

	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	// if the signature is valid, then store the batch
	err = e.storage.StoreBatch(batch)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not store batch. Cause: %w", err))
	}

	err = e.Validator().ExecuteStoredBatches()
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not execute batches. Cause: %w", err))
	}

	return nil
}

func (e *enclaveImpl) CreateBatch() common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested CreateBatch with the enclave stopping"))
	}

	callStart := time.Now()
	defer func() {
		e.logger.Info("CreateBatch call ended", log.DurationMilliKey, time.Since(callStart).Milliseconds())
	}()

	// todo - remove once the db operations are more atomic
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	err := e.Sequencer().CreateBatch()
	if err != nil {
		return responses.ToInternalError(err)
	}

	return nil
}

func (e *enclaveImpl) CreateRollup(fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GenerateRollup with the enclave stopping"))
	}

	callStart := time.Now()
	defer func() {
		e.logger.Info(fmt.Sprintf("CreateRollup call ended - start = %s duration %s", callStart.String(), time.Since(callStart).String()))
	}()

	// todo - remove once the db operations are more atomic
	e.mainMutex.Lock()
	defer e.mainMutex.Unlock()

	rollup, err := e.Sequencer().CreateRollup(fromSeqNo)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	return rollup, nil
}

// ObsCall handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_call)
func (e *enclaveImpl) ObsCall(encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested ObsCall with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_call params - %w", err)), nil
	}

	// Parameters are [ViewingKey, TransactionArgs, BlockNumber]
	if len(paramList) != 3 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}

	apiArgs, err := gethencoding.ExtractEthCall(paramList[1])
	if err != nil {
		err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return responses.AsPlaintextError(err), nil
	}

	// encryption will fail if no From address is provided
	if apiArgs.From == nil {
		err = fmt.Errorf("no from address provided")
		return responses.AsPlaintextError(err), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(apiArgs.From, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(paramList[2])
	if err != nil {
		err = fmt.Errorf("unable to extract requested block number - %w", err)
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	execResult, err := e.chain.ObsCall(apiArgs, blkNumber)
	if err != nil {
		e.logger.Debug("Failed eth_call.", log.ErrKey, err)

		// make sure it's not some internal error
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}

		// make sure to serialize any possible EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	// encrypt the result payload
	var encodedResult string
	if len(execResult.ReturnData) != 0 {
		encodedResult = hexutil.Encode(execResult.ReturnData)
	}

	e.logger.Info("Call result success ", "result", encodedResult)

	return responses.AsEncryptedResponse(&encodedResult, vkHandler), nil
}

func (e *enclaveImpl) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransactionCount with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_transactionCount params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Address]
	if len(paramList) < 2 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}
	addressStr, ok := paramList[1].(string)
	if !ok {
		return responses.AsPlaintextError(fmt.Errorf("unexpected address parameter")), nil
	}

	address := gethcommon.HexToAddress(addressStr)

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(&address, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	var nonce uint64
	l2Head, err := e.storage.FetchHeadBatch()
	if err == nil {
		// todo - we should return an error when head state is not available, but for current test situations with race
		//  conditions we allow it to return zero while head state is uninitialized
		s, err := e.storage.CreateStateDB(l2Head.Hash())
		if err != nil {
			return nil, responses.ToInternalError(err)
		}
		nonce = s.GetNonce(address)
	}

	encoded := hexutil.EncodeUint64(nonce)
	return responses.AsEncryptedResponse(&encoded, vkHandler), nil
}

func (e *enclaveImpl) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (*responses.TxByHash, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransaction with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_getTransaction params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Hash]
	if len(paramList) != 2 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}
	txHashStr, ok := paramList[1].(string)
	if !ok {
		return responses.AsPlaintextError(fmt.Errorf("unexpected tx hash parameter")), nil
	}
	txHash := gethcommon.HexToHash(txHashStr)

	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := e.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth, return an empty response when a not found tx is requested
			return responses.AsEmptyResponse(), nil
		}
		return responses.AsPlaintextError(err), nil
	}

	viewingKeyAddress, err := rpc.GetSender(tx)
	if err != nil {
		err = fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionByHash response. Cause: %w", err)
		return responses.AsPlaintextError(err), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(&viewingKeyAddress, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	// Unlike in the Geth impl, we hardcode the use of a London signer.
	// todo (#1553) - once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`
	signer := types.NewLondonSigner(tx.ChainId())
	rpcTx := newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer)

	return responses.AsEncryptedResponse(rpcTx, vkHandler), nil
}

func (e *enclaveImpl) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.TxReceipt, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTransactionReceipt with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_getTransaction params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Hash]
	if len(paramList) != 2 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}
	txHashStr, ok := paramList[1].(string)
	if !ok {
		return responses.AsPlaintextError(fmt.Errorf("unable to parse the tx hash")), nil
	}
	txHash := gethcommon.HexToHash(txHashStr)

	// todo - optimise these calls. This can be done with a single sql
	e.logger.Trace("Get receipt for ", "txHash", txHash)
	// We retrieve the transaction.
	tx, _, _, _, err := e.storage.GetTransaction(txHash) //nolint:dogsled
	if err != nil {
		e.logger.Trace("error getting tx ", "txHash", txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return responses.AsEmptyResponse(), nil
		}
		return responses.AsPlaintextError(err), nil
	}

	// We retrieve the sender's address.
	sender, err := rpc.GetSender(tx)
	if err != nil {
		e.logger.Trace("error getting sender tx ", "txHash", txHash, log.ErrKey, err)
		return responses.AsPlaintextError(fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionReceipt response. Cause: %w", err)), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(&sender, paramList[0])
	if err != nil {
		e.logger.Trace("error getting the vk ", "txHash", txHash, log.ErrKey, err)
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	// We retrieve the transaction receipt.
	txReceipt, err := e.storage.GetTransactionReceipt(txHash)
	if err != nil {
		e.logger.Trace("error getting tx receipt", "txHash", txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return responses.AsEmptyResponse(), nil
		}
		err = fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	// We filter out irrelevant logs.
	txReceipt.Logs, err = e.subscriptionManager.FilterLogsForReceipt(txReceipt, &sender)
	if err != nil {
		e.logger.Trace("error filter logs ", "txHash", txHash, log.ErrKey, err)
		return nil, responses.ToInternalError(err)
	}

	e.logger.Trace("Successfully retreived receipt for ", "txHash", txHash, "rec", txReceipt)

	return responses.AsEncryptedResponse(txReceipt, vkHandler), nil
}

func (e *enclaveImpl) Attestation() (*common.AttestationReport, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested ObsCall with the enclave stopping"))
	}

	if e.enclavePubKey == nil {
		return nil, responses.ToInternalError(fmt.Errorf("public key not initialized, we can't produce the attestation report"))
	}
	report, err := e.attestationProvider.GetReport(e.enclavePubKey, e.config.HostID, e.config.HostAddress)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not produce remote report. Cause %w", err))
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GenerateSecret with the enclave stopping"))
	}

	secret := crypto.GenerateEntropy(e.logger)
	err := e.storage.StoreSecret(secret)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	encSec, err := crypto.EncryptSecret(e.enclavePubKey, secret, e.logger)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("failed to encrypt secret. Cause: %w", err))
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s common.EncryptedSharedEnclaveSecret) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested InitEnclave with the enclave stopping"))
	}

	secret, err := crypto.DecryptSecret(s, e.enclaveKey)
	if err != nil {
		return responses.ToInternalError(err)
	}
	err = e.storage.StoreSecret(*secret)
	if err != nil {
		return responses.ToInternalError(fmt.Errorf("could not store secret. Cause: %w", err))
	}
	e.logger.Trace(fmt.Sprintf("Secret decrypted and stored. Secret: %v", secret))
	return nil
}

// GetBalance handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_getBalance)
func (e *enclaveImpl) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (*responses.Balance, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetBalance with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_getBalance params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Address, BlockNumber]
	if len(paramList) != 3 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}

	requestedAddress, err := gethencoding.ExtractAddress(paramList[1])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to extract requested address - %w", err)), nil
	}

	blockNumber, err := gethencoding.ExtractBlockNumber(paramList[2])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to extract requested block number - %w", err)), nil
	}

	// params are correct, fetch the balance of the requested address
	// If the accountAddress is a contract, encrypt with the address of the contract owner
	encryptAddress, balance, err := e.chain.GetBalance(*requestedAddress, blockNumber)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to get balance - %w", err)), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(encryptAddress, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	return responses.AsEncryptedResponse(balance, vkHandler), nil
}

func (e *enclaveImpl) GetCode(address gethcommon.Address, batchHash *common.L2BatchHash) ([]byte, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetCode with the enclave stopping"))
	}

	stateDB, err := e.storage.CreateStateDB(*batchHash)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not create stateDB. Cause: %w", err))
	}
	return stateDB.GetCode(address), nil
}

func (e *enclaveImpl) Subscribe(id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) common.SystemError {
	if e.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested SubscribeForExecutedBatches with the enclave stopping"))
	}

	return e.subscriptionManager.AddSubscription(id, encryptedSubscription)
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
			e.logger.Error("Could not profiler", log.ErrKey, err)
			return err
		}
	}

	if e.registry != nil {
		e.registry.UnsubscribeFromBatches()
	}

	time.Sleep(time.Second)
	err := e.storage.Close()
	if err != nil {
		e.logger.Error("Could not stop db", log.ErrKey, err)
		return err
	}

	return nil
}

// EstimateGas decrypts CallMsg data, runs the gas estimation for the data.
// Using the callMsg.From Viewing Key, returns the encrypted gas estimation
func (e *enclaveImpl) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (*responses.Gas, common.SystemError) {
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested EstimateGas with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_estimateGas params - %w", err)), nil
	}

	// Parameters are [ViewingKey, callMsg, block number (optional) ]
	if len(paramList) < 2 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}

	callMsg, err := gethencoding.ExtractEthCall(paramList[1])
	if err != nil {
		err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return responses.AsPlaintextError(err), nil
	}

	// encryption will fail if From address is not provided
	if callMsg.From == nil {
		err = fmt.Errorf("no from address provided")
		return responses.AsPlaintextError(err), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(callMsg.From, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	// extract optional block number - defaults to the latest block if not avail
	blockNumber, err := gethencoding.ExtractOptionalBlockNumber(paramList, 2)
	if err != nil {
		err = fmt.Errorf("unable to extract requested block number - %w", err)
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	executionGasEstimate, err := e.DoEstimateGas(callMsg, blockNumber, e.GlobalGasCap)
	if err != nil {
		err = fmt.Errorf("unable to estimate transaction - %w", err)

		// make sure it's not some internal error
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}

		// make sure to serialize any possible EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	return responses.AsEncryptedResponse(&executionGasEstimate, vkHandler), nil
}

func (e *enclaveImpl) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError) { //nolint
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetLogs with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_estimateGas params - %w", err)), nil
	}

	// Parameters are [ViewingKey, Filter, Address]
	if len(paramList) != 3 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}
	// We extract the arguments from the param bytes.
	filter, forAddress, err := extractGetLogsParams(paramList[1:])
	if err != nil {
		return responses.AsPlaintextError(err), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(forAddress, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	// todo logic to check that the filter is valid
	// can't have both from and blockhash
	// from <=to
	// todo (@stefan) - return user error
	if filter.BlockHash != nil && filter.FromBlock != nil {
		return responses.AsEncryptedError(fmt.Errorf("invalid filter. Cannot have both blockhash and fromBlock"), vkHandler), nil
	}

	from := filter.FromBlock
	if from != nil && from.Int64() < 0 {
		batch, err := e.storage.FetchHeadBatch()
		if err != nil {
			return responses.AsPlaintextError(fmt.Errorf("could not retrieve head batch. Cause: %w", err)), nil
		}
		from = batch.Number()
	}

	// Set from to the height of the block hash
	if from == nil && filter.BlockHash != nil {
		batch, err := e.storage.FetchBatchHeader(*filter.BlockHash)
		if err != nil {
			return nil, responses.ToInternalError(err)
		}
		from = batch.Number
	}

	to := filter.ToBlock
	// when to=="latest", don't filter on it
	if to != nil && to.Int64() < 0 {
		to = nil
	}

	if from != nil && to != nil && from.Cmp(to) > 0 {
		return responses.AsEncryptedError(fmt.Errorf("invalid filter. from (%d) > to (%d)", from, to), vkHandler), nil
	}

	// We retrieve the relevant logs that match the filter.
	filteredLogs, err := e.storage.FilterLogs(forAddress, from, to, nil, filter.Addresses, filter.Topics)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		err = fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
		return responses.AsEncryptedError(err, vkHandler), nil
	}

	return responses.AsEncryptedResponse(&filteredLogs, vkHandler), nil
}

// DoEstimateGas returns the estimation of minimum gas required to execute transaction
// This is a copy of https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1055
// there's a high complexity to the method due to geth business rules (which is mimic'd here)
// once the work of obscuro gas mechanics is established this method should be simplified
func (e *enclaveImpl) DoEstimateGas(args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, gasCap uint64) (hexutil.Uint64, common.SystemError) { //nolint: gocognit
	// Binary search the gas requirement, as it may be higher than the amount used
	var ( //nolint: revive
		lo  = params.TxGas - 1
		hi  uint64
		cap uint64 //nolint:predeclared
	)
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(gethcommon.Address)
	}
	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= params.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// todo (#627) - review this with the gas mechanics/tokenomics work
		/*
			//Retrieve the block to act as the gas ceiling
			block, err := b.BlockByNumberOrHash(ctx, blockNrOrHash)
			if err != nil {
				return 0, err
			}
			if block == nil {
				return 0, errors.New("block not found")
			}
			hi = block.GasLimit()
		*/
		hi = e.GlobalGasCap
	}
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = gethcommon.Big0
	}
	// Recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 { //nolint:nestif
		balance, err := e.chain.GetBalanceAtBlock(*args.From, blkNumber)
		if err != nil {
			return 0, fmt.Errorf("unable to fetch account balance - %w", err)
		}

		available := new(big.Int).Set(balance.ToInt())
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, feeCap)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			e.logger.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer.ToInt(), "maxFeePerGas", feeCap, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if gasCap != 0 && hi > gasCap {
		e.logger.Warn("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap
	}
	cap = hi //nolint: revive

	// Execute the binary search and hone in on an isGasEnough gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := e.isGasEnough(args, mid, blkNumber)
		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap { //nolint:nestif
		failed, result, err := e.isGasEnough(args, hi, blkNumber)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas { //nolint: errorlint
				if len(result.Revert()) > 0 {
					return 0, newRevertError(result)
				}
				return 0, result.Err
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hexutil.Uint64(hi), nil
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveImpl) HealthCheck() (bool, common.SystemError) {
	if e.stopControl.IsStopping() {
		return false, responses.ToInternalError(fmt.Errorf("requested HealthCheck with the enclave stopping"))
	}

	// check the storage health
	storageHealthy, err := e.storage.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Info("HealthCheck failed for the enclave storage", log.ErrKey, err)
		return false, nil
	}
	// todo (#1148) - enclave healthcheck operations
	enclaveHealthy := true
	return storageHealthy && enclaveHealthy, nil
}

func (e *enclaveImpl) DebugTraceTransaction(txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested DebugTraceTransaction with the enclave stopping"))
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, responses.ToInternalError(fmt.Errorf("debug namespace not enabled"))
	}

	jsonMsg, err := e.debugger.DebugTraceTransaction(context.Background(), txHash, config)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		// TODO *Pedro* MOVE THIS TO Enclave Response
		return json.RawMessage(err.Error()), nil
	}

	return jsonMsg, nil
}

func (e *enclaveImpl) DebugEventLogRelevancy(txHash gethcommon.Hash) (json.RawMessage, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested DebugEventLogRelevancy with the enclave stopping"))
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, responses.ToInternalError(fmt.Errorf("debug namespace not enabled"))
	}

	jsonMsg, err := e.debugger.DebugEventLogRelevancy(txHash)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		// TODO *Pedro* MOVE THIS TO Enclave Response
		return json.RawMessage(err.Error()), nil
	}

	return jsonMsg, nil
}

func (e *enclaveImpl) GetTotalContractCount() (*big.Int, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetTotalContractCount with the enclave stopping"))
	}

	return e.storage.GetContractCount()
}

func (e *enclaveImpl) GetCustomQuery(encryptedParams common.EncryptedParamsGetStorageAt) (*responses.PrivateQueryResponse, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetReceiptsByAddress with the enclave stopping"))
	}

	// decode the received request into a []interface
	paramList, err := e.decodeRequest(encryptedParams)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to decode eth_getStorageAt params - %w", err)), nil
	}

	// Parameters are [ViewingKey, PrivateCustomQueryHeader, PrivateCustomQueryArgs, null]
	if len(paramList) != 4 {
		return responses.AsPlaintextError(fmt.Errorf("unexpected number of parameters")), nil
	}

	privateCustomQuery, err := gethencoding.ExtractPrivateCustomQuery(paramList[1], paramList[2])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to extract query - %w", err)), nil
	}

	// extract, create and validate the VK encryption handler
	vkHandler, err := createVKHandler(&privateCustomQuery.Address, paramList[0])
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to create VK encryptor - %w", err)), nil
	}

	// params are correct, fetch the receipts of the requested address
	encryptReceipts, err := e.storage.GetReceiptsPerAddress(&privateCustomQuery.Address, &privateCustomQuery.Pagination)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to get storage - %w", err)), nil
	}

	receiptsCount, err := e.storage.GetReceiptsPerAddressCount(&privateCustomQuery.Address)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("unable to get storage - %w", err)), nil
	}

	return responses.AsEncryptedResponse(&common.PrivateQueryResponse{
		Receipts: encryptReceipts,
		Total:    receiptsCount,
	}, vkHandler), nil
}

func (e *enclaveImpl) GetPublicTransactionData(pagination *common.QueryPagination) (*common.TransactionListingResponse, common.SystemError) {
	// ensure the enclave is running
	if e.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested GetPublicTransactionData with the enclave stopping"))
	}

	paginatedData, err := e.storage.GetPublicTransactionData(pagination)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("unable to fetch data - %w", err))
	}

	// Todo eventually make this a cacheable method
	totalData, err := e.storage.GetPublicTransactionCount()
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("unable to fetch data - %w", err))
	}

	return &common.TransactionListingResponse{
		TransactionsData: paginatedData,
		Total:            totalData,
	}, nil
}

// Create a helper to check if a gas allowance results in an executable transaction
// isGasEnough returns whether the gaslimit should be raised, lowered, or if it was impossible to execute the message
func (e *enclaveImpl) isGasEnough(args *gethapi.TransactionArgs, gas uint64, blkNumber *gethrpc.BlockNumber) (bool, *gethcore.ExecutionResult, error) {
	args.Gas = (*hexutil.Uint64)(&gas)
	result, err := e.chain.ObsCallAtBlock(args, blkNumber)
	if err != nil {
		if errors.Is(err, gethcore.ErrIntrinsicGas) {
			return true, nil, nil // Special case, raise gas limit
		}
		return true, nil, err // Bail out
	}
	return result.Failed(), result, nil
}

func newRevertError(result *gethcore.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompasses an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

func (e *enclaveImpl) checkGas(tx *types.Transaction) error {
	txGasPrice := tx.GasPrice()
	if txGasPrice == nil {
		return fmt.Errorf("rejected transaction %s. No gas price was set", tx.Hash())
	}
	minGasPrice := e.config.MinGasPrice
	if txGasPrice.Cmp(minGasPrice) == -1 {
		return fmt.Errorf("rejected transaction %s. Gas price was only %d, wanted at least %d", tx.Hash(), txGasPrice, minGasPrice)
	}
	return nil
}

// Returns the params extracted from an eth_getLogs request.
func extractGetLogsParams(paramList []interface{}) (*filters.FilterCriteria, *gethcommon.Address, error) {
	// We extract the first param, the filter for the logs.
	// We marshal the filter criteria from a map to JSON, then back from JSON into a FilterCriteria. This is
	// because the filter criteria arrives as a map, and there is no way to convert it to a map directly into a
	// FilterCriteria.
	filterJSON, err := json.Marshal(paramList[0])
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}
	filter := filters.FilterCriteria{}
	err = filter.UnmarshalJSON(filterJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal filter criteria from JSON. Cause: %w", err)
	}

	// We extract the second param, the address the logs are for.
	forAddressHex, ok := paramList[1].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected second argument in GetLogs request to be of type string, but got %T", paramList[0])
	}
	forAddress := gethcommon.HexToAddress(forAddressHex)
	return &filter, &forAddress, nil
}

func (e *enclaveImpl) rejectBlockErr(cause error) *errutil.BlockRejectError {
	var hash common.L1BlockHash
	l1Head, err := e.storage.FetchHeadBlock()
	// todo - handle error
	if err == nil {
		hash = l1Head.Hash()
	}
	return &errutil.BlockRejectError{
		L1Head:  hash,
		Wrapped: cause,
	}
}

func (e *enclaveImpl) decodeRequest(tx []byte) ([]interface{}, error) {
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params - %w", err)
	}

	paramList, err := gethencoding.DecodeParamBytes(paramBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to decode params - %w", err)
	}

	return paramList, nil
}

func serializeEVMError(err error) ([]byte, error) {
	var errReturn interface{}

	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *errutil.EVMSerialisableError
	if ok := errors.As(err, &e); ok {
		errReturn = e
	} else {
		// it's a generic error, serialise it
		errReturn = &errutil.EVMSerialisableError{Err: err.Error()}
	}

	// serialise the error object returned by the evm into a json
	errSerializedBytes, marshallErr := json.Marshal(errReturn)
	if marshallErr != nil {
		return nil, marshallErr
	}
	return errSerializedBytes, nil
}

// this function looks at the batch chain and makes sure the resulting stateDB snapshots are available, replaying them if needed
// (if there had been a clean shutdown and all stateDB data was persisted this should do nothing)
func restoreStateDBCache(storage storage.Storage, producer components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
	batch, err := storage.FetchHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// there is no head batch, this is probably a new node - there is no state to rebuild
			logger.Info("no head batch found in DB after restart", log.ErrKey, err)
			return nil
		}
		return fmt.Errorf("unexpected error fetching head batch to resync- %w", err)
	}
	if !stateDBAvailableForBatch(storage, batch.Hash()) {
		logger.Info("state not available for latest batch after restart - rebuilding stateDB cache from batches")
		err = replayBatchesToValidState(storage, producer, gen, logger)
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
func stateDBAvailableForBatch(storage storage.Storage, hash common.L2BatchHash) bool {
	_, err := storage.CreateStateDB(hash)
	return err == nil
}

// replayBatchesToValidState is used to repopulate the stateDB cache with data from persisted batches. Two step process:
// 1. step backwards from head batch until we find a batch that is already in stateDB cache, builds list of batches to replay
// 2. iterate that list of batches from the earliest, process the transactions to calculate and cache the stateDB
// todo (#1416) - get unit test coverage around this (and L2 Chain code more widely, see ticket #1416 )
func replayBatchesToValidState(storage storage.Storage, batchExecutor components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
	// this slice will be a stack of batches to replay as we walk backwards in search of latest valid state
	// todo - consider capping the size of this batch list using FIFO to avoid memory issues, and then repeating as necessary
	var batchesToReplay []*core.Batch
	// `batchToReplayFrom` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	batchToReplayFrom, err := storage.FetchHeadBatch()
	if err != nil {
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(storage, batchToReplayFrom.Hash()) {
		batchesToReplay = append(batchesToReplay, batchToReplayFrom)
		if batchToReplayFrom.NumberU64() == 0 {
			// no more parents to check, replaying from genesis
			break
		}
		batchToReplayFrom, err = storage.FetchBatch(batchToReplayFrom.Header.ParentHash)
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
		_, err := batchExecutor.ExecuteBatch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

func createVKHandler(address *gethcommon.Address, vkIntf interface{}) (*vkhandler.VKHandler, error) {
	vkPubKeyHexBytes, accountSignatureHexBytes, err := gethencoding.ExtractViewingKey(vkIntf)
	if err != nil {
		return nil, fmt.Errorf("unable to decode viewing key - %w", err)
	}

	encryptor, err := vkhandler.New(address, vkPubKeyHexBytes, accountSignatureHexBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create vk encryption for request - %w", err)
	}
	return encryptor, nil
}
