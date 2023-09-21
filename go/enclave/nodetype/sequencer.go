package nodetype

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/measure"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
)

const RollupDelay = 2 // number of L1 blocks to exclude when creating a rollup. This will minimize compression reorg issues.

type SequencerSettings struct {
	MaxBatchSize      uint64
	MaxRollupSize     uint64
	GasPaymentAddress gethcommon.Address
	BatchGasLimit     *big.Int
	BaseFee           *big.Int
}

type sequencer struct {
	blockProcessor    components.L1BlockProcessor
	batchProducer     components.BatchExecutor
	batchRegistry     components.BatchRegistry
	rollupProducer    components.RollupProducer
	rollupConsumer    components.RollupConsumer
	rollupCompression *components.RollupCompression

	logger gethlog.Logger

	hostID                 gethcommon.Address
	chainConfig            *params.ChainConfig
	enclavePrivateKey      *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool                mempool.Manager
	storage                storage.Storage
	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	settings               SequencerSettings
}

func NewSequencer(
	blockProcessor components.L1BlockProcessor,
	batchExecutor components.BatchExecutor,
	registry components.BatchRegistry,
	rollupProducer components.RollupProducer,
	rollupConsumer components.RollupConsumer,
	rollupCompression *components.RollupCompression,

	logger gethlog.Logger,

	hostID gethcommon.Address,
	chainConfig *params.ChainConfig,
	enclavePrivateKey *ecdsa.PrivateKey, // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool mempool.Manager,
	storage storage.Storage,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	settings SequencerSettings,
) Sequencer {
	return &sequencer{
		blockProcessor:         blockProcessor,
		batchProducer:          batchExecutor,
		batchRegistry:          registry,
		rollupProducer:         rollupProducer,
		rollupConsumer:         rollupConsumer,
		rollupCompression:      rollupCompression,
		logger:                 logger,
		hostID:                 hostID,
		chainConfig:            chainConfig,
		enclavePrivateKey:      enclavePrivateKey,
		mempool:                mempool,
		storage:                storage,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		settings:               settings,
	}
}

func (s *sequencer) CreateBatch() error {
	hasGenesis, err := s.batchRegistry.HasGenesisBatch()
	if err != nil {
		return fmt.Errorf("unknown genesis batch state. Cause: %w", err)
	}

	// L1 Head is only updated when isLatest: true
	l1HeadBlock, err := s.blockProcessor.GetHead()
	if err != nil {
		return fmt.Errorf("failed retrieving l1 head. Cause: %w", err)
	}

	if !hasGenesis {
		return s.initGenesis(l1HeadBlock)
	}

	return s.createNewHeadBatch(l1HeadBlock)
}

