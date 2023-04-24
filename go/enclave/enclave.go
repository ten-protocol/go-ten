package enclave

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/debugger"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"

	"github.com/obscuronet/go-obscuro/go/enclave/l2chain"
	"github.com/obscuronet/go-obscuro/go/responses"

	"github.com/obscuronet/go-obscuro/go/enclave/genesis"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/events"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupmanager"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	_ "github.com/obscuronet/go-obscuro/go/common/tracers/native" // make sure the tracers are loaded

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type enclaveImpl struct {
	config               config.EnclaveConfig
	storage              db.Storage
	blockResolver        db.BlockResolver
	mempool              mempool.Manager
	l1Blockchain         *gethcore.BlockChain
	rpcEncryptionManager rpc.EncryptionManager
	rollupManager        rollupmanager.RollupManager
	subscriptionManager  *events.SubscriptionManager
	crossChainProcessors *crosschain.Processors

	chain *l2chain.ObscuroChain

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	attestationProvider AttestationProvider // interface for producing attestation reports and verifying them

	enclaveKey    *ecdsa.PrivateKey // this is a key specific to this enclave, which is included in the Attestation. Used for signing rollups and for encryption of the shared secret.
	enclavePubKey []byte            // the public key of the above

	transactionBlobCrypto crypto.TransactionBlobCrypto
	profiler              *profiler.Profiler
	debugger              *debugger.Debugger
	logger                gethlog.Logger

	stopInterrupt *int32
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(
	config config.EnclaveConfig,
	genesis *genesis.Genesis,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	logger gethlog.Logger,
) common.Enclave {
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
	backingDB, err := db.CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
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
	}
	storage := db.NewStorage(backingDB, &chainConfig, logger)

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
	var attestationProvider AttestationProvider
	if config.WillAttest {
		attestationProvider = &EgoAttestationProvider{}
	} else {
		logger.Info("WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &DummyAttestationProvider{}
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

	transactionBlobCrypto := crypto.NewTransactionBlobCryptoImpl(logger)

	memp := mempool.New(config.ObscuroChainID)

	crossChainProcessors := crosschain.New(&config.MessageBusAddress, storage, big.NewInt(config.ObscuroChainID), logger)

	subscriptionManager := events.NewSubscriptionManager(&rpcEncryptionManager, storage, logger)
	chain := l2chain.New(
		config.HostID,
		config.NodeType,
		storage,
		l1Blockchain,
		crossChainProcessors,
		memp,
		enclaveKey,
		&chainConfig,
		config.SequencerID,
		genesis,
		logger,
	)

	rollupManager := rollupmanager.New(
		mgmtContractLib,
		transactionBlobCrypto,
		config.ObscuroChainID,
		config.L1ChainID,
		storage,
		chain,
		logger,
	)

	// ensure cached chain state data is up-to-date using the persisted batch data
	err = chain.ResyncStateDB()
	if err != nil {
		logger.Crit("failed to resync L2 chain state DB after restart", log.ErrKey, err)
	}

	// TODO ensure debug is allowed/disallowed
	debug := debugger.New(chain, storage, &chainConfig)

	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Enclave service created with following config", log.CfgKey, string(jsonConfig))
	return &enclaveImpl{
		config:                config,
		storage:               storage,
		blockResolver:         storage,
		mempool:               memp,
		l1Blockchain:          l1Blockchain,
		rpcEncryptionManager:  rpcEncryptionManager,
		rollupManager:         rollupManager,
		subscriptionManager:   subscriptionManager,
		crossChainProcessors:  crossChainProcessors,
		chain:                 chain,
		mgmtContractLib:       mgmtContractLib,
		attestationProvider:   attestationProvider,
		enclaveKey:            enclaveKey,
		enclavePubKey:         serializedEnclavePubKey,
		transactionBlobCrypto: transactionBlobCrypto,
		profiler:              prof,
		logger:                logger,
		debugger:              debug,
		stopInterrupt:         new(int32),
	}
}

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status() (common.Status, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return common.Unavailable, nil
	}

	_, err := e.storage.FetchSecret()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return common.AwaitingSecret, nil
		}
		return common.Unavailable, err
	}
	return common.Running, nil // The enclave is local so it is always ready
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() error {
	return nil // The enclave is local so there is no client to stop
}

