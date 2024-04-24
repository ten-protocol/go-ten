package nodetype

import (
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
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ten-protocol/go-ten/go/common/compression"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
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

func (s *sequencer) CreateBatch(skipBatchIfEmpty bool) error {
	hasGenesis, err := s.batchRegistry.HasGenesisBatch()
	if err != nil {
		return fmt.Errorf("unknown genesis batch state. Cause: %w", err)
	}

	// L1 Head is only updated when isLatest: true
	l1HeadBlock, err := s.blockProcessor.GetHead()
	if err != nil {
		return fmt.Errorf("failed retrieving l1 head. Cause: %w", err)
	}

	// the sequencer creates the initial genesis batch if one does not exist yet
	if !hasGenesis {
		return s.createGenesisBatch(l1HeadBlock)
	}

	if running := s.mempool.Running(); !running {
		// the mempool can only be started after at least 1 block (the genesis) is in the blockchain object
		// if the node restarted the mempool must be started again
		err = s.mempool.Start()
		if err != nil {
			return err
		}
	}

	return s.createNewHeadBatch(l1HeadBlock, skipBatchIfEmpty)
}

// TODO - This is iffy, the producer commits the stateDB. The producer
// should only create batches and stateDBs but not commit them to the database,
// this is the responsibility of the sequencer. Refactor the code so genesis state
// won't be committed by the producer.
func (s *sequencer) createGenesisBatch(block *common.L1Block) error {
	s.logger.Info("Initializing genesis state", log.BlockHashKey, block.Hash())
	batch, msgBusTx, err := s.batchProducer.CreateGenesisState(
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

	if err := s.StoreExecutedBatch(batch, nil); err != nil {
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

func (s *sequencer) createNewHeadBatch(l1HeadBlock *common.L1Block, skipBatchIfEmpty bool) error {
	headBatchSeq := s.batchRegistry.HeadBatchSeq()
	if headBatchSeq == nil {
		headBatchSeq = big.NewInt(int64(common.L2GenesisSeqNo))
	}
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

	sequencerNo, err := s.storage.FetchCurrentSequencerNo()
	if err != nil {
		return err
	}

	// todo - time is set only here; take from l1 block?
	if _, err := s.produceBatch(sequencerNo.Add(sequencerNo, big.NewInt(1)), l1HeadBlock.Hash(), headBatch.Hash(), transactions, uint64(time.Now().Unix()), skipBatchIfEmpty); err != nil {
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
	sequencerNo *big.Int,
	l1Hash common.L1BlockHash,
	headBatch common.L2BatchHash,
	transactions common.L2Transactions,
	batchTime uint64,
	failForEmptyBatch bool,
) (*components.ComputedBatch, error) {
	cb, err := s.batchProducer.ComputeBatch(&components.BatchExecutionContext{
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

	if err := s.StoreExecutedBatch(cb.Batch, cb.Receipts); err != nil {
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
func (s *sequencer) StoreExecutedBatch(batch *core.Batch, receipts types.Receipts) error {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Registry StoreBatch() exit", log.BatchHashKey, batch.Hash())

	// Check if this batch is already stored.
	if _, err := s.storage.FetchBatchHeader(batch.Hash()); err == nil {
		s.logger.Warn("Attempted to store batch twice! This indicates issues with the batch processing loop")
		return nil
	}

	convertedHeader, err := s.gethEncoding.CreateEthHeaderForBatch(batch.Header)
	if err != nil {
		return err
	}

	if err := s.storage.StoreBatch(batch, convertedHeader.Hash()); err != nil {
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

	currentL1Head, err := s.blockProcessor.GetHead()
	if err != nil {
		return nil, err
	}
	upToL1Height := currentL1Head.NumberU64() - RollupDelay
	rollup, err := s.rollupProducer.CreateInternalRollup(lastBatchNo, upToL1Height, rollupLimiter)
	if err != nil {
		return nil, err
	}

	extRollup, err := s.rollupCompression.CreateExtRollup(rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to compress rollup: %w", err)
	}

	// todo - double-check that this signing approach is secure, and it properly includes the entire payload
	if err := s.signRollup(extRollup); err != nil {
		return nil, fmt.Errorf("failed to sign created rollup: %w", err)
	}

	return extRollup, nil
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
		// create the duplicate and store/broadcast it, recreate batch even if it was empty
		cb, err := s.produceBatch(sequencerNo, l1Head.ParentHash(), currentHead, orphanBatch.Transactions, orphanBatch.Header.Time, false)
		if err != nil {
			return fmt.Errorf("could not produce batch. Cause %w", err)
		}
		currentHead = cb.Batch.Hash()
		s.logger.Info("Duplicated batch", log.BatchHashKey, currentHead)
	}

	return nil
}

func (s *sequencer) SubmitTransaction(transaction *common.L2Tx) error {
	return s.mempool.Add(transaction)
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

func (s *sequencer) OnL1Block(_ types.Block, _ *components.BlockIngestionType) error {
	// nothing to do
	return nil
}

func (s *sequencer) Close() error {
	return s.mempool.Close()
}

func (s *sequencer) ExportCrossChainData(fromSeqNo uint64, toSeqNo uint64) (*common.ExtCrossChainBundle, error) {
	canonicalBatchesInRollup := make([]*core.Batch, 0)
	for i := fromSeqNo; i <= toSeqNo; i++ {
		batch, err := s.storage.FetchBatchBySeqNo(fromSeqNo)
		if err != nil {
			return nil, err
		}

		l1BlockHash := batch.Header.L1Proof
		block, err := s.storage.FetchBlock(l1BlockHash)
		if err != nil {
			return nil, err
		}

		canonicalBlock, err := s.storage.FetchCanonicaBlockByHeight(block.Header().Number)
		if err != nil {
			return nil, err
		}

		if canonicalBlock.Hash().Cmp(block.Hash()) != 0 {
			continue
		}

		// Only add batches that point to canonical blocks.
		canonicalBatchesInRollup = append(canonicalBatchesInRollup, batch)
	}

	if len(canonicalBatchesInRollup) == 0 {
		return nil, fmt.Errorf("no batches found for export of cross chain data")
	}

	// build a merkle tree of all the batches that are valid for the cannonical L1 chain
	// The proof is double inclusion - one for the message being in the batch's tree and one for

	smtValues := crosschain.MerkleBatches(canonicalBatchesInRollup).ForMerkleTree()

	tree, err := smt.Of(smtValues, []string{smt.SOL_BYTES32})
	if err != nil {
		return nil, err
	}

	batchesHash := tree.GetRoot()

	block, err := s.blockProcessor.GetHead()
	if err != nil {
		return nil, err
	}

	crossChainHashes := make([][]byte, 0)
	for _, batch := range canonicalBatchesInRollup {
		if batch.Header.TransfersTree != gethcommon.BigToHash(gethcommon.Big0) {
			crossChainHashes = append(crossChainHashes, batch.Header.TransfersTree.Bytes())
		}
	}

	bundle := &common.ExtCrossChainBundle{
		StateRootHash:    gethcommon.BytesToHash(batchesHash),
		L1BlockHash:      block.Hash(),
		L1BlockNum:       big.NewInt(0).Set(block.Header().Number),
		CrossChainHashes: crossChainHashes,
	}

	err = s.signCrossChainBundle(bundle)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}
