package params

import (
	"time"

	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
)

// SimParams are the parameters for setting up the simulation.
type SimParams struct {
	NumberOfNodes int

	// A critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	AvgBlockDuration  time.Duration
	AvgNetworkLatency time.Duration // artificial latency injected between sending and receiving messages on the mock network

	SimulationTime time.Duration // how long the simulations should run for

	L1EfficiencyThreshold float64
	L1BeaconPort          int

	ContractRegistryLib contractlib.ContractRegistryLib
	// ERC20ContractLib allows parsing ERC20Contract txs to and from the eth txs
	ERC20ContractLib erc20contractlib.ERC20ContractLib

	BlobResolver l1.BlobResolver
	L1TenData    *L1TenData

	// Contains all the wallets required by the simulation
	Wallets *SimWallets

	StartPort int  // The port from which to start allocating ports. Must be unique across all simulations.
	IsInMem   bool // Denotes that the sim does not have a full RPC layer.

	ReceiptTimeout time.Duration // How long to wait for transactions to be confirmed.

	StoppingDelay              time.Duration // How long to wait between injection and verification
	NodeWithInboundP2PDisabled int
	WithPrefunding             bool
}

type L1TenData struct {
	// TenStartBlock is the L1 block hash where the TEN network activity begins (e.g. mgmt contract deployment)
	TenStartBlock common.Hash
	// NetworkConfigAddress defines the network config contract address
	NetworkConfigAddress common.Address
	// EnclaveRegistryAddress defines the network enclave registry contract address
	EnclaveRegistryAddress common.Address
	// CrossChainContractAddress defines the cross chain contract address
	CrossChainContractAddress common.Address
	// DataAvailabilityRegistryAddress defines the rollup contract address
	DataAvailabilityRegistryAddress common.Address
	// ObxErc20Address - the address of the "TEN" ERC20
	ObxErc20Address common.Address
	// EthErc20Address - the address of the "ETH" ERC20
	EthErc20Address common.Address
	// MessageBusAddr - the address of the L1 message bus.
	MessageBusAddr common.Address
	// CrossChainMessengerAddress - the address of the L1 cross chain messenger.
	CrossChainMessengerAddress common.Address
	// BridgeAddress - the address of the L1 bridge.
	BridgeAddress common.Address
}
