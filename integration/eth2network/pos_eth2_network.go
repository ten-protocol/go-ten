package eth2network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/integration"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	_eth2BinariesRelPath = "../.build/eth2_bin"
	_gethBinaryName      = "geth"
)

var (
	gethFileNameVersion             = fmt.Sprintf("geth-%s-%s-%s", runtime.GOOS, runtime.GOARCH, _gethVersion)
	prysmBeaconChainFileNameVersion = fmt.Sprintf("beacon-chain-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
	prysmCTLFileNameVersion         = fmt.Sprintf("prysmctl-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
	prysmValidatorFileNameVersion   = fmt.Sprintf("validator-%s-%s-%s", _prysmVersion, runtime.GOOS, runtime.GOARCH)
)

type PosImpl struct {
	buildDir                 string
	binDir                   string
	chainID                  int
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	prysmValidatorBinaryPath string
	gethNetworkPort          int
	beaconP2PPort            int
	gethWSPort               int
	gethRPCPort              int
	gethHTTPPort             int
	beaconRPCPort            int
	gethLogFile              string
	prysmBeaconLogFile       string
	prysmValidatorLogFile    string
	gethdataDir              string
	beacondataDir            string
	validatordataDir         string
	gethGenesisBytes         []byte
	gethProcessID            int
	beaconProcessID          int
	validatorProcessID       int
	wallets                  []string
	timeout                  time.Duration
}

type PosEth2Network interface {
	Start() error
	Stop() error
	GenesisBytes() []byte
}

func NewPosEth2Network(binDir string, gethNetworkPort, beaconP2PPort, gethRPCPort, gethWSPort, gethHTTPPort, beaconRPCPort, chainID int, timeout time.Duration, walletsToFund ...string) PosEth2Network {
	println("STARTING GETH NET PORT: ", gethNetworkPort)
	println("STARTING GETH RPC PORT: ", gethRPCPort)
	println("STARTING GETH WS PORT: ", gethWSPort)
	println("STARTING GETH HTTP PORT: ", gethHTTPPort)
	println("STARTING BEACON PORT: ", beaconP2PPort)
	println("STARTING BEACON RPC PORT: ", beaconRPCPort)

	build, err := getBuildNumber()
	if err != nil {
		panic(fmt.Sprintf("could not get build number: %s", err.Error()))
	}
	buildString := strconv.Itoa(build)
	buildDir := path.Join(basepath, "../.build/eth2", buildString)

	gethBinaryPath := path.Join(binDir, gethFileNameVersion, _gethBinaryName)
	prysmBeaconBinaryPath := path.Join(binDir, prysmBeaconChainFileNameVersion)
	prysmBinaryPath := path.Join(binDir, prysmCTLFileNameVersion)
	prysmValidatorBinaryPath := path.Join(binDir, prysmValidatorFileNameVersion)

	// we overwrite when we exceed the max path length for eth2 logs
	if _, err := os.Stat(buildDir); err == nil {
		fmt.Printf("Folder %s already exists, overwriting\n", buildDir)
		err := os.RemoveAll(buildDir)
		if err != nil {
			panic(fmt.Sprintf("could not remove existing folder %s: %s", buildDir, err.Error()))
		}
	}

	err = os.MkdirAll(buildDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	gethLogFile := path.Join(buildDir, "geth.log")
	prysmBeaconLogFile := path.Join(buildDir, "beacon-chain.log")
	prysmValidatorLogFile := path.Join(buildDir, "validator.log")

	gethdataDir := path.Join(buildDir, "/gethdata")
	beacondataDir := path.Join(buildDir, "/beacondata")
	validatordataDir := path.Join(buildDir, "/validatordata")

	if err = os.MkdirAll(gethdataDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err = os.MkdirAll(beacondataDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err = os.MkdirAll(validatordataDir, os.ModePerm); err != nil {
		panic(err)
	}

	genesis, err := fundWallets(walletsToFund, buildDir, chainID)
	if err != nil {
		panic(fmt.Sprintf("could not generate genesis. cause: %s", err.Error()))
	}

	return &PosImpl{
		buildDir:                 buildDir,
		binDir:                   binDir,
		chainID:                  chainID,
		gethNetworkPort:          gethNetworkPort,
		beaconP2PPort:            beaconP2PPort,
		gethWSPort:               gethWSPort,
		gethRPCPort:              gethRPCPort,
		gethHTTPPort:             gethHTTPPort,
		beaconRPCPort:            beaconRPCPort,
		gethBinaryPath:           gethBinaryPath,
		prysmBinaryPath:          prysmBinaryPath,
		prysmBeaconBinaryPath:    prysmBeaconBinaryPath,
		prysmValidatorBinaryPath: prysmValidatorBinaryPath,
		gethLogFile:              gethLogFile,
		prysmBeaconLogFile:       prysmBeaconLogFile,
		prysmValidatorLogFile:    prysmValidatorLogFile,
		gethdataDir:              gethdataDir,
		beacondataDir:            beacondataDir,
		validatordataDir:         validatordataDir,
		gethGenesisBytes:         []byte(genesis),
		wallets:                  walletsToFund,
		timeout:                  timeout,
	}
}

func (n *PosImpl) Start() error {
	startTime := time.Now()
	var eg errgroup.Group
	if err := n.checkExistingNetworks(); err != nil {
		return err
	}

	err := eg.Wait()
	go func() {
		n.gethProcessID, n.beaconProcessID, n.validatorProcessID, err = startNetworkScript(n.gethNetworkPort, n.beaconP2PPort,
			n.gethRPCPort, n.gethHTTPPort, n.gethWSPort, n.beaconRPCPort, n.chainID, n.buildDir, n.prysmBeaconLogFile,
			n.prysmValidatorLogFile, n.gethLogFile, n.prysmBeaconBinaryPath, n.prysmBinaryPath, n.prysmValidatorBinaryPath,
			n.gethBinaryPath, n.gethdataDir, n.beacondataDir, n.validatordataDir)
		time.Sleep(time.Second)
	}()

	if err != nil {
		return fmt.Errorf("could not run the script to start l1 pos network. Cause: %s", err.Error())
	}
	return n.waitForMergeEvent(startTime)
}

func (n *PosImpl) Stop() error {
	kill(n.gethProcessID)
	kill(n.beaconProcessID)
	kill(n.validatorProcessID)
	time.Sleep(time.Second)
	return nil
}

func (n *PosImpl) checkExistingNetworks() error {
	port := n.gethWSPort
	beaconP2P := n.beaconP2PPort
	_, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", port))
	if err == nil {
		return fmt.Errorf("unexpected geth node is active before the network is started")
	}

	_, err = ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", beaconP2P))
	if err == nil {
		return fmt.Errorf("unexpected beacon node is active before the network is started")
	}
	return nil
}

// waitForMergeEvent connects to the geth node and waits until block 2 (the merge block) is reached
func (n *PosImpl) waitForMergeEvent(startTime time.Time) error {
	ctx := context.Background()
	dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPort))
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	number := uint64(0)
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

	if err = n.prefundedBalanceActive(dial); err != nil {
		fmt.Printf("Error prefunding account %s\n", err.Error())
		return err
	}
	return nil
}

func (n *PosImpl) prefundedBalanceActive(client *ethclient.Client) error {
	for _, addr := range n.wallets {
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

func (n *PosImpl) GenesisBytes() []byte {
	return n.gethGenesisBytes
}

func startNetworkScript(gethNetworkPort, beaconP2PPort, gethRPCPort, gethHTTPPort, gethWSPort, beaconRPCPort, chainID int, buildDir, beaconLogFile, validatorLogFile, gethLogFile,
	beaconBinary, prysmBinary, validatorBinary, gethBinary, gethdataDir, beacondataDir, validatordataDir string,
) (int, int, int, error) {
	startScript := filepath.Join(basepath, "start-pos-network.sh")
	gethNetworkPortStr := strconv.Itoa(gethNetworkPort)
	beaconP2PPortStr := strconv.Itoa(beaconP2PPort)
	beaconRPCPortStr := strconv.Itoa(beaconRPCPort)
	gethHTTPPortStr := strconv.Itoa(gethHTTPPort)
	gethWSPortStr := strconv.Itoa(gethWSPort)
	gethRPCPortStr := strconv.Itoa(gethRPCPort)
	chainStr := strconv.Itoa(chainID)

	println("STARTING GETH NET PORT: ", gethNetworkPortStr)
	println("STARTING GETH RPC PORT: ", gethRPCPortStr)
	println("STARTING GETH WS PORT: ", gethWSPortStr)
	println("STARTING GETH HTTP PORT: ", gethHTTPPortStr)
	println("STARTING BEACON PORT: ", beaconP2PPortStr)
	println("STARTING BEACON RPC PORT: ", beaconRPCPortStr)

	cmd := exec.Command("/bin/bash", startScript,
		"--geth-network", gethNetworkPortStr,
		"--beacon-p2p", beaconP2PPortStr,
		"--geth-http", gethHTTPPortStr,
		"--geth-ws", gethWSPortStr,
		"--geth-rpc", gethRPCPortStr,
		"--beacon-rpc", beaconRPCPortStr,
		"--chainid", chainStr,
		"--build-dir", buildDir,
		"--base-path", basepath,
		"--beacon-log", beaconLogFile,
		"--validator-log", validatorLogFile,
		"--geth-log", gethLogFile,
		"--beacon-binary", beaconBinary,
		"--prysmctl-binary", prysmBinary,
		"--validator-binary", validatorBinary,
		"--geth-binary", gethBinary,
		"--gethdata-dir", gethdataDir,
		"--beacondata-dir", beacondataDir,
		"--validatordata-dir", validatordataDir,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		panic(fmt.Errorf("failed to start network script: %w\nOutput: %s", err, out.String()))
	}

	output := out.String()
	return _parsePIDs(output)
}

// we have to run all the processes from the same script for the geth.ipc to work for some reason so this is an ugly
// workaround to capture the PIDs, so we can kill each process individually
func _parsePIDs(output string) (int, int, int, error) {
	lines := strings.Split(output, "\n")
	var gethPID, beaconPID, validatorPID int
	var err error

	// reverse order since PIDS are in the last three lines
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "GETH PID") {
			fields := strings.Fields(line)
			pidStr := fields[len(fields)-1]
			gethPID, err = strconv.Atoi(pidStr)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("failed to parse GETH PID: %w", err)
			}
		} else if strings.Contains(line, "BEACON PID") {
			fields := strings.Fields(line)
			pidStr := fields[len(fields)-1]
			beaconPID, err = strconv.Atoi(pidStr)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("failed to parse BEACON PID: %w", err)
			}
		} else if strings.Contains(line, "VALIDATOR PID") {
			fields := strings.Fields(line)
			pidStr := fields[len(fields)-1]
			validatorPID, err = strconv.Atoi(pidStr)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("failed to parse VALIDATOR PID: %w", err)
			}
		}

		// finish loop early when PIDs are found
		if gethPID != 0 && beaconPID != 0 && validatorPID != 0 {
			break
		}
	}

	if gethPID == 0 || beaconPID == 0 || validatorPID == 0 {
		return 0, 0, 0, fmt.Errorf("failed to find all required PIDs in script output")
	}

	return gethPID, beaconPID, validatorPID, nil
}

