package main

import "github.com/obscuronet/obscuro-playground/integration/gethnetwork"

func main() {
	gethBinaryPath := "/Users/joeldudley/Desktop/repos/go-ethereum/cmd/geth/geth"
	network := gethnetwork.NewGethNetwork(gethBinaryPath, 5, "/Users/joeldudley/Desktop/geth_network")
	println(network.IssueCommand(0, "net.peerCount"))
}
