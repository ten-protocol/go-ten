package main

import (
	"github.com/ten-protocol/go-ten/go/node"
)

func main() {
	// load default config
	var err error

	// set default, load overrides and cli flag config if applicable
	defaultConfig := LoadDefaultConfig()
	//overrideConfig := LoadOverrideConfig(defaultConfig)
	cliConfig := ParseConfigCLI(defaultConfig)
	// todo (#1618) - allow for multiple operation (start, stop, status)

	dockerNode := node.NewDockerNode(cliConfig)
	switch cliConfig.NodeAction {
	case startAction:
		// write the network-level config to disk for future restarts
		err = node.WriteNetworkConfigToDisk(cliConfig)
		if err != nil {
			panic(err)
		}
		err = dockerNode.Start()
	case upgradeAction:
		// load network-specific details from the initial node setup from disk
		var ntwCfg *node.NetworkConfig
		ntwCfg, err = node.ReadNetworkConfigFromDisk()
		if err != nil {
			panic(err)
		}

		err = dockerNode.Upgrade(ntwCfg)
	default:
		panic("unrecognized node action: " + cliConfig.NodeAction)
	}
	if err != nil {
		panic(err)
	}
}
