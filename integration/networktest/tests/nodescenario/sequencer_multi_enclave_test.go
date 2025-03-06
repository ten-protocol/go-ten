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
// Note: this is a happy path failover, we need to test for edge cases etc and test the failover in a live testnet
func TestHASequencerBackup(t *testing.T) {
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

			// wait for tx to complete before killing
			actions.SleepAction(5*time.Second),

			// kill the primary enclave
			actions.StopSequencerEnclave(0),

			// wait for failover to complete
			actions.SleepAction(5*time.Second),

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(3*time.Second),

			// two transfers should have happened so we verify double the amounts
			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: doubleTransferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(doubleTransferAmount)},
		),
	)
}

// TestHARestoringEnclaves tests that the enclaves can be brought back after being stopped and they continue to work
// The order for this test is:
// - tx when both encl up
// - stop encl 0 then tx
// - start encl 0 then tx
// - stop encl 1 then tx
// thus checking that encl 0 which was stopped successfully re-entered the pool of enclaves and was promoted again
func TestHARestoringEnclaves(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	totalTransferAmount := big.NewInt(1).Mul(big.NewInt(4), _transferAmount)
	networktest.Run(
		"ha-sequencer-enclaves-restored",
		t,
		env.LocalDevNetwork(devnetwork.WithHASequencer()),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete before killing
			actions.SleepAction(5*time.Second),

			// kill the primary enclave
			actions.StopSequencerEnclave(0),

			// wait for failover to complete
			actions.SleepAction(5*time.Second),

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(5*time.Second),

			// start the enclave again
			actions.StartSequencerEnclave(0),

			// wait for the enclave to start
			actions.SleepAction(5*time.Second),

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(5*time.Second),

			// stop the other enclave
			actions.StopSequencerEnclave(1),

			// wait for failover to complete
			actions.SleepAction(5*time.Second),

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			// wait for tx to complete
			actions.SleepAction(15*time.Second),

			// two transfers should have happened so we verify double the amounts
			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: totalTransferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(totalTransferAmount)},
		),
	)
}
