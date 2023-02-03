package networktest

import (
	"errors"
	"os"
	"testing"

	"github.com/obscuronet/go-obscuro/go/obsclient"
)

// IDEFlag is used as an environnment variable to allow tests to run that are designed not to run in CI
// todo: come up with a better method, perhaps using directory-based ignore/include mechanism for `ci` dir only
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
	if !health {
		return errors.New("node health check failed")
	}
	return nil
}

type namedTest struct {
	name string
	test func(network NetworkConnector) error
}

func (n *namedTest) Run(network NetworkConnector) error {
	return n.test(network)
}

func (n *namedTest) Name() string {
	return n.name
}

// CreateTest makes it easy to create a NetworkTest in one line rather than having to make a new struct for every test
// even when it doesn't have any state or config to warrant a struct
// `name` should be short and without spaces, it is used in log filename, `test` is the function to run (fail with error)
func CreateTest(name string, test func(network NetworkConnector) error) NetworkTest {
	return &namedTest{
		name: name,
		test: test,
	}
}
