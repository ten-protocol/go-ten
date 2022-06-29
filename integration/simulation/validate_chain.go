package simulation

import (
	"fmt"
	"sync"
	"testing"

	"github.com/obscuronet/obscuro-playground/go/enclave/bridge"

	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"

	"github.com/obscuronet/obscuro-playground/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/common"
)

const (
	// The threshold number of transactions below which we consider the simulation to have failed. We generally expect far
	// more than this, but this is a sanity check to ensure the simulation doesn't stop after a single transaction of each
	// type, for example.
	txThreshold = 5
	// The maximum number of blocks an Obscuro node can fall behind
	maxBlockDelay = 5
)

// After a simulation has run, check as much as possible that the outputs of the simulation are expected.
// For example, all injected transactions were processed correctly, the height of the rollup chain is a function of the total
// time of the simulation and the average block duration, that all Obscuro nodes are roughly in sync, etc
func checkNetworkValidity(t *testing.T, s *Simulation) {
	// ensure L1 and L2 txs were issued
	if len(s.TxInjector.Counter.L1Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d L1 transactions. At least %d expected", len(s.TxInjector.Counter.L1Transactions), txThreshold)
	}
	if len(s.TxInjector.Counter.TransferL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d transfer L2 transactions. At least %d expected", len(s.TxInjector.Counter.TransferL2Transactions), txThreshold)
	}
	if len(s.TxInjector.Counter.WithdrawalL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d withdrawal L2 transactions. At least %d expected", len(s.TxInjector.Counter.WithdrawalL2Transactions), txThreshold)
	}

	l1MaxHeight := checkEthereumBlockchainValidity(t, s)
	checkObscuroBlockchainValidity(t, s, l1MaxHeight)
}

// checkEthereumBlockchainValidity: sanity check of the state of all L1 nodes
// - the chain has a minimum number of blocks
// - the chain height is similar across all ethereum nodes
// - there are no duplicate txs
// - efficiency - number of created blocks/height
// - no reorgs
func checkEthereumBlockchainValidity(t *testing.T, s *Simulation) uint64 {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	heights := make([]uint64, len(s.EthClients))
	for i, node := range s.EthClients {
		heights[i] = checkBlockchainOfEthereumNode(t, node, minHeight, s)
	}

	min, max := minMax(heights)
	if max-min > max/10 {
		t.Logf("Node Heights: %v", heights)
		t.Errorf("There is a problem with the mock Ethereum chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}

	return max
}

// checkObscuroBlockchainValidity - perform the following checks
// - minimum height - the chain has a minimum number of rollups
// - check height is similar
// - check no duplicate txs
// - check efficiency - no of created blocks/ height
// - check amount in the system
// - check withdrawals/deposits
func checkObscuroBlockchainValidity(t *testing.T, s *Simulation, maxL1Height uint64) {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	// process the blockchain of each node in parallel to minimize the difference between them since they are still running
	heights := make([]uint64, len(s.ObscuroClients))
	var wg sync.WaitGroup
	for idx := range s.ObscuroClients {
		obscuroClient := s.ObscuroClients[idx]
		wg.Add(1)
		go checkBlockchainOfObscuroNode(t, obscuroClient, minHeight, maxL1Height, s, &wg, heights, idx)
	}
	wg.Wait()
	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	if max-min > max/10 {
		t.Errorf("There is a problem with the Obscuro chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethadapter.EthClient, minHeight uint64, s *Simulation) uint64 {
	nodeAddr := common.ShortAddress(node.Info().ID)
	head := node.FetchHeadBlock()
	height := head.NumberU64()

	if height < minHeight {
		t.Errorf("Node %d: There were only %d blocks mined. Expected at least: %d.", nodeAddr, height, minHeight)
	}

	deposits, rollups, totalDeposited, blockCount := ExtractDataFromEthereumChain(common.GenesisBlock, head, node, s)
	s.Stats.TotalL1Blocks = uint64(blockCount)

	if len(findHashDups(deposits)) > 0 {
		dups := findHashDups(deposits)
		t.Errorf("Node %d: Found Deposit duplicates: %v", nodeAddr, dups)
	}
	if len(findRollupDups(rollups)) > 0 {
		dups := findRollupDups(rollups)
		t.Errorf("Node %d: Found Rollup duplicates: %v", nodeAddr, dups)
	}
	if totalDeposited != s.Stats.TotalDepositedAmount {
		t.Errorf("Node %d: Deposit amounts don't match. Found %d , expected %d", nodeAddr, totalDeposited, s.Stats.TotalDepositedAmount)
	}

	efficiency := float64(s.Stats.TotalL1Blocks-height) / float64(s.Stats.TotalL1Blocks)
	if efficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: Efficiency in L1 is %f. Expected:%f. Number: %d.", nodeAddr, efficiency, s.Params.L1EfficiencyThreshold, height)
	}

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().ID]
	eff := float64(reorgs) / float64(height)
	if eff > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: The number of reorgs is too high: %d. ", nodeAddr, reorgs)
	}
	return height
}

// ExtractDataFromEthereumChain returns the deposits, rollups, total amount deposited and length of the blockchain
// between the start block and the end block.
func ExtractDataFromEthereumChain(startBlock *types.Block, endBlock *types.Block, node ethadapter.EthClient, s *Simulation) ([]gethcommon.Hash, []common.L2RootHash, uint64, int) {
	deposits := make([]gethcommon.Hash, 0)
	rollups := make([]common.L2RootHash, 0)
	totalDeposited := uint64(0)

	blockchain := node.BlocksBetween(startBlock, endBlock)
	for _, block := range blockchain {
		for _, tx := range block.Transactions() {
			t := s.Params.ERC20ContractLib.DecodeTx(tx)
			if t == nil {
				t = s.Params.MgmtContractLib.DecodeTx(tx)
			}

			if t == nil {
				continue
			}
			switch l1tx := t.(type) {
			case *ethadapter.L1DepositTx:
				deposits = append(deposits, tx.Hash())
				totalDeposited += l1tx.Amount
			case *ethadapter.L1RollupTx:
				r := common.DecodeRollupOrPanic(l1tx.Rollup)
				rollups = append(rollups, r.Hash())
				if node.IsBlockAncestor(block, r.Header.L1Proof) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					s.Stats.NewRollup(node.Info().ID, r)
				}
			}
		}
	}
	return deposits, rollups, totalDeposited, len(blockchain)
}

