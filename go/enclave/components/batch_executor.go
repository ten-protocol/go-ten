package components

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/holiman/uint256"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/system"

	gethcommon "github.com/ethereum/go-ethereum/common"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/evm"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
)

var ErrNoTransactionsToProcess = fmt.Errorf("no transactions to process")

// batchExecutor - the component responsible for executing batches
type batchExecutor struct {
	storage              storage.Storage
	batchRegistry        BatchRegistry
	config               enclaveconfig.EnclaveConfig
	gethEncodingService  gethencoding.EncodingService
	crossChainProcessors *crosschain.Processors
	genesis              *genesis.Genesis
	logger               gethlog.Logger
	gasOracle            gas.Oracle
	chainConfig          *params.ChainConfig
	systemContracts      system.SystemContractCallbacks
	entropyService       *crypto.EvmEntropyService
	// stateDBMutex - used to protect calls to stateDB.Commit as it is not safe for async access.
	stateDBMutex sync.Mutex

	batchGasLimit uint64 // max execution gas allowed in a batch
}

func NewBatchExecutor(
	storage storage.Storage,
	batchRegistry BatchRegistry,
	config enclaveconfig.EnclaveConfig,
	gethEncodingService gethencoding.EncodingService,
	cc *crosschain.Processors,
	genesis *genesis.Genesis,
	gasOracle gas.Oracle,
	chainConfig *params.ChainConfig,
	batchGasLimit uint64,
	systemContracts system.SystemContractCallbacks,
	entropyService *crypto.EvmEntropyService,
	logger gethlog.Logger,
) BatchExecutor {
	return &batchExecutor{
		storage:              storage,
		batchRegistry:        batchRegistry,
		config:               config,
		gethEncodingService:  gethEncodingService,
		crossChainProcessors: cc,
		genesis:              genesis,
		chainConfig:          chainConfig,
		logger:               logger,
		gasOracle:            gasOracle,
		stateDBMutex:         sync.Mutex{},
		batchGasLimit:        batchGasLimit,
		systemContracts:      systemContracts,
		entropyService:       entropyService,
	}
}

func (executor *batchExecutor) ComputeBatch(ctx context.Context, ec *BatchExecutionContext, failForEmptyBatch bool) (*ComputedBatch, error) {
	defer core.LogMethodDuration(executor.logger, measure.NewStopwatch(), "Batch context processed")

	ec.ctx = ctx
	if err := executor.verifyContext(ec); err != nil {
		return nil, err
	}

	if err := executor.prepareState(ec); err != nil {
		return nil, err
	}

	// the batch with seqNo==2 is by convention the batch where we deploy the system contracts
	if ec.SequencerNo.Uint64() == common.L2SysContractGenesisSeqNo {
		if err := executor.handleSysContractGenesis(ec); err != nil {
			return nil, err
		}
		// the sys genesis batch will not contain anything else
		return executor.execResult(ec)
	}

	// for every batch, the first thing, we execute the xChain messages
	if err := executor.readXChainMessages(ec); err != nil {
		return nil, err
	}
	if err := executor.execXChainMessages(ec); err != nil {
		return nil, err
	}

	// we execute the transactions included in the batch
	if err := executor.execBatchTransactions(ec); err != nil {
		return nil, err
	}

	// execute the callbacks
	if err := executor.execRegisteredCallbacks(ec); err != nil {
		return nil, err
	}

	// execute the system contract registered at the end of the block
	if err := executor.execOnBlockEndTx(ec); err != nil {
		return nil, err
	}

	// When the `failForEmptyBatch` flag is true, we skip if there is no transaction or xChain tx
	if failForEmptyBatch && len(ec.batchTxResults) == 0 && len(ec.xChainResults) == 0 {
		if ec.beforeProcessingSnap > 0 {
			//// revert any unexpected mutation to the statedb
			ec.stateDB.RevertToSnapshot(ec.beforeProcessingSnap)
		}
		return nil, ErrNoTransactionsToProcess
	}

	return executor.execResult(ec)
}

