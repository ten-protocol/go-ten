package networktest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type contextKey int

const (
	LogFileKey contextKey = 0
)

// Run provides a standardised way to run tests and provides a single place for changing logging/output styles, etc.
//
// The tests in `/tests` should typically only contain a single line, executing this method.
// The Environment and NetworkTest implementations and how they're configured define the test to be run.
//
// Example usage:
//
//	networktest.Run(t, env.DevTestnet(), tests.smokeTest())
//	networktest.Run(t, env.LocalDevNetwork(WithNumValidators(8)), traffic.RunnerTest(traffic.NativeFundsTransfers(), 30*time.Second)
func Run(testName string, t *testing.T, env Environment, action Action) {
	logFile := EnsureTestLogsSetUp(testName)
	network, envCleanup, err := env.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	initialCtx, cancelCtx := context.WithCancel(context.Background())
	ctx := context.WithValue(initialCtx, LogFileKey, logFile)
	defer func() {
		envCleanup()
		cancelCtx()
	}()
	fmt.Println("Started test:", testName)
	ctx, err = action.Run(ctx, network)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second) // allow time for latest test transactions to propagate todo (@matt) make network speeds readable from env to configure this
	fmt.Println("Verifying test:", testName)
	err = action.Verify(ctx, network)
	if err != nil {
		t.Fatal(err)
	}
}
