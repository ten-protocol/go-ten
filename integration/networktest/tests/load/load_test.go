package load

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/actions"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

func TestNativeTransfers(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"native-transfers-load-test",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			actions.CreateAndFundTestUsers(2),
			actions.GenerateUsersRandomisedTransferActionsInParallel(2, 10*time.Second),

			actions.VerifyUserBalancesSanity(),
		),
	)
}

//func TestERC20Transfers(t *testing.T) {
//	networktest.TestOnlyRunsInIDE(t)
//	networktest.Run(t, env.LocalDevNetwork(), traffic.DurationTest(traffic.ERC20Transfers(), 30*time.Second))
//}
