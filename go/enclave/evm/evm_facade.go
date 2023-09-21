package evm

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(
	txs []*common.L2Tx,
	s *state.StateDB,
	header *common.BatchHeader,
	storage storage.Storage,
	chainConfig *params.ChainConfig,
	fromTxIndex int,
	noBaseFee bool,
	logger gethlog.Logger,
) map[common.TxHash]interface{} {
	chain, vmCfg, gp := initParams(storage, noBaseFee, logger)
	zero := uint64(0)
	usedGas := &zero
	result := map[common.TxHash]interface{}{}

	ethHeader, err := gethencoding.CreateEthHeaderForBatch(header, secret(storage))
	if err != nil {
		logger.Crit("Could not convert to eth header", log.ErrKey, err)
		return nil
	}

	hash := header.Hash()
	for i, t := range txs {
		r, err := executeTransaction(
			s,
			chainConfig,
			chain,
			gp,
			ethHeader,
			t,
			usedGas,
			vmCfg,
			fromTxIndex+i,
			hash,
			header.Number.Uint64(),
		)
		if err != nil {
			result[t.Hash()] = err
			logger.Info("Failed to execute tx:", log.TxKey, t.Hash(), log.CtrErrKey, err)
			continue
		}
		result[t.Hash()] = r
		logReceipt(r, logger)
	}
	s.Finalise(true)
	return result
}

func executeTransaction(
	s *state.StateDB,
	cc *params.ChainConfig,
	chain *ObscuroChainContext,
	gp *gethcore.GasPool,
	header *types.Header,
	t *common.L2Tx,
	usedGas *uint64,
	vmCfg vm.Config,
	tCount int,
	batchHash common.L2BatchHash,
	batchHeight uint64,
) (*types.Receipt, error) {
	rules := cc.Rules(big.NewInt(0), true, 0)
	from, err := types.Sender(types.LatestSigner(cc), t)
	if err != nil {
		return nil, err
	}
	s.Prepare(rules, from, gethcommon.Address{}, t.To(), nil, nil)
	snap := s.Snapshot()
	s.SetTxContext(t.Hash(), tCount)

	before := header.MixDigest
	// calculate a random value per transaction
	header.MixDigest = crypto.CalculateTxRnd(before.Bytes(), tCount)
	receipt, err := gethcore.ApplyTransaction(cc, chain, nil, gp, s, header, t, usedGas, vmCfg)

	// adjust the receipt to point to the right batch hash
	if receipt != nil {
		receipt.Logs = s.GetLogs(t.Hash(), batchHeight, batchHash)
		receipt.BlockHash = batchHash
		receipt.BlockNumber = big.NewInt(int64(batchHeight))
		for _, l := range receipt.Logs {
			l.BlockHash = batchHash
		}
	}

	header.MixDigest = before
	if err != nil {
		s.RevertToSnapshot(snap)
		return receipt, err
	}

	return receipt, nil
}

func logReceipt(r *types.Receipt, logger gethlog.Logger) {
	logger.Trace("Receipt", log.TxKey, r.TxHash, "Result", gethlog.Lazy{Fn: func() string {
		receiptJSON, err := r.MarshalJSON()
		if err != nil {
			if r.Status == types.ReceiptStatusFailed {
				return "Unsuccessful (status != 1) (but could not print receipt as JSON)"
			}
			return "Successfully executed (but could not print receipt as JSON)"
		}
		if r.Status == types.ReceiptStatusFailed {
			return fmt.Sprintf("Unsuccessful (status != 1). Receipt: %s", string(receiptJSON))
		}
		return fmt.Sprintf("Successfully executed. Receipt: %s", string(receiptJSON))
	}})
}

// ExecuteObsCall - executes the eth_call call
func ExecuteObsCall(
	msg *gethcore.Message,
	s *state.StateDB,
	header *common.BatchHeader,
	storage storage.Storage,
	chainConfig *params.ChainConfig,
	logger gethlog.Logger,
) (*gethcore.ExecutionResult, error) {
	noBaseFee := true
	if header.BaseFee != nil && header.BaseFee.Cmp(gethcommon.Big0) != 0 && msg.GasPrice.Cmp(gethcommon.Big0) != 0 {
		noBaseFee = false
		logger.Info("ObsCall - with base fee ", "to", msg.To.Hex())
	}

	chain, vmCfg, gp := initParams(storage, noBaseFee, nil)
	ethHeader, err := gethencoding.CreateEthHeaderForBatch(header, secret(storage))
	if err != nil {
		return nil, err
	}
	blockContext := gethcore.NewEVMBlockContext(ethHeader, chain, nil)

	// sets TxKey.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, s, chainConfig, vmCfg)

	result, err := gethcore.ApplyMessage(vmenv, msg, gp)
	// Follow the same error check structure as in geth
	// 1 - vmError / stateDB err check
	// 2 - evm.Cancelled()  todo (#1576) - support the ability to cancel function call if it takes too long
	// 3 - error check the ApplyMessage

	// Read the error stored in the database.
	if dbErr := s.Error(); dbErr != nil {
		return nil, newErrorWithReasonAndCode(dbErr)
	}

	// If the result contains a revert reason, try to unpack and return it.
	if result != nil && len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}

	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		logger.Debug(fmt.Sprintf("Error applying msg %v:", msg), log.CtrErrKey, err)
		return result, err
	}

	logger.Info("ObsCall - with result ", "gas", result.UsedGas)
	return result, nil
}

func initParams(storage storage.Storage, noBaseFee bool, l gethlog.Logger) (*ObscuroChainContext, vm.Config, *gethcore.GasPool) {
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
	}
	gp := gethcore.GasPool(math.MaxUint64)
	return NewObscuroChainContext(storage, l), vmCfg, &gp
}

// todo (#1053) - this is currently just returning the shared secret
// it should not use it directly, but derive some entropy from it
func secret(storage storage.Storage) []byte {
	// todo (#1053) - handle secret not being found.
	secret, _ := storage.FetchSecret()
	return secret[:]
}

func newErrorWithReasonAndCode(err error) error {
	result := &errutil.EVMSerialisableError{
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

func newRevertError(result *gethcore.ExecutionResult) error {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &errutil.EVMSerialisableError{
		Err:    err.Error(),
		Reason: hexutil.Encode(result.Revert()),
		Code:   3, // todo - magic number, really needs thought around the value and made a constant
	}
}
