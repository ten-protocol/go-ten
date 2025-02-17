package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
)

type recordNewHeadsSubscriptionAction struct {
	duration    time.Duration
	gatewayUser int // -1 if not using gateway, else test user index to get gateway from

	recordedHeads []*types.Header
}

func (r *recordNewHeadsSubscriptionAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	// get gateway address for first user
	user, err := FetchTestUser(ctx, r.gatewayUser)
	if err != nil {
		return ctx, err
	}
	// verify user is a gateway user
	gwUser, ok := user.(*userwallet.GatewayUser)
	if !ok {
		return ctx, fmt.Errorf("user is not a gateway user")
	}
	ethClient, err := gwUser.WSClient()
	if err != nil {
		return ctx, err
	}
	headsCh := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(ctx, headsCh)
	if err != nil {
		return nil, err
	}
	startTime := time.Now()
	fmt.Println("Listening for new heads")
	// read from headsCh for duration or until subscription is closed
	for time.Since(startTime) < r.duration {
		select {
		case head := <-headsCh:
			// read and store head from headsCh, then continue listening if duration has not expired
			fmt.Printf("Received new head: %v\n", head.Number)
			r.recordedHeads = append(r.recordedHeads, head)
		case <-time.After(500 * time.Millisecond):
			// no new head received, continue listening if duration has not expired
		case <-sub.Err():
			// subscription closed
			return ctx, fmt.Errorf("subscription closed unexpectedly")
		case <-ctx.Done():
			sub.Unsubscribe()
			return ctx, fmt.Errorf("context cancelled")
		}
	}
	sub.Unsubscribe()
	return ctx, nil
}

func (r *recordNewHeadsSubscriptionAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	if len(r.recordedHeads) == 0 {
		return fmt.Errorf("no new heads received during the %s period", r.duration)
	}
	return nil
}

func RecordNewHeadsSubscription(duration time.Duration) networktest.Action {
	// for now this test expects a gateway user and tests via the gateway
	// todo: add support for testing without a gateway (need to add newHeads subscription to ObsClient)
	return &recordNewHeadsSubscriptionAction{duration: duration, gatewayUser: 0}
}
