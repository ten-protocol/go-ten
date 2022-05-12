package params

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/wallet"
)

// SimParams are the parameters for setting up the simulation.
type SimParams struct {
	NumberOfNodes          int
	NumberOfObscuroWallets int

	// A critical parameter of the simulation. The value should be as low as possible, as long as the test is still meaningful
	AvgBlockDuration  time.Duration
	AvgNetworkLatency time.Duration // artificial latency injected between sending and receiving messages on the mock network
	AvgGossipPeriod   time.Duration // POBI protocol setting

	SimulationTime time.Duration // how long the simulations should run for

	// EfficiencyThresholds represents an acceptable "dead blocks" percentage for this simulation.
	// dead blocks - Blocks that are produced and gossiped, but don't make it into the canonical chain.
	// We test the results against this threshold to catch eventual protocol errors.
	L1EfficiencyThreshold     float64
	L2EfficiencyThreshold     float64 // number of dead obscuro blocks
	L2ToL1EfficiencyThreshold float64 // number of ethereum blocks that don't include an obscuro node

	// TxHandler defines how the simulation should unpack transactions
	TxHandler mgmtcontractlib.TxHandler
	// MgmtContractAddr defines the management contract address
	MgmtContractAddr common.Address
	// EthWallets contains the wallets to use for the l1 nodes
	EthWallets        []wallet.Wallet
	ERC20ContractAddr common.Address
}
