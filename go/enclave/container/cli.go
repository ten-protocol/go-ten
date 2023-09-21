package container

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/config"
)

// EnclaveConfigToml is the structure that an enclave's .toml config is parsed into.
type EnclaveConfigToml struct {
	HostID                    string
	HostAddress               string
	Address                   string
	NodeType                  string
	L1ChainID                 int64
	ObscuroChainID            int64
	WillAttest                bool
	ValidateL1Blocks          bool
	ManagementContractAddress string
	LogLevel                  int
	LogPath                   string
	UseInMemoryDB             bool
	GenesisJSON               string
	EdgelessDBHost            string
	SqliteDBPath              string
	ProfilerEnabled           bool
	MinGasPrice               int64
	MessageBusAddress         string
	SequencerID               string
	ObscuroGenesis            string
	DebugNamespaceEnabled     bool
	MaxBatchSize              uint64
	MaxRollupSize             uint64
	GasPaymentAddress         string
	BaseFee                   uint64
	GasLimit                  uint64
}

// ParseConfig returns a config.EnclaveConfig based on either the file identified by the `config` flag, or the flags
// with specific defaults (if the `config` flag isn't specified).
func ParseConfig() (*config.EnclaveConfig, error) {
	cfg := config.DefaultEnclaveConfig()
	flagUsageMap := getFlagUsageMap()

	configPath := flag.String(configName, "", flagUsageMap[configName])
	hostID := flag.String(hostIDName, cfg.HostID.Hex(), flagUsageMap[hostIDName])
	hostAddress := flag.String(hostAddressName, cfg.HostAddress, flagUsageMap[hostAddressName])
	address := flag.String(addressName, cfg.Address, flagUsageMap[addressName])
	nodeTypeStr := flag.String(nodeTypeName, cfg.NodeType.String(), flagUsageMap[nodeTypeName])
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, flagUsageMap[l1ChainIDName])
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, flagUsageMap[obscuroChainIDName])
	willAttest := flag.Bool(willAttestName, cfg.WillAttest, flagUsageMap[willAttestName])
	validateL1Blocks := flag.Bool(validateL1BlocksName, cfg.ValidateL1Blocks, flagUsageMap[validateL1BlocksName])
	managementContractAddress := flag.String(ManagementContractAddressName, cfg.ManagementContractAddress.Hex(), flagUsageMap[ManagementContractAddressName])
	loglevel := flag.Int(logLevelName, cfg.LogLevel, flagUsageMap[logLevelName])
	logPath := flag.String(logPathName, cfg.LogPath, flagUsageMap[logPathName])
	useInMemoryDB := flag.Bool(useInMemoryDBName, cfg.UseInMemoryDB, flagUsageMap[useInMemoryDBName])
	edgelessDBHost := flag.String(edgelessDBHostName, cfg.EdgelessDBHost, flagUsageMap[edgelessDBHostName])
	sqliteDBPath := flag.String(sqliteDBPathName, cfg.SqliteDBPath, flagUsageMap[sqliteDBPathName])
	profilerEnabled := flag.Bool(profilerEnabledName, cfg.ProfilerEnabled, flagUsageMap[profilerEnabledName])
	minGasPrice := flag.Int64(minGasPriceName, cfg.MinGasPrice.Int64(), flagUsageMap[minGasPriceName])
	messageBusAddress := flag.String(messageBusAddressName, cfg.MessageBusAddress.Hex(), flagUsageMap[messageBusAddressName])
	sequencerID := flag.String(sequencerIDName, cfg.SequencerID.Hex(), flagUsageMap[sequencerIDName])
	obscuroGenesis := flag.String(obscuroGenesisName, cfg.ObscuroGenesis, flagUsageMap[obscuroGenesisName])
	debugNamespaceEnabled := flag.Bool(debugNamespaceEnabledName, cfg.DebugNamespaceEnabled, flagUsageMap[debugNamespaceEnabledName])
	maxBatchSize := flag.Uint64(maxBatchSizeName, cfg.MaxBatchSize, flagUsageMap[maxBatchSizeName])
	maxRollupSize := flag.Uint64(maxRollupSizeName, cfg.MaxRollupSize, flagUsageMap[maxRollupSizeName])
	baseFee := flag.Uint64("l2BaseFee", cfg.BaseFee.Uint64(), "")
	coinbaseAddress := flag.String("l2Coinbase", cfg.GasPaymentAddress.Hex(), "")
	gasLimit := flag.Uint64("l2GasLimit", cfg.GasLimit.Uint64(), "")

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	nodeType, err := common.ToNodeType(*nodeTypeStr)
	if err != nil {
		return nil, fmt.Errorf("unrecognised node type '%s'", *nodeTypeStr)
	}

	cfg.HostID = gethcommon.HexToAddress(*hostID)
	cfg.HostAddress = *hostAddress
	cfg.Address = *address
	cfg.NodeType = nodeType
	cfg.L1ChainID = *l1ChainID
	cfg.ObscuroChainID = *obscuroChainID
	cfg.WillAttest = *willAttest
	cfg.ValidateL1Blocks = *validateL1Blocks
	cfg.ManagementContractAddress = gethcommon.HexToAddress(*managementContractAddress)
	cfg.LogLevel = *loglevel
	cfg.LogPath = *logPath
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.EdgelessDBHost = *edgelessDBHost
	cfg.SqliteDBPath = *sqliteDBPath
	cfg.ProfilerEnabled = *profilerEnabled
	cfg.MinGasPrice = big.NewInt(*minGasPrice)
	cfg.MessageBusAddress = gethcommon.HexToAddress(*messageBusAddress)
	cfg.SequencerID = gethcommon.HexToAddress(*sequencerID)
	cfg.ObscuroGenesis = *obscuroGenesis
	cfg.DebugNamespaceEnabled = *debugNamespaceEnabled
	cfg.MaxBatchSize = *maxBatchSize
	cfg.MaxRollupSize = *maxRollupSize
	cfg.BaseFee = big.NewInt(0).SetUint64(*baseFee)
	cfg.GasPaymentAddress = gethcommon.HexToAddress(*coinbaseAddress)
	cfg.GasLimit = big.NewInt(0).SetUint64(*gasLimit)

	return cfg, nil
}

