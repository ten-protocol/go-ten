package nodescenario

import (
	"math/big"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

var _transferAmount = big.NewInt(100_000_000)

func TestMultiEnclaveSequencer(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"multi-enclave-sequencer",
		t,
		env.LocalDevNetwork(devnetwork.WithHASequencer()),
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

// This test runs with an HA sequencer, does a transfer then kills the active sequencer enclave,
// allows it time to failover then performs another transfer to check the failover was successful.
func TestHASequencerFailover(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	doubleTransferAmount := big.NewInt(2).Mul(big.NewInt(2), _transferAmount)
	networktest.Run(
		"ha-sequencer-failover",
		t,
		env.LocalDevNetwork(devnetwork.WithHASequencer()),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(5*time.Second), // allow time for shutdown/failover

			// kill the primary enclave
			actions.StopSequencerEnclave(0),

			// wait for failover to complete
			actions.SleepAction(5*time.Second), // allow time for shutdown/failover

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(3*time.Second), // allow time for shutdown/failover

			// two transfers should have happened so we verify double the amounts
			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: doubleTransferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(doubleTransferAmount)},
		),
	)
}
