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
	err := d.startEnclave()
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
		"-enclaveRPCAddress", fmt.Sprintf("enclave:%d", d.cfg.enclaveHTTPPort),
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-privateKey", d.cfg.privateKey,
		"-clientRPCHost", "0.0.0.0",
		"-logPath", "sys_out",
		"-logLevel", "4",
		"-isGenesis", fmt.Sprintf("%t", d.cfg.isGenesis),
		"-nodeType", d.cfg.nodeType,
		"-profilerEnabled", "false",
		"-p2pPublicAddress", fmt.Sprintf("%s:%d", d.cfg.hostP2PAddr, d.cfg.hostP2PPort),
	}

	_, err := docker.StartNewContainer("host", d.cfg.hostImage, cmd, nil, nil)

	return err
}

func (d *DockerNode) startEnclave() error {
	cmd := []string{
		"ego", "run", "/home/obscuro/go-obscuro/go/enclave/main/main",
		"-hostID", d.cfg.hostID,
		"-address", fmt.Sprintf("0.0.0.0:%d", d.cfg.enclaveHTTPPort),
		"-nodeType", d.cfg.nodeType,
		"-useInMemoryDB", "false",
		"-sqliteDBPath", "/data/sqlite.db",
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-hostAddress", fmt.Sprintf("host:%d", d.cfg.hostHTTPPort),
		"-profilerEnabled", "false",
		"-hostAddress", fmt.Sprintf("host:%d", d.cfg.hostP2PPort),
		"-logPath", "sys_out",
		"-logLevel", "2",
		"-sequencerID", d.cfg.sequencerID,
		"-messageBusAddress", d.cfg.messageBusContractAddress,
	}

	_, err := docker.StartNewContainer("enclave", d.cfg.enclaveImage, cmd, nil, nil)
	return err
}
