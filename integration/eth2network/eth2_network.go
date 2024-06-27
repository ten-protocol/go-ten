package eth2network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ten-protocol/go-ten/go/common/retry"

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

// https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.14.2-35b2d07f.tar.gz
var _gethFileNameVersion = fmt.Sprintf("geth-%s-%s-%s", runtime.GOOS, runtime.GOARCH, _gethVersion)

// https://github.com/prysmaticlabs/prysm/releases/download/v4.0.6/
var (
	_prysmBeaconChainFileNameVersion = fmt.Sprintf("beacon-chain-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
	_prysmCTLFileNameVersion         = fmt.Sprintf("prysmctl-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
	_prysmValidatorFileNameVersion   = fmt.Sprintf("validator-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
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
	addr := append(preFundedAddrs, preFundedMinerAddrs...)
	genesisStr, err := generateGenesis(blockTimeSecs, chainID, preFundedMinerAddrs, addr)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(gethGenesisPath, []byte(genesisStr), 0o600)
	if err != nil {
		panic(err)
	}

	// Write beacon config
	//beaconConf := fmt.Sprintf(_beaconConfig, chainID, chainID, secondsPerSlot, slotsPerEpoch)
	//beaconConf := fmt.Sprintf(_beaconConfig, chainID, chainID)
	err = os.WriteFile(prysmConfigPath, []byte(_beaconConfig), 0o600)
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
				n.gethAuthRPCPorts[nodeID],
				n.gethNetworkPorts[nodeID],
				n.gethHTTPPorts[nodeID],
				n.gethWSPorts[nodeID],
				dataDir,
				n.preFundedMinerPKs[nodeID])
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second)
		}()
	}

	// wait for each of the nodes to start
	for i := range n.dataDirs {
		nodeID := i
		eg.Go(func() error {
			return n.waitForNodeUp(nodeID, time.Minute)
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	// link nodes together by providing the enodes (eth node address) of the other nodes
	enodes := make([]string, len(n.dataDirs))
	for i := 0; i < len(n.dataDirs); i++ {
		enodes[i], err = n.gethGetEnode(i)
		if err != nil {
			return err
		}
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
			n.prysmBeaconProcesses[nodeID], err = n.prysmStartBeaconNode(
				n.gethAuthRPCPorts[nodeID],
				n.prysmBeaconHTTPPorts[nodeID],
				n.prysmBeaconP2PPorts[nodeID],
				dataDir,
			)
			if err != nil {
				panic(err)
			}
		}()
	}

	time.Sleep(5 * time.Second)
	// start each of the validator nodes
	for i, nodeDataDir := range n.dataDirs {
		nodeID := i
		dataDir := nodeDataDir
		go func() {
			n.prysmValidatorProcesses[nodeID], err = n.prysmStartValidator(n.prysmBeaconHTTPPorts[nodeID], dataDir)
			if err != nil {
				panic(err)
			}
		}()
	}

	// blocking wait until the network reaches the Merge
	return n.waitForMergeEvent(startTime)
}

// Stop stops the network
func (n *Impl) Stop() error {
	for i := 0; i < len(n.dataDirs); i++ {
		kill(n.gethProcesses[i].Process)
		kill(n.prysmBeaconProcesses[i].Process)
		kill(n.prysmValidatorProcesses[i].Process)
	}
	// wait a second for the kill signal
	time.Sleep(time.Second)
	return nil
}

func kill(p *os.Process) {
	killErr := p.Kill()
	if killErr != nil {
		fmt.Printf("Error killing process %s", killErr)
	}
	time.Sleep(200 * time.Millisecond)
	err := p.Release()
	if err != nil {
		fmt.Printf("Error releasing process %s", err)
	}
}

// GethGenesis returns the Genesis used in geth to boot up the network
func (n *Impl) GethGenesis() []byte {
	return n.gethGenesisBytes
}

func (n *Impl) gethInitGenesisData(dataDirPath string) error {
	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{_dataDirFlag, dataDirPath, "init", n.gethGenesisPath}
	fmt.Printf("gethInitGenesisData: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.gethLogFile
	cmd.Stderr = n.gethLogFile

	return cmd.Run()
}

func (n *Impl) gethImportMinerAccount(nodeID int) error {
	script := fmt.Sprintf(importAndUnlockScript, n.preFundedMinerPKs[nodeID])

	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		"--exec", script,
		"attach", fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[nodeID]),
	}

	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.gethLogFile
	cmd.Stderr = n.gethLogFile

	return cmd.Run()
}

