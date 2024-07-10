package eth2network

import (
	"bytes"
	"context"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	eth2BinariesRelPath = "../.build/eth2_bin"
	gethBinaryName      = "geth"
	prefundedAddr       = "0x123463a4B065722E99115D6c222f267d9cABb524"
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
	gethBinaryPath           string
	prysmBinaryPath          string
	prysmBeaconBinaryPath    string
	gethGenesisPath          string
	prysmGenesisPath         string
	prysmConfigPath          string
	prysmValidatorBinaryPath string
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
	timeout                  time.Duration
}

type PosEth2Network interface {
	Start() error
	Stop() error
}

func NewPosEth2Network(binDir string, gethRPCPort int, gethWSPort int, gethHTTPPort int, beaconRPCPort int, timeout time.Duration) PosEth2Network {

	// Build dirs are suffixed with a timestamp so multiple executions don't collide
	timestamp := strconv.FormatInt(time.Now().UnixMicro(), 10)

	// set the paths
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)

	gethBinaryPath := path.Join(binDir, gethFileNameVersion, _gethBinaryName)
	prysmBeaconBinaryPath := path.Join(binDir, prysmBeaconChainFileNameVersion)
	prysmBinaryPath := path.Join(binDir, prysmCTLFileNameVersion)
	prysmValidatorBinaryPath := path.Join(binDir, prysmValidatorFileNameVersion)

	if _, err := os.Stat(buildDir); err == nil {
		panic(fmt.Sprintf("folder %s already exists", buildDir))
	}

	err := os.MkdirAll(buildDir, os.ModePerm)
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

	return &PosImpl{
		buildDir:                 buildDir,
		binDir:                   binDir,
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
		timeout:                  timeout,
	}
}

func (n *PosImpl) Start() error {
	startTime := time.Now()
	if err := n.checkExistingNetworks(); err != nil {
		return err
	}

	err := startNetworkScript(n.gethHTTPPort, n.gethWSPort, n.beaconRPCPort, n.buildDir, n.prysmBeaconLogFile, n.prysmValidatorLogFile,
		n.gethLogFile, n.prysmBeaconBinaryPath, n.prysmBinaryPath, n.prysmValidatorBinaryPath, n.gethBinaryPath,
		n.gethdataDir, n.beacondataDir, n.validatordataDir)
	if err != nil {
		return fmt.Errorf("could not run the script to start l1 pos network. Cause: %s", err.Error())
	}
	return n.waitForMergeEvent(startTime)
}

func (n *PosImpl) Stop() error {
	return stopProcesses()
}

func (n *PosImpl) checkExistingNetworks() error {
	//port := n.gethWSPort
	//_, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", port))
	//if err == nil {
	//	return fmt.Errorf("unexpected geth node %d is active before the network is started")
	//}
	return nil
}

// waitForMergeEvent connects to the geth node and waits until block 2 (the merge block) is reached
func (n *PosImpl) waitForMergeEvent(startTime time.Time) error {
	ctx := context.Background()
	dial, err := ethclient.Dial(fmt.Sprintf("http://127.0.0.1:%d", n.gethHTTPPort))
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

	if err = n.prefundedBalanceActive(dial); err != nil {
		fmt.Printf("Error prefunding account %s\n", err.Error())
		return err
	}
	return nil
}

func (n *PosImpl) prefundedBalanceActive(client *ethclient.Client) error {
	balance, err := client.BalanceAt(context.Background(), gethcommon.HexToAddress(prefundedAddr), nil)
	if err != nil {
		return fmt.Errorf("unable to check balance for account %s - %w", prefundedAddr, err)
	}
	if balance.Cmp(gethcommon.Big0) == 0 {
		return fmt.Errorf("unexpected %s balance for account %s", balance.String(), prefundedAddr)
	}
	fmt.Printf("Account %s prefunded with %s\n", prefundedAddr, balance.String())

	return nil
}

func initialiseGethScript(gethBinary, prysmBinary, gethdataDir string) error {
	scriptPath := filepath.Join(".", "start-initialise-geth.sh")

	cmd := exec.Command("/bin/bash", scriptPath)

	//cmd := exec.Command("/bin/bash", scriptPath,
	//	"--geth-binary", gethBinary,
	//	"--prysmctl-binary", prysmBinary,
	//	"--gethdata-dir", gethdataDir,
	//)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start script: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("script execution failed: %w\nOutput: %s", err, out.String())
	}

	fmt.Printf("Script output: %s\n", out.String())
	return nil
}

func startNetworkScript(gethHTTPPort, gethWSPort, beaconRPCPort int, buildDir, beaconLogFile, validatorLogFile, gethLogFile,
	beaconBinary, prysmBinary, validatorBinary, gethBinary, gethdataDir, beacondataDir, validatordataDir string) error {
	scriptPath := filepath.Join(".", "start-pos-network.sh")

	beaconRPCPortStr := strconv.Itoa(beaconRPCPort)
	gethHTTPPortStr := strconv.Itoa(gethHTTPPort)
	gethWSPortStr := strconv.Itoa(gethWSPort)
	cmd := exec.Command("/bin/bash", scriptPath,
		"--geth-http", gethHTTPPortStr,
		"--geth-ws", gethWSPortStr,
		"--beacon-rpc", beaconRPCPortStr,
		"--build-dir", buildDir,
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

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start script: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("script execution failed: %w\nOutput: %s", err, out.String())
	}

	fmt.Printf("Script output: %s\n", out.String())
	return nil
}

func stopProcesses() error {
	cmd := exec.Command("/bin/bash", "-c", "lsof -i :12000 -i :30303 -t")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to find processes: %w\nOutput: %s", err, out.String())
	}

	pids := strings.Fields(out.String())
	if len(pids) == 0 {
		return fmt.Errorf("no processes found on UDP:12000 or TCP:30303")
	}

	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)
		err := killCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to kill process %s: %w", pid, err)
		}
		fmt.Printf("Killed process %s\n", pid)
	}

	return nil
}
