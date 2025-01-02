package config

import (
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/config"
)

// For now, this is the bridge between TenConfig and the config used internally by the enclave service.

// EnclaveConfig contains the full configuration for an Obscuro enclave service.
type EnclaveConfig struct {
	// The identity of the host the enclave service is tied to
	HostID gethcommon.Address
	// The public peer-to-peer IP address of the host the enclave service is tied to
	HostAddress string
	// The address on which to serve requests
	Address string
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
	// SystemContractOwner is the address that owns the system contracts
	SystemContractOwner gethcommon.Address
	// P2P address for validators to connect to the sequencer for live batch data
	SequencerP2PAddress string
	// A json string that specifies the prefunded addresses at the genesis of the TEN network
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

	// RPCTimeout - calls that are longer than this will be cancelled, to prevent resource starvation
	// normally, the context is propagated from the host, but in some cases ( like the evm, we have to create a context)
	RPCTimeout time.Duration

	// StoreExecutedTransactions is a flag that instructs the current enclave to store data required to answer RPC queries.
	StoreExecutedTransactions bool
}

func EnclaveConfigFromTenConfig(tenCfg *config.TenConfig) *EnclaveConfig {
	return &EnclaveConfig{
		HostID:                    tenCfg.Node.ID,
		HostAddress:               tenCfg.Node.HostAddress,
		WillAttest:                tenCfg.Enclave.EnableAttestation,
		StoreExecutedTransactions: tenCfg.Enclave.StoreExecutedTransactions,

		ObscuroChainID:      tenCfg.Network.ChainID,
		SequencerP2PAddress: tenCfg.Network.Sequencer.P2PAddress,

		Address:    tenCfg.Enclave.RPC.BindAddress,
		RPCTimeout: tenCfg.Enclave.RPC.Timeout,

		L1ChainID:                 tenCfg.Network.L1.ChainID,
		ValidateL1Blocks:          tenCfg.Enclave.L1.EnableBlockValidation,
		GenesisJSON:               tenCfg.Enclave.L1.GenesisJSON,
		ManagementContractAddress: tenCfg.Network.L1.L1Contracts.ManagementContract,
		MessageBusAddress:         tenCfg.Network.L1.L1Contracts.MessageBusContract,
		SystemContractOwner:       tenCfg.Network.Sequencer.SystemContractsUpgrader,
		LogLevel:                  tenCfg.Enclave.Log.Level,
		LogPath:                   tenCfg.Enclave.Log.Path,

		UseInMemoryDB:  tenCfg.Enclave.DB.UseInMemory,
		EdgelessDBHost: tenCfg.Enclave.DB.EdgelessDBHost,
		SqliteDBPath:   tenCfg.Enclave.DB.SqlitePath,

		ProfilerEnabled:       tenCfg.Enclave.Debug.EnableProfiler,
		DebugNamespaceEnabled: tenCfg.Enclave.Debug.EnableDebugNamespace,

		MinGasPrice:              tenCfg.Network.Gas.MinGasPrice,
		GasPaymentAddress:        tenCfg.Network.Gas.PaymentAddress,
		BaseFee:                  tenCfg.Network.Gas.BaseFee,
		GasBatchExecutionLimit:   tenCfg.Network.Gas.BatchExecutionLimit,
		GasLocalExecutionCapFlag: tenCfg.Network.Gas.LocalExecutionCap,

		TenGenesis:    tenCfg.Network.GenesisJSON,
		MaxBatchSize:  tenCfg.Network.Batch.MaxSize,
		MaxRollupSize: tenCfg.Network.Rollup.MaxSize,
	}
}
