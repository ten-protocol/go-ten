package hostrunner

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/naoina/toml"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/common"
)

// HostConfigToml is the structure that a host's .toml config is parsed into.
type HostConfigToml struct {
	ID                      string
	IsGenesis               bool
	GossipRoundNanos        int
	ClientRPCPortHTTP       uint
	ClientRPCPortWS         uint
	ClientRPCHost           string
	ClientRPCTimeoutSecs    int
	EnclaveRPCAddress       string
	EnclaveRPCTimeoutSecs   int
	P2PAddress              string
	L1NodeHost              string
	L1NodePort              uint
	L1ConnectionTimeoutSecs int
	RollupContractAddress   string
	LogPath                 string
	PrivateKey              string
	ChainID                 big.Int
}

// ParseConfig returns a config.HostConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.HostConfig {
	defaultConfig := config.DefaultHostConfig()

	configPath := flag.String(configName, "", configUsage)
	nodeID := flag.String(nodeIDName, defaultConfig.ID.Hex(), nodeIDUsage)
	isGenesis := flag.Bool(isGenesisName, defaultConfig.IsGenesis, isGenesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(defaultConfig.GossipRoundDuration), gossipRoundNanosUsage)
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPName, defaultConfig.ClientRPCPortHTTP, clientRPCPortHTTPUsage)
	clientRPCPortWS := flag.Uint64(clientRPCPortWSName, defaultConfig.ClientRPCPortWS, clientRPCPortWSUsage)
	clientRPCHost := flag.String(clientRPCHostName, defaultConfig.ClientRPCHost, clientRPCHostUsage)
	clientRPCTimeoutSecs := flag.Uint64(clientRPCTimeoutSecsName, uint64(defaultConfig.ClientRPCTimeout.Seconds()), clientRPCTimeoutSecsUsage)
	enclaveRPCAddress := flag.String(enclaveRPCAddressName, defaultConfig.EnclaveRPCAddress, enclaveRPCAddressUsage)
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, uint64(defaultConfig.EnclaveRPCTimeout.Seconds()), enclaveRPCTimeoutSecsUsage)
	p2pAddress := flag.String(p2pAddressName, defaultConfig.P2PAddress, p2pAddressUsage)
	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.L1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(defaultConfig.L1NodeWebsocketPort), l1NodePortUsage)
	l1ConnectionTimeoutSecs := flag.Uint64(l1ConnectionTimeoutSecsName, uint64(defaultConfig.L1ConnectionTimeout.Seconds()), l1ConnectionTimeoutSecsUsage)
	rollupContractAddress := flag.String(rollupContractAddrName, defaultConfig.RollupContractAddress.Hex(), rollupContractAddrUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID.Int64(), chainIDUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKeyString, privateKeyUsage)

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	defaultConfig.ID = common.HexToAddress(*nodeID)
	defaultConfig.IsGenesis = *isGenesis
	defaultConfig.GossipRoundDuration = time.Duration(*gossipRoundNanos)
	defaultConfig.HasClientRPCHTTP = true
	defaultConfig.ClientRPCPortHTTP = *clientRPCPortHTTP
	defaultConfig.HasClientRPCWebsockets = true
	defaultConfig.ClientRPCPortWS = *clientRPCPortWS
	defaultConfig.ClientRPCHost = *clientRPCHost
	defaultConfig.ClientRPCTimeout = time.Duration(*enclaveRPCTimeoutSecs) * time.Second
	defaultConfig.EnclaveRPCAddress = *enclaveRPCAddress
	defaultConfig.EnclaveRPCTimeout = time.Duration(*clientRPCTimeoutSecs) * time.Second
	defaultConfig.P2PAddress = *p2pAddress
	defaultConfig.L1NodeHost = *l1NodeHost
	defaultConfig.L1NodeWebsocketPort = uint(*l1NodePort)
	defaultConfig.L1ConnectionTimeout = time.Duration(*l1ConnectionTimeoutSecs) * time.Second
	defaultConfig.RollupContractAddress = common.HexToAddress(*rollupContractAddress)
	defaultConfig.PrivateKeyString = *privateKeyStr
	defaultConfig.LogPath = *logPath
	defaultConfig.ChainID = *big.NewInt(*chainID)

	return defaultConfig
}

// Parses the config from the .toml file at configPath.
func fileBasedConfig(configPath string) config.HostConfig {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var tomlConfig HostConfigToml
	err = toml.Unmarshal(bytes, &tomlConfig)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	return config.HostConfig{
		ID:                     common.HexToAddress(tomlConfig.ID),
		IsGenesis:              tomlConfig.IsGenesis,
		GossipRoundDuration:    time.Duration(tomlConfig.GossipRoundNanos),
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      uint64(tomlConfig.ClientRPCPortHTTP),
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        uint64(tomlConfig.ClientRPCPortWS),
		ClientRPCHost:          tomlConfig.ClientRPCHost,
		ClientRPCTimeout:       time.Duration(tomlConfig.ClientRPCTimeoutSecs) * time.Second,
		EnclaveRPCAddress:      tomlConfig.EnclaveRPCAddress,
		EnclaveRPCTimeout:      time.Duration(tomlConfig.EnclaveRPCTimeoutSecs) * time.Second,
		P2PAddress:             tomlConfig.P2PAddress,
		L1NodeHost:             tomlConfig.L1NodeHost,
		L1NodeWebsocketPort:    tomlConfig.L1NodePort,
		L1ConnectionTimeout:    time.Duration(tomlConfig.L1ConnectionTimeoutSecs) * time.Second,
		RollupContractAddress:  common.HexToAddress(tomlConfig.RollupContractAddress),
		LogPath:                tomlConfig.LogPath,
		PrivateKeyString:       tomlConfig.PrivateKey,
		ChainID:                tomlConfig.ChainID,
	}
}
