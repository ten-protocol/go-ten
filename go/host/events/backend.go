package events

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	// We only support a limited subset of the `PublicFilterAPI` APIs.
	errNotSupported = errors.New("this operation is not supported")

	// Loops indefinitely. Cannot return, as this causes the subscription to be unsubscribed and the processing of
	// events in `EventSystem.eventLoop` to end.
	nilProducer = func(<-chan struct{}) error {
		for {
			time.Sleep(time.Minute)
		}
	}
)

// Backend is a custom backend for Geth's `PublicFilterAPI`.
type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) ChainDb() ethdb.Database { //nolint:stylecheck
	return nil // Not implemented.
}

func (b Backend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	return nil, errNotSupported
}

func (b Backend) HeaderByHash(ctx context.Context, blockHash common.Hash) (*types.Header, error) {
	return nil, errNotSupported
}

func (b Backend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return nil, errNotSupported
}

func (b Backend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	return nil, errNotSupported
}

func (b Backend) SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	logsProducer := func(<-chan struct{}) error {
		for {
			// todo - joel - read logs from channel passed to backend
			log := types.Log{
				Topics: []common.Hash{},
				Data:   []byte("hello world"),
			}
			ch <- []*types.Log{&log}
			time.Sleep(time.Second)
		} // As soon as this method returns, the subscription is unsubscribed.
	}

	return event.NewSubscription(logsProducer)
}

func (b Backend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) BloomStatus() (uint64, uint64) {
	return 0, 0 // Not implemented.
}

func (b Backend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	// Not implemented.
}