func checkBlockchainOfObscuroNode(
	t *testing.T,
	nodeClient rpcclientlib.Client,
	minObscuroHeight uint64,
	maxEthereumHeight uint64,
	s *Simulation,
	wg *sync.WaitGroup,
	heights []uint64,
	nodeIdx int,
) {
	defer wg.Done()
	var nodeID gethcommon.Address
	err := nodeClient.Call(&nodeID, rpcclientlib.RPCGetID)
	if err != nil {
		t.Errorf("Could not retrieve Obscuro node's address when checking blockchain.")
	}
	nodeAddr := common.ShortAddress(nodeID)
	l1Height := getCurrentBlockHeadHeight(nodeClient)

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	if int64(maxEthereumHeight)-l1Height > maxBlockDelay {
		t.Errorf("Node %d: Obscuro node fell behind by %d blocks.", nodeAddr, maxEthereumHeight-uint64(l1Height))
	}

	// check that the height of the Rollup chain is higher than a minimum expected value.
	h := getCurrentRollupHead(nodeClient)

	if h == nil {
		t.Errorf("Node %d: No head rollup recorded. Skipping any further checks for this node.\n", nodeAddr)
		return
	}
	l2Height := h.Number
	if l2Height.Uint64() < minObscuroHeight {
		t.Errorf("Node %d: Node only mined %d rollups. Expected at least: %d.", l2Height, nodeAddr, minObscuroHeight)
	}

	totalL2Blocks := s.Stats.NoL2Blocks[nodeID]
	// in case the blockchain has advanced above what was collected, there is no longer a point to this check
	if l2Height.Uint64() <= totalL2Blocks {
		efficiencyL2 := float64(totalL2Blocks-l2Height.Uint64()) / float64(totalL2Blocks)
		if efficiencyL2 > s.Params.L2EfficiencyThreshold {
			t.Errorf("Node %d: Efficiency in L2 is %f. Expected:%f", nodeAddr, efficiencyL2, s.Params.L2EfficiencyThreshold)
		}
	}

	// check that the pobi protocol doesn't waste too many blocks.
	// todo- find the block where the genesis was published)
	efficiency := float64(uint64(l1Height)-l2Height.Uint64()) / float64(l1Height)
	if efficiency > s.Params.L2ToL1EfficiencyThreshold {
		t.Errorf("Node %d: L2 to L1 Efficiency is %f. Expected:%f", nodeAddr, efficiency, s.Params.L2ToL1EfficiencyThreshold)
	}

	notFoundTransfers, notFoundWithdrawals := FindNotIncludedL2Txs(nodeClient, s.TxInjector)
	if notFoundTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Transfer Txs not found in the enclave",
			nodeAddr, notFoundTransfers, len(s.TxInjector.Counter.TransferL2Transactions))
	}
	if notFoundWithdrawals > 0 {
		t.Errorf("Node %d: %d out of %d Withdrawal Txs not found in the enclave",
			nodeAddr, notFoundWithdrawals, len(s.TxInjector.Counter.WithdrawalL2Transactions))
	}

	totalSuccessfullyWithdrawn, numberOfWithdrawalRequests := extractWithdrawals(t, nodeClient, nodeAddr)

	// sanity check number of withdrawal transaction
	if numberOfWithdrawalRequests > len(s.TxInjector.Counter.GetL2WithdrawalRequests()) {
		t.Errorf("Node %d: found more transactions in the blockchain than the generated by the tx manager", nodeAddr)
	}

	injectorDepositedAmt := uint64(0)
	for _, tx := range s.TxInjector.Counter.GetL1Transactions() {
		if depTx, ok := tx.(*ethadapter.L1DepositTx); ok {
			injectorDepositedAmt += depTx.Amount
		}
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if totalSuccessfullyWithdrawn > s.Stats.TotalWithdrawalRequestedAmount {
		t.Errorf("Node %d: The amount withdrawn %d exceeds the actual amount requested %d", nodeAddr, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// sanity check that the injected withdrawals were mostly executed
	if totalSuccessfullyWithdrawn < s.Stats.TotalWithdrawalRequestedAmount/2 {
		t.Errorf("Node %d: The amount withdrawn %d is far smaller than the amount requested %d", nodeAddr, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// check that the sum of all balances matches the total amount of money that must be in the system
	totalAmountInSystem := s.Stats.TotalDepositedAmount - totalSuccessfullyWithdrawn
	total := uint64(0)
	for _, wallet := range s.Params.Wallets.SimObsWallets {
		total += balance(nodeClient, wallet.Address(), s.Params.Wallets.Tokens[bridge.BTC].L2ContractAddress)
	}

	if total != totalAmountInSystem {
		t.Errorf("Node %d: The amount of money in accounts does not match the amount deposited. Found %d , expected %d", nodeAddr, total, totalAmountInSystem)
	}
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// (execute deposits and transactions and compare to the state in the rollup)

	heights[nodeIdx] = l2Height.Uint64()
}

// FindNotIncludedL2Txs returns the number of transfers and withdrawals that were injected but are not present in the L2 blockchain.
func FindNotIncludedL2Txs(l2Client rpcclientlib.Client, txInjector *TransactionInjector) (int, int) {
	transfers, withdrawals := txInjector.Counter.GetL2Transactions()
	notFoundTransfers := 0
	for _, tx := range transfers {
		if l2tx := getTransaction(l2Client, tx.Hash()); l2tx == nil {
			notFoundTransfers++
		}
	}

	notFoundWithdrawals := 0
	for _, tx := range withdrawals {
		if l2tx := getTransaction(l2Client, tx.Hash()); l2tx == nil {
			notFoundWithdrawals++
		}
	}

	return notFoundTransfers, notFoundWithdrawals
}

func extractWithdrawals(t *testing.T, nodeClient rpcclientlib.Client, nodeAddr uint64) (totalSuccessfullyWithdrawn uint64, numberOfWithdrawalRequests int) {
	head := getCurrentRollupHead(nodeClient)

	if head == nil {
		panic(fmt.Sprintf("Node %d: The current head should not be nil", nodeAddr))
	}

	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for r := head; ; r = getRollupHeader(nodeClient, r.ParentHash) {
		if r != nil && r.Number.Uint64() == common.L1GenesisHeight {
			return
		}
		if r == nil {
			t.Errorf(fmt.Sprintf("Node %d: Reached a missing rollup", nodeAddr))
			return
		}
		for _, w := range r.Withdrawals {
			totalSuccessfullyWithdrawn += w.Amount
			numberOfWithdrawalRequests++
		}
	}
}
