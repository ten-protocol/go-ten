package config

import (
	"encoding/base64"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"math/big"
)

// EnclaveInputConfig used for parsing a default config or partial override from yaml file
type EnclaveInputConfig struct {
	HostID                    string `yaml:"hostID"`
	HostAddress               string `yaml:"hostAddress"`
	Address                   string `yaml:"address"`
	NodeType                  string `yaml:"nodeType"`
	L1ChainID                 int64  `yaml:"l1ChainID"`
	TenChainID                int64  `yaml:"tenChainID"`
	WillAttest                bool   `yaml:"willAttest"`
	ValidateL1Blocks          bool   `yaml:"validateL1Blocks"`
	GenesisJSON               string `yaml:"genesisJSON"`
	ManagementContractAddress string `yaml:"managementContractAddress"`
	LogLevel                  int    `yaml:"logLevel"`
	LogPath                   string `yaml:"logPath"`
	UseInMemoryDB             bool   `yaml:"useInMemoryDB"`
	EdgelessDBHost            string `yaml:"edgelessDBHost"`
	SqliteDBPath              string `yaml:"sqliteDBPath"`
	ProfilerEnabled           bool   `yaml:"profilerEnabled"`
	MinGasPrice               uint64 `yaml:"minGasPrice"`
	MessageBusAddress         string `yaml:"messageBusAddress"`
	SequencerID               string `yaml:"sequencerID"`
	TenGenesis                string `yaml:"tenGenesis"`
	DebugNamespaceEnabled     bool   `yaml:"debugNamespaceEnabled"`
	MaxBatchSize              uint64 `yaml:"maxBatchSize"`
	MaxRollupSize             uint64 `yaml:"maxRollupSize"`
	GasPaymentAddress         string `yaml:"gasPaymentAddress"`
	BaseFee                   uint64 `yaml:"l2BaseFee"`
	GasBatchExecutionLimit    uint64 `yaml:"gasBatchExecutionLimit"`
	GasLocalExecutionCap      uint64 `yaml:"gasLocalExecutionCap"`
}

// ToEnclaveConfig Generates an EnclaveConfig from flags or yaml to one with proper typing
func (p *EnclaveInputConfig) ToEnclaveConfig() (*EnclaveConfig, error) {
	// calculated
	nodeType, err := common.ToNodeType(p.NodeType)
	if err != nil {
		return nil, fmt.Errorf("unrecognized node type %s: %w", p.NodeType, err)
	}

	enclaveConfig := &EnclaveConfig{
		HostAddress:            p.HostAddress,
		Address:                p.Address,
		NodeType:               nodeType,
		L1ChainID:              p.L1ChainID,
		TenChainID:             p.TenChainID,
		WillAttest:             p.WillAttest,
		ValidateL1Blocks:       p.ValidateL1Blocks,
		LogLevel:               p.LogLevel,
		LogPath:                p.LogPath,
		UseInMemoryDB:          p.UseInMemoryDB,
		EdgelessDBHost:         p.EdgelessDBHost,
		SqliteDBPath:           p.SqliteDBPath,
		ProfilerEnabled:        p.ProfilerEnabled,
		TenGenesis:             p.TenGenesis,
		DebugNamespaceEnabled:  p.DebugNamespaceEnabled,
		MaxBatchSize:           p.MaxBatchSize,
		MaxRollupSize:          p.MaxRollupSize,
		GasBatchExecutionLimit: p.GasBatchExecutionLimit,
		GasLocalExecutionCap:   p.GasLocalExecutionCap,
	}

	// byte unmarshall
	decodedData, err := base64.StdEncoding.DecodeString(p.GenesisJSON)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshall %s: %w", p.GenesisJSON, err)
	}
	enclaveConfig.GenesisJSON = decodedData

	enclaveConfig.HostID = gethcommon.HexToAddress(p.HostID)
	enclaveConfig.ManagementContractAddress = gethcommon.HexToAddress(p.ManagementContractAddress)
	enclaveConfig.MessageBusAddress = gethcommon.HexToAddress(p.MessageBusAddress)
	enclaveConfig.SequencerID = gethcommon.HexToAddress(p.SequencerID)
	enclaveConfig.GasPaymentAddress = gethcommon.HexToAddress(p.GasPaymentAddress)

	// protocol params or override
	enclaveConfig.MinGasPrice = big.NewInt(params.InitialBaseFee)
	if p.MinGasPrice > 0 {
		enclaveConfig.MinGasPrice = new(big.Int).SetUint64(p.MinGasPrice)
	}
	enclaveConfig.BaseFee = big.NewInt(params.InitialBaseFee)
	if p.BaseFee > 0 {
		enclaveConfig.BaseFee = new(big.Int).SetUint64(p.BaseFee)
	}

	return enclaveConfig, nil
}

// EnclaveConfig contains the full configuration for a Ten enclave service.
type EnclaveConfig struct {
	// HostID, the identity of the host the enclave service is tied to
	HostID gethcommon.Address
	// HostAddress, the public peer-to-peer IP address of the host the enclave service is tied to
	HostAddress string
	// Address, the address on which to serve requests
	Address string
	// NodeType, The type of the node.
	NodeType common.NodeType
	// L1ChainID, the ID of the L1 chain
	L1ChainID int64
	// TenChainID, the ID of the TEN chain
	TenChainID int64
	// WillAttest, whether to produce a verified attestation report
	WillAttest bool
	// ValidateL1Blocks, whether to validate incoming L1 blocks
	ValidateL1Blocks bool
	// GenesisJSON, when validating incoming blocks, the genesis config for the L1 chain
	GenesisJSON []byte
	// ManagementContractAddress, the management contract address on the L1 network
	ManagementContractAddress gethcommon.Address
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// LogPath, the path that the enclave's logs are written to
	LogPath string
	// UseInMemoryDB, whether the enclave should use in-memory or persistent storage
	UseInMemoryDB bool
	// EdgelessDBHost address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled)
	EdgelessDBHost string
	// SqliteDBPath, filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	//	if using InMemory DB or if attestation is enabled)
	SqliteDBPath string
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
	// MinGasPrice is the minimum gas price for mining a transaction
	MinGasPrice *big.Int
	// MessageBusAddress L1 Address
	MessageBusAddress gethcommon.Address
	// SequencerID, the identity of the sequencer for the network
	SequencerID gethcommon.Address
	// TenGenesis, a json string that specifies the prefunded addresses at the genesis of the TEN network
	TenGenesis string
	// DebugNamespaceEnabled, whether debug calls are available
	DebugNamespaceEnabled bool
	// MaxBatchSize, maximum bytes a batch can be uncompressed.
	MaxBatchSize uint64
	// MaxRollupSize - configured to be close to what the ethereum clients
	// have configured as the maximum size a transaction can have. Note that this isn't
	// a protocol limit, but a miner imposed limit, and it might be hard to find someone
	// to include a transaction if it goes above it
	MaxRollupSize uint64
	// GasPaymentAddress address for covering L1 transaction fees
	GasPaymentAddress gethcommon.Address
	// BaseFee initial base fee for EIP-1559 blocks
	BaseFee *big.Int

	// Due to hiding L1 costs in the gas quantity, the gas limit needs to be huge
	// Arbitrum with the same approach has gas limit of 1,125,899,906,842,624,
	// whilst the usage is small. Should be ok since execution is paid for anyway.

	// GasBatchExecutionLimit maximum amount of gas that can be consumed by a single batch of transactions
	GasBatchExecutionLimit uint64
	// GasLocalExecutionCapFlag default is same value as `GasBatchExecutionLimit`
	GasLocalExecutionCap uint64
}
