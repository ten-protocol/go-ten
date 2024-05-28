package gateway

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

var _transferAmount = big.NewInt(100_000_000)

// TestGatewayHappyPath tests ths same functionality as the smoke_test but with the gateway:
// 1. Create two test users
// 2. Allocate funds to the first user
// 3. Send funds from the first user to the second
// 4. Verify the second user has the funds
// 5. Verify the first user has the funds deducted
// To run this test with a local network use the flag to start it with the gateway enabled.
func TestGatewayHappyPath(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gateway-happy-path",
		t,
		env.LocalDevNetwork(devnetwork.WithGateway()),
		actions.Series(
			&actions.CreateTestUser{UserID: 0, UseGateway: true},
			&actions.CreateTestUser{UserID: 1, UseGateway: true},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(_transferAmount)},

			// test net_version works through the gateway
			actions.VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
				user, err := actions.FetchTestUser(ctx, 0)
				if err != nil {
					return err
				}
				// verify user is a gateway user
				gwUser, ok := user.(*userwallet.GatewayUser)
				if !ok {
					return fmt.Errorf("user is not a gateway user")
				}
				ethClient := gwUser.Client()
				rpcClient := ethClient.Client()
				// check net_version response
				var result string
				err = rpcClient.CallContext(ctx, &result, "net_version")
				if err != nil {
					return fmt.Errorf("failed to get net_version: %w", err)
				}
				fmt.Println("net_version response:", result)
				expectedResult := "443"
				if result != expectedResult {
					return fmt.Errorf("expected net_version to be %s but got %s", expectedResult, result)
				}

				// check web3_clientVersion response
				var cvResult string
				err = rpcClient.CallContext(ctx, &cvResult, "web3_clientVersion")
				if err != nil {
					return fmt.Errorf("failed to get web3_clientVersion: %w", err)
				}
				fmt.Println("web3_clientVersion response:", cvResult)
				if cvResult == "" {
					return fmt.Errorf("expected web3_clientVersion to be non-empty")
				}

				return nil
			}),
		),
	)
}
