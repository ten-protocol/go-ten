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

func (e *EthBlockProvider) start(fromHeight *big.Int) {
	e.runningStatus = new(int32)
	go e.streamBlocks(fromHeight)
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
	if e.streamCh != nil {
		close(e.streamCh)
	}
	e.streamCh = make(chan *types.Block)
	if e.stopped() {
		// if the provider is stopped (or not yet started) then we kick off the streaming processes
		e.start(height)
	}
	return e.streamCh, nil
}

func (e *EthBlockProvider) Stop() {
	atomic.StoreInt32(e.runningStatus, int32(statusCodeStopped))
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
	atomic.StoreInt32(e.runningStatus, int32(statusCodeRunning))
	for !e.stopped() {
		// this will block if we're up-to-date with live blocks
		block, err := e.nextBlockToStream(fromHeight)
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

func (e *EthBlockProvider) nextBlockToStream(fromHeight *big.Int) (*types.Block, error) {
	if e.currentHead == nil {
		blk, err := e.ethClient.BlockByNumber(fromHeight)
		if err == nil {
			return blk, nil
		}
	}

	head, err := e.AwaitNewBlock()
	if err != nil {
		return nil, fmt.Errorf("no new block found from stream - %w", err)
	}

	// most common path should be: new head block that arrived is the next block, and needs to be sent
	if head.ParentHash == e.currentHead.Hash() {
		blk, err := e.ethClient.BlockByHash(head.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not fetch block with hash=%s - %w", head.Hash(), err)
		}
		return blk, nil
	}

	// and if not then, we find the latest canonical block we sent and try one after that
	latestCanon, err := e.latestCanonAncestor(e.currentHead.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not find ancestor on canonical chain for hash=%s - %w", e.currentHead.Hash(), err)
	}
	// and send the cannon block after the last sent (this may be a fork, or it may just be the next on the same branch)
	blk, err := e.ethClient.BlockByNumber(increment(latestCanon.Number()))
	if err != nil {
		return nil, fmt.Errorf("could not find block after canon fork branch, height=%s - %w", increment(latestCanon.Number()), err)
	}
	return blk, nil
}

// checkStopped checks the status for stopped code
func (e *EthBlockProvider) stopped() bool {
	return atomic.LoadInt32(e.runningStatus) == int32(statusCodeStopped)
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

// AwaitNewBlock takes a hash, it will block until it can return a head block with a different hash or error
// (note: this can currently only be used by one caller at a time - not an issue for current usage)
func (e *EthBlockProvider) AwaitNewBlock() (*types.Header, error) {
	// first we check if we're up-to-date
	l1Head, f := e.ethClient.FetchHeadBlock()
	if !f {
		return nil, errors.New("l1 head block not found")
	}
	if e.currentHead == nil || e.currentHead.Hash() != l1Head.Hash() {
		// we're behind, return the head (no need to wait for the next head, there's catching up to do)
		return l1Head.Header(), nil
	}

	// then we start streaming
	liveStream, streamSub := e.ethClient.BlockListener()

	// check again to make sure we didn't miss one
	l1Head, f = e.ethClient.FetchHeadBlock()
	if f && e.currentHead.Hash() != l1Head.Hash() {
		// we're behind now, return the head
		streamSub.Unsubscribe()
		return l1Head.Header(), nil
	}

	// and now we wait...
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

			// check head to make sure we didn't miss one
			l1Head, f = e.ethClient.FetchHeadBlock()
			if f && e.currentHead.Hash() != l1Head.Hash() {
				// we're behind now, return the head
				streamSub.Unsubscribe()
				return l1Head.Header(), nil
			}
		case <-time.After(waitingForBlockTimeout):
			return nil, fmt.Errorf("no block received from L1 client stream for over %s", waitingForBlockTimeout)
		}
	}
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
