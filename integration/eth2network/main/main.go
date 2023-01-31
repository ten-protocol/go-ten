package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/go-obscuro/integration/eth2network"
)

// Spins up a new eth 2 network.
func main() {
	config := parseCLIArgs()

	binariesPath, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}
	eth2Network := eth2network.NewEth2Network(
		binariesPath,
		config.gethHTTPStartPort,
		config.gethWSStartPort,
		config.gethAuthRPCStartPort,
		config.gethNetworkStartPort,
		config.prysmBeaconRPCStartPort,
		1337,
		1,
		config.blockTimeSecs,
		config.prefundedAddrs,
	)

	err = eth2Network.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("eth 2 network started..")

	handleInterrupt(eth2Network)
}

// Shuts down the Geth network when an interrupt is received.
func handleInterrupt(network eth2network.Eth2Network) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	network.Stop()
	fmt.Println("eth2 network stopping...")
	os.Exit(1)
}
