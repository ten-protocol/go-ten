package core

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventVisibilityConfig struct {
	AutoConfig                                  bool
	Public                                      bool  // everyone can see and query for this event
	Topic1CanView, Topic2CanView, Topic3CanView *bool // If the event is private, and this is true, it means that the address from topic1 is an EOA that can view this event
	SenderCanView                               *bool // if true, the tx signer will see this event. Default false
}

// ContractVisibilityConfig represents the configuration as defined by the dApp developer in the smart contract
type ContractVisibilityConfig struct {
	AutoConfig   bool // true - if the platform has to autodetect
	Transparent  *bool
	EventConfigs map[gethcommon.Hash]*EventVisibilityConfig
}

type TxExecResult struct {
	Receipt          *types.Receipt
	CreatedContracts map[gethcommon.Address]*ContractVisibilityConfig
	Err              error
}
