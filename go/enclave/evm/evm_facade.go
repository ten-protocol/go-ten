package evm

import (
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"

	core2 "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/enclave/db"
)

// The balance allocated to the sender when they perform a transaction, to ensure they have enough gas.
const prealloc = 750000000000000000

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*common.L2Tx, s *state.StateDB, header *common.Header, rollupResolver db.RollupResolver, chainID int64, fromTxIndex int) []*types.Receipt {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)
	zero := uint64(0)
	usedGas := &zero
	receipts := make([]*types.Receipt, 0)
	for i, t := range txs {
		r, err := executeTransaction(s, cc, chain, gp, header, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			log.Info("Error transaction %d: %s", common.ShortHash(t.Hash()), err)
			continue
		}
		receipts = append(receipts, r)
		if r.Status == types.ReceiptStatusFailed {
			log.Info("Unsuccessful (status != 1) tx %d.", common.ShortHash(t.Hash()))
		} else {
			log.Info("Successfully executed tx %d.", common.ShortHash(t.Hash()))
		}
	}
	s.Finalise(true)
	return receipts
}

func executeTransaction(s *state.StateDB, cc *params.ChainConfig, chain *ObscuroChainContext, gp *core2.GasPool, header *common.Header, t *common.L2Tx, usedGas *uint64, vmCfg vm.Config, tCount int) (*types.Receipt, error) {
	s.Prepare(t.Hash(), tCount)

	// Allocates a large balance to the sender, to allow them to perform any transaction.
	// TODO - Allocate balances properly.
	signer := types.NewLondonSigner(cc.ChainID)
	sender, err := types.Sender(signer, t)
	if err != nil {
		return nil, err
	}
	s.AddBalance(sender, big.NewInt(prealloc))

	snap := s.Snapshot()
	// todo - Author?
	receipt, err := core2.ApplyTransaction(cc, chain, nil, gp, s, convertToEthHeader(header), t, usedGas, vmCfg)
	if err == nil {
		return receipt, nil
	}
	s.RevertToSnapshot(snap)
	return nil, err
}

// ExecuteOffChainCall - executes the "data" command against the "to" smart contract
func ExecuteOffChainCall(from gethcommon.Address, to gethcommon.Address, data []byte, s *state.StateDB, header *common.Header, rollupResolver db.RollupResolver, chainID int64) (*core2.ExecutionResult, error) {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)

	blockContext := core2.NewEVMBlockContext(convertToEthHeader(header), chain, &header.Agg)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, s, cc, vmCfg)

	msg := types.NewMessage(from, &to, 0, gethcommon.Big0, 100_000, gethcommon.Big0, gethcommon.Big0, gethcommon.Big0, data, nil, true)
	result, err := core2.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	s.Finalise(true)
	return result, nil
}

func initParams(rollupResolver db.RollupResolver, chainID int64) (*ObscuroChainContext, *params.ChainConfig, vm.Config, *core2.GasPool) {
	chain := &ObscuroChainContext{rollupResolver: rollupResolver}
	cc := &params.ChainConfig{
		ChainID:     big.NewInt(chainID),
		LondonBlock: gethcommon.Big0,
	}
	vmCfg := vm.Config{
		NoBaseFee: true,
	}
	gp := core2.GasPool(math.MaxUint64)
	return chain, cc, vmCfg, &gp
}
