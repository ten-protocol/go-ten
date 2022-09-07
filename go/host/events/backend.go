package events

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
)

var nilProducer = func(<-chan struct{}) error {
	select {} // As soon as this method returns, the subscription is unsubscribed.
}

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) ChainDb() ethdb.Database { //nolint:stylecheck
	panic("not implemented")
}

func (b Backend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	panic("not implemented")
}

func (b Backend) HeaderByHash(ctx context.Context, blockHash common.Hash) (*types.Header, error) {
	panic("not implemented")
}

func (b Backend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	panic("not implemented")
}

func (b Backend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	panic("not implemented")
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
	var logsProducer = func(<-chan struct{}) error {
		for {
			log := types.Log{
				Topics: []common.Hash{},
				Data:   []byte("hello world"),
			}
			ch <- []*types.Log{&log}
			println("jjj just sent some logs")
			time.Sleep(time.Second)
		} // As soon as this method returns, the subscription is unsubscribed.
	}

	return event.NewSubscription(logsProducer)
}

func (b Backend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) BloomStatus() (uint64, uint64) {
	panic("not implemented")
}

func (b Backend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	panic("not implemented")
}
