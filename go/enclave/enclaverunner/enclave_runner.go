package enclaverunner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/enclave"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

// TODO - Replace with the genesis.json of Obscuro's L1 network.
const hardcodedGenesisJSON = "TODO - REPLACE ME"

// RunEnclave runs an Obscuro enclave as a standalone process.
func RunEnclave(config config.EnclaveConfig) {
	// todo - is this the right wiring?
	logger := log.New(log.EnclaveCmp, config.LogLevel, config.LogPath, log.NodeIDKey, config.HostID)

	contractAddr := config.ManagementContractAddress
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&contractAddr, logger)
	erc20ContractLib := erc20contractlib.NewERC20ContractLib(&contractAddr, config.ERC20ContractAddresses...)

	if config.ValidateL1Blocks {
		config.GenesisJSON = []byte(hardcodedGenesisJSON)
	}
	closeHandle, err := enclave.StartServer(config, mgmtContractLib, erc20ContractLib, logger)
	if err != nil {
		logger.Crit("could not start Obscuro enclave service.", log.ErrKey, err)
	}
	logger.Info("Obscuro enclave service started.")
	fmt.Println("Obscuro enclave service started.")

	handleInterrupt(closeHandle)
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
