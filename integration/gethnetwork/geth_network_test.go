//go:build geth
// +build geth

package gethnetwork

import (
	"strconv"
	"testing"
)

const (
	gethBinaryPath = "/Users/joeldudley/Desktop/repos/go-ethereum/cmd/geth/geth"
	numNodes       = 5
	peerCountCmd   = "net.peerCount"
)

func TestCreateNetworkAndIssueCommands(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes)

	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)
	if peerCount != numNodes-1 {
		t.Fatalf("wrong number of peers on the network. Found %d, expected %d", peerCount, numNodes)
	}
}
