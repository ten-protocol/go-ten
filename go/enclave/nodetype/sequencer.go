package nodetype

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/signature"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/common/compression"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"
)

const RollupDelay = 2 // number of L1 blocks to exclude when creating a rollup. This will minimize compression reorg issues.

type SequencerSettings struct {
	MaxBatchSize      uint64
	MaxRollupSize     uint64
	GasPaymentAddress gethcommon.Address
	BatchGasLimit     uint64
	BaseFee           *big.Int
}

type sequencer struct {
	blockProcessor    components.L1BlockProcessor
	batchProducer     components.BatchExecutor
	batchRegistry     components.BatchRegistry
	rollupProducer    components.RollupProducer
	rollupConsumer    components.RollupConsumer
	rollupCompression *components.RollupCompression
	gethEncoding      gethencoding.EncodingService

	logger gethlog.Logger

	chainConfig            *params.ChainConfig
	enclaveKey             *crypto.EnclaveKey
	mempool                *txpool.TxPool
	storage                storage.Storage
	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	settings               SequencerSettings
	blockchain             *ethchainadapter.EthChainAdapter
}

func NewSequencer(
	blockProcessor components.L1BlockProcessor,
	batchExecutor components.BatchExecutor,
	registry components.BatchRegistry,
	rollupProducer components.RollupProducer,
	rollupConsumer components.RollupConsumer,
	rollupCompression *components.RollupCompression,
	gethEncodingService gethencoding.EncodingService,
	logger gethlog.Logger,
	chainConfig *params.ChainConfig,
	enclavePrivateKey *crypto.EnclaveKey,
	mempool *txpool.TxPool,
	storage storage.Storage,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	settings SequencerSettings,
	blockchain *ethchainadapter.EthChainAdapter,
) Sequencer {
	return &sequencer{
		blockProcessor:         blockProcessor,
		batchProducer:          batchExecutor,
		batchRegistry:          registry,
		rollupProducer:         rollupProducer,
		rollupConsumer:         rollupConsumer,
		rollupCompression:      rollupCompression,
		gethEncoding:           gethEncodingService,
		logger:                 logger,
		chainConfig:            chainConfig,
		enclaveKey:             enclavePrivateKey,
		mempool:                mempool,
		storage:                storage,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		settings:               settings,
		blockchain:             blockchain,
	}
}

func (s *sequencer) CreateBatch(ctx context.Context, skipBatchIfEmpty bool) error {
	hasGenesis, err := s.batchRegistry.HasGenesisBatch()
	if err != nil {
		return fmt.Errorf("unknown genesis batch state. Cause: %w", err)
	}

	// L1 Head is only updated when isLatest: true
	l1HeadBlock, err := s.blockProcessor.GetHead(ctx)
	if err != nil {
		return fmt.Errorf("failed retrieving l1 head. Cause: %w", err)
	}

	// the sequencer creates the initial genesis batch if one does not exist yet
	if !hasGenesis {
		return s.createGenesisBatch(ctx, l1HeadBlock)
	}

	if running := s.mempool.Running(); !running {
		// the mempool can only be started after at least 1 block (the genesis) is in the blockchain object
		// if the node restarted the mempool must be started again
		err = s.mempool.Start()
		if err != nil {
			return err
		}
	}

	return s.createNewHeadBatch(ctx, l1HeadBlock, skipBatchIfEmpty)
}

