package obscuro

import "os"

var logFile *os.File

func SetLog(f *os.File) {
	logFile = f
}

func log(msg string) {
	_, err := logFile.WriteString(msg)
	if err != nil {
		panic("could not write to file")
	}
}
