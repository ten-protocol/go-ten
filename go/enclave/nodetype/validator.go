package nodetype

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
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

func (val *obsValidator) ValidateAndStoreBatch(incomingBatch *core.Batch) error {
	receipts, err := val.batchRegistry.ValidateBatch(incomingBatch)
	if err != nil {
		return err
	}

	return val.batchRegistry.StoreBatch(incomingBatch, receipts)
}

func (val *obsValidator) SubmitTransaction(transaction *common.L2Tx) error {
	val.logger.Trace(fmt.Sprintf("Transaction %s submitted to validator but there is nothing to do with it.", transaction.Hash().Hex()))
	return nil
}
