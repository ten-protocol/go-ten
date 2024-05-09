package subscription

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// this test goes via the gateway only for now
func TestGatewayNewHeadsSubscription(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gateway-new-heads-subscription",
		t,
		env.UATTestnet(), //env.LocalDevNetwork(devnetwork.WithGateway()),
		actions.Series(
			// user not technically needed, but we need a gateway address to use
			&actions.CreateTestUser{UserID: 0, UseGateway: true},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 1),

			// Record new heads for 30 seconds, verify that the subscription is working
			actions.RecordNewHeadsSubscription(30*time.Second),
		),
	)
}
