package eth2network

import (
	"context"
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"golang.org/x/sync/errgroup"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	_eth2BinariesRelPath = "../.build/eth2_bin"
	_dataDirFlag         = "--datadir"
	_gethBinaryName      = "geth"
)

// Docker service names
const (
	dockerComposeFile    = "docker-compose.yml"
	dockerNetworkName    = "polygon-net"
	beaconServiceName    = "beacon-chain"
	gethServiceName      = "geth"
	validatorServiceName = "validator"
)

type Impl struct {
	dataDirs                 []string
	buildDir                 string
	binDir                   string
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	gethGenesisPath          string
	prysmGenesisPath         string
	prysmConfigPath          string
	prysmValidatorBinaryPath string
	chainID                  int
	gethHTTPPorts            []int
	gethWSPorts              []int
	gethNetworkPorts         []int
	gethAuthRPCPorts         []int
	prysmBeaconHTTPPorts     []int
	prysmBeaconP2PPorts      []int
	gethProcesses            []*exec.Cmd
	prysmBeaconProcesses     []*exec.Cmd
	prysmValidatorProcesses  []*exec.Cmd
	gethLogFile              io.Writer
	prysmBeaconLogFile       io.Writer
	prysmValidtorLogFile     io.Writer
	preFundedMinerAddrs      []string
	preFundedMinerPKs        []string
	gethGenesisBytes         []byte
	timeout                  time.Duration
}

type Eth2Network interface {
	GethGenesis() []byte
	Start() error
	Stop() error
}

func NewEth2Network(
	binDir string,
	logToFile bool,
	gethHTTPPortStart int,
	gethWSPortStart int,
	gethAuthRPCPortStart int,
	gethNetworkPortStart int,
	prysmBeaconHTTPPortStart int,
	prysmBeaconP2PPortStart int,
	chainID int,
	numNodes int,
	blockTimeSecs int,
	slotsPerEpoch int,
	secondsPerSlot int,
	preFundedAddrs []string,
	timeout time.Duration,
) Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMicro(), 10)

	// set the paths
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)
	gethGenesisPath := path.Join(buildDir, "genesis.json")
	prysmGenesisPath := path.Join(buildDir, "genesis.ssz")
	prysmConfigPath := path.Join(buildDir, "prysm_chain_config.yml")

	gethBinaryPath := path.Join(binDir, _gethFileNameVersion, _gethBinaryName)
	prysmBeaconBinaryPath := path.Join(binDir, _prysmBeaconChainFileNameVersion)
	prysmBinaryPath := path.Join(binDir, _prysmCTLFileNameVersion)
	prysmValidatorBinaryPath := path.Join(binDir, _prysmValidatorFileNameVersion)

	// catch any issues due to folder collision early
	if _, err := os.Stat(buildDir); err == nil {
		panic(fmt.Sprintf("folder %s already exists", buildDir))
	}

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
	beaconConf := fmt.Sprintf(_beaconConfig, chainID, chainID, secondsPerSlot, slotsPerEpoch)
	err = os.WriteFile(prysmConfigPath, []byte(beaconConf), 0o600)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Geth nodes created in: %s\n", buildDir)

	gethProcesses := make([]*exec.Cmd, numNodes)
	prysmBeaconProcesses := make([]*exec.Cmd, numNodes)
	prysmValidatorProcesses := make([]*exec.Cmd, numNodes)
	dataDirs := make([]string, numNodes)
	gethHTTPPorts := make([]int, numNodes)
	gethWSPorts := make([]int, numNodes)
	gethAuthRPCPorts := make([]int, numNodes)
	gethNetworkPorts := make([]int, numNodes)
	prysmBeaconHTTPPorts := make([]int, numNodes)
	prysmBeaconP2PPorts := make([]int, numNodes)

	for i := 0; i < numNodes; i++ {
		dataDirs[i] = path.Join(buildDir, "n"+strconv.Itoa(i))
		gethHTTPPorts[i] = gethHTTPPortStart + i
		gethWSPorts[i] = gethWSPortStart + i
		gethAuthRPCPorts[i] = gethAuthRPCPortStart + i
		gethNetworkPorts[i] = gethNetworkPortStart + i
		prysmBeaconHTTPPorts[i] = prysmBeaconHTTPPortStart + i
		prysmBeaconP2PPorts[i] = prysmBeaconP2PPortStart + i
	}

	// create the log files
	gethLogFile := io.Writer(os.Stdout)
	prysmBeaconLogFile := io.Writer(os.Stdout)
	prysmValidatorLogFile := io.Writer(os.Stdout)

	if logToFile {
		gethLogFile, err = NewRotatingLogWriter(buildDir, "geth_logs", 10*1024*1024, 5)
		if err != nil {
			panic(err)
		}
		prysmBeaconLogFile, err = NewRotatingLogWriter(buildDir, "prysm_beacon_logs", 10*1024*1024, 5)
		if err != nil {
			panic(err)
		}
		prysmValidatorLogFile, err = NewRotatingLogWriter(buildDir, "prysm_validator_logs", 10*1024*1024, 5)
		if err != nil {
			panic(err)
		}
	}

	return &Impl{
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
		gethAuthRPCPorts:         gethAuthRPCPorts,
		prysmBeaconHTTPPorts:     prysmBeaconHTTPPorts,
		prysmBeaconP2PPorts:      prysmBeaconP2PPorts,
		gethBinaryPath:           gethBinaryPath,
		prysmBinaryPath:          prysmBinaryPath,
		prysmBeaconBinaryPath:    prysmBeaconBinaryPath,
		prysmConfigPath:          prysmConfigPath,
		prysmValidatorBinaryPath: prysmValidatorBinaryPath,
		gethGenesisPath:          gethGenesisPath,
		prysmGenesisPath:         prysmGenesisPath,
		gethLogFile:              gethLogFile,
		prysmBeaconLogFile:       prysmBeaconLogFile,
		prysmValidtorLogFile:     prysmValidatorLogFile,
		preFundedMinerAddrs:      preFundedMinerAddrs,
		preFundedMinerPKs:        preFundedMinerPKs,
		gethGenesisBytes:         []byte(genesisStr),
		timeout:                  timeout,
	}
}

