package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/config"
	l1cp "github.com/ten-protocol/go-ten/testnet/launcher/l1challengeperiod"
)

func main() {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}

	challengePeriodCfg := l1cp.NewChallengePeriodConfig(tenCfg)
	l1challengeperiod, err := l1cp.NewSetChallengePeriod(challengePeriodCfg)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1challengeperiod.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1challengeperiod.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for set challenge period script to finish - %w", err)
		os.Exit(1)
	}
	fmt.Println("L1 challenge period was successfully set...")
	os.Exit(0)
}
