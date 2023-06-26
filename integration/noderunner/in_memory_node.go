package noderunner

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/node"
	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/sanity-io/litter"

	enclavecontainer "github.com/obscuronet/go-obscuro/go/enclave/container"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
)

type InMemNode struct {
	cfg     *node.Config
	enclave *enclavecontainer.EnclaveContainer
	host    *hostcontainer.HostContainer
}

func NewInMemNode(cfg *node.Config) *InMemNode {
	return &InMemNode{
		cfg: cfg,
	}
}

func (d *InMemNode) Start() error {
	// TODO this should probably be removed in the future
	fmt.Printf("Starting Node with config: \n%s\n\n", litter.Sdump(*d.cfg))

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
	fmt.Printf("Upgrading node with config: %+v\n", d.cfg)

	err := d.Stop()
	if err != nil {
		return err
	}

	// update network configs
	d.cfg.UpdateNodeConfig(
		node.WithManagementContractAddress(networkCfg.ManagementContractAddress),
		node.WithManagementContractAddress(networkCfg.MessageBusAddress),
		node.WithL1Start(networkCfg.L1StartHash),
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

func (d *InMemNode) startHost() error {
	hostConfig := d.cfg.ToHostConfig()
	// calculate the host ID from the private key here, so we have it for the logger
	addr, err := wallet.RetrieveAddress(hostConfig.PrivateKeyString)
	if err != nil {
		panic("unable to calculate the Node ID")
	}
	logger := testlog.Logger().New(log.CmpKey, log.HostCmp, log.NodeIDKey, *addr)
	d.host = hostcontainer.NewHostContainerFromConfig(hostConfig, logger)
	return d.host.Start()
}

func (d *InMemNode) startEnclave() error {
	enclaveCfg := d.cfg.ToEnclaveConfig()
	logger := testlog.Logger().New(log.CmpKey, log.EnclaveCmp, log.NodeIDKey, enclaveCfg.HostID)

	// if not nil, the node will use the testlog.Logger - NewEnclaveContainerWithLogger will create one otherwise
	d.enclave = enclavecontainer.NewEnclaveContainerWithLogger(enclaveCfg, logger)
	return d.enclave.Start()
}
