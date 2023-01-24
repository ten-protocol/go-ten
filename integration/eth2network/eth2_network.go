package eth2network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

const (
	_dataDirFlag                     = "--datadir"
	_eth2BinariesRelPath             = "../.build/eth2_bin"
	_gethFileNameVersion             = "geth-" + _gethVersion
	_prysmBeaconChainFileNameVersion = "beacon-chain-" + _prysmVersion
	_prysmCTLFileNameVersion         = "prysmctl-" + _prysmVersion
	_prysmValidatorFileNameVersion   = "validator-" + _prysmVersion
)

type Eth2Network struct {
	dataDirs                 []string
	buildDir                 string
	binDir                   string
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	gethGenesisPath          string
	preloadScriptPath        string
	prysmGenesisPath         string
	prysmConfigPath          string
	prysmValidatorBinaryPath string
	chainID                  int
	gethHTTPPorts            []int
	gethWSPorts              []int
	gethNetworkPorts         []int
	prysmExecPorts           []int
	gethProcesses            []*exec.Cmd
	prysmBeaconProcesses     []*exec.Cmd
	prysmValidatorProcesses  []*exec.Cmd
	logFile                  io.Writer
	preFundedMinerAddrs      []string
	preFundedMinerPKs        []string
}

func NewEth2Network(
	binDir string,
	gethHTTPPortStart int,
	gethWSPortStart int,
	prysmExecPortStart int,
	gethNetworkPortStart int,
	chainID int,
	numNodes int,
	blockTimeSecs int,
	preFundedAddrs []string,
) *Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// set the paths
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)
	gethGenesisPath := path.Join(buildDir, "genesis.json")
	gethPreloadScriptPath := path.Join(buildDir, "preload-script.js")
	prysmGenesisPath := path.Join(buildDir, "genesis.ssz")
	prysmConfigPath := path.Join(buildDir, "prysm_chain_config.yml")
	logPath := path.Join(buildDir, "node_logs.txt")

	gethBinaryPath := path.Join(binDir, _gethFileNameVersion)
	prysmBeaconBinaryPath := path.Join(binDir, _prysmBeaconChainFileNameVersion)
	prysmBinaryPath := path.Join(binDir, _prysmCTLFileNameVersion)
	prysmValidatorBinaryPath := path.Join(binDir, _prysmValidatorFileNameVersion)

	// Nodes logs and execution related files are written in the build folder
	err := os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Generate pk pairs for miners
	preFundedMinerAddrs := make([]string, numNodes)
	preFundedMinerPKs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		w := datagenerator.RandomWallet(int64(chainID))
		preFundedMinerAddrs[i] = w.Address().Hex()
		preFundedMinerPKs[i] = fmt.Sprintf("%x", w.PrivateKey().D.Bytes())
	}
	// Generate and write genesis file
	genesisStr, err := generateGenesis(blockTimeSecs, chainID, preFundedMinerAddrs, append(preFundedAddrs, preFundedMinerAddrs...))
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(gethGenesisPath, []byte(genesisStr), 0o600)
	if err != nil {
		panic(err)
	}

	// Write beacon config
	err = os.WriteFile(prysmConfigPath, []byte(beaconConfig), 0o600)
	if err != nil {
		panic(err)
	}

	// Write geth js script
	err = os.WriteFile(gethPreloadScriptPath, []byte(gethPreloadJSONScript), 0o600)
	if err != nil {
		panic(err)
	}

	// Each node has a temp directory
	nodesDir, err := os.MkdirTemp("", timestamp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Geth nodes created in: %s\n", nodesDir)

	gethProcesses := make([]*exec.Cmd, numNodes)
	prysmBeaconProcesses := make([]*exec.Cmd, numNodes)
	prysmValidatorProcesses := make([]*exec.Cmd, numNodes)
	dataDirs := make([]string, numNodes)
	gethHTTPPorts := make([]int, numNodes)
	gethWSPorts := make([]int, numNodes)
	prysmExecPorts := make([]int, numNodes)
	gethNetworkPorts := make([]int, numNodes)
	for i := 0; i < numNodes; i++ {
		dataDirs[i] = path.Join(nodesDir, "node_datadir_"+strconv.Itoa(i+1))
		gethHTTPPorts[i] = gethHTTPPortStart + i
		gethWSPorts[i] = gethWSPortStart + i
		prysmExecPorts[i] = prysmExecPortStart + i
		gethNetworkPorts[i] = gethNetworkPortStart + i
	}

	// create the log file
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}

	return &Eth2Network{
		buildDir:                 buildDir,
		binDir:                   binDir,
		dataDirs:                 dataDirs,
		chainID:                  chainID,
		gethProcesses:            gethProcesses,
		prysmBeaconProcesses:     prysmBeaconProcesses,
		prysmValidatorProcesses:  prysmValidatorProcesses,
		gethHTTPPorts:            gethHTTPPorts,
		gethWSPorts:              gethWSPorts,
		gethNetworkPorts:         gethNetworkPorts,
		prysmExecPorts:           prysmExecPorts,
		gethBinaryPath:           gethBinaryPath,
		prysmBinaryPath:          prysmBinaryPath,
		prysmBeaconBinaryPath:    prysmBeaconBinaryPath,
		prysmConfigPath:          prysmConfigPath,
		prysmValidatorBinaryPath: prysmValidatorBinaryPath,
		gethGenesisPath:          gethGenesisPath,
		prysmGenesisPath:         prysmGenesisPath,
		logFile:                  logFile,
		preloadScriptPath:        gethPreloadScriptPath,
		preFundedMinerAddrs:      preFundedMinerAddrs,
		preFundedMinerPKs:        preFundedMinerPKs,
	}
}

