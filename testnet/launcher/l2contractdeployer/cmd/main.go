package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/config"
	l2cd "github.com/ten-protocol/go-ten/testnet/launcher/l2contractdeployer"
)

func main() {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}
	fmt.Println("Starting L2 contract deployer with the following TenConfig:")
	tenCfg.PrettyPrint() // dump config to stdout

	l2CDCfg := l2cd.NewContractDeployerConfig(tenCfg)
	l2ContractDeployer, err := l2cd.NewDockerContractDeployer(l2CDCfg)
	if err != nil {
		fmt.Println("unable to configure the l2 contract deployer - ", err)
		os.Exit(1)
	}

	err = l2ContractDeployer.Start()
	if err != nil {
		fmt.Println("unable to start the l2 contract deployer - ", err)
		os.Exit(1)
	}

	err = l2ContractDeployer.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for l2 contract deployer to finish - ", err)
		os.Exit(1)
	}
	fmt.Println("L2 Contracts were successfully deployed...")
	os.Exit(0)
}
