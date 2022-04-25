//go:build geth
// +build geth

package gethnetwork

import (
	"fmt"
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
		t.Fatalf("Wrong number of peers on the network. Found %d, expected %d.", peerCount, numNodes-1)
	}
}

func TestGenesisParamsAreUsed(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes)

	chainId := network.IssueCommand(0, chainIdCmd)
	if chainId != expectedChainId {
		t.Fatalf("Network not using chain ID specified in the genesis file. Found %s, expected %s.", chainId, expectedChainId)
	}
}

func TestTransactionCanBeSubmitted(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes)

	account := network.addresses[0]
	tx := fmt.Sprintf("{from: \"%s\", to: \"%s\", value: web3.toWei(0.001, \"ether\")}", account, account)
	txHash := network.IssueCommand(0, fmt.Sprintf("personal.sendTransaction(%s, \"%s\")", tx, password))
	status := network.IssueCommand(0, fmt.Sprintf("eth.getTransaction(\"%s\")", txHash))

	if status == "null" {
		t.Fatal("Could not issue transaction.")
	}
}
