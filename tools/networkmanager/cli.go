package networkmanager

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/obscuronet/go-obscuro/integration"

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

	// Flag names and usages.
	l1NodeHostName  = "l1NodeHost"
	l1NodeHostUsage = "The host on which to connect to the Ethereum client"

	l1NodeWebsocketPortName  = "l1NodeWebsocketPort"
	l1NodeWebsocketPortUsage = "The websocket port on which to connect to the Ethereum client"

	l1ConnectionTimeoutSecsName  = "l1ConnectionTimeoutSecs"
	l1ConnectionTimeoutSecsUsage = "The timeout for connecting to the Ethereum client"

	privateKeysName  = "privateKeys"
	privateKeysUsage = "The private keys for the L1 wallets, as a comma-separated list. These wallets should have been preallocated funds"

	ethereumChainIDName  = "ethereumChainID"
	ethereumChainIDUsage = "The ID of the L1 chain"

	obscuroChainIDName  = "obscuroChainID"
	obscuroChainIDUsage = "The ID of the L2 chain"

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
	privateKeys          []string
	l1ChainID            int64
	obscuroChainID       int64
	mgmtContractAddress  common.Address
	erc20ContractAddress common.Address
	obscuroClientAddress string
}

func defaultNetworkManagerConfig() Config {
	return Config{
		l1NodeHost:          "127.0.0.1",
		l1NodeWebsocketPort: 9000,
		l1ConnectionTimeout: time.Duration(defaultL1ConnectionTimeoutSecs) * time.Second,
		// Default chosen to not conflict with default private key used by host.
		privateKeys:          []string{"0000000000000000000000000000000000000000000000000000000000000002"},
		l1ChainID:            integration.EthereumChainID,
		obscuroChainID:       integration.ObscuroChainID,
		mgmtContractAddress:  common.BytesToAddress([]byte("")),
		erc20ContractAddress: common.BytesToAddress([]byte("")),
		obscuroClientAddress: "127.0.0.1:13000",
	}
}

// ParseCLIArgs returns the config, and any arguments to the command.
func ParseCLIArgs() (Config, []string) {
	defaultConfig := defaultNetworkManagerConfig()

	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.l1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodeWebsocketPortName, uint64(defaultConfig.l1NodeWebsocketPort), l1NodeWebsocketPortUsage)
	l1ConnectionTimeoutSecs := flag.Uint64(l1ConnectionTimeoutSecsName, uint64(defaultConfig.l1ConnectionTimeout.Seconds()), l1ConnectionTimeoutSecsUsage)
	ethereumChainID := flag.Int64(ethereumChainIDName, defaultConfig.l1ChainID, ethereumChainIDUsage)
	obscuroChainID := flag.Int64(obscuroChainIDName, defaultConfig.obscuroChainID, obscuroChainIDUsage)
	privateKeys := flag.String(privateKeysName, strings.Join(defaultConfig.privateKeys, ","), privateKeysUsage)
	mgmtContractAddress := flag.String(mgmtContractAddressName, defaultConfig.mgmtContractAddress.Hex(), mgmtContractAddressUsage)
	erc20ContractAddress := flag.String(erc20ContractAddressName, defaultConfig.erc20ContractAddress.Hex(), erc20ContractAddressUsage)
	obscuroClientAddress := flag.String(obscuroClientAddressName, defaultConfig.obscuroClientAddress, obscuroClientAddressUsage)

	flag.Parse()

	defaultConfig.l1NodeHost = *l1NodeHost
	defaultConfig.l1NodeWebsocketPort = uint(*l1NodePort)
	defaultConfig.l1ConnectionTimeout = time.Duration(*l1ConnectionTimeoutSecs) * time.Second
	defaultConfig.privateKeys = strings.Split(*privateKeys, ",")
	defaultConfig.l1ChainID = *ethereumChainID
	defaultConfig.obscuroChainID = *obscuroChainID
	defaultConfig.mgmtContractAddress = common.HexToAddress(*mgmtContractAddress)
	defaultConfig.erc20ContractAddress = common.HexToAddress(*erc20ContractAddress)
	defaultConfig.obscuroClientAddress = *obscuroClientAddress

	command := flag.Arg(0)
	var args []string
	switch command {
	case deployMgmtContractName:
		defaultConfig.Command = DeployMgmtContract
	case deployERC20ContractName:
		defaultConfig.Command = DeployERC20Contract
	case injectTxsName:
		defaultConfig.Command = InjectTxs
		numOfTxs := flag.Arg(1)
		args = append(args, numOfTxs)
	default:
		panic(fmt.Sprintf("unrecognised command %s", command))
	}

	return defaultConfig, args
}
