package actions

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/integration/networktest"
)

// standard snapshots to use as reference points across all tests
var (
	// SnapAfterAllocation used to record state after initialisation and faucet allocations to test users
	SnapAfterAllocation = "after-allocation"
)

// Snapshots are used for recording data at a point in the test with a string label to describe that stage
// (storing the data into the context or the snapshot struct for usage or verification later on)

func BalanceSnapshotKey(userID int, snapshot string) ActionKey {
	return ActionKey(fmt.Sprintf("bal-%s-%d", snapshot, userID))
}

func FetchBalanceAtSnapshot(ctx context.Context, userID int, snapshot string) (*big.Int, error) {
	bal, err := FetchBigInt(ctx, BalanceSnapshotKey(userID, snapshot))
	if err != nil {
		return nil, err
	}
	return bal, nil
}

// SnapshotUserBalances requests and records the curr users native balances in the context
// Note: when running this ensure that there are no transactions in flight if that will affect usage of this data
func SnapshotUserBalances(snapshot string) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		numUsers, err := FetchNumberOfTestUsers(ctx)
		if err != nil {
			return nil, err
		}
		for i := 0; i < numUsers; i++ {
			user, err := FetchTestUser(ctx, i)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch user %d - %w", i, err)
			}
			bal, err := user.NativeBalance(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch balance for user %d - %w", i, err)
			}
			ctx = context.WithValue(ctx, BalanceSnapshotKey(i, snapshot), bal)
		}
		return ctx, nil
	})
}
