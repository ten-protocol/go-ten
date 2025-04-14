package l1grantsequencers

import "github.com/ten-protocol/go-ten/go/config"

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	l1HTTPURL         string
	privateKey        string
	daRegistryAddress string
	dockerImage       string
	challengePeriod   int
}

func NewChallengePeriodConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		l1HTTPURL:         tenCfg.Deployment.L1Deploy.RPCAddress,
		privateKey:        tenCfg.Deployment.L1Deploy.DeployerPK,
		daRegistryAddress: tenCfg.Network.L1.L1Contracts.DataAvailabilityRegistry.Hex(),
		dockerImage:       tenCfg.Deployment.DockerImage,
		challengePeriod:   tenCfg.Deployment.L1Deploy.ChallengePeriod,
	}
}
