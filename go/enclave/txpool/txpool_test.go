package txpool

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
)

func TestTxPool_AddTransaction_Pending(t *testing.T) {
	chainID := datagenerator.RandomUInt64()
	mockStore := newMockStorage()
	mockRegistry := newMockBatchRegistry()
	w := datagenerator.RandomWallet(int64(chainID))

	genesisState, err := applyGenesisState(mockStore, []gethcommon.Address{w.Address()})
	require.NoError(t, err)
	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			ParentHash: common.L2BatchHash{},
			// L1Proof:          common.ha,
			Root:             genesisState,
			TxHash:           types.EmptyRootHash,
			Number:           big.NewInt(int64(0)),
			SequencerOrderNo: big.NewInt(int64(common.L2GenesisSeqNo)), // genesis batch has seq number 1
			ReceiptHash:      types.EmptyRootHash,
			TransfersTree:    types.EmptyRootHash,
			// Time:             timeNow,
			// Coinbase:         coinbase,
			// BaseFee:          baseFee,
			GasLimit: 1_000_000_000_000, // todo (@siliev) - does the batch header need uint64?
		},
		Transactions: []*common.L2Tx{},
	}

	err = mockStore.StoreExecutedBatch(genesisBatch, nil)
	require.NoError(t, err)

	mockRegistry.OnBatchExecuted(genesisBatch, nil)

	blockchain := ethchainadapter.NewEthChainAdapter(
		big.NewInt(int64(chainID)),
		mockRegistry,
		mockStore,
		testlog.Logger(),
	)
	err = blockchain.IngestNewBlock(genesisBatch)
	require.NoError(t, err)

	txPool, err := NewTxPool(blockchain, big.NewInt(1), testlog.Logger())
	require.NoError(t, err)

	// Start the TxPool
	err = txPool.Start()
	require.NoError(t, err)

	// Create and add a transaction
	randAddr := datagenerator.RandomAddress()
	transaction := &types.LegacyTx{
		Nonce:    0,
		Value:    big.NewInt(1_000_000_000),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
		To:       &randAddr,
	}
	signedTx, err := w.SignTransaction(transaction)
	require.NoError(t, err)

	err = txPool.Add(signedTx)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	time.Sleep(time.Second) // make sure the tx makes into the pool

	// Check if the transaction is in pending
	pendingTxs := txPool.PendingTransactions()
	require.Equal(t, len(pendingTxs), 1)
	require.Equal(t, pendingTxs[w.Address()][0].Hash.Hex(), signedTx.Hash().Hex())

	// TODO Mint a block and check if it's cleared from the pool
}

func applyGenesisState(storage *mockStorage, accounts []gethcommon.Address) (common.StateRoot, error) {
	statedb, err := state.New(types.EmptyRootHash, storage.stateDB, nil)
	if err != nil {
		return common.StateRoot{}, fmt.Errorf("could not create state DB. Cause: %w", err)
	}

	// set the accounts funds
	for _, acc := range accounts {
		statedb.SetBalance(acc, big.NewInt(1_000_000_000_000_00))
	}

	_ = statedb.IntermediateRoot(true)
	commit, err := statedb.Commit(0, true)
	if err != nil {
		return common.StateRoot{}, err
	}
	return commit, nil
}
