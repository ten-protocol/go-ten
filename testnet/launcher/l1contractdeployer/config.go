package l1contractdeployer

import "github.com/ten-protocol/go-ten/go/config"

// Config holds the properties that configure the package
type Config struct {
	L1HTTPURL            string
	PrivateKey           string
	DockerImage          string
	SequencerHostAddress string
	DebugEnabled         bool
	EtherscanAPIKey      string
	MaxGasGwei           string
	CheckGasPrice        string
}

func NewContractDeployerConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		L1HTTPURL:            tenCfg.Deployment.L1Deploy.RPCAddress,
		PrivateKey:           tenCfg.Deployment.L1Deploy.DeployerPK,
		DockerImage:          tenCfg.Deployment.DockerImage,
		SequencerHostAddress: tenCfg.Deployment.L1Deploy.InitialSeqAddress,
		DebugEnabled:         tenCfg.Deployment.DebugEnabled,
		EtherscanAPIKey:      tenCfg.Deployment.L1Deploy.EtherscanAPIKey,
		MaxGasGwei:           tenCfg.Deployment.L1Deploy.MaxGasGwei,
		CheckGasPrice:        tenCfg.Deployment.L1Deploy.CheckGasPrice,
	}
}
