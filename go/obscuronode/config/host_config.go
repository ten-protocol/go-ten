package config

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// HostConfig contains the full configuration for an Obscuro host.
type HostConfig struct {
	// The host's identity
	ID common.Address
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// Duration of the gossip round
	GossipRoundDuration time.Duration
	// Whether to serve client RPC requests
	HasClientRPC bool
	// Address on which to serve client RPC requests
	ClientRPCAddress string
	// Timeout duration for RPC requests from client applications
	ClientRPCTimeout time.Duration
	// Address on which to connect to the enclave
	EnclaveRPCAddress string
	// Timeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// Our network for P2P communication with peer Obscuro nodes
	P2PAddress string
	// The addresses of all the Obscuro nodes on the network
	AllP2PAddresses []string
	// The host of the connected L1 node
	L1NodeHost string
	// The websocket port of the connected L1 node
	L1NodeWebsocketPort uint
	// The rollup contract address on the L1 network
	RollupContractAddress common.Address
	// The path that the node's logs are written to
	LogPath string
	// The stringified private key for the host's L1 wallet
	PrivateKeyString string
	// The ID of the L1 chain
	ChainID big.Int
}
