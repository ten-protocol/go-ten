package nodetype

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/txpool"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

type obsValidator struct {
	blockProcessor components.L1BlockProcessor
	batchExecutor  components.BatchExecutor
	batchRegistry  components.BatchRegistry
	rollupConsumer components.RollupConsumer

	chainConfig *params.ChainConfig

	sequencerID  gethcommon.Address
	storage      storage.Storage
	sigValidator *components.SignatureValidator
	mempool      *txpool.TxPool

	logger gethlog.Logger
}

func NewValidator(consumer components.L1BlockProcessor, batchExecutor components.BatchExecutor, registry components.BatchRegistry, rollupConsumer components.RollupConsumer, chainConfig *params.ChainConfig, sequencerID gethcommon.Address, storage storage.Storage, sigValidator *components.SignatureValidator, mempool *txpool.TxPool, logger gethlog.Logger) ObsValidator {
	startMempool(registry, mempool)

	return &obsValidator{
		blockProcessor: consumer,
		batchExecutor:  batchExecutor,
		batchRegistry:  registry,
		rollupConsumer: rollupConsumer,
		chainConfig:    chainConfig,
		sequencerID:    sequencerID,
		storage:        storage,
		sigValidator:   sigValidator,
		mempool:        mempool,
		logger:         logger,
	}
}

func (val *obsValidator) SubmitTransaction(tx *common.L2Tx) error {
	headBatch := val.batchRegistry.HeadBatchSeq()
	if headBatch == nil || headBatch.Uint64() <= common.L2GenesisSeqNo+1 {
		return fmt.Errorf("not initialised")
	}
	err := val.mempool.Validate(tx)
	if err != nil {
		val.logger.Info("Error validating transaction.", log.ErrKey, err, log.TxKey, tx.Hash())
	}
	return err
}

func (val *obsValidator) OnL1Fork(_ *common.ChainFork) error {
	// nothing to do
	return nil
}

func (val *obsValidator) VerifySequencerSignature(b *core.Batch) error {
	return val.sigValidator.CheckSequencerSignature(b.Hash(), b.Header.Signature)
}

func (val *obsValidator) ExecuteStoredBatches() error {
	headBatchSeq := val.batchRegistry.HeadBatchSeq()
	if headBatchSeq == nil {
		headBatchSeq = big.NewInt(int64(common.L2GenesisSeqNo))
	}
	batches, err := val.storage.FetchCanonicalUnexecutedBatches(headBatchSeq)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		return err
	}

	startMempool(val.batchRegistry, val.mempool)

	for _, batch := range batches {
		if batch.IsGenesis() {
			if err = val.handleGenesis(batch); err != nil {
				return err
			}
		}

		// check batch execution prerequisites
		canExecute, err := val.executionPrerequisites(batch)
		if err != nil {
			return fmt.Errorf("could not determine the execution prerequisites for batch %s. Cause: %w", batch.Hash(), err)
		}

		if canExecute {
			receipts, err := val.batchExecutor.ExecuteBatch(batch)
			if err != nil {
				return fmt.Errorf("could not execute batch %s. Cause: %w", batch.Hash(), err)
			}
			err = val.storage.StoreExecutedBatch(batch, receipts)
			if err != nil {
				return fmt.Errorf("could not store executed batch %s. Cause: %w", batch.Hash(), err)
			}
			err = val.mempool.Chain.IngestNewBlock(batch)
			if err != nil {
				return fmt.Errorf("failed to feed batch into the virtual eth chain- %w", err)
			}
			val.batchRegistry.OnBatchExecuted(batch, receipts)
		}
	}
	return nil
}

func (val *obsValidator) executionPrerequisites(batch *core.Batch) (bool, error) {
	// 1.l1 block exists
	block, err := val.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil && errors.Is(err, errutil.ErrNotFound) {
		val.logger.Info("Error fetching block", log.BlockHashKey, batch.Header.L1Proof, log.ErrKey, err)
		return false, err
	}

	// 2. parent was executed
	parentExecuted, err := val.storage.BatchWasExecuted(batch.Header.ParentHash)
	if err != nil {
		val.logger.Info("Error reading execution status of batch", log.BatchHashKey, batch.Header.ParentHash, log.ErrKey, err)
		return false, err
	}

	return block != nil && parentExecuted, nil
}

func (val *obsValidator) handleGenesis(batch *core.Batch) error {
	genBatch, _, err := val.batchExecutor.CreateGenesisState(batch.Header.L1Proof, batch.Header.Time, batch.Header.Coinbase, batch.Header.BaseFee)
	if err != nil {
		return err
	}

	if genBatch.Hash() != batch.Hash() {
		return fmt.Errorf("received invalid genesis batch")
	}

	err = val.storage.StoreExecutedBatch(genBatch, nil)
	if err != nil {
		return err
	}
	val.batchRegistry.OnBatchExecuted(batch, nil)
	return nil
}

func (val *obsValidator) OnL1Block(_ types.Block, _ *components.BlockIngestionType) error {
	return val.ExecuteStoredBatches()
}

func (val *obsValidator) Close() error {
	return val.mempool.Close()
}

func startMempool(registry components.BatchRegistry, mempool *txpool.TxPool) {
	// the mempool can only be started when there are a couple of blocks already processed
	headBatchSeq := registry.HeadBatchSeq()
	if !mempool.Running() && headBatchSeq != nil && headBatchSeq.Uint64() > common.L2GenesisSeqNo+1 {
		err := mempool.Start()
		if err != nil {
			panic(fmt.Errorf("could not start mempool: %w", err))
		}
	}
}
