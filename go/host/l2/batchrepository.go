package l2

import (
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

const (
	// if request asks for batches from seq no. X we don't want to return potentially thousands of batches
	// so we limit the number of batches we return with this cap (they can request the next ones afterwards, but they should be catching up from rollups first)
	_maxBatchesInP2PResponse      = 20
	_timeoutWaitingForP2PResponse = 30 * time.Second
)

type Repository struct {
	subscribers []host.L2BatchHandler

	p2p host.P2P
	db  *db.DB

	// high watermark for batch sequence numbers seen so far. If we can't find batch for seq no < this, then we should ask peers for missing batches
	latestBatchSeqNo *big.Int
	latestSeqNoMutex sync.Mutex

	// the repository requests batches from peers asynchronously.
	// We don't want to repeatedly spam out requests if we haven't received a response yet,
	// but we don't want to wait forever if there's no response.
	p2pReqMutex          sync.Mutex
	p2pInFlightRequested *big.Int
	p2pInFlightReqTime   *time.Time

	running atomic.Bool
	logger  gethlog.Logger
}

func NewBatchRepository(_ *config.HostConfig, p2p host.P2P, database *db.DB, logger gethlog.Logger) *Repository {
	return &Repository{
		p2p:              p2p,
		db:               database,
		latestBatchSeqNo: big.NewInt(0),
		running:          atomic.Bool{},
		logger:           logger,
	}
}

func (r *Repository) Start() error {
	r.running.Store(true)

	// register ourselves for new batches from p2p
	r.p2p.SubscribeForBatches(r)
	r.p2p.SubscribeForBatchRequests(r)

	return nil
}

func (r *Repository) Stop() error {
	r.running.Store(false)
	return nil
}

func (r *Repository) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on last received batch or something
	errMsg := ""
	if !r.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (r *Repository) HandleBatches(batches []*common.ExtBatch, isLive bool) {
	// if these batches resolve our in-flight P2P request, clear it
	r.p2pReqMutex.Lock()
	if !isLive && len(batches) > 0 && r.p2pInFlightRequested != nil && batches[0].Header.SequencerOrderNo.Cmp(r.p2pInFlightRequested) == 0 {
		r.p2pInFlightRequested = nil
		r.p2pInFlightReqTime = nil
	}
	r.p2pReqMutex.Unlock()

	// try to add all the batches to the db, and notify subscribers if they are new and live
	for _, batch := range batches {
		err := r.AddBatch(batch)
		if err != nil {
			if !errors.Is(err, errutil.ErrAlreadyExists) {
				r.logger.Warn("unable to add p2p batch to L2 batch repository", log.ErrKey, err)
			}
			// we've already seen this batch or failed to store it for another reason - do not notify subscribers
			return
		}
		if isLive {
			// notify subscribers if the batch is new
			for _, subscriber := range r.subscribers {
				go subscriber.HandleBatch(batch)
			}
		}
	}
}

// HandleBatchRequest handles a request for a batch from a peer, sending batches to the requester asynchronously
// todo (#1625) - only allow requests for batches since last rollup, to avoid DoS attacks.
func (r *Repository) HandleBatchRequest(requesterID string, fromSeqNo *big.Int) {
	batches := make([]*common.ExtBatch, 0)
	nextSeqNum := fromSeqNo
	for len(batches) <= _maxBatchesInP2PResponse {
		batch, err := r.db.GetBatchBySequenceNumber(nextSeqNum)
		if err != nil {
			if !errors.Is(err, errutil.ErrNotFound) {
				r.logger.Warn("unexpected error fetching batches for peer req", "seqNo", nextSeqNum, log.ErrKey, err)
			}
			break // once one batch lookup fails we don't expect to find any of them
		}
		batches = append(batches, batch)
		nextSeqNum = nextSeqNum.Add(nextSeqNum, big.NewInt(1))
	}
	if len(batches) == 0 {
		return // nothing to send
	}

	err := r.p2p.RespondToBatchRequest(requesterID, batches)
	if err != nil {
		r.logger.Warn("unable to send batches to peer", "peer", requesterID, log.ErrKey, err)
	}
}

func (r *Repository) Subscribe(subscriber host.L2BatchHandler) {
	r.subscribers = append(r.subscribers, subscriber)
}

func (r *Repository) FetchBatchBySeqNo(seqNo *big.Int) (*common.ExtBatch, error) {
	b, err := r.db.GetBatchBySequenceNumber(seqNo)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) && seqNo.Cmp(r.latestBatchSeqNo) < 0 {
			// we haven't seen this batch before, but it is older than the latest batch we have seen so far
			// Request missing batches from peers (the batches from any response will be added asynchronously, so
			// we will return the not found error and hopefully future attempts will succeed)
			go r.requestMissingBatches(seqNo)
		}
		return nil, err
	}
	return b, nil
}

func (r *Repository) AddBatch(batch *common.ExtBatch) error {
	err := r.db.AddBatch(batch)
	if err != nil {
		return err
	}
	// atomically compare and swap latest batch sequence number if successfully added batch is newer
	r.latestSeqNoMutex.Lock()
	defer r.latestSeqNoMutex.Unlock()
	if batch.Header.SequencerOrderNo.Cmp(r.latestBatchSeqNo) > 0 {
		r.latestBatchSeqNo = batch.Header.SequencerOrderNo
	}
	return nil
}

// RequestMissingBatches requests batches from peers from the specified sequence number.
// It is an asynchronous request and the repository does not expect to be notified of the result.
func (r *Repository) requestMissingBatches(fromSeqNo *big.Int) {
	r.p2pReqMutex.Lock()
	defer r.p2pReqMutex.Unlock()
	if r.p2pInFlightReqTime != nil && time.Since(*r.p2pInFlightReqTime) < _timeoutWaitingForP2PResponse {
		// don't send request if we have sent one too recently
		r.logger.Trace("not requesting missing batches from sequencer - too soon since last request", "fromSeqNo", fromSeqNo, "lastReq", r.p2pInFlightReqTime)
		return
	}

	r.logger.Debug("requesting missing batches from sequencer", "fromSeqNo", fromSeqNo)
	err := r.p2p.RequestBatchesFromSequencer(fromSeqNo)
	if err != nil {
		r.logger.Warn("unable to request missing batches from sequencer", "fromSeqNo", fromSeqNo, log.ErrKey, err)
		return
	}
	now := time.Now()
	r.p2pInFlightReqTime = &now
}
