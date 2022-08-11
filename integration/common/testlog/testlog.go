package testlog

import (
	"fmt"
	"os"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"
)

type Cfg struct {
	LogDir      string // directory for the log file
	TestType    string // type of test (comes before timestamp in filename so sorted file list will block these together)
	TestSubtype string // test subtype (comes after timestamp in filename so sorted file list will show latest of different subtypes together)
}

// Setup will direct logs to a timestamped log file with a standard naming pattern, useful for simulations etc.
func Setup(cfg *Cfg) *os.File {
	// hardcode geth log level to error only
	gethlog.Root().SetHandler(gethlog.LvlFilterHandler(gethlog.LvlError, gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true))))

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
