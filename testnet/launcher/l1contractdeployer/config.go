package l1contractdeployer

import "github.com/ten-protocol/go-ten/go/config"

// Config holds the properties that configure the package
type Config struct {
	L1HTTPURL    string
	PrivateKey   string
	DockerImage  string
	DebugEnabled bool
}

func NewContractDeployerConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		L1HTTPURL:    tenCfg.Deployment.L1Deploy.RPCAddress,
		PrivateKey:   tenCfg.Deployment.L1Deploy.DeployerPK,
		DockerImage:  tenCfg.Deployment.DockerImage,
		DebugEnabled: tenCfg.Deployment.DebugEnabled,
	}
}
