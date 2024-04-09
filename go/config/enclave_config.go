package config

import (
	"fmt"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/log"
	"math/big"
	"os"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/flag"
)

// EnclaveConfig contains the full configuration for a Ten enclave service.
type EnclaveConfig struct {
	// The identity of the host the enclave service is tied to
	HostID gethcommon.Address
	// The public peer-to-peer IP address of the host the enclave service is tied to
	HostAddress string
	// The address on which to serve requests
	Address string
	// The type of the node.
	NodeType common.NodeType
	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	TenChainID int64
	// Whether to produce a verified attestation report
	WillAttest bool
	// Whether to validate incoming L1 blocks
	ValidateL1Blocks bool
	// When validating incoming blocks, the genesis config for the L1 chain
	GenesisJSON []byte
	// The management contract address on the L1 network
	ManagementContractAddress gethcommon.Address
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// The path that the enclave's logs are written to
	LogPath string
	// Whether the enclave should use in-memory or persistent storage
	UseInMemoryDB bool
	// host address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled)
	EdgelessDBHost string
	// filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	//	if using InMemory DB or if attestation is enabled)
	SqliteDBPath string
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
	// MinGasPrice is the minimum gas price for mining a transaction
	MinGasPrice *big.Int
	// MessageBus L1 Address
	MessageBusAddress gethcommon.Address
	// The identity of the sequencer for the network
	SequencerID gethcommon.Address
	// A json string that specifies the prefunded addresses at the genesis of the Obscuro network
	TenGenesis string
	// Whether debug calls are available
	DebugNamespaceEnabled bool
	// Maximum bytes a batch can be uncompressed.
	MaxBatchSize uint64
	// MaxRollupSize - configured to be close to what the ethereum clients
	// have configured as the maximum size a transaction can have. Note that this isn't
	// a protocol limit, but a miner imposed limit and it might be hard to find someone
	// to include a transaction if it goes above it
	MaxRollupSize uint64

	GasPaymentAddress        gethcommon.Address
	BaseFee                  *big.Int
	GasBatchExecutionLimit   uint64
	GasLocalExecutionCapFlag uint64
}

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() *EnclaveConfig {
	return &EnclaveConfig{
		HostID:                    gethcommon.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		NodeType:                  common.Sequencer,
		L1ChainID:                 1337,
		TenChainID:                443,
		WillAttest:                false, // todo (config) - attestation should be on by default before production release
		ValidateL1Blocks:          false,
		GenesisJSON:               nil,
		ManagementContractAddress: gethcommon.BytesToAddress([]byte("")),
		LogLevel:                  int(gethlog.LvlInfo),
		LogPath:                   log.SysOut,
		UseInMemoryDB:             true, // todo (config) - persistence should be on by default before production release
		EdgelessDBHost:            "",
		SqliteDBPath:              "",
		ProfilerEnabled:           false,
		MinGasPrice:               big.NewInt(params.InitialBaseFee),
		SequencerID:               gethcommon.BytesToAddress([]byte("")),
		TenGenesis:                "",
		DebugNamespaceEnabled:     false,
		MaxBatchSize:              1024 * 36,
		MaxRollupSize:             1024 * 64,
		GasPaymentAddress:         gethcommon.HexToAddress("0xd6C9230053f45F873Cb66D8A02439380a37A4fbF"),
		BaseFee:                   new(big.Int).SetUint64(params.InitialBaseFee),

		// Due to hiding L1 costs in the gas quantity, the gas limit needs to be huge
		// Arbitrum with the same approach has gas limit of 1,125,899,906,842,624,
		// whilst the usage is small. Should be ok since execution is paid for anyway.
		GasLocalExecutionCapFlag: 300_000_000_000,
		GasBatchExecutionLimit:   300_000_000_000,
	}
}

func NewConfigFromFlags(cliFlags map[string]*flag.TenFlag) (*EnclaveConfig, error) {
	productionMode := true

	// check if it's in production mode or not
	val := os.Getenv("EDG_TESTMODE")
	if val == "true" {
		productionMode = false
		fmt.Println("Using test mode flags")
	} else {
		fmt.Println("Using mandatory signed configurations.")
	}

	if productionMode {
		envFlags, err := retrieveEnvFlags()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve env flags - %w", err)
		}

		// fail if any restricted flag is set via the cli
		for _, envflag := range envFlags {
			if cliflag, ok := cliFlags[envflag.Name]; ok && cliflag.IsSet() {
				return nil, fmt.Errorf("restricted flag was set: %s", cliflag.Name)
			}
		}

		// create the final flag usage
		parsedFlags := map[string]*flag.TenFlag{}
		for flagName, cliflag := range cliFlags {
			parsedFlags[flagName] = cliflag
		}
		// env flags override CLI flags
		for flagName, envflag := range envFlags {
			parsedFlags[flagName] = envflag
		}

		return newConfig(parsedFlags)
	}
	return newConfig(cliFlags)
}

