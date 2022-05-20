package enclaverunner

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
)

// RunEnclave runs an Obscuro enclave as a standalone process.
func RunEnclave(config EnclaveConfig) {
	nodeAddress := common.BigToAddress(big.NewInt(config.NodeID))
	contractAddr := common.HexToAddress(config.ContractAddress)
	txHandler := mgmtcontractlib.NewEthMgmtContractTxHandler(contractAddr)

	// TODO - For now, genesisJSON is nil. This means that incoming L1 blocks are not validated by the enclave. In the
	//  future, we should allow the genesisJSON to be passed in somehow, with a default of the default genesis.
	closeHandle, err := enclave.StartServer(config.Address, nodeAddress, txHandler, false, nil, nil)
	if err != nil {
		log.Panic("could not start Obscuro enclave service. Cause: %s", err)
	}
	log.Info("Obscuro enclave service started.")
	fmt.Println("Obscuro enclave service started.")

	handleInterrupt(closeHandle)
}

// SetLogs sets the log file, defaulting to stdout if writeToLogs is false.
func SetLogs(writeToLogs bool, logPath string) {
	var logFile *os.File
	var err error

	if writeToLogs {
		logFile, err = os.Create(logPath)
		if err != nil {
			panic(fmt.Sprintf("could not create log file. Cause: %s", err))
		}
	} else {
		logFile = os.Stdout
	}

	log.SetLog(logFile)
}

// Shuts down the Obscuro enclave service when an interrupt is received.
func handleInterrupt(closeHandle func()) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	closeHandle()
	fmt.Println("Obscuro enclave service stopping...")
	os.Exit(1)
}
