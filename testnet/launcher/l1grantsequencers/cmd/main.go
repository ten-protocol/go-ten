package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/config"
	l1gs "github.com/ten-protocol/go-ten/testnet/launcher/l1grantsequencers"
)

func main() {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}

	grantSeqCfg := l1gs.NewGrantSequencerConfig(tenCfg)
	l1grantsequencers, err := l1gs.NewGrantSequencers(grantSeqCfg)
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
