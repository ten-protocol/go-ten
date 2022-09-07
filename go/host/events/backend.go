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
type Backend struct {
	logsCh <-chan []*types.Log
}

func NewBackend(logsCh <-chan []*types.Log) Backend {
	return Backend{
		logsCh: logsCh,
	}
}

func (b Backend) ChainDb() ethdb.Database { //nolint:stylecheck
	return nil // Not implemented.
}

func (b Backend) HeaderByNumber(context.Context, rpc.BlockNumber) (*types.Header, error) {
	return nil, errNotSupported
}

func (b Backend) HeaderByHash(context.Context, common.Hash) (*types.Header, error) {
	return nil, errNotSupported
}

func (b Backend) GetReceipts(context.Context, common.Hash) (types.Receipts, error) {
	return nil, errNotSupported
}

func (b Backend) GetLogs(context.Context, common.Hash) ([][]*types.Log, error) {
	return nil, errNotSupported
}

func (b Backend) SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribeChainEvent(chan<- core.ChainEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

// TODO - #453 - Handle removed logs.
func (b Backend) SubscribeRemovedLogsEvent(chan<- core.RemovedLogsEvent) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	logsProducer := func(<-chan struct{}) error {
		for {
			logs := <-b.logsCh
			ch <- logs
		}
	}

	return event.NewSubscription(logsProducer)
}

// TODO - #453 - Handle pending logs.
func (b Backend) SubscribePendingLogsEvent(chan<- []*types.Log) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) BloomStatus() (uint64, uint64) {
	return 0, 0 // Not implemented.
}

func (b Backend) ServiceFilter(context.Context, *bloombits.MatcherSession) {
	// Not implemented.
}
