package runner

import (
	"flag"
	"math/big"
)

// DefaultConfig stores the contract deployer default config
func DefaultConfig() *Config {
	return &Config{
		L1NodeHost:       "",
		L1NodePort:       0,
		MgmtContractPath: "",
		PrivateKey:       "",
		ChainID:          big.NewInt(1337),
	}
}

// Config is the structure that a contract deployer config is parsed into.
type Config struct {
	L1NodeHost       string
	L1NodePort       uint
	MgmtContractPath string
	PrivateKey       string
	ChainID          *big.Int
}

// ParseConfig returns a Config after parsing all available flags
func ParseConfig() *Config {
	defaultConfig := DefaultConfig()

	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.L1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(defaultConfig.L1NodePort), l1NodePortUsage)
	mgmtContractPath := flag.String(mgmtContractPathName, defaultConfig.MgmtContractPath, mgmtContractPathUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKey, privateKeyUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID.Int64(), chainIDUsage)

	flag.Parse()

	defaultConfig.L1NodeHost = *l1NodeHost
	defaultConfig.L1NodePort = uint(*l1NodePort)
	defaultConfig.MgmtContractPath = *mgmtContractPath
	defaultConfig.PrivateKey = *privateKeyStr
	defaultConfig.ChainID = big.NewInt(*chainID)

	return defaultConfig
}
