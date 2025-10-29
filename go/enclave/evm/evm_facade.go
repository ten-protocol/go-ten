package evm

// unsafe package imported in order to link to a private function in go-ethereum.
// This allows us to customize the message generated from a signed transaction and inject custom gas logic.
import (
	"context"
	"fmt"
	"math/big"
	_ "unsafe"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
)

type evmExecutor struct {
	chain            gethcore.ChainContext
	cc               *params.ChainConfig
	config           *enclaveconfig.EnclaveConfig
	gasEstimationCap uint64

	storage             storage.Storage
	gethEncodingService gethencoding.EncodingService
	visibilityReader    ContractVisibilityReader

	logger gethlog.Logger
}

func NewEVMExecutor(chain gethcore.ChainContext, cc *params.ChainConfig, config *enclaveconfig.EnclaveConfig, gasEstimationCap uint64, storage storage.Storage, gethEncodingService gethencoding.EncodingService, visibilityReader ContractVisibilityReader, logger gethlog.Logger) *evmExecutor {
	return &evmExecutor{
		chain:               chain,
		cc:                  cc,
		config:              config,
		gasEstimationCap:    gasEstimationCap,
		storage:             storage,
		gethEncodingService: gethEncodingService,
		visibilityReader:    visibilityReader,
		logger:              logger,
	}
}

func (exec *evmExecutor) ExecuteTx(tx *common.L2PricedTransaction, s *state.StateDB, header *types.Header, gp *gethcore.GasPool, usedGas *uint64, tCount int, noBaseFee bool) *core.TxExecResult {
	from, err := core.GetTxSigner(tx)
	if err != nil {
		return &core.TxExecResult{
			TxWithSender: &core.TxWithSender{Tx: tx.Tx, Sender: &from},
			Err:          err,
		}
	}

	snap := s.Snapshot()
	res, err := exec.execute(tx, from, s, header, gp, usedGas, tCount, noBaseFee)
	if err != nil {
		s.RevertToSnapshot(snap)
		return &core.TxExecResult{
			TxWithSender: &core.TxWithSender{Tx: tx.Tx, Sender: &from},
			Err:          err,
		}
	}
	return res
}

func (exec *evmExecutor) execute(tx *common.L2PricedTransaction, from gethcommon.Address, s *state.StateDB, header *types.Header, gp *gethcore.GasPool, usedGas *uint64, tCount int, noBaseFee bool) (*core.TxExecResult, error) {
	// a transaction can create multiple contracts.
	// we use a tracer hook to collect the addresses
	var createdContracts []*gethcommon.Address
	cfg := vm.Config{
		NoBaseFee: noBaseFee,
		Tracer: &tracing.Hooks{
			// called when the code of a contract changes.
			OnCodeChange: func(addr gethcommon.Address, prevCodeHash gethcommon.Hash, prevCode []byte, codeHash gethcommon.Hash, code []byte) {
				// only proceed for new deployments.
				if len(prevCode) > 0 {
					exec.logger.Debug("OnCodeChange: Skipping contract deployment", "address", addr.Hex())
					return
				}
				a := addr
				createdContracts = append(createdContracts, &a)
				exec.logger.Debug("OnCodeChange: Contract deployed", "address", addr.Hex())
			},
		},
	}

	hookedStateDb := state.NewHookedState(s, cfg.Tracer)

	blockContext := gethcore.NewEVMBlockContext(header, exec.chain, nil)
	evmEnv := vm.NewEVM(blockContext, hookedStateDb, exec.cc, cfg)

	msg, err := TransactionToMessageWithOverrides(tx, exec.cc, header)
	if err != nil {
		return nil, err
	}

	actualGasLimit := msg.GasLimit // before l1 cost is applied

	receipt, err := adjustPublishingCostGas(tx, msg, s, header, noBaseFee, func() (*types.Receipt, error) {
		s.SetTxContext(tx.Tx.Hash(), tCount)

		// the gas limit should never be higher than the max tx gas
		msg.GasLimit = min(msg.GasLimit, params.MaxTxGas-1)

		return gethcore.ApplyTransactionWithEVM(msg, gp, s, header.Number, header.Hash(), header.Time, tx.Tx, usedGas, evmEnv)
	})
	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, fmt.Errorf("internal: nil receipt returned from execution")
	}
	receipt.Logs = s.GetLogs(tx.Tx.Hash(), header.Number.Uint64(), header.Hash(), header.Time)

	contractsWithVisibility := make(map[gethcommon.Address]*core.ContractVisibilityConfig)

	// Compute leftover user gas after the main execution
	gasLeft := actualGasLimit
	if receipt.GasUsed > actualGasLimit {
		// ensure no bugs going overboard with the l1 publishing
		// or execution before we pick up the leftover gas.
		// SHOULD NOT HAPPEN  unless we have a bug in the code.
		// If this happens we bail otherwise we overcharge the user above the authorized max amount.
		return nil, fmt.Errorf("internal: gasUsed (%d) exceeds tx gasLimit (%d)", receipt.GasUsed, gasLeft)
	}
	gasLeft -= receipt.GasUsed

	for _, contractAddress := range createdContracts {
		if gasLeft == 0 {
			return nil, fmt.Errorf("out of gas while reading visibility for %s", contractAddress.Hex())
		}

		cfg, visUsed, err1 := exec.readVisibilityWithCap(context.Background(), evmEnv, gp, *contractAddress, gasLeft)
		if err1 != nil {
			exec.logger.Crit("metered visibility read failed", log.ErrKey, err1, "addr", contractAddress.Hex())
			return nil, fmt.Errorf("visibility read failed for %s: %w", contractAddress.Hex(), err1)
		}
		if noBaseFee {
			visUsed = 0
		}

		// Ensure we still respect the user's tx gas limit
		// this should never happen as the limit to go up to is the gasLeft.
		// thus if we hit the error we have a bug in the code and we error out of the transaction.
		// Otherwise we might end up overcharging a user - imagine signing a 0.1$ fee transaction
		// but we end up charging you 500$ because we do not respect the limit for whatever reason.
		if visUsed > gasLeft {
			return nil, fmt.Errorf("out of gas: visibility read used %d, leftover %d, addr %s", visUsed, gasLeft, contractAddress.Hex())
		}
		gasLeft -= visUsed

		// Fold into tx accounting so it looks like part of the tx
		receipt.GasUsed += visUsed
		if usedGas != nil {
			*usedGas += visUsed
		}

		// Charge baseFee for this portion only (no tip)
		if header != nil {
			// baseFee may be nil/zero in some test configs; the helper guards for that.
			exec.chargeBaseFeeOnly(s, header, from, visUsed)
		}

		contractsWithVisibility[*contractAddress] = cfg
	}

	return &core.TxExecResult{
		Receipt:          receipt,
		TxWithSender:     &core.TxWithSender{Tx: tx.Tx, Sender: &from},
		CreatedContracts: contractsWithVisibility,
	}, nil
}

