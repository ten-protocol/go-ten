package nodescenario

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest/actions"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

// restart both the sequencer and the validators (the entire network)
func TestRestartNetwork(t *testing.T) {
	//networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"restart-network",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			actions.CreateAndFundTestUsers(5),

			// short load test, build up some state
			actions.GenerateUsersRandomisedTransferActionsInParallel(4, 20*time.Second),

			// todo: this could be replaced by something that finds all the transaction IDs in context and waits for them to be mined
			actions.SleepAction(5*time.Second), // allow time for in-flight transactions

			// stop sequencer and validator
			actions.StopSequencerHost(),
			actions.StopSequencerEnclave(),
			actions.StopValidatorHost(0),
			actions.StopValidatorEnclave(0),
			actions.StopValidatorHost(1),
			actions.StopValidatorEnclave(1),
			actions.StopValidatorHost(2),
			actions.StopValidatorEnclave(2),

			actions.SleepAction(60*time.Second), // allow time for shutdowns, allow L1 to get a bit ahead

			// start sequencer and validator
			actions.StartValidatorEnclave(0),
			actions.StartValidatorHost(0),
			actions.StartValidatorEnclave(1),
			actions.StartValidatorHost(1),
			actions.StartValidatorEnclave(2),
			actions.StartValidatorHost(2),
			actions.StartSequencerEnclave(),
			actions.StartSequencerHost(),
			actions.WaitForValidatorHealthCheck(0, 30*time.Second),
			actions.WaitForValidatorHealthCheck(1, 30*time.Second),
			actions.WaitForValidatorHealthCheck(2, 30*time.Second),
			actions.WaitForSequencerHealthCheck(30*time.Second),

			// todo: we often see 1 transaction getting lost without this sleep after the node restarts.
			// 	This needs investigating but it suggests to me that the health check is succeeding prematurely
			actions.SleepAction(5*time.Second), // allow time for re-sync

			// resubmit user viewing keys (all users will have lost their "session")
			// todo: get rid of this once the enclave persists viewing keys correctly
			actions.AuthenticateAllUsers(),

			// another load test, check that the network is still working
			actions.GenerateUsersRandomisedTransferActionsInParallel(4, 60*time.Second),
		),
	)
}
