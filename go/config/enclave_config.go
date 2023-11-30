package config

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/flag"
	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// EnclaveConfig contains the full configuration for an Obscuro enclave service.
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
	ObscuroChainID int64
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
	ObscuroGenesis string
	// Whether debug calls are available
	DebugNamespaceEnabled bool
	// Maximum bytes a batch can be uncompressed.
	MaxBatchSize uint64
	// MaxRollupSize - configured to be close to what the ethereum clients
	// have configured as the maximum size a transaction can have. Note that this isn't
	// a protocol limit, but a miner imposed limit and it might be hard to find someone
	// to include a transaction if it goes above it
	MaxRollupSize uint64

	GasPaymentAddress gethcommon.Address
	BaseFee           *big.Int
	GasLimit          *big.Int
}

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() *EnclaveConfig {
	return &EnclaveConfig{
		HostID:                    gethcommon.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		NodeType:                  common.Sequencer,
		L1ChainID:                 1337,
		ObscuroChainID:            443,
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
		MinGasPrice:               big.NewInt(1),
		SequencerID:               gethcommon.BytesToAddress([]byte("")),
		ObscuroGenesis:            "",
		DebugNamespaceEnabled:     false,
		MaxBatchSize:              1024 * 25,
		MaxRollupSize:             1024 * 64,
		GasPaymentAddress:         gethcommon.HexToAddress("0xd6C9230053f45F873Cb66D8A02439380a37A4fbF"),
		BaseFee:                   new(big.Int).SetUint64(1),
		GasLimit:                  new(big.Int).SetUint64(params.MaxGasLimit / 6),
	}
}

func FromFlags(flagMap map[string]*flag.TenFlag) (*EnclaveConfig, error) {
	flagsTestMode := false

	// check if it's in test mode or not
	val := os.Getenv("EDG_TESTMODE")
	if val == "true" {
		flagsTestMode = true
	} else {
		fmt.Println("Using mandatory signed configurations.")
	}

	if !flagsTestMode {
		envFlags, err := retrieveEnvFlags()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve env flags - %w", err)
		}

		// create the final flag usage
		parsedFlags := map[string]*flag.TenFlag{}
		for flagName, cliflag := range flagMap {
			parsedFlags[flagName] = cliflag
		}
		// env flags override CLI flags
		for flagName, envflag := range envFlags {
			parsedFlags[flagName] = envflag
		}

		return newConfig(parsedFlags)
	}
	return newConfig(flagMap)
}

func retrieveEnvFlags() (map[string]*flag.TenFlag, error) {
	parsedFlags := map[string]*flag.TenFlag{}

	for _, eflag := range enclaveRestrictedFlags {
		val := os.Getenv("EDG_" + strings.ToUpper(eflag.Name))

		// all env flags must be set
		if val == "" {
			return nil, fmt.Errorf("env var not set: %s", eflag.Name)
		}

		switch eflag.FlagType {
		case "string":
			parsedFlag := flag.NewStringFlag(eflag.Name, "", "")
			parsedFlag.Value = val

			parsedFlags[eflag.Name] = parsedFlag
		case "int64":
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse flag %s - %w", eflag.Name, err)
			}

			parsedFlag := flag.NewIntFlag(eflag.Name, 0, "")
			parsedFlag.Value = i
			parsedFlags[eflag.Name] = parsedFlag
		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("unable to parse flag %s - %w", eflag.Name, err)
			}

			parsedFlag := flag.NewBoolFlag(eflag.Name, false, "")
			parsedFlag.Value = b
			parsedFlags[eflag.Name] = parsedFlag
		default:
			return nil, fmt.Errorf("unexpected type: %s", eflag.FlagType)
		}
	}
	return parsedFlags, nil
}

func newConfig(flags map[string]*flag.TenFlag) (*EnclaveConfig, error) {
	cfg := DefaultEnclaveConfig()

	nodeType, err := common.ToNodeType(flags[NodeTypeFlag].String())
	if err != nil {
		return nil, fmt.Errorf("unrecognised node type '%s'", flags[NodeTypeFlag].String())
	}

	cfg.HostID = gethcommon.HexToAddress(flags[HostIDFlag].String())
	cfg.HostAddress = flags[HostAddressFlag].String()
	cfg.Address = flags[AddressFlag].String()
	cfg.NodeType = nodeType
	cfg.L1ChainID = flags[L1ChainIDFlag].Int64()
	cfg.ObscuroChainID = flags[ObscuroChainIDFlag].Int64()
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
	cfg.ObscuroGenesis = flags[ObscuroGenesisFlag].String()
	cfg.DebugNamespaceEnabled = flags[DebugNamespaceEnabledFlag].Bool()
	cfg.MaxBatchSize = flags[MaxBatchSizeFlag].Uint64()
	cfg.MaxRollupSize = flags[MaxRollupSizeFlag].Uint64()
	cfg.BaseFee = big.NewInt(0).SetUint64(flags[L2BaseFeeFlag].Uint64())
	cfg.GasPaymentAddress = gethcommon.HexToAddress(flags[L2CoinbaseFlag].String())
	cfg.GasLimit = big.NewInt(0).SetUint64(flags[L2GasLimitFlag].Uint64())

	return cfg, nil
}
