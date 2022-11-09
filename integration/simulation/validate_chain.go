package simulation

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"

	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// The threshold number of transactions below which we consider the simulation to have failed. We generally expect far
	// more than this, but this is a sanity check to ensure the simulation doesn't stop after a single transaction of each
	// type, for example.
	txThreshold = 5
	// As above, but for the number of logs received via subscriptions.
	logsThreshold = 5
	// The maximum number of blocks an Obscuro node can fall behind
	maxBlockDelay = 5
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// After a simulation has run, check as much as possible that the outputs of the simulation are expected.
// For example, all injected transactions were processed correctly, the height of the rollup chain is a function of the total
// time of the simulation and the average block duration, that all Obscuro nodes are roughly in sync, etc
func checkNetworkValidity(t *testing.T, s *Simulation) {
	// ensure L1 and L2 txs were issued
	if len(s.TxInjector.TxTracker.L1Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d L1 transactions. At least %d expected", len(s.TxInjector.TxTracker.L1Transactions), txThreshold)
	}
	if len(s.TxInjector.TxTracker.TransferL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d transfer L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.TransferL2Transactions), txThreshold)
	}
	if len(s.TxInjector.TxTracker.WithdrawalL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d withdrawal L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.WithdrawalL2Transactions), txThreshold)
	}

	l1MaxHeight := checkEthereumBlockchainValidity(t, s)
	checkObscuroBlockchainValidity(t, s, l1MaxHeight)
	checkReceivedLogs(t, s)
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

	heights := make([]uint64, len(s.RPCHandles.EthClients))
	for i, node := range s.RPCHandles.EthClients {
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
	heights := make([]uint64, len(s.RPCHandles.ObscuroClients))
	var wg sync.WaitGroup
	for idx := range s.RPCHandles.ObscuroClients {
		wg.Add(1)
		go checkBlockchainOfObscuroNode(t, s.RPCHandles, minHeight, maxL1Height, s, &wg, heights, idx)
	}
	wg.Wait()
	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	if max-min > max/10 {
		t.Errorf("There is a problem with the Obscuro chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethadapter.EthClient, minHeight uint64, s *Simulation) uint64 {
	nodeAddr := common.ShortAddress(node.Info().L2ID)
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

	if s.Stats.TotalDepositedAmount.Cmp(totalDeposited) != 0 {
		t.Errorf("Node %d: Deposit amounts don't match. Found %d , expected %d", nodeAddr, totalDeposited, s.Stats.TotalDepositedAmount)
	}

	efficiency := float64(s.Stats.TotalL1Blocks-height) / float64(s.Stats.TotalL1Blocks)
	if efficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: Efficiency in L1 is %f. Expected:%f. Number: %d.", nodeAddr, efficiency, s.Params.L1EfficiencyThreshold, height)
	}

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().L2ID]
	reorgEfficiency := float64(reorgs) / float64(height)
	if reorgEfficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: The number of reorgs is too high: %d. ", nodeAddr, reorgs)
	}

	// Check that all the rollups are produced by aggregators.
	for _, rollup := range rollups {
		aggregatorID, err := strconv.ParseInt(rollup.Header.Agg.Hex()[2:], 16, 64)
		if err != nil {
			t.Errorf("Node %d: Could not parse node's integer ID. Cause: %s", nodeAddr, err)
			continue
		}
		if network.GetNodeType(int(aggregatorID)) != common.Aggregator {
			t.Errorf("Node %d: Found rollup produced by non-aggregator %d", nodeAddr, aggregatorID)
		}
	}

	return height
}

// ExtractDataFromEthereumChain returns the deposits, rollups, total amount deposited and length of the blockchain
// between the start block and the end block.
func ExtractDataFromEthereumChain(startBlock *types.Block, endBlock *types.Block, node ethadapter.EthClient, s *Simulation) ([]gethcommon.Hash, []*common.ExtRollupWithHash, *big.Int, int) {
	deposits := make([]gethcommon.Hash, 0)
	rollups := make([]*common.ExtRollupWithHash, 0)
	totalDeposited := big.NewInt(0)

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
				totalDeposited.Add(totalDeposited, l1tx.Amount)
			case *ethadapter.L1RollupTx:
				r, err := common.DecodeRollup(l1tx.Rollup)
				if err != nil {
					testlog.Logger().Crit("could not decode rollup. ", log.ErrKey, err)
				}
				rollups = append(rollups, r)
				if node.IsBlockAncestor(block, r.Header.L1Proof) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					s.Stats.NewRollup(node.Info().L2ID)
				}
			}
		}
	}
	return deposits, rollups, totalDeposited, len(blockchain)
}

