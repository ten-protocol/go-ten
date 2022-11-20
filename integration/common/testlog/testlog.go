package testlog

import (
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// must be set during setup
var (
	logFile = ""
	testlog gethlog.Logger
)

func Logger() gethlog.Logger {
	return testlog
}

func LogFile() string {
	return logFile
}

type Cfg struct {
	LogDir      string // directory for the log file
	TestType    string // type of test (comes before timestamp in filename so sorted file list will block these together)
	TestSubtype string // test subtype (comes after timestamp in filename so sorted file list will show latest of different subtypes together)
	LogLevel    gethlog.Lvl
}

// Setup will direct logs to a timestamped log file with a standard naming pattern, useful for simulations etc.
func Setup(cfg *Cfg) *os.File {
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
	logFile = f.Name()
	// hardcode geth log level to error only
	gethlog.Root().SetHandler(gethlog.LvlFilterHandler(cfg.LogLevel, gethlog.StreamHandler(f, gethlog.TerminalFormat(false))))
	testlog = gethlog.Root().New(log.CmpKey, log.TestLogCmp)
	return f
}
