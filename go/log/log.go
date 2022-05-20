package log

import (
	"fmt"
	"os"
	"time"
)

const (
	InfoLevel = iota
	DebugLevel
	TraceLevel
)

var (
	logFile   *os.File
	logStatus = InfoLevel
)

func SetLogLevel(level int) {
	if level < InfoLevel || level > TraceLevel {
		panic("unrecognized log level")
	}
	logStatus = level
}

func SetLog(f *os.File) {
	logFile = f
}

func Panic(msg string, args ...interface{}) {
	write("PANIC: "+msg, args...)
	logFile.Close()
	panic(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...interface{}) {
	write("ERROR: "+msg, args...)
}

func Trace(msg string, args ...interface{}) {
	if logStatus >= TraceLevel {
		write("TRACE: "+msg, args...)
	}
}

func Debug(msg string, args ...interface{}) {
	if logStatus >= DebugLevel {
		write("DEBUG: "+msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	write("INFO: "+msg, args...)
}

func write(msg string, args ...interface{}) {
	var wMsg string
	if len(args) == 0 {
		wMsg = msg
	} else {
		wMsg = fmt.Sprintf(msg, args...)
	}
	if logFile == nil {
		// defaults to outputting logs to stdout
		// things like unit tests don't require a logfile
		fmt.Println(msg)
		return
	}

	_, err := logFile.WriteString(fmt.Sprintf("%d %s\n", makeTimestamp(), wMsg))
	if err != nil {
		panic(fmt.Errorf("logger could not write to log file: %w", err))
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