func (n *Impl) gethStartNode(executionPort, networkPort, httpPort, wsPort int, dataDirPath string, minerAddress string) (*exec.Cmd, error) {
	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		_dataDirFlag, dataDirPath,
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", fmt.Sprintf("%d", httpPort),
		"--http.api", "admin,miner,engine,personal,eth,net,web3,debug",
		"--mine",
		"--miner.etherbase", minerAddress,
		"--ws",
		"--ws.addr", "0.0.0.0",
		"--ws.origins", "*",
		"--ws.port", fmt.Sprintf("%d", wsPort),
		"--ws.api", "admin,miner,engine,personal,eth,net,web3,debug",
		"--authrpc.addr", "0.0.0.0",
		"--authrpc.port", fmt.Sprintf("%d", executionPort),
		"--authrpc.jwtsecret", path.Join(dataDirPath, "geth", "jwtsecret"),
		"--port", fmt.Sprintf("%d", networkPort),
		"--networkid", fmt.Sprintf("%d", n.chainID),
		"--syncmode", "full", // sync mode to download and test all blocks and txs
		"--allow-insecure-unlock", // allows to use personal accounts over http/ws
		"--nodiscover",            // don't try and discover peers
		"--verbosity", "5",        // error log level
	}
	fmt.Printf("gethStartNode: %s %s\n", n.gethBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.gethLogFile
	cmd.Stderr = n.gethLogFile

	return cmd, cmd.Start()
}

func (n *Impl) prysmGenerateGenesis() error {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"testnet",
		"generate-genesis",
		"--fork", "deneb",
		"--num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		"--genesis-time-delay", "10",
		"--chain-config-file", n.prysmConfigPath,
		"--geth-genesis-json-in", n.gethGenesisPath,
		"--geth-genesis-json-out", n.gethGenesisPath,
		"--output-ssz", n.prysmGenesisPath,
	}
	fmt.Printf("prysmGenerateGenesis: %s %s\n", n.prysmBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBinaryPath, args...) //nolint
	cmd.Stdout = n.prysmBeaconLogFile
	cmd.Stderr = n.prysmBeaconLogFile

	return cmd.Run()
}

