package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ten-protocol/go-ten/lib/gethfork/log"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

// These are the keys of the log entries
const (
	ErrKey           = "err"
	CtrErrKey        = "ctr_err"
	SubIDKey         = "subscription_id"
	CfgKey           = "cfg"
	TxKey            = "tx"
	DurationKey      = "duration"
	DurationMilliKey = "durationMilli"
	BundleHashKey    = "bundle"
	BatchHashKey     = "batch"
	BatchHeightKey   = "batch_height"
	BatchSeqNoKey    = "batch_seq_num"
	RollupHashKey    = "rollup"
	CmpKey           = "cmp"
	NodeIDKey        = "node_id"
	EnclaveIDKey     = "enclave_id"
	NetworkIDKey     = "network_id"
	BlockHeightKey   = "block_height"
	BlockHashKey     = "block_hash"
	PackageKey       = "package"
)

// Logging is grouped by the component where it was initialised
const (
	EnclaveCmp      = "enclave"
	HostCmp         = "host"
	HostRPCCmp      = "host_rpc"
	TxInjectCmp     = "tx_inject"
	P2PCmp          = "p2p"
	RPCClientCmp    = "rpc_client"
	DeployerCmp     = "deployer"
	NetwMngCmp      = "network_manager"
	WalletExtCmp    = "wallet_extension"
	TestGethNetwCmp = "test_geth_network"
	EthereumL1Cmp   = "l1_host"
	TenscanCmp      = "tenscan"
	CrossChainCmp   = "cross_chain"
)

// SysOut - Used when the logger has to write to Sys.out
const (
	SysOut = "sys_out"
)

var (
	glogger       *gethlog.GlogHandler
	logOutputFile io.WriteCloser
)

func init() {
	glogger = gethlog.NewGlogHandler(gethlog.NewTerminalHandler(os.Stderr, false))
}

// New - helper function used to create a top level logger for a component.
// Note: this expects legacy geth log levels, you will get unexpected behaviour if you use gethlog.<LEVEL> directly.
func New(component string, level int, out string, ctx ...interface{}) gethlog.Logger {
	logFile := ""
	if out != SysOut {
		logFile = out
	}
	verbosity := gethlog.FromLegacyLevel(level)

	err := Setup("terminal", logFile, false, 0, 0, 0, false, false, verbosity, "")
	if err != nil {
		panic(err.Error())
	}

	context := append(ctx, CmpKey, component)
	l := gethlog.New(context...)

	return l
}

// Setup initializes profiling and logging based on the CLI flags.
// It should be called as early as possible in the program.
func Setup(logFmtFlag string, logFile string, rotation bool, maxSize int, maxBackups int, maxAge int, compress bool, logJson bool, verbosity slog.Level, vmodule string) error {
	var (
		handler        slog.Handler
		terminalOutput = io.Writer(os.Stderr)
		output         io.Writer
	)
	if len(logFile) > 0 {
		if err := validateLogLocation(filepath.Dir(logFile)); err != nil {
			return fmt.Errorf("failed to initiatilize file logger: %v", err)
		}
	}
	context := []interface{}{"rotate", rotation}
	if len(logFmtFlag) > 0 {
		context = append(context, "format", logFmtFlag)
	} else {
		context = append(context, "format", "terminal")
	}
	if rotation {
		// Lumberjack uses <processname>-lumberjack.log in is.TempDir() if empty.
		// so typically /tmp/geth-lumberjack.log on linux
		if len(logFile) > 0 {
			context = append(context, "location", logFile)
		} else {
			context = append(context, "location", filepath.Join(os.TempDir(), "geth-lumberjack.log"))
		}
		logOutputFile = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}
		output = io.MultiWriter(terminalOutput, logOutputFile)
	} else if logFile != "" {
		var err error
		if logOutputFile, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err != nil {
			return err
		}
		// output = io.MultiWriter(logOutputFile, terminalOutput)
		// we only want to print to the file
		output = logOutputFile
		context = append(context, "location", logFile)
	} else {
		output = terminalOutput
	}

	switch {
	case logJson:
		// Retain backwards compatibility with `--log.json` flag if `--log.format` not set
		defer gethlog.Warn("The flag '--log.json' is deprecated, please use '--log.format=json' instead")
		handler = gethlog.JSONHandlerWithLevel(output, gethlog.LevelInfo)
	case logFmtFlag == "json":
		handler = gethlog.JSONHandlerWithLevel(output, gethlog.LevelInfo)
	case logFmtFlag == "logfmt":
		handler = gethlog.LogfmtHandler(output)
	case logFmtFlag == "", logFmtFlag == "terminal":
		useColor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
		if useColor {
			terminalOutput = colorable.NewColorableStderr()
			if logOutputFile != nil {
				output = io.MultiWriter(logOutputFile, terminalOutput)
			} else {
				output = terminalOutput
			}
		}
		handler = log.NewTerminalHandler(output, useColor)
	default:
		// Unknown log format specified
		return fmt.Errorf("unknown log format: %v", logFmtFlag)
	}

	glogger = gethlog.NewGlogHandler(handler)

	// logging
	glogger.Verbosity(verbosity)
	glogger.Vmodule(vmodule)

	gethlog.SetDefault(gethlog.NewLogger(glogger))

	if len(logFile) > 0 || rotation {
		gethlog.Info("Logging configured", context...)
	}
	return nil
}

func validateLogLocation(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("error creating the directory: %w", err)
	}
	// Check if the path is writable by trying to create a temporary file
	tmp := filepath.Join(path, "tmp")
	if f, err := os.Create(tmp); err != nil {
		return err
	} else {
		f.Close()
	}
	return os.Remove(tmp)
}
