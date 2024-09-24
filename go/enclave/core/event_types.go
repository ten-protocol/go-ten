package core

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventVisibilityConfig - configuration per event by the dApp developer(DD)
// There are 4 cases:
// 1. DD doesn't configure anything. - ContractVisibilityConfig.AutoConfig=true
// 2. DD configures and  specifies the contract as transparent - ContractVisibilityConfig.Transparent=true
// 3. DD configures and specify the contract as non-transparent, but doesn't configure the event - Contract: false/false , EventVisibilityConfig.AutoConfig=true
// DD configures the contract as non-transparent, and also configures the topics for the event
type EventVisibilityConfig struct {
	AutoConfig                                  bool  // true for events that have no explicit configuration
	Public                                      bool  // everyone can see and query for this event
	Topic1CanView, Topic2CanView, Topic3CanView *bool // If the event is not public, and this is true, it means that the address from topicI is an EOA that can view this event
	SenderCanView                               *bool // if true, the tx signer will see this event. Default false
}

// ContractVisibilityConfig represents the configuration as defined by the dApp developer in the smart contract
type ContractVisibilityConfig struct {
	AutoConfig   bool                                       // true for contracts that have no explicit configuration
	Transparent  *bool                                      // users can configure contracts to be fully transparent. All events will be public, and it will expose the internal storage.
	EventConfigs map[gethcommon.Hash]*EventVisibilityConfig // map from the event log signature (topics[0]) to the settings
}

type TxExecResult struct {
	Receipt          *types.Receipt
	CreatedContracts map[gethcommon.Address]*ContractVisibilityConfig
	Err              error
}

type TxExecResults []*TxExecResult
