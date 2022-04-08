package simulation

import (
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"
)

const testLogs = "../.build/simulations/"

func setupTestLog() *os.File {
	// create a folder specific for the test
	err := os.MkdirAll(testLogs, 0o700)
	if err != nil {
		panic(err)
	}
	f, err := os.CreateTemp(testLogs, fmt.Sprintf("simulation-result-%d-*.txt", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	log.SetLog(f)
	return f
}

func minMax(arr []uint64) (min uint64, max uint64) {
	min = ^uint64(0)
	for _, no := range arr {
		if no < min {
			min = no
		}
		if no > max {
			max = no
		}
	}
	return
}
