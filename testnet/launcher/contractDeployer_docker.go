package launcher

type ContractDeployer struct {
	cfg *ContractDeployerConfig
}

func NewDockerContractDeployer(cfg *ContractDeployerConfig) (*ContractDeployer, error) {
	return &ContractDeployer{
		cfg: cfg,
	}, nil // todo: add validation
}

func (n *ContractDeployer) Start() error {

	return nil
}
