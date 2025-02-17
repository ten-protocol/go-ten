package components

import (
	"context"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcore "github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// EthChainAdapter is an obscuro wrapper around the ethereum core.Blockchain object
type EthChainAdapter struct {
	newHeadChan   chan gethcore.ChainHeadEvent
	batchRegistry BatchRegistry
	gethEncoding  gethencoding.EncodingService
	storage       storage.Storage
	config        *enclaveconfig.EnclaveConfig
	chainID       *big.Int
	logger        gethlog.Logger
}

// NewEthChainAdapter returns a new instance
func NewEthChainAdapter(chainID *big.Int, batchRegistry BatchRegistry, storage storage.Storage, gethEncoding gethencoding.EncodingService, config *enclaveconfig.EnclaveConfig, logger gethlog.Logger) *EthChainAdapter {
	return &EthChainAdapter{
		newHeadChan:   make(chan gethcore.ChainHeadEvent),
		batchRegistry: batchRegistry,
		storage:       storage,
		gethEncoding:  gethEncoding,
		config:        config,
		chainID:       chainID,
		logger:        logger,
	}
}

// Config retrieves the chain's fork configuration.
func (e *EthChainAdapter) Config() *params.ChainConfig {
	return ethchainadapter.ChainParams(e.chainID)
}

// CurrentBlock returns the current head of the chain.
func (e *EthChainAdapter) CurrentBlock() *gethtypes.Header {
	currentBatchSeqNo := e.batchRegistry.HeadBatchSeq()
	if currentBatchSeqNo == nil {
		return nil
	}
	currentBatch, err := e.storage.FetchBatchHeaderBySeqNo(context.Background(), currentBatchSeqNo.Uint64())
	if err != nil {
		e.logger.Warn("unable to retrieve batch seq no", "currentBatchSeqNo", currentBatchSeqNo, log.ErrKey, err)
		return nil
	}
	batch, err := e.gethEncoding.CreateEthHeaderForBatch(context.Background(), currentBatch)
	if err != nil {
		e.logger.Warn("unable to convert batch to eth header ", "currentBatchSeqNo", currentBatchSeqNo, log.ErrKey, err)
		return nil
	}
	return batch
}

func (e *EthChainAdapter) SubscribeChainHeadEvent(ch chan<- gethcore.ChainHeadEvent) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		for {
			select {
			case head := <-e.newHeadChan:
				select {
				case ch <- head:
				case <-quit:
					return nil
				}
			case <-quit:
				return nil
			}
		}
	})
}

// GetBlock retrieves a specific block, used during pool resets.
func (e *EthChainAdapter) GetBlock(_ gethcommon.Hash, number uint64) *gethtypes.Block {
	var batch *core.Batch
	ctx, cancelCtx := context.WithTimeout(context.Background(), e.config.RPCTimeout)
	defer cancelCtx()

	// to avoid a costly select to the db, check whether the batches requested are the last ones which are cached
	headBatch, err := e.storage.FetchBatchBySeqNo(ctx, e.batchRegistry.HeadBatchSeq().Uint64())
	if err != nil {
		e.logger.Error("unable to get head batch", log.ErrKey, err)
		return nil
	}
	if headBatch.Number().Uint64() == number {
		batch = headBatch
	} else if headBatch.Number().Uint64()-1 == number {
		batch, err = e.storage.FetchBatch(ctx, headBatch.Header.ParentHash)
		if err != nil {
			e.logger.Error("unable to get parent of head batch", log.ErrKey, err, log.BatchHashKey, headBatch.Header.ParentHash)
			return nil
		}
	} else {
		batch, err = e.storage.FetchBatchByHeight(ctx, number)
		if err != nil {
			e.logger.Error("unable to get batch by height", log.BatchHeightKey, number, log.ErrKey, err)
			return nil
		}
	}

	nfromBatch, err := e.gethEncoding.CreateEthBlockFromBatch(ctx, batch)
	if err != nil {
		e.logger.Error("unable to convert batch to eth block", log.ErrKey, err)
		return nil
	}

	return nfromBatch
}

// StateAt returns a state database for a given root hash (generally the head).
func (e *EthChainAdapter) StateAt(root gethcommon.Hash) (*state.StateDB, error) {
	if root.Hex() == gethtypes.EmptyCodeHash.Hex() {
		return nil, nil //nolint:nilnil
	}

	return state.New(root, e.storage.StateDB())
}

func (e *EthChainAdapter) IngestNewBlock(batch *core.Batch) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), e.config.RPCTimeout)
	defer cancelCtx()
	convertedBlock, err := e.gethEncoding.CreateEthBlockFromBatch(ctx, batch)
	if err != nil {
		return err
	}

	go func() {
		e.newHeadChan <- gethcore.ChainHeadEvent{Header: convertedBlock.Header()}
	}()

	return nil
}
