package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"math/big"
	"os"
)

func main() {
	config := parseCLIArgs()
	setLogs(*config.writeToLogs)

	nodeAddress := common.BigToAddress(big.NewInt(*config.nodeID))
	if err := enclave.StartServer(*config.address, nodeAddress, nil); err != nil {
		panic(err)
	}

	select {}
}

// Sets the log file, defaulting to stdout if writeToLogs is false.
func setLogs(writeToLogs bool) {
	var logFile *os.File
	var err error

	if writeToLogs {
		logFile, err = os.Create("enclave_logs.txt")
		if err != nil {
			panic(err)
		}
	} else {
		logFile = os.Stdout
	}

	log.SetLog(logFile)
}
