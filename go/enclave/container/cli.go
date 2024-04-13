package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func ParseConfig() (*config.EnclaveConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Enclave)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	cfg := inputCfg.(*config.EnclaveInputConfig) // assert

	flagUsageMap := getFlagUsageMap()

	hostID := flag.String(HostIDFlag, config.GetEnvString(HostIDFlag, cfg.HostID), flagUsageMap[HostIDFlag])
	hostAddress := flag.String(HostAddressFlag, config.GetEnvString(HostAddressFlag, cfg.HostAddress), flagUsageMap[HostAddressFlag])
	address := flag.String(AddressFlag, config.GetEnvString(AddressFlag, cfg.Address), flagUsageMap[AddressFlag])
	nodeType := flag.String(NodeTypeFlag, config.GetEnvString(NodeTypeFlag, cfg.NodeType), flagUsageMap[NodeTypeFlag])
	l1ChainID := flag.Int64(L1ChainIDFlag, config.GetEnvInt64(L1ChainIDFlag, cfg.L1ChainID), flagUsageMap[L1ChainIDFlag])
	tenChainID := flag.Int64(TenChainIDFlag, config.GetEnvInt64(TenChainIDFlag, cfg.TenChainID), flagUsageMap[TenGenesisFlag])
	willAttest := flag.Bool(WillAttestFlag, config.GetEnvBool(WillAttestFlag, cfg.WillAttest), flagUsageMap[WillAttestFlag])
	validateL1Block := flag.Bool(ValidateL1BlocksFlag, config.GetEnvBool(ValidateL1BlocksFlag, cfg.ValidateL1Blocks), flagUsageMap[ValidateL1BlocksFlag])
	managementContractAddress := flag.String(ManagementContractAddressFlag, config.GetEnvString(ManagementContractAddressFlag, cfg.ManagementContractAddress), flagUsageMap[ManagementContractAddressFlag])
	logLevel := flag.Int(LogLevelFlag, config.GetEnvInt(LogLevelFlag, cfg.LogLevel), flagUsageMap[LogLevelFlag])
	logPath := flag.String(LogPathFlag, config.GetEnvString(LogPathFlag, cfg.LogPath), flagUsageMap[LogPathFlag])
	useInMemoryDB := flag.Bool(UseInMemoryDBFlag, config.GetEnvBool(UseInMemoryDBFlag, cfg.UseInMemoryDB), flagUsageMap[UseInMemoryDBFlag])
	edgelessDBHost := flag.String(EdgelessDBHostFlag, config.GetEnvString(EdgelessDBHostFlag, cfg.EdgelessDBHost), flagUsageMap[EdgelessDBHostFlag])
	sqliteDBPath := flag.String(SQLiteDBPathFlag, config.GetEnvString(SQLiteDBPathFlag, cfg.SqliteDBPath), flagUsageMap[SQLiteDBPathFlag])
	profilerEnabled := flag.Bool(ProfilerEnabledFlag, config.GetEnvBool(ProfilerEnabledFlag, cfg.ProfilerEnabled), flagUsageMap[ProfilerEnabledFlag])
	minGasPrice := flag.Uint64(MinGasPriceFlag, config.GetEnvUint64(MinGasPriceFlag, cfg.MinGasPrice), flagUsageMap[MinGasPriceFlag])
	messageBusAddress := flag.String(MessageBusAddressFlag, config.GetEnvString(MessageBusAddressFlag, cfg.MessageBusAddress), flagUsageMap[MessageBusAddressFlag])
	sequencerID := flag.String(SequencerIDFlag, config.GetEnvString(SequencerIDFlag, cfg.SequencerID), flagUsageMap[SequencerIDFlag])
	tenGenesis := flag.String(TenGenesisFlag, config.GetEnvString(TenGenesisFlag, cfg.TenGenesis), flagUsageMap[TenGenesisFlag])
	debugNamespaceEnabled := flag.Bool(DebugNamespaceEnabledFlag, config.GetEnvBool(DebugNamespaceEnabledFlag, cfg.DebugNamespaceEnabled), flagUsageMap[DebugNamespaceEnabledFlag])
	maxBatchSize := flag.Uint64(MaxBatchSizeFlag, config.GetEnvUint64(MaxBatchSizeFlag, cfg.MaxBatchSize), flagUsageMap[MaxBatchSizeFlag])
	maxRollupSize := flag.Uint64(MaxRollupSizeFlag, config.GetEnvUint64(MaxRollupSizeFlag, cfg.MaxRollupSize), flagUsageMap[MaxRollupSizeFlag])
	gasPaymentAddress := flag.String(L2CoinbaseFlag, config.GetEnvString(L2CoinbaseFlag, cfg.GasPaymentAddress), flagUsageMap[L2CoinbaseFlag])
	baseFee := flag.Uint64(L2BaseFeeFlag, config.GetEnvUint64(L2BaseFeeFlag, cfg.BaseFee), flagUsageMap[L2BaseFeeFlag])
	gasBatchExecutionLimit := flag.Uint64(GasBatchExecutionLimit, config.GetEnvUint64(GasBatchExecutionLimit, cfg.GasBatchExecutionLimit), flagUsageMap[GasBatchExecutionLimit])
	gasLocalExecutionCap := flag.Uint64(GasLocalExecutionCapFlag, config.GetEnvUint64(GasLocalExecutionCapFlag, cfg.GasLocalExecutionCapFlag), flagUsageMap[GasLocalExecutionCapFlag])
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

			if strVal == "" {
				panic("Invalid default or EDG_ for " + eFlag)
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
