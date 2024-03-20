package clientapi

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/responses"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// FilterAPI exposes a subset of Geth's PublicFilterAPI operations.
type FilterAPI struct {
	host   host.Host
	logger gethlog.Logger
}

func NewFilterAPI(host host.Host, logger gethlog.Logger) *FilterAPI {
	return &FilterAPI{
		host:   host,
		logger: logger,
	}
}

// Logs returns a log subscription.
func (api *FilterAPI) Logs(ctx context.Context, encryptedParams common.EncryptedParamsLogSubscription) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, fmt.Errorf("creation of subscriptions is not supported")
	}
	subscription := notifier.CreateSubscription()

	logsFromSubscription := make(chan []byte)
	err := api.host.Subscribe(subscription.ID, encryptedParams, logsFromSubscription)
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
		api.host.Unsubscribe(subscription.ID)
		return nil, fmt.Errorf("could not send subscription ID to client on subscription %s", subscription.ID)
	}

	var unsubscribed atomic.Bool

	go func() {
		// to avoid unsubscribe deadlocks we have a 10 second delay between the unsubscribe command
		// and the moment we stop listening for messages
		for {
			select {
			case encryptedLog, ok := <-logsFromSubscription:
				if !ok {
					api.logger.Info("subscription channel closed", log.SubIDKey, subscription.ID)
					return
				}
				if unsubscribed.Load() {
					api.logger.Debug("subscription unsubscribed", log.SubIDKey, subscription.ID)
					return
				}
				idAndEncLog := common.IDAndEncLog{
					SubID:  subscription.ID,
					EncLog: encryptedLog,
				}
				err = notifier.Notify(subscription.ID, idAndEncLog)
				if err != nil {
					api.logger.Error("could not send encrypted log to client on subscription ", log.SubIDKey, subscription.ID)
				}
			case <-time.After(10 * time.Second):
				if unsubscribed.Load() {
					return
				}
			}
		}
	}()

	// unsubscribe commands are handled in a different go-routine to avoid deadlocking with the log processing
	go func() {
		<-subscription.Err()
		api.host.Unsubscribe(subscription.ID)
		unsubscribed.Store(true)
	}()

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
