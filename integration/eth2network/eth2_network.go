package eth2network

import (
	"context"
	"fmt"
	"io"
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
	_gethFileNameVersion             = "geth-v1.10.26"
	_prysmBeaconChainFileNameVersion = "beacon-chain-v3.2.0"
	_prysmCTLFileNameVersion         = "prysmctl-v3.2.0"
	_prysmValidatorFileNameVersion   = "validator-v3.2.0"
)

type Eth2Network struct {
	dataDirs                 []string
	buildDir                 string
	binDir                   string
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	logFile                  io.Writer
	gethGenesisPath          string
	preloadScriptPath        string
	prysmGenesisPath         string
	prysmConfigPath          string
	prysmValidatorBinaryPath string
	gethHTTPPorts            []int
	gethWSPorts              []int
	gethNetworkPorts         []int
	prysmExecPorts           []int
	gethProcesses            []*exec.Cmd
	prysmBeaconProcesses     []*exec.Cmd
	prysmValidatorProcesses  []*exec.Cmd
}

func NewEth2Network(
	binDir string,
	gethHTTPPortStart int,
	gethWSPortStart int,
	prysmExecPortStart int,
	gethNetworkPortStart int,
	numNodes int,
	blockTimeSecs int,
	preFundedAddrs []string,
) *Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
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

	// Nodes logs and execution related files are writen in the build folder
	err := os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Generate and write genesis file
	genesisStr, err := generateGenesis(blockTimeSecs, preFundedAddrs)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(gethGenesisPath, []byte(genesisStr), 0777)
	if err != nil {
		panic(err)
	}

	// Write beacon config
	err = os.WriteFile(prysmConfigPath, []byte(beaconConfig), 0777)
	if err != nil {
		panic(err)
	}

	// Write geth js script
	err = os.WriteFile(gethPreloadScriptPath, []byte(gethPreloadJsonScript), 0777)
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
	}
}

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
		panic(err)
	}

	// import miner account that helps to reach to POS on node 0
	err = n.gethImportMinerAccount(0)
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
	args := []string{_dataDirFlag, dataDirPath, "init", n.gethGenesisPath}
	fmt.Printf("gethInitGenesisData: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethImportMinerAccount(nodeID int) error {
	args := []string{
		"--exec", fmt.Sprintf("loadScript('%s');", n.preloadScriptPath),
		"attach", fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[nodeID]),
	}
	fmt.Printf("gethImportMinerAccount: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethStartNode(executionPort, networkPort, httpPort, wsPort int, dataDirPath string) (*exec.Cmd, error) {
	args := []string{
		"--http", "--http.port", fmt.Sprintf("%d", httpPort), "--http.api", "miner,engine,personal,eth,net,web3,debug",
		"--ws", "--ws.port", fmt.Sprintf("%d", wsPort),
		"--authrpc.port", fmt.Sprintf("%d", executionPort),
		"--port", fmt.Sprintf("%d", networkPort),
		_dataDirFlag, dataDirPath,
		"--allow-insecure-unlock",
		"--nodiscover", "--syncmode", "full",
		"--networkid", "32382",
	}
	fmt.Printf("gethStartNode: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd, cmd.Start()
}

func (n *Eth2Network) prysmGenerateGenesis() error {
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

func (n *Eth2Network) waitForMergeEvent(startTime time.Time) error {
	dial, err := ethclient.Dial("http://127.0.0.1:8545")
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
