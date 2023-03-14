package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/docker"
)

var (
	_hostDataDir      = "/data"
	_defaultHostMount = map[string]string{"host-persistence": _hostDataDir}
)

type DockerNode struct {
	cfg *Config
}

func NewDockerNode(cfg *Config) (Node, error) {
	return &DockerNode{
		cfg: cfg,
	}, nil // todo: add config validation
}

func (d *DockerNode) Start() error {
	// TODO this should probably be removed in the future
	fmt.Printf("Starting Node %s with config: %+v\n", d.cfg.nodeName, d.cfg)

	// write the network-level config to disk for future restarts
	err := WriteNetworkConfigToDisk(d.getNetworkConfig())
	if err != nil {
		return err
	}

	err = d.startEdgelessDB()
	if err != nil {
		return err
	}

	err = d.startEnclave()
	if err != nil {
		return err
	}

	err = d.startHost()
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerNode) Upgrade() error {
	// TODO this should probably be removed in the future
	fmt.Printf("Upgrading node %s with config: %+v\n", d.cfg.nodeName, d.cfg)

	// first we load network-specific details from the initial node setup from disk
	networkCfg, err := ReadNetworkConfigFromDisk()
	if err != nil {
		return err
	}
	d.updateConfigWithNetworkConfig(networkCfg)

	fmt.Println("Stopping existing host and enclave")
	err = docker.StopAndRemove(d.cfg.nodeName + "-host")
	if err != nil {
		return err
	}

	err = docker.StopAndRemove(d.cfg.nodeName + "-enclave")
	if err != nil {
		return err
	}

	fmt.Println("Starting upgraded host and enclave")
	err = d.startEnclave()
	if err != nil {
		return err
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
		"-l1NodeHost", d.cfg.l1Host,
		"-l1NodePort", fmt.Sprintf("%d", d.cfg.l1WSPort),
		"-enclaveRPCAddress", fmt.Sprintf("%s:%d", d.cfg.nodeName+"-enclave", d.cfg.enclaveWSPort),
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-privateKey", d.cfg.privateKey,
		"-clientRPCHost", "0.0.0.0",
		"-logPath", "sys_out",
		"-logLevel", fmt.Sprintf("%d", log.LvlInfo),
		fmt.Sprintf("-isGenesis=%t", d.cfg.isGenesis), // boolean are a special case where the = is required
		"-nodeType", d.cfg.nodeType,
		"-profilerEnabled=false",
		"-p2pPublicAddress", d.cfg.hostPublicP2PAddr,
		"-p2pBindAddress", fmt.Sprintf("0.0.0.0:%d", d.cfg.hostP2PPort),
		"-clientRPCPortHttp", fmt.Sprintf("%d", d.cfg.hostHTTPPort),
		"-clientRPCPortWs", fmt.Sprintf("%d", d.cfg.hostWSPort),
		// host persistence hardcoded to use /data dir within the container, this needs to be mounted
		"-useInMemoryDB=false",
		"-levelDBPath", _hostDataDir,
	}

	exposedPorts := []int{
		d.cfg.hostHTTPPort,
		d.cfg.hostWSPort,
		d.cfg.hostP2PPort,
	}

	_, err := docker.StartNewContainer(d.cfg.nodeName+"-host", d.cfg.hostImage, cmd, exposedPorts, nil, nil, _defaultHostMount)

	return err
}

func (d *DockerNode) startEnclave() error {
	devices := map[string]string{}
	exposedPorts := []int{}
	envs := map[string]string{
		"OE_SIMULATION": "1",
	}

	// default start of the enclave
	cmd := []string{
		"ego", "run", "/home/obscuro/go-obscuro/go/enclave/main/main",
	}

	if d.cfg.enclaveDebug {
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

	cmd = append(cmd,
		"-hostID", d.cfg.hostID,
		"-address", fmt.Sprintf("0.0.0.0:%d", d.cfg.enclaveWSPort), // todo review this 0.0.0.0 host bind
		"-nodeType", d.cfg.nodeType,
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-hostAddress", d.cfg.hostPublicP2PAddr,
		"-sequencerID", d.cfg.sequencerID,
		"-messageBusAddress", d.cfg.messageBusContractAddress,
		"-profilerEnabled=false",
		"-useInMemoryDB=false",
		"-logPath", "sys_out",
		"-logLevel", fmt.Sprintf("%d", log.LvlInfo),
	)

	if d.cfg.sgxEnabled {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"

		envs["OE_SIMULATION"] = "0"

		// prepend the entry.sh execution
		cmd = append([]string{"/home/obscuro/go-obscuro/go/enclave/main/entry.sh"}, cmd...)
		cmd = append(cmd,
			"-edgelessDBHost", d.cfg.nodeName+"-edgelessdb",
			"-willAttest=true",
		)
	} else {
		cmd = append(cmd,
			"-sqliteDBPath", "/data/sqlite.db",
		)
	}

	_, err := docker.StartNewContainer(d.cfg.nodeName+"-enclave", d.cfg.enclaveImage, cmd, exposedPorts, envs, devices, nil)
	return err
}

func (d *DockerNode) startEdgelessDB() error {
	if !d.cfg.sgxEnabled {
		// Non-SGX hardware use sqlite database so EdgelessDB is not required.
		return nil
	}

	envs := map[string]string{
		"EDG_EDB_CERT_DNS": d.cfg.nodeName + "-edgelessdb",
	}

	devices := map[string]string{
		"/dev/sgx_enclave":   "/dev/sgx_enclave",
		"/dev/sgx_provision": "/dev/sgx_provision",
	}

	// only set the pccsAddr env var if it's defined
	if d.cfg.pccsAddr != "" {
		envs["PCCS_ADDR"] = d.cfg.pccsAddr
	}

	_, err := docker.StartNewContainer(d.cfg.nodeName+"-edgelessdb", d.cfg.edgelessDBImage, nil, nil, envs, devices, nil)

	return err
}

func (d *DockerNode) getNetworkConfig() networkConfig {
	return networkConfig{
		managementContractAddress: d.cfg.managementContractAddr,
		messageBusAddress:         d.cfg.messageBusContractAddress,
	}
}

func (d *DockerNode) updateConfigWithNetworkConfig(networkCfg *networkConfig) {
	d.cfg.managementContractAddr = networkCfg.managementContractAddress
	d.cfg.messageBusContractAddress = networkCfg.messageBusAddress
}
