package networkmanager

import (
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Command indicates the command for the tool to run.
type Command uint8

const (
	defaultL1ConnectionTimeoutSecs = 15

	DeployMgmtContract Command = iota
	DeployERC20Contract
	InjectTxs
	deployMgmtContractName  = "deployMgmtContract"
	deployERC20ContractName = "deployERC20Contract"
	injectTxsName           = "injectTransactions"

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

	mgmtContractAddressName  = "managementContractAddress"
	mgmtContractAddressUsage = "The hex address of the management contract on the L1"

	erc20ContractAddressName  = "erc20ContractAddress"
	erc20ContractAddressUsage = "The hex address of the ERC20 contract on the L1"

	obscuroClientAddressName  = "obscuroClientAddress"
	obscuroClientAddressUsage = "The address at which clients connect to the Obscuro node"
)

type Config struct {
	Command              Command
	l1NodeHost           string
	l1NodeWebsocketPort  uint
	l1ConnectionTimeout  time.Duration
	privateKeyString     string
	chainID              big.Int
	mgmtContractAddress  common.Address
	erc20ContractAddress common.Address
	obscuroClientAddress string
}

func defaultNetworkManagerConfig() Config {
	return Config{
		l1NodeHost:           "127.0.0.1",
		l1NodeWebsocketPort:  8546,
		l1ConnectionTimeout:  time.Duration(defaultL1ConnectionTimeoutSecs) * time.Second,
		privateKeyString:     "0000000000000000000000000000000000000000000000000000000000000001",
		chainID:              *big.NewInt(1337),
		mgmtContractAddress:  common.BytesToAddress([]byte("")),
		erc20ContractAddress: common.BytesToAddress([]byte("")),
		obscuroClientAddress: "",
	}
}

func ParseCLIArgs() Config {
	defaultConfig := defaultNetworkManagerConfig()

	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.l1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(defaultConfig.l1NodeWebsocketPort), l1NodePortUsage)
	l1ConnectionTimeoutSecs := flag.Uint64(l1ConnectionTimeoutSecsName, uint64(defaultConfig.l1ConnectionTimeout.Seconds()), l1ConnectionTimeoutSecsUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.chainID.Int64(), chainIDUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.privateKeyString, privateKeyUsage)
	mgmtContractAddress := flag.String(mgmtContractAddressName, defaultConfig.mgmtContractAddress.Hex(), mgmtContractAddressUsage)
	erc20ContractAddress := flag.String(erc20ContractAddressName, defaultConfig.erc20ContractAddress.Hex(), erc20ContractAddressUsage)
	obscuroClientAddress := flag.String(obscuroClientAddressName, defaultConfig.obscuroClientAddress, obscuroClientAddressUsage)

	flag.Parse()

	defaultConfig.l1NodeHost = *l1NodeHost
	defaultConfig.l1NodeWebsocketPort = uint(*l1NodePort)
	defaultConfig.l1ConnectionTimeout = time.Duration(*l1ConnectionTimeoutSecs) * time.Second
	defaultConfig.privateKeyString = *privateKeyStr
	defaultConfig.chainID = *big.NewInt(*chainID)
	defaultConfig.mgmtContractAddress = common.HexToAddress(*mgmtContractAddress)
	defaultConfig.erc20ContractAddress = common.HexToAddress(*erc20ContractAddress)
	defaultConfig.obscuroClientAddress = *obscuroClientAddress

	command := flag.Arg(0)
	switch command {
	case deployMgmtContractName:
		defaultConfig.Command = DeployMgmtContract
	case deployERC20ContractName:
		defaultConfig.Command = DeployERC20Contract
	case injectTxsName:
		defaultConfig.Command = InjectTxs
	default:
		panic(fmt.Sprintf("unrecognised command %s", command))
	}

	return defaultConfig
}
