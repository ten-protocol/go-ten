package components

import (
	"errors"
	"fmt"
	"sort"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	storage      db.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	dataEncryptionService crypto.DataEncryptionService,
	dataCompressionService compression.DataCompressionService,
	obscuroChainID int64,
	ethereumChainID int64,
	storage db.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:        mgmtContractLib,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		ObscuroChainID:         obscuroChainID,
		EthereumChainID:        ethereumChainID,
		logger:                 logger,
		storage:                storage,
		sigValidator:           verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessL1Block(b *common.BlockAndReceipts) (*core.Rollup, error) {
	return rc.processRollup(b)
}

// extractRollups - returns a list of the rollups published in this block
func (rc *rollupConsumerImpl) extractRollups(br *common.BlockAndReceipts, blockResolver db.BlockResolver) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
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

		// Ignore rollups created with proofs from different L1 blocks
		// In case of L1 reorgs, rollups may end published on a fork
		if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
			rollup, err := core.ToRollup(r, rc.dataEncryptionService, rc.dataCompressionService)
			if err != nil {
				// todo - this should not fail, but generate the proof
				rc.logger.Crit("Failed to transform rollup", log.ErrKey, err)
			}
			rollups = append(rollups, rollup)
			rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.RollupHeightKey, r.Header.Number, log.BlockHashKey, b.Hash())
		} else {
			rc.logger.Warn("Ignored rollup from block, because it was produced on a fork", log.RollupHashKey, r.Hash(), log.RollupHeightKey, r.Header.Number, log.BlockHashKey, b.Hash())
		}
	}

	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	return rollups
}

// Validates and stores the rollup in a given block. Returns nil, nil when no rollup was found.
// todo (#718) - design a mechanism to detect a case where the rollup doesn't contain any batches (despite batches arriving via P2P)
func (rc *rollupConsumerImpl) processRollup(br *common.BlockAndReceipts) (*core.Rollup, error) {
	block := br.Block

	latestRollup, err := getLatestRollupBeforeBlock(block, rc.storage, rc.logger)
	if err != nil && !errors.Is(err, db.ErrNoRollups) {
		return nil, fmt.Errorf("unexpected error retrieving latest rollup for block %s. Cause: %w", block.Hash(), err)
	}

	rollups := rc.extractRollups(br, rc.storage)

	// If this is the first rollup we've ever received, we check that it's the genesis rollup.
	if latestRollup == nil && len(rollups) != 0 && !rollups[0].IsGenesis() {
		return nil, fmt.Errorf("received rollup with number %d but no genesis rollup is stored", rollups[0].Number())
	}

	if len(rollups) == 0 {
		return nil, nil //nolint:nilnil
	}

	var signedRollup *core.Rollup
	blockHash := block.Hash()

	// loop through the rollups, find the one that is signed, verify the signature, make sure it's the only one
	for _, rollup := range rollups {
		if err = rc.sigValidator.CheckSequencerSignature(rollup.Hash(), rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}
		if signedRollup != nil {
			// todo (@matt) - make sure this can't be used to DOS the network
			// we should never receive multiple signed rollups in a single block, the host should only ever publish one
			return nil, fmt.Errorf("received multiple signed rollups in single block %s", blockHash)
		}
		signedRollup = rollup
	}
	if signedRollup == nil {
		return nil, nil //nolint:nilnil
	}

	if err = rc.checkRollupsCorrectlyChained(signedRollup, latestRollup); err != nil {
		return nil, fmt.Errorf("rollup was not correctly chained. height=%d hash=%s Cause: %w",
			signedRollup.NumberU64(), signedRollup.Hash(), err)
	}

	// todo (@matt) - store batches from the rollup (important during catch-up)

	if err = rc.storage.StoreRollup(signedRollup); err != nil {
		// todo (@matt) - this seems catastrophic, how do we recover the lost rollup in this case?
		return nil, fmt.Errorf("could not store rollup. Cause: %w", err)
	}

	// we record the latest rollup published against this L1 block hash
	rollupHash := signedRollup.Hash()
	err = rc.storage.UpdateHeadRollup(&blockHash, &rollupHash)
	if err != nil {
		// todo (@matt) - this also seems catastrophic, would result in bad state unable to ingest further rollups?
		return nil, fmt.Errorf("unable to update head rollup - %w", err)
	}

	return signedRollup, nil
}

// Checks that the rollup:
//   - Has a number exactly 1 higher than the previous rollup
//   - Links to the previous rollup by hash
//   - Has a first batch whose parent is the head batch of the previous rollup
func (rc *rollupConsumerImpl) checkRollupsCorrectlyChained(rollup *core.Rollup, previousRollup *core.Rollup) error {
	if previousRollup == nil {
		// genesis rollup has no previous rollup to check
		return nil
	}

	if rollup.Hash() == previousRollup.Hash() {
		return ErrDuplicateRollup
	}

	if rollup.NumberU64()-previousRollup.NumberU64() > 1 {
		return fmt.Errorf("found gap in rollups between rollup %d and rollup %d",
			previousRollup.NumberU64(), rollup.NumberU64())
	}

	// In case we have published two rollups for the same height
	// This can happen when the first one takes too long to mine
	if rollup.NumberU64() == previousRollup.NumberU64() {
		if len(previousRollup.Batches) > len(rollup.Batches) {
			return fmt.Errorf("received duplicate rollup at height %d with less batches than previous rollup", rollup.NumberU64())
		}

		for idx, batch := range previousRollup.Batches {
			if rollup.Batches[idx].Hash() != batch.Hash() {
				return fmt.Errorf("duplicate rollup at height %d has different batches at position %d", rollup.NumberU64(), idx)
			}
		}

		return ErrDuplicateRollup
	}

	if rollup.NumberU64() <= previousRollup.NumberU64() {
		return fmt.Errorf("expected new rollup but rollup %d height was less than or equal to previous rollup %d",
			rollup.NumberU64(), previousRollup.NumberU64())
	}

	if rollup.Header.ParentHash != previousRollup.Hash() {
		return fmt.Errorf("found gap in rollups. Rollup %d did not reference rollup %d by hash",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	if len(rollup.Batches) != 0 && previousRollup.Header.HeadBatchHash != rollup.Batches[0].Header.ParentHash {
		return fmt.Errorf("found gap in rollup batches. Batches in rollup %d did not chain to batches in rollup %d",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	return nil
}
