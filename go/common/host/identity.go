package host

import gethcommon "github.com/ethereum/go-ethereum/common"

type Identity struct {
	ID               gethcommon.Address
	P2PPublicAddress string
}

func NewIdentity(id gethcommon.Address, p2pPublicAddress string) Identity {
	return Identity{
		ID:               id,
		P2PPublicAddress: p2pPublicAddress,
	}
}
