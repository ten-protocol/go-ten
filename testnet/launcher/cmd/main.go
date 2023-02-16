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

	fmt.Println("Starting a testnet with all the defaults...")
	testnet := launcher.NewTestnetLauncher(
		launcher.NewTestnetConfig(
			launcher.WithNumberNodes(cliConfig.numberNodes), // TODO: currently ignored flag
			// launcher.WithEnclaveDockerImage("testnetobscuronet.azurecr.io/obscuronet/enclave_debug:latest"),
			// launcher.WithEnclaveDebug(true),
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
