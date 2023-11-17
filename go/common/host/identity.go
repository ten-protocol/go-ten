package host

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
)

type Identity struct {
	ID               gethcommon.Address
	P2PPublicAddress string
	IsGenesis        bool
	IsSequencer      bool
}

func NewIdentity(cfg *config.HostConfig) Identity {
	return Identity{
		ID:               cfg.ID,
		P2PPublicAddress: cfg.P2PPublicAddress,
		IsGenesis:        cfg.IsGenesis,
		IsSequencer:      cfg.NodeType == common.Sequencer,
	}
}
