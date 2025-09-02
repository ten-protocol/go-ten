package config

import (
	"fmt"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/config"
	"kythe.io/kythe/go/util/datasize"
)

// For now, this is the bridge between TenConfig and the config used internally by the enclave service.

// EnclaveConfig contains the full configuration for an TEN enclave service.
type EnclaveConfig struct {
	// **Consensus configs - must be the same for all nodes. Included in the signed image.
	// Whether to produce a verified attestation report
	WillAttest bool

	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	TenChainID int64

	//// These L1 contracts must be already deployed before the TEN network is created
	// L1 contract address that maintains all the index of all system/ network contracts
	NetworkConfigAddress gethcommon.Address
	// Rollup contract L1 Address
	DataAvailabilityRegistryAddress gethcommon.Address
	// EnclaveRegistry L1 Address
	EnclaveRegistryAddress gethcommon.Address
	// MessageBus L1 Address
	MessageBusAddress gethcommon.Address
	// Bridge L1 Address
	L1BridgeAddress gethcommon.Address
	// SystemContractOwner is the address that owns the system contracts
	SystemContractOwner gethcommon.Address

	// Maximum bytes a batch can be uncompressed.
	MaxBatchSize uint64
	// MaxRollupSize - configured to be close to what the ethereum clients
	// have configured as the maximum size a transaction can have. Note that this isn't
	// a protocol limit, but a miner imposed limit and it might be hard to find someone
	// to include a transaction if it goes above it
	MaxRollupSize uint64
	// MinGasPrice is the minimum gas price for mining a transaction
	MinGasPrice *big.Int
	// A json string that specifies the prefunded addresses at the genesis of the TEN network
	TenGenesis             string
	GasPaymentAddress      gethcommon.Address
	MinBaseFee             *big.Int
	GasBatchExecutionLimit uint64

	// **Db configs
	// Whether the enclave should use in-memory or persistent storage
	UseInMemoryDB bool
	// host address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled)
	EdgelessDBHost string
	// filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	//	if using InMemory DB or if attestation is enabled)
	SqliteDBPath string

	// **Networking cfgs
	// The address on which to serve requests
	RPCAddress string
	// RPCTimeout - calls that are longer than this will be cancelled, to prevent resource starvation
	// normally, the context is propagated from the host, but in some cases ( like the evm, we have to create a context)
	RPCTimeout time.Duration

	// **Running config
	// Arbitrary identification of the Node. Usually derived from the L1 wallet address. Useful for logging.
	NodeID string
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// The path that the enclave's logs are written to
	LogPath string
	// StoreExecutedTransactions is a flag that instructs the current enclave to store data required to answer RPC queries.
	StoreExecutedTransactions bool
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled          bool
	DebugNamespaceEnabled    bool
	GasLocalExecutionCapFlag uint64
	DecompressionLimit       uint64

	// The public peer-to-peer IP address of the host the enclave service is tied to
	// This is required to advertise for node discovery, and we include it in the attestation
	// todo - should we really bind the physical address to the attestation.
	HostAddress                   string
	AttestationEDBSecurityVersion uint
	AttestationSignerID           string
	AttestationEDBProductID       uint16
}

func EnclaveConfigFromTenConfig(tenCfg *config.TenConfig) *EnclaveConfig {
	limit, err := datasize.Parse(tenCfg.Enclave.DecompressionLimit)
	if err != nil {
		panic(fmt.Sprintf("failed to parse decompression limit: %v", err))
	}

	return &EnclaveConfig{
		NodeID:                    tenCfg.Node.ID,
		HostAddress:               tenCfg.Node.HostAddress,
		WillAttest:                tenCfg.Enclave.EnableAttestation,
		StoreExecutedTransactions: tenCfg.Enclave.StoreExecutedTransactions,
		DecompressionLimit:        uint64(limit),

		TenChainID: tenCfg.Network.ChainID,

		RPCAddress: tenCfg.Enclave.RPC.BindAddress,
		RPCTimeout: tenCfg.Enclave.RPC.Timeout,

		L1ChainID:                       tenCfg.Network.L1.ChainID,
		NetworkConfigAddress:            tenCfg.Network.L1.L1Contracts.NetworkConfigContract,
		DataAvailabilityRegistryAddress: tenCfg.Network.L1.L1Contracts.DataAvailabilityRegistry,
		EnclaveRegistryAddress:          tenCfg.Network.L1.L1Contracts.EnclaveRegistryContract,
		MessageBusAddress:               tenCfg.Network.L1.L1Contracts.MessageBusContract,
		L1BridgeAddress:                 tenCfg.Network.L1.L1Contracts.BridgeContract,
		SystemContractOwner:             tenCfg.Network.Sequencer.SystemContractsUpgrader,
		LogLevel:                        tenCfg.Enclave.Log.Level,
		LogPath:                         tenCfg.Enclave.Log.Path,

		UseInMemoryDB:  tenCfg.Enclave.DB.UseInMemory,
		EdgelessDBHost: tenCfg.Enclave.DB.EdgelessDBHost,
		SqliteDBPath:   tenCfg.Enclave.DB.SqlitePath,

		ProfilerEnabled:       tenCfg.Enclave.Debug.EnableProfiler,
		DebugNamespaceEnabled: tenCfg.Enclave.Debug.EnableDebugNamespace,

		MinGasPrice:              tenCfg.Network.Gas.MinGasPrice,
		GasPaymentAddress:        tenCfg.Network.Gas.PaymentAddress,
		MinBaseFee:               tenCfg.Network.Gas.MinBaseFee,
		GasBatchExecutionLimit:   tenCfg.Network.Gas.BatchExecutionLimit,
		GasLocalExecutionCapFlag: tenCfg.Network.Gas.LocalExecutionCap,

		TenGenesis:    tenCfg.Network.GenesisJSON,
		MaxBatchSize:  tenCfg.Network.Batch.MaxSize,
		MaxRollupSize: tenCfg.Network.Rollup.MaxSize,

		AttestationSignerID:           tenCfg.Enclave.Attestation.SignerID,
		AttestationEDBSecurityVersion: tenCfg.Enclave.Attestation.EDBSecurityVersion,
		AttestationEDBProductID:       tenCfg.Enclave.Attestation.EDBProductID,
	}
}
