package l1grantsequencers

import "github.com/ten-protocol/go-ten/go/config"

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	L1HTTPURL              string
	PrivateKey             string
	EnclaveRegistryAddress string
	DockerImage            string
	SequencerURL           string

	// debugEnabled        bool
}

func NewGrantSequencerConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		L1HTTPURL:              tenCfg.Deployment.L1Deploy.RPCAddress,
		PrivateKey:             tenCfg.Deployment.L1Deploy.DeployerPK,
		EnclaveRegistryAddress: tenCfg.Network.L1.L1Contracts.EnclaveRegistryContract.Hex(),
		DockerImage:            tenCfg.Deployment.DockerImage,
		SequencerURL:           tenCfg.Deployment.L2Deploy.SequencerURL,
	}
}
