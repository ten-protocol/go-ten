package host

import (
	"github.com/ten-protocol/go-ten/go/common"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
)

type Identity struct {
	P2PPublicAddress string
	IsGenesis        bool
	IsSequencer      bool
}

func NewIdentity(cfg *hostconfig.HostConfig) Identity {
	return Identity{
		P2PPublicAddress: cfg.P2PPublicAddress,
		IsGenesis:        cfg.IsGenesis,
		IsSequencer:      cfg.NodeType == common.Sequencer,
	}
}
