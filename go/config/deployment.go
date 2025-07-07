package config

// DeploymentConfig contains the configuration for the deployment of the Ten network
// this is used by scripts in the initial deployment and maintenance of network contracts etc.
//
//	yaml: `deploy`
type DeploymentConfig struct {
	L1Deploy      *L1DeployConfig `mapstructure:"l1"`
	L2Deploy      *L2DeployConfig `mapstructure:"l2"`
	DebugEnabled  bool            `mapstructure:"debug"`
	DockerImage   string          `mapstructure:"dockerImage"` // the docker image to use for the hh deployment
	GithubPAT     string          `mapstructure:"githubPAT"`   // optional github personal access token to commit the deployment config
	OutputEnvFile string          `mapstructure:"outputEnv"`   // optional output env file to write the deployed contracts data
	NetworkName   string          `mapstructure:"networkName"` // the name of the testnet env, used in prefixes for KVs for example
}

// L1DeployConfig contains the configuration for the deployment of the L1 contracts
//
//	yaml: `deploy.l1`
type L1DeployConfig struct {
	RPCAddress        string `mapstructure:"rpcAddress"`        // an RPC address for the L1 network
	DeployerPK        string `mapstructure:"deployerPK"`        // the private key of the L1 deployer account
	InitialSeqAddress string `mapstructure:"initialSeqAddress"` // the initial sequencer EOA to expect for network initialization
	ChallengePeriod   int    `mapstructure:"challengePeriod"`   // the rollup challenge period in seconds
}

// L2DeployConfig contains the configuration for the deployment of the L2 contracts
//
//	yaml: `deploy.l2`
type L2DeployConfig struct {
	RPCAddress    string `mapstructure:"rpcAddress"`    // an RPC address for the L2 network
	HTTPPort      int    `mapstructure:"httpPort"`      // the port for the L2 network RPC (http)
	WSPort        int    `mapstructure:"wsPort"`        // the port for the L2 network RPC websocket
	DeployerPK    string `mapstructure:"deployerPK"`    // the private key of the L2 deployer account
	FaucetPrefund string `mapstructure:"faucetPrefund"` // initial amount of funds to pre-fund the faucet account
	SequencerURL  string `mapstructure:"sequencerURL"`  // RPC URL of the sequencer (used to fetch enclave IDs for permissioning)
}
