package enclaverunner

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
)

// EnclaveConfigToml is the structure that an enclave's .toml config is parsed into.
type EnclaveConfigToml struct {
	HostID                    string
	Address                   string
	ChainID                   int64
	WillAttest                bool
	ValidateL1Blocks          bool
	SpeculativeExecution      bool
	ManagementContractAddress string
	ERC20ContractAddresses    []string
	WriteToLogs               bool
	LogPath                   string
	UseInMemoryDb             bool
}

// ParseConfig returns a config.EnclaveConfig based on either the file identified by the `config` flag, or the flags
// with specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.EnclaveConfig {
	defaultConfig := config.DefaultEnclaveConfig()

	configPath := flag.String(configName, "", configUsage)
	hostID := flag.String(HostIDName, defaultConfig.HostID.Hex(), hostIDUsage)
	address := flag.String(AddressName, defaultConfig.Address, addressUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID, chainIDUsage)
	willAttest := flag.Bool(willAttestName, defaultConfig.WillAttest, willAttestUsage)
	validateL1Blocks := flag.Bool(validateL1BlocksName, defaultConfig.ValidateL1Blocks, validateL1BlocksUsage)
	speculativeExecution := flag.Bool(speculativeExecutionName, defaultConfig.SpeculativeExecution, speculativeExecutionUsage)
	managementContractAddress := flag.String(managementContractAddressName, defaultConfig.ManagementContractAddress.Hex(), managementContractAddressUsage)
	erc20ContractAddrs := flag.String(erc20contractAddrsName, "", erc20contractAddrsUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	useInMemoryDB := flag.Bool(useInMemoryDBName, defaultConfig.UseInMemoryDb, useInMemoryDBUsage)

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	parsedERC20ContractAddrs := strings.Split(*erc20ContractAddrs, ",")
	erc20contractAddresses := make([]*common.Address, len(parsedERC20ContractAddrs))
	if *erc20ContractAddrs != "" {
		for i, addr := range parsedERC20ContractAddrs {
			hexAddr := common.HexToAddress(addr)
			erc20contractAddresses[i] = &hexAddr
		}
	} else {
		// We handle the special case of an empty list.
		erc20contractAddresses = []*common.Address{}
	}

	defaultConfig.HostID = common.HexToAddress(*hostID)
	defaultConfig.Address = *address
	defaultConfig.ChainID = *chainID
	defaultConfig.WillAttest = *willAttest
	defaultConfig.ValidateL1Blocks = *validateL1Blocks
	defaultConfig.SpeculativeExecution = *speculativeExecution
	defaultConfig.ManagementContractAddress = common.HexToAddress(*managementContractAddress)
	defaultConfig.ERC20ContractAddresses = erc20contractAddresses
	defaultConfig.WriteToLogs = *writeToLogs
	defaultConfig.LogPath = *logPath
	defaultConfig.UseInMemoryDb = *useInMemoryDB

	return defaultConfig
}

// Parses the config from the .toml file at configPath.
func fileBasedConfig(configPath string) config.EnclaveConfig {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var tomlConfig EnclaveConfigToml
	err = toml.Unmarshal(bytes, &tomlConfig)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	erc20contractAddresses := make([]*common.Address, len(tomlConfig.ERC20ContractAddresses))
	for i, addr := range tomlConfig.ERC20ContractAddresses {
		hexAddr := common.HexToAddress(addr)
		erc20contractAddresses[i] = &hexAddr
	}

	return config.EnclaveConfig{
		HostID:                    common.HexToAddress(tomlConfig.HostID),
		Address:                   tomlConfig.Address,
		ChainID:                   tomlConfig.ChainID,
		WillAttest:                tomlConfig.WillAttest,
		ValidateL1Blocks:          tomlConfig.ValidateL1Blocks,
		SpeculativeExecution:      tomlConfig.SpeculativeExecution,
		ManagementContractAddress: common.HexToAddress(tomlConfig.ManagementContractAddress),
		ERC20ContractAddresses:    erc20contractAddresses,
		WriteToLogs:               tomlConfig.WriteToLogs,
		LogPath:                   tomlConfig.LogPath,
	}
}