func (n *Impl) prysmStartBeaconNode(gethAuthRPCPort, rpcPort, p2pPort int, nodeDataDir string) (*exec.Cmd, error) {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "beacondata"),
		"--min-sync-peers", "0",
		"--genesis-state", n.prysmGenesisPath,
		"--bootstrap-node", "",
		"--interop-eth1data-votes",
		"--chain-config-file", n.prysmConfigPath,
		"--contract-deployment-block", "0",
		"--chain-id", fmt.Sprintf("%d", n.chainID),
		"--accept-terms-of-use",
		"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
		//"--suggested-fee-recipient", n.preFundedMinerPKs["n0"]
		"--minimum-peers-per-subnet", "0",
		"--enable-debug-rpc-endpoints",
		"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", gethAuthRPCPort),

		//"--rpc-port", fmt.Sprintf("%d", rpcPort),
		//"--p2p-udp-port", fmt.Sprintf("%d", p2pPort),
		//"--min-sync-peers", fmt.Sprintf("%d", len(n.dataDirs)-1),
		//"--minimum-peers-per-subnet", fmt.Sprintf("%d", min(len(n.dataDirs)-1, 2)),
		//"--interop-num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		//"--config-file", n.prysmConfigPath,
		//"--grpc-gateway-corsdomain", "*",
		//"--grpc-gateway-port", fmt.Sprintf("%d", rpcPort+10),
		//"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", gethAuthRPCPort),
		//"--contract-deployment-block", "0",
		//"--verbosity", "trace",
		//"--enable-debug-rpc-endpoints",
		"--force-clear-db",
		"--verbosity", "trace",
	}
	//args := []string{
	//	"--datadir", path.Join(nodeDataDir, "prysm", "beacondata"),
	//	"--interop-eth1data-votes",
	//	"--accept-terms-of-use",
	//	"--no-discovery",
	//	"--rpc-port", fmt.Sprintf("%d", rpcPort),
	//	"--p2p-udp-port", fmt.Sprintf("%d", p2pPort),
	//	"--min-sync-peers", "0",
	//	"--minimum-peers-per-subnet", "0",
	//	//"--min-sync-peers", fmt.Sprintf("%d", len(n.dataDirs)-1),
	//	//"--minimum-peers-per-subnet", fmt.Sprintf("%d", min(len(n.dataDirs)-1, 2)),
	//	"--interop-num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
	//	"--genesis-state", n.prysmGenesisPath,
	//	"--chain-config-file", n.prysmConfigPath,
	//	"--config-file", n.prysmConfigPath,
	//	"--contract-deployment-block", "0",
	//	"--chain-id", fmt.Sprintf("%d", n.chainID),
	//	"--grpc-gateway-corsdomain", "*",
	//	"--grpc-gateway-port", fmt.Sprintf("%d", rpcPort+10),
	//	"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", gethAuthRPCPort),
	//	"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
	//	"--contract-deployment-block", "0",
	//	"--verbosity", "trace",
	//	"--enable-debug-rpc-endpoints",
	//	"--force-clear-db",
	//}

	fmt.Printf("prysmStartBeaconNode: %s %s\n", n.prysmBeaconBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmBeaconBinaryPath, args...) //nolint
	cmd.Stdout = n.prysmBeaconLogFile
	cmd.Stderr = n.prysmBeaconLogFile

	return cmd, cmd.Start()
}

func (n *Impl) prysmStartValidator(beaconHTTPPort int, nodeDataDir string) (*exec.Cmd, error) {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "validator"),
		"--accept-terms-of-use",
		"--interop-num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		"--chain-config-file", n.prysmConfigPath,
		"--verbosity", "trace",
		//"--force-clear-db",

		//"--beacon-rpc-provider", fmt.Sprintf("127.0.0.1:%d", beaconHTTPPort),
		//"--interop-start-index", "0",
		//"--config-file", n.prysmConfigPath,
		//"--suggested-fee-recipient", "0x52FfeB84540173B15eEC5a486FdB5c769F50400a", // random address to avoid a continuous warning
		//"--disable-account-metrics",
	}

	fmt.Printf("prysmStartValidator: %s %s\n", n.prysmValidatorBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmValidatorBinaryPath, args...) //nolint
	cmd.Stdout = n.prysmValidtorLogFile
	cmd.Stderr = n.prysmValidtorLogFile

	return cmd, cmd.Start()
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

	// wait for the merge block
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

// waitForNodeUp retries continuously for the node to respond to a http request
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

func (n *Impl) gethGetEnode(i int) (string, error) {
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

func (n *Impl) gethImportEnodes(enodes []string) error {
	for i, nodePort := range n.gethHTTPPorts {
		for j, enode := range enodes {
			if i == j {
				continue // same node, node 0 does not need to know node 0 enode
			}

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				fmt.Sprintf("http://127.0.0.1:%d", nodePort),
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

			err = response.Body.Close()
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second)
	}
	return nil
}

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

func (n *Impl) ensureNoDuplicatedNetwork() error {
	for nodeIdx, port := range n.gethWSPorts {
		_, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", port))
		if err == nil {
			return fmt.Errorf("unexpected geth node %d is active before the network is started", nodeIdx)
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