func (executor *batchExecutor) verifyContext(ec *BatchExecutionContext) error {
	// sanity check that the l1 block exists. We don't have to execute batches of forks.
	block, err := executor.storage.FetchBlock(ec.ctx, ec.BlockPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		return errutil.ErrBlockForBatchNotFound
	} else if err != nil {
		return fmt.Errorf("failed to retrieve block %s for batch. Cause: %w", ec.BlockPtr, err)
	}

	ec.l1block = block

	// These variables will be used to create the new batch
	parentBatch, err := executor.storage.FetchBatchHeader(ec.ctx, ec.ParentPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		executor.logger.Error(fmt.Sprintf("can't find parent batch %s. Seq %d", ec.ParentPtr, ec.SequencerNo))
		return errutil.ErrAncestorBatchNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve parent batch %s. Cause: %w", ec.ParentPtr, err)
	}
	ec.parentBatch = parentBatch

	parentBlock := block
	if parentBatch.L1Proof != block.Hash() {
		var err error
		parentBlock, err = executor.storage.FetchBlock(ec.ctx, parentBatch.L1Proof)
		if err != nil {
			executor.logger.Error(fmt.Sprintf("Could not retrieve a proof for batch %s", parentBatch.Hash()), log.ErrKey, err)
			return err
		}
	}
	ec.parentL1Block = parentBlock

	return nil
}

func (executor *batchExecutor) prepareState(ec *BatchExecutionContext) error {
	var err error
	// Create a new batch based on the provided context
	ec.currentBatch = core.DeterministicEmptyBatch(ec.parentBatch, ec.l1block, ec.AtTime, ec.SequencerNo, ec.BaseFee, ec.Creator)
	ec.stateDB, err = executor.batchRegistry.GetBatchState(ec.ctx, rpc.BlockNumberOrHash{BlockHash: &ec.currentBatch.Header.ParentHash})
	if err != nil {
		return fmt.Errorf("could not create stateDB. Cause: %w", err)
	}
	ec.beforeProcessingSnap = ec.stateDB.Snapshot()
	return nil
}

func (executor *batchExecutor) handleSysContractGenesis(ec *BatchExecutionContext) error {
	systemDeployerTx, err := system.SystemDeployerInitTransaction(executor.logger, *executor.systemContracts.SystemContractsUpgrader())
	if err != nil {
		executor.logger.Error("[SystemContracts] Failed to create system deployer contract", log.ErrKey, err)
		return err
	}

	transactions := common.L2PricedTransactions{
		common.L2PricedTransaction{
			Tx:             systemDeployerTx,
			PublishingCost: big.NewInt(0),
			SystemDeployer: true,
		},
	}

	sysCtrGenesisResult, err := executor.executeTxs(ec, 0, transactions, true)
	if err != nil {
		return fmt.Errorf("could not process system deployer transaction. Cause: %w", err)
	}

	if err = executor.verifySyntheticTransactionsSuccess(transactions, sysCtrGenesisResult); err != nil {
		return fmt.Errorf("batch computation failed due to system deployer reverting. Cause: %w", err)
	}

	ec.genesisSysCtrResult = sysCtrGenesisResult
	ec.genesisSysCtrResult.MarkSynthetic(true)
	return nil
}

func (executor *batchExecutor) readXChainMessages(ec *BatchExecutionContext) error {
	if ec.SequencerNo.Int64() > int64(common.L2SysContractGenesisSeqNo) {
		ec.xChainMsgs, ec.xChainValueMsgs = executor.crossChainProcessors.Local.RetrieveInboundMessages(ec.ctx, ec.parentL1Block, ec.l1block)
	}
	return nil
}

func (executor *batchExecutor) execXChainMessages(ec *BatchExecutionContext) error {
	crossChainTransactions := executor.crossChainProcessors.Local.CreateSyntheticTransactions(ec.ctx, ec.xChainMsgs, ec.stateDB)
	executor.crossChainProcessors.Local.ExecuteValueTransfers(ec.ctx, ec.xChainValueMsgs, ec.stateDB)
	xchainTxs := make(common.L2PricedTransactions, 0)
	for _, xTx := range crossChainTransactions {
		xchainTxs = append(xchainTxs, common.L2PricedTransaction{
			Tx:             xTx,
			PublishingCost: big.NewInt(0),
			FromSelf:       true,
		})
	}
	xChainResults, err := executor.executeTxs(ec, 0, xchainTxs, true)
	if err != nil {
		return fmt.Errorf("could not process cross chain messages. Cause: %w", err)
	}

	if len(xchainTxs) != len(xChainResults) {
		return fmt.Errorf("could not process cross chain messages. Some were excluded. Cause: %w", err)
	}
	ec.xChainResults = xChainResults
	ec.xChainResults.MarkSynthetic(true)
	return nil
}

