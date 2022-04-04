package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
)

func main() {
	println(os.Args)
	config := parseCLIArgs()
	setLogs(*config.writeToLogs)

	nodeAddress := common.BigToAddress(big.NewInt(*config.nodeID))
	if err := enclave.StartServer(*config.address, nodeAddress, nil); err != nil {
		panic(err)
	}
	fmt.Printf("Enclave server listening on address %s.\n", *config.address)

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
