package networktest

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
)

// EnsureTestLogsSetUp calls Setup if it hasn't already been called (some tests run tests within themselves, we don't want
// the log folder flipping around for every subtest, so we assume this is called for the top level test that is running
// and ignore subsequent calls
func EnsureTestLogsSetUp(testName string) *os.File {
	logger := testlog.Logger()
	if logger != nil {
		return nil // already setup, do not reconfigure
	}
	return testlog.Setup(&testlog.Cfg{
		// todo (@matt) - walk up the dir tree to find /integration/.build or find best practice solution
		// bit of a hack - tests need to be in a package nested within /tests to get logs in the right place
		LogDir:      "../../../.build/networktest/",
		TestType:    "net",
		TestSubtype: testName,
		LogLevel:    log.LvlInfo,
	})
}
