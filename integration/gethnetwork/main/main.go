package main

import (
	"fmt"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"strconv"
)

const (
	gethBinaryPath = "/path/to/geth/binary"
	numNodes       = 5
	peerCountCmd   = "net.peerCount"
)

func main() {
	network := gethnetwork.NewGethNetwork(gethBinaryPath, numNodes)

	// Example usage:
	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)
	if peerCount != numNodes-1 {
		panic(fmt.Errorf("wrong number of peers on the network. Found %d, expected %d", peerCount, numNodes))
	}
}
