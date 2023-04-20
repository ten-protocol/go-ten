package components

import (
	"errors"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type l1BlockProcessor struct {
	storage              db.Storage
	logger               gethlog.Logger
	crossChainProcessors *crosschain.Processors
}

func NewBlockConsumer(storage db.Storage, cc *crosschain.Processors, logger gethlog.Logger) L1BlockProcessor {
	return &l1BlockProcessor{
		storage:              storage,
		logger:               logger,
		crossChainProcessors: cc,
	}
}

func (bc *l1BlockProcessor) Process(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, error) {
	ingestion, err := bc.tryAndInsertBlock(br, isLatest)
	if err != nil {
		return nil, err
	}

	if !ingestion.PreGenesis {
		// This requires block to be stored first ... but can permanently fail a block
		err = bc.crossChainProcessors.Remote.StoreCrossChainMessages(br.Block, *br.Receipts)
		if err != nil {
			return nil, errors.New("failed to process cross chain messages")
		}
	}

	return ingestion, nil
}

func (bc *l1BlockProcessor) tryAndInsertBlock(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, error) {
	block := br.Block

	_, err := bc.storage.FetchBlock(block.Hash())
	if err == nil {
		return nil, common.ErrBlockAlreadyProcessed
	}

	if !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, err := bc.ingestBlock(block, isLatest)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return nil, err
	}
	bc.logger.Trace("block inserted successfully",
		"height", block.NumberU64(), "hash", block.Hash(), "ingestionType", ingestionType)

	bc.storage.StoreBlock(block)
	if isLatest {
		return ingestionType, bc.storage.UpdateL1Head(block.Hash())
	}

	return ingestionType, nil
}

func (bc *l1BlockProcessor) ingestBlock(block *common.L1Block, isLatest bool) (*BlockIngestionType, error) {
	// todo (#1056) - this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := bc.storage.FetchHeadBlock()

	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo (@matt) - we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &BlockIngestionType{IsLatest: isLatest, Fork: false, PreGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)

		// we do a basic sanity check, comparing the received block to the head block on the chain
	} else if block.ParentHash() != prevL1Head.Hash() {
		lcaBlock, err := gethutil.LCA(block, prevL1Head, bc.storage)
		if err != nil {
			bc.logger.Trace("parent not found",
				"blkHeight", block.NumberU64(), log.BlockHashKey, block.Hash(),
				"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
			)
			return nil, common.ErrBlockAncestorNotFound
		}

		// todo - this whole check is iffy, we found an lca, who cares where the head is set at.
		if lcaBlock.NumberU64() < prevL1Head.NumberU64() {
			// fork - least common ancestor for this block and l1 head is before the l1 head.
			return &BlockIngestionType{IsLatest: isLatest, Fork: true, PreGenesis: false}, nil
		}
	}
	// this is the typical, happy-path case. The ingested block's parent was the previously ingested block.
	return &BlockIngestionType{IsLatest: isLatest, Fork: false, PreGenesis: false}, nil
}

func (bc *l1BlockProcessor) GetHead() (*common.L1Block, error) {
	return bc.storage.FetchHeadBlock()
}
