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
	timeout                  time.Duration
}

type PosEth2Network interface {
	Start() error
	Stop() error
}

func NewPosEth2Network(binDir string, gethRPCPort int, gethWSPort int, gethHTTPPort int, beaconRPCPort int, timeout time.Duration) PosEth2Network {
	timestamp := strconv.FormatInt(time.Now().UnixMicro(), 10)
	basepath, _ := os.Getwd() // Ensure the base path is set correctly
	buildDir := path.Join(basepath, "../.build/eth2", timestamp)

	gethBinaryPath := path.Join(binDir, "geth-darwin-arm64-1.14.6", "geth")
	prysmBeaconBinaryPath := path.Join(binDir, "beacon-chain-v5.0.4-darwin-arm64")
	prysmCTLBinaryPath := path.Join(binDir, "prysmctl-v5.0.4-darwin-arm64")
	prysmValidatorBinaryPath := path.Join(binDir, "validator-v5.0.4-darwin-arm64")

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

	return &PosImpl{
		buildDir:                 buildDir,
		binDir:                   binDir,
		gethWSPort:               gethWSPort,
		gethRPCPort:              gethRPCPort,
		gethHTTPPort:             gethHTTPPort,
		beaconRPCPort:            beaconRPCPort,
		gethBinaryPath:           gethBinaryPath,
		prysmBinaryPath:          prysmCTLBinaryPath,
		prysmBeaconBinaryPath:    prysmBeaconBinaryPath,
		prysmValidatorBinaryPath: prysmValidatorBinaryPath,
		gethLogFile:              gethLogFile,
		prysmBeaconLogFile:       prysmBeaconLogFile,
		prysmValidatorLogFile:    prysmValidatorLogFile,
		timeout:                  timeout,
	}
}

func (n *PosImpl) Start() error {
	startTime := time.Now()
	if err := n.checkExistingNetworks(); err != nil {
		return err
	}

	err := startNetworkScript(n.gethRPCPort, n.gethWSPort, n.beaconRPCPort, n.prysmBeaconLogFile, n.prysmValidatorLogFile,
		n.gethLogFile, n.gethBinaryPath, n.prysmBeaconBinaryPath, n.prysmBinaryPath, n.prysmValidatorBinaryPath)
	if err != nil {
		return fmt.Errorf("could not run the script to start l1 pos network. cause: %s", err.Error())
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

func startNetworkScript(gethRPCPort, gethWSPort, beaconRPCPort int, gethBinary, beaconBinary, prysmctlBinary, validatorBinary, beaconLogFile, validatorLogFile, gethLogFile string) error {
	scriptPath := filepath.Join(".", "start-pos-network.sh")

	// Ensure the log file directories exist
	err := os.MkdirAll(filepath.Dir(beaconLogFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create log directory for beacon log: %w", err)
	}
	err = os.MkdirAll(filepath.Dir(validatorLogFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create log directory for validator log: %w", err)
	}
	err = os.MkdirAll(filepath.Dir(gethLogFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create log directory for geth log: %w", err)
	}
	beaconRPCPortStr := strconv.Itoa(beaconRPCPort)
	gethRPCPortStr := strconv.Itoa(gethRPCPort)
	gethWSPortStr := strconv.Itoa(gethWSPort)
	cmd := exec.Command("/bin/bash", scriptPath,
		"--geth-rpc", gethRPCPortStr,
		"--geth-ws", gethWSPortStr,
		"--beacon-rpc", beaconRPCPortStr,
		"--geth-binary", gethBinary,
		"--beacon-binary", beaconBinary,
		"--prysmctl-binary", prysmctlBinary,
		"--validator-binary", validatorBinary,
		"--beacon-log", beaconLogFile,
		"--validator-log", validatorLogFile,
		"--geth-log", gethLogFile,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
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
