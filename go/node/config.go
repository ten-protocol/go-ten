package node

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/integration/common/testlog"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	_localhost = "127.0.0.1"
)

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config represents the configurations default loaded from a YAML file overridable by CLI
type Config struct {
	NodeAction              string `yaml:"nodeAction"`
	NodeType                string `yaml:"nodeType"`
	IsGenesis               bool   `yaml:"isGenesis"`
	IsSGXEnabled            bool   `yaml:"isSGXEnabled"`
	EnclaveDockerImage      string `yaml:"enclaveDockerImage"`
	HostDockerImage         string `yaml:"hostDockerImage"`
	L1WebsocketURL          string `yaml:"l1WebsocketURL"`
	HostP2PPort             int    `yaml:"hostP2PPort"`
	HostP2PPublicAddr       string `yaml:"hostP2PPublicAddr"`
	EnclaveWSPort           int    `yaml:"enclaveWSPort"`
	PrivateKey              string `yaml:"privateKey"`
	HostID                  string `yaml:"hostID"`
	SequencerID             string `yaml:"sequencerID"`
	ManagementContractAddr  string `yaml:"managementContractAddr"`
	MessageBusContractAddr  string `yaml:"messageBusContractAddr"`
	L1Start                 string `yaml:"l1Start"`
	PccsAddr                string `yaml:"pccsAddr"`
	EdgelessDBImage         string `yaml:"edgelessDBImage"`
	HostHTTPPort            int    `yaml:"hostHTTPPort"`
	HostWSPort              int    `yaml:"hostWSPort"`
	NodeName                string `yaml:"nodeName"`
	IsDebugNamespaceEnabled bool   `yaml:"isDebugNamespaceEnabled"`
	LogLevel                int    `yaml:"logLevel"`
	IsInboundP2PDisabled    bool   `yaml:"isInboundP2PDisabled"`
	BatchInterval           string `yaml:"batchInterval"`
	MaxBatchInterval        string `yaml:"maxBatchInterval"`
	RollupInterval          string `yaml:"rollupInterval"`
	L1ChainID               int    `yaml:"l1ChainID"`
	ProfilerEnabled         bool   `yaml:"profilerEnabled"`
	MetricsEnabled          bool   `yaml:"metricsEnabled"`
	CoinbaseAddress         string `yaml:"coinbaseAddress"`
	L1BlockTime             int    `yaml:"l1BlockTime"`
	TenGenesis              string `yaml:"tenGenesis"`
	EnclaveDebug            bool   `yaml:"enclaveDebug"`
	HostInMemDB             bool   `yaml:"hostInMemDB"`
	HostExternalDBHost      string `yaml:"hostExternalDBHost"`
	HostExternalDBUser      string `yaml:"hostExternalDBUser"`
	HostExternalDBPass      string `yaml:"hostExternalDBPass"`
}

// LoadConfig reads configuration from a file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	defaultConfig := &Config{}

	// Read YAML configuration
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

// NewNodeConfig takes and initial default config and sets options on
func NewNodeConfig(config Config, opts ...Option) *Config {
	initialConfig := &config
	for _, opt := range opts {
		opt(initialConfig)
	}

	return initialConfig
}

func (c *Config) ApplyOverrides(o *Config) {
	// Obtain reflect.Value objects for both structs.
	cVal := reflect.ValueOf(c).Elem()
	oVal := reflect.ValueOf(o).Elem()

	// Iterate over each field in the override struct.
	for i := 0; i < oVal.NumField(); i++ {
		oField := oVal.Field(i)
		cField := cVal.Field(i)

		// Check if the field in the override struct is set (non-default).
		if isFieldSet(oField) {
			cField.Set(oField)
		}
	}
}

// isFieldSet determines whether the provided reflect.Value holds a non-default value.
func isFieldSet(field reflect.Value) bool {
	// Handle based on the field kind.
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() != 0
	case reflect.String:
		return field.String() != ""
	case reflect.Bool:
		return field.Bool()
	default:
		panic("unhandled default case")
	}

	// For struct or other complex types, you might need a more sophisticated approach.
	return false
}

func (c *Config) ToEnclaveConfig() *config.EnclaveConfig {
	cfg := config.DefaultEnclaveConfig()

	if c.NodeType == "validator" {
		cfg.NodeType = common.Validator
	}

	cfg.MessageBusAddress = gethcommon.HexToAddress(c.MessageBusContractAddr)
	cfg.ManagementContractAddress = gethcommon.HexToAddress(c.ManagementContractAddr)
	cfg.SequencerID = gethcommon.HexToAddress(c.SequencerID)
	cfg.HostID = gethcommon.HexToAddress(c.HostID)
	cfg.HostAddress = fmt.Sprintf("127.0.0.1:%d", c.HostP2PPort)
	cfg.LogPath = testlog.LogFile()
	cfg.LogLevel = c.LogLevel
	cfg.Address = fmt.Sprintf("%s:%d", _localhost, c.EnclaveWSPort)
	cfg.DebugNamespaceEnabled = c.IsDebugNamespaceEnabled
	cfg.TenGenesis = c.TenGenesis

	if c.NodeType == "sequencer" && c.CoinbaseAddress != "" {
		cfg.GasPaymentAddress = gethcommon.HexToAddress(c.CoinbaseAddress)
	}

	return cfg
}

