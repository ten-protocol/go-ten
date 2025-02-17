package helpful

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

// Smoke tests are useful for checking a network is live or checking basic functionality is not broken

var _transferAmount = big.NewInt(100_000_000)

// Transaction with insufficient gas limit for the intrinsic cost. It should result in no difference
// to user balances, but network should remain operational.
// Used to automatically detect batch desync based on transaction inclusion.
// Sequencer and Validator will process different transactions, but state should be identical.
func TestExecuteNativeFundsTransferNoGas(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gas-underlimit-test",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			&actions.CreateTestUser{UserID: 0},
			&actions.CreateTestUser{UserID: 1},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)
			&actions.SendNativeFunds{
				FromUser:   0,
				ToUser:     1,
				Amount:     _transferAmount,
				GasLimit:   big.NewInt(11_000),
				SkipVerify: true,
			},
			&actions.VerifyBalanceAfterTest{
				UserID:          1,
				ExpectedBalance: common.Big0,
			},
			actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
				logFile, ok := (ctx.Value(networktest.LogFileKey)).(*os.File)
				if !ok {
					return fmt.Errorf("log file not provided in context")
				}
				fmt.Println(logFile.Name())

				f, err := os.Open(logFile.Name())
				if err != nil {
					return err
				}

				scanner := bufio.NewScanner(f)

				// https://golang.org/pkg/bufio/#Scanner.Scan
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "Error validating batch") {
						return fmt.Errorf("found bad batches in test logs")
					}
				}

				if err := scanner.Err(); err != nil {
					// Handle the error
					return err
				}

				return nil
			}),
		),
	)
}
