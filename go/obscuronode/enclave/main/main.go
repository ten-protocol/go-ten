package main

import (
	"fmt"
	"os"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
)

func main() {
	setLogs()
	config := parseCLIArgs()

	nodeAddress := common.BytesToAddress([]byte(*config.nodeID))
	if err := enclave.StartServer(*config.port, nodeAddress, nil); err != nil {
		panic(err)
	}
	fmt.Printf("Enclave server listening on port %d.\n", *config.port)

	select {}
}

// Sets the log file.
func setLogs() {
	logFile, err := os.Create("enclave_logs.txt")
	if err != nil {
		panic(err)
	}
	log.SetLog(logFile)
}
