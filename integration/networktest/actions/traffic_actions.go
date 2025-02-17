package actions

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

const _minTransferAmt = 1_000_000

// functions in here are used to generate actions to run in parallel or series which
// simulate "random" user activity over a period of time

type parallelFundsTransferTraffic struct {
	txPerSec int
	duration time.Duration

	// parallelFundsTransferTraffic just wraps a ParallelAction, but it doesn't initialise that until Run() is called because it
	// relies on data from the context
	parallelAction networktest.Action
}

// Run builds the parallel action series to be run before delegating its call to the parallel runner
func (p *parallelFundsTransferTraffic) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	// find number of test users available
	numUsers, err := FetchNumberOfTestUsers(ctx)
	if err != nil {
		return nil, err
	}
	avgTimeBetweenUserTx := float32(numUsers) / float32(p.txPerSec)
	txPerUser := int(float64(p.txPerSec) * p.duration.Seconds() / float64(numUsers))

	// create a series of actions for each test user that can be run in parallel
	var allUsersActionSeries []networktest.Action // series of actions per user
	for i := 0; i < numUsers; i++ {
		var userActionSeries []networktest.Action
		for j := 0; j < txPerUser; j++ {
			// sleep up to 2x avg time between tx (avgTimeBetweenUserTx is a float in seconds, important to multiply it up before casting it to duration)
			maxWait := time.Duration(float32(2000)*avgTimeBetweenUserTx) * time.Millisecond
			userActionSeries = append(userActionSeries, RandomSleepAction(time.Millisecond, maxWait),
				&SendNativeFunds{
					FromUser: i,
					ToUser:   getRandomTargetUser(numUsers, i),
					Amount:   big.NewInt(int64(common.RndBtw(_minTransferAmt, 50*_minTransferAmt))),
				})
		}
		allUsersActionSeries = append(allUsersActionSeries, NamedSeries(fmt.Sprintf("native transfers - user %d", i), userActionSeries...))
	}
	// create a parallel action for the user serieses and then delegate running to that action
	p.parallelAction = NamedParallel(p.String(), allUsersActionSeries...)
	return p.parallelAction.Run(ctx, network)
}

func (p *parallelFundsTransferTraffic) String() string {
	return fmt.Sprintf("user actions - %d TPS", p.txPerSec)
}

func getRandomTargetUser(numUsers int, fromUser int) int {
	rndIdx := rand.Intn(numUsers) //nolint:gosec
	for rndIdx == fromUser {
		rndIdx = rand.Intn(numUsers) //nolint:gosec
	}
	return rndIdx
}

func (p *parallelFundsTransferTraffic) Verify(ctx context.Context, network networktest.NetworkConnector) error {
	return p.parallelAction.Verify(ctx, network)
}

func GenerateUsersRandomisedTransferActionsInParallel(txPerSec int, duration time.Duration) networktest.Action {
	return &parallelFundsTransferTraffic{txPerSec: txPerSec, duration: duration}
}

// VerifyUserBalancesSanity expects a balances SnapAfterAllocation snapshot
// It sums up the user balances then and now to make sure that the total hasn't increased and that it's decreased but not drastically (gas fees)
func VerifyUserBalancesSanity() networktest.Action {
	return VerifyOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) error {
		numUsers, err := FetchNumberOfTestUsers(ctx)
		if err != nil {
			return fmt.Errorf("expected number of test users to be set - %w", err)
		}
		initTotalBal := big.NewInt(0)
		currTotalBal := big.NewInt(0)
		for i := 0; i < numUsers; i++ {
			initBalance, err := FetchBalanceAtSnapshot(ctx, i, SnapAfterAllocation)
			if err != nil {
				return err
			}
			initTotalBal = initTotalBal.Add(initTotalBal, initBalance)
			user, err := FetchTestUser(ctx, i)
			if err != nil {
				return err
			}
			currBal, err := user.NativeBalance(ctx)
			if err != nil {
				return err
			}
			currTotalBal = currTotalBal.Add(currTotalBal, currBal)
		}

		if currTotalBal.Cmp(initTotalBal) >= 0 {
			return fmt.Errorf("expected total balances to have slightly decreased due to gas costs, but initTotal=%d currTotal=%d", initTotalBal, currTotalBal)
		}
		// this test is a bit wooly, if it starts failing in certain use cases then reconsider it but could be a useful sanity check
		if currTotalBal.Cmp(initTotalBal.Div(initTotalBal, big.NewInt(2))) < 0 {
			return fmt.Errorf("expected total balances to still be at least half of original balances but initTotal=%d currTotal=%d", initTotalBal, currTotalBal)
		}
		return nil
	})
}
