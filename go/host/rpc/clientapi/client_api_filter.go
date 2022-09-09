package clientapi

import (
	"context"
	"fmt"
	"time"

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

func NewFilterAPI(host host.Host, logsCh chan []*types.Log) *FilterAPI {
	return &FilterAPI{
		host:          host,
		gethFilterAPI: filters.NewPublicFilterAPI(events.NewBackend(logsCh), false, 5*time.Minute),
	}
}

// Logs returns a log subscription.
// TODO - #453 - Handle termination of the corresponding host -> enclave subscription when a client subscription is
//  terminated. This is non-trivial, as since we are leveraging Geth's subscription framework, we have no way of
//  determining when an unsubscribe request is made. The cleanest option appears to be detecting the `eth_unsubscribe`
//  request in the wallet extension, and firing an additional unsubscribe call to the host.
func (api *FilterAPI) Logs(ctx context.Context, encryptedLogSubscription common.EncryptedLogSubscription) (*rpc.Subscription, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("could not generate new UUID for subscription. Cause: %w", err)
	}

	err = api.host.Subscribe(id, encryptedLogSubscription)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}

	return api.gethFilterAPI.Logs(ctx, emptyFilter)
}
