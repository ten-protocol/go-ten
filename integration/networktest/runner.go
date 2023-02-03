package networktest

import (
	"fmt"
	"testing"
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
func Run(t *testing.T, env Environment, test NetworkTest) {
	EnsureTestLogsSetUp(test.Name())
	network, envCleanup, err := env.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	defer envCleanup()
	fmt.Println("Preparing to run test:", test.Name())
	err = test.Run(network)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Test succeeded:", test.Name())
}
