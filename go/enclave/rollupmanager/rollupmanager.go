package rollupmanager

// TODO: once the cross chain messages based bridge is implemented remove this completely.

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	crypto2 "github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/l2chain"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

// RollupManager encapsulates the logic of decoding rollup transactions submitted to the L1 and resolving them
// to rollups that the enclave can process.
type RollupManager struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	TransactionBlobCrypto crypto2.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	l2chain *l2chain.ObscuroChain
	storage db.Storage
}

func New(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	transactionBlobCrypto crypto2.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
	storage db.Storage,
	l2chain *l2chain.ObscuroChain,
	logger gethlog.Logger,
) *RollupManager {
	return &RollupManager{
		MgmtContractLib:       mgmtContractLib,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
		l2chain:               l2chain,
		storage:               storage,
	}
}

func (re *RollupManager) fetchLatestRolluп() (*core.Rollup, error) {
	b, err := re.storage.FetchHeadBlock()
	if b != nil {
		return re.getLatestRollupBeforeBlock(b)
	} else {
		return nil, err
	}
}

func NextRollup(rollup *core.Rollup, batches []*core.Batch) *core.Rollup {
	headBatch := batches[len(batches)-1]

	rh := headBatch.Header.ToRollupHeader()

	if rollup != nil {
		rh.ParentHash = rollup.Header.Hash()
	} else { // genesis
		rh.ParentHash = gethcommon.Hash{}
	}

	rh.CrossChainMessages = make([]MessageBus.StructsCrossChainMessage, 0)
	for _, b := range batches {
		rh.CrossChainMessages = append(rh.CrossChainMessages, b.Header.CrossChainMessages...)
	}

	rollupHeight := big.NewInt(0)
	if rollup != nil {
		rollupHeight = rollup.Header.Number
		rollupHeight.Add(rollupHeight, gethcommon.Big1)
	}

	rh.Number = rollupHeight
	rh.HeadBatchHash = headBatch.Header.Hash()

	return &core.Rollup{
		Header:  rh,
		Batches: batches,
	}
}

func (re *RollupManager) CreateRollup() (*core.Rollup, error) {
	rollup, err := re.fetchLatestRolluп()
	if err != nil && !errors.Is(err, db.ErrNoRollups) {
		return nil, err
	}

	hash := gethcommon.Hash{}
	if rollup != nil {
		hash = rollup.Header.HeadBatchHash
	}

	batches, err := re.l2chain.BatchesAfter(hash)
	if err != nil {
		return nil, err
	}

	if len(batches) == 0 {
		return nil, fmt.Errorf("no batches for rollup")
	}

	if batches[len(batches)-1].Header.Hash() == hash {
		return nil, fmt.Errorf("current head batch matches the rollup head bash")
	}

	newRollup := NextRollup(rollup, batches)

	re.l2chain.SignRollup(newRollup)
	return newRollup, nil
}

func (re *RollupManager) ProcessL1Block(b *types.Block) ([]*core.Rollup, error) {
	return re.processRollups(b)
}

// ExtractRollups - returns a list of the rollups published in this block
func (re *RollupManager) ExtractRollups(b *types.Block, blockResolver db.BlockResolver) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := re.MgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if rolTx, ok := t.(*ethadapter.L1RollupTx); ok {
			r, err := common.DecodeRollup(rolTx.Rollup)
			if err != nil {
				re.logger.Crit("could not decode rollup.", log.ErrKey, err)
				return nil
			}

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, core.ToRollup(r, re.TransactionBlobCrypto))
				re.logger.Trace(fmt.Sprintf("Extracted Rollup r_%d from block b_%d",
					common.ShortHash(r.Hash()),
					common.ShortHash(b.Hash()),
				))
			}
		}
	}

	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	return rollups
}

// MOVED FROM L2 CHAIN

// Validates and stores the rollup in a given block.
// TODO - #718 - Design a mechanism to detect a case where the rollups never contain any batches (despite batches arriving via P2P).
func (re *RollupManager) processRollups(block *common.L1Block) ([]*core.Rollup, error) {
	latestRollup, err := re.getLatestRollupBeforeBlock(block)
	if err != nil && !errors.Is(err, db.ErrNoRollups) {
		return nil, fmt.Errorf("unexpected error retrieving latest rollup for block %s. Cause: %w", block.Hash(), err)
	}

	rollups := re.ExtractRollups(block, re.storage)

	// If this is the first rollup we've ever received, we check that it's the genesis rollup.
	if latestRollup == nil && len(rollups) != 0 && !rollups[0].IsGenesis() {
		return nil, fmt.Errorf("received rollup with number %d but no genesis rollup is stored", rollups[0].Number())
	}

	if len(rollups) == 0 {
		return nil, nil
	}

	blockHash := block.Hash()
	for idx, rollup := range rollups {
		if err = re.l2chain.CheckSequencerSignature(rollup.Hash(), &rollup.Header.Agg, rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		if !rollup.IsGenesis() {
			previousRollup := latestRollup
			if idx != 0 {
				previousRollup = rollups[idx-1]
			}
			if err = re.checkRollupsCorrectlyChained(rollup, previousRollup); err != nil {
				return nil, err
			}
		}

		for _, batch := range rollup.Batches {
			if err = re.l2chain.CheckAndStoreBatch(batch); err != nil {
				return nil, fmt.Errorf("could not store batch. Cause: %w", err)
			}
		}

		if err = re.storage.StoreRollup(rollup); err != nil {
			return nil, fmt.Errorf("could not store rollup. Cause: %w", err)
		}
	}

	// we record the latest rollup published against this L1 block hash
	rollupHash := rollups[len(rollups)-1].Header.Hash()

	err = re.storage.UpdateHeadRollup(&blockHash, &rollupHash)
	if err != nil {
		return nil, fmt.Errorf("unable to update head rollup - %w", err)
	}

	return rollups, nil
}

// Given a block, returns the latest rollup in the canonical chain for that block (excluding those in the block itself).
func (re *RollupManager) getLatestRollupBeforeBlock(block *common.L1Block) (*core.Rollup, error) {
	for {
		blockParentHash := block.ParentHash()
		latestRollup, err := re.storage.FetchHeadRollupForBlock(&blockParentHash)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not fetch current L2 head rollup - %w", err)
		}
		if latestRollup != nil {
			return latestRollup, nil
		}

		// we scan backwards now to the prev block in the chain and we will lookup to see if that has an entry
		// todo: is this still required for safety, even though we're storing an entry for every L1 block?
		block, err = re.storage.FetchBlock(block.ParentHash())
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

// Checks that the rollup:
//   - Has a number exactly 1 higher than the previous rollup
//   - Links to the previous rollup by hash
//   - Has a first batch whose parent is the head batch of the previous rollup
func (re *RollupManager) checkRollupsCorrectlyChained(rollup *core.Rollup, previousRollup *core.Rollup) error {
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
