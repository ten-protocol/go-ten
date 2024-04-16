package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/testnet/launcher"
	"os"
)

var configs = map[string]string{
	sequencer: "./testnet/config/s1_node.yaml",
	validator: "./testnet/config/v1_node.yaml",
}

const (
	sequencer = "sequencer"
	validator = "validator"
)

func main() {
	var err error
	// load flags with defaults from config / sub-configs
	_, _, err = config.LoadFlagStrings(config.Node)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting a testnet with 1 sequencer and 2 validator...")
	configMap, err := ParseConfig(configs)
	if err != nil {
		panic(err)
	}
	nodeConf := &launcher.Config{
		Nodes: configMap,
	}

	testnet := launcher.NewTestnetLauncher(nodeConf)

	err = testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
}
