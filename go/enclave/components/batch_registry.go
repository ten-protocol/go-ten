package components

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/measure"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/state"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/async"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type batchRegistry struct {
	storage      storage.Storage
	logger       gethlog.Logger
	headBatchSeq atomic.Pointer[big.Int] // keep track of the last executed batch to optimise db access

	batchesCallback   func(*core.Batch, types.Receipts)
	callbackMutex     sync.RWMutex
	healthTimeout     time.Duration
	lastExecutedBatch *async.Timestamp
	ethChainAdapter   *EthChainAdapter
}

func NewBatchRegistry(storage storage.Storage, config *enclaveconfig.EnclaveConfig, gethEncodingService gethencoding.EncodingService, logger gethlog.Logger) BatchRegistry {
	var headBatchSeq *big.Int
	headBatch, err := storage.FetchHeadBatchHeader(context.Background())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			headBatchSeq = nil
		} else {
			logger.Crit("Could not create batch registry", log.ErrKey, err)
			return nil
		}
	} else {
		headBatchSeq = headBatch.SequencerOrderNo
	}
	br := &batchRegistry{
		storage:           storage,
		logger:            logger,
		healthTimeout:     time.Minute,
		lastExecutedBatch: async.NewAsyncTimestamp(time.Now().Add(-time.Minute)),
	}
	br.headBatchSeq.Store(headBatchSeq)

	br.ethChainAdapter = NewEthChainAdapter(big.NewInt(config.TenChainID), br, storage, gethEncodingService, *config, logger)
	return br
}

func (br *batchRegistry) EthChain() *EthChainAdapter {
	return br.ethChainAdapter
}

func (br *batchRegistry) HeadBatchSeq() *big.Int {
	return br.headBatchSeq.Load()
}

func (br *batchRegistry) SubscribeForExecutedBatches(callback func(*core.Batch, types.Receipts)) {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()
	br.batchesCallback = callback
}

func (br *batchRegistry) UnsubscribeFromBatches() {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()

	br.batchesCallback = nil
}

func (br *batchRegistry) OnL1Reorg(_ *BlockIngestionType) {
	// refresh the cached head batch from the database because there was an L1 reorg
	headBatch, err := br.storage.FetchHeadBatchHeader(context.Background())
	if err != nil {
		br.logger.Error("Could not fetch head batch", log.ErrKey, err)
		return
	}
	br.headBatchSeq.Store(headBatch.SequencerOrderNo)
}

func (br *batchRegistry) OnBatchExecuted(batchHeader *common.BatchHeader, txExecResults []*core.TxExecResult) error {
	defer core.LogMethodDuration(br.logger, measure.NewStopwatch(), "OnBatchExecuted", log.BatchHashKey, batchHeader.Hash())

	txs, err := br.storage.FetchBatchTransactionsBySeq(context.Background(), batchHeader.SequencerOrderNo.Uint64())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		// this function is called after a batch was successfully executed. This is a catastrophic failure
		br.logger.Crit("should not happen. cannot get transactions. ", log.ErrKey, err)
	}
	batch := &core.Batch{
		Header:       batchHeader,
		Transactions: txs,
	}
	err = br.ethChainAdapter.IngestNewBlock(batch)
	if err != nil {
		return fmt.Errorf("failed to feed batch into the virtual eth chain. cause %w", err)
	}

	br.headBatchSeq.Store(batchHeader.SequencerOrderNo)

	br.callbackMutex.RLock()
	callback := br.batchesCallback
	br.callbackMutex.RUnlock()

	if callback != nil {
		txReceipts := make([]*types.Receipt, len(txExecResults))
		for i, txExecResult := range txExecResults {
			txReceipts[i] = txExecResult.Receipt
		}
		callback(batch, txReceipts)
	}

	br.lastExecutedBatch.Mark()
	return nil
}

func (br *batchRegistry) HasGenesisBatch() (bool, error) {
	return br.HeadBatchSeq() != nil, nil
}

