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

func Log(msg string) {
	_, err := logFile.WriteString(fmt.Sprintf("%d %s\n", makeTimestamp(), msg))
	if err != nil {
		if logFile == nil {
			panic("logger could not write as log file not set")
		}
		panic(fmt.Errorf("logger could not write to log file: %w", err))
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