func checkBlockchainOfObscuroNode(t *testing.T, rpcHandles *network.RPCHandles, minObscuroHeight uint64, maxEthereumHeight uint64, s *Simulation, wg *sync.WaitGroup, heights []uint64, nodeIdx int) {
	defer wg.Done()
	var nodeID gethcommon.Address
	nodeClient := rpcHandles.ObscuroClients[nodeIdx]
	obsClient := obsclient.NewObsClient(nodeClient)

	err := nodeClient.Call(&nodeID, rpc.GetID)
	if err != nil {
		t.Errorf("Could not retrieve Obscuro node's address when checking blockchain.")
	}
	nodeAddr := common.ShortAddress(nodeID)

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	l1Height, err := obsClient.BlockNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve L1 height. Cause: %s", nodeAddr, err)
	}
	if int(maxEthereumHeight)-int(l1Height) > maxBlockDelay {
		t.Errorf("Node %d: Obscuro node fell behind by %d blocks.", nodeAddr, maxEthereumHeight-uint64(l1Height))
	}

	// check that the height of the Rollup chain is higher than a minimum expected value.
	h := getHeadRollupHeader(obsClient)
	if h == nil {
		t.Errorf("Node %d: No head rollup recorded. Skipping any further checks for this node.\n", nodeAddr)
		return
	}
	l2Height := h.Number
	if l2Height.Uint64() < minObscuroHeight {
		t.Errorf("Node %d: Node only mined %d rollups. Expected at least: %d.", nodeAddr, l2Height, minObscuroHeight)
	}

	// check that the height from the rollup header is consistent with the height returned by eth_blockNumber.
	l2HeightFromRollupNumber, err := obsClient.RollupNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve block number. Cause: %s", nodeAddr, err)
	}
	if l2HeightFromRollupNumber != l2Height.Uint64() {
		t.Errorf("Node %d: Node's head rollup had a height %d, but %s height was %d", nodeAddr, l2Height, rpc.BlockNumber, l2HeightFromRollupNumber)
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

	notFoundTransfers, notFoundWithdrawals := FindNotIncludedL2Txs(s.ctx, nodeIdx, rpcHandles, s.TxInjector)
	if notFoundTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Transfer Txs not found in the enclave",
			nodeAddr, notFoundTransfers, len(s.TxInjector.TxTracker.TransferL2Transactions))
	}
	if notFoundWithdrawals > 0 {
		t.Errorf("Node %d: %d out of %d Withdrawal Txs not found in the enclave",
			nodeAddr, notFoundWithdrawals, len(s.TxInjector.TxTracker.WithdrawalL2Transactions))
	}

	checkTransactionReceipts(s.ctx, t, nodeIdx, rpcHandles, s.TxInjector)

	totalSuccessfullyWithdrawn, numberOfWithdrawalRequests := extractWithdrawals(t, nodeClient, nodeAddr)

	// sanity check number of withdrawal transaction
	if numberOfWithdrawalRequests > len(s.TxInjector.TxTracker.GetL2WithdrawalRequests()) {
		t.Errorf("Node %d: found more transactions in the blockchain than the generated by the tx manager", nodeAddr)
	}

	injectorDepositedAmt := big.NewInt(0)
	for _, tx := range s.TxInjector.TxTracker.GetL1Transactions() {
		if depTx, ok := tx.(*ethadapter.L1DepositTx); ok {
			injectorDepositedAmt.Add(injectorDepositedAmt, depTx.Amount)
		}
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if totalSuccessfullyWithdrawn.Cmp(s.Stats.TotalWithdrawalRequestedAmount) > 0 {
		t.Errorf("Node %d: The amount withdrawn %d exceeds the actual amount requested %d", nodeAddr, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// sanity check that the injected withdrawals were mostly executed

	halfRequestedWithdrawalAmt := big.NewInt(0).Div(s.Stats.TotalWithdrawalRequestedAmount, big.NewInt(2))
	if totalSuccessfullyWithdrawn.Cmp(halfRequestedWithdrawalAmt) < 0 {
		t.Errorf("Node %d: The amount withdrawn %d is far smaller than the amount requested %d", nodeAddr, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// check that the sum of all balances matches the total amount of money that must be in the system
	totalAmountInSystem := big.NewInt(0).Sub(s.Stats.TotalDepositedAmount, totalSuccessfullyWithdrawn)
	total := big.NewInt(0)
	for _, wallet := range s.Params.Wallets.SimObsWallets {
		client := rpcHandles.ObscuroWalletClient(wallet.Address(), nodeIdx)
		bal := balance(s.ctx, client, wallet.Address(), s.Params.Wallets.Tokens[bridge.HOC].L2ContractAddress)
		total.Add(total, bal)
	}

	if total.Cmp(totalAmountInSystem) != 0 {
		t.Errorf("Node %d: The amount of money in accounts does not match the amount deposited. Found %d , expected %d", nodeAddr, total, totalAmountInSystem)
	}
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// (execute deposits and transactions and compare to the state in the rollup)

	heights[nodeIdx] = l2Height.Uint64()
}

// FindNotIncludedL2Txs returns the number of transfers and withdrawals that were injected but are not present in the L2 blockchain.
func FindNotIncludedL2Txs(ctx context.Context, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) (int, int) {
	transfers, withdrawals := txInjector.TxTracker.GetL2Transactions()
	notFoundTransfers := 0
	for _, tx := range transfers {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil {
			panic(err)
		}
		if l2tx == nil {
			notFoundTransfers++
		}
	}

	notFoundWithdrawals := 0
	for _, tx := range withdrawals {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil {
			panic(err)
		}
		if l2tx == nil {
			notFoundWithdrawals++
		}
	}

	return notFoundTransfers, notFoundWithdrawals
}

func getSender(tx *common.L2Tx) gethcommon.Address {
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		panic(fmt.Errorf("couldn't find sender to verify transaction - %w", err))
	}
	return msg.From()
}

// Checks that there is a receipt available for each L2 transaction.
func checkTransactionReceipts(ctx context.Context, t *testing.T, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) {
	l2Txs := append(txInjector.TxTracker.TransferL2Transactions, txInjector.TxTracker.WithdrawalL2Transactions...)

	for _, tx := range l2Txs {
		sender := getSender(tx)

		// We check that there is a receipt available for each transaction
		receipt, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			t.Errorf("could not retrieve receipt for transaction %s. Cause: %s", tx.Hash().Hex(), err)
			continue
		}

		// We check that the receipt's logs are relevant to the sender.
		for _, retrievedLog := range receipt.Logs {
			assertRelevantLogsOnly(t, sender.Hex(), *retrievedLog)
		}

		if receipt.Status == types.ReceiptStatusFailed {
			testlog.Logger().Info("Transaction failed.", log.TxKey, tx.Hash().Hex())
		}
	}
}

func extractWithdrawals(t *testing.T, nodeClient rpc.Client, nodeAddr uint64) (totalSuccessfullyWithdrawn *big.Int, numberOfWithdrawalRequests int) {
	totalSuccessfullyWithdrawn = big.NewInt(0)
	head := getHeadRollupHeader(obsclient.NewObsClient(nodeClient))

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
			totalSuccessfullyWithdrawn.Add(totalSuccessfullyWithdrawn, w.Amount)
			numberOfWithdrawalRequests++
		}
	}
}

