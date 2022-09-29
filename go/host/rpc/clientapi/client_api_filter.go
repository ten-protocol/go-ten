package clientapi

import (
	"context"
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/host/events"
)

// FilterAPI exposes a subset of Geth's PublicFilterAPI operations.
type FilterAPI struct {
	host          host.Host
	gethFilterAPI *filters.PublicFilterAPI
}

func NewFilterAPI(host host.Host, logsCh chan *types.Log) *FilterAPI {
	return &FilterAPI{
		host:          host,
		gethFilterAPI: filters.NewPublicFilterAPI(events.NewBackend(logsCh), false, 5*time.Minute),
	}
}

// Logs returns a log subscription.
func (api *FilterAPI) Logs(ctx context.Context, encryptedParams common.EncryptedParamsLogSubscription) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		panic("jjj not supported") // todo - joel - better handling
	}
	subscription := notifier.CreateSubscription()

	matchedLogs := make(chan []byte)
	// todo - joel - feed the logs from the host to this channel based on the subscription IDs
	err := api.host.Subscribe(subscription.ID, encryptedParams, matchedLogs)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			matchedLogs <- []byte("hey joel jjj")
		}
	}()

	go func() {
		for {
			select {
			case encryptedLog := <-matchedLogs:
				println("jjj notifying")
				notifier.Notify(subscription.ID, encryptedLog)

			case <-subscription.Err(): // client send an unsubscribe request
				// todo - joel - call unsubscribe
				return

			case <-notifier.Closed(): // connection dropped
				// todo - joel - call unsubscribe
				return
			}
		}
	}()

	return subscription, nil
}
