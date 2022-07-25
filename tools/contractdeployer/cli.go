package contractdeployer

import (
	"flag"
	"math/big"
)

// DefaultConfig stores the contract deployer default config
func DefaultConfig() *Config {
	return &Config{
		NodeHost:   "",
		NodePort:   0,
		PrivateKey: "",
		ChainID:    big.NewInt(1337),
		Contract:   "",
	}
}

// Config is the structure that a contract deployer config is parsed into.
type Config struct {
	NodeHost   string   // host for the client connection
	NodePort   uint     // port for client connection
	PrivateKey string   // private key to be used for the contract deployer address
	ChainID    *big.Int // chain ID we're deploying too
	Contract   string   // the name of the contract to deploy (e.g. ERC20 or MGMT)
}

// ParseConfig returns a Config after parsing all available flags
func ParseConfig() *Config {
	defaultConfig := DefaultConfig()

	nodeHost := flag.String(nodeHostName, defaultConfig.NodeHost, nodeHostUsage)
	nodePort := flag.Uint64(nodePortName, uint64(defaultConfig.NodePort), nodePortUsage)
	contract := flag.String(contractName, defaultConfig.Contract, contractUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKey, privateKeyUsage)
	ChainID := flag.Int64(chainIDName, defaultConfig.ChainID.Int64(), chainIDUsage)

	flag.Parse()

	defaultConfig.NodeHost = *nodeHost
	defaultConfig.NodePort = uint(*nodePort)
	defaultConfig.PrivateKey = *privateKeyStr
	defaultConfig.ChainID = big.NewInt(*ChainID)
	defaultConfig.Contract = *contract

	return defaultConfig
}