// SubmitL1Block is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil //nolint:nilnil
	}

	e.logger.Info("SubmitL1Block", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash())

	// If the block and receipts do not match, reject the block.
	br, err := common.ParseBlockAndReceipts(&block, &receipts, e.crossChainProcessors.Enabled())
	if err != nil {
		return nil, e.rejectBlockErr(fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}

	// We update the enclave state based on the L1 block.
	newL2Head, producedBatch, err := e.chain.ProcessL1Block(block, receipts, isLatest)
	if err != nil {
		e.logger.Warn("ProcessL1Block failed", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash(), log.ErrKey, err)
		return nil, e.rejectBlockErr(fmt.Errorf("could not submit L1 block. Cause: %w", err))
	}
	e.logger.Info("ProcessL1Block successful", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash())

	_, err = e.rollupManager.ProcessL1Block(br)
	if err != nil {
		e.logger.Error("Error processing L1 block for rollups", log.ErrKey, err)
	}

	// We prepare the block submission response.
	// todo (@stefan) - fix subscribed logs for validators who are being synchronized only through L1
	blockSubmissionResponse := e.produceBlockSubmissionResponse(newL2Head, producedBatch)

	if producedBatch != nil && (producedBatch.Header.Number.Uint64()%e.config.Cadence == 0) {
		rollup, err := e.rollupManager.CreateRollup()
		if err != nil {
			e.logger.Error("Failed to produce rollup", log.ErrKey, err)
		} else {
			blockSubmissionResponse.ProducedRollup = rollup.ToExtRollup(e.transactionBlobCrypto)
		}
	}

	e.logger.Info("produceBlockSubmissionResponse successful", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash(),
		"newBatch", describeBSR(blockSubmissionResponse))
	blockSubmissionResponse.ProducedSecretResponses = e.processNetworkSecretMsgs(br)

	// We remove any transactions considered immune to re-orgs from the mempool.
	if blockSubmissionResponse.ProducedBatch != nil {
		err = e.removeOldMempoolTxs(blockSubmissionResponse.ProducedBatch.Header)
		if err != nil {
			e.logger.Error("removeOldMempoolTxs fail", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash(), log.ErrKey, err)
			return nil, e.rejectBlockErr(fmt.Errorf("could not remove transactions from mempool. Cause: %w", err))
		}
	}

	return blockSubmissionResponse, nil
}

// useful description of the BSR for debugging
func describeBSR(response *common.BlockSubmissionResponse) string {
	if response.RejectError != nil {
		return fmt.Sprintf("BlockSubmissionResponse failed with err=%s", response.RejectError.Error())
	}
	producedBatch := "no batch produced"
	if response.ProducedBatch != nil {
		producedBatch = fmt.Sprintf("newBatch{num=%d, numTx=%d, hash=%s}",
			response.ProducedBatch.Header.Number, len(response.ProducedBatch.TxHashes), response.ProducedBatch.Hash())
	}
	producedRollup := "no rollup produced"
	if response.ProducedRollup != nil {
		producedRollup = fmt.Sprintf("newRollup{num=%d, numBatches=%d, hash=%s}",
			response.ProducedRollup.Header.Number, len(response.ProducedRollup.Batches), response.ProducedRollup.Hash())
	}
	return fmt.Sprintf("%s, %s", producedBatch, producedRollup)
}

func (e *enclaveImpl) SubmitTx(tx common.EncryptedTx) responses.RawTx {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested SubmitTx with the enclave stopping")))
	}

	encodedTx, err := e.rpcEncryptionManager.DecryptBytes(tx)
	if err != nil {
		err = fmt.Errorf("could not decrypt params in eth_sendRawTransaction request. Cause: %w", err)
		return e.handlePlainTextError(err)
	}
	decryptedTx, err := rpc.ExtractTx(encodedTx)
	if err != nil {
		e.logger.Info("could not decrypt transaction. ", log.ErrKey, err)
		return e.handlePlainTextError(fmt.Errorf("could not decrypt transaction. Cause: %w", err))
	}

	e.logger.Info(fmt.Sprintf("Submitted transaction = %s", decryptedTx.Hash().Hex()))

	viewingKeyAddress, err := rpc.GetSender(decryptedTx)
	if err != nil {
		if errors.Is(err, types.ErrInvalidSig) {
			return e.handlePlainTextError(fmt.Errorf("transaction contains invalid signature"))
		}
		return e.handlePlainTextError(fmt.Errorf("could not recover from address. Cause: %w", err))
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(viewingKeyAddress)

	if e.crossChainProcessors.Local.IsSyntheticTransaction(*decryptedTx) {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("synthetic transaction coming from external rpc")))
	}
	if err = e.checkGas(decryptedTx); err != nil {
		e.logger.Info("", log.ErrKey, err.Error())
		return e.handleEncryptedError(err, encryptor)
	}

	// Only the sequencer needs to maintain a transaction mempool. Other node types can return early.
	if e.config.NodeType == common.Sequencer {
		if err = e.mempool.AddMempoolTx(decryptedTx); err != nil {
			return e.handleEncryptedError(err, encryptor)
		}
	}

	hash := decryptedTx.Hash().Hex()
	return responses.AsEncryptedResponse(&hash, encryptor)
}

