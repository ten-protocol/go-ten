package gethnetwork

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	genesisFileName = "genesis.json"
	nodeFolderName  = "node%d_datadir"
	ipcFileName     = "geth.ipc"
	tempDirPrefix   = "geth_nodes"
	buildDir        = "../.build/geth"
	nodeLogs        = "node_logs.txt"
	startPort       = 30303

	addPeerCmd = "admin.addPeer(%s)"
	attachCmd  = "attach"
	enodeCmd   = "admin.nodeInfo.enode"
	initCmd    = "init"

	websocketFlag    = "--ws" // Enables websocket connections to the node.
	dataDirFlag      = "--datadir"
	execFlag         = "--exec"
	portFlag         = "--port"
	mineFlag         = "--mine"
	minerThreadsFlag = "--miner.threads=1"
	minerEthbaseFlag = "--miner.etherbase=0x0000000000000000000000000000000000000001"
	rpcFeeCapFlag    = "--rpc.txfeecap=0" // Disables the 1 ETH cap for RPC transactions.

	// We pre-allocate a single wallet, which is shared by all nodes.
	genesisConfig = `{
	  "config": {
		"chainId": 777,
		"homesteadBlock": 0,
		"eip150Block": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"byzantiumBlock": 0,
		"constantinopleBlock": 0,
		"petersburgBlock": 0,
		"istanbulBlock": 0,
		"berlinBlock": 0,
		"londonBlock": 0
	  },
	  "alloc": {
		"0x0000000000000000000000000000000000000001": {
		  "balance": "111111111"
		}
	  },
	  "coinbase": "0x0000000000000000000000000000000000000000",
	  "difficulty": "0x20000",
	  "extraData": "",
	  "gasLimit": "0x2fefd8",
	  "nonce": "0x0000000000000042",
	  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "timestamp": "0x00"
  }`
)

// GethNetwork is a network of Geth nodes, built using the provided Geth binary.
type GethNetwork struct {
	gethBinaryPath  string
	genesisFilePath string
	dataDirs        []string
	logFile         *os.File
}

// NewGethNetwork using the provided Geth binary to create a private Ethereum network with numNodes Geth nodes.
func NewGethNetwork(gethBinaryPath string, numNodes int) GethNetwork {
	nodesDir, err := ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		panic(err)
	}

	// We write out the `genesis.json` file to be used by the network.
	genesisFilePath := path.Join(nodesDir, genesisFileName)
	err = os.WriteFile(genesisFilePath, []byte(genesisConfig), 0o600)
	if err != nil {
		panic(err)
	}

	// Each Geth node needs its own data directory.
	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		nodeFolder := fmt.Sprintf(nodeFolderName, i+1)
		dataDirs[i] = path.Join(nodesDir, nodeFolder)
	}

	// All the Geth node logs are pushed to a single log file.
	err = os.MkdirAll(buildDir, 0o700)
	if err != nil {
		panic(err)
	}
	logFile, _ := os.Create(path.Join(buildDir, nodeLogs))

	network := GethNetwork{
		gethBinaryPath:  gethBinaryPath,
		genesisFilePath: genesisFilePath,
		dataDirs:        dataDirs,
		logFile:         logFile,
	}

	for i, dataDir := range dataDirs {
		_ = os.Remove(path.Join(dataDir, ipcFileName)) // We delete leftover IPC files from previous runs.
		go network.createMiner(dataDir, startPort+i)
	}

	// We need to manually tell the nodes about one another.
	network.joinNodesToNetwork(dataDirs)

	return network
}

// IssueCommand sends the command via RPC to the nodeIdx'th node in the network.
func (network *GethNetwork) IssueCommand(nodeIdx int, command string) string {
	dataDir := network.dataDirs[nodeIdx]

	args := []string{dataDirFlag, dataDir, attachCmd, path.Join(dataDir, ipcFileName), execFlag, command}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stderr = network.logFile

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(output))
}

// Initialises and starts a Geth node.
func (network *GethNetwork) createMiner(dataDir string, port int) {
	// The node must create its initial config based on the network's genesis file before it can be started.
	network.initNode(dataDir)
	network.startMiner(dataDir, port)
}

// Initialises a Geth node based on the network genesis file.
func (network *GethNetwork) initNode(dataDirPath string) {
	args := []string{dataDirFlag, dataDirPath, initCmd, network.genesisFilePath}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// Starts a Geth node.
func (network *GethNetwork) startMiner(dataDirPath string, port int) {
	args := []string{websocketFlag, dataDirFlag, dataDirPath, fmt.Sprintf("%s=%d", portFlag, port), mineFlag, minerThreadsFlag, minerEthbaseFlag, rpcFeeCapFlag}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

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
	counter := 0
	for {
		ipcFilePath := path.Join(dataDir, ipcFileName)
		_, err := os.Stat(ipcFilePath)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)

		if counter > 20 {
			fmt.Printf("Waiting for .ipc file of node at %s", dataDir)
			counter = 0
		}
		counter++
	}
}
