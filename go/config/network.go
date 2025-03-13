package config

import (
	"github.com/ten-protocol/go-ten/go/common"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// NetworkConfig contains the static configuration for this instance of the Ten network
//
//	yaml: `network`
type NetworkConfig struct {
	// ChainID is the chainID for the Ten network
	ChainID int64 `mapstructure:"chainId"`
	// GenesisJSON is a json string that specifies the prefunded addresses at the genesis of the Ten network
	GenesisJSON string `mapstructure:"genesis"`

	Batch      *BatchConfig      `mapstructure:"batch"`
	Gas        *GasConfig        `mapstructure:"gas"`
	L1         *L1Config         `mapstructure:"l1"`
	Rollup     *RollupConfig     `mapstructure:"rollup"`
	Sequencer  *Sequencer        `mapstructure:"sequencer"`
	CrossChain *CrossChainConfig `mapstructure:"crossChain"`
}

// BatchConfig contains the configuration for the batch processing on the Ten network
//
//	yaml: `network.batch`
type BatchConfig struct {
	// Interval is the time between batches being created
	Interval time.Duration `mapstructure:"interval"`
	// MaxInterval is the maximum time between batches being created (if this is set higher than Batch Interval, the host will
	// not create empty batches until the MaxBatchInterval is reached or a transaction is received)
	MaxInterval time.Duration `mapstructure:"maxInterval"`
	// MaxSize is the maximum bytes a batch can be uncompressed
	MaxSize uint64 `mapstructure:"maxSize"`
}

// GasConfig contains the gas configuration for the Ten network
//
//	yaml: `network.gas`
type GasConfig struct {
	BaseFee *big.Int `mapstructure:"baseFee"`
	// MinGasPrice is the minimum gas price for mining a transaction
	MinGasPrice         *big.Int           `mapstructure:"minGasPrice"`
	PaymentAddress      gethcommon.Address `mapstructure:"paymentAddress"`
	BatchExecutionLimit uint64             `mapstructure:"batchExecutionLimit"`
	LocalExecutionCap   uint64             `mapstructure:"localExecutionCap"`
}

// L1Config contains config about the L1 network that the Ten network is rolling up to
//
//	yaml: `network.l1`
type L1Config struct {
	ChainID   int64           `mapstructure:"chainId"`   // chainID for the L1 network
	BlockTime time.Duration   `mapstructure:"blockTime"` // average expected block time for the L1 network
	StartHash gethcommon.Hash `mapstructure:"startHash"` // hash of the first block on the L1 network relevant to the Ten network

	L1Contracts *L1Contracts `mapstructure:"contracts"`
}

// L1Contracts contains the addresses of Ten contracts on the L1 network
//
//	yaml: `network.l1.contracts`
type L1Contracts struct {
	//FIXME add to yamls
	NetworkConfigContract   common.NetworkConfigAddress `mapstructure:"networkConfig"` //this might be the only one we need
	CrossChainContract      common.CrossChainAddress    `mapstructure:"crossChain"`
	RollupContract          common.RollupAddress        `mapstructure:"rollup"`
	EnclaveRegistryContract common.RollupAddress        `mapstructure:"enclaveRegistry"`
	MessageBusContract      gethcommon.Address          `mapstructure:"messageBus"`
	BridgeContract          gethcommon.Address          `mapstructure:"bridge"`
}

// RollupConfig contains the configuration for the rollup processing on the Ten network
//
//	yaml: `network.rollup`
type RollupConfig struct {
	// Interval is the time between sequencer checking if it should produce rollups
	Interval time.Duration `mapstructure:"interval"`
	// MaxInterval is the max time between sequencer rollups (it will create even if the rollup won't be full)
	MaxInterval time.Duration `mapstructure:"maxInterval"`
	// MaxSize - (in bytes) configured to be close to what the ethereum clients
	// have configured as the maximum size a transaction can have. Note that this isn't
	// a protocol limit, but a miner imposed limit and it might be hard to find someone
	// to include a transaction if it goes above it
	MaxSize uint64 `mapstructure:"maxSize"`
}

// Sequencer contains the configuration for how the L2 sequencer will operate for the Ten network
//
//	yaml: `network.sequencer`
type Sequencer struct {
	// P2PAddress is the address that the sequencer will listen on for incoming P2P connections
	P2PAddress              string             `mapstructure:"p2pAddress"`
	SystemContractsUpgrader gethcommon.Address `mapstructure:"systemContractsUpgrader"`
}

// CrossChainConfig contains the configuration for the cross chain processing on the Ten network
//
//	yaml: `network.crossChain`
type CrossChainConfig struct {
	// Interval is the time between sequencer checking if it should produce cross chain messages
	Interval time.Duration `mapstructure:"interval"`
}
