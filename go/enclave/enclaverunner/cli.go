package enclaverunner

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
	"github.com/obscuronet/obscuro-playground/go/config"
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
	LogLevel                  string
	LogPath                   string
	UseInMemoryDB             bool
	ViewingKeysEnabled        bool
}

// ParseConfig returns a config.EnclaveConfig based on either the file identified by the `config` flag, or the flags
// with specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.EnclaveConfig {
	cfg := config.DefaultEnclaveConfig()

	configPath := flag.String(configName, "", configUsage)
	hostID := flag.String(HostIDName, cfg.HostID.Hex(), hostIDUsage)
	hostAddress := flag.String(HostAddressName, cfg.HostAddress, hostAddressUsage)
	address := flag.String(AddressName, cfg.Address, addressUsage)
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, l1ChainIDUsage)
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, obscuroChainIDUsage)
	willAttest := flag.Bool(willAttestName, cfg.WillAttest, willAttestUsage)
	validateL1Blocks := flag.Bool(validateL1BlocksName, cfg.ValidateL1Blocks, validateL1BlocksUsage)
	speculativeExecution := flag.Bool(speculativeExecutionName, cfg.SpeculativeExecution, speculativeExecutionUsage)
	managementContractAddress := flag.String(ManagementContractAddressName, cfg.ManagementContractAddress.Hex(), managementContractAddressUsage)
	erc20ContractAddrs := flag.String(Erc20ContractAddrsName, "", erc20ContractAddrsUsage)
	loglevel := flag.String(logLevelName, cfg.LogLevel, logLevelUsage)
	logPath := flag.String(logPathName, cfg.LogPath, logPathUsage)
	useInMemoryDB := flag.Bool(useInMemoryDBName, cfg.UseInMemoryDB, useInMemoryDBUsage)
	viewingKeysEnabled := flag.Bool(ViewingKeysEnabledName, cfg.ViewingKeysEnabled, ViewingKeysEnabledUsage)
	edgelessDBHost := flag.String(edgelessDBHostName, cfg.EdgelessDBHost, edgelessDBHostUsage)
	sqliteDBPath := flag.String(sqliteDBPathName, cfg.SqliteDBPath, sqliteDBPathUsage)

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

	cfg.HostID = common.HexToAddress(*hostID)
	cfg.HostAddress = *hostAddress
	cfg.Address = *address
	cfg.L1ChainID = *l1ChainID
	cfg.ObscuroChainID = *obscuroChainID
	cfg.WillAttest = *willAttest
	cfg.ValidateL1Blocks = *validateL1Blocks
	cfg.SpeculativeExecution = *speculativeExecution
	cfg.ManagementContractAddress = common.HexToAddress(*managementContractAddress)
	cfg.ERC20ContractAddresses = erc20contractAddresses
	cfg.LogLevel = *loglevel
	cfg.LogPath = *logPath
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.ViewingKeysEnabled = *viewingKeysEnabled
	cfg.EdgelessDBHost = *edgelessDBHost
	cfg.SqliteDBPath = *sqliteDBPath

	return cfg
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
		LogLevel:                  tomlConfig.LogLevel,
		LogPath:                   tomlConfig.LogPath,
		UseInMemoryDB:             tomlConfig.UseInMemoryDB,
		ViewingKeysEnabled:        tomlConfig.ViewingKeysEnabled,
	}
}
