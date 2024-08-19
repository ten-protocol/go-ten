package components

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ten-protocol/go-ten/go/common/measure"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	rollupCompression *RollupCompression
	batchRegistry     BatchRegistry
	blobResolver      BlobResolver
	logger            gethlog.Logger

	storage      storage.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
	blobResolver BlobResolver,
	rollupCompression *RollupCompression,
	storage storage.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:   mgmtContractLib,
		batchRegistry:     batchRegistry,
		blobResolver:      blobResolver,
		rollupCompression: rollupCompression,
		logger:            logger,
		storage:           storage,
		sigValidator:      verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessRollupsInBlock(ctx context.Context, b *common.BlockAndReceipts) error {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed block", log.BlockHashKey, b.Block.Hash())

	rollups := rc.extractRollups(ctx, b)
	if len(rollups) == 0 {
		return nil
	}

	rollups, err := rc.getSignedRollup(rollups)
	if err != nil {
		return err
	}

	if len(rollups) > 1 {
		// todo - we need to sort this out
		rc.logger.Warn(fmt.Sprintf("Multiple rollups %d in block %s", len(rollups), b.Block.Hash()))
	}

	for _, rollup := range rollups {
		l1CompressionBlock, err := rc.storage.FetchBlock(ctx, rollup.Header.CompressionL1Head)
		if err != nil {
			rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(ctx, l1CompressionBlock.Number())
		if err != nil {
			return err
		}
		if canonicalBlockByHeight.Hash() != l1CompressionBlock.Hash() {
			rc.logger.Warn("Skipping rollup because it was compressed on top of a non-canonical rollup", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		// read batch data from rollup, verify and store it
		internalHeader, err := rc.rollupCompression.ProcessExtRollup(ctx, rollup)
		if err != nil {
			rc.logger.Error("Failed processing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			// todo - issue challenge as a validator
			return err
		}
		if err := rc.storage.StoreRollup(ctx, rollup, internalHeader); err != nil {
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
func (rc *rollupConsumerImpl) extractRollups(ctx context.Context, br *common.BlockAndReceipts) []*common.ExtRollup {
	rollups := make([]*common.ExtRollup, 0)
	b := br.Block

	for _, tx := range *br.SuccessfulTransactions() {
		decodedTx := rc.MgmtContractLib.DecodeTx(tx)
		if decodedTx == nil {
			continue
		}

		rollupHashes, ok := decodedTx.(*ethadapter.L1RollupHashes)
		if !ok {
			continue
		}

		blobs, err := rc.blobResolver.FetchBlobs(ctx, br.Block.Header(), rollupHashes.BlobHashes)
		if err != nil {
			rc.logger.Crit("could not fetch blobs consumer", log.ErrKey, err)
			return nil
		}
		r, err := reconstructRollup(blobs)
		if err != nil {
			rc.logger.Crit("could not recreate rollup from blobs.", log.ErrKey, err)
		}

		rollups = append(rollups, r)
		rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.BlockHashKey, b.Hash())
	}
	return rollups
}

// Function to reconstruct rollup from blobs
func reconstructRollup(blobs []*kzg4844.Blob) (*common.ExtRollup, error) {
	data, err := ethadapter.DecodeBlobs(blobs)
	if err != nil {
		fmt.Println("Error decoding rollup blob:", err)
	}
	var rollup common.ExtRollup
	if err := rlp.DecodeBytes(data, &rollup); err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}
	return &rollup, nil
}
