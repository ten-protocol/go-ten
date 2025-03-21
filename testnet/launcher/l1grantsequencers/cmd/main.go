package main

import (
	"fmt"
	"os"

	l1gs "github.com/ten-protocol/go-ten/testnet/launcher/l1grantsequencers"
)

func main() {
	cliConfig := ParseConfigCLI()

	l1grantsequencers, err := l1gs.NewGrantSequencers(
		l1gs.NewGrantSequencerConfig(
			l1gs.WithL1HTTPURL(cliConfig.l1HTTPURL),
			l1gs.WithPrivateKey(cliConfig.privateKey),
			l1gs.WithDockerImage(cliConfig.dockerImage),
			l1gs.WithEnclaveContractAddress(cliConfig.enclaveRegistryAddr),
			l1gs.WithEnclaveIDs(cliConfig.enclaveIDs),
			l1gs.WithSequencerURL(cliConfig.sequencerURL),
		),
	)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1grantsequencers.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1grantsequencers.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for grant sequnecer permission script to finish - %w", err)
		os.Exit(1)
	}
	fmt.Println("L1 Sequencer permissions were successfully granted...")
	os.Exit(0)
}