// ExecuteCall - executes the eth_call call
func (exec *evmExecutor) ExecuteCall(ctx context.Context, msg *gethcore.Message, s *state.StateDB, header *common.BatchHeader, isEstimateGas bool) (*gethcore.ExecutionResult, error, common.SystemError) {
	defer core.LogMethodDuration(exec.logger, measure.NewStopwatch(), "evm_facade.go:Call()")

	var initslots []gethcommon.Hash
	if msg.To != nil {
		reader, err := s.Database().Reader(header.Root)
		if err != nil {
			exec.logger.Error("evmf: could not get state reader", log.ErrKey, err)
			return nil, nil, nil
		}

		var i = big.NewInt(0)
		for {
			k := gethcommon.Hash{}
			k.SetBytes(i.Bytes())

			slot, err := reader.Storage(*msg.To, k)
			if err != nil {
				exec.logger.Error("evmf: could not get account", log.ErrKey, err)
				return nil, nil, nil
			}
			if slot == (gethcommon.Hash{}) {
				break
			}
			initslots = append(initslots, slot)
			i = i.Add(i, big.NewInt(1))
		}
	}

	vmCfg := vm.Config{
		NoBaseFee: true,
	}

	// a call can create multiple contracts; capture via tracer when estimating
	var createdContracts []*gethcommon.Address
	if isEstimateGas {
		vmCfg.Tracer = &tracing.Hooks{
			OnCodeChange: func(addr gethcommon.Address, prevCodeHash gethcommon.Hash, prevCode []byte, codeHash gethcommon.Hash, code []byte) {
				if len(prevCode) > 0 {
					return
				}
				a := addr
				createdContracts = append(createdContracts, &a)
			},
		}
	}

	ethHeader, err := exec.gethEncodingService.CreateEthHeaderForBatch(ctx, header)
	if err != nil {
		return nil, nil, fmt.Errorf("evmf: could not convert to eth header: %w", err)
	}

	gp := gethcore.GasPool(exec.gasEstimationCap)
	gp.SetGas(exec.gasEstimationCap)

	cleanState := createCleanState(s, msg, ethHeader, exec.cc)

	blockContext := gethcore.NewEVMBlockContext(ethHeader, exec.chain, nil)
	// Use hooked state so tracer fires during estimation; otherwise use clean state
	var evmState vm.StateDB
	if vmCfg.Tracer != nil {
		evmState = state.NewHookedState(cleanState, vmCfg.Tracer)
	} else {
		evmState = cleanState
	}
	// sets TxKey.origin
	vmenv := vm.NewEVM(blockContext, evmState, exec.cc, vmCfg)

	// Monitor the outer context and interrupt the EVM upon cancellation. To avoid
	// a dangling goroutine until the outer estimation finishes, create an internal
	// context for the lifetime of this method call.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		vmenv.Cancel()
	}()

	result, err := gethcore.ApplyMessage(vmenv, msg, &gp)
	// Read the error stored in the database.
	if vmerr := cleanState.Error(); vmerr != nil {
		return nil, vmerr, nil
	}

	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		exec.logger.Debug(fmt.Sprintf("Error applying msg %v:", msg), log.CtrErrKey, err)
		return result, fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit), nil
	}

	// Meter visibility reads for newly created contracts (estimation)
	if isEstimateGas && result != nil && len(createdContracts) > 0 {
		// leftover from the user's gas limit
		var gasLeft uint64 = msg.GasLimit
		if result.UsedGas > gasLeft {
			return result, fmt.Errorf("internal: estimate gasUsed exceeds limit"), nil
		}
		gasLeft -= result.UsedGas

		var extra uint64
		for _, ca := range createdContracts {
			if gasLeft == 0 {
				return result, fmt.Errorf("out of gas while estimating visibility read for %s", ca.Hex()), nil
			}
			_, visUsed, err := exec.readVisibilityWithCap(context.Background(), vmenv, &gp, *ca, gasLeft)
			if err != nil {
				return result, fmt.Errorf("visibility read (estimate) failed for %s: %w", ca.Hex(), err), nil
			}

			if visUsed > gasLeft {
				return result, fmt.Errorf("out of gas after visibility read (estimate) %s", ca.Hex()), nil
			}
			gasLeft -= visUsed
			extra += visUsed
		}
		result.UsedGas += extra
		exec.logger.Debug("estimate: added visibility-read gas", "created", len(createdContracts), "extraGas", extra, "totalUsedGas", result.UsedGas)
	}

	if msg.To != nil {
		reader, err := s.Database().Reader(header.Root)
		if err != nil {
			exec.logger.Error("evmf: could not get state reader", log.ErrKey, err)
			return nil, nil, nil
		}

		var i = big.NewInt(0)
		for {
			k := gethcommon.Hash{}
			k.SetBytes(i.Bytes())

			slot, err := reader.Storage(*msg.To, k)
			if err != nil {
				exec.logger.Error("evmf: could not get account", log.ErrKey, err)
				return nil, nil, nil
			}
			if slot == (gethcommon.Hash{}) {
				break
			}
			initSlot := initslots[i.Uint64()]
			if slot != initSlot {
				exec.logger.Error("evmf: storage slot changed", "slot", i, "initSlot", initSlot, "Slot", slot, "msg.To", msg.To.Hex(), "msg.From", msg.From.Hex())
				return nil, nil, nil
			}
			i = i.Add(i, big.NewInt(1))
		}
	}

	return result, nil, nil
}

