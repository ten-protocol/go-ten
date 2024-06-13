package components

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/holiman/uint256"

	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/common/gethencoding"

	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/core/state"
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
	config               config.EnclaveConfig
	gethEncodingService  gethencoding.EncodingService
	crossChainProcessors *crosschain.Processors
	genesis              *genesis.Genesis
	logger               gethlog.Logger
	gasOracle            gas.Oracle
	chainConfig          *params.ChainConfig

	// stateDBMutex - used to protect calls to stateDB.Commit as it is not safe for async access.
	stateDBMutex sync.Mutex

	batchGasLimit uint64 // max execution gas allowed in a batch
}

func NewBatchExecutor(
	storage storage.Storage,
	batchRegistry BatchRegistry,
	config config.EnclaveConfig,
	gethEncodingService gethencoding.EncodingService,
	cc *crosschain.Processors,
	genesis *genesis.Genesis,
	gasOracle gas.Oracle,
	chainConfig *params.ChainConfig,
	batchGasLimit uint64,
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
	}
}

// filterTransactionsWithSufficientFunds - this function estimates hte l1 fees for the transaction in a given batch execution context. It does so by taking the price of the
// pinned L1 block and using it as the cost per gas for the estimated gas of the calldata encoding of a transaction. It filters out any transactions that cannot afford to pay for their L1
// publishing cost.
func (executor *batchExecutor) filterTransactionsWithSufficientFunds(ctx context.Context, stateDB *state.StateDB, context *BatchExecutionContext) (common.L2PricedTransactions, common.L2PricedTransactions) {
	transactions := make(common.L2PricedTransactions, 0)
	freeTransactions := make(common.L2PricedTransactions, 0)
	block, _ := executor.storage.FetchBlock(ctx, context.BlockPtr)

	for _, tx := range context.Transactions {
		sender, err := core.GetAuthenticatedSender(context.ChainConfig.ChainID.Int64(), tx)
		if err != nil {
			executor.logger.Error("Unable to extract sender for tx. Should not happen at this point.", log.TxKey, tx.Hash(), log.ErrKey, err)
			continue
		}
		accBalance := stateDB.GetBalance(*sender)

		cost, err := executor.gasOracle.EstimateL1StorageGasCost(tx, block)
		if err != nil {
			executor.logger.Error("Unable to get gas cost for tx. Should not happen at this point.", log.TxKey, tx.Hash(), log.ErrKey, err)
			continue
		}

		// Transactions that are created inside the enclave can have no GasPrice set.
		// External transactions are always required to have a gas price set. Thus we filter
		// those transactions for separate processing than the normal ones and we run them through the EVM
		// with a flag that disables the baseFee logic and wont fail them for having price lower than the base fee.
		isFreeTransaction := tx.GasFeeCap().Cmp(gethcommon.Big0) == 0
		isFreeTransaction = isFreeTransaction && tx.GasPrice().Cmp(gethcommon.Big0) == 0

		if isFreeTransaction {
			freeTransactions = append(freeTransactions, common.L2PricedTransaction{
				Tx:             tx,
				PublishingCost: big.NewInt(0),
			})
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
	return transactions, freeTransactions
}

func (executor *batchExecutor) ComputeBatch(ctx context.Context, context *BatchExecutionContext, failForEmptyBatch bool) (*ComputedBatch, error) { //nolint:gocognit
	defer core.LogMethodDuration(executor.logger, measure.NewStopwatch(), "Batch context processed")

	// sanity check that the l1 block exists. We don't have to execute batches of forks.
	block, err := executor.storage.FetchBlock(ctx, context.BlockPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		return nil, errutil.ErrBlockForBatchNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %s for batch. Cause: %w", context.BlockPtr, err)
	}

	// These variables will be used to create the new batch
	parentBatch, err := executor.storage.FetchBatchHeader(ctx, context.ParentPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		executor.logger.Error(fmt.Sprintf("can't find parent batch %s. Seq %d", context.ParentPtr, context.SequencerNo))
		return nil, errutil.ErrAncestorBatchNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parent batch %s. Cause: %w", context.ParentPtr, err)
	}

	parentBlock := block
	if parentBatch.L1Proof != block.Hash() {
		var err error
		parentBlock, err = executor.storage.FetchBlock(ctx, parentBatch.L1Proof)
		if err != nil {
			executor.logger.Error(fmt.Sprintf("Could not retrieve a proof for batch %s", parentBatch.Hash()), log.ErrKey, err)
			return nil, err
		}
	}

	// Create a new batch based on the fromBlock of inclusion of the previous, including all new transactions
	batch := core.DeterministicEmptyBatch(parentBatch, block, context.AtTime, context.SequencerNo, context.BaseFee, context.Creator)

	stateDB, err := executor.batchRegistry.GetBatchState(ctx, &batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}
	snap := stateDB.Snapshot()

	var messages common.CrossChainMessages
	var transfers common.ValueTransferEvents
	if context.SequencerNo.Int64() > int64(common.L2GenesisSeqNo+1) {
		messages, transfers = executor.crossChainProcessors.Local.RetrieveInboundMessages(ctx, parentBlock, block, stateDB)
	}

	crossChainTransactions := executor.crossChainProcessors.Local.CreateSyntheticTransactions(ctx, messages, stateDB)
	executor.crossChainProcessors.Local.ExecuteValueTransfers(ctx, transfers, stateDB)

	transactionsToProcess, freeTransactions := executor.filterTransactionsWithSufficientFunds(ctx, stateDB, context)

	xchainTxs := make(common.L2PricedTransactions, 0)
	for _, xTx := range crossChainTransactions {
		xchainTxs = append(xchainTxs, common.L2PricedTransaction{
			Tx:             xTx,
			PublishingCost: big.NewInt(0),
		})
	}

	syntheticTransactions := append(xchainTxs, freeTransactions...)

	// fromTxIndex - Here we start from the 0 index. This will be the same for a validator.
	successfulTxs, excludedTxs, txReceipts, err := executor.processTransactions(ctx, batch, 0, transactionsToProcess, stateDB, context.ChainConfig, false)
	if err != nil {
		return nil, fmt.Errorf("could not process transactions. Cause: %w", err)
	}

	// fromTxIndex - Here we start from the len of the successful transactions; As long as we have the exact same successful transactions in a batch,
	// we will start from the same place.
	ccSuccessfulTxs, _, ccReceipts, err := executor.processTransactions(ctx, batch, len(successfulTxs), syntheticTransactions, stateDB, context.ChainConfig, true)
	if err != nil {
		return nil, err
	}

	if err = executor.verifyInboundCrossChainTransactions(syntheticTransactions, ccSuccessfulTxs, ccReceipts); err != nil {
		return nil, fmt.Errorf("batch computation failed due to cross chain messages. Cause: %w", err)
	}

	if failForEmptyBatch &&
		len(txReceipts) == 0 &&
		len(ccReceipts) == 0 &&
		len(transactionsToProcess)-len(excludedTxs) == 0 &&
		len(crossChainTransactions) == 0 &&
		len(messages) == 0 &&
		len(transfers) == 0 {
		if snap > 0 {
			//// revert any unexpected mutation to the statedb
			stateDB.RevertToSnapshot(snap)
		}
		return nil, ErrNoTransactionsToProcess
	}

	// we need to copy the batch to reset the internal hash cache
	copyBatch := *batch
	copyBatch.Header.Root = stateDB.IntermediateRoot(false)
	copyBatch.Transactions = append(successfulTxs, freeTransactions.ToTransactions()...)
	copyBatch.ResetHash()

	if err = executor.populateOutboundCrossChainData(ctx, &copyBatch, block, txReceipts); err != nil {
		return nil, fmt.Errorf("failed adding cross chain data to batch. Cause: %w", err)
	}

	allReceipts := append(txReceipts, ccReceipts...)
	executor.populateHeader(&copyBatch, allReceipts)

	// the logs and receipts produced by the EVM have the wrong hash which must be adjusted
	for _, receipt := range allReceipts {
		receipt.BlockHash = copyBatch.Hash()
		for _, l := range receipt.Logs {
			l.BlockHash = copyBatch.Hash()
		}
	}

	return &ComputedBatch{
		Batch:    &copyBatch,
		Receipts: allReceipts,
		Commit: func(deleteEmptyObjects bool) (gethcommon.Hash, error) {
			executor.stateDBMutex.Lock()
			defer executor.stateDBMutex.Unlock()
			h, err := stateDB.Commit(copyBatch.Number().Uint64(), deleteEmptyObjects)
			if err != nil {
				return gethcommon.Hash{}, fmt.Errorf("commit failure for batch %d. Cause: %w", batch.SeqNo(), err)
			}
			trieDB := executor.storage.TrieDB()
			err = trieDB.Commit(h, false)
			return h, err
		},
	}, nil
}

func (executor *batchExecutor) ExecuteBatch(ctx context.Context, batch *core.Batch) (types.Receipts, error) {
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

	return cb.Receipts, nil
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

	// todo (#1577) - figure out a better way to bootstrap the system contracts
	deployTx, err := executor.crossChainProcessors.Local.GenerateMessageBusDeployTx()
	if err != nil {
		executor.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	executor.logger.Info("L2 Bus deploy", log.TxKey, deployTx.Hash())

	if err = executor.genesis.CommitGenesisState(executor.storage); err != nil {
		return nil, nil, fmt.Errorf("could not apply genesis preallocation. Cause: %w", err)
	}
	return genesisBatch, deployTx, nil
}

func (executor *batchExecutor) populateOutboundCrossChainData(ctx context.Context, batch *core.Batch, block *types.Block, receipts types.Receipts) error {
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

	var xchainHash gethcommon.Hash = gethcommon.BigToHash(gethcommon.Big0)
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
		executor.logger.Info("[CrossChain] adding messages to batch")
	}
	batch.Header.CrossChainMessages = crossChainMessages
	batch.Header.CrossChainRoot = xchainHash

	executor.logger.Trace(fmt.Sprintf("Added %d cross chain messages to batch.",
		len(batch.Header.CrossChainMessages)), log.CmpKey, log.CrossChainCmp)

	batch.Header.LatestInboundCrossChainHash = block.Hash()
	batch.Header.LatestInboundCrossChainHeight = block.Number()

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

func (executor *batchExecutor) verifyInboundCrossChainTransactions(transactions common.L2PricedTransactions, executedTxs types.Transactions, receipts types.Receipts) error {
	if len(transactions) != executedTxs.Len() {
		return fmt.Errorf("some synthetic transactions have not been executed")
	}

	for _, rec := range receipts {
		if rec.Status == 1 {
			continue
		}
		return fmt.Errorf("found a failed receipt for a synthetic transaction: %s", rec.TxHash.Hex())
	}
	return nil
}

func (executor *batchExecutor) processTransactions(
	ctx context.Context,
	batch *core.Batch,
	tCount int,
	txs common.L2PricedTransactions,
	stateDB *state.StateDB,
	cc *params.ChainConfig,
	noBaseFee bool,
) ([]*common.L2Tx, []*common.L2Tx, []*types.Receipt, error) {
	var executedTransactions []*common.L2Tx
	var excludedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt
	txResults := evm.ExecuteTransactions(
		ctx,
		txs,
		stateDB,
		batch.Header,
		executor.storage,
		executor.gethEncodingService,
		cc,
		executor.config,
		tCount,
		noBaseFee,
		executor.batchGasLimit,
		executor.logger,
	)
	for _, tx := range txs {
		result, f := txResults[tx.Tx.Hash()]
		if !f {
			return nil, nil, nil, fmt.Errorf("there should be an entry for each transaction")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx.Tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			excludedTransactions = append(excludedTransactions, tx.Tx)
			executor.logger.Debug("Excluding transaction from batch", log.TxKey, tx.Tx.Hash(), log.BatchHashKey, batch.Hash(), "cause", result)
		}
	}
	sort.Sort(sortByTxIndex(txReceipts))

	return executedTransactions, excludedTransactions, txReceipts, nil
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }
