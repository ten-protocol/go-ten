package components

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ten-protocol/go-ten/go/common/measure"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	rollupCompression *RollupCompression
	batchRegistry     BatchRegistry

	logger gethlog.Logger

	storage      storage.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
	rollupCompression *RollupCompression,
	storage storage.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:   mgmtContractLib,
		batchRegistry:     batchRegistry,
		rollupCompression: rollupCompression,
		logger:            logger,
		storage:           storage,
		sigValidator:      verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessRollupsInBlock(b *common.BlockAndReceipts) error {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed block", log.BlockHashKey, b.Block.Hash())

	rollups := rc.extractRollups(b)
	if len(rollups) == 0 {
		return nil
	}

	rollups, err := rc.getSignedRollup(rollups)
	if err != nil {
		return err
	}

	for _, rollup := range rollups {
		l1CompressionBlock, err := rc.storage.FetchBlock(rollup.Header.CompressionL1Head)
		if err != nil {
			rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(l1CompressionBlock.Number())
		if err != nil {
			return err
		}
		if canonicalBlockByHeight.Hash() != l1CompressionBlock.Hash() {
			rc.logger.Warn("Skipping rollup because it was compressed on top of a non-canonical rollup", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		// read batch data from rollup, verify and store it
		internalHeader, err := rc.rollupCompression.ProcessExtRollup(rollup)
		if err != nil {
			rc.logger.Error("Failed processing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			// todo - issue challenge as a validator
			return err
		}
		if err := rc.storage.StoreRollup(rollup, internalHeader); err != nil {
			rc.logger.Error("Failed storing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			return err
		}
	}

	return nil
}

func (rc *rollupConsumerImpl) getSignedRollup(rollups []*common.ExtRollup) ([]*common.ExtRollup, error) {
	signedRollup := make([]*common.ExtRollup, 0)

	// loop through the rollups, find the one that is signed, verify the signature, make sure it's the only one
	for _, rollup := range rollups {
		if err := rc.sigValidator.CheckSequencerSignature(rollup.Hash(), rollup.Header.Signature); err != nil {
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
