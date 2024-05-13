package subscription

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

func TestGatewayNewHeadsSubscription(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gateway-new-heads-subscription",
		t,
		env.LocalDevNetwork(devnetwork.WithGateway()),
		actions.Series(
			// user not technically needed, but we need a gateway address to use
			&actions.CreateTestUser{UserID: 0, UseGateway: true},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 1),

			// Record new heads for specified duration, verify that the subscription is working
			actions.RecordNewHeadsSubscription(20*time.Second),
		),
	)
}
