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

var (
	one               = big.NewInt(1)
	statusCodeStopped = int32(0)
	statusCodeRunning = int32(1)
)

func NewEthBlockProvider(ethClient EthClient, logger gethlog.Logger) *EthBlockProvider {
	return &EthBlockProvider{
		ethClient:     ethClient,
		ctx:           context.TODO(),
		streamCh:      make(chan *types.Block),
		errCh:         make(chan error),
		runningStatus: new(int32),
		logger:        logger,
	}
}

type EthBlockProvider struct {
	ethClient EthClient
	ctx       context.Context

	liveBlocks *LiveBlocksMonitor // process that streams live blocks to track head and notify waiting callers
	liveCancel context.CancelFunc // cancel for the live monitor process

	streamCh      chan *types.Block
	errCh         chan error
	runningStatus *int32 // 0 = stopped, 1 = running

	latestSent *types.Header // most recently sent block (reset if streamFrom is reset)
	streamFrom *big.Int      // height most-recently requested to stream from

	logger gethlog.Logger
}

func (e *EthBlockProvider) start() {
	e.runningStatus = new(int32)
	liveCtx, liveCancel := context.WithCancel(e.ctx)
	e.liveCancel = liveCancel
	e.liveBlocks = &LiveBlocksMonitor{
		ctx:       liveCtx,
		ethClient: e.ethClient,
		logger:    e.logger,
	}
	go e.streamBlocks()
	go e.liveBlocks.Start() // process to handle incoming stream of L1 blocks
}

func (e *EthBlockProvider) Err() <-chan error {
	return e.errCh
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
	e.streamFrom = height
	if e.streamCh != nil {
		close(e.streamCh)
	}
	e.streamCh = make(chan *types.Block)
	if e.stopped() {
		// if the provider is stopped (or not yet started) then we kick off the streaming processes
		e.start()
	}
	return e.streamCh, nil
}

func (e *EthBlockProvider) Stop() {
	if e.liveCancel != nil {
		e.liveCancel()
	}
	atomic.StoreInt32(e.runningStatus, statusCodeStopped)
}

func (e *EthBlockProvider) IsLive(h gethcommon.Hash) bool {
	if e.liveBlocks.latest == nil {
		// live block streaming not working, we're probably running the in-mem simulation where that's not implemented.
		// So check manually:
		block := e.ethClient.FetchHeadBlock()
		return h == block.Hash()
	}
	return h == e.liveBlocks.latest.Hash()
}

// streamBlocks should be run in a separate go routine. It will stream catch-up blocks from requested height until it
// reaches the latest live block, then it will block until next live block arrives
// It blocks when:
// - publishing a block, it blocks on the outbound channel until the block is consumed
// - awaiting a live block, when consumer is completely up-to-date it waits for a live block to arrive
func (e *EthBlockProvider) streamBlocks() {
	atomic.StoreInt32(e.runningStatus, statusCodeRunning)
	for !e.stopped() {
		// this will block if we're up-to-date with live blocks
		block, err := e.nextBlockToStream()
		if err != nil {
			e.logger.Error("unexpected error while preparing block to stream, will retry in 1 sec", log.ErrKey, err)
			time.Sleep(time.Second)
			continue
		}
		e.logger.Trace("blockProvider streaming block", "height", block.Number(), "hash", block.Hash())
		e.streamCh <- block // we block here until consumer takes it
		// update stream state
		e.latestSent = block.Header()
	}
}

