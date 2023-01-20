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
	nodePorts                []int
}

func NewEth2Network(
	// gethBinaryPath string,
	// prysmBinaryPath string,
	httpPortStart int,
	// websocketPortStart int,
	numNodes int,
	blockTimeSecs int,
	// preFundedAddrs []string,
) *Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)
	binDir := path.Join(basepath, "../.build/eth2_bin")

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
	genesisStr, err := generateGenesis(blockTimeSecs, nil)
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

	dataDirs := make([]string, numNodes)
	nodePorts := make([]int, numNodes)
	for i := 0; i < numNodes; i++ {
		dataDirs[i] = path.Join(nodesDir, "node_datadir_"+strconv.Itoa(i+1))
		nodePorts[i] = httpPortStart + i
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
		nodePorts:                nodePorts,
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
			err := n.gethStartNode(8551+nodeID, 30303+nodeID, n.nodePorts[nodeID], dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// wait for each of the nodes to start
	for i := range n.dataDirs {
		nodeID := i
		eg.Go(func() error {
			return n.waitForNodeUp(nodeID, 15*time.Second)
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
			err = n.prysmStartBeaconNode(8551+nodeID, 12000+nodeID, 4000+nodeID, dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// start each of the validator nodes
	for _, nodeDataDir := range n.dataDirs {
		dataDir := nodeDataDir
		go func() {
			err = n.prysmStartValidator(dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	err = n.waitForMergeEvent(startTime)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)

	return nil
}

func (n *Eth2Network) Stop() {
}

func (n *Eth2Network) gethInitGenesisData(dataDirPath string) error {
	args := []string{_dataDirFlag, dataDirPath, "init", n.gethGenesisPath}
	fmt.Printf("gethInitGenesisData: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

//// Creates an account for a Geth node.
//func (n *Eth2Network) createAccount(dataDirPath string) error {
//	args := []string{dataDirFlag, dataDirPath, "account", "new", "--password", n.passwordPath}
//	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
//	cmd.Stdout = n.logFile
//	cmd.Stderr = n.logFile
//
//	return cmd.Run()
//}

func (n *Eth2Network) gethImportMinerAccount(nodeID int) error {
	args := []string{
		"--exec", fmt.Sprintf("loadScript('%s');", n.preloadScriptPath),
		"attach", fmt.Sprintf("http://127.0.0.1:%d", n.nodePorts[nodeID]),
	}
	fmt.Printf("gethImportMinerAccount: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethStartNode(executionPort, networkPort, httpPort int, dataDirPath string) error {
	args := []string{
		"--http", "--http.port", fmt.Sprintf("%d", httpPort), "--http.api", "miner,engine,personal,eth,net,web3,debug",
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

	return cmd.Run()
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

func (n *Eth2Network) prysmStartBeaconNode(executionPort, p2pPort, rpcPort int, nodeDataDir string) error {
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
		"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", executionPort),
		"--accept-terms-of-use",
		"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
	}

	fmt.Printf("prysmStartBeaconNode: %s %s\n", n.prysmBeaconBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBeaconBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) prysmStartValidator(nodeDataDir string) error {
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "validator"),
		"--accept-terms-of-use",
		"--interop-num-validators", "64",
		"--interop-start-index", "0",
		"--force-clear-db",
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
	}

	fmt.Printf("prysmStartValidator: %s %s\n", n.prysmValidatorBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmValidatorBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
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

	for ; number != 3; time.Sleep(time.Second) {
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
		dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.nodePorts[nodeID]))
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
