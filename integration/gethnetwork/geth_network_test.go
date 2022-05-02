package gethnetwork

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"
)

const (
	shCmd          = "sh"
	gethBinarySh   = "./build_geth_binary.sh"
	gethBinaryPath = "./geth-release-1.10.17"

	numNodes        = 3
	expectedChainID = "777"

	peerCountCmd = "net.peerCount"
	chainIDCmd   = "admin.nodeInfo.protocols.eth.config.chainId"
)

func init() { //nolint:gochecknoinits
	_, err := exec.Command(shCmd, gethBinarySh).Output()
	if err != nil {
		panic(err)
	}
}

func TestAllNodesJoinSameNetwork(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes, 1)
	defer network.StopNodes()

	peerCountStr := network.IssueCommand(0, peerCountCmd)
	peerCount, _ := strconv.Atoi(peerCountStr)
	if peerCount != numNodes-1 {
		t.Fatalf("Wrong number of peers on the network. Found %d, expected %d.", peerCount, numNodes-1)
	}
}

func TestGenesisParamsAreUsed(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes, 1)
	defer network.StopNodes()

	chainID := network.IssueCommand(0, chainIDCmd)
	if chainID != expectedChainID {
		t.Fatalf("Network not using chain ID specified in the genesis file. Found %s, expected %s.", chainID, expectedChainID)
	}
}

func TestTransactionCanBeSubmitted(t *testing.T) {
	network := NewGethNetwork(gethBinaryPath, numNodes, 1)
	defer network.StopNodes()

	account := network.addresses[0]
	tx := fmt.Sprintf("{from: \"%s\", to: \"%s\", value: web3.toWei(0.001, \"ether\")}", account, account)
	txHash := network.IssueCommand(0, fmt.Sprintf("personal.sendTransaction(%s, \"%s\")", tx, password))
	status := network.IssueCommand(0, fmt.Sprintf("eth.getTransaction(\"%s\")", txHash))

	if status == "null" {
		t.Fatal("Could not issue transaction.")
	}
}
