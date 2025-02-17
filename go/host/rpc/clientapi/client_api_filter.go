package clientapi

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	subscriptioncommon "github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// FilterAPI exposes a subset of Geth's PublicFilterAPI operations.
type FilterAPI struct {
	host            host.Host
	logger          gethlog.Logger
	NewHeadsService *subscriptioncommon.NewHeadsService
}

func NewFilterAPI(host host.Host, logger gethlog.Logger) *FilterAPI {
	return &FilterAPI{
		host:   host,
		logger: logger,
		NewHeadsService: subscriptioncommon.NewNewHeadsService(
			func() (chan *common.BatchHeader, <-chan error, error) {
				return host.NewHeadsChan(), nil, nil
			},
			false,
			logger,
			nil,
		),
	}
}

func (api *FilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()
	api.NewHeadsService.RegisterNotifier(notifier, subscription)
	return subscription, nil
}

// Logs exposes the "logs" rpc endpoint.
func (api *FilterAPI) Logs(ctx context.Context, encryptedParams common.EncryptedParamsLogSubscription) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()

	logsFromSubscription := make(chan []byte)
	err := api.host.SubscribeLogs(subscription.ID, encryptedParams, logsFromSubscription)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}

	var unsubscribed atomic.Bool
	go subscriptioncommon.ForwardFromChannels(
		[]chan []byte{logsFromSubscription},
		func(elem []byte) error {
			return notifier.Notify(subscription.ID, elem)
		},
		nil,
		nil,
		&unsubscribed,
		12*time.Hour,
		api.logger,
	)
	go subscriptioncommon.HandleUnsubscribe(subscription, func() {
		// first exit the forwarding go-routine
		unsubscribed.Store(true)
		time.Sleep(100 * time.Millisecond)
		// and then close the channel
		api.host.UnsubscribeLogs(subscription.ID)
	})
	return subscription, nil
}
