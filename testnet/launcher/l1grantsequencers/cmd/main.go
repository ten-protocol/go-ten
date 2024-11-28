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
			l1gs.WithL1HTTPURL("http://eth2network:8025"),
			l1gs.WithPrivateKey("f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1gs.WithDockerImage(cliConfig.dockerImage),
			l1gs.WithMgmtContractAddress(cliConfig.mgmtContractAddress),
			l1gs.WithEnclaveIDs(cliConfig.enclaveIDs),
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

	os.Exit(0)
}
