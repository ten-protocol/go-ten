package events

import (
	"context"

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
	return nil
}

func (b Backend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	return nil, nil //nolint:nilnil
}

func (b Backend) HeaderByHash(ctx context.Context, blockHash common.Hash) (*types.Header, error) {
	return nil, nil //nolint:nilnil
}

func (b Backend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return nil, nil
}

func (b Backend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	return nil, nil
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
	return event.NewSubscription(nilProducer)
}

func (b Backend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return event.NewSubscription(nilProducer)
}

func (b Backend) BloomStatus() (uint64, uint64) {
	return 0, 0
}
func (b Backend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {}
