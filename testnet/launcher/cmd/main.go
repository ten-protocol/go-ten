package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/go-obscuro/testnet/launcher"
)

func main() {
	cliConfig := ParseConfigCLI()

	fmt.Println("Starting a testnet with 1 sequencer and 1 validator...")
	testnet := launcher.NewTestnetLauncher(
		launcher.NewTestnetConfig(
			launcher.WithValidatorEnclaveDockerImage(cliConfig.validatorEnclaveDockerImage),
			launcher.WithValidatorEnclaveDebug(cliConfig.validatorEnclaveDebug),
			launcher.WithSequencerEnclaveDockerImage(cliConfig.sequencerEnclaveDockerImage),
			launcher.WithSequencerEnclaveDebug(cliConfig.sequencerEnclaveDebug),
		),
	)
	err := testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press ctrl+c to stop...")
	<-done // Will block here until user hits ctrl+c
	// TODO add clean up / teardown
	os.Exit(0)
}