// Parses the config from the .toml file at configPath.
func fileBasedConfig(configPath string) (*config.EnclaveConfig, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var tomlConfig EnclaveConfigToml
	err = toml.Unmarshal(bytes, &tomlConfig)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	nodeType, err := common.ToNodeType(tomlConfig.NodeType)
	if err != nil {
		return nil, fmt.Errorf("unrecognised node type '%s'", tomlConfig.NodeType)
	}

	return &config.EnclaveConfig{
		HostID:                    gethcommon.HexToAddress(tomlConfig.HostID),
		HostAddress:               tomlConfig.HostAddress,
		Address:                   tomlConfig.Address,
		NodeType:                  nodeType,
		L1ChainID:                 tomlConfig.L1ChainID,
		ObscuroChainID:            tomlConfig.ObscuroChainID,
		WillAttest:                tomlConfig.WillAttest,
		ValidateL1Blocks:          tomlConfig.ValidateL1Blocks,
		ManagementContractAddress: gethcommon.HexToAddress(tomlConfig.ManagementContractAddress),
		LogLevel:                  tomlConfig.LogLevel,
		LogPath:                   tomlConfig.LogPath,
		UseInMemoryDB:             tomlConfig.UseInMemoryDB,
		GenesisJSON:               []byte(tomlConfig.GenesisJSON),
		EdgelessDBHost:            tomlConfig.EdgelessDBHost,
		SqliteDBPath:              tomlConfig.SqliteDBPath,
		ProfilerEnabled:           tomlConfig.ProfilerEnabled,
	}, nil
}
