package main

import (
	"fmt"
	"os"

	funds "github.com/ten-protocol/go-ten/testnet/launcher/fundsrecovery"
)

func main() {
	cliConfig := ParseConfigCLI()

	fundsRecovery, err := funds.NewFundsRecovery(
		funds.NewFundsRecoveryConfig(
			funds.WithL1HTTPURL(cliConfig.l1HTTPURL),
			funds.WithL1PrivateKey(cliConfig.privateKey),
			funds.WithDockerImage(cliConfig.dockerImage),
		),
	)
	if err != nil {
		fmt.Println("unable to configure the funds recovery - ", err)
		os.Exit(1)
	}

	err = fundsRecovery.Start()
	if err != nil {
		fmt.Println("unable to start the funds recovery - ", err)
		os.Exit(1)
	}

	err = fundsRecovery.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for funds recovery to finish - ", err)
		os.Exit(1)
	}
	fmt.Println("Funds recovery was successful...")
	os.Exit(0)
}
