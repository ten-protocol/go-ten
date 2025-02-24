package evm

// unsafe package imported in order to link to a private function in go-ethereum.
// This allows us to customize the message generated from a signed transaction and inject custom gas logic.
import (
	"context"
	"fmt"
	_ "unsafe"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"

	"github.com/ethereum/go-ethereum/core/tracing"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
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

	blockContext := gethcore.NewEVMBlockContext(header, exec.chain, nil)
	// todo - why do we need the blobhashes here?
	evmEnv := vm.NewEVM(blockContext, vm.TxContext{BlobHashes: tx.Tx.BlobHashes(), GasPrice: header.BaseFee}, s, exec.cc, cfg)

	msg, err := TransactionToMessageWithOverrides(tx, exec.cc, header)
	if err != nil {
		return nil, err
	}

	receipt, err := adjustPublishingCostGas(tx, msg, s, header, noBaseFee, func() (*types.Receipt, error) {
		s.SetTxContext(tx.Tx.Hash(), tCount)
		return gethcore.ApplyTransactionWithEVM(msg, exec.cc, gp, s, header.Number, header.Hash(), tx.Tx, usedGas, evmEnv)
	})
	if err != nil {
		return nil, err
	}
	receipt.Logs = s.GetLogs(tx.Tx.Hash(), header.Number.Uint64(), header.Hash())

	contractsWithVisibility := make(map[gethcommon.Address]*core.ContractVisibilityConfig)
	for _, contractAddress := range createdContracts {
		var err1 error
		contractsWithVisibility[*contractAddress], err1 = exec.visibilityReader.ReadVisibilityConfig(context.Background(), evmEnv, *contractAddress)
		if err1 != nil {
			exec.logger.Crit("could not read visibility config. Should not happen", log.ErrKey, err1)
			return nil, err1
		}
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

	snapshot := s.Snapshot()
	defer s.RevertToSnapshot(snapshot) // Always revert after simulation

	// todo - figure out the noBaseFee logic
	noBaseFee := true
	if header.BaseFee != nil && header.BaseFee.Cmp(gethcommon.Big0) != 0 && msg.GasPrice.Cmp(gethcommon.Big0) != 0 {
		noBaseFee = false
	}
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
	}

	ethHeader, err := exec.gethEncodingService.CreateEthHeaderForBatch(ctx, header)
	if err != nil {
		return nil, nil, err
	}

	gp := gethcore.GasPool(exec.gasEstimationCap)
	gp.SetGas(exec.gasEstimationCap)

	cleanState := createCleanState(s, msg, ethHeader, exec.cc)

	blockContext := gethcore.NewEVMBlockContext(ethHeader, exec.chain, nil)
	// sets TxKey.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, cleanState, exec.cc, vmCfg)
	result, err := gethcore.ApplyMessage(vmenv, msg, &gp)
	// Follow the same error check structure as in geth
	// 1 - vmError / stateDB err check
	// 2 - evm.Cancelled()  todo (#1576) - support the ability to cancel function call if it takes too long
	// 3 - error check the ApplyMessage

	// Read the error stored in the database.
	if vmerr := cleanState.Error(); vmerr != nil {
		return nil, vmerr, nil
	}

	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		exec.logger.Debug(fmt.Sprintf("Error applying msg %v:", msg), log.CtrErrKey, err)
		return result, fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit), nil
	}

	return result, nil, nil
}

func createCleanState(s *state.StateDB, msg *gethcore.Message, ethHeader *types.Header, chainConfig *params.ChainConfig) *state.StateDB {
	cleanState := s.Copy()
	cleanState.Prepare(chainConfig.Rules(ethHeader.Number, true, 0), msg.From, ethHeader.Coinbase, msg.To, nil, msg.AccessList)
	return cleanState
}