// TODO - This is iffy, the producer commits the stateDB. The producer
// should only create batches and stateDBs but not commit them to the database,
// this is the responsibility of the sequencer. Refactor the code so genesis state
// won't be committed by the producer.
func (s *sequencer) initGenesis(block *common.L1Block) error {
	s.logger.Info("Initializing genesis state", log.BlockHashKey, block.Hash())
	batch, msgBusTx, err := s.batchProducer.CreateGenesisState(
		block.Hash(),
		uint64(time.Now().Unix()),
		s.settings.GasPaymentAddress,
		s.settings.BaseFee,
		s.settings.BatchGasLimit,
	)
	if err != nil {
		return err
	}

	if err = s.mempool.AddMempoolTx(msgBusTx); err != nil {
		return fmt.Errorf("failed to queue message bus creation transaction to genesis. Cause: %w", err)
	}

	if err := s.signBatch(batch); err != nil {
		return fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.StoreExecutedBatch(batch, nil); err != nil {
		return fmt.Errorf("1. failed storing batch. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) createNewHeadBatch(l1HeadBlock *common.L1Block) error {
	headBatchSeq := s.batchRegistry.HeadBatchSeq()
	headBatch, err := s.storage.FetchBatchBySeqNo(headBatchSeq.Uint64())
	if err != nil {
		return err
	}

	// todo - sanity check that the headBatch.Header.L1Proof is an ancestor of the l1HeadBlock
	b, err := s.storage.FetchBlock(headBatch.Header.L1Proof)
	if err != nil {
		return err
	}
	if !s.storage.IsAncestor(l1HeadBlock, b) {
		return fmt.Errorf("attempted to create batch on top of batch=%s. With l1 head=%s", headBatch.Hash(), l1HeadBlock.Hash())
	}

	stateDB, err := s.storage.CreateStateDB(headBatch.Hash())
	if err != nil {
		return fmt.Errorf("unable to create stateDB for selecting transactions. Batch: %s Cause: %w", headBatch.Hash(), err)
	}

	// todo (@stefan) - limit on receipts too
	limiter := limiters.NewBatchSizeLimiter(s.settings.MaxBatchSize)
	transactions, err := s.mempool.CurrentTxs(stateDB, limiter)
	if err != nil {
		return err
	}

	sequencerNo, err := s.storage.FetchCurrentSequencerNo()
	if err != nil {
		return err
	}

	// todo - time is set only here; take from l1 block?
	if _, err := s.produceBatch(sequencerNo.Add(sequencerNo, big.NewInt(1)), l1HeadBlock.Hash(), headBatch.Hash(), transactions, uint64(time.Now().Unix())); err != nil {
		return fmt.Errorf(" failed producing batch. Cause: %w", err)
	}

	if err := s.mempool.RemoveTxs(transactions); err != nil {
		return fmt.Errorf("could not remove transactions from mempool. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) produceBatch(sequencerNo *big.Int, l1Hash common.L1BlockHash, headBatch common.L2BatchHash, transactions common.L2Transactions, batchTime uint64) (*core.Batch, error) {
	cb, err := s.batchProducer.ComputeBatch(&components.BatchExecutionContext{
		BlockPtr:     l1Hash,
		ParentPtr:    headBatch,
		Transactions: transactions,
		AtTime:       batchTime,
		Creator:      s.settings.GasPaymentAddress,
		BaseFee:      s.settings.BaseFee,
		ChainConfig:  s.chainConfig,
		SequencerNo:  sequencerNo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed computing batch. Cause: %w", err)
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("failed committing batch state. Cause: %w", err)
	}

	if err := s.signBatch(cb.Batch); err != nil {
		return nil, fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.StoreExecutedBatch(cb.Batch, cb.Receipts); err != nil {
		return nil, fmt.Errorf("2. failed storing batch. Cause: %w", err)
	}

	s.logger.Info("Produced new batch", log.BatchHashKey, cb.Batch.Hash(),
		"height", cb.Batch.Number(), "numTxs", len(cb.Batch.Transactions), log.BatchSeqNoKey, cb.Batch.SeqNo(), "parent", cb.Batch.Header.ParentHash)

	return cb.Batch, nil
}

// StoreExecutedBatch - stores an executed batch in one go. This can be done for the sequencer because it is guaranteed
// that all dependencies are in place for the execution to be successful.
func (s *sequencer) StoreExecutedBatch(batch *core.Batch, receipts types.Receipts) error {
	defer s.logger.Info("Registry StoreBatch() exit", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	// Check if this batch is already stored.
	if _, err := s.storage.FetchBatchHeader(batch.Hash()); err == nil {
		s.logger.Warn("Attempted to store batch twice! This indicates issues with the batch processing loop")
		return nil
	}

	if err := s.storage.StoreBatch(batch); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	if err := s.storage.StoreExecutedBatch(batch, receipts); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	s.batchRegistry.OnBatchExecuted(batch, receipts)

	return nil
}

func (s *sequencer) CreateRollup(lastBatchNo uint64) (*common.ExtRollup, error) {
	rollupLimiter := limiters.NewRollupLimiter(s.settings.MaxRollupSize)

	currentL1Head, err := s.storage.FetchHeadBlock()
	if err != nil {
		return nil, err
	}
	upToL1Height := currentL1Head.NumberU64() - RollupDelay
	rollup, err := s.rollupProducer.CreateRollup(lastBatchNo, upToL1Height, rollupLimiter)
	if err != nil {
		return nil, err
	}

	if err := s.signRollup(rollup); err != nil {
		return nil, fmt.Errorf("failed to sign created rollup: %w", err)
	}

	s.logger.Info("Created new head rollup", log.RollupHashKey, rollup.Hash(), "numBatches", len(rollup.Batches))

	return s.rollupCompression.CreateExtRollup(rollup)
}

func (s *sequencer) duplicateBatches(l1Head *types.Block, nonCanonicalL1Path []common.L1BlockHash) error {
	batchesToDuplicate := make([]*core.Batch, 0)

	// read the batches attached to these blocks
	for _, l1BlockHash := range nonCanonicalL1Path {
		batches, err := s.storage.FetchBatchesByBlock(l1BlockHash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				continue
			}
			return fmt.Errorf("could not FetchBatchesByBlock %s. Cause %w", l1BlockHash, err)
		}
		batchesToDuplicate = append(batchesToDuplicate, batches...)
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
		sequencerNo, err := s.storage.FetchCurrentSequencerNo()
		if err != nil {
			return fmt.Errorf("could not fetch sequencer no. Cause %w", err)
		}
		sequencerNo = sequencerNo.Add(sequencerNo, big.NewInt(1))
		// create the duplicate and store/broadcast it
		b, err := s.produceBatch(sequencerNo, l1Head.ParentHash(), currentHead, orphanBatch.Transactions, orphanBatch.Header.Time)
		if err != nil {
			return fmt.Errorf("could not produce batch. Cause %w", err)
		}
		currentHead = b.Hash()
		s.logger.Info("Duplicated batch", log.BatchHashKey, currentHead)
	}

	return nil
}

func (s *sequencer) SubmitTransaction(transaction *common.L2Tx) error {
	return s.mempool.AddMempoolTx(transaction)
}

func (s *sequencer) OnL1Fork(fork *common.ChainFork) error {
	if !fork.IsFork() {
		return nil
	}

	err := s.duplicateBatches(fork.NewCanonical, fork.NonCanonicalPath)
	if err != nil {
		return fmt.Errorf("could not duplicate batches. Cause %w", err)
	}

	rollup, err := s.storage.FetchReorgedRollup(fork.NonCanonicalPath)
	if err == nil {
		s.logger.Error("Reissue rollup", log.RollupHashKey, rollup)
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
	batch.Header.R, batch.Header.S, err = ecdsa.Sign(rand.Reader, s.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (s *sequencer) signRollup(rollup *core.Rollup) error {
	var err error
	h := rollup.Header.Hash()
	rollup.Header.R, rollup.Header.S, err = ecdsa.Sign(rand.Reader, s.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (s *sequencer) OnL1Block(_ types.Block, _ *components.BlockIngestionType) error {
	// nothing to do
	return nil
}
