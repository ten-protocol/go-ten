package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

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
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
)

type sequencer struct {
	blockConsumer  components.L1BlockProcessor
	batchProducer  components.BatchProducer
	batchRegistry  components.BatchRegistry
	rollupProducer components.RollupProducer
	rollupConsumer components.RollupConsumer

	logger gethlog.Logger

	hostID            gethcommon.Address
	chainConfig       *params.ChainConfig
	enclavePrivateKey *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool           mempool.Manager
	storage           db.Storage
	encryption        crypto.TransactionBlobCrypto

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
	encryption crypto.TransactionBlobCrypto,
) Sequencer {
	return &sequencer{
		blockConsumer:     consumer,
		batchProducer:     producer,
		batchRegistry:     registry,
		rollupProducer:    rollupProducer,
		rollupConsumer:    rollupConsumer,
		logger:            logger,
		hostID:            hostID,
		chainConfig:       chainConfig,
		enclavePrivateKey: enclavePrivateKey,
		mempool:           mempool,
		storage:           storage,
		encryption:        encryption,
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
	// when a l1HeadBlock is specified it will override this and allow
	// building batches for unfinished forks.
	l1HeadBlock, err := s.blockConsumer.GetHead()
	if err != nil {
		return fmt.Errorf("failed retrieving l1 head. Cause: %w", err)
	}

	if !hasGenesis {
		return s.initGenesis(l1HeadBlock)
	}

	return s.createNewHeadBatch(l1HeadBlock)
}

// TODO - This is iffy, the producer commits the stateDB
func (s *sequencer) initGenesis(block *common.L1Block) error {
	batch, msgBusTx, err := s.batchProducer.CreateGenesisState(block.Hash(), s.hostID, uint64(time.Now().Unix()))
	if err != nil {
		return err
	}

	if err := s.mempool.AddMempoolTx(msgBusTx); err != nil {
		return fmt.Errorf("failed to queue message bus creation transaction to genesis. Cause: %w", err)
	}

	if err := s.signBatch(batch); err != nil {
		return fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.batchRegistry.StoreBatch(batch, nil); err != nil {
		return fmt.Errorf("failed storing batch. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) createNewHeadBatch(l1HeadBlock *common.L1Block) error {
	headBatch, err := s.batchRegistry.GetHeadBatch()
	if err != nil {
		return err
	}

	// We get the latest known batch for this block's chain of parents.
	// This includes the block so if there are any head batches linked to it
	// they will get picked up.
	ancestralBatch, err := s.batchRegistry.FindAncestralBatchFor(l1HeadBlock)
	if err != nil {
		return err
	}

	// If the l1 head block is not a fork then headBatch should
	// be equal to ancestralBatch. Thus they numbers should also be
	// the same. Difference in numbers means ancestor was built on
	// a different chain.
	if ancestralBatch.NumberU64() != headBatch.NumberU64() {
		if err := s.handleFork(l1HeadBlock); err != nil {
			return fmt.Errorf("failed handling fork: Cause: %w", err)
		}
		return s.createNewHeadBatch(l1HeadBlock)
	}

	// After we have determined that the ancestral batch we have is identical to head
	// batch (which can be on another fork) we set the head batch to it as it is guaranteed
	// to be in our chain.
	headBatch = ancestralBatch

	transactions, err := s.mempool.CurrentTxs(headBatch, s.storage)
	if err != nil {
		return err
	}

	// As we are incrementing the chain to a new max height, across all forks,
	// we generate the randomness for this batch.
	rand, err := crypto.GeneratePublicRandomness()
	if err != nil {
		return err
	}

	cb, err := s.batchProducer.ComputeBatch(&components.BatchContext{
		BlockPtr:     l1HeadBlock.Hash(),
		ParentPtr:    *headBatch.Hash(),
		Transactions: transactions,
		AtTime:       uint64(time.Now().Unix()), // todo - time is set only here; take from l1 block?
		Randomness:   gethcommon.BytesToHash(rand),
		Creator:      s.hostID,
		ChainConfig:  s.chainConfig,
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
		return fmt.Errorf("failed storing batch. Cause: %w", err)
	}

	if err := s.mempool.RemoveMempoolTxs(cb.Batch, s.storage); err != nil {
		return fmt.Errorf("could not remove transactions from mempool. Cause: %w", err)
	}

	return nil
}

func (s *sequencer) CreateRollup() (*common.ExtRollup, error) {
	rollup, err := s.rollupProducer.CreateRollup()
	if err != nil {
		return nil, err
	}

	if err := s.signRollup(rollup); err != nil {
		return nil, err
	}

	return rollup.ToExtRollup(s.encryption), nil
}

func (s *sequencer) ReceiveBlock(br *common.BlockAndReceipts, isLatest bool) (*components.BlockIngestionType, error) {
	ingestion, err := s.blockConsumer.Process(br, isLatest)
	if err != nil {
		return nil, err
	}

	if _, err := s.rollupConsumer.ProcessL1Block(br); err != nil {
		s.logger.Error("Encountered error while processing l1 block", log.ErrKey, err)
		// Unsure what to do here; block has been stored
	}

	if !ingestion.Fork {
		return ingestion, nil
	}

	return ingestion, nil
}

func (s *sequencer) handleFork(block *common.L1Block) error {
	headBatch, err := s.batchRegistry.GetHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed retrieving head batch. Cause: %w", err)
	}

	ancestralBatch, err := s.batchRegistry.FindAncestralBatchFor(block)
	if err != nil {
		return fmt.Errorf("failed to find ancestral batch for block: %s", block.Hash())
	}

	if bytes.Equal(headBatch.Header.Hash().Bytes(), ancestralBatch.Hash().Bytes()) {
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

		// Extend the chain with identical cousin batches
		cb, err := s.batchProducer.ComputeBatch(&components.BatchContext{
			BlockPtr:     block.Hash(),
			ParentPtr:    *currHeadPtr.Hash(),
			Transactions: orphan.Transactions,
			AtTime:       orphan.Header.Time,
			Randomness:   orphan.Header.MixDigest,
			Creator:      s.hostID,
			ChainConfig:  s.chainConfig,
		})
		if err != nil {
			s.logger.Crit("Error recalculating l2chain for forked block", log.ErrKey, err)
			return err
		}

		if _, err := cb.Commit(true); err != nil {
			return fmt.Errorf("failed committing stateDB for computed batch. Cause: %w", err)
		}

		if err := s.signBatch(cb.Batch); err != nil {
			return fmt.Errorf("failed signing batch. Cause: %w", err)
		}

		if err := s.batchRegistry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
			return fmt.Errorf("failed storing batch. Cause: %w", err)
		}
		currHeadPtr = cb.Batch

		// i equals 0 at the highest batch number
		if i == 0 {
			if err := s.storage.SetHeadBatchPointer(cb.Batch); err != nil {
				return fmt.Errorf("failed setting head batch ptr. Cause: %w", err)
			}
			return s.storage.UpdateHeadBatch(block.Hash(), cb.Batch, cb.Receipts)
		}
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
