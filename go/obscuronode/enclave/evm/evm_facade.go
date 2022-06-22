package evm

import (
	"math"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	core2 "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*nodecommon.L2Tx, s *state.StateDB, header *nodecommon.Header, rollupResolver db.RollupResolver, chainID int64, fromTxIndex int) []*types.Receipt {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)
	zero := uint64(0)
	usedGas := &zero
	receipts := make([]*types.Receipt, 0)
	for i, t := range txs {
		r, err := executeTransaction(s, cc, chain, gp, header, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			log.Info("Error transaction %d: %s", obscurocommon.ShortHash(t.Hash()), err)
			continue
		}
		receipts = append(receipts, r)
		if r.Status == types.ReceiptStatusFailed {
			log.Info("Unsuccessful (status != 1) tx %d. Logs: %+v", obscurocommon.ShortHash(t.Hash()), r.Logs)
		} else {
			log.Info("Successfully executed tx %d.", obscurocommon.ShortHash(t.Hash()))
		}
	}
	s.Finalise(true)
	return receipts
}

func executeTransaction(s *state.StateDB, cc *params.ChainConfig, chain *ObscuroChainContext, gp *core2.GasPool, header *nodecommon.Header, t *nodecommon.L2Tx, usedGas *uint64, vmCfg vm.Config, tCount int) (*types.Receipt, error) {
	s.Prepare(t.Hash(), tCount)

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
func ExecuteOffChainCall(from common.Address, to common.Address, data []byte, s *state.StateDB, header *nodecommon.Header, rollupResolver db.RollupResolver, chainID int64) (*core2.ExecutionResult, error) {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)

	blockContext := core2.NewEVMBlockContext(convertToEthHeader(header), chain, &header.Agg)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, s, cc, vmCfg)

	msg := types.NewMessage(from, &to, 0, common.Big0, 100_000, common.Big0, common.Big0, common.Big0, data, nil, true)
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
		LondonBlock: common.Big0,
	}
	vmCfg := vm.Config{
		NoBaseFee: true,
	}
	gp := core2.GasPool(math.MaxUint64)
	return chain, cc, vmCfg, &gp
}
