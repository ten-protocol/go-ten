package nodetype

import (
	"bytes"
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type obsValidator struct {
	blockProcessor components.L1BlockProcessor
	batchProducer  components.BatchProducer
	batchRegistry  components.BatchRegistry
	rollupConsumer components.RollupConsumer

	chainConfig *params.ChainConfig

	sequencerID  gethcommon.Address
	storage      db.Storage
	sigValidator *components.SignatureValidator
	logger       gethlog.Logger
}

func NewValidator(
	consumer components.L1BlockProcessor,
	producer components.BatchProducer,
	registry components.BatchRegistry,
	rollupConsumer components.RollupConsumer,

	chainConfig *params.ChainConfig,

	sequencerID gethcommon.Address,
	storage db.Storage,
	sigValidator *components.SignatureValidator,
	logger gethlog.Logger,
) ObsValidator {
	return &obsValidator{
		blockProcessor: consumer,
		batchProducer:  producer,
		batchRegistry:  registry,
		rollupConsumer: rollupConsumer,
		chainConfig:    chainConfig,
		sequencerID:    sequencerID,
		storage:        storage,
		sigValidator:   sigValidator,
		logger:         logger,
	}
}

func (val *obsValidator) handleGenesisBatch(incomingBatch *core.Batch) (bool, error) {
	batch, _, err := val.batchProducer.CreateGenesisState(incomingBatch.Header.L1Proof, incomingBatch.Header.Time)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(incomingBatch.Hash().Bytes(), batch.Hash().Bytes()) {
		return false, fmt.Errorf("received bad genesis batch")
	}

	return true, val.batchRegistry.StoreBatch(incomingBatch, nil)
}

func (val *obsValidator) ValidateAndStoreBatch(incomingBatch *core.Batch) error {
	if incomingBatch.NumberU64() == 0 {
		if handled, err := val.handleGenesisBatch(incomingBatch); handled {
			return err
		}
	}

	defer val.logger.Info("Validator processed batch", log.BatchHashKey, incomingBatch.Hash(), log.DurationKey, measure.NewStopwatch())

	if batch, err := val.batchRegistry.GetBatch(incomingBatch.Hash()); err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return err
	} else if batch != nil {
		return nil // already know about this one
	}

	if err := val.sigValidator.CheckSequencerSignature(incomingBatch.Hash(), incomingBatch.Header.R, incomingBatch.Header.S); err != nil {
		return err
	}
	fmt.Printf("%+v\n", incomingBatch.Header)

	// Validators recompute the entire batch using the same batch context
	// if they have all necessary prerequisites like having the l1 block processed
	// and the parent hash. This recomputed batch is then checked against the incoming batch.
	// If the sequencer has tampered with something the hash will not add up and validation will
	// produce an error.
	cb, err := val.batchProducer.ComputeBatch(&components.BatchExecutionContext{
		BlockPtr:     incomingBatch.Header.L1Proof,
		ParentPtr:    incomingBatch.Header.ParentHash,
		Transactions: incomingBatch.Transactions,
		AtTime:       incomingBatch.Header.Time,
		ChainConfig:  val.chainConfig,
		SequencerNo:  incomingBatch.Header.SequencerOrderNo,
	})
	if err != nil {
		return fmt.Errorf("failed recomputing batch %s. Cause: %w", incomingBatch.Hash(), err)
	}

	if !bytes.Equal(cb.Batch.Hash().Bytes(), incomingBatch.Hash().Bytes()) {
		// todo @stefan - generate a validator challenge here and return it
		return fmt.Errorf("batch is in invalid state. Incoming hash: %s  Computed hash: %s", incomingBatch.Hash().Hex(), cb.Batch.Hash().Hex())
	}

	if _, err := cb.Commit(true); err != nil {
		return fmt.Errorf("cannot commit stateDB for incoming valid batch %s. Cause: %w", incomingBatch.Hash(), err)
	}

	return val.batchRegistry.StoreBatch(incomingBatch, cb.Receipts)
}

func (val *obsValidator) ReceiveBlock(br *common.BlockAndReceipts, isLatest bool) (*components.BlockIngestionType, error) {
	ingestion, err := val.blockProcessor.Process(br, isLatest)
	if err != nil {
		return nil, err
	}

	rollup, err := val.rollupConsumer.ProcessL1Block(br)
	if err != nil && !errors.Is(err, components.ErrDuplicateRollup) {
		// todo - log err?
		val.logger.Error("Encountered error processing l1 block", log.ErrKey, err)
		return ingestion, nil
	}

	if rollup != nil {
		// read batch data from rollup, verify and store it
		if err = val.verifyRollup(rollup); err != nil {
			return nil, err
		}
	}

	return ingestion, nil
}

func (val *obsValidator) verifyRollup(rollup *core.Rollup) error {
	stopwatch := measure.NewStopwatch()
	defer val.logger.Info("Rollup processed", log.RollupHashKey, rollup.Hash(), log.DurationKey, stopwatch)

	for _, batch := range rollup.Batches {
		if err := val.ValidateAndStoreBatch(batch); err != nil {
			val.logger.Error("Attempted to store incorrect batch", log.BatchHashKey, batch.Hash(), log.ErrKey, err)
			return fmt.Errorf("failed validating and storing batch. Cause: %w", err)
		}
	}
	return nil
}

func (val *obsValidator) SubmitTransaction(transaction *common.L2Tx) error {
	val.logger.Trace(fmt.Sprintf("Transaction %s submitted to validator but there is nothing to do with it.", transaction.Hash().Hex()))
	return nil
}
