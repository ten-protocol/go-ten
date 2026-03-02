package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/testnet/launcher/erc20deployer"
)

func main() {
	cliConfig := ParseConfigCLI()

	deployer, err := erc20deployer.NewERC20Deployer(
		erc20deployer.NewConfig(
			erc20deployer.WithTokenName(cliConfig.tokenName),
			erc20deployer.WithTokenSymbol(cliConfig.tokenSymbol),
			erc20deployer.WithTokenDecimals(cliConfig.tokenDecimals),
			erc20deployer.WithTokenSupply(cliConfig.tokenSupply),
			erc20deployer.WithL1HTTPURL(cliConfig.l1HTTPURL),
			erc20deployer.WithPrivateKey(cliConfig.privateKey),
			erc20deployer.WithDockerImage(cliConfig.dockerImage),
			erc20deployer.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure ERC20 deployer: %s\n", err)
		os.Exit(1)
	}

	err = deployer.Start()
	if err != nil {
		fmt.Printf("unable to start ERC20 deployer: %s\n", err)
		os.Exit(1)
	}

	err = deployer.WaitForFinish()
	if err != nil {
		fmt.Printf("ERC20 deployment failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("ERC20 token deployed and registered successfully!")
	os.Exit(0)
}
