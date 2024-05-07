package config

type TestnetConfig struct {
	TestNetSettings    TestNetSettings          `yaml:"testNetSettings"`
	Eth2Network        Eth2NetworkConfig        `yaml:"eth2Network"`
	L1ContractDeployer L1ContractDeployerConfig `yaml:"l1ContractDeployer"`
	Network            NetworkInputConfig       `yaml:"networkConfig"`
	Nodes              []*NodeConfig            `yaml:"nodes"`
	L2ContractDeployer L2ContractDeployerConfig `yaml:"l2ContractDeployer"`
	Faucet             faucetConfig             `yaml:"faucet"`
}

type TestNetSettings struct {
	Eth2Network        bool `yaml:"eth2Network"`
	L1ContractDeployer bool `yaml:"l1ContractDeployer"`
	Nodes              bool `yaml:"nodes"`
	L2ContractDeployer bool `yaml:"l2ContractDeployer"`
	Faucet             bool `yaml:"faucet"`
	ReplaceRunning     bool `yaml:"replaceRunning"`
	DryRun             bool `yaml:"dryRun"`
}

// Eth2NetworkConfig represents the configurations passed into the eth2network (L1) over CLI
type Eth2NetworkConfig struct {
	GethHTTPPort           int      `yaml:"gethHTTPPort"`
	GethWebsocketPort      int      `yaml:"gethWebsocketPort"`
	GethPrefundedAddresses []string `yaml:"gethPrefundedAddresses"`
	GethNumNodes           int      `yaml:"gethNumNodes"`
	GethImage              string   `yaml:"gethImage"`
	ContainerName          string   `yaml:"containerName"`
}

// L1ContractDeployerConfig represents the configurations passed into the deployer over CLI
type L1ContractDeployerConfig struct {
	L1HTTPurl             string `yaml:"l1HttpUrl"`
	PrivateKey            string `yaml:"privateKey"`
	L1DeployerImage       string `yaml:"l1DeployerImage"`
	ContractsEnvFile      string `yaml:"contractsEnvFile"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
	ContainerName         string `yaml:"containerName"`
}

// L2ContractDeployerConfig represents the configurations passed into the deployer over CLI
type L2ContractDeployerConfig struct {
	L1HTTPURL             string `yaml:"l1HttpURL"`
	L1PrivateKey          string `yaml:"l1PrivateKey"`
	L2DeployerImage       string `yaml:"l2DeployerImage"`
	L2WebsocketURL        string `yaml:"l2WebsocketURL"`
	L2PrivateKey          string `yaml:"l2PrivateKey"`
	L2HOCPrivateKey       string `yaml:"l2HOCPrivateKey"`
	L2POCPrivateKey       string `yaml:"l2POCPrivateKey"`
	FaucetFunding         string `yaml:"faucetFunding"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
	ContainerName         string `yaml:"containerName"`
}

type faucetConfig struct {
	TenNodeHost   string `yaml:"tenNodeHost"`
	TenNodePort   int    `yaml:"tenNodePort"`
	FaucetPort    int    `yaml:"faucetPort"`
	PrivateKey    string `yaml:"privateKey"`
	FaucetImage   string `yaml:"faucetImage"`
	ContainerName string `yaml:"containerName"`
}
