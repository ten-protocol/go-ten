package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

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
	SequencerP2PAddress       string `yaml:"sequencerP2PAddress"`
	TenGenesis                string `yaml:"tenGenesis"`
	DebugNamespaceEnabled     bool   `yaml:"debugNamespaceEnabled"`
	MaxBatchSize              uint64 `yaml:"maxBatchSize"`
	MaxRollupSize             uint64 `yaml:"maxRollupSize"`
	GasPaymentAddress         string `yaml:"gasPaymentAddress"`
	BaseFee                   uint64 `yaml:"l2BaseFee"`
	GasBatchExecutionLimit    uint64 `yaml:"gasBatchExecutionLimit"`
	GasLocalExecutionCap      uint64 `yaml:"gasLocalExecutionCap"`
	RPCTimeout                int    `yaml:"rpcTimeout"`
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
		SequencerP2PAddress:    p.SequencerP2PAddress,
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

	enclaveConfig.RPCTimeout = time.Duration(p.RPCTimeout) * time.Second

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
	// P2P address for validators to connect to the sequencer for live batch data
	SequencerP2PAddress string
	// A json string that specifies the prefunded addresses at the genesis of the TEN network
	TenGenesis string
	// Whether debug calls are available
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

	// RPCTimeout - calls that are longer than this will be cancelled, to prevent resource starvation
	// normally, the context is propagated from the host, but in some cases ( like the evm, we have to create a context)
	RPCTimeout time.Duration
}

// SGX enclave.json configuration structs
type EnclaveConfigJson struct {
	Exe             string   `json:"exe"`
	Key             string   `json:"key"`
	Debug           bool     `json:"debug"`
	HeapSize        int      `json:"heapSize"`
	ExecutableHeap  bool     `json:"executableHeap"`
	ProductID       int      `json:"productID"`
	SecurityVersion int      `json:"securityVersion"`
	Mounts          []Mount  `json:"mounts"`
	Env             []EnvVar `json:"env"`
}

type Mount struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Type     string `json:"type"`
	ReadOnly bool   `json:"readOnly"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ToEnclaveConfigJson writes enclave.json from EnclaveInputConfig
func (p *EnclaveInputConfig) ToEnclaveConfigJson(filePath string) error {
	enclaveFile := filePath
	enclaveData, err := os.ReadFile(enclaveFile)
	if err != nil {
		fmt.Printf("Error reading enclave.json: %v\n", err)
		return err
	}

	var enclaveConfigJson EnclaveConfigJson
	err = json.Unmarshal(enclaveData, &enclaveConfigJson)
	if err != nil {
		fmt.Printf("Error unmarshalling enclave.json: %v\n", err)
		return err
	}

	// Inject variables into enclave.json using reflection
	configMap := createConfigMap(*p)

	for key, value := range configMap {
		if !isZeroValue(value) {
			enclaveConfigJson.Env = append(enclaveConfigJson.Env, EnvVar{Name: key, Value: fmt.Sprintf("%v", value)})
		}
	}
	// Write the updated enclave.json
	updatedData, err := json.MarshalIndent(enclaveConfigJson, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling updated enclave.json: %v\n", err)
		return err
	}

	err = os.WriteFile(enclaveFile, updatedData, 0644)
	if err != nil {
		fmt.Printf("Error writing updated enclave.json: %v\n", err)
		return err
	}

	fmt.Println("enclave.json updated successfully")

	return nil
}

// createConfigMap uses reflection to create a map from struct fields
func createConfigMap(config EnclaveInputConfig) map[string]interface{} {
	configMap := make(map[string]interface{})
	v := reflect.ValueOf(config)
	t := reflect.TypeOf(config)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if !isZeroValue(value) {
			configMap[strings.ToUpper(field.Name)] = value
		}
	}

	return configMap
}

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
