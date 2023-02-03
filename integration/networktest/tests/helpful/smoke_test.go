package helpful

import (
	"math/big"
	"testing"

	"github.com/obscuronet/go-obscuro/integration/networktest/actions"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

// Smoke tests are useful for checking a network is live or checking basic functionality is not broken

var _transferAmount = big.NewInt(100_000_000)

func TestExecuteNativeFundsTransfer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"native-funds-smoketest",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			&actions.AllocateFaucetFunds{UserID: 0},
			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
		),
	)
}