func createCleanState(s *state.StateDB, msg *gethcore.Message, ethHeader *types.Header, chainConfig *params.ChainConfig) *state.StateDB {
	cleanState := s.Copy()
	cleanState.Prepare(chainConfig.Rules(ethHeader.Number, true, 0), msg.From, ethHeader.Coinbase, msg.To, nil, msg.AccessList)
	return cleanState
}

func (exec *evmExecutor) chargeBaseFeeOnly(s *state.StateDB, header *types.Header, payer gethcommon.Address, gasUsed uint64) {
	if gasUsed == 0 || header.BaseFee == nil || header.BaseFee.Sign() == 0 {
		return
	}
	weiBig := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), header.BaseFee)
	if weiBig.Sign() <= 0 {
		return
	}
	wei, _ := uint256.FromBig(weiBig)
	s.SubBalance(payer, wei, tracing.BalanceDecreaseGasBuy)
	s.AddBalance(header.Coinbase, wei, tracing.BalanceIncreaseRewardTransactionFee)
}

// readVisibilityWithCap reserves cap gas from the GasPool, invokes the visibility reader,
// and refunds any unused portion back to the pool. On error, it refunds the full cap if reserved.
func (exec *evmExecutor) readVisibilityWithCap(ctx context.Context, evmEnv *vm.EVM, gp *gethcore.GasPool, addr gethcommon.Address, cap uint64) (*core.ContractVisibilityConfig, uint64, error) {
	if err := gp.SubGas(cap); err != nil {
		return nil, 0, err
	}
	cfg, used, err := exec.visibilityReader.ReadVisibilityConfig(ctx, evmEnv, addr, cap)
	if err != nil {
		gp.AddGas(cap) // error out tx is not in batch, do not make the gas pool different between validators
		return nil, used, err
	}
	if cap >= used {
		gp.AddGas(cap - used)
	}
	return cfg, used, nil
}
