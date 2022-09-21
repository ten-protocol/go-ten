package clientapi

import (
	"context"
	"fmt"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/google/uuid"

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
	subscriptionID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("could not generate new UUID for subscription. Cause: %w", err)
	}

	err = api.host.Subscribe(subscriptionID, encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe for logs. Cause: %w", err)
	}

	subscriptionIDBytes, err := subscriptionID.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("could not marshal subscription ID to bytes. Cause: %w", err)
	}

	// The "real" filtering is performed on the enclave side instead, based on the filters in the
	// `common.EncryptedLogSubscription`. Here, we simply filter using a padded subscription ID. This allows us to
	// reuse Geth's subscription machinery to automatically filter out logs for other subscription IDs.
	filter := filters.FilterCriteria{
		Topics: [][]gethcommon.Hash{{gethcommon.BytesToHash(subscriptionIDBytes)}},
	}

	subscription, err := api.gethFilterAPI.Logs(ctx, filter)
	if err != nil {
		return nil, err
	}

	go func() {
		<-subscription.Err() // This channel's sole purpose is to be closed when the subscription is unsubscribed.
		err = api.host.Unsubscribe(subscriptionID)
		if err != nil {
			log.Error("could not unsubscribe from subscription %s", subscriptionID)
		}
	}()

	return subscription, nil
}
