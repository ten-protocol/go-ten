package node

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	_localhost = "127.0.0.1"
)

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	isGenesis                 bool
	sgxEnabled                bool
	enclaveImage              string
	hostImage                 string
	nodeType                  string
	l1WSURL                   string
	sequencerID               string
	privateKey                string
	hostP2PPort               int
	hostID                    string
	hostHTTPPort              int
	hostWSPort                int
	enclaveWSPort             int
	messageBusContractAddress string
	managementContractAddr    string
	l1Start                   string
	hostP2PHost               string
	hostPublicP2PAddr         string
	pccsAddr                  string
	edgelessDBImage           string
	enclaveDebug              bool
	nodeName                  string
	hostInMemDB               bool
	debugNamespaceEnabled     bool
	profilerEnabled           bool
	coinbaseAddress           string
	logLevel                  int
	isInboundP2PDisabled      bool
	l1BlockTime               time.Duration
	batchInterval             string
	rollupInterval            string
}

func NewNodeConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func (c *Config) ToEnclaveConfig() *config.EnclaveConfig {
	cfg := config.DefaultEnclaveConfig()

	if c.nodeType == "validator" {
		cfg.NodeType = common.Validator
	}

	cfg.MessageBusAddress = gethcommon.HexToAddress(c.messageBusContractAddress)
	cfg.ManagementContractAddress = gethcommon.HexToAddress(c.managementContractAddr)
	cfg.SequencerID = gethcommon.HexToAddress(c.sequencerID)
	cfg.HostID = gethcommon.HexToAddress(c.hostID)
	cfg.HostAddress = fmt.Sprintf("127.0.0.1:%d", c.hostP2PPort)
	cfg.LogPath = testlog.LogFile()
	cfg.LogLevel = c.logLevel
	cfg.Address = fmt.Sprintf("%s:%d", _localhost, c.enclaveWSPort)
	cfg.DebugNamespaceEnabled = c.debugNamespaceEnabled

	if c.nodeType == "sequencer" && c.coinbaseAddress != "" {
		cfg.GasPaymentAddress = gethcommon.HexToAddress(c.coinbaseAddress)
	}

	return cfg
}

func (c *Config) ToHostConfig() *config.HostInputConfig {
	cfg := config.DefaultHostParsedConfig()

	if c.nodeType == "validator" {
		cfg.NodeType = common.Validator
	}

	cfg.IsGenesis = c.isGenesis
	cfg.PrivateKeyString = c.privateKey
	cfg.EnclaveRPCAddress = fmt.Sprintf("127.0.0.1:%d", c.enclaveWSPort)
	cfg.ClientRPCPortWS = uint64(c.hostWSPort)
	cfg.ClientRPCPortHTTP = uint64(c.hostHTTPPort)

	cfg.P2PPublicAddress = fmt.Sprintf("127.0.0.1:%d", c.hostP2PPort)
	cfg.P2PBindAddress = c.hostPublicP2PAddr

	cfg.L1WebsocketURL = c.l1WSURL
	cfg.ManagementContractAddress = gethcommon.HexToAddress(c.managementContractAddr)
	cfg.MessageBusAddress = gethcommon.HexToAddress(c.messageBusContractAddress)
	cfg.LogPath = testlog.LogFile()
	cfg.ProfilerEnabled = c.profilerEnabled
	cfg.MetricsEnabled = false
	cfg.DebugNamespaceEnabled = c.debugNamespaceEnabled
	cfg.LogLevel = c.logLevel
	cfg.SequencerID = gethcommon.HexToAddress(c.sequencerID)
	cfg.IsInboundP2PDisabled = c.isInboundP2PDisabled
	cfg.L1BlockTime = c.l1BlockTime

	return cfg
}

func (c *Config) UpdateNodeConfig(opts ...Option) *Config {
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithCoinbase(s string) Option {
	return func(c *Config) {
		c.coinbaseAddress = s
	}
}

func WithNodeName(s string) Option {
	return func(c *Config) {
		c.nodeName = s
	}
}

func WithNodeType(nodeType string) Option {
	return func(c *Config) {
		c.nodeType = nodeType
	}
}

func WithGenesis(b bool) Option {
	return func(c *Config) {
		c.isGenesis = b
	}
}

func WithSGXEnabled(b bool) Option {
	return func(c *Config) {
		c.sgxEnabled = b
	}
}

func WithEnclaveImage(s string) Option {
	return func(c *Config) {
		c.enclaveImage = s
	}
}

func WithEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.enclaveDebug = b
	}
}

func WithHostImage(s string) Option {
	return func(c *Config) {
		c.hostImage = s
	}
}

func WithMessageBusContractAddress(s string) Option {
	return func(c *Config) {
		c.messageBusContractAddress = s
	}
}

func WithManagementContractAddress(s string) Option {
	return func(c *Config) {
		c.managementContractAddr = s
	}
}

func WithSequencerID(s string) Option {
	return func(c *Config) {
		c.sequencerID = s
	}
}

func WithHostID(s string) Option {
	return func(c *Config) {
		c.hostID = s
	}
}

func WithPrivateKey(s string) Option {
	return func(c *Config) {
		c.privateKey = s
	}
}

func WithEnclaveWSPort(i int) Option {
	return func(c *Config) {
		c.enclaveWSPort = i
	}
}

func WithL1Start(blockHash string) Option {
	return func(c *Config) {
		c.l1Start = blockHash
	}
}

func WithL1WebsocketURL(addr string) Option {
	return func(c *Config) {
		c.l1WSURL = addr
	}
}

func WithHostP2PPort(i int) Option {
	return func(c *Config) {
		c.hostP2PPort = i
	}
}

func WithHostP2PHost(s string) Option {
	return func(c *Config) {
		c.hostP2PHost = s
	}
}

func WithHostPublicP2PAddr(s string) Option {
	return func(c *Config) {
		c.hostPublicP2PAddr = s
	}
}

func WithHostHTTPPort(i int) Option {
	return func(c *Config) {
		c.hostHTTPPort = i
	}
}

func WithHostWSPort(i int) Option {
	return func(c *Config) {
		c.hostWSPort = i
	}
}

func WithEdgelessDBImage(s string) Option {
	return func(c *Config) {
		c.edgelessDBImage = s
	}
}

func WithPCCSAddr(s string) Option {
	return func(c *Config) {
		c.pccsAddr = s
	}
}

func WithInMemoryHostDB(b bool) Option {
	return func(c *Config) {
		c.hostInMemDB = b
	}
}

func WithDebugNamespaceEnabled(b bool) Option {
	return func(c *Config) {
		c.debugNamespaceEnabled = b
	}
}

func WithProfiler(b bool) Option {
	return func(c *Config) {
		c.profilerEnabled = b
	}
}

func WithLogLevel(i int) Option {
	return func(c *Config) {
		c.logLevel = i
	}
}

func WithInboundP2PDisabled(b bool) Option {
	return func(c *Config) {
		c.isInboundP2PDisabled = b
	}
}

func WithL1BlockTime(d time.Duration) Option {
	return func(c *Config) {
		c.l1BlockTime = d
	}
}

func WithBatchInterval(d string) Option {
	return func(c *Config) {
		c.batchInterval = d
	}
}

func WithRollupInterval(d string) Option {
	return func(c *Config) {
		c.rollupInterval = d
	}
}
