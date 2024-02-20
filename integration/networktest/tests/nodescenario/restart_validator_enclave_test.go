package nodescenario

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest/actions"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

func TestRestartValidatorEnclave(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"restart-enclave",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			actions.CreateAndFundTestUsers(5),

			// short load test, build up some state
			actions.GenerateUsersRandomisedTransferActionsInParallel(4, 10*time.Second),

			// restart enclave on a validator
			actions.StopValidatorEnclave(1),
			actions.SleepAction(5*time.Second), // allow time for shutdown
			actions.StartValidatorEnclave(1),
			actions.WaitForValidatorHealthCheck(1, 30*time.Second),

			// todo (@matt) - we often see 1 transaction getting lost without this sleep after the node restarts.
			// 	This needs investigating but it suggests to me that the health check is succeeding prematurely
			actions.SleepAction(5*time.Second), // allow time for re-sync

			// another load test (important that at least one of the users will be using the validator with restarted enclave)
			actions.GenerateUsersRandomisedTransferActionsInParallel(4, 10*time.Second),
		),
	)
}