// Start starts the network
func (n *Eth2Network) Start() error {
	startTime := time.Now()
	var eg errgroup.Group

	// initialize the genesis data on the nodes
	for _, nodeDataDir := range n.dataDirs {
		dataDir := nodeDataDir
		eg.Go(func() error {
			return n.gethInitGenesisData(dataDir)
		})
	}
	err := eg.Wait()
	if err != nil {
		return err
	}

	// start each of the nodes
	for i, nodeDataDir := range n.dataDirs {
		dataDir := nodeDataDir
		nodeID := i
		go func() {
			n.gethProcesses[nodeID], err = n.gethStartNode(
				n.prysmExecPorts[nodeID],
				n.gethNetworkPorts[nodeID],
				n.gethHTTPPorts[nodeID],
				n.gethWSPorts[nodeID],
				dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// wait for each of the nodes to start
	for i := range n.dataDirs {
		nodeID := i
		eg.Go(func() error {
			return n.waitForNodeUp(nodeID, 30*time.Second)
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	// link nodes together by providing node 0 the enode (eth node address) of the other nodes
	var enodes []string
	for i := 1; i < len(n.dataDirs); i++ {
		enode, err := n.gethGetEnode(i)
		if err != nil {
			return err
		}
		enodes = append(enodes, enode)
	}
	err = n.gethImportEnodes(enodes)
	if err != nil {
		return err
	}

	// import prefunded key to each node and start mining
	for i := range n.dataDirs {
		nodeID := i
		eg.Go(func() error {
			return n.gethImportMinerAccount(nodeID)
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	// generate the genesis using the node 0
	err = n.prysmGenerateGenesis()
	if err != nil {
		return err
	}

	// start each of the beacon nodes
	for i, nodeDataDir := range n.dataDirs {
		nodeID := i
		dataDir := nodeDataDir
		go func() {
			n.prysmBeaconProcesses[nodeID], err = n.prysmStartBeaconNode(n.prysmExecPorts[nodeID], 12000+nodeID, 4000+nodeID, dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// start each of the validator nodes
	for i, nodeDataDir := range n.dataDirs {
		nodeID := i
		dataDir := nodeDataDir
		go func() {
			n.prysmValidatorProcesses[nodeID], err = n.prysmStartValidator(4000+nodeID, dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// this locks the process waiting for the event to happen
	return n.waitForMergeEvent(startTime)
}

func (n *Eth2Network) Stop() {
	for i := 0; i < len(n.dataDirs); i++ {
		err := n.gethProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill geth node - %s\n", err.Error())
		}
		err = n.prysmBeaconProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill prysm beacon node - %s\n", err.Error())
		}
		err = n.prysmValidatorProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill prysm validator node - %s\n", err.Error())
		}
	}
}

func (n *Eth2Network) gethInitGenesisData(dataDirPath string) error {
	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{_dataDirFlag, dataDirPath, "init", n.gethGenesisPath}
	fmt.Printf("gethInitGenesisData: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethImportMinerAccount(nodeID int) error {
	startScript := fmt.Sprintf(gethStartupScriptJS, n.preFundedMinerPKs[nodeID])

	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		"--exec", fmt.Sprintf("%s", startScript),
		"attach", fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[nodeID]),
	}
	fmt.Printf("gethImportMinerAccount: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethStartNode(executionPort, networkPort, httpPort, wsPort int, dataDirPath string) (*exec.Cmd, error) {
	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		_dataDirFlag, dataDirPath,
		"--http",
		"--http.port", fmt.Sprintf("%d", httpPort),
		"--http.api", "admin,miner,engine,personal,eth,net,web3,debug",
		"--ws",
		"--ws.port", fmt.Sprintf("%d", wsPort),
		"--authrpc.port", fmt.Sprintf("%d", executionPort),
		"--port", fmt.Sprintf("%d", networkPort),
		"--networkid", fmt.Sprintf("%d", n.chainID),
		"--syncmode", "full", // sync mode to download and test all blocks and txs
		"--allow-insecure-unlock", // allows to use personal accounts over http/ws
		"--nodiscover",            // don't try and discover peers
	}
	fmt.Printf("gethStartNode: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd, cmd.Start()
}

func (n *Eth2Network) prysmGenerateGenesis() error {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"testnet", "generate-genesis",
		"--num-validators", "64", "--output-ssz", n.prysmGenesisPath,
		"--chain-config-file", n.prysmConfigPath,
	}
	fmt.Printf("prysmGenerateGenesis: %s %s\n", n.prysmBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) prysmStartBeaconNode(executionPort, p2pPort, rpcPort int, nodeDataDir string) (*exec.Cmd, error) {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "beacondata"),
		"--rpc-port", fmt.Sprintf("%d", rpcPort),
		"--p2p-udp-port", fmt.Sprintf("%d", p2pPort),
		"--min-sync-peers", "0",
		"--interop-genesis-state", n.prysmGenesisPath,
		"--interop-eth1data-votes",
		"--bootstrap-node", "",
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
		"--chain-id", "32382",
		"--grpc-gateway-port", fmt.Sprintf("%d", rpcPort+100),
		"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", executionPort),
		"--accept-terms-of-use",
		"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
	}

	fmt.Printf("prysmStartBeaconNode: %s %s\n", n.prysmBeaconBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBeaconBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd, cmd.Start()
}

func (n *Eth2Network) prysmStartValidator(rpcPort int, nodeDataDir string) (*exec.Cmd, error) {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "validator"),
		"--accept-terms-of-use",
		"--interop-num-validators", "64",
		"--interop-start-index", "0",
		"--force-clear-db",
		"--beacon-rpc-gateway-provider", fmt.Sprintf("%d", rpcPort+100),
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
	}

	fmt.Printf("prysmStartValidator: %s %s\n", n.prysmValidatorBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmValidatorBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd, cmd.Start()
}

// waitForMergeEvent connects to the geth node and waits until block 2 (the merge block) is reached
func (n *Eth2Network) waitForMergeEvent(startTime time.Time) error {
	dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[0]))
	if err != nil {
		return err
	}
	number, err := dial.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	for ; number != 2; time.Sleep(time.Second) {
		number, err = dial.BlockNumber(context.Background())
		if err != nil {
			return err
		}
	}

	fmt.Printf("Reached the merge block after %s\n", time.Since(startTime))
	return nil
}

// waitForNodeUp retries continuously for the node to respond to a http request
func (n *Eth2Network) waitForNodeUp(nodeID int, timeout time.Duration) error {
	for startTime := time.Now(); time.Now().Before(startTime.Add(timeout)); time.Sleep(time.Second) {
		dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[nodeID]))
		if err != nil {
			continue
		}
		_, err = dial.BlockNumber(context.Background())
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("node not responsive after %s", timeout)
}

func (n *Eth2Network) gethGetEnode(i int) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[i]),
		strings.NewReader(`{"jsonrpc": "2.0", "method": "admin_nodeInfo", "params": [], "id": 1}`))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	return res["result"].(map[string]interface{})["enode"].(string), nil
}

func (n *Eth2Network) gethImportEnodes(enodes []string) error {
	for _, enode := range enodes {
		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[0]),
			strings.NewReader(
				fmt.Sprintf(`{"jsonrpc": "2.0", "method": "admin_addPeer", "params": ["%s"], "id": 1}`, enode),
			),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	return nil
}
