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
	defaultL1RPCTimeoutSecs = 15

	DeployMgmtContract Command = iota
	DeployERC20Contract
	InjectTxs
	deployMgmtContractName  = "deployMgmtContract"
	deployERC20ContractName = "deployERC20Contract"
	injectTxsName           = "injectTransactions"

	// Flag names and usages.
	l1NodeAddressName  = "l1NodeAddress"
	l1NodeAddressUsage = "The address on which to connect to the Ethereum client"

	l1RPCTimeoutSecsName  = "l1RPCTimeoutSecs"
	l1RPCTimeoutSecsUsage = "The timeout for connecting to, and communicating with, the Ethereum client"

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

	erc20TokenName  = "erc20Token" //nolint:gosec
	erc20TokenUsage = "The name of the ERC20 token. Default: TST"
)

type Config struct {
	Command              Command
	l1NodeAddress        string
	l1RPCTimeout         time.Duration
	privateKeys          []string
	l1ChainID            int64
	obscuroChainID       int64
	mgmtContractAddress  common.Address
	erc20ContractAddress common.Address
	obscuroClientAddress string
	erc20Token           string
}

func defaultNetworkManagerConfig() Config {
	return Config{
		l1NodeAddress: "ws://127.0.0.1:9000",
		l1RPCTimeout:  time.Duration(defaultL1RPCTimeoutSecs) * time.Second,
		// Default chosen to not conflict with default private key used by host.
		privateKeys:          []string{"0000000000000000000000000000000000000000000000000000000000000002"},
		l1ChainID:            integration.EthereumChainID,
		obscuroChainID:       integration.ObscuroChainID,
		mgmtContractAddress:  common.BytesToAddress([]byte("")),
		erc20ContractAddress: common.BytesToAddress([]byte("")),
		obscuroClientAddress: "http://127.0.0.1:13001",
		erc20Token:           "TST",
	}
}

// ParseCLIArgs returns the config, and any arguments to the command.
func ParseCLIArgs() (Config, []string) {
	defaultConfig := defaultNetworkManagerConfig()

	l1NodeAddress := flag.String(l1NodeAddressName, defaultConfig.l1NodeAddress, l1NodeAddressUsage)
	l1RPCTimeoutSecs := flag.Uint64(l1RPCTimeoutSecsName, uint64(defaultConfig.l1RPCTimeout.Seconds()), l1RPCTimeoutSecsUsage)
	ethereumChainID := flag.Int64(ethereumChainIDName, defaultConfig.l1ChainID, ethereumChainIDUsage)
	obscuroChainID := flag.Int64(obscuroChainIDName, defaultConfig.obscuroChainID, obscuroChainIDUsage)
	privateKeys := flag.String(privateKeysName, strings.Join(defaultConfig.privateKeys, ","), privateKeysUsage)
	mgmtContractAddress := flag.String(mgmtContractAddressName, defaultConfig.mgmtContractAddress.Hex(), mgmtContractAddressUsage)
	erc20ContractAddress := flag.String(erc20ContractAddressName, defaultConfig.erc20ContractAddress.Hex(), erc20ContractAddressUsage)
	obscuroClientAddress := flag.String(obscuroClientAddressName, defaultConfig.obscuroClientAddress, obscuroClientAddressUsage)
	erc20Token := flag.String(erc20TokenName, defaultConfig.obscuroClientAddress, erc20TokenUsage)

	flag.Parse()

	defaultConfig.l1NodeAddress = *l1NodeAddress
	defaultConfig.l1RPCTimeout = time.Duration(*l1RPCTimeoutSecs) * time.Second
	defaultConfig.privateKeys = strings.Split(*privateKeys, ",")
	defaultConfig.l1ChainID = *ethereumChainID
	defaultConfig.obscuroChainID = *obscuroChainID
	defaultConfig.mgmtContractAddress = common.HexToAddress(*mgmtContractAddress)
	defaultConfig.erc20ContractAddress = common.HexToAddress(*erc20ContractAddress)
	defaultConfig.obscuroClientAddress = *obscuroClientAddress
	defaultConfig.erc20Token = *erc20Token

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
