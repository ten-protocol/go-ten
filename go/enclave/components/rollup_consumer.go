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

func (rc *rollupConsumerImpl) ExtractAndVerifyRollupData(rollupTx *common.L1TxData) (*common.ExtRollup, error) {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer verified rollup data", &core.RelaxedThresholds)
	// extract blob hashes, signatures and recreate rollup
	rollup, compositeHash, blobHashes, signatures, err := rc.extractRollupData(rollupTx)
	if err != nil {
		return nil, err
	}

	err = rc.verifySequencerSignature(rollup, *compositeHash, signatures)
	if err != nil {
		return nil, fmt.Errorf("invalid sequencer signature: %w", err)
	}

	err = rc.verifyBlobHashes(rollupTx, blobHashes)
	if err != nil {
		// critical error as the sequencer has signed this rollup
		return nil, fmt.Errorf("rollup hash verification failed: %w", errutil.ErrCriticalRollupProcessing)
	}

	return rollup, nil
}

// ProcessRollup - processes the rollup found in the block and stores it. The verification of the rollup data happens
// before calling this function.
func (rc *rollupConsumerImpl) ProcessRollup(ctx context.Context, rollup *common.ExtRollup) (*common.ExtRollupMetadata, error) {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed rollup", &core.RelaxedThresholds)

	l1CompressionBlock, err := rc.storage.FetchBlock(ctx, rollup.Header.CompressionL1Head)
	if err != nil {
		rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
		return nil, nil
	}
	canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(ctx, l1CompressionBlock.Number)
	if err != nil {
		return nil, err
	}
	if canonicalBlockByHeight.Hash() != l1CompressionBlock.Hash() {
		rc.logger.Warn("Skipping rollup because it was compressed on top of a non-canonical block", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
		return nil, nil
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

	serializedTree, err := rc.ExportAndVerifyCrossChainData(ctx, rollup.Header.FirstBatchSeqNo, rollup.Header.LastBatchSeqNo, rollup.Header.CrossChainRoot)
	if err != nil {
		rc.logger.Error("Failed exporting and verifying cross chain data", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
		return nil, err
	}

	rollupMetadata := common.ExtRollupMetadata{
		CrossChainTree: serializedTree,
	}

	return &rollupMetadata, nil
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

// extractRollupData - extracts the data required to verify and process the rollup transaction.
// 1. Extracts blobs and signatures from the transaction
// 2. Computes blob hashes using KZG commitments
// 3. Reconstructs the rollup from blob data
//
// Note: All errors are considered non-critical as they occur prior to signature verification
// and could be due to malformed or invalid input data. We don't want to prevent blocks from being processed if this is
// the case.
func (rc *rollupConsumerImpl) extractRollupData(rollupTx *common.L1TxData) (*common.ExtRollup, *gethcommon.Hash, []gethcommon.Hash, [][]byte, error) {
	blobs := make([]*kzg4844.Blob, 0)
	signatures := make([][]byte, 0)
	for _, blobWithSig := range rollupTx.BlobsWithSignature {
		blobs = append(blobs, blobWithSig.Blob)
		signatures = append(signatures, blobWithSig.Signature)
	}

	_, blobHashes, err := ethadapter.MakeSidecar(blobs, rc.MgmtContractLib.BlobHasher())
	if err != nil {
		// non-critical as signature not verified - could be bad data
		return nil, nil, nil, nil, fmt.Errorf("could not get blob hashes from blobs. Cause: %w", err)
	}

	rollup, err := ethadapter.ReconstructRollup(blobs)
	if err != nil {
		// non-critical as signature not verified - could be bad data
		return nil, nil, nil, nil, fmt.Errorf("could not recreate rollup from blobs. Cause: %w", err)
	}

	// TODO would there ever be more than one blob hash and signature?
	compositeHash := common.ComputeCompositeHash(rollup.Header, blobHashes[0])
	return rollup, &compositeHash, blobHashes, signatures, nil
}

// verifySequencerSignature - verifies the sequencer signature using a composite hash of the rollup header and blob hash
func (rc *rollupConsumerImpl) verifySequencerSignature(rollup *common.ExtRollup, compositeHash gethcommon.Hash, signatures [][]byte) error {
	if err := rc.sigValidator.CheckSequencerSignature(compositeHash, signatures[0]); err != nil {
		// non-critical as signature not verified
		return fmt.Errorf("rollup signature was invalid. Cause: %w", err)
	}
	rc.logger.Info("Extracted rollup from block with valid sequencer signature", log.RollupHashKey, rollup.Hash(), log.BlockHashKey, rollup.Header.CompressionL1Head.Hex())

	return nil
}

// verifyBlobHashes - verifies that all blob hashes referenced in a rollup transaction
// exist in the block's blob hash list. Since multiple rollups can be included in a single
// block, the rollup's blob hashes should be a subset of the block's total blob hashes.
//
// The function creates an efficient hash lookup map and verifies each rollup blob hash
// exists in the block's blob hash set.
func (rc *rollupConsumerImpl) verifyBlobHashes(rollupTx *common.L1TxData, blobHashes []gethcommon.Hash) error {
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
