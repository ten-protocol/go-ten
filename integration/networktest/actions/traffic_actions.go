package actions

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
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
