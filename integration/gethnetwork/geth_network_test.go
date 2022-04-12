//go:build geth
// +build geth

package gethnetwork

import (
	"strconv"
	"testing"
)

const (
	// The Geth binary can be built using the instructions here: https://github.com/ethereum/go-ethereum#building-the-source.
	gethBinaryPath  = "path/to/geth/binary"
	numNodes        = 3
	expectedChainId = "777"

	peerCountCmd = "net.peerCount"
	chainIdCmd   = "admin.nodeInfo.protocols.eth.config.chainId"
)

func TestAllNodesJoinSameNetwork(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes)

	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)
	if peerCount != numNodes-1 {
		t.Fatalf("Wrong number of peers on the network. Found %d, expected %d.", peerCount, numNodes)
	}
}

func TestGenesisParamsAreUsed(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes)

	chainId := network.IssueCommand(0, chainIdCmd)
	if chainId != expectedChainId {
		t.Fatalf("Network not using chain ID specified in the genesis file. Found %s, expected %s.", chainId, expectedChainId)
	}
}
