package load

import (
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
	"github.com/obscuronet/go-obscuro/integration/networktest/traffic"
)

func TestNativeTransfers(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(t, env.LocalDevNetwork(), traffic.DurationTest(traffic.NativeFundsTransfers(), 30*time.Second))
}

func TestERC20Transfers(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(t, env.LocalDevNetwork(), traffic.DurationTest(traffic.ERC20Transfers(), 30*time.Second))
}