// filterTransactionsWithSufficientFunds - this function estimates hte l1 fees for the transaction in a given batch execution context. It does so by taking the price of the
// pinned L1 block and using it as the cost per gas for the estimated gas of the calldata encoding of a transaction. It filters out any transactions that cannot afford to pay for their L1
// publishing cost.
func (executor *batchExecutor) filterTransactionsWithSufficientFunds(ec *BatchExecutionContext) common.L2PricedTransactions {
	transactions := make(common.L2PricedTransactions, 0)
	block, _ := executor.storage.FetchBlock(ec.ctx, ec.BlockPtr)

	for _, tx := range ec.Transactions {
		sender, err := core.GetAuthenticatedSender(ec.ChainConfig.ChainID.Int64(), tx)
		if err != nil {
			executor.logger.Error("Unable to extract sender for tx. Should not happen at this point.", log.TxKey, tx.Hash(), log.ErrKey, err)
			continue
		}
		accBalance := ec.stateDB.GetBalance(*sender)

		cost, err := executor.gasOracle.EstimateL1StorageGasCost(tx, block)
		if err != nil {
			executor.logger.Error("Unable to get gas cost for tx. Should not happen at this point.", log.TxKey, tx.Hash(), log.ErrKey, err)
			continue
		}

		if accBalance.Cmp(uint256.MustFromBig(cost)) == -1 {
			executor.logger.Info(fmt.Sprintf("insufficient account balance for tx - want: %d have: %d", cost, accBalance), log.TxKey, tx.Hash(), "addr", sender.Hex())
			continue
		}

		transactions = append(transactions, common.L2PricedTransaction{
			Tx:             tx,
			PublishingCost: big.NewInt(0).Set(cost),
		})
	}
	return transactions
}

func (executor *batchExecutor) execBatchTransactions(ec *BatchExecutionContext) error {
	transactionsToProcess := executor.filterTransactionsWithSufficientFunds(ec)

	txResults, err := executor.executeTxs(ec, len(ec.xChainResults), transactionsToProcess, false)
	if err != nil {
		return fmt.Errorf("could not process transactions. Cause: %w", err)
	}

	ec.batchTxResults = txResults
	return nil
}

func (executor *batchExecutor) execRegisteredCallbacks(ec *BatchExecutionContext) error {
	// Create and process public callback transaction if needed
	publicCallbackTx, err := executor.systemContracts.CreatePublicCallbackHandlerTransaction(ec.ctx, ec.stateDB)
	if err != nil {
		return fmt.Errorf("could not create public callback transaction. Cause: %w", err)
	}

	if publicCallbackTx == nil {
		return nil
	}

	publicCallbackPricedTxes := common.L2PricedTransactions{
		common.L2PricedTransaction{
			Tx:             publicCallbackTx,
			PublishingCost: big.NewInt(0),
			FromSelf:       true,
		},
	}
	offset := len(ec.batchTxResults) + len(ec.xChainResults)
	publicCallbackTxResult, err := executor.executeTxs(ec, offset, publicCallbackPricedTxes, true)
	if err != nil {
		return fmt.Errorf("could not process public callback transaction. Cause: %w", err)
	}
	// Ensure the public callback transaction is successful. It should NEVER fail.
	if err = executor.verifySyntheticTransactionsSuccess(publicCallbackPricedTxes, publicCallbackTxResult); err != nil {
		return fmt.Errorf("batch computation failed due to public callback reverting. Cause: %w", err)
	}
	ec.callbackTxResults = publicCallbackTxResult
	ec.callbackTxResults.MarkSynthetic(true)
	return nil
}

