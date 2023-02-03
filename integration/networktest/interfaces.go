package networktest

import (
	"github.com/ethereum/go-ethereum/common"
)

// NetworkConnector represents the network being tested against, e.g. testnet, dev-testnet, dev-sim
//
// It provides network details (standard contract addresses) and easy client setup for sim users
type NetworkConnector interface {
	ChainID() int64
	// AllocateFaucetFunds uses the networks default faucet mechanism for allocating funds to a test account
	AllocateFaucetFunds(account common.Address) error
	SequencerRPCAddress() string
	ValidatorRPCAddress(idx int) string
	NumValidators() int
	GetSequencerNode() NodeOperator
	GetValidatorNode(idx int) NodeOperator
}

// NetworkTest defines a test, it can be run against an arbitrary network connector (same tests can run against different environments)
type NetworkTest interface {
	Run(network NetworkConnector) error
	Name() string // used for logfile naming
}

// Environment abstraction allows the test runner to prepare the network connector with optional config and steps
// and to handle the clean-up so that different types of NetworkConnector can be configured in the same style (see runner.go)
// (local network requires setup and teardown for example, whereas a connector to a Testnet is ready to go)
type Environment interface {
	Prepare() (NetworkConnector, func(), error)
}

// NodeOperator is used by the DevNetwork for orchestrating different scenarios
//
// It attempts to encapsulate the data, monitoring and possible actions that would be available to a real NodeOperator
// in a live permissionless Obscuro network
type NodeOperator interface {
	Start()
	Stop()

	StartEnclave()
	StopEnclave()
	StartHost()
	StopHost()

	HostRPCAddress() string
}
