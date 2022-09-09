package clientapi

import (
	"context"
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/google/uuid"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/host/events"
)

// Filters out nothing. The filtering is performed on the enclave side instead, based on the filters in the
// `common.EncryptedLogSubscription`.
var emptyFilter = filters.FilterCriteria{}

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
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("could not generate new UUID for subscription. Cause: %w", err)
	}

	err = api.host.Subscribe(id, encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}

	// TODO - #453 - The empty filter here is incorrect. The API retrieves all logs, then filters them based on the
	//  filter. If we pass an empty filter, we'll get back all the logs for all the subscription IDs. We could always
	//  filter this on the wallet extension side, but we could also use custom topics in the filter (e.g. a topic that
	//  was the subscription ID).
	subscription, err := api.gethFilterAPI.Logs(ctx, emptyFilter)
	if err != nil {
		return nil, err
	}

	go func() {
		<-subscription.Err() // This channel's sole purpose is to be closed when the subscription is unsubscribed.
		err = api.host.Unsubscribe(id)
		if err != nil {
			log.Error("could not unsubscribe from subscription %s", id)
		}
	}()

	return subscription, nil
}