// TODO - This is iffy, the producer commits the stateDB. The producer
// should only create batches and stateDBs but not commit them to the database,
// this is the responsibility of the sequencer. Refactor the code so genesis state
// won't be committed by the producer.
func (s *sequencer) createGenesisBatch(ctx context.Context, block *common.L1Block) error {
	s.logger.Info("Initializing genesis state", log.BlockHashKey, block.Hash())
	batch, msgBusTx, err := s.batchProducer.CreateGenesisState(
		ctx,
		block.Hash(),
		uint64(time.Now().Unix()),
		s.settings.GasPaymentAddress,
		s.settings.BaseFee,
	)
	if err != nil {
		return err
	}

	if err := s.signBatch(batch); err != nil {
		return fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.StoreExecutedBatch(ctx, batch, nil); err != nil {
		return fmt.Errorf("1. failed storing batch. Cause: %w", err)
	}

	// this is the actual first block produced in chain
	err = s.blockchain.IngestNewBlock(batch)
	if err != nil {
		return fmt.Errorf("failed to feed batch into the virtual eth chain - %w", err)
	}

	// the mempool can only be started after at least 1 block is in the blockchain object
	err = s.mempool.Start()
	if err != nil {
		return err
	}

	// errors in unit test seem to suggest that batch 2 was received before batch 1
	// this ensures that there is enough gap so that batch 1 is issued before batch 2
	time.Sleep(time.Second)
	// produce batch #2 which has the message bus and any other system contracts
	cb, err := s.produceBatch(
		ctx,
		big.NewInt(0).Add(batch.Header.SequencerOrderNo, big.NewInt(1)),
		block.Hash(),
		batch.Hash(),
		common.L2Transactions{msgBusTx},
		uint64(time.Now().Unix()),
		false,
	)
	if err != nil {
		if errors.Is(err, components.ErrNoTransactionsToProcess) {
			// skip batch production when there are no transactions to process
			// todo: this might be a useful event to track for metrics (skipping batch production because empty batch)
			s.logger.Debug("Skipping batch production, no transactions to execute")
			return nil
		}
		return fmt.Errorf(" failed producing batch. Cause: %w", err)
	}

	if len(cb.Receipts) == 0 || cb.Receipts[0].TxHash.Hex() != msgBusTx.Hash().Hex() {
		err = fmt.Errorf("message Bus contract not minted - no receipts in batch")
		s.logger.Error(err.Error())
		return err
	}

	s.logger.Info("Message Bus Contract minted successfully", "address", cb.Receipts[0].ContractAddress.Hex())

	return nil
}

func (s *sequencer) createNewHeadBatch(ctx context.Context, l1HeadBlock *common.L1Block, skipBatchIfEmpty bool) error {
	headBatchSeq := s.batchRegistry.HeadBatchSeq()
	if headBatchSeq == nil {
		headBatchSeq = big.NewInt(int64(common.L2GenesisSeqNo))
	}
	headBatch, err := s.storage.FetchBatchBySeqNo(ctx, headBatchSeq.Uint64())
	if err != nil {
		return err
	}

	// todo - sanity check that the headBatch.Header.L1Proof is an ancestor of the l1HeadBlock
	b, err := s.storage.FetchBlock(ctx, headBatch.Header.L1Proof)
	if err != nil {
		return err
	}
	if !s.storage.IsAncestor(ctx, l1HeadBlock, b) {
		return fmt.Errorf("attempted to create batch on top of batch=%s. With l1 head=%s", headBatch.Hash(), l1HeadBlock.Hash())
	}

	// todo (@stefan) - limit on receipts too
	limiter := limiters.NewBatchSizeLimiter(s.settings.MaxBatchSize)
	pendingTransactions := s.mempool.PendingTransactions()
	var transactions []*types.Transaction
	for _, group := range pendingTransactions {
		// lazily resolve transactions until the batch runs out of space
		for _, lazyTx := range group {
			if tx := lazyTx.Resolve(); tx != nil {
				err = limiter.AcceptTransaction(tx)
				if err != nil {
					s.logger.Info("Unable to accept transaction", log.TxKey, tx.Hash(), log.ErrKey, err)
					if errors.Is(err, limiters.ErrInsufficientSpace) { // Batch ran out of space
						break
					}
					// Limiter encountered unexpected error
					return fmt.Errorf("limiter encountered unexpected error - %w", err)
				}
				transactions = append(transactions, tx)
			}
		}
	}

	sequencerNo, err := s.storage.FetchCurrentSequencerNo(ctx)
	if err != nil {
		return err
	}

	// todo - time is set only here; take from l1 block?
	if _, err := s.produceBatch(ctx, sequencerNo.Add(sequencerNo, big.NewInt(1)), l1HeadBlock.Hash(), headBatch.Hash(), transactions, uint64(time.Now().Unix()), skipBatchIfEmpty); err != nil {
		if errors.Is(err, components.ErrNoTransactionsToProcess) {
			// skip batch production when there are no transactions to process
			// todo: this might be a useful event to track for metrics (skipping batch production because empty batch)
			s.logger.Debug("Skipping batch production, no transactions to execute")
			return nil
		}
		return fmt.Errorf(" failed producing batch. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) produceBatch(
	ctx context.Context,
	sequencerNo *big.Int,
	l1Hash common.L1BlockHash,
	headBatch common.L2BatchHash,
	transactions common.L2Transactions,
	batchTime uint64,
	failForEmptyBatch bool,
) (*components.ComputedBatch, error) {
	cb, err := s.batchProducer.ComputeBatch(ctx,
		&components.BatchExecutionContext{
			BlockPtr:     l1Hash,
			ParentPtr:    headBatch,
			Transactions: transactions,
			AtTime:       batchTime,
			Creator:      s.settings.GasPaymentAddress,
			BaseFee:      s.settings.BaseFee,
			ChainConfig:  s.chainConfig,
			SequencerNo:  sequencerNo,
		}, failForEmptyBatch)
	if err != nil {
		return nil, fmt.Errorf("failed computing batch. Cause: %w", err)
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("failed committing batch state. Cause: %w", err)
	}

	if err := s.signBatch(cb.Batch); err != nil {
		return nil, fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.StoreExecutedBatch(ctx, cb.Batch, cb.Receipts); err != nil {
		return nil, fmt.Errorf("2. failed storing batch. Cause: %w", err)
	}

	s.logger.Info("Produced new batch", log.BatchHashKey, cb.Batch.Hash(),
		"height", cb.Batch.Number(), "numTxs", len(cb.Batch.Transactions), log.BatchSeqNoKey, cb.Batch.SeqNo(), "parent", cb.Batch.Header.ParentHash)

	// add the batch to the chain so it can remove pending transactions from the pool
	err = s.blockchain.IngestNewBlock(cb.Batch)
	if err != nil {
		return nil, fmt.Errorf("failed to feed batch into the virtual eth chain - %w", err)
	}

	return cb, nil
}

// StoreExecutedBatch - stores an executed batch in one go. This can be done for the sequencer because it is guaranteed
// that all dependencies are in place for the execution to be successful.
func (s *sequencer) StoreExecutedBatch(ctx context.Context, batch *core.Batch, receipts types.Receipts) error {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Registry StoreBatch() exit", log.BatchHashKey, batch.Hash())

	// Check if this batch is already stored.
	if _, err := s.storage.FetchBatchHeader(ctx, batch.Hash()); err == nil {
		s.logger.Warn("Attempted to store batch twice! This indicates issues with the batch processing loop")
		return nil
	}

	convertedHeader, err := s.gethEncoding.CreateEthHeaderForBatch(ctx, batch.Header)
	if err != nil {
		return err
	}

	if err := s.storage.StoreBatch(ctx, batch, convertedHeader.Hash()); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	if err := s.storage.StoreExecutedBatch(ctx, batch, receipts); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	s.batchRegistry.OnBatchExecuted(batch, receipts)

	return nil
}

func (s *sequencer) CreateRollup(ctx context.Context, lastBatchNo uint64) (*common.ExtRollup, error) {
	rollupLimiter := limiters.NewRollupLimiter(s.settings.MaxRollupSize)

	currentL1Head, err := s.blockProcessor.GetHead(ctx)
	if err != nil {
		return nil, err
	}
	upToL1Height := currentL1Head.NumberU64() - RollupDelay
	rollup, err := s.rollupProducer.CreateInternalRollup(ctx, lastBatchNo, upToL1Height, rollupLimiter)
	if err != nil {
		return nil, err
	}

	extRollup, err := s.rollupCompression.CreateExtRollup(ctx, rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to compress rollup: %w", err)
	}

	// todo - double-check that this signing approach is secure, and it properly includes the entire payload
	if err := s.signRollup(extRollup); err != nil {
		return nil, fmt.Errorf("failed to sign created rollup: %w", err)
	}

	return extRollup, nil
}

func (s *sequencer) duplicateBatches(ctx context.Context, l1Head *types.Block, nonCanonicalL1Path []common.L1BlockHash, canonicalL1Path []common.L1BlockHash) error {
	batchesToDuplicate := make([]*core.Batch, 0)
	batchesToExclude := make(map[uint64]*core.Batch, 0)

	// read the batches attached to these blocks
	for _, l1BlockHash := range nonCanonicalL1Path {
		batches, err := s.storage.FetchBatchesByBlock(ctx, l1BlockHash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				continue
			}
			return fmt.Errorf("could not FetchBatchesByBlock %s. Cause %w", l1BlockHash, err)
		}
		batchesToDuplicate = append(batchesToDuplicate, batches...)
	}

	// check whether there are already batches on the canonical branch
	for _, l1BlockHash := range canonicalL1Path {
		batches, err := s.storage.FetchBatchesByBlock(ctx, l1BlockHash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				continue
			}
			return fmt.Errorf("could not FetchBatchesByBlock %s. Cause %w", l1BlockHash, err)
		}
		for _, batch := range batches {
			batchesToExclude[batch.NumberU64()] = batch
		}
	}

	if len(batchesToDuplicate) == 0 {
		return nil
	}

	// sort by height
	sort.Slice(batchesToDuplicate, func(i, j int) bool {
		return batchesToDuplicate[i].Number().Cmp(batchesToDuplicate[j].Number()) == -1
	})

	currentHead := batchesToDuplicate[0].Header.ParentHash

	// find all batches for that path
	for i, orphanBatch := range batchesToDuplicate {
		// sanity check that all these batches are consecutive
		if i > 0 && batchesToDuplicate[i].Header.ParentHash != batchesToDuplicate[i-1].Hash() {
			s.logger.Crit("the batches that must be duplicated are invalid")
		}
		if batchesToExclude[orphanBatch.NumberU64()] != nil {
			s.logger.Info("Not duplicating batch because there is already a canonical batch on that height", log.BatchSeqNoKey, orphanBatch.SeqNo())
			currentHead = batchesToExclude[orphanBatch.NumberU64()].Hash()
			continue
		}
		sequencerNo, err := s.storage.FetchCurrentSequencerNo(ctx)
		if err != nil {
			return fmt.Errorf("could not fetch sequencer no. Cause %w", err)
		}
		sequencerNo = sequencerNo.Add(sequencerNo, big.NewInt(1))
		// create the duplicate and store/broadcast it, recreate batch even if it was empty
		cb, err := s.produceBatch(ctx, sequencerNo, l1Head.ParentHash(), currentHead, orphanBatch.Transactions, orphanBatch.Header.Time, false)
		if err != nil {
			return fmt.Errorf("could not produce batch. Cause %w", err)
		}
		currentHead = cb.Batch.Hash()
		s.logger.Info("Duplicated batch", log.BatchHashKey, currentHead)
	}

	// useful for debugging
	//start := batchesToDuplicate[0].SeqNo().Uint64()
	//batches, err := s.storage.FetchNonCanonicalBatchesBetween(ctx, start-1, start+uint64(len(batchesToDuplicate))+1)
	//if err != nil {
	//	panic(err)
	//}
	//for _, batch := range batches {
	//	s.logger.Info("After duplication. Noncanonical", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.Header.SequencerOrderNo)
	//}

	return nil
}

func (s *sequencer) SubmitTransaction(transaction *common.L2Tx) error {
	return s.mempool.Add(transaction)
}

func (s *sequencer) OnL1Fork(ctx context.Context, fork *common.ChainFork) error {
	if !fork.IsFork() {
		return nil
	}

	err := s.duplicateBatches(ctx, fork.NewCanonical, fork.NonCanonicalPath, fork.CanonicalPath)
	if err != nil {
		return fmt.Errorf("could not duplicate batches. Cause %w", err)
	}

	rollup, err := s.storage.FetchReorgedRollup(ctx, fork.NonCanonicalPath)
	if err == nil {
		s.logger.Error("Reissue rollup", log.RollupHashKey, rollup)
		// todo - tudor - finalise the logic to reissue a rollup when the block used for compression was reorged
		return nil
	}
	if !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not call FetchReorgedRollup. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) signBatch(batch *core.Batch) error {
	var err error
	h := batch.Hash()
	batch.Header.Signature, err = signature.Sign(h.Bytes(), s.enclaveKey.PrivateKey())
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (s *sequencer) signRollup(rollup *common.ExtRollup) error {
	var err error
	h := rollup.Header.Hash()
	rollup.Header.Signature, err = signature.Sign(h.Bytes(), s.enclaveKey.PrivateKey())
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (s *sequencer) signCrossChainBundle(bundle *common.ExtCrossChainBundle) error {
	var err error
	h := bundle.HashPacked()
	bundle.Signature, err = signature.Sign(h.Bytes(), s.enclaveKey.PrivateKey())
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (s *sequencer) OnL1Block(ctx context.Context, block *types.Block, result *components.BlockIngestionType) error {
	// nothing to do
	return nil
}

func (s *sequencer) Close() error {
	return s.mempool.Close()
}

func (s *sequencer) ExportCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, error) {
	bundle, err := ExportCrossChainData(ctx, s.storage, fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	err = s.signCrossChainBundle(bundle)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}
