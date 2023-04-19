package components

// todo (@stefan) - once the cross chain messages based bridge is implemented remove this completely

import (
	"errors"
	"fmt"
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// rollupProducer encapsulates the logic of decoding rollup transactions submitted to the L1 and resolving them
// to rollups that the enclave can process.
type rollupProducer struct {
	// TransactionBlobCrypto- This contains the required properties to encrypt rollups.
	TransactionBlobCrypto crypto.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	storage db.Storage

	registry      BatchRegistry
	blockConsumer BlockConsumer
}

func NewRollupProducer(
	transactionBlobCrypto crypto.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
	storage db.Storage,
	registry BatchRegistry,
	blockConsumer BlockConsumer,
	logger gethlog.Logger,
) RollupProducer {
	return &rollupProducer{
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
		registry:              registry,
		blockConsumer:         blockConsumer,
		storage:               storage,
	}
}

// fetchLatestRollup - Will pull the latest rollup based on the current head block from the database or return null
func (re *rollupProducer) fetchLatestRollup() (*core.Rollup, error) {
	b, err := re.blockConsumer.GetHead()
	if err != nil {
		return nil, err
	}
	return re.getLatestRollupBeforeBlock(b)
}

func (re *rollupProducer) CreateRollup() (*core.Rollup, error) {
	rollup, err := re.fetchLatestRollup()
	if err != nil && !errors.Is(err, db.ErrNoRollups) {
		return nil, err
	}

	hash := gethcommon.Hash{}
	if rollup != nil {
		hash = rollup.Header.HeadBatchHash
	}

	batches, err := re.registry.BatchesAfter(hash)
	if err != nil {
		return nil, err
	}

	if len(batches) == 0 {
		return nil, fmt.Errorf("no batches for rollup")
	}

	if batches[len(batches)-1].Header.Hash() == hash {
		return nil, fmt.Errorf("current head batch matches the rollup head bash")
	}

	newRollup := createNextRollup(rollup, batches)
	return newRollup, nil
}

// getLatestRollupBeforeBlock - Given a block, returns the latest rollup in the canonical chain for that block (excluding those in the block itself).
func (re *rollupProducer) getLatestRollupBeforeBlock(block *common.L1Block) (*core.Rollup, error) {
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
		// todo (@stefan) - is this still required for safety, even though we're storing an entry for every L1 block?
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

// createNextRollup - based on a previous rollup and batches will create a new rollup that encapsulate the state
// transition from the old rollup to the new one's head batch.
func createNextRollup(rollup *core.Rollup, batches []*core.Batch) *core.Rollup {
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
