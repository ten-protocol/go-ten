package config

type TestnetConfig struct {
	Eth2Network        Eth2NetworkConfig
	L1ContractDeployer L1ContractDeployerConfig
	Network            NetworkInputConfig
	Nodes              []NodeConfig
	L2ContractDeployer L2ContractDeployerConfig
	Faucet             faucetConfig
}

// Eth2NetworkConfig represents the configurations passed into the eth2network (L1) over CLI
type Eth2NetworkConfig struct {
	GethHTTPPort           int      `yaml:"gethHTTPPort"`
	GethWebsocketPort      int      `yaml:"gethWebsocketPort"`
	GethPrefundedAddresses []string `yaml:"gethPrefundedAddresses"`
	GethNumNodes           int      `yaml:"gethNumNodes"`
	GethImage              string   `yaml:"gethImage"`
}

// L1ContractDeployerConfig represents the configurations passed into the deployer over CLI
type L1ContractDeployerConfig struct {
	L1HTTPURL             string `yaml:"l1HTTPURL"`
	PrivateKey            string `yaml:"privateKey"`
	L1DeployerImage       string `yaml:"l1DeployerImage"`
	ContractsEnvFile      string `yaml:"contractsEnvFile"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
}

// L2ContractDeployerConfig represents the configurations passed into the deployer over CLI
type L2ContractDeployerConfig struct {
	L1HTTPURL             string `yaml:"l1HTTPURL"`
	L1PrivateKey          string `yaml:"l1PrivateKey"`
	L2DeployerImage       string `yaml:"l2DeployerImage"`
	L2WebsocketURL        string `yaml:"l2WebsocketURL"`
	L2PrivateKey          string `yaml:"l2PrivateKey"`
	L2HOCPrivateKey       string `yaml:"l2HOCPrivateKey"`
	L2POCPrivateKey       string `yaml:"l2POCPrivateKey"`
	FaucetFunding         string `yaml:"faucetFunding"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
}

type faucetConfig struct {
	TenNodeHost string `yaml:"tenNodeHost"`
	TenNodePort int    `yaml:"tenNodePort"`
	FaucetPort  int    `yaml:"faucetPort"`
	PrivateKey  string `yaml:"privateKey"`
	FaucetImage string `yaml:"faucetImage"`
}
