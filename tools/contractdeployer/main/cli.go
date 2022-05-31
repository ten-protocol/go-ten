package main

import (
	"flag"
	"fmt"
	"math/big"
	"time"
)

// ContractType indicates the type of contract to deploy.
type ContractType uint8

const (
	defaultL1ConnectionTimeoutSecs = 15

	management ContractType = iota
	erc20
	managementName = "management"
	erc20Name      = "erc20"

	// Flag names, defaults and usages.
	l1NodeHostName  = "l1NodeHost"
	l1NodeHostUsage = "The host on which to connect to the Ethereum client"

	l1NodePortName  = "l1NodePort"
	l1NodePortUsage = "The port on which to connect to the Ethereum client"

	l1ConnectionTimeoutSecsName  = "l1ConnectionTimeoutSecs"
	l1ConnectionTimeoutSecsUsage = "The timeout for connecting to the Ethereum client"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the L1 chain"
)

type contractDeployerConfig struct {
	l1NodeHost          string
	l1NodeWebsocketPort uint
	l1ConnectionTimeout time.Duration
	privateKeyString    string
	chainID             big.Int
	contractType        ContractType
}

func defaultContractDeployerConfig() contractDeployerConfig {
	return contractDeployerConfig{
		l1NodeHost:          "127.0.0.1",
		l1NodeWebsocketPort: 8546,
		l1ConnectionTimeout: time.Duration(defaultL1ConnectionTimeoutSecs) * time.Second,
		privateKeyString:    "0000000000000000000000000000000000000000000000000000000000000001",
		chainID:             *big.NewInt(1337),
	}
}

func parseCLIArgs() contractDeployerConfig {
	defaultConfig := defaultContractDeployerConfig()

	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.l1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(defaultConfig.l1NodeWebsocketPort), l1NodePortUsage)
	l1ConnectionTimeoutSecs := flag.Uint64(l1ConnectionTimeoutSecsName, uint64(defaultConfig.l1ConnectionTimeout.Seconds()), l1ConnectionTimeoutSecsUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.chainID.Int64(), chainIDUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.privateKeyString, privateKeyUsage)

	flag.Parse()

	defaultConfig.l1NodeHost = *l1NodeHost
	defaultConfig.l1NodeWebsocketPort = uint(*l1NodePort)
	defaultConfig.l1ConnectionTimeout = time.Duration(*l1ConnectionTimeoutSecs) * time.Second
	defaultConfig.privateKeyString = *privateKeyStr
	defaultConfig.chainID = *big.NewInt(*chainID)

	contractType := flag.Arg(0)
	switch contractType {
	case managementName:
		defaultConfig.contractType = management
	case erc20Name:
		defaultConfig.contractType = erc20
	default:
		panic(fmt.Sprintf("unrecognised contract type %s. Expected either %s or %s", contractType, managementName, erc20Name))
	}

	return defaultConfig
}
