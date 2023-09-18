package main

import (
	"fmt"
	"os"

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
			launcher.WithSGXEnabled(cliConfig.isSGXEnabled),
		),
	)
	err := testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
}