func (e *enclaveImpl) SubmitBatch(extBatch *common.ExtBatch) error {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil
	}

	e.logger.Info("SubmitBatch", "height", extBatch.Header.Number, "hash", extBatch.Hash(), "l1", extBatch.Header.L1Proof)
	batch := core.ToBatch(extBatch, e.transactionBlobCrypto)
	if err := e.chain.UpdateL2Chain(batch); err != nil {
		return fmt.Errorf("could not update L2 chain based on batch. Cause: %w", err)
	}

	return nil
}

func (e *enclaveImpl) GenerateRollup() (*common.ExtRollup, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil //nolint:nilnil
	}

	if e.config.NodeType != common.Sequencer {
		return nil, errors.New("only sequencer can generate rollups")
	}

	rollup, err := e.rollupManager.CreateRollup()
	if err != nil {
		return nil, err
	}

	return rollup.ToExtRollup(e.transactionBlobCrypto), nil
}

// ObsCall handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_call)
func (e *enclaveImpl) ObsCall(encryptedParams common.EncryptedParamsCall) responses.Call {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested ObsCall with the enclave stopping")))
	}

	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		err = fmt.Errorf("could not decrypt params in eth_call request. Cause: %w", err)
		return e.handlePlainTextError(err)
	}

	// extract params from byte slice to array of strings
	var paramList []interface{}
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		err = fmt.Errorf("unable to decode eth_call params - %w", err)
		return e.handlePlainTextError(err)
	}

	// params are [TransactionArgs, BlockNumber]
	if len(paramList) != 2 {
		err = fmt.Errorf("required exactly two params, but received %d", len(paramList))
		return e.handlePlainTextError(err)
	}

	apiArgs, err := gethencoding.ExtractEthCall(paramList[0])
	if err != nil {
		err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return e.handlePlainTextError(err)
	}

	// encryption will fail if no From address is provided
	if apiArgs.From == nil {
		err = fmt.Errorf("no from address provided")
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(*apiArgs.From)

	blkNumber, err := gethencoding.ExtractBlockNumber(paramList[1])
	if err != nil {
		err = fmt.Errorf("unable to extract requested block number - %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	execResult, err := e.chain.ObsCall(apiArgs, blkNumber)
	if err != nil {
		e.logger.Info("Could not execute off chain call.", log.ErrKey, err)

		// make sure it's not some internal error
		if errors.Is(err, errutil.SystemError{}) {
			return e.handlePlainTextError(err)
		}

		// make sure to serialize any possible EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		return e.handleEncryptedError(err, encryptor)
	}

	// encrypt the result payload
	var encodedResult string
	if len(execResult.ReturnData) != 0 {
		encodedResult = hexutil.Encode(execResult.ReturnData)
	}

	return responses.AsEncryptedResponse(&encodedResult, encryptor)
}

func (e *enclaveImpl) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) responses.TxCount {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested GetTransactionCount with the enclave stopping")))
	}

	var nonce uint64
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return e.handlePlainTextError(err)
	}

	address, err := rpc.ExtractAddress(paramBytes)
	if err != nil {
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(address)

	l2Head, err := e.storage.FetchHeadBatch()
	if err == nil {
		// todo - we should return an error when head state is not available, but for current test situations with race
		//  conditions we allow it to return zero while head state is uninitialized
		s, err := e.storage.CreateStateDB(*l2Head.Hash())
		if err != nil {
			return e.handlePlainTextError(fmt.Errorf("could not create stateDB. Cause: %w", err))
		}
		nonce = s.GetNonce(address)
	}

	encoded := hexutil.EncodeUint64(nonce)
	return responses.AsEncryptedResponse(&encoded, encryptor)
}

