package enclaverunner

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/config"
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
	GenesisJSON               string
	EdgelessDBHost            string
	SqliteDBPath              string
	ProfilerEnabled           bool
}

// ParseConfig returns a config.EnclaveConfig based on either the file identified by the `config` flag, or the flags
// with specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.EnclaveConfig {
	cfg := config.DefaultEnclaveConfig()
	flagUsageMap := getFlagUsageMap()

	configPath := flag.String(configName, "", flagUsageMap[configName])
	hostID := flag.String(HostIDName, cfg.HostID.Hex(), flagUsageMap[HostIDName])
	hostAddress := flag.String(HostAddressName, cfg.HostAddress, flagUsageMap[HostAddressName])
	address := flag.String(AddressName, cfg.Address, flagUsageMap[AddressName])
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, flagUsageMap[l1ChainIDName])
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, flagUsageMap[obscuroChainIDName])
	willAttest := flag.Bool(willAttestName, cfg.WillAttest, flagUsageMap[willAttestName])
	validateL1Blocks := flag.Bool(validateL1BlocksName, cfg.ValidateL1Blocks, flagUsageMap[validateL1BlocksName])
	speculativeExecution := flag.Bool(speculativeExecutionName, cfg.SpeculativeExecution, flagUsageMap[speculativeExecutionName])
	managementContractAddress := flag.String(ManagementContractAddressName, cfg.ManagementContractAddress.Hex(), flagUsageMap[ManagementContractAddressName])
	erc20ContractAddrs := flag.String(Erc20ContractAddrsName, "", flagUsageMap[Erc20ContractAddrsName])
	loglevel := flag.String(logLevelName, cfg.LogLevel, flagUsageMap[logLevelName])
	logPath := flag.String(logPathName, cfg.LogPath, flagUsageMap[logPathName])
	useInMemoryDB := flag.Bool(useInMemoryDBName, cfg.UseInMemoryDB, flagUsageMap[useInMemoryDBName])
	edgelessDBHost := flag.String(edgelessDBHostName, cfg.EdgelessDBHost, flagUsageMap[edgelessDBHostName])
	sqliteDBPath := flag.String(sqliteDBPathName, cfg.SqliteDBPath, flagUsageMap[sqliteDBPathName])
	profilerEnabled := flag.Bool(profilerEnabledName, cfg.ProfilerEnabled, flagUsageMap[profilerEnabledName])

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
	cfg.EdgelessDBHost = *edgelessDBHost
	cfg.SqliteDBPath = *sqliteDBPath
	cfg.ProfilerEnabled = *profilerEnabled

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
		GenesisJSON:               []byte(tomlConfig.GenesisJSON),
		EdgelessDBHost:            tomlConfig.EdgelessDBHost,
		SqliteDBPath:              tomlConfig.SqliteDBPath,
		ProfilerEnabled:           tomlConfig.ProfilerEnabled,
	}
}
