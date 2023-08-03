package nodetype

import (
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

type obsValidator struct {
	blockProcessor components.L1BlockProcessor
	batchProducer  components.BatchProducer
	batchRegistry  components.BatchRegistry
	rollupConsumer components.RollupConsumer

	chainConfig *params.ChainConfig

	sequencerID  gethcommon.Address
	storage      storage.Storage
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
	storage storage.Storage,
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

func (val *obsValidator) SubmitTransaction(transaction *common.L2Tx) error {
	val.logger.Trace(fmt.Sprintf("Transaction %s submitted to validator but there is nothing to do with it.", transaction.Hash().Hex()))
	return nil
}

func (val *obsValidator) OnL1Fork(_ *common.ChainFork) error {
	// nothing to do
	return nil
}

func (val *obsValidator) VerifySequencerSignature(*core.Batch) error {
	// todo
	return nil
}

func (val *obsValidator) ExecuteBatches() error {
	batches, err := val.storage.FetchUnexecutedBatches()
	if err != nil {
		return err
	}

	for _, batch := range batches {

		if batch.IsGenesis() {
			genBatch, _, err := val.batchProducer.CreateGenesisState(batch.Header.L1Proof, batch.Header.Time)
			if err != nil {
				return err
			}

			if genBatch.Hash() != batch.Hash() {
				return fmt.Errorf("received invalid genesis batch")
			}

			return val.storage.StoreExecutedBatch(genBatch, nil)
		}

		// check prerequisites
		// l1 block exists
		block, err := val.storage.FetchBlock(batch.Header.L1Proof)
		if err != nil && errors.Is(err, errutil.ErrNotFound) {
			val.logger.Info("Error fetching block", log.BlockHashKey, batch.Header.L1Proof, log.ErrKey, err)
			continue
		}

		// parent was executed
		parentExecuted, err := val.storage.BatchWasExecuted(batch.Header.ParentHash)
		if err != nil {
			val.logger.Info("Error reading execution status of batch", log.BatchHashKey, batch.Header.ParentHash, log.ErrKey, err)
			continue
		}

		if block != nil && parentExecuted {
			receipts, err := val.batchRegistry.ExecuteBatch(batch)
			if err != nil {
				// todo
				return err
			}
			err = val.storage.StoreExecutedBatch(batch, receipts)
			if err != nil {
				// todo
				return err
			}
		}
	}

	return nil
}