func retrieveEnvFlags() (map[string]*flag.TenFlag, error) {
	parsedFlags := map[string]*flag.TenFlag{}

	for _, eflag := range enclaveRestrictedFlags {
		val := os.Getenv("EDG_" + strings.ToUpper(eflag))

		// all env flags must be set
		if val == "" {
			return nil, fmt.Errorf("env var not set: %s", eflag)
		}

		switch EnclaveFlags[eflag].FlagType {
		case "string":
			parsedFlag := flag.NewStringFlag(eflag, "", "")
			parsedFlag.Value = val

			parsedFlags[eflag] = parsedFlag
		case "int64":
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse flag %s - %w", eflag, err)
			}

			parsedFlag := flag.NewIntFlag(eflag, 0, "")
			parsedFlag.Value = i
			parsedFlags[eflag] = parsedFlag
		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("unable to parse flag %s - %w", eflag, err)
			}

			parsedFlag := flag.NewBoolFlag(eflag, false, "")
			parsedFlag.Value = b
			parsedFlags[eflag] = parsedFlag
		default:
			return nil, fmt.Errorf("unexpected type: %s", EnclaveFlags[eflag].FlagType)
		}
	}
	return parsedFlags, nil
}

func newConfig(flags map[string]*flag.TenFlag) (*EnclaveConfig, error) {
	cfg := &EnclaveConfig{}

	nodeType, err := common.ToNodeType(flags[NodeTypeFlag].String())
	if err != nil {
		return nil, fmt.Errorf("unrecognised node type '%s'", flags[NodeTypeFlag].String())
	}

	cfg.HostID = gethcommon.HexToAddress(flags[HostIDFlag].String())
	cfg.HostAddress = flags[HostAddressFlag].String()
	cfg.Address = flags[AddressFlag].String()
	cfg.NodeType = nodeType
	cfg.L1ChainID = flags[L1ChainIDFlag].Int64()
	cfg.TenChainID = flags[ObscuroChainIDFlag].Int64()
	cfg.WillAttest = flags[WillAttestFlag].Bool()
	cfg.ValidateL1Blocks = flags[ValidateL1BlocksFlag].Bool()
	cfg.ManagementContractAddress = gethcommon.HexToAddress(flags[ManagementContractAddressFlag].String())
	cfg.LogLevel = flags[LogLevelFlag].Int()
	cfg.LogPath = flags[LogPathFlag].String()
	cfg.UseInMemoryDB = flags[UseInMemoryDBFlag].Bool()
	cfg.EdgelessDBHost = flags[EdgelessDBHostFlag].String()
	cfg.SqliteDBPath = flags[SQLiteDBPathFlag].String()
	cfg.ProfilerEnabled = flags[ProfilerEnabledFlag].Bool()
	cfg.MinGasPrice = big.NewInt(flags[MinGasPriceFlag].Int64())
	cfg.MessageBusAddress = gethcommon.HexToAddress(flags[MessageBusAddressFlag].String())
	cfg.SequencerID = gethcommon.HexToAddress(flags[SequencerIDFlag].String())
	cfg.TenGenesis = flags[ObscuroGenesisFlag].String()
	cfg.DebugNamespaceEnabled = flags[DebugNamespaceEnabledFlag].Bool()
	cfg.MaxBatchSize = flags[MaxBatchSizeFlag].Uint64()
	cfg.MaxRollupSize = flags[MaxRollupSizeFlag].Uint64()
	cfg.BaseFee = big.NewInt(0).SetUint64(flags[L2BaseFeeFlag].Uint64())
	cfg.GasPaymentAddress = gethcommon.HexToAddress(flags[L2CoinbaseFlag].String())
	cfg.GasBatchExecutionLimit = flags[GasBatchExecutionLimit].Uint64()
	cfg.GasLocalExecutionCapFlag = flags[GasLocalExecutionCapFlag].Uint64()

	return cfg, nil
}
