package log

import (
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common/log"
)

type TestLogCfg struct {
	LogDir      string // directory for the log file
	TestType    string // type of test (comes before timestamp in filename so sorted file list will block these together)
	TestSubtype string // test subtype (comes after timestamp in filename so sorted file list will show latest of different subtypes together)
}

// SetupTestLog will direct logs to a timestamped log file with a standard naming pattern, useful for simulations etc.
func SetupTestLog(cfg *TestLogCfg) *os.File {
	// Create a folder for the logs if none exists.
	err := os.MkdirAll(cfg.LogDir, 0o700)
	if err != nil {
		panic(err)
	}
	timeFormatted := time.Now().Format("2006-01-02_15-04-05")
	f, err := os.CreateTemp(cfg.LogDir, fmt.Sprintf("%s-%s-%s-*.txt", cfg.TestType, timeFormatted, cfg.TestSubtype))
	if err != nil {
		panic(err)
	}
	log.OutputToFile(f)
	return f
}
