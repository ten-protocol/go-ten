package bridge

import (
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/actions"
	"github.com/obscuronet/go-obscuro/integration/networktest/actions/l1"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

func TestImportantContractsLookup(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"important-contracts-lookup",
		t,
		env.LocalDevNetwork(),
		actions.Series(
			l1.SetImportantContract("L1TestContract", gethcommon.HexToAddress("0x64")),
		),
	)
}
