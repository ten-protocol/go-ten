package evm

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	core2 "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/obscuro-playground/contracts"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.

// These are hardcoded values necessary as an intermediary step
var Erc20OwnerKey, _ = crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")
var Erc20OwnerAddress = crypto.PubkeyToAddress(Erc20OwnerKey.PublicKey)

var Erc20ContractAddress = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))

// WithdrawalAddress Custom address used for exiting Obscuro
// Todo - This should be the address of a Bridge contract.
var WithdrawalAddress = common.HexToAddress("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
func ExecuteTransactions(txs []nodecommon.L2Tx, s *state.StateDB, header *nodecommon.Header, rollupResolver db.RollupResolver) map[common.Hash]*types.Receipt {
	chain := &ObscuroChainContext{rollupResolver: rollupResolver}

	cc := &params.ChainConfig{
		ChainID:     big.NewInt(ChainID),
		LondonBlock: big.NewInt(0),
	}
	vmCfg := vm.Config{
		NoBaseFee: true,
	}
	usedGas := uint64(0)
	gp := core2.GasPool(1_000_000_000_000)
	receipts := make(map[common.Hash]*types.Receipt, len(txs))
	for i := range txs {
		t := txs[i]
		snap := s.Snapshot()

		receipt, err := core2.ApplyTransaction(cc, chain, nil, &gp, s, convertToEthHeader(header), &t, &usedGas, vmCfg)
		if err != nil {
			fmt.Printf("Failed to exec tx %d: %s\n", obscurocommon.ShortHash(t.Hash()), err)
			s.RevertToSnapshot(snap)
		} else {
			//if len(t.Data()) > 1000 {
			//	fmt.Printf(">>Status: %d.  %s\n", receipt.Status, receipt.ContractAddress.Hex())
			//}
			//log.Log(fmt.Sprintf("Executed tx %d: %d", obscurocommon.ShortHash(t.Hash()), receipt.Status))
			if receipt.Status != 1 {
				fmt.Printf("Failed tx %d. Receipt: %+v\n", obscurocommon.ShortHash(t.Hash()), receipt)
			}
			receipts[t.Hash()] = receipt
		}
	}
	return receipts
}

// Used in tests
// todo - create a generic version
func BalanceOfErc20(s vm.StateDB, address common.Address, header *nodecommon.Header, rollupResolver db.RollupResolver) uint64 {
	chain := &ObscuroChainContext{rollupResolver: rollupResolver}

	cc := &params.ChainConfig{
		ChainID:     big.NewInt(ChainID),
		LondonBlock: big.NewInt(0),
	}
	vmCfg := vm.Config{
		NoBaseFee: true,
	}
	// usedGas := uint64(0)
	gp := new(core2.GasPool).AddGas(math.MaxUint64)

	blockContext := core2.NewEVMBlockContext(convertToEthHeader(header), chain, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, s, cc, vmCfg)

	balanceData, err := contracts.PedroERC20ContractABIJSON.Pack("balanceOf", address)
	if err != nil {
		panic(err)
	}
	//tx := types.NewTx(&types.LegacyTx{
	//	Value:    common.Big0,
	//	Gas:      1_000_000,
	//	GasPrice: common.Big0,
	//	Data:     balanceData,
	//	To:       &ERC20_ADDRESS,
	//})

	msg := types.NewMessage(address, &Erc20ContractAddress, 0, common.Big0, 100_000, common.Big0, common.Big0, common.Big0, balanceData, nil, true)
	result, err := core2.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		fmt.Printf("Balance err: %s\n", err)
	} else {
		if result.Failed() {
			fmt.Printf("result falied: %v\n", result)
			return 0
		}
		r := new(big.Int)
		r = r.SetBytes(result.ReturnData)
		fmt.Printf("result success: %d\n", r.Uint64())
		return r.Uint64()
	}

	return 0
}
