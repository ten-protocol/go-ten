package enclave

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

var (
	singleCoinERC20ContractAddress = common.HexToAddress("0x1")
	withdrawalAddress              = common.HexToAddress("0x2")
)

// mutates the State
// header - the header of the rollup where this transaction will be included
func executeTransactions(txs []nodecommon.L2Tx, s vm.StateDB, header *nodecommon.Header) {
	for _, tx := range txs {
		executeTx(s, tx, header)
	}
	// fmt.Printf("w1: %v\n", is.w)
}

// mutates the state
func executeTx(s vm.StateDB, tx nodecommon.L2Tx, header *nodecommon.Header) {
	// snap := s.Snapshot()
	// log.Log(fmt.Sprintf("snap: %d", snap))
	txType := core.TxData(&tx).Type
	switch txType {
	case core.TransferTx:
		executeTransfer(s, tx)
	case core.WithdrawalTx:
		executeWithdrawal(s, tx, header)
	case core.DepositTx:
		executeDeposit(s, tx)
	default:
		msg := fmt.Sprintf("attempted to execute unrecognised transaction type %s", reflect.TypeOf(txType))
		log.Error(msg)
		panic(msg)
	}
}

func addrToHash(a common.Address) common.Hash {
	return common.BytesToHash(a.Bytes())
}

func getBalance(s vm.StateDB, address common.Address) uint64 {
	// This is necessary so that the account actually gets stored
	s.SetNonce(singleCoinERC20ContractAddress, 1)

	result := s.GetState(singleCoinERC20ContractAddress, addrToHash(address))
	return binary.BigEndian.Uint64(result.Bytes()[24:32])
}

func setBalance(s vm.StateDB, address common.Address, balance uint64) {
	bal := make([]byte, 8)
	binary.BigEndian.PutUint64(bal, balance)
	h := common.BytesToHash(bal)
	s.SetState(singleCoinERC20ContractAddress, addrToHash(address), h)
}

func addWithdrawal(s vm.StateDB, _ common.Hash, hash obscurocommon.TxHash) {
	// This is necessary so that the account actually gets stored
	s.SetNonce(withdrawalAddress, 1)

	s.SetState(withdrawalAddress, hash, hash)
}

func withdrawals(s vm.StateDB, _ common.Hash) []obscurocommon.TxHash {
	var withdrawalTxs []obscurocommon.TxHash
	err := s.ForEachStorage(withdrawalAddress, func(k common.Hash, v common.Hash) bool {
		if (v != common.Hash{}) {
			withdrawalTxs = append(withdrawalTxs, v)
		}
		return true
	})
	if err != nil {
		log.Error(fmt.Sprintf("could not retrieve withdrawals. Cause: %s", err))
		panic(err)
	}
	// fmt.Printf(">>>withd %d\n", len(withdrawalTxs))
	return withdrawalTxs
}

//func clearWithdrawals(s vm.StateDB, w []obscurocommon.TxHash) {
//	for _, hash := range w {
//		s.SetState(withdrawalAddress, hash, common.Hash{})
//	}
//}

func executeWithdrawal(s vm.StateDB, tx nodecommon.L2Tx, header *nodecommon.Header) {
	txData := core.TxData(&tx)
	from := getBalance(s, txData.From)
	if from >= txData.Amount {
		setBalance(s, txData.From, from-txData.Amount)
		addWithdrawal(s, header.Hash(), tx.Hash())
	}
}

func executeTransfer(s vm.StateDB, tx nodecommon.L2Tx) {
	txData := core.TxData(&tx)
	from := getBalance(s, txData.From)
	to := getBalance(s, txData.To)

	if from >= txData.Amount {
		setBalance(s, txData.From, from-txData.Amount)
		setBalance(s, txData.To, to+txData.Amount)
	}
}

func executeDeposit(s vm.StateDB, tx nodecommon.L2Tx) {
	txData := core.TxData(&tx)
	to := getBalance(s, txData.To)
	log.Info(fmt.Sprintf("Tx=%d; Process deposit %d into %s(old=%d).", obscurocommon.ShortHash(tx.Hash()), txData.Amount, txData.To, to))
	setBalance(s, txData.To, to+txData.Amount)
}
