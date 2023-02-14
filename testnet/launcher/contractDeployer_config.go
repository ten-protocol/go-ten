package launcher

type ContractDeployerOption = func(c *ContractDeployerConfig)

type ContractDeployerConfig struct {
}

func NewContractDeployerConfig(opts ...ContractDeployerOption) *ContractDeployerConfig {
	defaultConfig := &ContractDeployerConfig{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}
