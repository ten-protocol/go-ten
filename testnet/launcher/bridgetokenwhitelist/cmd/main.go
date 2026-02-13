package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/testnet/launcher/bridgetokenwhitelist"
)

func main() {
	cliConfig := ParseConfigCLI()

	whitelister, err := bridgetokenwhitelist.NewBridgeTokenWhitelister(
		bridgetokenwhitelist.NewConfig(
			bridgetokenwhitelist.WithTokenAddress(cliConfig.tokenAddress),
			bridgetokenwhitelist.WithTokenName(cliConfig.tokenName),
			bridgetokenwhitelist.WithTokenSymbol(cliConfig.tokenSymbol),
			bridgetokenwhitelist.WithNetworkEnv(cliConfig.networkEnv),
			bridgetokenwhitelist.WithL1HTTPURL(cliConfig.l1HTTPURL),
			bridgetokenwhitelist.WithL2GatewayURL(cliConfig.l2GatewayURL),
			bridgetokenwhitelist.WithPrivateKey(cliConfig.privateKey),
			bridgetokenwhitelist.WithDockerImage(cliConfig.dockerImage),
			bridgetokenwhitelist.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure bridge token whitelister: %s\n", err)
		os.Exit(1)
	}

	err = whitelister.Start()
	if err != nil {
		fmt.Printf("unable to start bridge token whitelister: %s\n", err)
		os.Exit(1)
	}

	err = whitelister.WaitForFinish()
	if err != nil {
		fmt.Printf("bridge token whitelisting failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Token whitelisted on bridge and registered in NetworkConfig successfully!")
	os.Exit(0)
}
