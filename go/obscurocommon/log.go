package obscurocommon

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
		panic("could not write to file")
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
