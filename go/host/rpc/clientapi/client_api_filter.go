package clientapi

import (
	"context"
	"fmt"
	"time"

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

func NewFilterAPI(host host.Host, logsCh chan []*types.Log) *FilterAPI {
	return &FilterAPI{
		host:          host,
		gethFilterAPI: filters.NewPublicFilterAPI(events.NewBackend(logsCh), false, 5*time.Minute),
	}
}

// Logs returns a log subscription.
func (api *FilterAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	err := api.host.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}
	// TODO - #453 - Manage termination of host -> enclave subscriptions when client subscriptions are terminated.
	return api.gethFilterAPI.Logs(ctx, crit)
}
