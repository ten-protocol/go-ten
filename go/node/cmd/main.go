package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/node"
	"os"
	"strings"
)

func main() {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, nodeFlags, err := config.LoadFlagStrings(config.Node)
	if err != nil {
		panic(err)
	}

	if !validateNodeAction(rParams[node.Action]) {
		if rParams[node.Action] == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), rParams[node.Action])
		}
		os.Exit(1)
	}

	// retrieve node-config
	nodeConfig, err := ParseConfig(rParams)
	if err != nil {
		panic(err)
	}

	dockerNode := node.NewDockerNode(rParams, nodeConfig, nodeFlags)

	//// NETWORK CONFIGS INCLUDING NETWORK + NODE LEVEL details
	switch dockerNode.Action {
	case startAction:
		// write the network-level config to disk for future restarts
		err = node.WriteNetworkConfigToDisk(dockerNode.Cfg)
		if err != nil {
			panic(err)
		}
		err = dockerNode.Start()
	case upgradeAction:
		// load network-specific details from the initial node setup from disk
		var ntwCfg *config.NetworkInputConfig
		ntwCfg, err = node.ReadNetworkConfigFromDisk()
		if err != nil {
			panic(err)
		}

		err = dockerNode.Upgrade(ntwCfg)
	default:
		panic("unrecognized node action: " + dockerNode.Action)
	}
	if err != nil {
		panic(err)
	}
}