func (e *enclaveImpl) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) responses.TxByHash {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested GetTransaction with the enclave stopping")))
	}

	hashBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		err = fmt.Errorf("could not decrypt encrypted RPC request params. Cause: %w", err)
		return e.handlePlainTextError(err)
	}
	var paramList []string
	err = json.Unmarshal(hashBytes, &paramList)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
		return e.handlePlainTextError(err)
	}
	if len(paramList) == 0 {
		err = fmt.Errorf("required at least one param, but received zero")
		return e.handlePlainTextError(err)
	}
	txHash := gethcommon.HexToHash(paramList[0])

	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := e.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not found tx is requested
			return responses.AsEmptyResponse()
		}
		return e.handlePlainTextError(err)
	}

	viewingKeyAddress, err := rpc.GetSender(tx)
	if err != nil {
		err = fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionByHash response. Cause: %w", err)
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(viewingKeyAddress)

	// Unlike in the Geth impl, we hardcode the use of a London signer.
	// todo (#1553) - once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`
	signer := types.NewLondonSigner(tx.ChainId())
	rpcTx := newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer)

	return responses.AsEncryptedResponse(rpcTx, encryptor)
}

func (e *enclaveImpl) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) responses.TxReceipt {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested GetTransactionReceipt with the enclave stopping")))
	}

	// We decrypt the transaction bytes.
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return e.handlePlainTextError(fmt.Errorf("could not decrypt params in eth_getTransactionReceipt request. Cause: %w", err))
	}
	txHash, err := rpc.ExtractTxHash(paramBytes)
	if err != nil {
		return e.handlePlainTextError(err)
	}

	// We retrieve the transaction.
	tx, txBatchHash, txBatchHeight, _, err := e.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return responses.AsEmptyResponse()
		}
		return e.handlePlainTextError(err)
	}

	// We retrieve the sender's address.
	sender, err := rpc.GetSender(tx)
	if err != nil {
		return e.handlePlainTextError(fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionReceipt response. Cause: %w", err))
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(sender)

	// Only return receipts for transactions included in the canonical chain.
	r, err := e.storage.FetchBatchByHeight(txBatchHeight)
	if err != nil {
		err = fmt.Errorf("could not retrieve batch containing transaction. Cause: %w", err)
		return e.handlePlainTextError(err)
	}
	if !bytes.Equal(r.Hash().Bytes(), txBatchHash.Bytes()) {
		err = fmt.Errorf("transaction not included in the canonical chain")
		return e.handleEncryptedError(err, encryptor)
	}

	// We retrieve the transaction receipt.
	txReceipt, err := e.storage.GetTransactionReceipt(txHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth return an empty response when a not-found tx is requested
			return responses.AsEmptyResponse()
		}
		err := fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	// We filter out irrelevant logs.
	txReceipt.Logs, err = e.subscriptionManager.FilterLogs(txReceipt.Logs, txBatchHash, &sender, &filters.FilterCriteria{})
	if err != nil {
		return e.handleEncryptedError(fmt.Errorf("could not filter logs. Cause: %w", err), encryptor)
	}

	return responses.AsEncryptedResponse(txReceipt, encryptor)
}

func (e *enclaveImpl) Attestation() (*common.AttestationReport, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil //nolint:nilnil
	}

	if e.enclavePubKey == nil {
		e.logger.Error("public key not initialized, we can't produce the attestation report")
		return nil, fmt.Errorf("public key not initialized, we can't produce the attestation report")
	}
	report, err := e.attestationProvider.GetReport(e.enclavePubKey, e.config.HostID, e.config.HostAddress)
	if err != nil {
		e.logger.Error("could not produce remote report")
		return nil, fmt.Errorf("could not produce remote report")
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() (common.EncryptedSharedEnclaveSecret, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil
	}

	secret := crypto.GenerateEntropy(e.logger)
	err := e.storage.StoreSecret(secret)
	if err != nil {
		return nil, fmt.Errorf("could not store secret. Cause: %w", err)
	}
	encSec, err := crypto.EncryptSecret(e.enclavePubKey, secret, e.logger)
	if err != nil {
		e.logger.Error("failed to encrypt secret.", log.ErrKey, err)
		return nil, fmt.Errorf("failed to encrypt secret. Cause: %w", err)
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s common.EncryptedSharedEnclaveSecret) error {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil
	}

	secret, err := crypto.DecryptSecret(s, e.enclaveKey)
	if err != nil {
		return err
	}
	err = e.storage.StoreSecret(*secret)
	if err != nil {
		return fmt.Errorf("could not store secret. Cause: %w", err)
	}
	e.logger.Trace(fmt.Sprintf("Secret decrypted and stored. Secret: %v", secret))
	return nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (e *enclaveImpl) verifyAttestationAndEncryptSecret(att *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil
	}

	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := e.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, fmt.Errorf("unable to verify report - %w", err)
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = VerifyIdentity(data, att); err != nil {
		return nil, fmt.Errorf("unable to verify identity - %w", err)
	}
	e.logger.Info(fmt.Sprintf("Successfully verified attestation and identity. Owner: %s", att.Owner))

	secret, err := e.storage.FetchSecret()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve secret; this should not happen. Cause: %w", err)
	}
	return crypto.EncryptSecret(att.PubKey, *secret, e.logger)
}

