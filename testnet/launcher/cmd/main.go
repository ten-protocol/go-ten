package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/ten-protocol/go-ten/testnet/launcher"
)

func main() {
	// Try to overload environment variables if .env file exists
	errEnv := godotenv.Overload()
	if errEnv != nil {
		log.Printf("Optional .env file not supplied: %v", errEnv)
	}

	cliConfig := ParseConfigCLI()

	fmt.Println("Starting a testnet with 1 sequencer and 1 validator...")
	testnet := launcher.NewTestnetLauncher(
		launcher.NewTestnetConfig(
			launcher.WithValidatorEnclaveDockerImage(cliConfig.validatorEnclaveDockerImage),
			launcher.WithValidatorEnclaveDebug(cliConfig.validatorEnclaveDebug),
			launcher.WithSequencerEnclaveDockerImage(cliConfig.sequencerEnclaveDockerImage),
			launcher.WithSequencerEnclaveDebug(cliConfig.sequencerEnclaveDebug),
			launcher.WithContractDeployerDebug(cliConfig.contractDeployerDebug),
			launcher.WithContractDeployerDockerImage(cliConfig.contractDeployerDockerImage),
			launcher.WithSGXEnabled(cliConfig.isSGXEnabled),
			launcher.WithLogLevel(cliConfig.logLevel),
		),
	)
	err := testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
}
