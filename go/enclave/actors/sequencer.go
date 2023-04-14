package actors

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
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
	consumer       components.BlockConsumer
	producer       components.BatchProducer
	registry       components.BatchRegistry
	rollupProducer components.RollupProducer
	rollupConsumer components.RollupConsumer

	logger gethlog.Logger

	hostID            gethcommon.Address
	chainConfig       *params.ChainConfig
	enclavePrivateKey *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	mempool           mempool.Manager
	storage           db.Storage
	encryption        crypto.TransactionBlobCrypto
}

func NewSequencer(
	consumer components.BlockConsumer,
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
		consumer:          consumer,
		producer:          producer,
		registry:          registry,
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

func (s *sequencer) IsReady() bool {
	return false
}

func (s *sequencer) CreateBatch(block *common.L1Block) (*core.Batch, error) {
	hasGenesis, err := s.registry.HasGenesisBatch()
	if err != nil {
		return nil, fmt.Errorf("unknown genesis batch state. Cause: %w", err)
	}

	if block == nil {
		// L1 Head is only updated when isLatest: true
		// when a block is specified it will override this and allow
		// building batches for unfinished forks.
		block, err = s.consumer.GetHead()
		if err != nil {
			return nil, fmt.Errorf("failed retrieving l1 head. Cause: %w", err)
		}
	}

	if !hasGenesis {
		return s.initGenesis(block)
	}

	return s.extendHead(block)
}

// TODO - This is iffy, the producer commits the stateDB
func (s *sequencer) initGenesis(block *common.L1Block) (*core.Batch, error) {
	batch, msgBusTx, err := s.producer.CreateGenesisState(block.Hash(), s.hostID, uint64(time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	if err := s.mempool.AddMempoolTx(msgBusTx); err != nil {
		return nil, fmt.Errorf("failed to queue message bus creation transaction to genesis. Cause: %w", err)
	}

	if err := s.signBatch(batch); err != nil {
		return nil, fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.registry.StoreBatch(batch, nil); err != nil {
		return nil, fmt.Errorf("failed storing batch. Cause: %w", err)
	}

	return batch, nil
}

func (s *sequencer) extendHead(block *common.L1Block) (*core.Batch, error) {
	headBatch, err := s.registry.GetHeadBatch()
	if err != nil {
		return nil, err
	}

	// We find the ancestral batch for the block as it might be different
	// than the current head; Alternatively we might have just processed a fork
	// which would've updated the head to an unfinished fork so we need to get
	// the correct batch that is building on the latest known final chain
	ancestralBatch, err := s.registry.FindAncestralBatchFor(block)
	if err != nil {
		return nil, err
	}

	if ancestralBatch.NumberU64() != headBatch.NumberU64() {
		return nil, fmt.Errorf("cant continue head for block that has decadent ancestor. Cause: %w", err)
	}

	// After we have determined that the ancestral batch we have is identical to head
	// batch (which can be on another fork) we set the head batch to it as it is guaranteed
	// to be in our chain.
	headBatch = ancestralBatch

	transactions, err := s.mempool.CurrentTxs(headBatch, s.storage)
	if err != nil {
		return nil, err
	}

	// As we are incrementing the chain to a new max height, across all forks,
	// we generate the randomness for this batch.
	rand, err := crypto.GeneratePublicRandomness()
	if err != nil {
		return nil, err
	}

	cb, err := s.producer.ComputeBatch(&components.BatchContext{
		BlockPtr:     block.Hash(),
		ParentPtr:    *headBatch.Hash(),
		Transactions: transactions,
		AtTime:       uint64(time.Now().Unix()), // todo - time is set only here; take from l1 block?
		Randomness:   gethcommon.BytesToHash(rand),
		Creator:      s.hostID,
		ChainConfig:  s.chainConfig,
	})

	if err != nil {
		return nil, fmt.Errorf("failed computing batch. Cause: %w", err)
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("failed commiting batch state. Cause: %w", err)
	}

	if err := s.signBatch(cb.Batch); err != nil {
		return nil, fmt.Errorf("failed signing created batch. Cause: %w", err)
	}

	if err := s.registry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
		return nil, fmt.Errorf("failed storing batch. Cause: %w", err)
	}

	return cb.Batch, nil
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
	ingestion, err := s.consumer.ConsumeBlock(br, isLatest)
	if err != nil {
		return nil, err
	}

	s.rollupConsumer.ProcessL1Block(br)

	if !ingestion.Fork {
		return ingestion, nil
	}

	if err := s.handleFork(br); err != nil {
		return nil, fmt.Errorf("failed handling fork: Cause: %w", err)
	}

	return ingestion, nil
}

func (s *sequencer) handleFork(br *common.BlockAndReceipts) error {
	headBatch, err := s.registry.GetHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed retrieving head batch. Cause: %w", err)
	}

	ancestralBatch, err := s.registry.FindAncestralBatchFor(br.Block)
	if err != nil {
		return fmt.Errorf("failed to find ancestral batch for block: %s", br.Block.Hash())
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
		currHead, err = s.registry.GetBatch(currHead.Header.ParentHash)
		if err != nil {
			s.logger.Crit("Failure while looking for previously stored batch!", log.ErrKey, err)
			return err
		}
	}

	currHeadPtr := ancestralBatch
	for i := len(orphanedBatches) - 1; i >= 0; i-- {
		orphan := orphanedBatches[i]

		// Extend the chain with identical cousin batches
		cb, err := s.producer.ComputeBatch(&components.BatchContext{
			BlockPtr:     br.Block.Hash(),
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
			return fmt.Errorf("failed commiting stateDB for computed batch. Cause: %w", err)
		}

		if err := s.signBatch(cb.Batch); err != nil {
			return fmt.Errorf("failed signing batch. Cause: %w", err)
		}

		if err := s.registry.StoreBatch(cb.Batch, cb.Receipts); err != nil {
			return fmt.Errorf("failed storing batch. Cause: %w", err)
		}
		currHeadPtr = cb.Batch

		if i == 0 {
			s.storage.SetHeadBatchPointer(cb.Batch)
			return s.storage.UpdateHeadBatch(br.Block.Hash(), cb.Batch, cb.Receipts)
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
