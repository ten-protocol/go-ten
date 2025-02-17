package ci

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// TestSimulation spins up a local network and runs some activity on it, verifying the results.
// This is useful for testing the network in a more realistic scenario, we will extend this to include more complex tests.
func TestSimulation(t *testing.T) {
	networktest.Run(
		"ci-simulation-test",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			actions.CreateAndFundTestUsers(25),
			actions.GenerateUsersRandomisedTransferActionsInParallel(5, 10*time.Second),
			actions.VerifyUserBalancesSanity(),
		),
	)
}
