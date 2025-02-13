package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/common/merkle"
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
	sigValidator SequencerSignatureVerifier
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
	rollupCompression *RollupCompression,
	storage storage.Storage,
	logger gethlog.Logger,
	verifier SequencerSignatureVerifier,
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

func (rc *rollupConsumerImpl) ProcessRollupData(ctx context.Context, processed *common.ProcessedL1Data) ([]common.ExtRollupMetadata, error) {
	rollupTxs := processed.GetEvents(common.RollupTx)
	rollups := make([]*common.ExtRollup, 0, len(rollupTxs))
	txsSeen := make(map[gethcommon.Hash]bool)

	for _, rtx := range rollupTxs {
		// verify signature and decode rollups
		r, hashes, err := rc.verifySequencerSignature(rtx)
		if err != nil {
			return nil, fmt.Errorf("invalid sequencer signature: %w", err)
		}

		// prevent the case where someone pushes a blob to the same slot. multiple rollups can be found in a block,
		// but they must come from unique transactions
		if txsSeen[rtx.Transaction.Hash()] {
			return nil, fmt.Errorf("multiple rollups from same transaction: %s. Err: %w", rtx.Transaction.Hash(), errutil.ErrCriticalRollupProcessing)
		}

		err = rc.verifyBlobHashes(rtx, hashes)
		if err != nil {
			// critical error as the sequencer has signed this rollup
			return nil, fmt.Errorf("rollup hash verification failed: %w", errutil.ErrCriticalRollupProcessing)
		}
		rollups = append(rollups, r)
		txsSeen[rtx.Transaction.Hash()] = true
	}

	if len(rollups) == 0 {
		rc.logger.Warn("No rollups found in block when rollupTxs present", log.BlockHashKey, processed.BlockHeader.Hash())
		return nil, nil
	}

	if len(rollups) > 1 {
		// this is allowed as long as they come from unique transactions
		rc.logger.Trace(fmt.Sprintf("Multiple rollups %d in block %s", len(rollups), processed.BlockHeader.Hash()))
	}

	// process rollup through compression service
	metadata, err := rc.ProcessRollups(ctx, rollups)
	if err != nil {
		return nil, fmt.Errorf("failed to process rollup: %w", errutil.ErrCriticalRollupProcessing)
	}

	return metadata, nil
}

// ProcessRollups - processes the rollups found in the block, verifies the rollups and stores them
func (rc *rollupConsumerImpl) ProcessRollups(ctx context.Context, rollups []*common.ExtRollup) ([]common.ExtRollupMetadata, error) {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed rollups", &core.RelaxedThresholds)

	rollupMetadata := make([]common.ExtRollupMetadata, len(rollups))
	for idx, rollup := range rollups {
		l1CompressionBlock, err := rc.storage.FetchBlock(ctx, rollup.Header.CompressionL1Head)
		if err != nil {
			rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(ctx, l1CompressionBlock.Number)
		if err != nil {
			return nil, err
		}
		if canonicalBlockByHeight.Hash() != l1CompressionBlock.Hash() {
			rc.logger.Warn("Skipping rollup because it was compressed on top of a non-canonical block", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		// read batch data from rollup, verify and store it
		internalHeader, err := rc.rollupCompression.ProcessExtRollup(ctx, rollup)
		if err != nil {
			rc.logger.Error("Failed processing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			// todo - issue challenge as a validator
			return nil, err
		}
		if err := rc.storage.StoreRollup(ctx, rollup, internalHeader); err != nil {
			rc.logger.Error("Failed storing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			return nil, err
		}

		serializedTree, err := rc.ExportAndVerifyCrossChainData(ctx, internalHeader.FirstBatchSequence.Uint64(), rollup.Header.LastBatchSeqNo, rollup.Header.CrossChainRoot)
		if err != nil {
			rc.logger.Error("Failed exporting and verifying cross chain data", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			return nil, err
		}

		rollupMetadata[idx] = common.ExtRollupMetadata{
			CrossChainTree: serializedTree,
		}
	}

	if len(rollupMetadata) < len(rollups) {
		return nil, fmt.Errorf("missing metadata for some rollups")
	}

	return rollupMetadata, nil
}

func (rc *rollupConsumerImpl) ExportAndVerifyCrossChainData(ctx context.Context, fromSeqNo uint64, toSeqNo uint64, publishedCrossChainRoot gethcommon.Hash) (common.SerializedCrossChainTree, error) {
	batches, err := rc.storage.FetchCanonicalBatchesBetween(ctx, fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	localCrossChainRoot, serializedTree, err := merkle.ComputeCrossChainRootFromBatches(batches)
	if err != nil {
		return nil, err
	}

	if localCrossChainRoot != publishedCrossChainRoot {
		return nil, errutil.ErrCrossChainRootMismatch
	}

	return serializedTree, nil
}

func (rc *rollupConsumerImpl) verifySequencerSignature(rollupTx *common.L1TxData) (*common.ExtRollup, []gethcommon.Hash, error) {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed rollup sequencer", &core.RelaxedThresholds)
	blobs := make([]*kzg4844.Blob, 0)
	signatures := make([][]byte, 0)
	for _, blobWithSig := range rollupTx.BlobsWithSignature {
		blobs = append(blobs, blobWithSig.Blob)
		signatures = append(signatures, blobWithSig.Signature)
	}

	_, blobHashes, err := ethadapter.MakeSidecar(blobs, rc.MgmtContractLib.BlobHasher())
	if err != nil {
		// non-critical as signature not verified - could be bad data
		return nil, nil, fmt.Errorf("could not create blob sidecar and blob hashes. Cause: %w", err)
	}

	rollup, err := ethadapter.ReconstructRollup(blobs)
	if err != nil {
		// non-critical as signature not verified - could be bad data
		return nil, nil, fmt.Errorf("could not recreate rollup from blobs. Cause: %w", err)
	}

	// TODO would there ever be more than one blob hash and signature?
	compositeHash := common.ComputeCompositeHash(rollup.Header, blobHashes[0])
	if err := rc.sigValidator.CheckSequencerSignature(compositeHash, signatures[0]); err != nil {
		// non-critical as signature not verified
		return nil, nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
	}
	rc.logger.Info("Extracted rollup from block with valid sequencer signature", log.RollupHashKey, rollup.Hash(), log.BlockHashKey, rollup.Header.CompressionL1Head.Hex())

	return rollup, blobHashes, nil
}

// verifyBlobHashes -
// there may be many rollups in one block so the blobHashes array, so it is possible that the rollupHashes array is a
// subset of the blobHashes array
func (rc *rollupConsumerImpl) verifyBlobHashes(rollupTx *common.L1TxData, blobHashes []gethcommon.Hash) error {
	// more efficient lookup
	blobHashSet := make(map[gethcommon.Hash]struct{}, len(blobHashes))
	for _, h := range blobHashes {
		blobHashSet[h] = struct{}{}
	}

	t, err := rc.MgmtContractLib.DecodeTx(rollupTx.Transaction)
	if err != nil {
		return fmt.Errorf("could not decode tx. Cause: %s", err)
	}
	if t == nil {
		return fmt.Errorf("decoded transaction should not be nil at this point")
	}

	rollupHashes, ok := t.(*common.L1RollupHashes)
	if !ok {
		return fmt.Errorf("decoded transaction should contain blob hashes")
	}

	for i, rollupHash := range rollupHashes.BlobHashes {
		if _, exists := blobHashSet[rollupHash]; !exists {
			return fmt.Errorf(
				"rollupHash at index %d (%s) not found in blobHashes",
				i,
				rollupHash.Hex(),
			)
		}
	}
	return nil
}
