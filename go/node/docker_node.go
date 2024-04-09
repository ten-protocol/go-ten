package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/sanity-io/litter"

	"github.com/ten-protocol/go-ten/go/common/docker"
)

var (
	_hostDataDir    = "/data"        // this is how the directory is referenced within the host container
	_enclaveDataDir = "/enclavedata" // this is how the directory is references within the enclave container
)

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
	fmt.Printf("Starting Node %s with config: \n%s\n\n", d.cfg.NodeName, litter.Sdump(*d.cfg))

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
	err := docker.StopAndRemove(d.cfg.NodeName + "-host")
	if err != nil {
		return err
	}

	err = docker.StopAndRemove(d.cfg.NodeName + "-enclave")
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerNode) Upgrade(networkCfg *NetworkConfig) error {
	// TODO this should probably be removed in the future
	fmt.Printf("Upgrading node %s with config: %+v\n", d.cfg.NodeName, d.cfg)

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
		"/home/ten/go-ten/go/host/main/main",
		"-l1WSURL", d.cfg.L1WebsocketURL,
		"-enclaveRPCAddress", fmt.Sprintf("%s:%d", d.cfg.NodeName+"-enclave", d.cfg.EnclaveWSPort),
		"-managementContractAddress", d.cfg.ManagementContractAddr,
		"-messageBusContractAddress", d.cfg.MessageBusContractAddr,
		"-l1Start", d.cfg.L1Start,
		"-sequencerID", d.cfg.SequencerID,
		"-privateKey", d.cfg.PrivateKey,
		"-clientRPCHost", "0.0.0.0",
		"-logPath", "sys_out",
		"-logLevel", fmt.Sprintf("%d", log.LvlInfo),
		fmt.Sprintf("-isGenesis %t", d.cfg.IsGenesis), // boolean are a special case where the = is required
		"-nodeType", d.cfg.NodeType,
		"-profilerEnabled false",
		"-p2pPublicAddress", d.cfg.HostP2PPublicAddr,
		"-p2pBindAddress", fmt.Sprintf("0.0.0.0:%d", d.cfg.HostP2PPort),
		"-clientRPCPortHttp", fmt.Sprintf("%d", d.cfg.HostHTTPPort),
		"-clientRPCPortWs", fmt.Sprintf("%d", d.cfg.HostWSPort),
		"-maxRollupSize 65536",
		// host persistence hardcoded to use /data dir within the container, this needs to be mounted
		fmt.Sprintf("-useInMemoryDB %t", d.cfg.HostInMemDB),
		fmt.Sprintf("-debugNamespaceEnabled %t", d.cfg.IsDebugNamespaceEnabled),
		// todo (@stefan): once the limiter is in, increase it back to 5 or 10s
		fmt.Sprintf("-batchInterval %s", d.cfg.BatchInterval),
		fmt.Sprintf("-maxBatchInterval %s", d.cfg.MaxBatchInterval),
		fmt.Sprintf("-rollupInterval %s", d.cfg.RollupInterval),
		fmt.Sprintf("-logLevel %d", d.cfg.LogLevel),
		fmt.Sprintf("-isInboundP2PDisabled %t", d.cfg.IsInboundP2PDisabled),
		fmt.Sprintf("-l1ChainID %d", d.cfg.L1ChainID),
	}
	if !d.cfg.HostInMemDB {
		cmd = append(cmd, "-levelDBPath", _hostDataDir)
	}

	exposedPorts := []int{
		d.cfg.HostHTTPPort,
		d.cfg.HostWSPort,
		d.cfg.HostP2PPort,
	}

	hostVolume := map[string]string{d.cfg.NodeName + "-host-volume": _hostDataDir}

	_, err := docker.StartNewContainer(d.cfg.NodeName+"-host", d.cfg.HostDockerImage, cmd, exposedPorts, nil, nil, hostVolume)

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
		"ego", "run", "/home/ten/go-ten/go/enclave/main/main",
	}

	if d.cfg.EnclaveDebug {
		cmd = []string{
			"dlv",
			"--listen :2345",
			"--headless true",
			"--log true",
			"--api-version 2",
			"debug",
			"/home/ten/go-ten/go/enclave/main",
			"--",
		}
		exposedPorts = append(exposedPorts, 2345)
	}

	cmd = append(cmd,
		"-hostID", d.cfg.HostID,
		"-address", fmt.Sprintf("0.0.0.0:%d", d.cfg.EnclaveWSPort), // todo (@pedro) - review this 0.0.0.0 host bind
		"-nodeType", d.cfg.NodeType,
		"-managementContractAddress", d.cfg.ManagementContractAddr,
		"-hostAddress", d.cfg.HostP2PPublicAddr,
		"-sequencerID", d.cfg.SequencerID,
		"-messageBusAddress", d.cfg.MessageBusContractAddr,
		"-profilerEnabled false",
		"-useInMemoryDB false",
		"-logPath", "sys_out",
		fmt.Sprintf("-debugNamespaceEnabled %t", d.cfg.IsDebugNamespaceEnabled),
		"-maxBatchSize 36864",
		"-maxRollupSize 65536",
		fmt.Sprintf("-logLevel %d", d.cfg.LogLevel),
		"-obscuroGenesis", "{}",
	)

	if d.cfg.IsSGXEnabled {
		devices["/dev/sgx_enclave"] = "/dev/sgx_enclave"
		devices["/dev/sgx_provision"] = "/dev/sgx_provision"

		envs["OE_SIMULATION"] = "0"

		// prepend the entry.sh execution
		cmd = append([]string{"/home/ten/go-ten/go/enclave/main/entry.sh"}, cmd...)
		cmd = append(cmd,
			"-edgelessDBHost", d.cfg.NodeName+"-edgelessdb",
			"-willAttest=true",
		)
	} else {
		cmd = append(cmd,
			"-sqliteDBPath", "/data/sqlite.db",
		)
	}

	enclaveVolume := map[string]string{d.cfg.NodeName + "-enclave-volume": _enclaveDataDir}

	_, err := docker.StartNewContainer(d.cfg.NodeName+"-enclave", d.cfg.EnclaveDockerImage, cmd, exposedPorts, envs, devices, enclaveVolume)
	return err
}

func (d *DockerNode) startEdgelessDB() error {
	if !d.cfg.IsSGXEnabled {
		// Non-SGX hardware use sqlite database so EdgelessDB is not required.
		return nil
	}

	envs := map[string]string{
		"EDG_EDB_CERT_DNS": d.cfg.NodeName + "-edgelessdb",
	}

	devices := map[string]string{
		"/dev/sgx_enclave":   "/dev/sgx_enclave",
		"/dev/sgx_provision": "/dev/sgx_provision",
	}

	// only set the pccsAddr env var if it's defined
	if d.cfg.PccsAddr != "" {
		envs["PCCS_ADDR"] = d.cfg.PccsAddr
	}

	_, err := docker.StartNewContainer(d.cfg.NodeName+"-edgelessdb", d.cfg.EdgelessDBImage, nil, nil, envs, devices, nil)

	return err
}
