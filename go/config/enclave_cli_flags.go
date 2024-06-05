package config

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/flag"
)

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
	ObscuroGenesisFlag            = "obscuroGenesis"
	DebugNamespaceEnabledFlag     = "debugNamespaceEnabled"
	MaxBatchSizeFlag              = "maxBatchSize"
	MaxRollupSizeFlag             = "maxRollupSize"
	L2BaseFeeFlag                 = "l2BaseFee"
	L2CoinbaseFlag                = "l2Coinbase"
	GasBatchExecutionLimit        = "gasBatchExecutionLimit"
	GasLocalExecutionCapFlag      = "gasLocalExecutionCap"
)

// EnclaveFlags are the flags that the enclave can receive
var EnclaveFlags = map[string]*flag.TenFlag{
	HostIDFlag:                    flag.NewStringFlag(HostIDFlag, "", "The 20 bytes of the address of the Obscuro host this enclave serves"),
	HostAddressFlag:               flag.NewStringFlag(HostAddressFlag, "127.0.0.1:10000", "The peer-to-peer IP address of the Obscuro host this enclave serves"),
	AddressFlag:                   flag.NewStringFlag(AddressFlag, "127.0.0.1:11000", "The address on which to serve the Obscuro enclave service"),
	NodeTypeFlag:                  flag.NewStringFlag(NodeTypeFlag, common.Sequencer.String(), "The node's type (e.g. sequencer, validator)"),
	WillAttestFlag:                flag.NewBoolFlag(WillAttestFlag, false, "Whether the enclave will produce a verified attestation report"),
	ValidateL1BlocksFlag:          flag.NewBoolFlag(ValidateL1BlocksFlag, false, "Whether to validate incoming blocks using the hardcoded L1 genesis.json config"),
	ManagementContractAddressFlag: flag.NewStringFlag(ManagementContractAddressFlag, "", "The management contract address on the L1"),
	LogLevelFlag:                  flag.NewIntFlag(LogLevelFlag, 3, "The verbosity level of logs. (Defaults to Info)"),
	LogPathFlag:                   flag.NewStringFlag(LogPathFlag, "stdout", "The path to use for the enclave service's log file"),
	EdgelessDBHostFlag:            flag.NewStringFlag(EdgelessDBHostFlag, "", "Host address for the edgeless DB instance (can be empty if useInMemoryDB is true or if not using attestation"),
	SQLiteDBPathFlag:              flag.NewStringFlag(SQLiteDBPathFlag, "", "Filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB or if using attestation/EdgelessDB)"),
	MinGasPriceFlag:               flag.NewInt64Flag(MinGasPriceFlag, 1, "The minimum gas price for mining a transaction"),
	MessageBusAddressFlag:         flag.NewStringFlag(MessageBusAddressFlag, "", "The address of the L1 message bus contract owned by the management contract."),
	MaxBatchSizeFlag:              flag.NewUint64Flag(MaxBatchSizeFlag, 1024*55, "The maximum size a batch is allowed to reach uncompressed"),
	MaxRollupSizeFlag:             flag.NewUint64Flag(MaxRollupSizeFlag, 1024*128, "The maximum size a rollup is allowed to reach"),
	L2BaseFeeFlag:                 flag.NewUint64Flag(L2BaseFeeFlag, params.InitialBaseFee, ""),
	L2CoinbaseFlag:                flag.NewStringFlag(L2CoinbaseFlag, "0xd6C9230053f45F873Cb66D8A02439380a37A4fbF", ""),
	GasBatchExecutionLimit:        flag.NewUint64Flag(GasBatchExecutionLimit, 3_000_000_000, "Max gas that can be executed in a single batch"),
	ObscuroGenesisFlag:            flag.NewStringFlag(ObscuroGenesisFlag, "", "The json string with the obscuro genesis"),
	L1ChainIDFlag:                 flag.NewInt64Flag(L1ChainIDFlag, 1337, "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)"),
	ObscuroChainIDFlag:            flag.NewInt64Flag(ObscuroChainIDFlag, 443, "An integer representing the unique chain id of the Obscuro chain (default 443)"),
	UseInMemoryDBFlag:             flag.NewBoolFlag(UseInMemoryDBFlag, true, "Whether the enclave will use an in-memory DB rather than persist data"),
	ProfilerEnabledFlag:           flag.NewBoolFlag(ProfilerEnabledFlag, false, "Runs a profiler instance (Defaults to false)"),
	DebugNamespaceEnabledFlag:     flag.NewBoolFlag(DebugNamespaceEnabledFlag, false, "Whether the debug namespace is enabled"),
	GasLocalExecutionCapFlag:      flag.NewUint64Flag(GasLocalExecutionCapFlag, 4_000_000_000, "Max gas usage when executing local transactions"),
}

// enclaveRestrictedFlags are the flags that the enclave can receive ONLY over the Ego signed enclave.json
var enclaveRestrictedFlags = []string{
	L1ChainIDFlag,
	ObscuroChainIDFlag,
	ObscuroGenesisFlag,
	UseInMemoryDBFlag,
	ProfilerEnabledFlag,
	DebugNamespaceEnabledFlag,
}
