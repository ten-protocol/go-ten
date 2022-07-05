package enclaverunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/obscuronet/obscuro-playground/go/config"

	"github.com/obscuronet/obscuro-playground/go/enclave"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethadapter/mgmtcontractlib"
)

// TODO - Replace with the genesis.json of Obscuro's L1 network.
const hardcodedGenesisJSON = "TODO - REPLACE ME"

// RunEnclave runs an Obscuro enclave as a standalone process.
func RunEnclave(config config.EnclaveConfig) {
	contractAddr := config.ManagementContractAddress
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&contractAddr)
	erc20ContractLib := erc20contractlib.NewERC20ContractLib(&contractAddr, config.ERC20ContractAddresses...)

	log.SetLogLevel(log.ParseLevel(config.LogLevel))

	if config.ValidateL1Blocks {
		config.GenesisJSON = []byte(hardcodedGenesisJSON)
	}
	closeHandle, err := enclave.StartServer(config, mgmtContractLib, erc20ContractLib, nil)
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

	log.OutputToFile(logFile)
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
