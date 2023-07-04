package components

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/measure"

	"github.com/obscuronet/go-obscuro/go/common/compression"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	dataEncryptionService  crypto.DataEncryptionService
	dataCompressionService compression.DataCompressionService
	batchRegistry          BatchRegistry

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger

	storage      db.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
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
		batchRegistry:          batchRegistry,
		dataEncryptionService:  dataEncryptionService,
		dataCompressionService: dataCompressionService,
		ObscuroChainID:         obscuroChainID,
		EthereumChainID:        ethereumChainID,
		logger:                 logger,
		storage:                storage,
		sigValidator:           verifier,
	}
}

func (rc *rollupConsumerImpl) ProcessL1Block(b *common.BlockAndReceipts) (*common.ExtRollup, error) {
	stopwatch := measure.NewStopwatch()
	defer rc.logger.Info("Rollup consumer processed block", log.BlockHashKey, b.Block.Hash(), log.DurationKey, stopwatch)

	rollups := rc.extractRollups(b, rc.storage)
	if len(rollups) == 0 {
		return nil, nil //nolint:nilnil
	}

	rollup, err := rc.getCanonicalRollup(rollups, b)
	if err != nil {
		return nil, err
	}

	err = rc.processRollup(b.Block, rollup)
	if err != nil {
		return nil, err
	}
	return rollup, nil
}

func (rc *rollupConsumerImpl) getCanonicalRollup(rollups []*common.ExtRollup, b *common.BlockAndReceipts) (*common.ExtRollup, error) {
	var signedRollup *common.ExtRollup

	// loop through the rollups, find the one that is signed, verify the signature, make sure it's the only one
	for _, rollup := range rollups {
		if err := rc.sigValidator.CheckSequencerSignature(rollup.Hash(), rollup.Header.R, rollup.Header.S); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}
		if signedRollup != nil {
			// todo (@matt) - make sure this can't be used to DOS the network
			// we should never receive multiple signed rollups in a single block, the host should only ever publish one
			return nil, fmt.Errorf("received multiple signed rollups in single block %s", b.Block.Hash())
		}
		signedRollup = rollup
	}
	return signedRollup, nil
}

// extractRollups - returns a list of the rollups published in this block
func (rc *rollupConsumerImpl) extractRollups(br *common.BlockAndReceipts, blockResolver db.BlockResolver) []*common.ExtRollup {
	rollups := make([]*common.ExtRollup, 0)
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

		rollups = append(rollups, r)
		rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.RollupHeightKey, r.Header.Number, log.BlockHashKey, b.Hash())
	}

	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	return rollups
}

// IsGenesis indicates whether the rollup is the genesis rollup.
// todo (#718) - Change this to a check against a hardcoded genesis hash.
func IsGenesis(rollupHeader *common.RollupHeader) bool {
	return rollupHeader.Number.Cmp(big.NewInt(int64(common.L2GenesisHeight))) == 0
}

// Validates and stores the rollup in a given block. Returns nil, nil when no rollup was found.
// todo (#718) - design a mechanism to detect a case where the rollup doesn't contain any batches (despite batches arriving via P2P)
func (rc *rollupConsumerImpl) processRollup(block *types.Block, rollup *common.ExtRollup) error {
	// do we need to store the entire rollup?
	// Should be enough to store the header and the batches
	if err := rc.storage.StoreRollup(rollup); err != nil {
		// todo (@matt) - this seems catastrophic, how do we recover the lost rollup in this case?
		return fmt.Errorf("could not store rollup. Cause: %w", err)
	}

	// we record the latest rollup published against this L1 block hash
	rollupHash := rollup.Header.Hash()
	blockHash := block.Hash()
	err := rc.storage.UpdateHeadRollup(&blockHash, &rollupHash)
	if err != nil {
		// todo (@matt) - this also seems catastrophic, would result in bad state unable to ingest further rollups?
		return fmt.Errorf("unable to update head rollup - %w", err)
	}

	return nil
}

func (rc *rollupConsumerImpl) ProcessRollup(rollup *common.ExtRollup) error {
	// todo logic to decompress the rollups on the fly
	r, err := core.ToRollup(rollup, rc.dataEncryptionService, rc.dataCompressionService)
	if err != nil {
		return err
	}

	for _, batch := range r.Batches {
		rc.logger.Info("Processing batch from rollup", log.BatchHashKey, batch.Hash(), "seqNo", batch.SeqNo())
		_, batchFoundErr := rc.batchRegistry.GetBatch(batch.Hash())
		// Process and store a batch only if it wasn't already processed via p2p.
		if batchFoundErr != nil && !errors.Is(batchFoundErr, errutil.ErrNotFound) {
			return batchFoundErr
		}
		receipts, err := rc.batchRegistry.ValidateBatch(batch)
		if errors.Is(err, errutil.ErrBlockForBatchNotFound) || errors.Is(err, errutil.ErrAncestorBatchNotFound) {
			rc.logger.Warn("Unable to validate batch due to it being on a different chain.", log.BatchHashKey, batch.Hash())
			continue
		}
		if err != nil {
			rc.logger.Error("Failed validating batch", log.BatchHashKey, batch.Hash(), log.ErrKey, err)
			return fmt.Errorf("failed validating and storing batch. Cause: %w", err)
		}

		err = rc.batchRegistry.StoreBatch(batch, receipts)
		if err != nil {
			return err
		}
	}
	return nil
}
