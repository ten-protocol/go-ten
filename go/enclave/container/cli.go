package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func ParseConfig() (*config.EnclaveConfig, error) {
	cfg, err := loadDefaultEnclaveInputConfig()
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	flagUsageMap := getFlagUsageMap()

	hostID := flag.String(HostIDFlag, cfg.HostID, flagUsageMap[HostIDFlag])
	hostAddress := flag.String(HostAddressFlag, cfg.HostAddress, flagUsageMap[HostAddressFlag])
	address := flag.String(AddressFlag, cfg.Address, flagUsageMap[AddressFlag])
	nodeType := flag.String(NodeTypeFlag, cfg.NodeType, flagUsageMap[NodeTypeFlag])
	l1ChainID := flag.Int64(L1ChainIDFlag, cfg.L1ChainID, flagUsageMap[L1ChainIDFlag])
	tenChainID := flag.Int64(TenChainIDFlag, cfg.TenChainID, flagUsageMap[TenGenesisFlag])
	willAttest := flag.Bool(WillAttestFlag, cfg.WillAttest, flagUsageMap[WillAttestFlag])
	validateL1Block := flag.Bool(ValidateL1BlocksFlag, cfg.ValidateL1Blocks, flagUsageMap[ValidateL1BlocksFlag])
	managementContractAddress := flag.String(ManagementContractAddressFlag, cfg.ManagementContractAddress, flagUsageMap[ManagementContractAddressFlag])
	logLevel := flag.Int(LogLevelFlag, cfg.LogLevel, flagUsageMap[LogLevelFlag])
	logPath := flag.String(LogPathFlag, cfg.LogPath, flagUsageMap[LogPathFlag])
	useInMemoryDB := flag.Bool(UseInMemoryDBFlag, cfg.UseInMemoryDB, flagUsageMap[UseInMemoryDBFlag])
	edgelessDBHost := flag.String(EdgelessDBHostFlag, cfg.EdgelessDBHost, flagUsageMap[EdgelessDBHostFlag])
	sqliteDBPath := flag.String(SQLiteDBPathFlag, cfg.SqliteDBPath, flagUsageMap[SQLiteDBPathFlag])
	profilerEnabled := flag.Bool(ProfilerEnabledFlag, cfg.ProfilerEnabled, flagUsageMap[ProfilerEnabledFlag])
	minGasPrice := flag.Uint64(MinGasPriceFlag, cfg.MinGasPrice, flagUsageMap[MinGasPriceFlag])
	messageBusAddress := flag.String(MessageBusAddressFlag, cfg.MessageBusAddress, flagUsageMap[MessageBusAddressFlag])
	sequencerID := flag.String(SequencerIDFlag, cfg.SequencerID, flagUsageMap[SequencerIDFlag])
	tenGenesis := flag.String(TenGenesisFlag, cfg.TenGenesis, flagUsageMap[TenGenesisFlag])
	debugNamespaceEnabled := flag.Bool(DebugNamespaceEnabledFlag, cfg.DebugNamespaceEnabled, flagUsageMap[DebugNamespaceEnabledFlag])
	maxBatchSize := flag.Uint64(MaxBatchSizeFlag, cfg.MaxBatchSize, flagUsageMap[MaxBatchSizeFlag])
	maxRollupSize := flag.Uint64(MaxRollupSizeFlag, cfg.MaxRollupSize, flagUsageMap[MaxRollupSizeFlag])
	gasPaymentAddress := flag.String(L2CoinbaseFlag, cfg.GasPaymentAddress, flagUsageMap[L2CoinbaseFlag])
	baseFee := flag.Uint64(L2BaseFeeFlag, cfg.BaseFee, flagUsageMap[L2BaseFeeFlag])
	gasBatchExecutionLimit := flag.Uint64(GasBatchExecutionLimit, cfg.GasBatchExecutionLimit, flagUsageMap[GasBatchExecutionLimit])
	gasLocalExecutionCap := flag.Uint64(GasLocalExecutionCapFlag, cfg.GasLocalExecutionCapFlag, flagUsageMap[GasLocalExecutionCapFlag])

	flag.Parse()

	cfg.HostID = *hostID
	cfg.HostAddress = *hostAddress
	cfg.Address = *address
	cfg.NodeType = *nodeType
	cfg.L1ChainID = *l1ChainID
	cfg.TenChainID = *tenChainID
	cfg.WillAttest = *willAttest
	cfg.ValidateL1Blocks = *validateL1Block
	cfg.ManagementContractAddress = *managementContractAddress
	cfg.LogLevel = *logLevel
	cfg.LogPath = *logPath
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.EdgelessDBHost = *edgelessDBHost
	cfg.SqliteDBPath = *sqliteDBPath
	cfg.ProfilerEnabled = *profilerEnabled
	cfg.MinGasPrice = *minGasPrice
	cfg.MessageBusAddress = *messageBusAddress
	cfg.SequencerID = *sequencerID
	cfg.TenGenesis = *tenGenesis
	cfg.DebugNamespaceEnabled = *debugNamespaceEnabled
	cfg.MaxBatchSize = *maxBatchSize
	cfg.MaxRollupSize = *maxRollupSize
	cfg.GasPaymentAddress = *gasPaymentAddress
	cfg.BaseFee = *baseFee
	cfg.GasBatchExecutionLimit = *gasBatchExecutionLimit
	cfg.GasLocalExecutionCapFlag = *gasLocalExecutionCap

	cfg, err = retrieveOrSetEnclaveRestrictedFlags(cfg)

	enclaveConfig, err := cfg.ToEnclaveConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert EnclaveInputConfig to EnclaveConfig")
	}
	return enclaveConfig, nil
}

