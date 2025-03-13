package bridge

import (
	"testing"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/actions/l1"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

func TestImportantContractsLookup(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"important-contracts-lookup",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			// TODO we no longer have this function but may want to add a test around the new functionality
			// l1.SetImportantContract("L1TestContract", gethcommon.HexToAddress("0x64")),
			// Verify that the L2 Message Bus address is made available by the host (it is deployed with a synthetic tx)
			l1.VerifyL2MessageBusAddressAvailable(),
		),
	)
}