func (br *batchRegistry) BatchesAfter(ctx context.Context, batchSeqNo uint64, upToL1Height uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, []*types.Header, error) {
	// sanity check
	headBatch, err := br.storage.FetchBatchHeaderBySeqNo(ctx, br.HeadBatchSeq().Uint64())
	if err != nil {
		return nil, nil, err
	}

	if headBatch.SequencerOrderNo.Uint64() < batchSeqNo {
		return nil, nil, fmt.Errorf("head batch height %d is in the past compared to requested batch %d", headBatch.SequencerOrderNo.Uint64(), batchSeqNo)
	}

	resultBatches := make([]*core.Batch, 0)
	resultBlocks := make([]*types.Header, 0)

	currentBatchSeq := batchSeqNo
	var currentBlock *types.Header
	for currentBatchSeq <= headBatch.SequencerOrderNo.Uint64() {
		batch, err := br.storage.FetchBatchBySeqNo(ctx, currentBatchSeq)
		if err != nil {
			return nil, nil, fmt.Errorf("could not retrieve batch by sequence number %d. Cause: %w", currentBatchSeq, err)
		}

		// check the block height
		// if it's the same block as the previous batch there is no reason to check
		if currentBlock == nil || currentBlock.Hash() != batch.Header.L1Proof {
			block, err := br.storage.FetchBlock(ctx, batch.Header.L1Proof)
			if err != nil {
				return nil, nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
			}
			currentBlock = block
			if block.Number.Uint64() > upToL1Height {
				break
			}
			resultBlocks = append(resultBlocks, block)
		}

		// check the limiter
		didAcceptBatch, err := rollupLimiter.AcceptBatch(batch)
		if err != nil {
			return nil, nil, err
		}
		if !didAcceptBatch {
			break
		}

		resultBatches = append(resultBatches, batch)
		br.logger.Info("Added batch to rollup", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SeqNo(), log.BatchHeightKey, batch.Number(), "l1_proof", batch.Header.L1Proof)

		currentBatchSeq++
	}

	if len(resultBatches) > 0 {
		// Sanity check that the rollup includes consecutive batches (according to the seqNo)
		current := resultBatches[0].SeqNo().Uint64()
		for i, b := range resultBatches {
			if current+uint64(i) != b.SeqNo().Uint64() {
				return nil, nil, fmt.Errorf("created invalid rollup with batches out of sequence")
			}
		}
	}

	return resultBatches, resultBlocks, nil
}

func (br *batchRegistry) GetBatchState(ctx context.Context, blockNumberOrHash gethrpc.BlockNumberOrHash) (*state.StateDB, error) {
	if blockNumberOrHash.BlockHash != nil {
		return getBatchState(ctx, br.storage, *blockNumberOrHash.BlockHash)
	}
	if blockNumberOrHash.BlockNumber != nil {
		return br.GetBatchStateAtHeight(ctx, blockNumberOrHash.BlockNumber)
	}
	return nil, fmt.Errorf("block number or block hash does not exist")
}

func (br *batchRegistry) GetBatchStateAtHeight(ctx context.Context, blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the batch of interest.
	batch, err := br.GetBatchAtHeight(ctx, *blockNumber)
	if err != nil {
		return nil, err
	}

	return getBatchState(ctx, br.storage, batch.Hash())
}

func getBatchState(ctx context.Context, storage storage.Storage, batchHash common.L2BatchHash) (*state.StateDB, error) {
	blockchainState, err := storage.CreateStateDB(ctx, batchHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for batch %s", batchHash.Hex())
	}

	return blockchainState, err
}

func (br *batchRegistry) GetBatchAtHeight(ctx context.Context, height gethrpc.BlockNumber) (*core.Batch, error) {
	if br.HeadBatchSeq() == nil {
		return nil, fmt.Errorf("chain not initialised")
	}
	var batch *core.Batch
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisBatch, err := br.storage.FetchBatchByHeight(ctx, 0)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		batch = genesisBatch
	// note: our API currently treats all these block statuses the same for obscuro batches
	case gethrpc.SafeBlockNumber, gethrpc.FinalizedBlockNumber, gethrpc.LatestBlockNumber, gethrpc.PendingBlockNumber:
		headBatch, err := br.storage.FetchBatchBySeqNo(ctx, br.HeadBatchSeq().Uint64())
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d was not found. Cause: %w", height, err)
		}
		batch = headBatch
	default:
		maybeBatch, err := br.storage.FetchBatchByHeight(ctx, uint64(height))
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		batch = maybeBatch
	}
	return batch, nil
}

// HealthCheck checks if the last executed batch was more than healthTimeout ago
func (br *batchRegistry) HealthCheck() (bool, error) {
	lastExecutedBatchTime := br.lastExecutedBatch.LastTimestamp()
	if time.Now().After(lastExecutedBatchTime.Add(br.healthTimeout)) {
		return false, fmt.Errorf("last executed batch was %s ago", time.Since(lastExecutedBatchTime))
	}

	if br.HeadBatchSeq() == nil {
		return false, fmt.Errorf("head batch seq is nil")
	}

	return true, nil
}