func (e *enclaveImpl) AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil
	}

	return e.rpcEncryptionManager.AddViewingKey(encryptedViewingKeyBytes, signature)
}

// storeAttestation stores the attested keys of other nodes so we can decrypt their rollups
func (e *enclaveImpl) storeAttestation(att *common.AttestationReport) error {
	e.logger.Info(fmt.Sprintf("Store attestation. Owner: %s", att.Owner))
	// Store the attestation
	key, err := gethcrypto.DecompressPubkey(att.PubKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key %w", err)
	}
	err = e.storage.StoreAttestedKey(att.Owner, key)
	if err != nil {
		return fmt.Errorf("could not store attested key. Cause: %w", err)
	}
	return nil
}

// GetBalance handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_getBalance)
func (e *enclaveImpl) GetBalance(encryptedParams common.EncryptedParamsGetBalance) responses.Balance {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested GetBalance with the enclave stopping")))
	}

	// Decrypt the request.
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		err = fmt.Errorf("could not decrypt params in eth_getBalance request. Cause: %w", err)
		return e.handlePlainTextError(err)
	}

	// Extract the params from the request.
	var paramList []string
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
		return e.handlePlainTextError(err)
	}
	if len(paramList) != 2 {
		err = fmt.Errorf("required exactly two params, but received %d", len(paramList))
		return e.handlePlainTextError(err)
	}

	accountAddress, err := gethencoding.ExtractAddress(paramList[0])
	if err != nil {
		err = fmt.Errorf("unable to extract requested address - %w", err)
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(*accountAddress)

	blockNumber, err := gethencoding.ExtractBlockNumber(paramList[1])
	if err != nil {
		err = fmt.Errorf("unable to extract requested block number - %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	encryptAddress, balance, err := e.chain.GetBalance(*accountAddress, blockNumber)
	if err != nil {
		err = fmt.Errorf("unable to get balance - %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	encryptor = e.rpcEncryptionManager.CreateEncryptorFor(*encryptAddress)

	return responses.AsEncryptedResponse(balance, encryptor)
}

func (e *enclaveImpl) GetCode(address gethcommon.Address, batchHash *common.L2BatchHash) ([]byte, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil
	}

	stateDB, err := e.storage.CreateStateDB(*batchHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}
	return stateDB.GetCode(address), nil
}

func (e *enclaveImpl) Subscribe(id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil
	}

	return e.subscriptionManager.AddSubscription(id, encryptedSubscription)
}

func (e *enclaveImpl) Unsubscribe(id gethrpc.ID) error {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil
	}

	e.subscriptionManager.RemoveSubscription(id)
	return nil
}

func (e *enclaveImpl) Stop() error {
	// block all requests
	atomic.StoreInt32(e.stopInterrupt, 1)

	if e.profiler != nil {
		return e.profiler.Stop()
	}

	time.Sleep(time.Second)
	err := e.storage.Close()
	if err != nil {
		e.logger.Error("Could not stop db", log.ErrKey, err)
	}
	return nil
}

// EstimateGas decrypts CallMsg data, runs the gas estimation for the data.
// Using the callMsg.From Viewing Key, returns the encrypted gas estimation
func (e *enclaveImpl) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) responses.Gas {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested EstimateGas with the enclave stopping")))
	}

	// decrypt the input with the enclave PK
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		err = fmt.Errorf("unable to decrypt params in EstimateGas request. Cause: %w", err)
		return e.handlePlainTextError(err)
	}

	// extract params from byte slice to array of strings
	var paramList []interface{}
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		err = fmt.Errorf("unable to decode EthCall params - %w", err)
		return e.handlePlainTextError(err)
	}

	// params are [callMsg, block number (optional) ]
	if len(paramList) < 1 {
		err = fmt.Errorf("required at least 1 params, but received %d", len(paramList))
		return e.handlePlainTextError(err)
	}

	callMsg, err := gethencoding.ExtractEthCall(paramList[0])
	if err != nil {
		err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return e.handlePlainTextError(err)
	}

	// encryption will fail if From address is not provided
	if callMsg.From == nil {
		err = fmt.Errorf("no from address provided")
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(*callMsg.From)

	// extract optional block number - defaults to the latest block if not avail
	blockNumber, err := gethencoding.ExtractOptionalBlockNumber(paramList, 1)
	if err != nil {
		err = fmt.Errorf("unable to extract requested block number - %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	// todo (@pedro) - hook the correct blockNumber from the API call (paramList[1])
	gasEstimate, err := e.DoEstimateGas(callMsg, blockNumber, e.chain.GlobalGasCap)
	if err != nil {
		err = fmt.Errorf("unable to estimate transaction - %w", err)

		// make sure it's not some internal error
		if errors.Is(err, errutil.SystemError{}) {
			return e.handlePlainTextError(err)
		}

		// make sure to serialize any possible EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		return e.handleEncryptedError(err, encryptor)
	}

	return responses.AsEncryptedResponse(&gasEstimate, encryptor)
}

func (e *enclaveImpl) GetLogs(encryptedParams common.EncryptedParamsGetLogs) responses.Logs {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return e.handlePlainTextError(errutil.NewSystemErr(fmt.Errorf("requested GetLogs with the enclave stopping")))
	}

	// We decrypt the params.
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		err = fmt.Errorf("unable to decrypt params in GetLogs request. Cause: %w", err)
		return e.handlePlainTextError(err)
	}

	// We extract the arguments from the param bytes.
	filter, forAddress, err := extractGetLogsParams(paramBytes)
	if err != nil {
		return e.handlePlainTextError(err)
	}

	encryptor := e.rpcEncryptionManager.CreateEncryptorFor(*forAddress)
	// todo logic to check that the filter is valid
	// can't have both from and blockhash
	// from <=to
	// todo (@stefan) - return user error
	if filter.BlockHash != nil && filter.FromBlock != nil {
		return e.handleEncryptedError(fmt.Errorf("invalid filter. Cannot have both blockhash and fromBlock"), encryptor)
	}

	from := filter.FromBlock
	if from != nil && from.Int64() < 0 {
		batch, err := e.storage.FetchHeadBatch()
		if err != nil {
			return e.handlePlainTextError(fmt.Errorf("could not retrieve head batch. Cause: %w", err))
		}
		from = batch.Number()
	}

	// Set from to the height of the block hash
	if from == nil && filter.BlockHash != nil {
		batch, err := e.storage.FetchBatch(*filter.BlockHash)
		if err != nil {
			return e.handleEncryptedError(fmt.Errorf("could not retrieve batch with hash %s. Cause: %w", filter.BlockHash.Hex(), err), encryptor)
		}
		from = batch.Number()
	}

	to := filter.ToBlock
	// when to=="latest", don't filter on it
	if to != nil && to.Int64() < 0 {
		to = nil
	}

	if from != nil && to != nil && from.Cmp(to) > 0 {
		return e.handleEncryptedError(fmt.Errorf("invalid filter. from (%d) > to (%d)", from, to), encryptor)
	}

	// We retrieve the relevant logs that match the filter.
	filteredLogs, err := e.storage.FilterLogs(forAddress, from, to, nil, filter.Addresses, filter.Topics)
	if err != nil {
		err = fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
		return e.handleEncryptedError(err, encryptor)
	}

	return responses.AsEncryptedResponse(&filteredLogs, encryptor)
}

