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
	HostAddress               string
	Address                   string
	L1ChainID                 int64
	ObscuroChainID            int64
	WillAttest                bool
	ValidateL1Blocks          bool
	SpeculativeExecution      bool
	ManagementContractAddress string
	ERC20ContractAddresses    []string
	WriteToLogs               bool
	LogPath                   string
	UseInMemoryDB             bool
	ViewingKeysEnabled        bool
}

// ParseConfig returns a config.EnclaveConfig based on either the file identified by the `config` flag, or the flags
// with specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.EnclaveConfig {
	defaultConfig := config.DefaultEnclaveConfig()

	configPath := flag.String(configName, "", configUsage)
	hostID := flag.String(HostIDName, defaultConfig.HostID.Hex(), hostIDUsage)
	hostAddress := flag.String(hostAddressName, defaultConfig.HostAddress, hostAddressUsage)
	address := flag.String(AddressName, defaultConfig.Address, addressUsage)
	l1ChainID := flag.Int64(l1ChainIDName, defaultConfig.L1ChainID, l1ChainIDUsage)
	obscuroChainID := flag.Int64(obscuroChainIDName, defaultConfig.ObscuroChainID, obscuroChainIDUsage)
	willAttest := flag.Bool(willAttestName, defaultConfig.WillAttest, willAttestUsage)
	validateL1Blocks := flag.Bool(validateL1BlocksName, defaultConfig.ValidateL1Blocks, validateL1BlocksUsage)
	speculativeExecution := flag.Bool(speculativeExecutionName, defaultConfig.SpeculativeExecution, speculativeExecutionUsage)
	managementContractAddress := flag.String(ManagementContractAddressName, defaultConfig.ManagementContractAddress.Hex(), managementContractAddressUsage)
	erc20ContractAddrs := flag.String(Erc20ContractAddrsName, "", erc20ContractAddrsUsage)
	writeToLogs := flag.Bool(writeToLogsName, defaultConfig.WriteToLogs, writeToLogsUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	useInMemoryDB := flag.Bool(useInMemoryDBName, defaultConfig.UseInMemoryDB, useInMemoryDBUsage)
	viewingKeysEnabled := flag.Bool(ViewingKeysEnabledName, defaultConfig.ViewingKeysEnabled, ViewingKeysEnabledUsage)

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	fmt.Printf("addr: %s\n", *erc20ContractAddrs)
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
	defaultConfig.HostAddress = *hostAddress
	defaultConfig.Address = *address
	defaultConfig.L1ChainID = *l1ChainID
	defaultConfig.ObscuroChainID = *obscuroChainID
	defaultConfig.WillAttest = *willAttest
	defaultConfig.ValidateL1Blocks = *validateL1Blocks
	defaultConfig.SpeculativeExecution = *speculativeExecution
	defaultConfig.ManagementContractAddress = common.HexToAddress(*managementContractAddress)
	defaultConfig.ERC20ContractAddresses = erc20contractAddresses
	defaultConfig.WriteToLogs = *writeToLogs
	defaultConfig.LogPath = *logPath
	defaultConfig.UseInMemoryDB = *useInMemoryDB
	defaultConfig.ViewingKeysEnabled = *viewingKeysEnabled

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
		HostAddress:               tomlConfig.HostAddress,
		Address:                   tomlConfig.Address,
		L1ChainID:                 tomlConfig.L1ChainID,
		ObscuroChainID:            tomlConfig.ObscuroChainID,
		WillAttest:                tomlConfig.WillAttest,
		ValidateL1Blocks:          tomlConfig.ValidateL1Blocks,
		SpeculativeExecution:      tomlConfig.SpeculativeExecution,
		ManagementContractAddress: common.HexToAddress(tomlConfig.ManagementContractAddress),
		ERC20ContractAddresses:    erc20contractAddresses,
		WriteToLogs:               tomlConfig.WriteToLogs,
		LogPath:                   tomlConfig.LogPath,
		UseInMemoryDB:             tomlConfig.UseInMemoryDB,
		ViewingKeysEnabled:        tomlConfig.ViewingKeysEnabled,
	}
}
