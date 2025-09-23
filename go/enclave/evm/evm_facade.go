package evm

// unsafe package imported in order to link to a private function in go-ethereum.
// This allows us to customize the message generated from a signed transaction and inject custom gas logic.
import (
	"context"
	"fmt"
	"math/big"
	_ "unsafe"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"

	"github.com/ethereum/go-ethereum/core/tracing"
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
	chain            *TenChainContext
	cc               *params.ChainConfig
	config           *enclaveconfig.EnclaveConfig
	gasEstimationCap uint64

	storage             storage.Storage
	gethEncodingService gethencoding.EncodingService
	visibilityReader    ContractVisibilityReader

	logger gethlog.Logger
}

func NewEVMExecutor(chain *TenChainContext, cc *params.ChainConfig, config *enclaveconfig.EnclaveConfig, gasEstimationCap uint64, storage storage.Storage, gethEncodingService gethencoding.EncodingService, visibilityReader ContractVisibilityReader, logger gethlog.Logger) *evmExecutor {
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
				createdContracts = append(createdContracts, &addr)
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
		return gethcore.ApplyTransactionWithEVM(msg, gp, s, header.Number, header.Hash(), header.Time, tx.Tx, usedGas, evmEnv)
	})
	if err != nil {
		return nil, err
	}
	receipt.Logs = s.GetLogs(tx.Tx.Hash(), header.Number.Uint64(), header.Hash(), header.Time)

	contractsWithVisibility := make(map[gethcommon.Address]*core.ContractVisibilityConfig)

	// Compute leftover user gas after the main execution
	var gasLeft uint64 = actualGasLimit
	if receipt != nil {
		if receipt.GasUsed > gasLeft {
			return nil, fmt.Errorf("internal: gasUsed (%d) exceeds tx gasLimit (%d)", receipt.GasUsed, gasLeft)
		}
		gasLeft -= receipt.GasUsed
	}

	for _, contractAddress := range createdContracts {
		if gasLeft == 0 {
			return nil, fmt.Errorf("out of gas while reading visibility for %s", contractAddress.Hex())
		}

		// Cap this read by current leftover (and by maxGasForVisibility inside the reader)
		cap := gasLeft

		// Reserve block gas up-front, refund after we know how much was actually used.
		if err := gp.SubGas(cap); err != nil {
			return nil, fmt.Errorf("block gas depleted before visibility read for %s: %w", contractAddress.Hex(), err)
		}

		cfg, visUsed, err1 := exec.visibilityReader.(*contractVisibilityReader).ReadVisibilityConfigMetered(context.Background(), evmEnv, *contractAddress, cap)
		if err1 != nil {
			// Return block gas we reserved (all of it) to be safe; execution will revert anyway.
			gp.AddGas(cap)
			exec.logger.Crit("metered visibility read failed", log.ErrKey, err1, "addr", contractAddress.Hex())
			return nil, err1
		}

		// Refund any unused portion to the block gas pool
		if cap >= visUsed {
			gp.AddGas(cap - visUsed)
		}

		// Ensure we still respect the user's tx gas limit
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
func (exec *evmExecutor) ExecuteCall(ctx context.Context, msg *gethcore.Message, s *state.StateDB, header *common.BatchHeader) (*gethcore.ExecutionResult, error, common.SystemError) {
	defer core.LogMethodDuration(exec.logger, measure.NewStopwatch(), "evm_facade.go:Call()")

	vmCfg := vm.Config{
		NoBaseFee: true,
	}

	// a call can create multiple contracts; capture via tracer
	var createdContracts []*gethcommon.Address
	vmCfg.Tracer = &tracing.Hooks{
		OnCodeChange: func(addr gethcommon.Address, prevCodeHash gethcommon.Hash, prevCode []byte, codeHash gethcommon.Hash, code []byte) {
			if len(prevCode) > 0 {
				return
			}
			a := addr
			createdContracts = append(createdContracts, &a)
		},
	}

	ethHeader, err := exec.gethEncodingService.CreateEthHeaderForBatch(ctx, header)
	if err != nil {
		return nil, nil, fmt.Errorf("evmf: could not convert to eth header: %w", err)
	}

	gp := gethcore.GasPool(exec.gasEstimationCap)
	gp.SetGas(exec.gasEstimationCap)

	cleanState := createCleanState(s, msg, ethHeader, exec.cc)
	snapshot := cleanState.Snapshot()
	defer cleanState.RevertToSnapshot(snapshot) // Always revert after simulation

	blockContext := gethcore.NewEVMBlockContext(ethHeader, exec.chain, nil)
	// Use hooked state so tracer fires during estimation as well
	hookedState := state.NewHookedState(cleanState, vmCfg.Tracer)
	// sets TxKey.origin
	vmenv := vm.NewEVM(blockContext, hookedState, exec.cc, vmCfg)

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
	if result != nil && len(createdContracts) > 0 {
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
			cap := gasLeft

			// Reserve in the estimation GasPool; refund after call.
			if err := gp.SubGas(cap); err != nil {
				return result, fmt.Errorf("block gas depleted (estimate) before visibility read for %s: %w", ca.Hex(), err), nil
			}

			cfg, visUsed, visErr := exec.visibilityReader.(*contractVisibilityReader).ReadVisibilityConfigMetered(ctx, vmenv, *ca, cap)
			if visErr != nil {
				gp.AddGas(cap) // refund all on error
				return result, fmt.Errorf("visibility read (estimate) failed for %s: %w", ca.Hex(), visErr), nil
			}
			_ = cfg // we don't return it from estimate, but reading it proves ABI/paths work

			if cap >= visUsed {
				gp.AddGas(cap - visUsed) // refund unused block gas
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