func (c *Config) ToHostConfig() *config.HostConfig {
	cfg := config.DefaultHostConfig()

	if c.NodeType == "validator" {
		cfg.NodeType = common.Validator
	}

	cfg.IsGenesis = c.IsGenesis
	cfg.PrivateKeyString = c.PrivateKey
	cfg.EnclaveRPCAddress = fmt.Sprintf("127.0.0.1:%d", c.EnclaveWSPort)
	cfg.ClientRPCPortWS = uint64(c.HostWSPort)
	cfg.ClientRPCPortHTTP = uint64(c.HostHTTPPort)

	cfg.P2PPublicAddress = fmt.Sprintf("127.0.0.1:%d", c.HostP2PPort)
	cfg.P2PBindAddress = c.HostP2PPublicAddr

	cfg.L1WebsocketURL = c.L1WebsocketURL
	cfg.ManagementContractAddress = gethcommon.HexToAddress(c.ManagementContractAddr)
	cfg.MessageBusAddress = gethcommon.HexToAddress(c.MessageBusContractAddr)
	cfg.LogPath = testlog.LogFile()
	cfg.ProfilerEnabled = c.ProfilerEnabled
	cfg.MetricsEnabled = c.MetricsEnabled
	cfg.DebugNamespaceEnabled = c.IsDebugNamespaceEnabled
	cfg.LogLevel = c.LogLevel
	cfg.SequencerID = gethcommon.HexToAddress(c.SequencerID)
	cfg.IsInboundP2PDisabled = c.IsInboundP2PDisabled
	cfg.L1BlockTime = time.Second * time.Duration(c.L1BlockTime)
	cfg.L1ChainID = int64(c.L1ChainID)

	return cfg
}

func (c *Config) UpdateNodeConfig(opts ...Option) *Config {
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithNodeType(nodeType string) Option {
	return func(c *Config) {
		c.NodeType = nodeType
	}
}

func WithGenesis(b bool) Option {
	return func(c *Config) {
		c.IsGenesis = b
	}
}

func WithSGXEnabled(b bool) Option {
	return func(c *Config) {
		c.IsSGXEnabled = b
	}
}

func WithEnclaveImage(s string) Option {
	return func(c *Config) {
		c.EnclaveDockerImage = s
	}
}

func WithHostImage(s string) Option {
	return func(c *Config) {
		c.HostDockerImage = s
	}
}

func WithL1WebsocketURL(addr string) Option {
	return func(c *Config) {
		c.L1WebsocketURL = addr
	}
}

func WithHostP2PPort(i int) Option {
	return func(c *Config) {
		c.HostP2PPort = i
	}
}

func WithHostPublicP2PAddr(s string) Option {
	return func(c *Config) {
		c.HostP2PPublicAddr = s
	}
}

func WithEnclaveWSPort(i int) Option {
	return func(c *Config) {
		c.EnclaveWSPort = i
	}
}

func WithPrivateKey(s string) Option {
	return func(c *Config) {
		c.PrivateKey = s
	}
}

func WithHostID(s string) Option {
	return func(c *Config) {
		c.HostID = s
	}
}

func WithSequencerID(s string) Option {
	return func(c *Config) {
		c.SequencerID = s
	}
}

func WithManagementContractAddress(s string) Option {
	return func(c *Config) {
		c.ManagementContractAddr = s
	}
}

func WithMessageBusContractAddress(s string) Option {
	return func(c *Config) {
		c.MessageBusContractAddr = s
	}
}

func WithL1Start(blockHash string) Option {
	return func(c *Config) {
		c.L1Start = blockHash
	}
}

func WithPCCSAddr(s string) Option {
	return func(c *Config) {
		c.PccsAddr = s
	}
}

func WithEdgelessDBImage(s string) Option {
	return func(c *Config) {
		c.EdgelessDBImage = s
	}
}

func WithHostHTTPPort(i int) Option {
	return func(c *Config) {
		c.HostHTTPPort = i
	}
}

func WithHostWSPort(i int) Option {
	return func(c *Config) {
		c.HostWSPort = i
	}
}

func WithNodeName(s string) Option {
	return func(c *Config) {
		c.NodeName = s
	}
}

func WithDebugNamespaceEnabled(b bool) Option {
	return func(c *Config) {
		c.IsDebugNamespaceEnabled = b
	}
}

func WithLogLevel(i int) Option {
	return func(c *Config) {
		c.LogLevel = i
	}
}

func WithInboundP2PDisabled(b bool) Option {
	return func(c *Config) {
		c.IsInboundP2PDisabled = b
	}
}

func WithBatchInterval(d string) Option {
	return func(c *Config) {
		c.BatchInterval = d
	}
}

func WithMaxBatchInterval(d string) Option {
	return func(c *Config) {
		c.MaxBatchInterval = d
	}
}

func WithRollupInterval(d string) Option {
	return func(c *Config) {
		c.RollupInterval = d
	}
}

func WithL1ChainID(i int) Option {
	return func(c *Config) {
		c.L1ChainID = i
	}
}

func WithProfiler(b bool) Option {
	return func(c *Config) {
		c.ProfilerEnabled = b
	}
}

func WithMetricsEnabled(b bool) Option {
	return func(c *Config) {
		c.MetricsEnabled = b
	}
}

func WithCoinbase(s string) Option {
	return func(c *Config) {
		c.CoinbaseAddress = s
	}
}

func WithL1BlockTime(d int) Option {
	return func(c *Config) {
		c.L1BlockTime = d
	}
}

func WithTenGenesis(s string) Option {
	return func(c *Config) {
		c.TenGenesis = s
	}
}

func WithEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.EnclaveDebug = b
	}
}

func WithInMemoryHostDB(b bool) Option {
	return func(c *Config) {
		c.HostInMemDB = b
	}
}

func WithHostExternalDBHost(s string) Option {
	return func(c *Config) {
		c.HostExternalDBHost = s
	}
}

func WithHostExternalDBUser(s string) Option {
	return func(c *Config) {
		c.HostExternalDBUser = s
	}
}

func WithHostExternalDBPass(s string) Option {
	return func(c *Config) {
		c.HostExternalDBPass = s
	}
}
