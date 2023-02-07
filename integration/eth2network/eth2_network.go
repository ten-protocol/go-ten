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
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"golang.org/x/sync/errgroup"
)

const (
	_dataDirFlag                     = "--datadir"
	_eth2BinariesRelPath             = "../.build/eth2_bin"
	_gethFileNameVersion             = "geth-v" + _gethVersion
	_prysmBeaconChainFileNameVersion = "beacon-chain-" + _prysmVersion
	_prysmCTLFileNameVersion         = "prysmctl-" + _prysmVersion
	_prysmValidatorFileNameVersion   = "validator-" + _prysmVersion
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
}

type Eth2Network interface {
	GethGenesis() []byte
	Start() error
	Stop() error
}

func NewEth2Network(
	binDir string,
	gethHTTPPortStart int,
	gethWSPortStart int,
	gethAuthRPCPortStart int,
	gethNetworkPortStart int,
	prysmBeaconHTTPPortStart int,
	prysmBeaconP2PPortStart int,
	chainID int,
	numNodes int,
	blockTimeSecs int,
	preFundedAddrs []string,
) Eth2Network {
	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// set the paths
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)
	gethGenesisPath := path.Join(buildDir, "genesis.json")
	prysmGenesisPath := path.Join(buildDir, "genesis.ssz")
	prysmConfigPath := path.Join(buildDir, "prysm_chain_config.yml")
	gethLogPath := path.Join(buildDir, "geth_logs.txt")
	prysmBeaconLogPath := path.Join(buildDir, "prysm_beacon_logs.txt")
	prysmValidatorLogPath := path.Join(buildDir, "prysm_validator_logs.txt")

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
	beaconConf := fmt.Sprintf(_beaconConfig, chainID, chainID)
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
	gethLogFile, err := os.Create(gethLogPath)
	if err != nil {
		panic(err)
	}
	prysmBeaconLogFile, err := os.Create(prysmBeaconLogPath)
	if err != nil {
		panic(err)
	}
	prysmValidatorLogFile, err := os.Create(prysmValidatorLogPath)
	if err != nil {
		panic(err)
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
	}
}

// Start starts the network
func (n *Impl) Start() error {
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
				n.gethAuthRPCPorts[nodeID],
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

	// this locks the process waiting for the event to happen
	return n.waitForMergeEvent(startTime)
}

// Stop stops the network
func (n *Impl) Stop() error {
	for i := 0; i < len(n.dataDirs); i++ {
		err := n.gethProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill geth node - %s\n", err.Error())
			return err
		}
		err = n.prysmBeaconProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill prysm beacon node - %s\n", err.Error())
			return err
		}
		err = n.prysmValidatorProcesses[i].Process.Kill()
		if err != nil {
			fmt.Printf("unable to kill prysm validator node - %s\n", err.Error())
			return err
		}
	}
	return nil
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
	startScript := fmt.Sprintf(gethStartupScriptJS, n.preFundedMinerPKs[nodeID])

	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		"--exec", startScript,
		"attach", fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[nodeID]),
	}

	cmd := exec.Command(n.gethBinaryPath, args...) //nolint
	cmd.Stdout = n.gethLogFile
	cmd.Stderr = n.gethLogFile

	return cmd.Run()
}

func (n *Impl) gethStartNode(executionPort, networkPort, httpPort, wsPort int, dataDirPath string) (*exec.Cmd, error) {
	// full command list at https://geth.ethereum.org/docs/fundamentals/command-line-options
	args := []string{
		_dataDirFlag, dataDirPath,
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.vhosts", "*",
		"--http.port", fmt.Sprintf("%d", httpPort),
		"--http.corsdomain", "*",
		"--http.api", "admin,miner,engine,personal,eth,net,web3,debug",
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
		"--ipcdisable",            // avoid geth erroring bc the ipc path is too long
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

func (n *Impl) prysmStartBeaconNode(gethAuthRPCPort, rpcPort, p2pPort int, nodeDataDir string) (*exec.Cmd, error) {
	// full command list at https://docs.prylabs.network/docs/prysm-usage/parameters
	args := []string{
		"--datadir", path.Join(nodeDataDir, "prysm", "beacondata"),
		"--interop-eth1data-votes",
		"--accept-terms-of-use",
		"--rpc-port", fmt.Sprintf("%d", rpcPort),
		"--p2p-udp-port", fmt.Sprintf("%d", p2pPort),
		"--min-sync-peers", fmt.Sprintf("%d", len(n.dataDirs)-1),
		"--interop-num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		"--interop-genesis-state", n.prysmGenesisPath,
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
		"--chain-id", fmt.Sprintf("%d", n.chainID),
		"--grpc-gateway-corsdomain", "*",
		"--grpc-gateway-port", fmt.Sprintf("%d", rpcPort+10),
		"--execution-endpoint", fmt.Sprintf("http://127.0.0.1:%d", gethAuthRPCPort),
		"--jwt-secret", path.Join(nodeDataDir, "geth", "jwtsecret"),
		"--contract-deployment-block", "0",
		"--verbosity", "debug",
	}

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
		//"--beacon-rpc-gateway-provider", fmt.Sprintf("127.0.0.1:%d", prysmBeaconHTTPPort+10),
		"--beacon-rpc-provider", fmt.Sprintf("127.0.0.1:%d", beaconHTTPPort),
		"--interop-num-validators", fmt.Sprintf("%d", len(n.dataDirs)),
		"--interop-start-index", "0",
		"--chain-config-file", n.prysmConfigPath,
		"--config-file", n.prysmConfigPath,
		"--force-clear-db",
		"--disable-account-metrics",
		"--accept-terms-of-use",
	}

	fmt.Printf("prysmStartValidator: %s %s\n", n.prysmValidatorBinaryPath, strings.Join(args, " "))
	cmd := exec.Command(n.prysmValidatorBinaryPath, args...) //nolint
	cmd.Stdout = n.prysmValidtorLogFile
	cmd.Stderr = n.prysmValidtorLogFile

	return cmd, cmd.Start()
}

// waitForMergeEvent connects to the geth node and waits until block 2 (the merge block) is reached
func (n *Impl) waitForMergeEvent(startTime time.Time) error {
	dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPorts[0]))
	if err != nil {
		return err
	}
	number, err := dial.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	for ; number != 7; time.Sleep(time.Second) {
		number, err = dial.BlockNumber(context.Background())
		if err != nil {
			return err
		}
	}

	fmt.Printf("Reached the merge block after %s\n", time.Since(startTime))
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
