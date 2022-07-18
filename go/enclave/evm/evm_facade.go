package evm

import (
	"math"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	core2 "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// The balance allocated to the sender when they perform a transaction, to ensure they have enough gas.
const prealloc = 750000000000000000

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*common.L2Tx, s *state.StateDB, header *common.Header, rollupResolver db.RollupResolver, chainID int64, fromTxIndex int) map[common.TxHash]interface{} {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)
	zero := uint64(0)
	usedGas := &zero
	result := map[common.TxHash]interface{}{}
	for i, t := range txs {
		r, err := executeTransaction(s, cc, chain, gp, header, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			result[t.Hash()] = err
			common.LogTXExecution(t.Hash(), "Error: %s", err)
			continue
		}
		result[t.Hash()] = r
		if r.Status == types.ReceiptStatusFailed {
			common.LogTXExecution(t.Hash(), "Unsuccessful (status != 1).")
		} else {
			common.LogTXExecution(t.Hash(), "Successfully executed. Address: %s", r.ContractAddress.Hex())
		}
	}
	s.Finalise(true)
	return result
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

	snap := s.Snapshot() //nolint

	// Add some balance to the sender to avoid gas issues.
	// Todo - this has to be removed once the gas logic is sorted.
	s.AddBalance(sender, big.NewInt(prealloc))

	// todo - Author?
	receipt, err := core2.ApplyTransaction(cc, chain, nil, gp, s, convertToEthHeader(header), t, usedGas, vmCfg)
	if err != nil {
		s.RevertToSnapshot(snap)
		return nil, err
	}

	return receipt, nil
}

// ExecuteOffChainCall - executes the "data" command against the "to" smart contract
func ExecuteOffChainCall(from gethcommon.Address, to gethcommon.Address, data []byte, s *state.StateDB, header *common.Header, rollupResolver db.RollupResolver, chainID int64) (*core2.ExecutionResult, error) {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)

	blockContext := core2.NewEVMBlockContext(convertToEthHeader(header), chain, &header.Agg)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, s, cc, vmCfg)
	// use ToMessage
	msg := types.NewMessage(from, &to, 0, gethcommon.Big0, 100_000, gethcommon.Big0, gethcommon.Big0, gethcommon.Big0, data, nil, true)
	result, err := core2.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	// todo - find out why this was called since it's not being called in geth
	// s.Finalise(true)
	return result, nil
}

func initParams(rollupResolver db.RollupResolver, chainID int64) (*ObscuroChainContext, *params.ChainConfig, vm.Config, *core2.GasPool) {
	chain := &ObscuroChainContext{rollupResolver: rollupResolver}
	// TODO - Consolidate this config with the one used in storage.go.
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
