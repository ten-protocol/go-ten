package components

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
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

	storage storage.Storage

	batchRegistry  BatchRegistry
	blockProcessor L1BlockProcessor
}

func NewRollupProducer(sequencerID gethcommon.Address, transactionBlobCrypto crypto.DataEncryptionService, obscuroChainID int64, ethereumChainID int64, storage storage.Storage, batchRegistry BatchRegistry, blockProcessor L1BlockProcessor, logger gethlog.Logger) RollupProducer {
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

func (re *rollupProducerImpl) CreateRollup(fromBatchNo uint64, limiter limiters.RollupLimiter) (*core.Rollup, error) {
	batches, err := re.batchRegistry.BatchesAfter(fromBatchNo, limiter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch 'from' batch (seqNo=%d) for rollup: %w", fromBatchNo, err)
	}

	hasBatches := len(batches) != 0

	if !hasBatches {
		return nil, fmt.Errorf("no batches for rollup")
	}

	newRollup := re.createNextRollup(batches)

	re.logger.Info(fmt.Sprintf("Created new rollup %s with %d batches. From %d to %d", newRollup.Hash(), len(newRollup.Batches), batches[0].SeqNo(), batches[len(batches)-1].SeqNo()))

	return newRollup, nil
}

// createNextRollup - based on a previous rollup and batches will create a new rollup that encapsulate the state
// transition from the old rollup to the new one's head batch.
func (re *rollupProducerImpl) createNextRollup(batches []*core.Batch) *core.Rollup {
	lastBatch := batches[len(batches)-1]

	rh := common.RollupHeader{}
	rh.L1Proof = lastBatch.Header.L1Proof
	b, err := re.storage.FetchBlock(rh.L1Proof)
	if err != nil {
		re.logger.Crit("Could not fetch block. Should not happen", log.ErrKey, err)
	}
	rh.L1ProofNumber = b.Number()
	rh.Coinbase = re.sequencerID

	rh.CrossChainMessages = make([]MessageBus.StructsCrossChainMessage, 0)
	for _, b := range batches {
		rh.CrossChainMessages = append(rh.CrossChainMessages, b.Header.CrossChainMessages...)
	}

	rh.LastBatchSeqNo = lastBatch.SeqNo().Uint64()
	return &core.Rollup{
		Header:  &rh,
		Batches: batches,
	}
}
