package evm

import (
	"errors"
	"fmt"
	"math/big"
	_ "unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(
	txs common.L2PricedTransactions,
	s *state.StateDB,
	header *common.BatchHeader,
	storage storage.Storage,
	chainConfig *params.ChainConfig,
	fromTxIndex int,
	noBaseFee bool,
	batchGasLimit uint64,
	logger gethlog.Logger,
) map[common.TxHash]interface{} {
	chain, vmCfg := initParams(storage, noBaseFee, logger)
	gp := gethcore.GasPool(batchGasLimit)
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
			&gp,
			ethHeader,
			t,
			usedGas,
			vmCfg,
			fromTxIndex+i,
			hash,
			header.Number.Uint64(),
		)
		if err != nil {
			result[t.Tx.Hash()] = err
			logger.Info("Failed to execute tx:", log.TxKey, t.Tx.Hash(), log.CtrErrKey, err)
			continue
		}
		result[t.Tx.Hash()] = r
		logReceipt(r, logger)
	}
	s.Finalise(true)
	return result
}

//go:linkname applyTransaction github.com/ethereum/go-ethereum/core.applyTransaction
func applyTransaction(msg *gethcore.Message, config *params.ChainConfig, gp *gethcore.GasPool, statedb *state.StateDB, blockNumber *big.Int, blockHash gethcommon.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, error)

func executeTransaction(
	s *state.StateDB,
	cc *params.ChainConfig,
	chain *ObscuroChainContext,
	gp *gethcore.GasPool,
	header *types.Header,
	t common.L2PricedTransaction,
	usedGas *uint64,
	vmCfg vm.Config,
	tCount int,
	batchHash common.L2BatchHash,
	batchHeight uint64,
) (*types.Receipt, error) {
	rules := cc.Rules(big.NewInt(0), true, 0)
	from, err := types.Sender(types.LatestSigner(cc), t.Tx)
	if err != nil {
		return nil, err
	}
	s.Prepare(rules, from, gethcommon.Address{}, t.Tx.To(), nil, nil)
	snap := s.Snapshot()
	s.SetTxContext(t.Tx.Hash(), tCount)

	before := header.MixDigest
	// calculate a random value per transaction
	header.MixDigest = crypto.CalculateTxRnd(before.Bytes(), tCount)

	applyTx := func(
		config *params.ChainConfig,
		bc gethcore.ChainContext,
		author *gethcommon.Address,
		gp *gethcore.GasPool,
		statedb *state.StateDB,
		header *types.Header,
		tx common.L2PricedTransaction,
		usedGas *uint64,
		cfg vm.Config,
	) (*types.Receipt, error) {
		msg, err := gethcore.TransactionToMessage(tx.Tx, types.MakeSigner(config, header.Number, header.Time), header.BaseFee)
		if err != nil {
			return nil, err
		}
		l1cost := tx.PublishingCost
		l1Gas := big.NewInt(0)
		hasL1Cost := l1cost.Cmp(big.NewInt(0)) != 0

		if hasL1Cost {
			l1Gas.Div(l1cost, header.BaseFee)
			l1Gas.Add(l1Gas, big.NewInt(0).Mod(l1cost, header.BaseFee))

			if msg.GasLimit < l1Gas.Uint64() {
				return nil, fmt.Errorf("gas limit for tx too low. Want at least: %d have: %d", l1Gas, msg.GasLimit)
			}
			msg.GasLimit -= l1Gas.Uint64()

			statedb.SubBalance(msg.From, l1cost)
			statedb.AddBalance(header.Coinbase, l1cost)

		}

		// Create a new context to be used in the EVM environment
		blockContext := gethcore.NewEVMBlockContext(header, bc, author)
		vmenv := vm.NewEVM(blockContext, vm.TxContext{BlobHashes: tx.Tx.BlobHashes()}, statedb, config, cfg)
		receipt, err := applyTransaction(msg, config, gp, statedb, header.Number, header.Hash(), tx.Tx, usedGas, vmenv)
		if err != nil {
			if hasL1Cost {
				statedb.SubBalance(header.Coinbase, l1cost)
				statedb.AddBalance(msg.From, l1cost)
			}

			return receipt, err
		}

		// Do not increase the balance of zero address as it is the contract deployment address.
		// Doing so might cause weird interactions.
		if header.Coinbase.Big().Cmp(gethcommon.Big0) != 0 {
			gasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
			executionGasCost := big.NewInt(0).Mul(gasUsed, header.BaseFee)
			// As the baseFee is burned, we add it back to the coinbase.
			// Geth should automatically add the tips.
			statedb.AddBalance(header.Coinbase, executionGasCost)
		}
		receipt.GasUsed += l1Gas.Uint64()

		return receipt, err
	}

	receipt, err := applyTx(cc, chain, nil, gp, s, header, t, usedGas, vmCfg)

	// adjust the receipt to point to the right batch hash
	if receipt != nil {
		receipt.Logs = s.GetLogs(t.Tx.Hash(), batchHeight, batchHash)
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
	gasEstimationCap uint64,
	logger gethlog.Logger,
) (*gethcore.ExecutionResult, error) {
	noBaseFee := true
	if header.BaseFee != nil && header.BaseFee.Cmp(gethcommon.Big0) != 0 && msg.GasPrice.Cmp(gethcommon.Big0) != 0 {
		noBaseFee = false
	}

	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "evm_facade.go:ObsCall()")

	gp := gethcore.GasPool(gasEstimationCap)
	gp.SetGas(gasEstimationCap)
	chain, vmCfg := initParams(storage, noBaseFee, nil)
	ethHeader, err := gethencoding.CreateEthHeaderForBatch(header, secret(storage))
	if err != nil {
		return nil, err
	}
	blockContext := gethcore.NewEVMBlockContext(ethHeader, chain, nil)

	// sets TxKey.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, s, chainConfig, vmCfg)

	result, err := gethcore.ApplyMessage(vmenv, msg, &gp)
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

	return result, nil
}

func initParams(storage storage.Storage, noBaseFee bool, l gethlog.Logger) (*ObscuroChainContext, vm.Config) {
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
	}
	return NewObscuroChainContext(storage, l), vmCfg
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
