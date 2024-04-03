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
		NewHeadsService: subscriptioncommon.NewNewHeadsService(host.NewHeadsChan(), false, logger),
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

	// We send the ID of the newly-created subscription, before sending any log events. This is because the wallet
	// extension needs to return the subscription ID to the end client, but this information is not exposed to it
	// (since the subscription ID is automatically converted to a subscription object).
	err = notifier.Notify(subscription.ID, common.IDAndEncLog{
		SubID: subscription.ID,
	})
	if err != nil {
		api.host.UnsubscribeLogs(subscription.ID)
		return nil, fmt.Errorf("could not send subscription ID to client on subscription %s", subscription.ID)
	}

	var unsubscribed atomic.Bool
	go subscriptioncommon.ForwardFromChannels([]chan []byte{logsFromSubscription}, &unsubscribed, func(elem []byte) error {
		msg := &common.IDAndEncLog{
			SubID:  subscription.ID,
			EncLog: elem,
		}
		return notifier.Notify(subscription.ID, msg)
	})
	go subscriptioncommon.HandleUnsubscribe(subscription, &unsubscribed, func() {
		api.host.UnsubscribeLogs(subscription.ID)
	})
	return subscription, nil
}

// GetLogs returns the logs matching the filter.
func (api *FilterAPI) GetLogs(_ context.Context, encryptedParams common.EncryptedParamsGetLogs) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetLogs(encryptedParams)
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