func (e *EthBlockProvider) nextBlockToStream() (*types.Block, error) {
	if e.latestSent == nil {
		blk, err := e.ethClient.BlockByNumber(e.streamFrom)
		if err != nil {
			return nil, fmt.Errorf("no block at requested height=%s - %w", e.streamFrom, err)
		}

		return blk, nil
	}

	head, err := e.liveBlocks.AwaitNewBlock(context.TODO(), e.latestSent.Hash())
	if err != nil {
		return nil, fmt.Errorf("no new block found from stream - %w", err)
	}

	// most common path should be: new head block that arrived is the next block, and needs to be sent
	if head.ParentHash == e.latestSent.Hash() {
		blk, err := e.ethClient.BlockByHash(head.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not fetch block with hash=%s - %w", head.Hash(), err)
		}
		return blk, nil
	}

	// and if not then, we find the latest canonical block we sent and try one after that
	latestCanon, err := e.latestCanonAncestor(e.latestSent.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not find ancestor on canonical chain for hash=%s - %w", e.latestSent.Hash(), err)
	}
	// and send the cannon block after the fork
	blk, err := e.ethClient.BlockByNumber(increment(latestCanon.Number()))
	if err != nil {
		return nil, fmt.Errorf("could not find block after canon fork branch, height=%s - %w", increment(latestCanon.Number()), err)
	}
	return blk, nil
}

// checkStopped checks the stopCh for a signal to stop. This seems dumb, what's the correct way to do this?
func (e *EthBlockProvider) stopped() bool {
	return atomic.LoadInt32(e.runningStatus) == statusCodeStopped
}

// LiveBlocksMonitor manages a process that queues up the latest X blocks (where X=queueCap) in a FIFO queue, streamed from
// an L1 client. If it has X and receives another one it will drop the oldest.
type LiveBlocksMonitor struct {
	ctx       context.Context
	logger    gethlog.Logger
	ethClient EthClient

	latest *types.Header

	awaitChan chan *types.Header // recreated when user starts awaiting, receives true when new head
}

func (l *LiveBlocksMonitor) Start() {
	blkHeadChan, blkSubs := l.ethClient.BlockListener()
	for {
		select {
		case <-l.ctx.Done():
			return

		case blkHead := <-blkHeadChan:
			l.logger.Trace("received new L1 head", "height", blkHead.Number, "hash", blkHead.Hash())
			l.latest = blkHead
			if l.awaitChan != nil {
				l.awaitChan <- blkHead // notify
			}

		case err := <-blkSubs.Err():
			l.logger.Error("L1 block monitoring error", log.ErrKey, err)
			l.logger.Info("Restarting L1 block Monitoring...")
			// this disconnect could result in a gap in the LiveBlocksMonitor queue, but the block provider is responsible
			// for sending a coherent ordering of blocks and will look up blocks if parent not sent
			blkHeadChan, blkSubs = l.ethClient.BlockListener()
		}
	}
}

// AwaitNewBlock takes a hash, it will block until it can return a head block with a different hash or error
// (note: this can currently only be used by one caller at a time - not an issue for current usage)
func (l *LiveBlocksMonitor) AwaitNewBlock(ctx context.Context, afterBlkHash gethcommon.Hash) (*types.Header, error) {
	awaitChan := l.registerAwait()
	streamTimeout := 1 * time.Minute
	if l.latest == nil {
		// note: the in-mem client does not support streaming from ethClient.BlockListener() currently, this allows it to still function
		l.logger.Debug("latest is nil, blockprovider will timeout after 1sec without a block as client streaming may not be working")
		streamTimeout = 1 * time.Second
	} else {
		l.logger.Trace("latest=%s, afterBl=%s\n", l.latest.Hash(), afterBlkHash)
	}
	if l.latest != nil && l.latest.Hash() != afterBlkHash {
		// return immediately if caller is not at latest head
		return l.latest, nil
	}

	// otherwise wait for a block to show up
	for {
		select {
		case blkHead := <-awaitChan:
			return blkHead, nil
		case <-ctx.Done():
			return nil, errors.New("context expired before block found")
		case <-time.After(streamTimeout):
			if l.latest == nil {
				// not seen any streamed blocks, see if we can get head from client
				// note: in-mem sim client doesn't support block streaming, it will always use this path
				block := l.ethClient.FetchHeadBlock()
				if block.Hash() != afterBlkHash {
					return block.Header(), nil
				}
				continue // not seen any streamed block and head isn't new, wait for the next timeout.
			}
			return nil, fmt.Errorf("no live block received after timeout=%s", streamTimeout)
		}
	}
}

func (l *LiveBlocksMonitor) registerAwait() <-chan *types.Header {
	l.awaitChan = make(chan *types.Header)
	return l.awaitChan
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

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
