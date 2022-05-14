package log

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func SetLog(f *os.File) {
	logFile = f
}

func Trace(msg string, args ...interface{}) {
	// fmt.Printf(msg+"\n", args...)
}

func Log(msg string) {
	if logFile == nil {
		// defaults to outputting logs to stdout
		// things like unit tests don't require a logfile
		fmt.Println(msg)
		return
	}

	_, err := logFile.WriteString(fmt.Sprintf("%d %s\n", makeTimestamp(), msg))
	if err != nil {
		panic(fmt.Errorf("logger could not write to log file: %w", err))
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
