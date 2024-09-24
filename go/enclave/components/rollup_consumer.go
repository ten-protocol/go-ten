package components

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/ethadapter"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
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

// ProcessBlobsInBlock - FIXME
func (rc *rollupConsumerImpl) ProcessBlobsInBlock(ctx context.Context, b *common.BlockAndReceipts, blobs []*kzg4844.Blob) error {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed blobs", log.BlockHashKey, b.BlockHeader.Hash())

	rollups, err := rc.extractAndVerifyRollups(b, blobs)
	if err != nil {
		rc.logger.Error("Failed to extract rollups from block", log.BlockHashKey, b.BlockHeader.Hash(), log.ErrKey, err)
		return err
	}
	if len(rollups) == 0 {
		rc.logger.Info("No rollups found in block", log.BlockHashKey, b.BlockHeader.Hash(), log.ErrKey, err)
		return nil
	}

	rollups, err = rc.getSignedRollup(rollups)
	if err != nil {
		return err
	}

	if len(rollups) > 1 {
		// todo - we need to sort this out
		rc.logger.Warn(fmt.Sprintf("Multiple rollups %d in block %s", len(rollups), b.BlockHeader.Hash()))
	}

	for _, rollup := range rollups {
		l1CompressionBlock, err := rc.storage.FetchBlock(ctx, rollup.Header.CompressionL1Head)
		if err != nil {
			rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(ctx, l1CompressionBlock.Number)
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
func (rc *rollupConsumerImpl) extractAndVerifyRollups(br *common.BlockAndReceipts, blobs []*kzg4844.Blob) ([]*common.ExtRollup, error) {
	rollups := make([]*common.ExtRollup, 0)
	b := br.BlockHeader

	for i, tx := range *br.RelevantTransactions() {
		// go through all rollup transactions
		t := rc.MgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		rollupHashes, ok := t.(*ethadapter.L1RollupHashes)
		if !ok {
			continue
		}

		var blobHashes []gethcommon.Hash
		var err error
		if _, blobHashes, err = ethadapter.MakeSidecar(blobs); err != nil {
			return nil, fmt.Errorf("could not create blob sidecar and blob hashes. Cause: %w", err)
		}

		if err := verifyBlobHashes(rollupHashes, blobHashes); err != nil {
			rc.logger.Info(fmt.Sprintf("blob hashes in rollup at index %d do not match the rollup blob hashes", i))
			continue
		}

		r, err := ethadapter.ReconstructRollup(blobs)
		if err != nil {
			return nil, fmt.Errorf("could not recreate rollup from blobs. Cause: %w", err)
		}

		rollups = append(rollups, r)
		rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.BlockHashKey, b.Hash())
	}

	return rollups, nil
}

func verifyBlobHashes(rollupHashes *ethadapter.L1RollupHashes, blobHashes []gethcommon.Hash) error {
	if len(rollupHashes.BlobHashes) != len(blobHashes) {
		return fmt.Errorf("hash count mismatch: rollupHashes (%d) and blobHashes (%d)", len(rollupHashes.BlobHashes), len(blobHashes))
	}

	for i, hash := range rollupHashes.BlobHashes {
		if hash != blobHashes[i] {
			return fmt.Errorf("hash mismatch at index %d: rollupHash (%s) != blobHash (%s)", i, hash.Hex(), blobHashes[i].Hex())
		}
	}
	return nil
}
