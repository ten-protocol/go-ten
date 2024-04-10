package subscription

import (
	"math/big"
	"sync"
	"sync/atomic"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// NewHeadsService multiplexes new batch header messages from an input channel into multiple subscribers
// also handles unsubscribe
// Note: this is a service which must be Started and Stopped
type NewHeadsService struct {
	inputCh            chan *common.BatchHeader
	convertToEthHeader bool
	notifiersMutex     *sync.RWMutex
	newHeadNotifiers   map[rpc.ID]*rpc.Notifier
	onMessage          func(*common.BatchHeader) error
	stopped            *atomic.Bool
	logger             gethlog.Logger
}

func NewNewHeadsService(inputCh chan *common.BatchHeader, convertToEthHeader bool, logger gethlog.Logger, onMessage func(*common.BatchHeader) error) *NewHeadsService {
	return &NewHeadsService{
		inputCh:            inputCh,
		convertToEthHeader: convertToEthHeader,
		onMessage:          onMessage,
		logger:             logger,
		stopped:            &atomic.Bool{},
		newHeadNotifiers:   make(map[rpc.ID]*rpc.Notifier),
		notifiersMutex:     &sync.RWMutex{},
	}
}

func (nhs *NewHeadsService) Start() error {
	go ForwardFromChannels([]chan *common.BatchHeader{nhs.inputCh}, nhs.stopped, func(head *common.BatchHeader) error {
		if nhs.onMessage != nil {
			err := nhs.onMessage(head)
			if err != nil {
				nhs.logger.Info("failed invoking onMessage callback.", log.ErrKey, err)
			}
		}

		var msg any = head
		if nhs.convertToEthHeader {
			msg = convertBatchHeader(head)
		}

		nhs.notifiersMutex.RLock()
		defer nhs.notifiersMutex.RUnlock()

		// for each new head, notify all registered subscriptions
		for id, notifier := range nhs.newHeadNotifiers {
			if nhs.stopped.Load() {
				return nil
			}
			err := notifier.Notify(id, msg)
			if err != nil {
				// on error, remove the notification
				nhs.logger.Info("failed to notify newHead subscription", log.ErrKey, err, log.SubIDKey, id)
				nhs.notifiersMutex.Lock()
				delete(nhs.newHeadNotifiers, id)
				nhs.notifiersMutex.Unlock()
			}
		}
		return nil
	})
	return nil
}

func (nhs *NewHeadsService) RegisterNotifier(notifier *rpc.Notifier, subscription *rpc.Subscription) {
	nhs.notifiersMutex.Lock()
	defer nhs.notifiersMutex.Unlock()
	nhs.newHeadNotifiers[subscription.ID] = notifier

	go HandleUnsubscribe(subscription, nil, func() {
		nhs.notifiersMutex.Lock()
		defer nhs.notifiersMutex.Unlock()
		delete(nhs.newHeadNotifiers, subscription.ID)
	})
}

func (nhs *NewHeadsService) Stop() error {
	nhs.stopped.Store(true)
	return nil
}

func (nhs *NewHeadsService) HealthStatus() host.HealthStatus {
	return &host.BasicErrHealthStatus{}
}

func convertBatchHeader(head *common.BatchHeader) *types.Header {
	return &types.Header{
		ParentHash:  head.ParentHash,
		UncleHash:   gethcommon.Hash{},
		Root:        head.Root,
		TxHash:      head.TxHash,
		ReceiptHash: head.ReceiptHash,
		Bloom:       types.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      head.Number,
		GasLimit:    head.GasLimit,
		GasUsed:     head.GasUsed,
		Time:        head.Time,
		Extra:       make([]byte, 0),
		BaseFee:     head.BaseFee,
	}
}