// DoEstimateGas returns the estimation of minimum gas required to execute transaction
// This is a copy of https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1055
// there's a high complexity to the method due to geth business rules (which is mimic'd here)
// once the work of obscuro gas mechanics is established this method should be simplified
func (e *enclaveImpl) DoEstimateGas(args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, gasCap uint64) (hexutil.Uint64, error) { //nolint: gocognit
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
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
		hi = e.chain.GlobalGasCap
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
	cap = hi

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
func (e *enclaveImpl) HealthCheck() (bool, error) {
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return false, nil
	}

	// check the storage health
	storageHealthy, err := e.storage.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Error("unable to HealthCheck enclave storage", log.ErrKey, err)
		return false, nil
	}
	// todo (#1148) - enclave healthcheck operations
	enclaveHealthy := true
	return storageHealthy && enclaveHealthy, nil
}

func (e *enclaveImpl) DebugTraceTransaction(txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, error) {
	// ensure the enclave is running
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, fmt.Errorf("debug namespace not enabled")
	}

	return e.debugger.DebugTraceTransaction(context.Background(), txHash, config)
}

func (e *enclaveImpl) DebugEventLogRelevancy(txHash gethcommon.Hash) (json.RawMessage, error) {
	// ensure the enclave is running
	if atomic.LoadInt32(e.stopInterrupt) == 1 {
		return nil, nil
	}

	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, fmt.Errorf("debug namespace not enabled")
	}

	return e.debugger.DebugEventLogRelevancy(txHash)
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