// Start starts the network
func (n *Impl) Start() error {
	startTime := time.Now()
	var eg errgroup.Group

	if err := n.ensureNoDuplicatedNetwork(); err != nil {
		return err
	}

	// Run docker compose to start the network
	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "up", "-d")
	fmt.Printf("Starting Docker Compose: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker Compose: %w", err)
	}

	// Wait for the nodes to start
	for i := range n.dataDirs {
		nodeID := i
		eg.Go(func() error {
			return n.waitForNodeUp(nodeID, time.Minute)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	// Generate the genesis using the node 0
	if err := n.prysmGenerateGenesis(); err != nil {
		return err
	}

	// Blocking wait until the network reaches the Merge
	return n.waitForMergeEvent(startTime)
}

// Stop stops the network
func (n *Impl) Stop() error {
	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "down")
	fmt.Printf("Stopping Docker Compose: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Docker Compose: %w", err)
	}
	// Wait a second for the kill signal
	time.Sleep(time.Second)
	return nil
}

// waitForNodeUp retries continuously for the node to respond to an HTTP request
func (n *Impl) waitForNodeUp(nodeID int, timeout time.Duration) error {
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
	fmt.Printf("Geth node error:\n%s\n", n.gethProcesses[nodeID].Stderr)
	return fmt.Errorf("node not responsive after %s", timeout)
}

// prysmGenerateGenesis generates the genesis for the beacon chain
func (n *Impl) prysmGenerateGenesis() error {
	// Full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"testnet", "generate-genesis",
		"--num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		"--output-ssz", n.prysmGenesisPath,
		"--config-name", "interop",
		"--chain-config-file", n.prysmConfigPath,
	}
	fmt.Printf("prysmGenerateGenesis: %s %s\n", n.prysmBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBinaryPath, args...) //nolint
	cmd.Stdout = n.prysmBeaconLogFile
	cmd.Stderr = n.prysmBeaconLogFile

	return cmd.Run()
}

// waitForMergeEvent connects to the geth node and waits until block 2 (the merge block) is reached
func (n *Impl) waitForMergeEvent(startTime time.Time) error {
	ctx := context.Background()
	dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[0]))
	if err != nil {
		return err
	}
	number, err := dial.BlockNumber(ctx)
	if err != nil {
		return err
	}

	// Wait for the merge block
	err = retry.Do(
		func() error {
			number, err = dial.BlockNumber(ctx)
			if err != nil {
				return err
			}
			if number <= 7 {
				return fmt.Errorf("has not arrived at The Merge")
			}
			return nil
		},
		retry.NewTimeoutStrategy(n.timeout, time.Second),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Reached the merge block after %s\n", time.Since(startTime))

	if err = n.prefundedBalancesActive(dial); err != nil {
		fmt.Printf("Error prefunding accounts %s\n", err.Error())
		return err
	}
	return nil
}

// ensureNoDuplicatedNetwork ensures no duplicated network
func (n *Impl) ensureNoDuplicatedNetwork() error {
	for nodeIdx, port := range n.gethWSPorts {
		_, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", port))
		if err == nil {
			return fmt.Errorf("unexpected geth node %d is active before the network is started", nodeIdx)
		}
	}
	return nil
}

// prefundedBalancesActive checks prefunded balances are active
func (n *Impl) prefundedBalancesActive(client *ethclient.Client) error {
	for _, addr := range n.preFundedMinerAddrs {
		balance, err := client.BalanceAt(context.Background(), gethcommon.HexToAddress(addr), nil)
		if err != nil {
			return fmt.Errorf("unable to check balance for account %s - %w", addr, err)
		}
		if balance.Cmp(gethcommon.Big0) == 0 {
			return fmt.Errorf("unexpected %s balance for account %s", balance.String(), addr)
		}
		fmt.Printf("Account %s prefunded with %s\n", addr, balance.String())
	}

	return nil
}

// GethGenesis returns the Genesis used in geth to boot up the network
func (n *Impl) GethGenesis() []byte {
	return n.gethGenesisBytes
}
