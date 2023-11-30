package config

import "github.com/ten-protocol/go-ten/go/common/flag"

// EnclaveFlags are the flags that the enclave can receive
func EnclaveFlags() map[string]*flag.TenFlag {
	return map[string]*flag.TenFlag{
		HostIDFlag:                    flag.NewStringFlag(HostIDFlag, "", FlagDescriptionMap[HostIDFlag]),
		HostAddressFlag:               flag.NewStringFlag(HostAddressFlag, "", FlagDescriptionMap[HostAddressFlag]),
		AddressFlag:                   flag.NewStringFlag(AddressFlag, "", FlagDescriptionMap[HostAddressFlag]),
		NodeTypeFlag:                  flag.NewStringFlag(NodeTypeFlag, "", FlagDescriptionMap[NodeTypeFlag]),
		WillAttestFlag:                flag.NewBoolFlag(WillAttestFlag, false, FlagDescriptionMap[WillAttestFlag]),
		ValidateL1BlocksFlag:          flag.NewBoolFlag(ValidateL1BlocksFlag, false, FlagDescriptionMap[ValidateL1BlocksFlag]),
		ManagementContractAddressFlag: flag.NewStringFlag(ManagementContractAddressFlag, "", FlagDescriptionMap[ManagementContractAddressFlag]),
		LogLevelFlag:                  flag.NewIntFlag(LogLevelFlag, 0, FlagDescriptionMap[LogLevelFlag]),
		LogPathFlag:                   flag.NewStringFlag(LogPathFlag, "", FlagDescriptionMap[LogPathFlag]),
		EdgelessDBHostFlag:            flag.NewStringFlag(EdgelessDBHostFlag, "", FlagDescriptionMap[EdgelessDBHostFlag]),
		SQLiteDBPathFlag:              flag.NewStringFlag(SQLiteDBPathFlag, "", FlagDescriptionMap[SQLiteDBPathFlag]),
		MinGasPriceFlag:               flag.NewInt64Flag(MinGasPriceFlag, 0, FlagDescriptionMap[MinGasPriceFlag]),
		MessageBusAddressFlag:         flag.NewStringFlag(MessageBusAddressFlag, "", FlagDescriptionMap[MessageBusAddressFlag]),
		SequencerIDFlag:               flag.NewStringFlag(SequencerIDFlag, "", FlagDescriptionMap[SequencerIDFlag]),
		MaxBatchSizeFlag:              flag.NewUint64Flag(MaxBatchSizeFlag, 0, FlagDescriptionMap[MaxBatchSizeFlag]),
		MaxRollupSizeFlag:             flag.NewUint64Flag(MaxRollupSizeFlag, 0, FlagDescriptionMap[MaxRollupSizeFlag]),
		L2BaseFeeFlag:                 flag.NewUint64Flag(L2BaseFeeFlag, 0, ""),
		L2CoinbaseFlag:                flag.NewStringFlag(L2CoinbaseFlag, "", ""),
		L2GasLimitFlag:                flag.NewUint64Flag(L2GasLimitFlag, 0, ""),
		ObscuroGenesisFlag:            flag.NewStringFlag(ObscuroGenesisFlag, "", FlagDescriptionMap[ObscuroGenesisFlag]),
		L1ChainIDFlag:                 flag.NewInt64Flag(L1ChainIDFlag, 0, FlagDescriptionMap[L1ChainIDFlag]),
		ObscuroChainIDFlag:            flag.NewInt64Flag(ObscuroChainIDFlag, 0, FlagDescriptionMap[ObscuroChainIDFlag]),
		UseInMemoryDBFlag:             flag.NewBoolFlag(UseInMemoryDBFlag, false, FlagDescriptionMap[UseInMemoryDBFlag]),
		ProfilerEnabledFlag:           flag.NewBoolFlag(ProfilerEnabledFlag, false, FlagDescriptionMap[ProfilerEnabledFlag]),
		DebugNamespaceEnabledFlag:     flag.NewBoolFlag(DebugNamespaceEnabledFlag, false, FlagDescriptionMap[DebugNamespaceEnabledFlag]),
	}
}

