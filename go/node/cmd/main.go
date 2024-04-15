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
	action, cPaths, nodeFlags, err := config.LoadFlagStrings(config.Node)
	if err != nil {
		panic(err)
	}

	if !validateNodeAction(action) {
		if action == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), action)
		}
		os.Exit(1)
	}

	// retrieve node-config
	nodeConfig, err := ParseConfig(cPaths)
	if err != nil {
		panic(err)
	}
	//
	//// Change anything that was statically defined in the nodeConf to a set of envs for
	//envs := nodeConfig.GetConfigAsEnvVars(config.Host)
	//envs = config.MergeEnvMaps(envs, nodeConfig.GetConfigAsEnvVars(config.Enclave))

	// override with flags
	//print(envs)

	// inside runs then split to service

	dockerNode := node.NewDockerNode(action, nodeConfig, nodeFlags)

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
