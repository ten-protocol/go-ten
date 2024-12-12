package node

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/common/docker"
)

var _enclaveDataDir = "/enclavedata" // this is how the directory is references within the enclave container

type DockerNode struct {
	cfg              *config.TenConfig
	hostImage        string
	enclaveImage     string
	edgelessDBImage  string
	enclaveDebugMode bool
	pccsAddr         string // optional specified PCCS address
	numEnclaves      int    // number of enclaves to start for the node as an HA setup
}

func NewDockerNode(cfg *config.TenConfig, hostImage, enclaveImage, edgelessDBImage string, enclaveDebug bool, pccsAddr string, numEnclaves int) *DockerNode {
	return &DockerNode{
		cfg:              cfg,
		hostImage:        hostImage,
		enclaveImage:     enclaveImage,
		edgelessDBImage:  edgelessDBImage,
		enclaveDebugMode: enclaveDebug,
		pccsAddr:         pccsAddr,
		numEnclaves:      numEnclaves,
	}
}

func (d *DockerNode) Start() error {
	// todo (@pedro) - this should probably be removed in the future
	d.cfg.PrettyPrint() // dump config to stdout

	var err error
	for i := 0; i < d.numEnclaves; i++ {
		err = d.startEdgelessDB(i)
		if err != nil {
			return fmt.Errorf("failed to start edgelessdb: %w", err)
		}

		err = d.startEnclave(i)
		if err != nil {
			return fmt.Errorf("failed to start enclave: %w", err)
		}
	}

	err = d.startHost()
	if err != nil {
		return fmt.Errorf("failed to start host: %w", err)
	}

	return nil
}

func (d *DockerNode) Stop() error {
	fmt.Println("Stopping existing host and enclave")
	err := docker.StopAndRemove(d.cfg.Node.Name + "-host")
	if err != nil {
		return err
	}

	for i := 0; i < d.numEnclaves; i++ {
		err = docker.StopAndRemove(d.cfg.Node.Name + "-enclave-" + strconv.Itoa(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DockerNode) Upgrade(networkCfg *NetworkConfig) error {
	// TODO this should probably be removed in the future
	fmt.Printf("Upgrading node %s with config: %+v\n", d.cfg.Node.Name, d.cfg)

	err := d.Stop()
	if err != nil {
		return err
	}

	// update the config with the existing network config
	d.cfg.Network.L1.L1Contracts.ManagementContract = common.HexToAddress(networkCfg.ManagementContractAddress)
	d.cfg.Network.L1.L1Contracts.MessageBusContract = common.HexToAddress(networkCfg.MessageBusAddress)
	d.cfg.Network.L1.StartHash = common.HexToHash(networkCfg.L1StartHash)

	fmt.Println("Starting upgraded host and enclaves")
	for i := 0; i < d.numEnclaves; i++ {
		err = d.startEnclave(i)
		if err != nil {
			return err
		}
	}

	err = d.startHost()
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerNode) startHost() error {
	cmd := []string{
		"/home/obscuro/go-obscuro/go/host/main/main",
	}

	// split on ":" to extract p2p port from bind address
	p2pPortStr := d.cfg.Host.P2P.BindAddress[strings.LastIndex(d.cfg.Host.P2P.BindAddress, ":")+1:]
	// convert to int
	p2pPort, err := strconv.Atoi(p2pPortStr)
	if err != nil {
		return fmt.Errorf("failed to convert p2p port to int: %w", err)
	}

	exposedPorts := []int{
		int(d.cfg.Host.RPC.HTTPPort),
		int(d.cfg.Host.RPC.WSPort),
		p2pPort,
	}

	envVariables := d.cfg.ToEnvironmentVariables()

	_, err = docker.StartNewContainer(d.cfg.Node.Name+"-host", d.hostImage, cmd, exposedPorts, envVariables, nil, nil, true)

	return err
}

func (d *DockerNode) startEnclave(enclaveIdx int) error {
	devices := map[string]string{}
	exposedPorts := []int{}

	// default start of the enclave
	cmd := []string{
		"ego", "run", "/home/obscuro/go-obscuro/go/enclave/main/main",
	}

	if d.enclaveDebugMode {
		if d.numEnclaves > 1 {
			return fmt.Errorf("cannot run multiple enclaves in debug mode")
		}
		cmd = []string{
			"dlv",
			"--listen=:2345",
			"--headless=true",
			"--log=true",
			"--api-version=2",
			"debug",
			"/home/obscuro/go-obscuro/go/enclave/main",
			"--",
		}
		exposedPorts = append(exposedPorts, 2345)
	}

	envVariables := d.cfg.ToEnvironmentVariables()

	if d.cfg.Enclave.EnableAttestation {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"

		envVariables["OE_SIMULATION"] = "0"

		// prepend the entry.sh execution
		cmd = append([]string{"/home/obscuro/go-obscuro/go/enclave/main/entry.sh"}, cmd...)
		cmd = append(cmd, "-willAttest=true")
	} else {
		envVariables["OE_SIMULATION"] = "1"
		cmd = append(cmd, "-willAttest=false")
	}

	volumeName := fmt.Sprintf("%s-enclave-volume-%d", d.cfg.Node.Name, enclaveIdx)
	containerName := fmt.Sprintf("%s-enclave-%d", d.cfg.Node.Name, enclaveIdx)

	// we need the enclave volume to store the db credentials
	enclaveVolume := map[string]string{volumeName: _enclaveDataDir}
	_, err := docker.StartNewContainer(containerName, d.enclaveImage, cmd, exposedPorts, envVariables, devices, enclaveVolume, true)

	return err
}

func (d *DockerNode) startEdgelessDB(enclaveIdx int) error {
	containerName := fmt.Sprintf("%s-edgelessdb-%d", d.cfg.Node.Name, enclaveIdx)
	envs := map[string]string{
		"EDG_EDB_CERT_DNS": containerName,
	}
	devices := map[string]string{}

	if d.cfg.Enclave.EnableAttestation {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"
	} else {
		envs["OE_SIMULATION"] = "1"
	}

	// only set the pccsAddr env var if it's defined
	if d.pccsAddr != "" {
		envs["PCCS_ADDR"] = d.pccsAddr
	}

	// todo - do we need this volume?
	//dbVolume := map[string]string{d.cfg.Node.Name + "-db-volume": "/data"}
	//_, err := docker.StartNewContainer(d.cfg.Node.Name+"-edgelessdb", d.cfg.edgelessDBImage, nil, nil, envs, devices, dbVolume)

	_, err := docker.StartNewContainer(containerName, d.edgelessDBImage, nil, nil, envs, devices, nil, true)

	return err
}
