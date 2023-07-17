package nodetype

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
)

type SequencerSettings struct {
	MaxBatchSize  uint64
	MaxRollupSize uint64
}

type sequencer struct {
	blockProcessor components.L1BlockProcessor
	batchProducer  components.BatchProducer
	batchRegistry  components.BatchRegistry
	rollupProducer components.RollupProducer
	rollupConsumer components.RollupConsumer

	logger gethlog.Logger

	hostID                 gethcommon.Address
	chainConfig            *params.ChainConfig
	enclavePrivateKey      *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool                mempool.Manager
	storage                db.Storage
	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	settings               SequencerSettings

	// This is used to coordinate creating
	// new batches and creating fork batches.
	batchProductionMutex sync.Mutex
}

func NewSequencer(
	consumer components.L1BlockProcessor,
	producer components.BatchProducer,
	registry components.BatchRegistry,
	rollupProducer components.RollupProducer,
	rollupConsumer components.RollupConsumer,

	logger gethlog.Logger,

	hostID gethcommon.Address,
	chainConfig *params.ChainConfig,
	enclavePrivateKey *ecdsa.PrivateKey, // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool mempool.Manager,
	storage db.Storage,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	settings SequencerSettings,
) Sequencer {
	return &sequencer{
		blockProcessor:         consumer,
		batchProducer:          producer,
		batchRegistry:          registry,
		rollupProducer:         rollupProducer,
		rollupConsumer:         rollupConsumer,
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
	s.batchProductionMutex.Lock()
	defer s.batchProductionMutex.Unlock()

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
	batch, msgBusTx, err := s.batchProducer.CreateGenesisState(block.Hash(), uint64(time.Now().Unix()))
	if err != nil {
		return err
	}

	if err = s.mempool.AddMempoolTx(msgBusTx); err != nil {
		return fmt.Errorf("failed to queue message bus creation transaction to genesis. Cause: %w", err)
	}

	if err := s.signBatch(batch); err != nil {
		return fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.batchRegistry.StoreBatch(batch, nil); err != nil {
		return fmt.Errorf("1. failed storing batch. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) createNewHeadBatch(l1HeadBlock *common.L1Block) error {
	headBatch, err := s.batchRegistry.GetHeadBatch()
	if err != nil {
		return err
	}

	// todo - as an optimisation we could check here that the current headBatch is on the incoming l1 chain
	// this is the happy path after all
	// if s.storage.IsAncestor() headBatch.Header.L1Proof

	// We get the latest known batch for this block's chain of parents.
	// This includes the block so if there are any head batches linked to it
	// they will get picked up.
	ancestralBatch, err := s.batchRegistry.FindAncestralBatchFor(l1HeadBlock)
	if err != nil {
		return err
	}

	// If the l1 head block is not a fork then headBatch should
	// be equal to ancestralBatch. Thus, their numbers should also be
	// the same. If there is a difference, it means the ancestor was built on
	// a different chain.
	if ancestralBatch.NumberU64() != headBatch.NumberU64() {
		if err := s.handleFork(l1HeadBlock, ancestralBatch); err != nil {
			return fmt.Errorf("failed handling fork: Cause: %w", err)
		}
		return s.createNewHeadBatch(l1HeadBlock)
	}

	// After we have determined that the ancestral batch we have is identical to head
	// batch (which can be on another fork) we set the head batch to it as it is guaranteed
	// to be in our chain.
	headBatch = ancestralBatch

	stateDB, err := s.storage.CreateStateDB(headBatch.Hash())
	if err != nil {
		return fmt.Errorf("unable to create stateDB for selecting transactions. Cause: %w", err)
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
	cb, err := s.batchProducer.ComputeBatch(&components.BatchExecutionContext{
		BlockPtr:     l1HeadBlock.Hash(),
		ParentPtr:    headBatch.Hash(),
		Transactions: transactions,
		AtTime:       uint64(time.Now().Unix()), // todo - time is set only here; take from l1 block?
		Creator:      s.hostID,
		ChainConfig:  s.chainConfig,
		SequencerNo:  sequencerNo.Add(sequencerNo, big.NewInt(1)),
	})
	if err != nil {
		return fmt.Errorf("failed computing batch. Cause: %w", err)
	}

	if _, err := cb.Commit(true); err != nil {
		return fmt.Errorf("failed committing batch state. Cause: %w", err)
	}

	if err := s.signBatch(cb.Batch); err != nil {
		return fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.batchRegistry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
		return fmt.Errorf("2. failed storing batch. Cause: %w", err)
	}

	if err := s.mempool.RemoveTxs(transactions); err != nil {
		return fmt.Errorf("could not remove transactions from mempool. Cause: %w", err)
	}

	s.logger.Info("Created new head batch", log.BatchHashKey, cb.Batch.Hash(),
		"height", cb.Batch.Number(), "numTxs", len(cb.Batch.Transactions), "seqNo", cb.Batch.SeqNo())

	return nil
}

func (s *sequencer) CreateRollup() (*common.ExtRollup, error) {
	// todo @stefan - move this somewhere else, it shouldn't be in the batch registry.
	rollupLimiter := limiters.NewRollupLimiter(s.settings.MaxRollupSize)

	rollup, err := s.rollupProducer.CreateRollup(rollupLimiter)
	if err != nil {
		return nil, err
	}

	if err := s.signRollup(rollup); err != nil {
		return nil, err
	}

	s.logger.Info("Created new head rollup", log.RollupHashKey, rollup.Hash(), log.RollupHeightKey, rollup.Number(), "numBatches", len(rollup.Batches))

	return rollup.ToExtRollup(s.dataEncryptionService, s.dataCompressionService)
}

func (s *sequencer) handleFork(block *common.L1Block, ancestralBatch *core.Batch) error {
	headBatch, err := s.batchRegistry.GetHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed retrieving head batch. Cause: %w", err)
	}

	if bytes.Equal(headBatch.Hash().Bytes(), ancestralBatch.Hash().Bytes()) {
		return nil
	}

	if headBatch.NumberU64() < ancestralBatch.NumberU64() {
		panic("fork should never resolve to a higher height batch...")
	}

	currHead := headBatch
	orphanedBatches := make([]*core.Batch, 0)
	for currHead.NumberU64() > ancestralBatch.NumberU64() {
		orphanedBatches = append(orphanedBatches, currHead)
		currHead, err = s.batchRegistry.GetBatch(currHead.Header.ParentHash)
		if err != nil {
			s.logger.Crit("Failure while looking for previously stored batch!", log.ErrKey, err)
			return err
		}
	}

	currHeadPtr := ancestralBatch
	for i := len(orphanedBatches) - 1; i >= 0; i-- {
		orphan := orphanedBatches[i]

		sequencerNo, err := s.storage.FetchCurrentSequencerNo()
		if err != nil {
			return err
		}

		// Extend the chain with identical cousin batches
		cb, err := s.batchProducer.ComputeBatch(&components.BatchExecutionContext{
			BlockPtr:     block.Hash(),
			ParentPtr:    currHeadPtr.Hash(),
			Transactions: orphan.Transactions,
			AtTime:       orphan.Header.Time,
			Creator:      s.hostID,
			ChainConfig:  s.chainConfig,
			SequencerNo:  sequencerNo.Add(sequencerNo, big.NewInt(1)),
		})
		if err != nil {
			s.logger.Crit("Error recalculating l2chain for forked block", log.ErrKey, err)
			return err
		}

		s.logger.Info(fmt.Sprintf("Produced fork batch %s with seqNo %d", cb.Batch.Hash(), cb.Batch.SeqNo().Uint64()))

		if _, err := cb.Commit(true); err != nil {
			return fmt.Errorf("failed committing stateDB for computed batch. Cause: %w", err)
		}

		if err := s.signBatch(cb.Batch); err != nil {
			return fmt.Errorf("failed signing batch. Cause: %w", err)
		}

		if err := s.batchRegistry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
			return fmt.Errorf("3. failed storing batch. Cause: %w", err)
		}
		currHeadPtr = cb.Batch
	}

	return nil
}

func (s *sequencer) SubmitTransaction(transaction *common.L2Tx) error {
	return s.mempool.AddMempoolTx(transaction)
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
