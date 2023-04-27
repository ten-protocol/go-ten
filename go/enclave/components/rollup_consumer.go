package components

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type rollupConsumer struct {
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
	return &rollupConsumer{
		MgmtContractLib:       mgmtContractLib,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
		storage:               storage,
		verifier:              verifier,
	}
}

func (rc *rollupConsumer) ProcessL1Block(b *common.BlockAndReceipts) ([]*core.Rollup, error) {
	return rc.processRollups(b)
}

// getLatestRollupBeforeBlock - Given a block, returns the latest rollup in the canonical chain for that block (excluding those in the block itself).
func (rc *rollupConsumer) getLatestRollupBeforeBlock(block *common.L1Block) (*core.Rollup, error) {
	for {
		blockParentHash := block.ParentHash()
		latestRollup, err := rc.storage.FetchHeadRollupForBlock(&blockParentHash)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not fetch current L2 head rollup - %w", err)
		}
		if latestRollup != nil {
			return latestRollup, nil
		}

		// we scan backwards now to the prev block in the chain and we will lookup to see if that has an entry
		// todo (@stefan) - is this still required for safety, even though we're storing an entry for every L1 block?
		block, err = rc.storage.FetchBlock(block.ParentHash())
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				// No more blocks available (enclave does not read the L1 chain from genesis if it knows
				// when management contract was deployed, so we don't keep going to block zero, we just stop when the blocks run out)
				// We have now checked through the entire (relevant) history of the L1 and no rollups were found.
				return nil, db.ErrNoRollups
			}
			return nil, fmt.Errorf("could not fetch parent block - %w", err)
		}
	}
}

// extractRollups - returns a list of the rollups published in this block
func (rc *rollupConsumer) extractRollups(br *common.BlockAndReceipts, blockResolver db.BlockResolver) []*core.Rollup {
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
func (rc *rollupConsumer) processRollups(br *common.BlockAndReceipts) ([]*core.Rollup, error) {
	block := br.Block

	latestRollup, err := rc.getLatestRollupBeforeBlock(block)
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

	blockHash := block.Hash()
	for idx, rollup := range rollups {
		if err = rc.verifier.CheckSequencerSignature(rollup.Hash(), &rollup.Header.Agg, rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		if !rollup.IsGenesis() {
			previousRollup := latestRollup
			if idx != 0 {
				previousRollup = rollups[idx-1]
			}
			if err = rc.checkRollupsCorrectlyChained(rollup, previousRollup); err != nil {
				return nil, err
			}
		}

		for _, batch := range rollup.Batches {
			b, _ := rc.storage.FetchBatch(*batch.Hash())
			// only store the batch if not found in the db
			// todo (@matt) - this needs to be clarified if still exists after the FF refactor
			if b != nil {
				continue
			}
		}

		if err = rc.storage.StoreRollup(rollup); err != nil {
			return nil, fmt.Errorf("could not store rollup. Cause: %w", err)
		}
	}

	// we record the latest rollup published against this L1 block hash
	rollupHash := rollups[len(rollups)-1].Header.Hash()

	err = rc.storage.UpdateHeadRollup(&blockHash, &rollupHash)
	if err != nil {
		return nil, fmt.Errorf("unable to update head rollup - %w", err)
	}

	return rollups, nil
}

// Checks that the rollup:
//   - Has a number exactly 1 higher than the previous rollup
//   - Links to the previous rollup by hash
//   - Has a first batch whose parent is the head batch of the previous rollup
func (rc *rollupConsumer) checkRollupsCorrectlyChained(rollup *core.Rollup, previousRollup *core.Rollup) error {
	if big.NewInt(0).Sub(rollup.Header.Number, previousRollup.Header.Number).Cmp(big.NewInt(1)) != 0 {
		return fmt.Errorf("found gap in rollups between rollup %d and rollup %d",
			previousRollup.Header.Number, rollup.Header.Number)
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