func (executor *batchExecutor) execOnBlockEndTx(ec *BatchExecutionContext) error {
	onBlockTx, err := executor.systemContracts.CreateOnBatchEndTransaction(ec.ctx, ec.stateDB, ec.batchTxResults)
	if err != nil && !errors.Is(err, system.ErrNoTransactions) {
		return fmt.Errorf("could not create on block end transaction. Cause: %w", err)
	}
	if onBlockTx == nil {
		return nil
	}
	onBlockPricedTx := common.L2PricedTransactions{
		common.L2PricedTransaction{
			Tx:             onBlockTx,
			PublishingCost: big.NewInt(0),
			FromSelf:       true,
		},
	}
	offset := len(ec.callbackTxResults) + len(ec.batchTxResults) + len(ec.xChainResults)
	onBlockTxResult, err := executor.executeTxs(ec, offset, onBlockPricedTx, true)
	if err != nil {
		return fmt.Errorf("could not process on block end transaction hook. Cause: %w", err)
	}
	// Ensure the onBlock callback transaction is successful. It should NEVER fail.
	if err = executor.verifySyntheticTransactionsSuccess(onBlockPricedTx, onBlockTxResult); err != nil {
		return fmt.Errorf("batch computation failed due to onBlock hook reverting. Cause: %w", err)
	}
	ec.blockEndResult = onBlockTxResult
	ec.blockEndResult.MarkSynthetic(true)
	// todo - stefan not clear what this does
	//result := onBlockTxResult[0]
	//if ok, err := executor.systemContracts.VerifyOnBlockReceipt(successfulTxs, result.Receipt); !ok || err != nil {
	//	executor.logger.Error("VerifyOnBlockReceipt failed", "error", err, "ok", ok)
	//	return fmt.Errorf("VerifyOnBlockReceipt failed")
	//}
	return nil
}

func (executor *batchExecutor) execResult(ec *BatchExecutionContext) (*ComputedBatch, error) {
	// we need to copy the batch to reset the internal hash cache
	batch := *ec.currentBatch
	batch.Header.Root = ec.stateDB.IntermediateRoot(false)
	batch.Transactions = ec.batchTxResults.Transactions()
	batch.ResetHash()

	txReceipts := ec.batchTxResults.Receipts()
	if err := executor.populateOutboundCrossChainData(ec.ctx, &batch, ec.l1block, txReceipts); err != nil {
		return nil, fmt.Errorf("failed adding cross chain data to batch. Cause: %w", err)
	}

	xChainReceipts := ec.xChainResults.Receipts()
	allReceipts := append(txReceipts, xChainReceipts...)
	executor.populateHeader(&batch, allReceipts)

	// the logs and receipts produced by the EVM have the wrong hash which must be adjusted
	//for _, receipt := range allReceipts {
	//	receipt.BlockHash = batch.Hash()
	//	for _, l := range receipt.Logs {
	//		l.BlockHash = batch.Hash()
	//	}
	//}

	commitFunc := func(deleteEmptyObjects bool) (gethcommon.Hash, error) {
		executor.stateDBMutex.Lock()
		defer executor.stateDBMutex.Unlock()
		h, err := ec.stateDB.Commit(batch.Number().Uint64(), deleteEmptyObjects)
		if err != nil {
			return gethcommon.Hash{}, fmt.Errorf("commit failure for batch %d. Cause: %w", ec.currentBatch.SeqNo(), err)
		}
		trieDB := executor.storage.TrieDB()
		err = trieDB.Commit(h, false)

		// When system contract deployment genesis batch is committed, initialize executor's addresses for the hooks.
		// Further restarts will call into Load() which will take the receipts for batch number 2 (which should never be deleted)
		// and reinitialize them.
		if err == nil && ec.currentBatch.Header.SequencerOrderNo.Uint64() == common.L2SysContractGenesisSeqNo {
			if len(ec.genesisSysCtrResult) == 0 {
				return h, fmt.Errorf("failed to instantiate system contracts: expected receipt for system deployer transaction, but no receipts found in batch")
			}

			return h, executor.systemContracts.Initialize(&batch, *ec.genesisSysCtrResult.Receipts()[0], executor.crossChainProcessors.Local)
		}

		return h, err
	}

	allResults := append(append(append(append(ec.xChainResults, ec.batchTxResults...), ec.callbackTxResults...), ec.blockEndResult...), ec.genesisSysCtrResult...)
	return &ComputedBatch{
		Batch:         &batch,
		TxExecResults: allResults,
		Commit:        commitFunc,
	}, nil
}

