package gethnetwork

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	nodeFolderName = "node_datadir_"
	tempDirPrefix  = "geth_nodes"
	buildDir       = "../.build/geth"
	keystoreDir    = "keystore"

	genesisFileName = "genesis.json"
	ipcFileName     = "geth.ipc"
	logFile         = "node_logs.txt"
	passwordFile    = "password.txt"

	startPort   = 30303
	wsStartPort = 8546
	password    = "password"

	accountCmd    = "account"
	accountNewCmd = "new"
	addPeerCmd    = "admin.addPeer(%s)"
	attachCmd     = "attach"
	enodeCmd      = "admin.nodeInfo.enode"
	initCmd       = "init"

	dataDirFlag        = "--datadir"
	execFlag           = "--exec"
	mineFlag           = "--mine"
	passwordFlag       = "--password"
	portFlag           = "--port"
	rpcFeeCapFlag      = "--rpc.txfeecap=0" // Disables the 1 ETH cap for RPC transactions.
	unlockFlag         = "--unlock"
	unlockInsecureFlag = "--allow-insecure-unlock"
	websocketFlag      = "--ws" // Enables websocket connections to the node.
	wsPortFlag         = "--ws.port"

	// We pre-allocate a wallet matching the private key used in the tests, plus an account per clique member.
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
		  "period": 0,
		  "epoch": 30000
		}
	  },
	  "alloc": {
		"0x323AefbFC16159655514846a9e5433C457de9389": {
		  "balance": "1000000000000000000000"
		},
%s
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
	allocBlockTemplate = `		"0x%s": {
		  "balance": "1000000000000000000000"
		}`
	genesisJSONAddrKey = "address"
)

// GethNetwork is a network of Geth nodes, built using the provided Geth binary.
type GethNetwork struct {
	gethBinaryPath   string
	genesisFilePath  string
	dataDirs         []string
	addresses        []string // The public keys of the nodes' accounts.
	logFile          *os.File
	passwordFilePath string // The path to the file storing the password to unlock node accounts.
}

// NewGethNetwork using the provided Geth binary to create a private Ethereum network with numNodes Geth nodes.
func NewGethNetwork(gethBinaryPath string, numNodes int) GethNetwork {
	// We create a data directory for each node.
	nodesDir, err := ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		panic(err)
	}
	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		nodeFolder := nodeFolderName + strconv.Itoa(i+1)
		dataDirs[i] = path.Join(nodesDir, nodeFolder)
	}

	// We push all the node logs to a single file.
	err = os.MkdirAll(buildDir, 0o700)
	if err != nil {
		panic(err)
	}
	logFile, err := os.Create(path.Join(buildDir, logFile))
	if err != nil {
		panic(err)
	}

	// We create a password file to unlock the node accounts.
	passwordFile, _ := os.Create(path.Join(nodesDir, passwordFile))
	if err != nil {
		panic(err)
	}
	_, err = passwordFile.WriteString(password)
	if err != nil {
		panic(err)
	}

	network := GethNetwork{
		gethBinaryPath:   gethBinaryPath,
		dataDirs:         dataDirs,
		logFile:          logFile,
		passwordFilePath: passwordFile.Name(),
	}

	// We create an account for each node and generate the genesis config file accordingly.
	for _, dataDir := range dataDirs {
		network.createAccount(dataDir)
		accountAddress := network.retrieveAccount(dataDir)
		network.addresses = append(network.addresses, accountAddress)
	}
	allocs := make([]string, len(network.addresses))
	for idx, addr := range network.addresses {
		allocs[idx] = fmt.Sprintf(allocBlockTemplate, addr)
	}
	genesisJSON := fmt.Sprintf(genesisJSONTemplate, strings.Join(allocs, ",\r\n"), strings.Join(network.addresses, ""))

	// We write out the `genesis.json` file to be used by the network.
	genesisFilePath := path.Join(buildDir, genesisFileName)
	err = os.WriteFile(genesisFilePath, []byte(genesisJSON), 0o600)
	if err != nil {
		panic(err)
	}
	network.genesisFilePath = genesisFilePath

	// We start the miners.
	for idx, dataDir := range dataDirs {
		go network.createMiner(dataDir, idx)
	}

	// We need to manually tell the nodes about one another.
	network.joinNodesToNetwork()

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

// Initialises and starts a miner.
func (network *GethNetwork) createMiner(dataDir string, idx int) {
	// We delete the leftover IPC file from the previous run, if it exists.
	_ = os.Remove(path.Join(dataDir, ipcFileName))
	// The node must create its initial config based on the network's genesis file before it can be started.
	network.initNode(dataDir)
	network.startMiner(dataDir, idx)
}

// Creates an account for a Geth node.
func (network *GethNetwork) createAccount(dataDirPath string) {
	args := []string{dataDirFlag, dataDirPath, accountCmd, accountNewCmd, passwordFlag, network.passwordFilePath}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// Adds a Geth node's account public key to the `network` object.
func (network *GethNetwork) retrieveAccount(dataDirPath string) string {
	dir := path.Join(dataDirPath, keystoreDir)
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
		return contents[genesisJSONAddrKey].(string) // We return as we only expect one account per node.
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

// Starts a Geth miner.
func (network *GethNetwork) startMiner(dataDirPath string, idx int) {
	args := []string{
		websocketFlag, wsPortFlag, strconv.Itoa(wsStartPort + idx), dataDirFlag, dataDirPath, portFlag,
		strconv.Itoa(startPort + idx), unlockInsecureFlag, unlockFlag, network.addresses[idx], passwordFlag,
		network.passwordFilePath, mineFlag, rpcFeeCapFlag,
	}
	cmd := exec.Command(network.gethBinaryPath, args...) // nolint
	cmd.Stdout = network.logFile
	cmd.Stderr = network.logFile

	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

// Tells the network's nodes about one another.
func (network *GethNetwork) joinNodesToNetwork() {
	enodeAddrs := make([]string, len(network.dataDirs))

	for i, dataDir := range network.dataDirs {
		waitForIPC(dataDir) // We cannot issue RPC commands until the IPC files are available.
		enodeAddrs[i] = network.IssueCommand(i, enodeCmd)
	}

	for _, enodeAddr := range enodeAddrs {
		for i := range network.dataDirs {
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
