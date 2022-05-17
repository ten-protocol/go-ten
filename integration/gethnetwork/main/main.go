package main

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
)

// Spins up a new Geth network.
func main() {
	config := ParseCLIArgs()

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}
	gethnetwork.NewGethNetwork(config.StartPort, gethBinaryPath, config.NumNodes, 1, config.PrefundedAddrs)
	fmt.Println("Geth network started.")

	select {}
}
