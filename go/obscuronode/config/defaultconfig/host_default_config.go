package defaultconfig

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"math/big"
	"time"
)

const (
	defaultRPCTimeoutSecs          = 10
	defaultL1ConnectionTimeoutSecs = 15
)

// DefaultHostConfig returns a HostConfig with default values.
func DefaultHostConfig() config.HostConfig {
	return config.HostConfig{
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
