package load

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

func TestNativeTransfers(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"native-transfers-load-test",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			actions.CreateAndFundTestUsers(25),
			actions.GenerateUsersRandomisedTransferActionsInParallel(10, 60*time.Second),

			actions.VerifyUserBalancesSanity(),
		),
	)
}

//func TestERC20Transfers(t *testing.T) {
//	networktest.TestOnlyRunsInIDE(t)
//	networktest.Run(t, env.LocalDevNetwork(), traffic.DurationTest(traffic.ERC20Transfers(), 30*time.Second))
//}
