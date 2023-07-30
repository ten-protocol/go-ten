package nodetype

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"

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
	storage                storage.Storage
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
	storage storage.Storage,
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

	// todo - time is set only here; take from l1 block?
	// when creating a new head batch, it is pointing to the parent of the current l1 head
	// the reason for this is to minimize the chance of creating batches on top of blocks that will not be seen by the network
	// todo - to fix in a follow up PR
	if _, err := s.produceBatch(sequencerNo.Add(sequencerNo, big.NewInt(1)), l1HeadBlock.ParentHash(), headBatch.Hash(), transactions, uint64(time.Now().Unix())); err != nil {
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
		Creator:      s.hostID,
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

	if err := s.batchRegistry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
		return nil, fmt.Errorf("2. failed storing batch. Cause: %w", err)
	}

	s.logger.Info("Produced new batch", log.BatchHashKey, cb.Batch.Hash(),
		"height", cb.Batch.Number(), "numTxs", len(cb.Batch.Transactions), "seqNo", cb.Batch.SeqNo())

	return cb.Batch, nil
}

func (s *sequencer) CreateRollup(lastBatchNo uint64) (*common.ExtRollup, error) {
	// todo @stefan - move this somewhere else, it shouldn't be in the batch registry.
	rollupLimiter := limiters.NewRollupLimiter(s.settings.MaxRollupSize)

	rollup, err := s.rollupProducer.CreateRollup(lastBatchNo, rollupLimiter)
	if err != nil {
		return nil, err
	}

	if err := s.signRollup(rollup); err != nil {
		return nil, errors.Wrap(err, "failed to sign created rollup")
	}

	s.logger.Info("Created new head rollup", log.RollupHashKey, rollup.Hash(), "numBatches", len(rollup.Batches))

	return rollup.ToExtRollup(s.dataEncryptionService, s.dataCompressionService)
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

	sequencerNo, err := s.storage.FetchCurrentSequencerNo()
	if err != nil {
		return fmt.Errorf("could not fetch sequencer no. Cause %w", err)
	}

	currentHead := batchesToDuplicate[0].Header.ParentHash

	// find all batches for that path
	for i, orphanBatch := range batchesToDuplicate {
		// sanity check that all these batches are consecutive
		if i > 0 && !bytes.Equal(batchesToDuplicate[i].Header.ParentHash.Bytes(), batchesToDuplicate[i-1].Hash().Bytes()) {
			s.logger.Crit("the batches that must be duplicated are invalid")
		}
		sequencerNo = sequencerNo.Add(sequencerNo, big.NewInt(1))
		// create the duplicate and store/broadcast it
		b, err := s.produceBatch(sequencerNo, l1Head.ParentHash(), currentHead, orphanBatch.Transactions, orphanBatch.Header.Time)
		currentHead = b.Hash()
		if err != nil {
			return fmt.Errorf("could not produce batch. Cause %w", err)
		}
		s.logger.Info("Duplicated batch", log.BatchHashKey, currentHead)
	}

	return nil
}

func (s *sequencer) SubmitTransaction(transaction *common.L2Tx) error {
	return s.mempool.AddMempoolTx(transaction)
}

func (s *sequencer) OnL1Fork(fork *common.ChainFork) error {
	if fork.IsFork() {
		err := s.duplicateBatches(fork.NewCanonical, fork.NonCanonicalPath)
		if err != nil {
			return fmt.Errorf("could not duplicate batches. Cause %w", err)
		}
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
