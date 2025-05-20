package l2contractdeployer

import (
	"github.com/ten-protocol/go-ten/go/config"
)

// Config holds the properties that configure the package
type Config struct {
	L1HTTPURL              string
	L1PrivateKey           string
	L2Host                 string
	L2HTTPPort             int
	L2WSPort               int
	L2PrivateKey           string
	EnclaveRegistryAddress string
	CrossChainAddress      string
	DaRegistryAddress      string
	NetworkConfigAddress   string
	MessageBusAddress      string
	DockerImage            string
	FaucetPrefundAmount    string
	DebugEnabled           bool
}

func NewContractDeployerConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		L1HTTPURL:              tenCfg.Deployment.L1Deploy.RPCAddress,
		L1PrivateKey:           tenCfg.Deployment.L1Deploy.DeployerPK,
		L2HTTPPort:             tenCfg.Deployment.L2Deploy.HTTPPort,
		L2WSPort:               tenCfg.Deployment.L2Deploy.WSPort,
		L2Host:                 tenCfg.Deployment.L2Deploy.RPCAddress,
		L2PrivateKey:           tenCfg.Deployment.L2Deploy.DeployerPK,
		EnclaveRegistryAddress: tenCfg.Network.L1.L1Contracts.EnclaveRegistryContract.Hex(),
		CrossChainAddress:      tenCfg.Network.L1.L1Contracts.CrossChainContract.Hex(),
		DaRegistryAddress:      tenCfg.Network.L1.L1Contracts.DataAvailabilityRegistry.Hex(),
		NetworkConfigAddress:   tenCfg.Network.L1.L1Contracts.NetworkConfigContract.Hex(),
		MessageBusAddress:      tenCfg.Network.L1.L1Contracts.MessageBusContract.Hex(),
		DockerImage:            tenCfg.Deployment.DockerImage,
		FaucetPrefundAmount:    tenCfg.Deployment.L2Deploy.FaucetPrefund,
		DebugEnabled:           tenCfg.Deployment.DebugEnabled,
	}
}
