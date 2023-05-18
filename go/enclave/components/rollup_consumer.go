package components

import (
	"errors"
	"fmt"
	"sort"

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

	// TransactionBlobCrypto- This contains the required properties to encrypt rollups.
	TransactionBlobCrypto crypto.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	storage  db.Storage
	verifier *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	transactionBlobCrypto crypto.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
	storage db.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:       mgmtContractLib,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
		storage:               storage,
		verifier:              verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessL1Block(b *common.BlockAndReceipts) ([]*core.Rollup, error) {
	return rc.processRollups(b)
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
			rollups = append(rollups, core.ToRollup(r, rc.TransactionBlobCrypto))
			rc.logger.Info(fmt.Sprintf("Extracted Rollup r_%d from block b_%d",
				common.ShortHash(r.Hash()),
				common.ShortHash(b.Hash()),
			))
		} else {
			rc.logger.Warn(fmt.Sprintf("Ignored rollup r_%d from block b_%d, because it was produced on a fork",
				common.ShortHash(r.Hash()),
				common.ShortHash(b.Hash()),
			))
		}
	}

	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	return rollups
}

// Validates and stores the rollup in a given block.
// todo (#718) - design a mechanism to detect a case where the rollups never contain any batches (despite batches arriving via P2P)
func (rc *rollupConsumerImpl) processRollups(br *common.BlockAndReceipts) ([]*core.Rollup, error) {
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
		return nil, nil
	}

	validRollups := make([]*core.Rollup, 0)

	blockHash := block.Hash()
	for _, rollup := range rollups {
		if err = rc.verifier.CheckSequencerSignature(rollup.Hash(), &rollup.Header.Agg, rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		if !rollup.IsGenesis() {
			if err = rc.checkRollupsCorrectlyChained(rollup, latestRollup); err != nil {
				rc.logger.Warn("Rollup was not correctly chained",
					"height", rollup.NumberU64(), "hash", rollup.Hash(), log.ErrKey, err)
				// we need to check the other rollups in this block still otherwise a correct one could get missed
				continue
			}
		}

		// todo (@matt) - store batches from the rollup (important during catch-up)

		// rollup has been verified as the next valid rollup. If there are more rollups to process then we build then on top of this one
		validRollups = append(validRollups, rollup)
		latestRollup = rollup
		if err = rc.storage.StoreRollup(rollup); err != nil {
			// todo (@matt) - this seems catastrophic, how do we recover the lost rollup in this case?
			return nil, fmt.Errorf("could not store rollup. Cause: %w", err)
		}
	}
	if len(validRollups) == 0 {
		// no valid rollups were found in this block
		return nil, nil
	}

	// we record the latest rollup published against this L1 block hash
	rollupHash := latestRollup.Header.Hash()

	err = rc.storage.UpdateHeadRollup(&blockHash, &rollupHash)
	if err != nil {
		// todo (@matt) - this also seems catastrophic, would result in bad state unable to ingest further rollups?
		return nil, fmt.Errorf("unable to update head rollup - %w", err)
	}

	return validRollups, nil
}

// Checks that the rollup:
//   - Has a number exactly 1 higher than the previous rollup
//   - Links to the previous rollup by hash
//   - Has a first batch whose parent is the head batch of the previous rollup
func (rc *rollupConsumerImpl) checkRollupsCorrectlyChained(rollup *core.Rollup, previousRollup *core.Rollup) error {
	if rollup.NumberU64()-previousRollup.NumberU64() > 1 {
		return fmt.Errorf("found gap in rollups between rollup %d and rollup %d",
			previousRollup.NumberU64(), rollup.NumberU64())
	}
	if rollup.NumberU64() <= previousRollup.NumberU64() {
		return fmt.Errorf("expected new rollup but rollup %d height was less than or equal to previous rollup %d",
			rollup.NumberU64(), previousRollup.NumberU64())
	}

	if rollup.Header.ParentHash != *previousRollup.Hash() {
		return fmt.Errorf("found gap in rollups. Rollup %d did not reference rollup %d by hash",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	if len(rollup.Batches) != 0 && previousRollup.Header.HeadBatchHash != rollup.Batches[0].Header.ParentHash {
		return fmt.Errorf("found gap in rollup batches. Batches in rollup %d did not chain to batches in rollup %d",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	return nil
}
