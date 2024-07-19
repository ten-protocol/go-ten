package eth2network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/integration"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	_eth2BinariesRelPath = "../.build/eth2_bin"
	_gethBinaryName      = "geth"
	udpPort              = "12000"
	tcpPort              = "30303"
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
	networkProcess           *exec.Cmd
	timeout                  time.Duration
}

type PosEth2Network interface {
	Start() error
	Stop() error
	GenesisBytes() []byte
}

func NewPosEth2Network(binDir string, gethRPCPort, gethWSPort, gethHTTPPort, beaconRPCPort, chainID int, timeout time.Duration, walletsToFund ...string) PosEth2Network {
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

	if _, err := os.Stat(buildDir); err == nil {
		panic(fmt.Sprintf("folder %s already exists", buildDir))
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

	genesis, err := fundWallets(walletsToFund, chainID)
	if err != nil {
		panic(fmt.Sprintf("could not generate genesis. cause: %s", err.Error()))
	}

	return &PosImpl{
		buildDir:                 buildDir,
		binDir:                   binDir,
		chainID:                  chainID,
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
		startNetworkScript(n.gethHTTPPort, n.gethWSPort, n.beaconRPCPort, n.chainID, n.buildDir, n.prysmBeaconLogFile, n.prysmValidatorLogFile,
			n.gethLogFile, n.prysmBeaconBinaryPath, n.prysmBinaryPath, n.prysmValidatorBinaryPath, n.gethBinaryPath,
			n.gethdataDir, n.beacondataDir, n.validatordataDir)
		time.Sleep(time.Second)
	}()
	if err != nil {
		return fmt.Errorf("could not run the script to start l1 pos network. Cause: %s", err.Error())
	}
	return n.waitForMergeEvent(startTime)
}

func (n *PosImpl) Stop() error {
	kill(n.networkProcess.Process)
	time.Sleep(5 * time.Second)
	return nil
}

func (n *PosImpl) checkExistingNetworks() error {
	port := n.gethWSPort
	_, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%d", port))
	if err == nil {
		return fmt.Errorf("unexpected geth node is active before the network is started")
	}
	checkBindAddresses(udpPort, tcpPort)
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
	balance, err := client.BalanceAt(context.Background(), gethcommon.HexToAddress(integration.GethNodeAddress), nil)
	if err != nil {
		return fmt.Errorf("unable to check balance for account %s - %w", integration.GethNodeAddress, err)
	}
	if balance.Cmp(gethcommon.Big0) == 0 {
		return fmt.Errorf("unexpected %s balance for account %s", balance.String(), integration.GethNodeAddress)
	}
	fmt.Printf("Account %s prefunded with %s\n", integration.GethNodeAddress, balance.String())

	return nil
}

func (n *PosImpl) GenesisBytes() []byte {
	return n.gethGenesisBytes
}

func startNetworkScript(gethHTTPPort, gethWSPort, beaconRPCPort, chainID int, buildDir, beaconLogFile, validatorLogFile, gethLogFile,
	beaconBinary, prysmBinary, validatorBinary, gethBinary, gethdataDir, beacondataDir, validatordataDir string,
) {
	startScript := filepath.Join(basepath, "start-pos-network.sh")
	beaconRPCPortStr := strconv.Itoa(beaconRPCPort)
	gethHTTPPortStr := strconv.Itoa(gethHTTPPort)
	gethWSPortStr := strconv.Itoa(gethWSPort)
	chainStr := strconv.Itoa(chainID)

	// TODO move all this to a config file
	cmd := exec.Command("/bin/bash", startScript,
		"--geth-http", gethHTTPPortStr,
		"--geth-ws", gethWSPortStr,
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

// we parse the wallet addresses and append them to the genesis json, using an intermediate file which is cleaned up
// at the end of the network script. genesis bytes are returned to be parsed to the enclave config
func fundWallets(walletsToFund []string, chainID int) (string, error) {
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

	newFile := filepath.Join(basepath, "genesis-updated.json")
	err = os.WriteFile(newFile, formattedGenesisBytes, 0o644) //nolint:gosec
	if err != nil {
		return "", err
	}

	return string(formattedGenesisBytes), nil
}

func kill(p *os.Process) {
	killErr := p.Kill()
	if killErr != nil {
		fmt.Printf("Error killing process %s", killErr)
	}
	time.Sleep(200 * time.Millisecond)
	err := p.Release()

	_ = stopProcesses()
	checkBindAddresses("12000", "30303")
	if err != nil {
		fmt.Printf("Error releasing process %s", err)
	}
}

// checkBindAddresses checks that no processes exist at the specified UDP and TCP ports
func checkBindAddresses(udpPort, tcpPort string) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+udpPort)
	if err != nil {
		fmt.Printf("Error resolving UDP address: %s\n", err)
	} else {
		conn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			fmt.Printf("UDP port %s is in use\n", udpPort)
			_ = stopProcesses()
		} else {
			fmt.Printf("UDP port %s is available\n", udpPort)
			conn.Close()
		}
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+tcpPort)
	if err != nil {
		fmt.Printf("Error resolving TCP address: %s\n", err)
	} else {
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			fmt.Printf("TCP port %s is in use\n", tcpPort)
			_ = stopProcesses()
		} else {
			fmt.Printf("TCP port %s is available\n", tcpPort)
			listener.Close()
		}
	}
}

func stopProcesses() error {
	stopScript := filepath.Join(basepath, "stop-processes.sh")
	cmd := exec.Command("/bin/bash", stopScript)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stop processes: %w\nOutput: %s", err, out.String())
	}

	fmt.Println(out.String())
	return nil
}
