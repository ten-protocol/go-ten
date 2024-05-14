package subscription

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

func TestGatewayNewHeadsSubscription(t *testing.T) {
	// networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gateway-new-heads-subscription",
		t,
		env.LongRunningLocalNetwork("http://127.0.0.1:3000"),
		actions.Series(
			// user not technically needed, but we need a gateway address to use
			&actions.CreateTestUser{UserID: 0, UseGateway: true},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 1),

			// Record new heads for specified duration, verify that the subscription is working
			actions.RecordNewHeadsSubscription(500*time.Second),
		),
	)
}