func (executor *batchExecutor) ExecuteBatch(ctx context.Context, batch *core.Batch) ([]*core.TxExecResult, error) {
	defer core.LogMethodDuration(executor.logger, measure.NewStopwatch(), "Executed batch", log.BatchHashKey, batch.Hash())

	// Validators recompute the entire batch using the same batch context
	// if they have all necessary prerequisites like having the l1 block processed
	// and the parent hash. This recomputed batch is then checked against the incoming batch.
	// If the sequencer has tampered with something the hash will not add up and validation will
	// produce an error.
	cb, err := executor.ComputeBatch(ctx, &BatchExecutionContext{
		BlockPtr:     batch.Header.L1Proof,
		ParentPtr:    batch.Header.ParentHash,
		Transactions: batch.Transactions,
		AtTime:       batch.Header.Time,
		ChainConfig:  executor.chainConfig,
		SequencerNo:  batch.Header.SequencerOrderNo,
		Creator:      batch.Header.Coinbase,
		BaseFee:      batch.Header.BaseFee,
	}, false) // this execution is not used when first producing a batch, we never want to fail for empty batches
	if err != nil {
		return nil, fmt.Errorf("failed computing batch %s. Cause: %w", batch.Hash(), err)
	}

	if cb.Batch.Hash() != batch.Hash() {
		// todo @stefan - generate a validator challenge here and return it
		executor.logger.Error(fmt.Sprintf("Error validating batch. Calculated: %+v    Incoming: %+v", cb.Batch.Header, batch.Header))
		return nil, fmt.Errorf("batch is in invalid state. Incoming hash: %s  Computed hash: %s", batch.Hash(), cb.Batch.Hash())
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("cannot commit stateDB for incoming valid batch %s. Cause: %w", batch.Hash(), err)
	}

	return cb.TxExecResults, nil
}

func (executor *batchExecutor) CreateGenesisState(
	ctx context.Context,
	blkHash common.L1BlockHash,
	timeNow uint64,
	coinbase gethcommon.Address,
	baseFee *big.Int,
) (*core.Batch, *types.Transaction, error) {
	preFundGenesisState, err := executor.genesis.GetGenesisRoot(executor.storage)
	if err != nil {
		return nil, nil, err
	}

	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			ParentHash:       common.L2BatchHash{},
			L1Proof:          blkHash,
			Root:             *preFundGenesisState,
			TxHash:           types.EmptyRootHash,
			Number:           big.NewInt(int64(0)),
			SequencerOrderNo: big.NewInt(int64(common.L2GenesisSeqNo)), // genesis batch has seq number 1
			ReceiptHash:      types.EmptyRootHash,
			CrossChainRoot:   types.EmptyRootHash,
			Time:             timeNow,
			Coinbase:         coinbase,
			BaseFee:          baseFee,
			GasLimit:         executor.batchGasLimit,
		},
		Transactions: []*common.L2Tx{},
	}

	if err = executor.genesis.CommitGenesisState(executor.storage); err != nil {
		return nil, nil, fmt.Errorf("could not apply genesis preallocation. Cause: %w", err)
	}
	return genesisBatch, nil, nil
}

