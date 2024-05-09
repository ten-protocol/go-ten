package helpful

import (
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// Smoke tests are useful for checking a network is live or checking basic functionality is not broken

var _transferAmount = big.NewInt(100_000_000)

func TestExecuteNativeFundsTransfer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"native-funds-smoketest",
		t,
		env.UATTestnet(),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(_transferAmount)},
		),
	)
}
