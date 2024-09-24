package noderunner

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/config2"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/node"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

type InMemNode struct {
	tenCfg  *config2.TenConfig
	enclave *enclavecontainer.EnclaveContainer
	host    *hostcontainer.HostContainer
}

func NewInMemNode(cfg *config2.TenConfig) *InMemNode {
	return &InMemNode{
		tenCfg: cfg,
	}
}

func (d *InMemNode) Start() error {
	// TODO this should probably be removed in the future
	d.tenCfg.PrettyPrint()

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

func (d *InMemNode) Stop() error {
	fmt.Println("Stopping existing host and enclave")
	if err := d.host.Stop(); err != nil {
		return err
	}

	return d.enclave.Stop()
}

func (d *InMemNode) Upgrade(networkCfg *node.NetworkConfig) error {
	// TODO this should probably be removed in the future
	d.tenCfg.PrettyPrint()

	err := d.Stop()
	if err != nil {
		return err
	}

	d.tenCfg.Network.L1.L1Contracts.ManagementContract = common.HexToAddress(networkCfg.ManagementContractAddress)
	d.tenCfg.Network.L1.L1Contracts.MessageBusContract = common.HexToAddress(networkCfg.MessageBusAddress)
	d.tenCfg.Network.L1.StartHash = common.HexToHash(networkCfg.L1StartHash)

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

func (d *InMemNode) startHost() error {
	hostConfig := hostconfig.HostConfigFromTenConfig(d.tenCfg)

	logger := testlog.Logger().New(log.CmpKey, log.HostCmp, log.NodeIDKey, d.tenCfg.Node.ID)
	d.host = hostcontainer.NewHostContainerFromConfig(hostConfig, logger)
	return d.host.Start()
}

func (d *InMemNode) startEnclave() error {
	enclaveCfg := enclaveconfig.EnclaveConfigFromTenConfig(d.tenCfg)
	logger := testlog.Logger().New(log.CmpKey, log.EnclaveCmp, log.NodeIDKey, d.tenCfg.Node.ID)
	enclaveCfg.LogPath = testlog.LogFile()

	// if not nil, the node will use the testlog.Logger - NewEnclaveContainerWithLogger will create one otherwise
	d.enclave = enclavecontainer.NewEnclaveContainerWithLogger(enclaveCfg, logger)
	return d.enclave.Start()
}