// we parse the wallet addresses and append them to the genesis json, using an intermediate file which is cleaned up
// at the end of the network script. genesis bytes are returned to be parsed to the enclave config
func fundWallets(walletsToFund []string, buildDir string, chainID int) (string, error) {
	filePath := filepath.Join(basepath, "genesis-init.json")
	genesis, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var genesisJSON map[string]interface{}
	err = json.Unmarshal(genesis, &genesisJSON)
	if err != nil {
		return "", err
	}

	walletsToFund = append(walletsToFund, integration.GethNodeAddress)
	for _, account := range walletsToFund {
		genesisJSON["alloc"].(map[string]interface{})[account] = map[string]string{"balance": "7500000000000000000000000000000"}
	}

	// set the chain ID
	genesisJSON["config"].(map[string]interface{})["chainId"] = chainID

	formattedGenesisBytes, err := json.MarshalIndent(genesisJSON, "", "  ")
	if err != nil {
		return "", err
	}

	newFile := filepath.Join(buildDir, "genesis.json")
	err = os.WriteFile(newFile, formattedGenesisBytes, 0o644) //nolint:gosec
	if err != nil {
		return "", err
	}

	return string(formattedGenesisBytes), nil
}

func kill(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Error finding process with PID %d: %v\n", pid, err)
		return
	}

	killErr := process.Kill()
	if killErr != nil {
		fmt.Printf("Error killing process with PID %d: %v\n", pid, killErr)
		return
	}

	time.Sleep(200 * time.Millisecond)
	err = process.Release()
	if err != nil {
		fmt.Printf("Error releasing process with PID %d: %v\n", pid, err)
	}
}
