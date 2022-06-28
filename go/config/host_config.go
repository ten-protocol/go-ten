package config

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	defaultRPCTimeoutSecs          = 10
	defaultL1ConnectionTimeoutSecs = 15
)

// HostConfig contains the full configuration for an Obscuro host.
type HostConfig struct {
	// The host's identity
	ID common.Address
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// Duration of the gossip round
	GossipRoundDuration time.Duration
	// Whether to serve client RPC requests over HTTP
	HasClientRPCHTTP bool
	// Port on which to handle HTTP client RPC requests
	ClientRPCPortHTTP uint64
	// Whether to serve client RPC requests over websockets
	HasClientRPCWebsockets bool
	// Port on which to handle websocket client RPC requests
	ClientRPCPortWS uint64
	// Host on which to handle client RPC requests
	ClientRPCHost string
	// Timeout duration for RPC requests from client applications
	ClientRPCTimeout time.Duration
	// Address on which to connect to the enclave
	EnclaveRPCAddress string
	// Timeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// Our network for P2P communication with peer Obscuro nodes
	P2PAddress string
	// The host of the connected L1 node
	L1NodeHost string
	// The websocket port of the connected L1 node
	L1NodeWebsocketPort uint
	// Timeout duration for connecting to the L1 node
	L1ConnectionTimeout time.Duration
	// The rollup contract address on the L1 network
	RollupContractAddress common.Address
	// The path that the node's logs are written to
	LogPath string
	// The stringified private key for the host's L1 wallet
	PrivateKeyString string
	// The ID of the L1 chain
	ChainID big.Int
}

// DefaultHostConfig returns a HostConfig with default values.
func DefaultHostConfig() HostConfig {
	return HostConfig{
		ID:                     common.BytesToAddress([]byte("")),
		IsGenesis:              true,
		GossipRoundDuration:    8333,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      13000,
		HasClientRPCWebsockets: false,
		ClientRPCPortWS:        13001,
		ClientRPCHost:          "127.0.0.1",
		ClientRPCTimeout:       time.Duration(defaultRPCTimeoutSecs) * time.Second,
		EnclaveRPCAddress:      "127.0.0.1:11000",
		EnclaveRPCTimeout:      time.Duration(defaultRPCTimeoutSecs) * time.Second,
		P2PAddress:             "127.0.0.1:10000",
		L1NodeHost:             "127.0.0.1",
		L1NodeWebsocketPort:    8546,
		L1ConnectionTimeout:    time.Duration(defaultL1ConnectionTimeoutSecs) * time.Second,
		RollupContractAddress:  common.BytesToAddress([]byte("")),
		LogPath:                "host_logs.txt",
		PrivateKeyString:       "0000000000000000000000000000000000000000000000000000000000000001",
		ChainID:                *big.NewInt(1337),
	}
}
