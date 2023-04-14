package components

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
)

type batchProducer struct {
	storage              db.Storage
	crossChainProcessors *crosschain.Processors
	genesis              *genesis.Genesis
	logger               gethlog.Logger
}

func NewBatchProducer(storage db.Storage, cc *crosschain.Processors, genesis *genesis.Genesis, logger gethlog.Logger) BatchProducer {
	return &batchProducer{
		storage:              storage,
		crossChainProcessors: cc,
		genesis:              genesis,
		logger:               logger,
	}
}

func (bp *batchProducer) populateOutboundCrossChainData(batch *core.Batch, block *types.Block, receipts types.Receipts) error {
	crossChainMessages, err := bp.crossChainProcessors.Local.ExtractOutboundMessages(receipts)
	if err != nil {
		bp.logger.Error("Extracting messages L2->L1 failed", err, log.CmpKey, log.CrossChainCmp)
		return fmt.Errorf("could not extract cross chain messages. Cause: %w", err)
	}

	batch.Header.CrossChainMessages = crossChainMessages

	bp.logger.Trace(fmt.Sprintf("Added %d cross chain messages to batch.",
		len(batch.Header.CrossChainMessages)), log.CmpKey, log.CrossChainCmp)

	batch.Header.LatestInboundCrossChainHash = block.Hash()
	batch.Header.LatestInboundCrossChainHeight = block.Number()

	return nil
}

func (bp *batchProducer) populateHeader(batch *core.Batch, receipts types.Receipts) {
	if len(receipts) == 0 {
		batch.Header.ReceiptHash = types.EmptyRootHash
	} else {
		batch.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		batch.Header.Bloom = types.CreateBloom(receipts)
	}

	if len(batch.Transactions) == 0 {
		batch.Header.TxHash = types.EmptyRootHash
	} else {
		batch.Header.TxHash = types.DeriveSha(types.Transactions(batch.Transactions), trie.NewStackTrie(nil))
	}
}

func (bp *batchProducer) ComputeBatch(context *BatchContext) (*ComputedBatch, error) {
	// These variables will be used to create the new batch

	parent, err := bp.storage.FetchBatch(context.ParentPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parent batch. Cause: %w", err)
	}

	block, err := bp.storage.FetchBlock(context.BlockPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block for batch. Cause: %w", err)
	}

	var parentBlock *types.Block = block
	if parent.Header.L1Proof != block.Hash() {
		var err error
		parentBlock, err = bp.storage.FetchBlock(parent.Header.L1Proof)
		if err != nil {
			bp.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", parent.Hash()), log.ErrKey, err)
		}
	}

	// Create a new batch based on the fromBlock of inclusion of the previous, including all new transactions
	batch := core.DeterministicEmptyBatch(context.Creator, parent.Header, block, context.Randomness, context.AtTime)

	stateDB, err := bp.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	messages := bp.crossChainProcessors.Local.RetrieveInboundMessages(parentBlock, block, stateDB)
	crossChainTransactions := bp.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)

	successfulTxs, txReceipts, err := bp.processTransactions(batch, 0, context.Transactions, stateDB, context.ChainConfig)
	if err != nil {
		return nil, err
	}

	ccSuccessfulTxs, ccReceipts, err := bp.processTransactions(batch, len(successfulTxs), crossChainTransactions, stateDB, context.ChainConfig)
	if err != nil {
		return nil, err
	}

	if err := bp.verifyInboundCrossChainTransactions(crossChainTransactions, ccSuccessfulTxs, ccReceipts); err != nil {
		return nil, fmt.Errorf("batch computation failed due to cross chain messages. Cause: %w", err)
	}

	batch.Header.Root = stateDB.IntermediateRoot(false)
	batch.Transactions = successfulTxs

	if err := bp.populateOutboundCrossChainData(batch, block, txReceipts); err != nil {
		return nil, fmt.Errorf("failed adding cross chain data to batch. Cause: %w", err)
	}

	bp.populateHeader(batch, allReceipts(txReceipts, ccReceipts))

	return &ComputedBatch{
		Batch:    batch,
		Receipts: txReceipts,
		Commit:   stateDB.Commit,
	}, nil
}

func (bp *batchProducer) verifyInboundCrossChainTransactions(transactions types.Transactions, executedTxs types.Transactions, receipts types.Receipts) error {
	if transactions.Len() != executedTxs.Len() {
		return fmt.Errorf("some synthetic transactions have not been executed")
	}

	for _, rec := range receipts {
		if rec.Status == 1 {
			continue
		}
		return fmt.Errorf("found a failed receipt for a synthetic transaction")
	}
	return nil
}

func (bp *batchProducer) processTransactions(batch *core.Batch, tCount int, txs []*common.L2Tx, stateDB *state.StateDB, cc *params.ChainConfig) ([]*common.L2Tx, []*types.Receipt, error) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, batch.Header, bp.storage, cc, 0, bp.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			return nil, nil, fmt.Errorf("there should be an entry for each transaction")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			bp.logger.Info(fmt.Sprintf("Excluding transaction %s from batch b_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(*batch.Hash()), result))
		}
	}
	sort.Sort(sortByTxIndex(txReceipts))

	return executedTransactions, txReceipts, nil
}

func (bp *batchProducer) CreateGenesisState(blkHash common.L1BlockHash, aggregatorAddress common.L2Address, timeNow uint64) (*core.Batch, *types.Transaction, error) {
	preFundGenesisState, err := bp.genesis.GetGenesisRoot(bp.storage)
	if err != nil {
		return nil, nil, err
	}

	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			Agg:         aggregatorAddress,
			ParentHash:  common.L2BatchHash{},
			L1Proof:     blkHash,
			Root:        *preFundGenesisState,
			TxHash:      types.EmptyRootHash,
			Number:      big.NewInt(int64(0)),
			ReceiptHash: types.EmptyRootHash,
			Time:        timeNow,
		},
		Transactions: []*common.L2Tx{},
	}

	// todo (#1577) - figure out a better way to bootstrap the system contracts
	deployTx, err := bp.crossChainProcessors.Local.GenerateMessageBusDeployTx()
	if err != nil {
		bp.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	if err = bp.genesis.CommitGenesisState(bp.storage); err != nil {
		return nil, nil, fmt.Errorf("could not apply genesis preallocation. Cause: %w", err)
	}
	return genesisBatch, deployTx, nil
}
