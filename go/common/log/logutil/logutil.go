package logutil

import (
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/common/log"
)

type TestLogCfg struct {
	LogFile     *os.File // the log file to use, or nil if a new log file should be created
	LogDir      string   // directory for the new log file
	TestType    string   // type of test for the new log file (comes before timestamp in filename so sorted file list will block these together)
	TestSubtype string   // test subtype for the new log file (comes after timestamp in filename so sorted file list will show latest of different subtypes together)
}

// SetupTestLog will direct logs to a timestamped log file with a standard naming pattern, useful for simulations etc.
func SetupTestLog(cfg *TestLogCfg) *os.File {
	if cfg.LogFile != nil {
		log.OutputToFile(cfg.LogFile)
		return cfg.LogFile
	}

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
