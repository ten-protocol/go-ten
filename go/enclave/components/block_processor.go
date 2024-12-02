package components

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/async"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
)

type l1BlockProcessor struct {
	storage              storage.Storage
	gasOracle            gas.Oracle
	logger               gethlog.Logger
	crossChainProcessors *crosschain.Processors

	// we store the l1 head to avoid expensive db access
	// the host is responsible to always submitting the head l1 block
	currentL1Head     *common.L1BlockHash
	healthTimeout     time.Duration
	lastIngestedBlock *async.Timestamp
}

func NewBlockProcessor(storage storage.Storage, cc *crosschain.Processors, gasOracle gas.Oracle, logger gethlog.Logger) L1BlockProcessor {
	var l1BlockHash *common.L1BlockHash
	head, err := storage.FetchHeadBlock(context.Background())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			logger.Crit("Cannot fetch head block", log.ErrKey, err)
		}
	} else {
		h := head.Hash()
		l1BlockHash = &h
	}

	return &l1BlockProcessor{
		storage:              storage,
		logger:               logger,
		gasOracle:            gasOracle,
		crossChainProcessors: cc,
		currentL1Head:        l1BlockHash,
		healthTimeout:        time.Minute,
		lastIngestedBlock:    async.NewAsyncTimestamp(time.Now().Add(-time.Minute)),
	}
}

func (bp *l1BlockProcessor) Process(ctx context.Context, processed *common.ProcessedL1Data) (*BlockIngestionType, error) {
	defer core.LogMethodDuration(bp.logger, measure.NewStopwatch(), "L1 block processed", log.BlockHashKey, processed.BlockHeader.Hash())
	header := processed.BlockHeader
	ingestion, err := bp.tryAndInsertBlock(ctx, processed.BlockHeader)

	if len(processed.Events) > 0 {
		println("EVENTS PRESENT")
	}
	if err != nil {
		return nil, err
	}

	if !ingestion.PreGenesis {
		// This requires block to be stored first ... but can permanently fail a block
		err = bp.crossChainProcessors.Remote.StoreCrossChainMessages(ctx, header, processed)
		if err != nil {
			return nil, errors.New("failed to process cross chain messages")
		}

		err = bp.crossChainProcessors.Remote.StoreCrossChainValueTransfers(ctx, header, processed)
		if err != nil {
			return nil, fmt.Errorf("failed to process cross chain transfers. Cause: %w", err)
		}
	}

	// todo @siliev - not sure if this is the best way to update the price, will pick up random stale blocks from forks?
	bp.gasOracle.ProcessL1Block(header)

	h := header.Hash()
	bp.currentL1Head = &h
	bp.lastIngestedBlock.Mark()
	return ingestion, nil
}

// HealthCheck checks if the last ingested block was more than healthTimeout ago
func (bp *l1BlockProcessor) HealthCheck() (bool, error) {
	lastIngestedBlockTime := bp.lastIngestedBlock.LastTimestamp()
	if time.Now().After(lastIngestedBlockTime.Add(bp.healthTimeout)) {
		return false, fmt.Errorf("last ingested block was %s ago", time.Since(lastIngestedBlockTime))
	}

	return true, nil
}

func (bp *l1BlockProcessor) tryAndInsertBlock(ctx context.Context, block *types.Header) (*BlockIngestionType, error) {
	// We insert the block into the L1 chain and store it.
	// in case the block already exists in the database, this will be treated like a fork, because the head changes to
	// the block that was already saved
	ingestionType, err := bp.ingestBlock(ctx, block)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return nil, err
	}

	if ingestionType.OldCanonicalBlock {
		return nil, errutil.ErrBlockAlreadyProcessed
	}

	err = bp.storage.StoreBlock(ctx, block, ingestionType.ChainFork)
	if err != nil {
		return nil, fmt.Errorf("1. could not store block. Cause: %w", err)
	}

	bp.logger.Trace("BlockHeader inserted successfully",
		log.BlockHeightKey, block.Number, log.BlockHashKey, block.Hash(), "parentHash", block.ParentHash, "ingestionType", ingestionType)

	return ingestionType, nil
}

func (bp *l1BlockProcessor) ingestBlock(ctx context.Context, block *types.Header) (*BlockIngestionType, error) {
	// todo (#1056) - this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := bp.GetHead(ctx)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo (@matt) - we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &BlockIngestionType{PreGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}

	if prevL1Head.Hash() == block.Hash() {
		return &BlockIngestionType{OldCanonicalBlock: true}, nil
	}

	// we do a basic sanity check, comparing the received block to the head block on the chain
	if block.ParentHash != prevL1Head.Hash() {
		isCanon, err := bp.storage.IsBlockCanonical(ctx, block.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not check if block is canonical. Cause: %w", err)
		}
		if isCanon {
			return &BlockIngestionType{OldCanonicalBlock: true}, nil
		}

		chainFork, err := gethutil.LCA(ctx, block, prevL1Head, bp.storage)
		if err != nil {
			bp.logger.Trace("cannot calculate the fork for received block",
				"blkHeight", block.Number.Uint64(), log.BlockHashKey, block.Hash(),
				"l1HeadHeight", prevL1Head.Number.Uint64(), "l1HeadHash", prevL1Head.Hash(),
				log.ErrKey, err,
			)
			return nil, errutil.ErrBlockAncestorNotFound
		}

		if chainFork.IsFork() {
			bp.logger.Info("Fork detected in the l1 chain", "can", chainFork.CommonAncestor.Hash(), "noncan", prevL1Head.Hash())
		} else {
			bp.logger.Error("Should not happen. Weird Fork detected in the l1 chain", "fork", chainFork)
		}
		return &BlockIngestionType{ChainFork: chainFork, PreGenesis: false}, nil
	}

	return &BlockIngestionType{ChainFork: nil, PreGenesis: false}, nil
}

func (bp *l1BlockProcessor) GetHead(ctx context.Context) (*types.Header, error) {
	if bp.currentL1Head == nil {
		return nil, errutil.ErrNotFound
	}
	return bp.storage.FetchBlock(ctx, *bp.currentL1Head)
}

func (bp *l1BlockProcessor) GetCrossChainContractAddress() *gethcommon.Address {
	return bp.crossChainProcessors.Remote.GetBusAddress()
}
