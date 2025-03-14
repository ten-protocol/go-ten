package networktest

import (
	"context"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/wallet"

	"github.com/ethereum/go-ethereum/common"
)

// NetworkConnector represents the network being tested against, e.g. testnet, dev-testnet, dev-sim
//
// # It provides network details (standard contract addresses) and easy client setup for sim users
//
// Note: some of these methods may not be available for some networks (e.g. MC Owner wallet for live testnets)
type NetworkConnector interface {
	ChainID() int64
	// AllocateFaucetFunds uses the networks default faucet mechanism for allocating funds to a test account
	AllocateFaucetFunds(ctx context.Context, account common.Address) error
	SequencerRPCAddress() string
	ValidatorRPCAddress(idx int) string
	NumValidators() int
	GetSequencerNode() NodeOperator
	GetValidatorNode(idx int) NodeOperator
	GetL1Client() (ethadapter.EthClient, error)
	GetContractOwnerWallet() (wallet.Wallet, error) // wallet that owns the management contract (network admin)
	GetGatewayURL() (string, error)
	GetGatewayWSURL() (string, error)
}

// Action is any step in a test, they will typically be either minimally small steps in the test or they will be containers
// that coordinate the running of multiple sub-actions (e.g. SeriesAction/ParallelAction)
//
// With these action containers a tree of actions is built to form a test.
//
// A test will call Run on an action (triggering the run of the tree of subactions),
// and then Verify (orchestrating the verification step on all sub actions)
//
// Conventions:
//   - if an action name begins with `Verify` then its `Run` method will be a no-op, these should be at the end of a test run (since they only test the post-test state)
type Action interface {
	Run(ctx context.Context, network NetworkConnector) (context.Context, error)
	Verify(ctx context.Context, network NetworkConnector) error
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
	Start() error
	Stop() error

	StartEnclave(idx int) error
	StopEnclave(idx int) error
	StartHost() error
	StopHost() error

	HostRPCHTTPAddress() string
	HostRPCWSAddress() string
}
