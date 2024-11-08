package main

import (
	"github.com/ten-protocol/go-ten/go/node"
)

func main() {
	cliConfig := ParseConfigCLI()

	tenCfg := NodeCLIConfigToTenConfig(cliConfig)

	dockerNode := node.NewDockerNode(tenCfg, cliConfig.hostDockerImage, cliConfig.enclaveDockerImage, cliConfig.edgelessDBImage, false, cliConfig.pccsAddr)
	var err error
	switch cliConfig.nodeAction {
	case startAction:
		// write the network-level config to disk for future restarts
		err = node.WriteNetworkConfigToDisk(tenCfg)
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
		panic("unrecognized node action: " + cliConfig.nodeAction)
	}
	if err != nil {
		panic(err)
	}
}
