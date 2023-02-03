package ci

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
	"github.com/obscuronet/go-obscuro/integration/networktest/traffic"
)

// TestFullNetworkSim reimplements the test case of our current "sim" tests (the ERC20 transaction injector)
func TestFullNetworkSim(t *testing.T) {
	// todo: once this test is fully functional, remove this line so the test runs regularly in CI builds
	networktest.TestOnlyRunsInIDE(t)
	// todo: this is a placeholder, need to add more validation etc. to achieve parity with the existing full network sim
	networktest.Run(t, env.LocalDevNetwork(), traffic.DurationTest(traffic.ERC20Transfers(), 30*time.Second))
}
