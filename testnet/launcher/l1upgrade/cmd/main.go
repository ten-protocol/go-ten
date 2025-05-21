package main

import (
	"fmt"
	"os"

	l1upgrade "github.com/ten-protocol/go-ten/testnet/launcher/l1upgrade"
)

func main() {
	cliConfig := ParseConfigCLI()

	l1Upgrade, err := l1upgrade.NewUpgradeContracts(
		l1upgrade.NewUpgradeContractsConfig(
			l1upgrade.WithL1HTTPURL(cliConfig.l1HTTPURL),
			l1upgrade.WithPrivateKey(cliConfig.privateKey),
			l1upgrade.WithDockerImage(cliConfig.dockerImage),
			l1upgrade.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
		),
	)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1Upgrade.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1Upgrade.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for grant sequnecer permission script to finish - %w", err)
		os.Exit(1)
	}
	fmt.Println("L1 upgrades were successfully completed...")
	os.Exit(0)
}