// Terminates all subscriptions and validates the received events.
func checkReceivedLogs(t *testing.T, s *Simulation) {
	logsFromSnapshots := 0
	for _, clients := range s.RPCHandles.AuthObsClients {
		for _, client := range clients {
			logsFromSnapshots += checkSnapshotLogs(t, client)
		}
	}
	if logsFromSnapshots < logsThreshold {
		t.Errorf("only received %d logs from snapshots, expected at least %d", logsFromSnapshots, logsThreshold)
	}

	// In-memory clients cannot handle subscriptions for now.
	if !s.Params.IsInMem {
		for _, sub := range s.Subscriptions {
			sub.Unsubscribe()
		}

		logsFromSubscriptions := 0
		for owner, channels := range s.LogChannels {
			for _, channel := range channels {
				logsFromSubscriptions += checkSubscribedLogs(t, owner, channel)
			}
		}
		if logsFromSubscriptions < logsThreshold {
			t.Errorf("only received %d logs from subscriptions, expected at least %d", logsFromSubscriptions, logsThreshold)
		}
	}
}

// Checks that a subscription has received the expected logs.
func checkSubscribedLogs(t *testing.T, owner string, channel chan common.IDAndLog) int {
	var logs []*types.Log

out:
	for {
		select {
		case idAndLog := <-channel:
			logs = append(logs, idAndLog.Log)

		// The logs will have built up on the channel throughout the simulation, so they should arrive immediately.
		// However, if we use a `default` case, only the first one arrives. Some minimal wait is required.
		case <-time.After(time.Millisecond):
			break out
		}
	}

	assertLogsValid(t, owner, logs)
	return len(logs)
}

