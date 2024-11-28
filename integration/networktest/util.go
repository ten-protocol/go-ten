package networktest

import (
	"errors"
	"os"
	"testing"

	"github.com/ten-protocol/go-ten/go/obsclient"
)

// IDEFlag is used as an environnment variable to allow tests to run that are designed not to run in CI
// todo (@matt) - come up with a better method, perhaps using directory-based ignore/include mechanism for `ci` dir only
const IDEFlag = "IDE"

func TestOnlyRunsInIDE(t *testing.T) {
	// test is skipped by default to avoid breaking CI - set env flag in `Run Configurations` to run it in IDE
	if os.Getenv(IDEFlag) == "" {
		t.Skipf("set flag %s to run this test in the IDE", IDEFlag)
	}
}

func NodeHealthCheck(rpcAddress string) error {
	client, err := obsclient.Dial(rpcAddress)
	if err != nil {
		return err
	}
	health, err := client.Health()
	if err != nil {
		return err
	}
	if !health.OverallHealth {
		return errors.New("node health check failed")
	}
	return nil
}
