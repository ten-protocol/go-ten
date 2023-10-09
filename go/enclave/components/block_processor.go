package components

import (
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	"github.com/obscuronet/go-obscuro/go/enclave/gas"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
)

type l1BlockProcessor struct {
	storage              storage.Storage
	gasOracle            gas.Oracle
	logger               gethlog.Logger
	crossChainProcessors *crosschain.Processors

	// we store the l1 head to avoid expensive db access
	// the host is responsible to always submitting the head l1 block
	currentL1Head *common.L1BlockHash
}

func NewBlockProcessor(storage storage.Storage, cc *crosschain.Processors, gasOracle gas.Oracle, logger gethlog.Logger) L1BlockProcessor {
	var l1BlockHash *common.L1BlockHash
	head, err := storage.FetchHeadBlock()
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
	}
}

func (bp *l1BlockProcessor) Process(br *common.BlockAndReceipts) (*BlockIngestionType, error) {
	defer core.LogMethodDuration(bp.logger, measure.NewStopwatch(), "L1 block processed", log.BlockHashKey, br.Block.Hash())

	ingestion, err := bp.tryAndInsertBlock(br)
	if err != nil {
		return nil, err
	}

	if !ingestion.PreGenesis {
		// This requires block to be stored first ... but can permanently fail a block
		err = bp.crossChainProcessors.Remote.StoreCrossChainMessages(br.Block, *br.Receipts)
		if err != nil {
			return nil, errors.New("failed to process cross chain messages")
		}

		err = bp.crossChainProcessors.Remote.StoreCrossChainValueTransfers(br.Block, *br.Receipts)
		if err != nil {
			return nil, fmt.Errorf("failed to process cross chain transfers. Cause: %w", err)
		}
	}

	// todo @siliev - not sure if this is the best way to update the price, will pick up random stale blocks from forks?
	bp.gasOracle.ProcessL1Block(br.Block)

	h := br.Block.Hash()
	bp.currentL1Head = &h
	return ingestion, nil
}

func (bp *l1BlockProcessor) tryAndInsertBlock(br *common.BlockAndReceipts) (*BlockIngestionType, error) {
	block := br.Block

	_, err := bp.storage.FetchBlock(block.Hash())
	if err == nil {
		return nil, errutil.ErrBlockAlreadyProcessed
	}

	if !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, err := bp.ingestBlock(block)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return nil, err
	}
	bp.logger.Trace("Block inserted successfully",
		log.BlockHeightKey, block.NumberU64(), log.BlockHashKey, block.Hash(), "ingestionType", ingestionType)

	err = bp.storage.StoreBlock(block, ingestionType.ChainFork)
	if err != nil {
		return nil, fmt.Errorf("1. could not store block. Cause: %w", err)
	}

	return ingestionType, nil
}

func (bp *l1BlockProcessor) ingestBlock(block *common.L1Block) (*BlockIngestionType, error) {
	// todo (#1056) - this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := bp.GetHead()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo (@matt) - we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &BlockIngestionType{PreGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	// we do a basic sanity check, comparing the received block to the head block on the chain
	if block.ParentHash() != prevL1Head.Hash() {
		chainFork, err := gethutil.LCA(block, prevL1Head, bp.storage)
		if err != nil {
			bp.logger.Trace("parent not found",
				"blkHeight", block.NumberU64(), log.BlockHashKey, block.Hash(),
				"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
				log.ErrKey, err,
			)
			return nil, errutil.ErrBlockAncestorNotFound
		}

		if chainFork.IsFork() {
			bp.logger.Info("Fork detected in the l1 chain", "can", chainFork.CommonAncestor.Hash().Hex(), "noncan", prevL1Head.Hash().Hex())
		}
		return &BlockIngestionType{ChainFork: chainFork, PreGenesis: false}, nil
	}

	return &BlockIngestionType{ChainFork: nil, PreGenesis: false}, nil
}

func (bp *l1BlockProcessor) GetHead() (*common.L1Block, error) {
	if bp.currentL1Head == nil {
		return nil, errutil.ErrNotFound
	}
	return bp.storage.FetchBlock(*bp.currentL1Head)
}

func (bp *l1BlockProcessor) GetCrossChainContractAddress() *gethcommon.Address {
	return bp.crossChainProcessors.Remote.GetBusAddress()
}
