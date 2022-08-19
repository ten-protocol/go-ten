package evm

import (
	"math"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*common.L2Tx, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig, fromTxIndex int) map[common.TxHash]interface{} {
	chain, vmCfg, gp := initParams(storage, false)
	zero := uint64(0)
	usedGas := &zero
	result := map[common.TxHash]interface{}{}

	ethHeader := convertToEthHeader(header, secret(storage))

	for i, t := range txs {
		r, err := executeTransaction(s, chainConfig, chain, gp, ethHeader, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			result[t.Hash()] = err
			common.ErrorTXExecution(t.Hash(), "Error: %s", err)
			continue
		}
		result[t.Hash()] = r
		if r.Status == types.ReceiptStatusFailed {
			common.ErrorTXExecution(t.Hash(),
				"Unsuccessful (status != 1)."+
					"\n To: %s"+
					"\n Data: %x"+
					"\n Gas: %d",
				t.To().Hex(),
				t.Data(),
				t.Gas())
		} else {
			common.TraceTXExecution(t.Hash(), "Successfully executed. Address: %s", r.ContractAddress.Hex())
		}
	}
	s.Finalise(true)
	return result
}

func executeTransaction(s *state.StateDB, cc *params.ChainConfig, chain *ObscuroChainContext, gp *gethcore.GasPool, header *types.Header, t *common.L2Tx, usedGas *uint64, vmCfg vm.Config, tCount int) (*types.Receipt, error) {
	s.Prepare(t.Hash(), tCount)
	snap := s.Snapshot()

	before := header.MixDigest
	// calculate a random value per transaction
	header.MixDigest = gethcommon.BytesToHash(crypto.PerTransactionRnd(before.Bytes(), tCount))
	// todo - Author?
	receipt, err := gethcore.ApplyTransaction(cc, chain, nil, gp, s, header, t, usedGas, vmCfg)
	header.MixDigest = before
	if err != nil {
		s.RevertToSnapshot(snap)
		return nil, err
	}

	return receipt, nil
}

// ExecuteOffChainCall - executes the "data" command against the "to" smart contract
func ExecuteOffChainCall(from gethcommon.Address, to *gethcommon.Address, data []byte, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig) (*gethcore.ExecutionResult, error) {
	chain, vmCfg, gp := initParams(storage, true)

	blockContext := gethcore.NewEVMBlockContext(convertToEthHeader(header, secret(storage)), chain, &header.Agg)
	// todo use ToMessage
	// 100_000_000_000 is just a huge number gasLimit for making sure the local tx doesn't fail with lack of gas
	msg := types.NewMessage(from, to, 0, gethcommon.Big0, 100_000_000_000, gethcommon.Big0, gethcommon.Big0, gethcommon.Big0, data, nil, true)

	// sets Tx.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, s, chainConfig, vmCfg)

	result, err := gethcore.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	// todo - find out why this was called since it's not being called in geth
	// s.Finalise(true)
	return result, nil
}

func initParams(storage db.Storage, noBaseFee bool) (*ObscuroChainContext, vm.Config, *gethcore.GasPool) {
	chain := &ObscuroChainContext{storage: storage}
	vmCfg := vm.Config{NoBaseFee: noBaseFee}
	gp := gethcore.GasPool(math.MaxUint64)
	return chain, vmCfg, &gp
}

// Todo - this is currently just returning the shared secret
// it should not use it directly, but derive some entropy from it
func secret(storage db.Storage) []byte {
	secret := storage.FetchSecret()
	return secret[:]
}
