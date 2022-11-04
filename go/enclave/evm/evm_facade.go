package evm

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"math"

	gethlog "github.com/ethereum/go-ethereum/log"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*common.L2Tx, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig, fromTxIndex int, logger gethlog.Logger) map[common.TxHash]interface{} {
	chain, vmCfg, gp := initParams(storage, false)
	zero := uint64(0)
	usedGas := &zero
	result := map[common.TxHash]interface{}{}

	ethHeader := convertToEthHeader(header, secret(storage))

	for i, t := range txs {
		r, err := executeTransaction(s, chainConfig, chain, gp, ethHeader, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			result[t.Hash()] = err
			logger.Error("!TxKey", log.TxKey, t.Hash().Hex(), log.ErrKey, err)
			continue
		}
		result[t.Hash()] = r
		logReceipt(r, logger)
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
	receipt, err := gethcore.ApplyTransaction(cc, chain, nil, gp, s, header, t, usedGas, vmCfg)
	header.MixDigest = before
	if err != nil {
		s.RevertToSnapshot(snap)
		return nil, err
	}

	return receipt, nil
}

func logReceipt(r *types.Receipt, logger gethlog.Logger) {
	receiptJSON, err := r.MarshalJSON()
	if err != nil {
		if r.Status == types.ReceiptStatusFailed {
			logger.Error("Unsuccessful (status != 1) (but could not print receipt as JSON)", log.TxKey, r.TxHash.Hex())
		} else {
			logger.Trace("Successfully executed (but could not print receipt as JSON)", log.TxKey, r.TxHash.Hex())
		}
	}

	if r.Status == types.ReceiptStatusFailed {
		logger.Error(fmt.Sprintf("Unsuccessful (status != 1). Receipt: %s", string(receiptJSON)), log.TxKey, r.TxHash.Hex())
	} else {
		logger.Trace(fmt.Sprintf("Successfully executed. Receipt: %s", string(receiptJSON)), log.TxKey, r.TxHash.Hex())
	}
}

// ExecuteOffChainCall - executes the "data" command against the "to" smart contract
func ExecuteOffChainCall(call *ethereum.CallMsg, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig, logger gethlog.Logger) (*gethcore.ExecutionResult, error) {
	chain, vmCfg, gp := initParams(storage, true)

	blockContext := gethcore.NewEVMBlockContext(convertToEthHeader(header, secret(storage)), chain, &header.Agg)
	// todo use ToMessage
	// call values can't be nil or they will fail downstream
	if call.Value == nil {
		call.Value = gethcommon.Big0
	}
	msg := types.NewMessage(call.From, call.To, 0, call.Value, call.Gas, gethcommon.Big0, gethcommon.Big0, gethcommon.Big0, call.Data, nil, true)

	// sets TxKey.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, s, chainConfig, vmCfg)

	result, err := gethcore.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		logger.Error("ErrKey applying msg:", log.ErrKey, err)
		return nil, err
	}

	// Read the error stored in the database.
	err = s.Error()
	if err != nil {
		return nil, newErrorWithReasonAndCode(err)
	}

	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}
	return result, nil
}

func initParams(storage db.Storage, noBaseFee bool) (*ObscuroChainContext, vm.Config, *gethcore.GasPool) {
	chain := &ObscuroChainContext{storage: storage}

	// Todo - temporarily enable the evm tracer to check what sort of extra info we receive
	tracer := logger.NewStructLogger(&logger.Config{Debug: true})
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
		Debug:     false,
		Tracer:    tracer,
	}
	gp := gethcore.GasPool(math.MaxUint64)
	return chain, vmCfg, &gp
}

// Todo - this is currently just returning the shared secret
// it should not use it directly, but derive some entropy from it
func secret(storage db.Storage) []byte {
	secret := storage.FetchSecret()
	return secret[:]
}

func newErrorWithReasonAndCode(err error) SerialisableError {
	result := SerialisableError{
		Err: err.Error(),
	}

	var e gethrpc.Error
	ok := errors.As(err, &e)
	if ok {
		result.Code = e.ErrorCode()
	}
	var de gethrpc.DataError
	ok = errors.As(err, &de)
	if ok {
		result.Reason = de.ErrorData()
	}
	return result
}

func newRevertError(result *gethcore.ExecutionResult) SerialisableError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return SerialisableError{
		Err:    err.Error(),
		Reason: hexutil.Encode(result.Revert()),
		Code:   3, // todo - magic number
	}
}

// SerialisableError is an API error that encompasses an EVM error with a code and a reason
type SerialisableError struct {
	Err    string
	Reason interface{}
	Code   int
}

func (e SerialisableError) Error() string {
	return e.Err
}

func (e SerialisableError) ErrorCode() int {
	return e.Code
}

func (e SerialisableError) ErrorData() interface{} {
	return e.Reason
}