func (executor *batchExecutor) populateOutboundCrossChainData(ctx context.Context, batch *core.Batch, block *types.Header, receipts types.Receipts) error {
	crossChainMessages, err := executor.crossChainProcessors.Local.ExtractOutboundMessages(ctx, receipts)
	if err != nil {
		executor.logger.Error("Failed extracting L2->L1 messages", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
		return fmt.Errorf("could not extract cross chain messages. Cause: %w", err)
	}

	valueTransferMessages, err := executor.crossChainProcessors.Local.ExtractOutboundTransfers(ctx, receipts)
	if err != nil {
		executor.logger.Error("Failed extracting L2->L1 messages value transfers", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
		return fmt.Errorf("could not extract cross chain value transfers. Cause: %w", err)
	}

	xchainTree := make([][]interface{}, 0)

	hasMessages := false
	if len(valueTransferMessages) > 0 {
		transfers := crosschain.ValueTransfers(valueTransferMessages).ForMerkleTree()
		xchainTree = append(xchainTree, transfers...)
		hasMessages = true
	}

	if len(crossChainMessages) > 0 {
		messages := crosschain.MessageStructs(crossChainMessages).ForMerkleTree()
		xchainTree = append(xchainTree, messages...)
		hasMessages = true
	}

	xchainHash := gethcommon.BigToHash(gethcommon.Big0)
	if hasMessages {
		tree, err := smt.Of(xchainTree, crosschain.CrossChainEncodings)
		if err != nil {
			executor.logger.Error("Unable to create merkle tree for cross chain messages", log.ErrKey, err)
			return fmt.Errorf("unable to create merkle tree for cross chain messages. Cause: %w", err)
		}

		encodedTree, err := json.Marshal(xchainTree)
		if err != nil {
			panic(err) // todo: figure out what to do
		}

		batch.Header.CrossChainTree = encodedTree
		xchainHash = gethcommon.BytesToHash(tree.GetRoot())
		executor.logger.Debug("[CrossChain] adding messages to batch", "encodedTree", encodedTree)
	}
	batch.Header.CrossChainMessages = crossChainMessages
	batch.Header.CrossChainRoot = xchainHash

	executor.logger.Debug(fmt.Sprintf("Added %d cross chain messages to batch.",
		len(batch.Header.CrossChainMessages)), log.CmpKey, log.CrossChainCmp)

	batch.Header.LatestInboundCrossChainHash = block.Hash()
	batch.Header.LatestInboundCrossChainHeight = block.Number

	return nil
}

func (executor *batchExecutor) populateHeader(batch *core.Batch, receipts types.Receipts) {
	if len(receipts) == 0 {
		batch.Header.ReceiptHash = types.EmptyRootHash
	} else {
		batch.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
	}

	if len(batch.Transactions) == 0 {
		batch.Header.TxHash = types.EmptyRootHash
	} else {
		batch.Header.TxHash = types.DeriveSha(types.Transactions(batch.Transactions), trie.NewStackTrie(nil))
	}
}

func (executor *batchExecutor) verifySyntheticTransactionsSuccess(transactions common.L2PricedTransactions, results core.TxExecResults) error {
	if len(transactions) != len(results) {
		return fmt.Errorf("some synthetic transactions have not been executed")
	}

	for _, rec := range results {
		if rec.Receipt.Status == 1 {
			continue
		}
		return fmt.Errorf("found a failed receipt for a synthetic transaction: %s", rec.Receipt.TxHash.Hex())
	}
	return nil
}

func (executor *batchExecutor) executeTxs(ec *BatchExecutionContext, offset int, txs common.L2PricedTransactions, noBaseFee bool) (core.TxExecResults, error) {
	txResultsMap, err := evm.ExecuteTransactions(
		ec.ctx,
		executor.entropyService,
		txs,
		ec.stateDB,
		ec.currentBatch.Header,
		executor.storage,
		executor.gethEncodingService,
		ec.ChainConfig,
		executor.config,
		offset,
		noBaseFee,
		executor.batchGasLimit,
		executor.logger,
	)
	if err != nil {
		return nil, err
	}
	txResults := make(core.TxExecResults, 0)
	for _, tx := range txs {
		result, f := txResultsMap[tx.Tx.Hash()]
		if !f {
			return nil, fmt.Errorf("there should be an entry for each transaction")
		}
		if result.Receipt != nil {
			txResults = append(txResults, result)
		} else {
			// Exclude failed transactions
			executor.logger.Debug("Excluding transaction from batch", log.TxKey, tx.Tx.Hash(), log.BatchHashKey, ec.currentBatch.Hash(), "cause", result.Err)
		}
	}

	// populate the derived fields in the receipt
	txReceipts := make(types.Receipts, 0)
	for _, txResult := range txResults {
		txReceipts = append(txReceipts, txResult.Receipt)
	}
	batch := ec.currentBatch
	err = txReceipts.DeriveFields(executor.chainConfig, batch.Hash(), batch.NumberU64(), batch.Header.Time, batch.Header.BaseFee, nil, txResults.Transactions())
	if err != nil {
		return nil, fmt.Errorf("could not derive receipts. Cause: %w", err)
	}
	for i, txResult := range txResults {
		txResult.Receipt = txReceipts[i]
	}

	sort.Sort(sortByTxIndex(txResults))
	return txResults, nil
}

type sortByTxIndex []*core.TxExecResult

func (c sortByTxIndex) Len() int      { return len(c) }
func (c sortByTxIndex) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool {
	return c[i].Receipt.TransactionIndex < c[j].Receipt.TransactionIndex
}
