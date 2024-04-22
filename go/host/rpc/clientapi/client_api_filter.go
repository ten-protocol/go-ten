package clientapi

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ten-protocol/go-ten/go/common/host"
	subscriptioncommon "github.com/ten-protocol/go-ten/go/common/subscription"
	"github.com/ten-protocol/go-ten/go/responses"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common"

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
		host:            host,
		logger:          logger,
		NewHeadsService: subscriptioncommon.NewNewHeadsService(host.NewHeadsChan(), false, logger, nil),
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
	go subscriptioncommon.ForwardFromChannels([]chan []byte{logsFromSubscription}, &unsubscribed, func(elem []byte) error {
		return notifier.Notify(subscription.ID, elem)
	})
	go subscriptioncommon.HandleUnsubscribe(subscription, &unsubscribed, func() {
		api.host.UnsubscribeLogs(subscription.ID)
	})
	return subscription, nil
}

// GetLogs returns the logs matching the filter.
func (api *FilterAPI) GetLogs(ctx context.Context, encryptedParams common.EncryptedParamsGetLogs) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetLogs(ctx, encryptedParams)
	if sysError != nil {
		return api.handleSysError("GetLogs", sysError)
	}
	return *enclaveResponse, nil
}

func (api *FilterAPI) handleSysError(function string, sysError common.SystemError) (responses.EnclaveResponse, error) {
	api.logger.Error(fmt.Sprintf("Enclave System Error. Function %s", function), log.ErrKey, sysError)
	return responses.EnclaveResponse{
		Err: &responses.InternalErrMsg,
	}, nil
}
