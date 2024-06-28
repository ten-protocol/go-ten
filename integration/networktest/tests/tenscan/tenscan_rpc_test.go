package tenscan

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/actions/publicdata"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

var _transferAmount = big.NewInt(100_000_000)

// Verify and debug the RPC endpoints that Tenscan relies on for data in various environments

func TestRPC(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"tenscan-rpc-data",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			publicdata.VerifyBatchesDataAction(),
		),
	)
}

// Test the personal transactions endpoint in various environments (it uses getStorageAt so it can run through MM etc.)
// 1. create user
// 2. send some transactions
// 3. verify transactions are returned by the personal transactions endpoint that tenscan uses
func TestPersonalTransactions(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"tenscan-personal-transactions",
		t,
		env.LocalDevNetwork(devnetwork.WithGateway()),
		actions.Series(
			// create 3 users
			&actions.CreateTestUser{UserID: 0, UseGateway: true}, // <-- this user makes the PersonalTransactions request, choose gateway or not here
			&actions.CreateTestUser{UserID: 1},
			&actions.CreateTestUser{UserID: 2},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 3),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			// user 0 sends funds to users 1 and 2
			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},
			&actions.SendNativeFunds{FromUser: 0, ToUser: 2, Amount: _transferAmount},

			// after the test we will verify the other users received them
			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
			&actions.VerifyBalanceAfterTest{UserID: 2, ExpectedBalance: _transferAmount},

			// verify the personal transactions endpoint returns the two txs
			actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
				user, err := actions.FetchTestUser(ctx, 0)
				if err != nil {
					return err
				}

				pagination := common.QueryPagination{
					Offset: 0,
					Size:   20,
				}
				personalTxs, total, err := user.GetPersonalTransactions(ctx, pagination)
				if err != nil {
					return fmt.Errorf("unable to get personal transactions - %w", err)
				}

				// verify the transactions
				if len(personalTxs) != 2 {
					return fmt.Errorf("expected 2 transactions, got %d", len(personalTxs))
				}

				// verify total set
				if total != 2 {
					return fmt.Errorf("expected total receipts to be at least 2, got %d", total)
				}
				return nil
			}),
		),
	)
}
