package node

import (
	"fmt"

	"github.com/sanity-io/litter"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/docker"
)

var _enclaveDataDir = "/enclavedata" // this is how the directory is references within the enclave container

type DockerNode struct {
	cfg *Config
}

func NewDockerNode(cfg *Config) *DockerNode {
	return &DockerNode{
		cfg: cfg,
	}
}

func (d *DockerNode) Start() error {
	// todo (@pedro) - this should probably be removed in the future
	fmt.Printf("Starting Node %s with config: \n%s\n\n", d.cfg.nodeName, litter.Sdump(*d.cfg))

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

func (d *DockerNode) Stop() error {
	fmt.Println("Stopping existing host and enclave")
	err := docker.StopAndRemove(d.cfg.nodeName + "-host")
	if err != nil {
		return err
	}

	err = docker.StopAndRemove(d.cfg.nodeName + "-enclave")
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerNode) Upgrade(networkCfg *NetworkConfig) error {
	// TODO this should probably be removed in the future
	fmt.Printf("Upgrading node %s with config: %+v\n", d.cfg.nodeName, d.cfg)

	err := d.Stop()
	if err != nil {
		return err
	}

	// update network configs
	d.cfg.UpdateNodeConfig(
		WithManagementContractAddress(networkCfg.ManagementContractAddress),
		WithMessageBusContractAddress(networkCfg.MessageBusAddress),
		WithL1Start(networkCfg.L1StartHash),
	)

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
		"-l1WSURL", d.cfg.l1WSURL,
		"-enclaveRPCAddresses", fmt.Sprintf("%s:%d", d.cfg.nodeName+"-enclave", d.cfg.enclaveWSPort),
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-messageBusContractAddress", d.cfg.messageBusContractAddress,
		"-l1Start", d.cfg.l1Start,
		"-sequencerP2PAddress", d.cfg.sequencerP2PAddr,
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
		"-maxRollupSize=131072",
		// host persistence hardcoded to use /data dir within the container, this needs to be mounted
		fmt.Sprintf("-useInMemoryDB=%t", d.cfg.hostInMemDB),
		fmt.Sprintf("-debugNamespaceEnabled=%t", d.cfg.debugNamespaceEnabled),
		// todo (@stefan): once the limiter is in, increase it back to 5 or 10s
		fmt.Sprintf("-batchInterval=%s", d.cfg.batchInterval),
		fmt.Sprintf("-maxBatchInterval=%s", d.cfg.maxBatchInterval),
		fmt.Sprintf("-rollupInterval=%s", d.cfg.rollupInterval),
		fmt.Sprintf("-logLevel=%d", d.cfg.logLevel),
		fmt.Sprintf("-isInboundP2PDisabled=%t", d.cfg.isInboundP2PDisabled),
		fmt.Sprintf("-l1ChainID=%d", d.cfg.l1ChainID),
		fmt.Sprintf("-l1BeaconUrl=%s", d.cfg.l1BeaconUrl),
	}
	if !d.cfg.hostInMemDB {
		cmd = append(cmd, "-postgresDBHost", d.cfg.postgresDB)
	}

	exposedPorts := []int{
		d.cfg.hostHTTPPort,
		d.cfg.hostWSPort,
		d.cfg.hostP2PPort,
	}

	_, err := docker.StartNewContainer(d.cfg.nodeName+"-host", d.cfg.hostImage, cmd, exposedPorts, nil, nil, nil, true)

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
		"-address", fmt.Sprintf("0.0.0.0:%d", d.cfg.enclaveWSPort), // todo (@pedro) - review this 0.0.0.0 host bind
		"-nodeType", d.cfg.nodeType,
		"-managementContractAddress", d.cfg.managementContractAddr,
		"-hostAddress", d.cfg.hostPublicP2PAddr,
		"-messageBusAddress", d.cfg.messageBusContractAddress,
		"-profilerEnabled=false",
		"-useInMemoryDB=false",
		"-logPath", "sys_out",
		"-logLevel", fmt.Sprintf("%d", log.LvlInfo),
		fmt.Sprintf("-debugNamespaceEnabled=%t", d.cfg.debugNamespaceEnabled),
		"-maxBatchSize=56320",
		"-maxRollupSize=131072",
		fmt.Sprintf("-logLevel=%d", d.cfg.logLevel),
		"-tenGenesis", "{}",
		"-edgelessDBHost", d.cfg.nodeName+"-edgelessdb",
	)

	if d.cfg.sgxEnabled {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"

		envs["OE_SIMULATION"] = "0"

		// prepend the entry.sh execution
		cmd = append([]string{"/home/obscuro/go-obscuro/go/enclave/main/entry.sh"}, cmd...)
		cmd = append(cmd, "-willAttest=true")
	} else {
		cmd = append(cmd, "-willAttest=false")
	}

	// we need the enclave volume to store the db credentials
	enclaveVolume := map[string]string{d.cfg.nodeName + "-enclave-volume": _enclaveDataDir}
	_, err := docker.StartNewContainer(d.cfg.nodeName+"-enclave", d.cfg.enclaveImage, cmd, exposedPorts, envs, devices, enclaveVolume, true)

	return err
}

func (d *DockerNode) startEdgelessDB() error {
	envs := map[string]string{
		"EDG_EDB_CERT_DNS": d.cfg.nodeName + "-edgelessdb",
	}
	devices := map[string]string{}

	if d.cfg.sgxEnabled {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"
	} else {
		envs["OE_SIMULATION"] = "1"
	}

	// only set the pccsAddr env var if it's defined
	if d.cfg.pccsAddr != "" {
		envs["PCCS_ADDR"] = d.cfg.pccsAddr
	}

	// todo - do we need this volume?
	//dbVolume := map[string]string{d.cfg.nodeName + "-db-volume": "/data"}
	//_, err := docker.StartNewContainer(d.cfg.nodeName+"-edgelessdb", d.cfg.edgelessDBImage, nil, nil, envs, devices, dbVolume)

	_, err := docker.StartNewContainer(d.cfg.nodeName+"-edgelessdb", d.cfg.edgelessDBImage, nil, nil, envs, devices, nil, true)

	return err
}
