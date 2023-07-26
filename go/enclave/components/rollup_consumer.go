package components

import (
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/obscuronet/go-obscuro/go/common/measure"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	batchRegistry          BatchRegistry

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	storage      storage.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	obscuroChainID int64,
	ethereumChainID int64,
	storage storage.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:        mgmtContractLib,
		batchRegistry:          batchRegistry,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		ObscuroChainID:         obscuroChainID,
		EthereumChainID:        ethereumChainID,
		logger:                 logger,
		storage:                storage,
		sigValidator:           verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessRollupsInBlock(b *common.BlockAndReceipts) error {
	stopwatch := measure.NewStopwatch()
	defer rc.logger.Info("Rollup consumer processed block", log.BlockHashKey, b.Block.Hash(), log.DurationKey, stopwatch)

	rollups := rc.extractRollups(b)
	if len(rollups) == 0 {
		return nil
	}

	rollups, err := rc.getSignedRollup(rollups)
	if err != nil {
		return err
	}
	if len(rollups) > 0 {
		for _, rollup := range rollups {
			// read batch data from rollup, verify and store it
			if err := rc.ProcessRollup(rollup); err != nil {
				rc.logger.Error("Failed processing rollup", log.ErrKey, err)
				return err
			}
		}
	}

	return nil
}

func (rc *rollupConsumerImpl) getSignedRollup(rollups []*common.ExtRollup) ([]*common.ExtRollup, error) {
	signedRollup := make([]*common.ExtRollup, 0)

	// loop through the rollups, find the one that is signed, verify the signature, make sure it's the only one
	for _, rollup := range rollups {
		if err := rc.sigValidator.CheckSequencerSignature(rollup.Hash(), rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		signedRollup = append(signedRollup, rollup)
	}
	return signedRollup, nil
}

// todo - when processing the rollup, instead of looking up batches one by one, compare the last sequence number from the db with the ones in the rollup
// extractRollups - returns a list of the rollups published in this block
func (rc *rollupConsumerImpl) extractRollups(br *common.BlockAndReceipts) []*common.ExtRollup {
	rollups := make([]*common.ExtRollup, 0)
	b := br.Block

	for _, tx := range *br.SuccessfulTransactions() {
		// go through all rollup transactions
		t := rc.MgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		rolTx, ok := t.(*ethadapter.L1RollupTx)
		if !ok {
			continue
		}

		r, err := common.DecodeRollup(rolTx.Rollup)
		if err != nil {
			rc.logger.Crit("could not decode rollup.", log.ErrKey, err)
			return nil
		}

		rollups = append(rollups, r)
		rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.BlockHashKey, b.Hash())
	}

	return rollups
}

func (rc *rollupConsumerImpl) ProcessRollup(rollup *common.ExtRollup) error {
	// todo logic to decompress the rollups on the fly
	r, err := core.ToRollup(rollup, rc.dataEncryptionService, rc.dataCompressionService)
	if err != nil {
		return err
	}

	for _, batch := range r.Batches {
		rc.logger.Info("Processing batch from rollup", log.BatchHashKey, batch.Hash(), "seqNo", batch.SeqNo())
		b, batchFoundErr := rc.batchRegistry.GetBatch(batch.Hash())
		// Process and store a batch only if it wasn't already processed via p2p.
		if batchFoundErr != nil && !errors.Is(batchFoundErr, errutil.ErrNotFound) {
			return batchFoundErr
		}
		// if the batch is already stored, skip
		if b != nil {
			continue
		}
		receipts, err := rc.batchRegistry.ValidateBatch(batch)
		if errors.Is(err, errutil.ErrBlockForBatchNotFound) {
			rc.logger.Warn("Unable to validate batch due to it being on a different chain.", log.BatchHashKey, batch.Hash())
			continue
		}
		if err != nil {
			rc.logger.Error("Failed validating batch", log.BatchHashKey, batch.Hash(), log.ErrKey, err)
			return fmt.Errorf("failed validating and storing batch. Cause: %w", err)
		}

		err = rc.batchRegistry.StoreBatch(batch, receipts)
		if err != nil {
			return err
		}
	}
	return nil
}
