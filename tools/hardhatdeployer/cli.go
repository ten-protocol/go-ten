package contractdeployer

import (
	"flag"
	"math/big"
	"strings"
)

const (
	// We check the flags against these default Go values, to see if we should use the provided defaults instead.
	chainIDPlaceholder           = 0
	constructorParamsPlaceholder = ""
)

var (
	defaultL1ChainID = big.NewInt(1337)
	defaultL2ChainID = big.NewInt(443)
)

// DefaultConfig stores the contract client default config
func DefaultConfig() *Config {
	return &Config{
		NodeHost:          "",
		NodePort:          0,
		IsL1Deployment:    false,
		PrivateKey:        "",
		ChainID:           defaultL2ChainID,
		ContractName:      "",
		ConstructorParams: []string{},
	}
}

// Config is the structure that a contract deployer config is parsed into.
type Config struct {
	NodeHost          string   // host for the client connection
	NodePort          uint     // port for client connection
	IsL1Deployment    bool     // flag for L1/Eth contract deployment (rather than Obscuro/L2 deployment)
	PrivateKey        string   // private key to be used for the contract deployer address
	ChainID           *big.Int // chain ID we're deploying too
	ContractName      string   // the name of the contract to deploy (e.g. ERC20 or MGMT)
	ConstructorParams []string // parameters sent to the constructor
}

// ParseConfig returns a Config after parsing all available flags
func ParseConfig() *Config {
	defaultConfig := DefaultConfig()

	nodeHost := flag.String(nodeHostName, defaultConfig.NodeHost, nodeHostUsage)
	nodePort := flag.Uint64(nodePortName, uint64(defaultConfig.NodePort), nodePortUsage)
	isL1Deployment := flag.Bool(isL1DeploymentName, defaultConfig.IsL1Deployment, isL1DeploymentUsage)
	contractName := flag.String(contractNameName, defaultConfig.ContractName, contractNameUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKey, privateKeyUsage)
	// if this flag has a non-zero value it will be used instead of the default chain IDs
	overrideChainID := flag.Int64(chainIDName, chainIDPlaceholder, chainIDUsage)
	constructorParams := flag.String(constructorParamsName, constructorParamsPlaceholder, constructorParamsUsage)

	flag.Parse()

	defaultConfig.NodeHost = *nodeHost
	defaultConfig.NodePort = uint(*nodePort)
	defaultConfig.IsL1Deployment = *isL1Deployment
	defaultConfig.PrivateKey = *privateKeyStr
	defaultConfig.ContractName = *contractName

	if defaultConfig.IsL1Deployment {
		// for L1 deployment we default the chain ID to the L1 chain (it will still be overridden if arg was set by caller)
		defaultConfig.ChainID = defaultL1ChainID
	}
	if *overrideChainID != chainIDPlaceholder {
		defaultConfig.ChainID = big.NewInt(*overrideChainID)
	}
	if *constructorParams != constructorParamsPlaceholder {
		defaultConfig.ConstructorParams = strings.Split(*constructorParams, ",")
	}

	return defaultConfig
}
