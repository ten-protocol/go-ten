package tenscan

import (
	"testing"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/actions/publicdata"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

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