// enclaveRestrictedFlags are the flags that the enclave can receive ONLY over the Ego signed enclave.json
var enclaveRestrictedFlags = map[string]*flag.TenFlag{
	L1ChainIDFlag:             flag.NewInt64Flag(L1ChainIDFlag, 0, FlagDescriptionMap[L1ChainIDFlag]),
	ObscuroChainIDFlag:        flag.NewInt64Flag(ObscuroChainIDFlag, 0, FlagDescriptionMap[ObscuroChainIDFlag]),
	ObscuroGenesisFlag:        flag.NewStringFlag(ObscuroGenesisFlag, "", FlagDescriptionMap[ObscuroGenesisFlag]),
	UseInMemoryDBFlag:         flag.NewBoolFlag(UseInMemoryDBFlag, false, FlagDescriptionMap[UseInMemoryDBFlag]),
	ProfilerEnabledFlag:       flag.NewBoolFlag(ProfilerEnabledFlag, false, FlagDescriptionMap[ProfilerEnabledFlag]),
	DebugNamespaceEnabledFlag: flag.NewBoolFlag(DebugNamespaceEnabledFlag, false, FlagDescriptionMap[DebugNamespaceEnabledFlag]),
}

// Flag names.
const (
	HostIDFlag                    = "hostID"
	HostAddressFlag               = "hostAddress"
	AddressFlag                   = "address"
	NodeTypeFlag                  = "nodeType"
	L1ChainIDFlag                 = "l1ChainID"
	ObscuroChainIDFlag            = "obscuroChainID"
	WillAttestFlag                = "willAttest"
	ValidateL1BlocksFlag          = "validateL1Blocks"
	ManagementContractAddressFlag = "managementContractAddress"
	LogLevelFlag                  = "logLevel"
	LogPathFlag                   = "logPath"
	UseInMemoryDBFlag             = "useInMemoryDB"
	EdgelessDBHostFlag            = "edgelessDBHost"
	SQLiteDBPathFlag              = "sqliteDBPath"
	ProfilerEnabledFlag           = "profilerEnabled"
	MinGasPriceFlag               = "minGasPrice"
	MessageBusAddressFlag         = "messageBusAddress"
	SequencerIDFlag               = "sequencerID"
	ObscuroGenesisFlag            = "obscuroGenesis"
	DebugNamespaceEnabledFlag     = "debugNamespaceEnabled"
	MaxBatchSizeFlag              = "maxBatchSize"
	MaxRollupSizeFlag             = "maxRollupSize"
	L2BaseFeeFlag                 = "l2BaseFee"
	L2CoinbaseFlag                = "l2Coinbase"
	L2GasLimitFlag                = "l2GasLimit"
)

var FlagDescriptionMap = map[string]string{
	HostIDFlag:                    "The 20 bytes of the address of the Obscuro host this enclave serves",
	HostAddressFlag:               "The peer-to-peer IP address of the Obscuro host this enclave serves",
	AddressFlag:                   "The address on which to serve the Obscuro enclave service",
	NodeTypeFlag:                  "The node's type (e.g. sequencer, validator)",
	L1ChainIDFlag:                 "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)",
	ObscuroChainIDFlag:            "An integer representing the unique chain id of the Obscuro chain (default 443)",
	WillAttestFlag:                "Whether the enclave will produce a verified attestation report",
	ValidateL1BlocksFlag:          "Whether to validate incoming blocks using the hardcoded L1 genesis.json config",
	ManagementContractAddressFlag: "The management contract address on the L1",
	LogLevelFlag:                  "The verbosity level of logs. (Defaults to Info)",
	LogPathFlag:                   "The path to use for the enclave service's log file",
	UseInMemoryDBFlag:             "Whether the enclave will use an in-memory DB rather than persist data",
	EdgelessDBHostFlag:            "Host address for the edgeless DB instance (can be empty if useInMemoryDB is true or if not using attestation",
	SQLiteDBPathFlag:              "Filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB or if using attestation/EdgelessDB)",
	ProfilerEnabledFlag:           "Runs a profiler instance (Defaults to false)",
	MinGasPriceFlag:               "The minimum gas price for mining a transaction",
	MessageBusAddressFlag:         "The address of the L1 message bus contract owned by the management contract.",
	SequencerIDFlag:               "The 20 bytes of the address of the sequencer for this network",
	ObscuroGenesisFlag:            "The json string with the obscuro genesis",
	DebugNamespaceEnabledFlag:     "Whether the debug namespace is enabled",
	MaxBatchSizeFlag:              "The maximum size a batch is allowed to reach uncompressed",
	MaxRollupSizeFlag:             "The maximum size a rollup is allowed to reach",
}
