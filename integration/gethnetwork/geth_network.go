package gethnetwork

import (
	"encoding/json"
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

	accountCmd    = "account"
	accountNewCmd = "new"
	addPeerCmd    = "admin.addPeer(%s)"
	attachCmd     = "attach"
	enodeCmd      = "admin.nodeInfo.enode"
	initCmd       = "init"

	websocketFlag = "--ws" // Enables websocket connections to the node.
	dataDirFlag   = "--datadir"
	execFlag      = "--exec"
	portFlag      = "--port"
	mineFlag      = "--mine"
	rpcFeeCapFlag = "--rpc.txfeecap=0" // Disables the 1 ETH cap for RPC transactions.

	// We pre-allocate a single wallet, the miner etherbase account above.
	genesisJSONTemplate = `{
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
		"londonBlock": 0,
		"clique": {
		  "period": 5,
		  "epoch": 30000
		}
	  },
	  "alloc": {
		"0x323AefbFC16159655514846a9e5433C457de9389": {
		  "balance": "10000000000"
		}
	  },
	  "coinbase": "0x0000000000000000000000000000000000000000",
	  "difficulty": "0x20000",
	  "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000%s0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	  "gasLimit": "0x77359400",
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
	addresses       []string // The public keys of the nodes' accounts.
	logFile         *os.File
}

// NewGethNetwork using the provided Geth binary to create a private Ethereum network with numNodes Geth nodes.
func NewGethNetwork(gethBinaryPath string, numNodes int) GethNetwork {
	// We create a data directory for each Geth node.
	nodesDir, err := ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		panic(err)
	}
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
		gethBinaryPath: gethBinaryPath,
		dataDirs:       dataDirs,
		logFile:        logFile,
	}

	for _, dataDir := range dataDirs {
		network.createAccount(dataDir)
		network.retrieveAccount(dataDir)
	}
	genesisJSON := fmt.Sprintf(genesisJSONTemplate, strings.Join(network.addresses, ""))

	// We write out the `genesis.json` file to be used by the network.
	genesisFilePath := path.Join(nodesDir, genesisFileName)
	err = os.WriteFile(genesisFilePath, []byte(genesisJSON), 0o600)
	if err != nil {
		panic(err)
	}
	network.genesisFilePath = genesisFilePath

	// todo - joel - fold loop into method
	for idx, dataDir := range dataDirs {
		_ = os.Remove(path.Join(dataDir, ipcFileName)) // We delete leftover IPC files from previous runs.
		go network.createMiner(dataDir, idx)
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
func (network *GethNetwork) createMiner(dataDir string, idx int) {
	// The node must create its initial config based on the network's genesis file before it can be started.
	network.initNode(dataDir)
	network.startMiner(dataDir, idx)
}

// Creates an account for a Geth node.
func (network *GethNetwork) createAccount(dataDirPath string) {
	args := []string{dataDirFlag, dataDirPath, accountCmd, accountNewCmd}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	// todo - joel - pass password via command line
	x, _ := cmd.StdinPipe()
	x.Write([]byte("wefwefwef\r\n"))
	x.Write([]byte("wefwefwef\r\n"))

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// Adds a Geth node's account public key to the `network` object.
func (network *GethNetwork) retrieveAccount(dataDirPath string) {
	dir := path.Join(dataDirPath, "keystore") // todo - joel - use constant
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		// `ReadDir` returns the folder itself, as well as the files within it.
		if file.IsDir() {
			continue
		}

		y, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			panic(err)
		}
		contents := make(map[string]interface{})
		err = json.Unmarshal(y, &contents)
		if err != nil {
			panic(err)
		}
		network.addresses = append(network.addresses, contents["address"].(string)) // todo - joel - use constant
		return                                                                      // We return as we only expect one account per node.
	}

	panic(fmt.Sprintf("could not find account for node at %s", dataDirPath))
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
// todo - joel - rename miner references
func (network *GethNetwork) startMiner(dataDirPath string, idx int) {
	args := []string{websocketFlag, dataDirFlag, dataDirPath, fmt.Sprintf("%s=%d", portFlag, startPort+idx), "--allow-insecure-unlock", "--unlock", network.addresses[idx], mineFlag, rpcFeeCapFlag} // todo - joel - use constants
	cmd := exec.Command(network.gethBinaryPath, args...)                                                                                                                                             // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	// todo - joel - pass password via command line
	x, _ := cmd.StdinPipe()
	x.Write([]byte("wefwefwef\r\n"))
	x.Write([]byte("wefwefwef\r\n"))

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
