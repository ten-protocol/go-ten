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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// These are hardcoded values necessary as an intermediary step.
// The assumption is that there is a single ERC20 which represents "The balance"
// Todo - this has to be changed to mapping of "supported ERC20 Ethereum address - Obscuro address" ( eg.: USDT address -> Obscuro WUSDT address)
// Todo - also on depositing, there has to be a minting step
var (
	Erc20OwnerKey, _     = crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")
	Erc20OwnerAddress    = crypto.PubkeyToAddress(Erc20OwnerKey.PublicKey)
	Erc20ContractAddress = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))
)

// WithdrawalAddress Custom address used for exiting Obscuro
// Todo - This should be the address of a Bridge contract.
var WithdrawalAddress = common.HexToAddress("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
func ExecuteTransactions(txs []nodecommon.L2Tx, s *state.StateDB, header *nodecommon.Header, rollupResolver db.RollupResolver, chainID int64) map[common.Hash]*types.Receipt {
	chain, cc, vmCfg, gp := initParams(rollupResolver, chainID)
	zero := uint64(0)
	usedGas := &zero
	receipts := make(map[common.Hash]*types.Receipt, len(txs))
	for _, t := range txs {
		r, err := executeTransaction(s, cc, chain, gp, header, t, usedGas, vmCfg)
		if err != nil {
			log.Info("Error transaction %d: %s", obscurocommon.ShortHash(t.Hash()), err)
			continue
		}
		receipts[t.Hash()] = r
		if r.Status != 1 {
			log.Info("Unsuccessful (status != 1) tx %d.", obscurocommon.ShortHash(t.Hash()))
		} else {
			log.Info("Successfully executed tx %d", obscurocommon.ShortHash(t.Hash()))
		}
	}
	s.Finalise(true)
	return receipts
}

func executeTransaction(s *state.StateDB, cc *params.ChainConfig, chain *ObscuroChainContext, gp *core2.GasPool, header *nodecommon.Header, t nodecommon.L2Tx, usedGas *uint64, vmCfg vm.Config) (*types.Receipt, error) {
	snap := s.Snapshot()
	// todo - Author?
	receipt, err := core2.ApplyTransaction(cc, chain, nil, gp, s, convertToEthHeader(header), &t, usedGas, vmCfg)
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
