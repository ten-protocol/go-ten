package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/testnet/launcher"
	"os"
)

func main() {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, _, err := config.LoadFlagStrings(config.Node)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting a testnet with 1 sequencer and 2 validator...")
	testNetConfig, err := ParseConfig(rParams)
	if err != nil {
		panic(err)
	}

	testnet := launcher.NewTestnetLauncher(testNetConfig)

	err = testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
}
