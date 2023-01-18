package gethnetwork

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

const (
	dataDirFlag = "--datadir"
)

type Eth2Network struct {
	buildDir                 string
	dataDirs                 []string
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	logFile                  io.Writer
	gethGenesisPath          string
	binDir                   string
	preloadScriptPath        string
	prysmGenesisPath         string
	prysmConfigPath          string
	prysmValidatorBinaryPath string
}

func NewEth2Network(
	//gethBinaryPath string,
	//prysmBinaryPath string,
	//portStart int,
	//websocketPortStart int,
	numNodes int,
	//blockTimeSecs int,
	//preFundedAddrs []string,
) *Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)
	binDir := path.Join(basepath, "../.build/eth2_bin")

	gethGenesisPath := path.Join(binDir, "genesis.json")
	prysmGenesisPath := path.Join(binDir, "genesis.ssz")
	gethBinaryPath := path.Join(binDir, "geth-v1.10.26")
	prysmBeaconBinaryPath := path.Join(binDir, "beacon-chain-v3.2.0-darwin-arm64")
	prysmBinaryPath := path.Join(binDir, "prysmctl-v3.2.0-darwin-arm64")
	prysmConfigPath := path.Join(binDir, "prysm_chain_config.yml")
	prysmValidatorBinaryPath := path.Join(binDir, "validator-v3.2.0-darwin-arm64")
	preloadScriptPath := path.Join(binDir, "preload-script.js")

	// Each node has a temp directory
	nodesDir, err := os.MkdirTemp("", timestamp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Geth nodes created in: %s\n", nodesDir)

	dataDirs := make([]string, numNodes)
	for i := 0; i < numNodes; i++ {
		dataDirs[i] = path.Join(nodesDir, "node_datadir_"+strconv.Itoa(i+1))
	}

	// Nodes logs are written to the build directory
	err = os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// TODO HOOK LOGS
	//logPath := path.Join(buildDir, "node_logs.txt")
	//logFile, err := os.Create(logPath)
	//if err != nil {
	//	panic(err)
	//}

	return &Eth2Network{
		buildDir:                 buildDir,
		binDir:                   binDir,
		dataDirs:                 dataDirs,
		gethBinaryPath:           gethBinaryPath,
		prysmBinaryPath:          prysmBinaryPath,
		prysmBeaconBinaryPath:    prysmBeaconBinaryPath,
		prysmConfigPath:          prysmConfigPath,
		prysmValidatorBinaryPath: prysmValidatorBinaryPath,
		gethGenesisPath:          gethGenesisPath,
		prysmGenesisPath:         prysmGenesisPath,
		logFile:                  os.Stdout,
		preloadScriptPath:        preloadScriptPath,
	}
}

func (n *Eth2Network) Start() error {

	// initialize the genesis data on the node
	err := n.gethInitGenesisData(n.dataDirs[0])
	if err != nil {
		return err
	}

	// start the node
	go func() {
		err = n.gethStartNode(n.dataDirs[0])
		if err != nil {
			panic(err)
		}
	}()

	// TODO dont use sleep ensure the node is up instead with polling
	time.Sleep(15 * time.Second)

	// import miner account that helps to reach to POS
	err = n.gethImportMinerAccount()
	if err != nil {
		return err
	}

	err = n.prysmGenerateGenesis()
	if err != nil {
		return err
	}

	err = n.prysmStartBeaconNode(n.dataDirs[0])
	if err != nil {
		return err
	}

	err = n.prysmStartValidator(n.dataDirs[0])
	if err != nil {
		return err
	}

	return nil
}

func (n *Eth2Network) Stop() {

}

func (n *Eth2Network) gethInitGenesisData(dataDirPath string) error {
	args := []string{dataDirFlag, dataDirPath, "init", n.gethGenesisPath}
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

func (n *Eth2Network) gethImportMinerAccount() error {
	args := []string{
		"--exec", fmt.Sprintf("loadScript('%s');", n.preloadScriptPath),
		"attach", "http://127.0.0.1:8545",
	}
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) gethStartNode(dataDirPath string) error {
	args := []string{
		"--http", "--http.api", "miner,engine,personal,eth,net,web3,debug",
		dataDirFlag, dataDirPath,
		"--allow-insecure-unlock",
		"--networkid", "1999",
	}
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
	cmd := exec.Command(n.prysmBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()
}

func (n *Eth2Network) prysmStartBeaconNode(nodeDataDir string) error {
	args := []string{
		"--datadir", "beacondata",
		"--min-sync-peers", "0",
		"--interop-genesis-state", n.prysmGenesisPath,
		"--interop-eth1data-votes",
		"--bootstrap-node", "",
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
		"--chain-id", "32382",
		"--execution-endpoint", "http://localhost:8551",
		"--accept-terms-of-use",
		"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
	}
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
	cmd := exec.Command(n.prysmValidatorBinaryPath, args...) //nolint
	cmd.Stdout = n.logFile
	cmd.Stderr = n.logFile

	return cmd.Run()

}
