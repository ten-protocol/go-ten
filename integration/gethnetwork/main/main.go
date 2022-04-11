package main

import (
	"fmt"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"strconv"
)

const (
	// TODO - Use agnostic paths.
	gethBinaryPath = "/Users/joeldudley/Desktop/repos/go-ethereum/cmd/geth/geth"
	nodesDir       = "../.build/12847t/"
	numNodes       = 5
	peerCountCmd   = "net.peerCount"
)

func main() {
	network := gethnetwork.NewGethNetwork(gethBinaryPath, nodesDir, numNodes)

	// Example usage:
	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)
	if peerCount != numNodes-1 {
		panic(fmt.Errorf("wrong number of peers on the network. Found %d, expected %d", peerCount, numNodes))
	}
}