// retrieveOrSetEnclaveRestrictedFlags, ensures relevant flags are able to pass into the enclave in scenarios where
// an `ego sign` procedure isn't enabled - and no env[] array used. In this case, it will take the EDG_ env vars as
// first or the default config values as fall-back.
func retrieveOrSetEnclaveRestrictedFlags(cfg *config.EnclaveInputConfig) (*config.EnclaveInputConfig, error) {
	val := os.Getenv("EDG_TESTMODE")
	if val == "true" {
		fmt.Println("Using test mode flags")
		return cfg, nil
	} else {
		fmt.Println("Using mandatory signed configurations.")
	}

	v := reflect.ValueOf(cfg).Elem() // Get the reflect.Value of the struct

	for eFlag, flagType := range enclaveRestrictedFlags {
		eFlag = capitalizeFirst(eFlag)
		targetEnvVar := "EDG_" + strings.ToUpper(eFlag)
		val := os.Getenv(targetEnvVar)
		if val == "" {
			fieldVal := v.FieldByName(eFlag) // Access the struct field by name
			if !fieldVal.IsValid() {
				panic("No valid field found for flag " + eFlag)
			}

			var strVal string
			switch flagType {
			case "int64":
				strVal = strconv.FormatInt(fieldVal.Int(), 10)
			case "string":
				strVal = fieldVal.String()
			case "bool":
				strVal = strconv.FormatBool(fieldVal.Bool())
			default:
				panic("Unsupported type for field " + eFlag)
			}

			if err := os.Setenv(targetEnvVar, strVal); err != nil {
				panic("Failed to set environment variable " + targetEnvVar)
			}
			fmt.Printf("Set %s to %s from default configuration.\n", targetEnvVar, strVal)
		}
	}
	return cfg, nil
}

// capitalizeFirst capitalizes the first letter of the given string. handles mismatch between flag and config struct
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// loadDefaultEnclaveInputConfig parses optional or default configuration file and returns struct.
func loadDefaultEnclaveInputConfig() (*config.EnclaveInputConfig, error) {
	flagUsageMap := getFlagUsageMap()
	configPath := flag.String("config", "./go/config/templates/default_enclave_config.yaml", flagUsageMap["configFlag"])
	overridePath := flag.String("override", "", flagUsageMap["overrideFlag"])

	// Parse only once capturing all necessary flags
	flag.Parse()

	var err error
	conf, err := loadEnclaveConfigFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	// Apply overrides if the override path is provided
	if *overridePath != "" {
		overridesConf, err := loadEnclaveConfigFromFile(*overridePath)
		if err != nil {
			panic(err)
		}

		config.ApplyOverrides(conf, overridesConf)
	}

	return conf, nil
}

// loadEnclaveConfigFromFile reads configuration from a file and environment variables
func loadEnclaveConfigFromFile(configPath string) (*config.EnclaveInputConfig, error) {
	defaultConfig := &config.EnclaveInputConfig{}
	// Read YAML configuration
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
