package ethadapter

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

type blockProviderStatus int32

const (
	statusCodeStopped blockProviderStatus = iota
	statusCodeRunning
)

const (
	waitingForBlockTimeout = 30 * time.Second
)

var one = big.NewInt(1)

func NewEthBlockProvider(ethClient EthClient, logger gethlog.Logger) *EthBlockProvider {
	return &EthBlockProvider{
		ethClient:     ethClient,
		ctx:           context.TODO(),
		streamCh:      make(chan *types.Block),
		runningStatus: new(int32),
		logger:        logger,
	}
}

// EthBlockProvider streams blocks from the ethereum L1 client in the order expected by the enclave (consecutive, canonical blocks)
type EthBlockProvider struct {
	ethClient EthClient
	ctx       context.Context
	logger    gethlog.Logger
	streamCh  chan *types.Block

	// the state of the block provider
	currentHead   *types.Header // most recently sent block (reset to nil when a `StartStreaming...` method is called)
	runningStatus *int32        // 0 = stopped, 1 = running
}

// StartStreamingFromHash will look up the hash block, find the appropriate height (LCA if there have been forks) and
// then call StartStreamingFromHeight based on that
func (e *EthBlockProvider) StartStreamingFromHash(latestHash gethcommon.Hash) (<-chan *types.Block, error) {
	ancestorBlk, err := e.latestCanonAncestor(latestHash)
	if err != nil {
		return nil, err
	}
	return e.StartStreamingFromHeight(increment(ancestorBlk.Number()))
}

// StartStreamingFromHeight will (re)start streaming from the given height, closing out any existing stream channel and
// returning the fresh channel - the next block will be the requested height
func (e *EthBlockProvider) StartStreamingFromHeight(height *big.Int) (<-chan *types.Block, error) {
	// block heights start at 1
	if height.Cmp(one) < 0 {
		height = one
	}
	e.streamCh = make(chan *types.Block)
	if e.getStatus() == statusCodeStopped {
		// if the provider is stopped (or not yet started) then we kick off the streaming processes
		e.setStatus(statusCodeRunning)
		go e.streamBlocks(height)
	}
	return e.streamCh, nil
}

func (e *EthBlockProvider) Stop() {
	e.setStatus(statusCodeStopped)
}

func (e *EthBlockProvider) IsLive(h gethcommon.Hash) bool {
	l1Head, f := e.ethClient.FetchHeadBlock()
	return f && h == l1Head.Hash()
}

// streamBlocks is the main loop. It should be run in a separate go routine. It will stream catch-up blocks from requested height until it
// reaches the latest live block, then it will block until next live block arrives
// It blocks when:
// - publishing a block, it blocks on the outbound channel until the block is consumed
// - awaiting a live block, when consumer is completely up-to-date it waits for a live block to arrive
func (e *EthBlockProvider) streamBlocks(fromHeight *big.Int) {
	for e.getStatus() != statusCodeStopped {
		// this will block if we're up-to-date with live blocks
		block, err := e.fetchNextCanonicalBlock(fromHeight)
		if err != nil {
			e.logger.Error("unexpected error while preparing block to stream, will retry in 1 sec", log.ErrKey, err)
			time.Sleep(time.Second)
			continue
		}
		e.logger.Trace("blockProvider streaming block", "height", block.Number(), "hash", block.Hash())
		e.streamCh <- block // we block here until consumer takes it
		// update stream state
		e.currentHead = block.Header()
	}
}

func (e *EthBlockProvider) fetchNextCanonicalBlock(fromHeight *big.Int) (*types.Block, error) {
	if e.currentHead == nil {
		blk, err := e.ethClient.BlockByNumber(fromHeight)
		if err == nil {
			return blk, nil
		}
	}

	l1Block, f := e.ethClient.FetchHeadBlock()
	if !f {
		return nil, errors.New("l1 head block not found")
	}
	l1Head := l1Block.Header()

	// check if the block provider is already up-to-date with the latest L1 block (if it is, then wait for a new block)
	if e.currentHead != nil && l1Head.Hash() == e.currentHead.Hash() {
		var err error
		// block provider's current head is up-to-date, we will wait for the next L1 block to arrive
		l1Head, err = e.awaitNewBlock()
		if err != nil {
			return nil, fmt.Errorf("no new block found from stream - %w", err)
		}
	}

	// most common path should be: new head block that arrived is the next block, and needs to be sent
	if l1Head.ParentHash == e.currentHead.Hash() {
		blk, err := e.ethClient.BlockByHash(l1Head.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not fetch block with hash=%s - %w", l1Head.Hash(), err)
		}
		return blk, nil
	}

	// and if not then, we walk back up the tree from the current head, to find the most recent **canonical** block we sent
	latestCanon, err := e.latestCanonAncestor(e.currentHead.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not find ancestor on canonical chain for hash=%s - %w", e.currentHead.Hash(), err)
	}
	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := e.ethClient.BlockByNumber(increment(latestCanon.Number()))
	if err != nil {
		return nil, fmt.Errorf("could not find block after canon fork branch, height=%s - %w", increment(latestCanon.Number()), err)
	}
	return blk, nil
}

func (e *EthBlockProvider) latestCanonAncestor(blkHash gethcommon.Hash) (*types.Block, error) {
	blk, err := e.ethClient.BlockByHash(blkHash)
	if err != nil {
		return nil, err
	}
	canonAtSameHeight, err := e.ethClient.BlockByNumber(blk.Number())
	if err != nil {
		return nil, err
	}
	if blk.Hash() != canonAtSameHeight.Hash() {
		return e.latestCanonAncestor(blk.ParentHash())
	}
	return blk, nil
}

// awaitNewBlock will block, waiting for the next live L1 block from the L1 client.
// It is used when the block provider is up-to-date and waiting for a new block to forward to its consumer.
func (e *EthBlockProvider) awaitNewBlock() (*types.Header, error) {
	liveStream, streamSub := e.ethClient.BlockListener()

	for {
		select {
		case blkHead := <-liveStream:
			e.logger.Trace("received new L1 head", "height", blkHead.Number, "hash", blkHead.Hash())
			streamSub.Unsubscribe()
			return blkHead, nil

		case <-e.ctx.Done():
			return nil, fmt.Errorf("context closed before block was received")

		case err := <-streamSub.Err():
			if errors.Is(err, ErrSubscriptionNotSupported) {
				return nil, err
			}
			e.logger.Error("L1 block monitoring error", log.ErrKey, err)

			e.logger.Info("Restarting L1 block Monitoring...")
			liveStream, streamSub = e.ethClient.BlockListener()

		case <-time.After(waitingForBlockTimeout):
			return nil, fmt.Errorf("no block received from L1 client stream for over %s", waitingForBlockTimeout)
		}
	}
}

func (e *EthBlockProvider) getStatus() blockProviderStatus {
	return blockProviderStatus(atomic.LoadInt32(e.runningStatus))
}
func (e *EthBlockProvider) setStatus(status blockProviderStatus) {
	atomic.StoreInt32(e.runningStatus, int32(status))
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