// processNetworkSecretMsgs we watch for all messages that are requesting or receiving the secret and we store the nodes attested keys
func (e *enclaveImpl) processNetworkSecretMsgs(br *common.BlockAndReceipts) []*common.ProducedSecretResponse {
	var responses []*common.ProducedSecretResponse
	transactions := br.SuccessfulTransactions()
	block := br.Block
	for _, tx := range *transactions {
		t := e.mgmtContractLib.DecodeTx(tx)

		// this transaction is for a node that has joined the network and needs to be sent the network secret
		if scrtReqTx, ok := t.(*ethadapter.L1RequestSecretTx); ok {
			e.logger.Info(fmt.Sprintf("Process shared secret request. Block: %d. TxKey: %d",
				block.NumberU64(), common.ShortHash(tx.Hash())))
			resp, err := e.processSecretRequest(scrtReqTx)
			if err != nil {
				e.logger.Error("Failed to process shared secret request.", log.ErrKey, err)
				continue
			}
			responses = append(responses, resp)
		}

		// this transaction was created by the genesis node, we need to store their attested key to decrypt their rollup
		if initSecretTx, ok := t.(*ethadapter.L1InitializeSecretTx); ok {
			// todo (#1580) - ensure that we don't accidentally skip over the real `L1InitializeSecretTx` message. Otherwise
			//  our node will never be able to speak to other nodes.
			// there must be a way to make sure that this transaction can only be sent once.
			att, err := common.DecodeAttestation(initSecretTx.Attestation)
			if err != nil {
				e.logger.Error("Could not decode attestation report", log.ErrKey, err)
			}

			err = e.storeAttestation(att)
			if err != nil {
				e.logger.Error("Could not store the attestation report.", log.ErrKey, err)
			}
		}
	}
	return responses
}

func (e *enclaveImpl) processSecretRequest(req *ethadapter.L1RequestSecretTx) (*common.ProducedSecretResponse, error) {
	att, err := common.DecodeAttestation(req.Attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation - %w", err)
	}

	e.logger.Info("received attestation", "attestation", att)
	secret, err := e.verifyAttestationAndEncryptSecret(att)
	if err != nil {
		return nil, fmt.Errorf("secret request failed, no response will be published - %w", err)
	}

	// Store the attested key only if the attestation process succeeded.
	err = e.storeAttestation(att)
	if err != nil {
		return nil, fmt.Errorf("could not store attestation, no response will be published. Cause: %w", err)
	}

	e.logger.Trace("Processed secret request.", "owner", att.Owner)
	return &common.ProducedSecretResponse{
		Secret:      secret,
		RequesterID: att.Owner,
		HostAddress: att.HostAddress,
	}, nil
}

// Returns the params extracted from an eth_getLogs request.
func extractGetLogsParams(paramBytes []byte) (*filters.FilterCriteria, *gethcommon.Address, error) {
	// We verify the params.
	var paramsList []interface{}
	err := json.Unmarshal(paramBytes, &paramsList)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode GetLogs params. Cause: %w", err)
	}
	if len(paramsList) != 2 {
		return nil, nil, fmt.Errorf("expected 2 params in GetLogs request, but received %d", len(paramsList))
	}

	// We extract the first param, the filter for the logs.
	// We marshal the filter criteria from a map to JSON, then back from JSON into a FilterCriteria. This is
	// because the filter criteria arrives as a map, and there is no way to convert it to a map directly into a
	// FilterCriteria.
	filterJSON, err := json.Marshal(paramsList[0])
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}
	filter := filters.FilterCriteria{}
	err = filter.UnmarshalJSON(filterJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal filter criteria from JSON. Cause: %w", err)
	}

	// We extract the second param, the address the logs are for.
	forAddressHex, ok := paramsList[1].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected second argument in GetLogs request to be of type string, but got %T", paramsList[0])
	}
	forAddress := gethcommon.HexToAddress(forAddressHex)
	return &filter, &forAddress, nil
}

// Removes transactions from the mempool that are considered immune to re-orgs (i.e. over X batches deep).
func (e *enclaveImpl) removeOldMempoolTxs(batchHeader *common.BatchHeader) error {
	if batchHeader == nil {
		return nil
	}

	hr, err := e.storage.FetchBatch(batchHeader.Hash())
	if err != nil {
		return fmt.Errorf("could not retrieve batch. This should not happen because this batch was just processed. Cause: %w", err)
	}
	err = e.mempool.RemoveMempoolTxs(hr, e.storage)
	if err != nil {
		return fmt.Errorf("could not remove transactions from mempool. Cause: %w", err)
	}

	return nil
}

