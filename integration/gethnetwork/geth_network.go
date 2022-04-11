package gethnetwork

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TODO - Use platform-agnostic separators, instead of forward-slashes.

const (
	genesisFileRelPath = "/integration/gethnetwork/genesis.json"
	nodeFolderName     = "/node%d_datadir/"
	ipcFileName        = "geth.ipc"

	dataDirFlag = "--datadir"
	execFlag    = "--exec"
	portFlag    = "--port"

	addPeerCmd = "admin.addPeer(%s)"
	attachCmd  = "attach"
	enodeCmd   = "admin.nodeInfo.enode"
	initCmd    = "init"
)

// GethNetwork is a network of Geth nodes, built using the provided Geth binary.
type GethNetwork struct {
	gethBinaryPath string
	dataDirs       []string
}

// NewGethNetwork using the provided Geth binary to create a private Ethereum network with numNodes Geth nodes.
func NewGethNetwork(gethBinaryPath string, nodesDir string, numNodes int) GethNetwork {
	err := os.MkdirAll(nodesDir, 0o700)
	if err != nil {
		panic(err)
	}

	// Each Geth node needs its own data directory.
	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		dataDirs[i] = fmt.Sprintf(nodesDir+nodeFolderName, i+1)
	}

	network := GethNetwork{
		gethBinaryPath: gethBinaryPath,
		dataDirs:       dataDirs,
	}

	for i, dataDir := range dataDirs {
		_ = os.Remove(dataDir + ipcFileName) // We delete leftover IPC files from previous runs.
		go network.createGethNode(dataDir, 30303+i)
	}

	// We need to manually tell the nodes about one another.
	network.joinNodesToNetwork(dataDirs)

	return network
}

// IssueCommand sends the command via RPC to the nodeIdx'th node in the network.
func (network *GethNetwork) IssueCommand(nodeIdx int, command string) string {
	dataDir := network.dataDirs[nodeIdx]

	args := []string{dataDirFlag, dataDir, attachCmd, dataDir + ipcFileName, execFlag, command}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(output))
}

// Initialises and starts a Geth node.
func (network *GethNetwork) createGethNode(dataDir string, port int) {
	network.initGethNode(dataDir)
	network.startGethNode(dataDir, port)
}

// Initialises a Geth node based on the network genesis file.
func (network *GethNetwork) initGethNode(dataDirPath string) {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	args := []string{dataDirFlag, dataDirPath, initCmd, workingDir + genesisFileRelPath}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Starts a Geth node.
func (network *GethNetwork) startGethNode(dataDirPath string, port int) {
	args := []string{dataDirFlag, dataDirPath, fmt.Sprintf("%s=%d", portFlag, port)}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint

	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

// Tells the network's nodes about one another.
func (network *GethNetwork) joinNodesToNetwork(dataDirs []string) {
	enodeAddrs := make([]string, len(dataDirs))

	for i, dataDir := range network.dataDirs {
		waitForIPC(dataDir) // We cannot issue RPC commands until the IPC files are available.
		enodeAddrs[i] = network.IssueCommand(i, enodeCmd)
	}

	for _, enodeAddr := range enodeAddrs {
		for i := range dataDirs {
			// As part of this loop, we also try and peer a node with itself, but Geth ignores this.
			network.IssueCommand(i, fmt.Sprintf(addPeerCmd, enodeAddr))
		}
	}
}

// Waits for a node's IPC file to exist.
func waitForIPC(dataDir string) {
	for {
		_, err := os.Stat(dataDir + ipcFileName)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
