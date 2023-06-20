package components

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// rollupProducerImpl encapsulates the logic of decoding rollup transactions submitted to the L1 and resolving them
// to rollups that the enclave can process.
type rollupProducerImpl struct {
	// TransactionBlobCrypto- This contains the required properties to encrypt rollups.
	TransactionBlobCrypto crypto.DataEncryptionService

	ObscuroChainID  int64
	EthereumChainID int64

	sequencerID gethcommon.Address

	logger gethlog.Logger

	storage db.Storage

	batchRegistry  BatchRegistry
	blockProcessor L1BlockProcessor
}

func NewRollupProducer(sequencerID gethcommon.Address, transactionBlobCrypto crypto.DataEncryptionService, obscuroChainID int64, ethereumChainID int64, storage db.Storage, batchRegistry BatchRegistry, blockProcessor L1BlockProcessor, logger gethlog.Logger) RollupProducer {
	return &rollupProducerImpl{
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		sequencerID:           sequencerID,
		logger:                logger,
		batchRegistry:         batchRegistry,
		blockProcessor:        blockProcessor,
		storage:               storage,
	}
}

// fetchLatestRollup - Will pull the latest rollup based on the current head block from the database or return null
func (re *rollupProducerImpl) fetchLatestRollup() (*core.Rollup, error) {
	b, err := re.blockProcessor.GetHead()
	if err != nil {
		return nil, err
	}
	return getLatestRollupBeforeBlock(b, re.storage, re.logger)
}

func (re *rollupProducerImpl) CreateRollup(limiter limiters.RollupLimiter) (*core.Rollup, error) {
	rollup, err := re.fetchLatestRollup()
	if err != nil && !errors.Is(err, db.ErrNoRollups) {
		return nil, err
	}

	hash := gethcommon.Hash{}
	if rollup != nil {
		hash = rollup.HeadBatchHash()
	}

	batches, err := re.batchRegistry.BatchesAfter(hash, limiter)
	if err != nil {
		return nil, err
	}

	if len(batches) == 0 {
		return nil, fmt.Errorf("no batches for rollup")
	}

	if batches[len(batches)-1].Hash() == hash {
		return nil, fmt.Errorf("current head batch matches the rollup head bash")
	}

	newRollup := re.createNextRollup(rollup, batches)
	return newRollup, nil
}

// createNextRollup - based on a previous rollup and batches will create a new rollup that encapsulate the state
// transition from the old rollup to the new one's head batch.
func (re *rollupProducerImpl) createNextRollup(parentRollup *core.Rollup, batches []*core.Batch) *core.Rollup {
	lastBatch := batches[len(batches)-1]

	rh := common.RollupHeader{}
	rh.L1Proof = lastBatch.Header.L1Proof
	b, err := re.storage.FetchBlock(rh.L1Proof)
	if err != nil {
		re.logger.Crit("Could not fetch block. Should not happen", log.ErrKey, err)
	}
	rh.L1ProofNumber = b.Number()
	rh.Time = lastBatch.Header.Time
	rh.Coinbase = re.sequencerID

	if parentRollup != nil {
		rh.ParentHash = parentRollup.Header.Hash()
		rh.Number = big.NewInt(parentRollup.Number().Int64() + 1)
	} else { // genesis
		rh.ParentHash = gethcommon.Hash{}
		rh.Number = big.NewInt(0)
	}

	rh.CrossChainMessages = make([]MessageBus.StructsCrossChainMessage, 0)
	for _, b := range batches {
		rh.CrossChainMessages = append(rh.CrossChainMessages, b.Header.CrossChainMessages...)
	}

	rollupHeight := big.NewInt(0)
	if parentRollup != nil {
		rollupHeight = parentRollup.Header.Number
		rollupHeight.Add(rollupHeight, gethcommon.Big1)
	}
	rh.Number = rollupHeight

	return &core.Rollup{
		Header:  &rh,
		Batches: batches,
	}
}

// getLatestRollupBeforeBlock - Given a block, returns the latest rollup in the canonical chain for that block (excluding those in the block itself).
func getLatestRollupBeforeBlock(block *common.L1Block, storage db.Storage, logger gethlog.Logger) (*core.Rollup, error) {
	scanBackCount := 0
	for {
		blockParentHash := block.ParentHash()
		latestRollup, err := storage.FetchHeadRollupForBlock(&blockParentHash)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not fetch current L2 head rollup - %w", err)
		}

		// we found a rollup so we return
		if latestRollup != nil {
			return latestRollup, nil
		}

		// we scan backwards now to the prev block in the chain and we will lookup to see if that has an entry
		block, err = storage.FetchBlock(block.ParentHash())
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				// No more blocks available (enclave does not read the L1 chain from genesis if it knows
				// when management contract was deployed, so we don't keep going to block zero, we just stop when the blocks run out)
				// We have now checked through the entire (relevant) history of the L1 and no rollups were found.
				return nil, db.ErrNoRollups
			}
			return nil, fmt.Errorf("could not fetch parent block - %w", err)
		}
		scanBackCount++
		// todo (@matt) - remove this when we are confident that we are not scanning backwards any more
		if scanBackCount%100 == 0 {
			// if we are scanning a long way backwards (when we don't think we need to, and it might be expensive) we want to know about it
			logger.Warn(fmt.Sprintf("Scanning backwards for rollup, scanned %d blocks backwards so far...", scanBackCount), log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash())
		}
	}
}
