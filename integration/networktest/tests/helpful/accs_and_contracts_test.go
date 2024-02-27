package helpful

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

var _accountToFund = common.HexToAddress("0xD19f62b5A721747A04b969C90062CBb85D4aAaA8")

// Run this test to fund an account with native funds
func TestSendFaucetFunds(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"send-faucet-funds",
		t,
		env.LongRunningLocalNetwork(""),
		&actions.AllocateFaucetFunds{Account: &_accountToFund},
	)
}
