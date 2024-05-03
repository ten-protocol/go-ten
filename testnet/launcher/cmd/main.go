package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/testnet/launcher"
	"os"
)

func main() {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, _, err := config.LoadFlagStrings(config.Testnet)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting a testnet with 1 sequencer and 2 validator...")
	testNetConfig, err := ParseConfig(rParams)
	if err != nil {
		panic(err)
	}

	if testNetConfig.TestNetSettings.ReplaceRunning {
		fmt.Println("Replacing any running testnet elements selected for in config...")
		if testNetConfig.TestNetSettings.Nodes {
			for _, node := range testNetConfig.Nodes {
				err := docker.RemoveContainer(node.NodeDetails.NodeName + "-host")
				if err != nil {
					fmt.Println(err)
				}
				err = docker.RemoveContainer(node.NodeDetails.NodeName + "-enclave")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if testNetConfig.TestNetSettings.Faucet {
			err = docker.RemoveContainer(testNetConfig.Faucet.ContainerName)
			if err != nil {
				fmt.Println(err)
			}
		}
		if testNetConfig.TestNetSettings.L2ContractDeployer {
			err = docker.RemoveContainer(testNetConfig.L2ContractDeployer.ContainerName)
			if err != nil {
				fmt.Println(err)
			}
		}
		if testNetConfig.TestNetSettings.L1ContractDeployer {
			err = docker.RemoveContainer(testNetConfig.L1ContractDeployer.ContainerName)
			if err != nil {
				fmt.Println(err)
			}
		}
		if testNetConfig.TestNetSettings.Eth2Network {
			err = docker.RemoveContainer(testNetConfig.Eth2Network.ContainerName)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	testnet := launcher.NewTestnetLauncher(testNetConfig)

	err = testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully!")
}
