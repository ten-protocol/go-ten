package eth2network

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRotatingLogWriter(t *testing.T) {
	temp, err := os.MkdirTemp("", "*")
	assert.Nil(t, err)

	// Create a rotating log writer with a maximum size of 1 MB
	writer, err := NewRotatingLogWriter(temp, "derp", 1024*1024, 5)
	assert.Nil(t, err)

	createdFile := writer.currentFile.Name()

	// Use the rotating log writer with the standard log package
	data := make([]byte, 1024*1024)
	rand.Read(data) //nolint:gosec
	_, err = writer.Write(data)
	assert.Nil(t, err)

	// ensure the file name is different
	time.Sleep(2 * time.Second)

	// Use the rotating log writer with the standard log package
	rand.Read(data) //nolint:gosec
	_, err = writer.Write(data)
	assert.Nil(t, err)

	assert.NotEqual(t, createdFile, writer.currentFile.Name())
}
