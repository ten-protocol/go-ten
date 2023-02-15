package eth2network

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type RotatingLogWriter struct {
	dirPath         string
	fileStartName   string
	maxFileSize     int64
	maxNumFiles     int
	currentFile     *os.File
	currentSize     int64
	numFilesRotated int
}

func NewRotatingLogWriter(dirPath, fileStartName string, maxFileSize int64, maxNumFiles int) (*RotatingLogWriter, error) {
	writer := &RotatingLogWriter{
		dirPath:       dirPath,
		fileStartName: fileStartName,
		maxFileSize:   maxFileSize,
		maxNumFiles:   maxNumFiles,
	}

	// Open the current log file
	err := writer.openNextFile()
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingLogWriter) openNextFile() error {
	if w.currentFile != nil {
		err := w.currentFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	w.currentSize = 0
	w.numFilesRotated = 0

	// Construct the file name for the next log file
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := w.fileStartName + "-" + timestamp + ".log"
	filePath := filepath.Join(w.dirPath, filename)

	// Open the next log file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w.currentFile = file
	return nil
}

func (w *RotatingLogWriter) Write(p []byte) (int, error) {
	if w.currentFile == nil {
		return 0, io.ErrClosedPipe
	}

	// Rotate the log file if it exceeds the maximum file size
	if w.currentSize+int64(len(p)) > w.maxFileSize {
		err := w.rotateLogFile()
		if err != nil {
			return 0, err
		}
	}

	// Write the log message to the current log file
	n, err := w.currentFile.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *RotatingLogWriter) rotateLogFile() error {
	err := w.openNextFile()
	if err != nil {
		return err
	}

	w.numFilesRotated++

	// Delete old log files if the maximum number of files has been exceeded
	if w.maxNumFiles > 0 && w.numFilesRotated >= w.maxNumFiles {
		err = w.deleteOldLogFiles()
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *RotatingLogWriter) deleteOldLogFiles() error {
	// List all log files in the directory
	files, err := filepath.Glob(filepath.Join(w.dirPath, w.fileStartName+"-*.log"))
	if err != nil {
		return err
	}

	// Sort the log files by name (oldest first)
	sort.Strings(files)

	// Delete the oldest log files
	for i := 0; i < len(files)-w.maxNumFiles+1; i++ {
		err := os.Remove(files[i])
		if err != nil {
			return err
		}
	}

	return nil
}
