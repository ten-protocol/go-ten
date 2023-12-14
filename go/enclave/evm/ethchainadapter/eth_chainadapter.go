package ethchainadapter

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcore "github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// EthChainAdapter is an obscuro wrapper around the ethereum core.Blockchain object
type EthChainAdapter struct {
	newHeadChan   chan gethcore.ChainHeadEvent
	batchRegistry components.BatchRegistry
	storage       storage.Storage
	chainID       *big.Int
	logger        gethlog.Logger
}

// NewEthChainAdapter returns a new instance
func NewEthChainAdapter(chainID *big.Int, batchRegistry components.BatchRegistry, storage storage.Storage, logger gethlog.Logger) *EthChainAdapter {
	return &EthChainAdapter{
		newHeadChan:   make(chan gethcore.ChainHeadEvent),
		batchRegistry: batchRegistry,
		storage:       storage,
		chainID:       chainID,
		logger:        logger,
	}
}

// Config retrieves the chain's fork configuration.
func (e *EthChainAdapter) Config() *params.ChainConfig {
	return ChainParams(e.chainID)
}

// CurrentBlock returns the current head of the chain.
func (e *EthChainAdapter) CurrentBlock() *gethtypes.Header {
	currentBatchSeqNo := e.batchRegistry.HeadBatchSeq()
	if currentBatchSeqNo == nil {
		return nil
	}
	currentBatch, err := e.storage.FetchBatchBySeqNo(currentBatchSeqNo.Uint64())
	if err != nil {
		e.logger.Warn("unable to retrieve batch seq no: %d", "currentBatchSeqNo", currentBatchSeqNo, log.ErrKey, err)
		return nil
	}
	batch, err := gethencoding.CreateEthHeaderForBatch(currentBatch.Header, secret(e.storage))
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
func (e *EthChainAdapter) GetBlock(_ common.Hash, number uint64) *gethtypes.Block {
	nbatch, err := e.storage.FetchBatchByHeight(number)
	if err != nil {
		e.logger.Warn("unable to get batch by height", "number", number, log.ErrKey, err)
		return nil
	}

	nfromBatch, err := gethencoding.CreateEthBlockFromBatch(nbatch)
	if err != nil {
		e.logger.Error("unable to convert batch to eth block", log.ErrKey, err)
		return nil
	}

	return nfromBatch
}

// StateAt returns a state database for a given root hash (generally the head).
func (e *EthChainAdapter) StateAt(root common.Hash) (*state.StateDB, error) {
	if root.Hex() == gethtypes.EmptyCodeHash.Hex() {
		return nil, nil //nolint:nilnil
	}

	currentBatchSeqNo := e.batchRegistry.HeadBatchSeq()
	if currentBatchSeqNo == nil {
		return nil, fmt.Errorf("not ready yet")
	}
	currentBatch, err := e.storage.FetchBatchBySeqNo(currentBatchSeqNo.Uint64())
	if err != nil {
		e.logger.Warn("unable to get batch by height", "currentBatchSeqNo", currentBatchSeqNo, log.ErrKey, err)
		return nil, nil //nolint:nilnil
	}

	return e.storage.CreateStateDB(currentBatch.Hash())
}

func (e *EthChainAdapter) IngestNewBlock(batch *core.Batch) error {
	convertedBlock, err := gethencoding.CreateEthBlockFromBatch(batch)
	if err != nil {
		return err
	}

	go func() {
		e.newHeadChan <- gethcore.ChainHeadEvent{Block: convertedBlock}
	}()

	return nil
}

func NewLegacyPoolConfig() legacypool.Config {
	return legacypool.Config{
		Locals:       nil,
		NoLocals:     false,
		Journal:      "",
		Rejournal:    0,
		PriceLimit:   legacypool.DefaultConfig.PriceLimit,
		PriceBump:    legacypool.DefaultConfig.PriceBump,
		AccountSlots: legacypool.DefaultConfig.AccountSlots,
		GlobalSlots:  legacypool.DefaultConfig.GlobalSlots,
		AccountQueue: legacypool.DefaultConfig.AccountQueue,
		GlobalQueue:  legacypool.DefaultConfig.GlobalQueue,
		Lifetime:     legacypool.DefaultConfig.Lifetime,
	}
}

func secret(storage storage.Storage) []byte {
	// todo (#1053) - handle secret not being found.
	secret, _ := storage.FetchSecret()
	return secret[:]
}