func checkSnapshotLogs(t *testing.T, client *obsclient.AuthObsClient) int {
	// To exercise the filtering mechanism, we get a snapshot for HOC events only, ignoring POC events.
	hocFilter := common.FilterCriteriaJSON{
		Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + bridge.HOCAddr)},
	}
	logs, err := client.GetLogs(context.Background(), hocFilter)
	if err != nil {
		t.Errorf("could not retrieve logs for client. Cause: %s", err)
	}

	assertLogsValid(t, client.Address().Hex(), logs)

	return len(logs)
}

// Asserts that the logs meet various criteria.
func assertLogsValid(t *testing.T, owner string, logs []*types.Log) {
	for _, receivedLog := range logs {
		assertRelevantLogsOnly(t, owner, *receivedLog)

		logAddrHex := receivedLog.Address.Hex()
		if logAddrHex != "0x"+bridge.HOCAddr {
			t.Errorf("due to filter, expected logs from the HOC contract only, but got a log from %s", logAddrHex)
		}
	}

	assertNoDupeLogs(t, logs)
}

// Asserts that the log is relevant to the recipient (either a lifecycle event or a relevant user event).
func assertRelevantLogsOnly(t *testing.T, owner string, receivedLog types.Log) {
	// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can (potentially)
	// be user addresses. We filter these out. In theory, we should also check whether the topics exist in the state DB
	// and are contract addresses, but we cannot do this as part of chain validation.
	var userAddrs []string
	for idx, topic := range receivedLog.Topics {
		// The first topic is always the hash of the event.
		if idx == 0 {
			continue
		}

		// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can
		// (potentially) be user addresses.
		topicHex := topic.Hex()
		if topicHex[2:len(zeroBytesHex)+2] == zeroBytesHex {
			userAddrs = append(userAddrs, gethcommon.HexToAddress(topicHex).Hex())
		}
	}

	// If there are no potential user addresses, this is a lifecycle event, and is therefore relevant to everyone.
	if len(userAddrs) == 0 {
		return
	}

	for _, addr := range userAddrs {
		if addr == owner {
			return
		}
	}

	// If we've fallen through to here, it means the log was not relevant.
	t.Errorf("received log that was not relevant (neither a lifecycle event nor relevant to the client's account)")
}

// Asserts that there are no duplicate logs in the provided list.
func assertNoDupeLogs(t *testing.T, logs []*types.Log) {
	logCount := make(map[string]int)

	for _, item := range logs {
		logJSON, err := item.MarshalJSON()
		if err != nil {
			t.Errorf("could not marshal log to JSON to check for duplicate logs")
			continue
		}

		// check if the item/element exist in the duplicate_frequency map
		_, exist := logCount[string(logJSON)]
		if exist {
			logCount[string(logJSON)]++ // increase counter by 1 if already in the map
		} else {
			logCount[string(logJSON)] = 1 // else start counting from 1
		}
	}

	for logJSON, count := range logCount {
		if count > 1 {
			t.Errorf("received duplicate log with body %s", logJSON)
		}
	}
}