func (e *enclaveImpl) produceBlockSubmissionResponse(l2Head *common.L2BatchHash, producedBatch *core.Batch) *common.BlockSubmissionResponse {
	if l2Head == nil {
		// not an error state, we ingested a block but no rollup head found
		return &common.BlockSubmissionResponse{}
	}

	var producedExtBatch *common.ExtBatch
	if producedBatch != nil {
		producedExtBatch = producedBatch.ToExtBatch(e.transactionBlobCrypto)
	}

	batch, err := e.storage.FetchBatch(*l2Head)
	if err != nil {
		e.logger.Crit("Failed to retrieve batch. Should not happen", log.ErrKey, err)
		return nil
	}
	logs, err := e.subscriptionLogs(batch.Number())
	if err != nil {
		e.logger.Error("Could not fetch logs", log.ErrKey, err)
		return nil
	}
	return &common.BlockSubmissionResponse{
		ProducedBatch:  producedExtBatch,
		SubscribedLogs: logs,
	}
}

// Retrieves and encrypts the logs for the block.
func (e *enclaveImpl) subscriptionLogs(upToBatchNr *big.Int) (map[gethrpc.ID][]byte, error) {
	result := map[gethrpc.ID][]*types.Log{}

	// Go through each subscription and collect the logs
	err := e.subscriptionManager.ForEachSubscription(func(id gethrpc.ID, subscription *common.LogSubscription, previousHead *big.Int) error {
		// 1. fetch the logs since the last request
		var from *big.Int
		to := upToBatchNr

		if previousHead == nil || previousHead.Int64() <= 0 {
			// when the subscription is initialised, default from the latest batch
			batch, err := e.storage.FetchHeadBatch()
			if err != nil {
				return err
			}
			from = batch.Number()
		} else {
			from = big.NewInt(previousHead.Int64() + 1)
		}

		if from.Cmp(to) > 0 {
			e.logger.Warn(fmt.Sprintf("Skipping subscription step id=%s: [%d, %d]", id, from, to))
			return nil
		}

		logs, err := e.storage.FilterLogs(subscription.Account, from, to, nil, subscription.Filter.Addresses, subscription.Filter.Topics)
		e.logger.Info(fmt.Sprintf("Subscription id=%s: [%d, %d]. Logs %d, Err: %s", id, from, to, len(logs), err))
		if err != nil {
			return err
		}

		// 2.  store the current l2Head in the Subscription
		e.subscriptionManager.SetLastHead(id, to)
		result[id] = logs
		return nil
	})
	if err != nil {
		e.logger.Error("Could not retrieve subscription logs", log.ErrKey, err)
		return nil, err
	}

	// Encrypt the results
	return e.subscriptionManager.EncryptLogs(result)
}

func (e *enclaveImpl) rejectBlockErr(cause error) *common.BlockRejectError {
	var hash common.L1BlockHash
	l1Head, err := e.storage.FetchHeadBlock()
	// todo - handle error
	if err == nil {
		hash = l1Head.Hash()
	}
	return &common.BlockRejectError{
		L1Head:  hash,
		Wrapped: cause,
	}
}

func serializeEVMError(err error) ([]byte, error) {
	var errReturn interface{}

	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *evm.SerialisableError
	if ok := errors.As(err, &e); ok {
		errReturn = e
	} else {
		// it's a generic error, serialise it
		errReturn = &evm.SerialisableError{Err: err.Error()}
	}

	// serialise the error object returned by the evm into a json
	errSerializedBytes, marshallErr := json.Marshal(errReturn)
	if marshallErr != nil {
		return nil, marshallErr
	}
	return errSerializedBytes, nil
}

// handlePlainTextError ensures that SystemError errors are properly handled for PlaintextErrors
func (e *enclaveImpl) handlePlainTextError(err error) responses.EnclaveResponse {
	if errors.Is(err, errutil.SystemError{}) {
		e.logger.Error("enclave system err - ", log.ErrKey, err)
		return responses.AsSystemErr()
	}

	return responses.AsPlaintextError(err)
}

// handleEncryptedError ensures that SystemError errors are properly handled for possible EncryptedErrors
func (e *enclaveImpl) handleEncryptedError(err error, encryptor responses.ViewingKeyEncryptor) responses.EnclaveResponse {
	if errors.Is(err, errutil.SystemError{}) {
		e.logger.Error("enclave system err - ", log.ErrKey, err)
		return responses.AsSystemErr()
	}

	return responses.AsEncryptedError(err, encryptor)
}
