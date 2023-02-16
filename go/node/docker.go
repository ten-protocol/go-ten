package node

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/docker"
)

type DockerNode struct {
	cfg *Config
}

func NewDockerNode(cfg *Config) (*DockerNode, error) {
	return &DockerNode{
		cfg: cfg,
	}, nil // todo: add config validation
}

func (d *DockerNode) Start() error {
	err := d.startEdgelessDB()
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

func (d *DockerNode) startHost() error {
	cmd := []string{
		"/home/obscuro/go-obscuro/go/host/main/main",
		"-l1NodeHost", d.cfg.l1Host,
		"-l1NodePort", fmt.Sprintf("%d", d.cfg.l1WSPort),
		"-enclaveRPCAddress", fmt.Sprintf("enclave:%d", d.cfg.enclaveWSPort),
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-privateKey", d.cfg.privateKey,
		"-clientRPCHost", "0.0.0.0",
		"-logPath", "sys_out",
		"-logLevel", "4",
		fmt.Sprintf("-isGenesis=%t", d.cfg.isGenesis), // boolean are a special case where the = is required
		"-nodeType", d.cfg.nodeType,
		"-profilerEnabled", "false",
		"-p2pPublicAddress", fmt.Sprintf("%s:%d", d.cfg.hostP2PHost, d.cfg.hostP2PPort),
		"-p2pBindAddress", fmt.Sprintf("0.0.0.0:%d", d.cfg.hostP2PPort),
		"-clientRPCPortHttp", fmt.Sprintf("%d", d.cfg.hostHTTPPort),
		"-clientRPCPortWs", fmt.Sprintf("%d", d.cfg.hostWSPort),
	}

	exposedPorts := []int{
		d.cfg.hostHTTPPort,
		d.cfg.hostWSPort,
	}

	_, err := docker.StartNewContainer("host", d.cfg.hostImage, cmd, exposedPorts, nil, nil)

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
		"-useInMemoryDB", "false",
		"-sqliteDBPath", "/data/sqlite.db",
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-hostAddress", fmt.Sprintf("host:%d", d.cfg.hostP2PPort),
		"-sequencerID", d.cfg.sequencerID,
		"-messageBusAddress", d.cfg.messageBusContractAddress,
		"-profilerEnabled", "false",
		"-logPath", "sys_out",
		"-logLevel", "2",
	)

	if d.cfg.sgxEnabled {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"

		envs["OE_SIMULATION"] = "0"
	}

	_, err := docker.StartNewContainer("enclave", d.cfg.enclaveImage, cmd, exposedPorts, envs, devices)
	return err
}

func (d *DockerNode) startEdgelessDB() error {
	if !d.cfg.sgxEnabled {
		return nil
	}

	envs := map[string]string{
		"EDG_EDB_CERT_DNS": "edgelessdb",
	}

	devices := map[string]string{
		"/dev/sgx_enclave":   "/dev/sgx_enclave",
		"/dev/sgx_provision": "/dev/sgx_provision",
	}

	// only set the pccsAddr env var if it's defined
	if d.cfg.pccsAddr != "" {
		envs["PCCS_ADDR"] = d.cfg.pccsAddr
	}

	_, err := docker.StartNewContainer("edgelessdb", d.cfg.edgelessDBImage, nil, nil, envs, devices)

	return err
}
