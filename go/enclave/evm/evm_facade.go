package evm

// unsafe package imported in order to link to a private function in go-ethereum.
// This allows us to customize the message generated from a signed transaction and inject custom gas logic.
import (
	"context"
	"errors"
	"fmt"
	"math/big"
	_ "unsafe"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/holiman/uint256"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
)

var ErrGasNotEnoughForL1 = errors.New("gas limit too low to pay for execution and l1 fees")

const (
	BalanceDecreaseL1Payment       tracing.BalanceChangeReason = 100
	BalanceIncreaseL1Payment       tracing.BalanceChangeReason = 101
	BalanceRevertDecreaseL1Payment tracing.BalanceChangeReason = 102
	BalanceRevertIncreaseL1Payment tracing.BalanceChangeReason = 103
)

func ExecuteTransaction(
	tx *common.L2PricedTransaction,
	s *state.StateDB,
	header *types.Header,
	chain *TenChainContext,
	cc *params.ChainConfig,
	gp *gethcore.GasPool,
	usedGas *uint64,
	vmCfg vm.Config,
	tCount int,
	logger gethlog.Logger,
) *core.TxExecResult {
	var createdContracts []*gethcommon.Address
	rules := cc.Rules(big.NewInt(0), true, 0)
	from, err := core.GetTxSigner(tx)
	if err != nil {
		return &core.TxExecResult{
			TxWithSender: &core.TxWithSender{Tx: tx.Tx, Sender: &from},
			Err:          err,
		}
	}
	s.Prepare(rules, from, gethcommon.Address{}, tx.Tx.To(), nil, nil)
	snap := s.Snapshot()
	s.SetTxContext(tx.Tx.Hash(), tCount)

	vmCfg.Tracer = &tracing.Hooks{
		// called when the code of a contract changes.
		OnCodeChange: func(addr gethcommon.Address, prevCodeHash gethcommon.Hash, prevCode []byte, codeHash gethcommon.Hash, code []byte) {
			// only proceed for new deployments.
			if len(prevCode) > 0 {
				logger.Debug("OnCodeChange: Skipping contract deployment", "address", addr.Hex())
				return
			}
			createdContracts = append(createdContracts, &addr)
			logger.Debug("OnCodeChange: Contract deployed", "address", addr.Hex())
		},
	}

	defer func() { vmCfg.Tracer = nil }()

	var vmenv *vm.EVM
	applyTx := func(
		config *params.ChainConfig,
		bc gethcore.ChainContext,
		author *gethcommon.Address,
		gp *gethcore.GasPool,
		statedb *state.StateDB,
		header *types.Header,
		tx *common.L2PricedTransaction,
		usedGas *uint64,
		cfg vm.Config,
	) (*types.Receipt, error) {
		msg, err := TransactionToMessageWithOverrides(tx, config, header)
		if err != nil {
			return nil, err
		}
		l1cost := tx.PublishingCost
		l1Gas := big.NewInt(0)
		hasL1Cost := l1cost.Cmp(big.NewInt(0)) != 0

		// If a transaction has to be published on the l1, it will have an l1 cost
		if hasL1Cost {
			l1Gas.Div(l1cost, header.BaseFee) // TotalCost/CostPerGas = Gas
			l1Gas.Add(l1Gas, big.NewInt(1))   // Cover from leftover from the division

			// The gas limit of the transaction (evm message) should always be higher than the gas overhead
			// used to cover the l1 cost
			// todo - this check has to be added to the mempool as well
			if msg.GasLimit < l1Gas.Uint64() {
				return nil, fmt.Errorf("%w. Want at least: %d have: %d", ErrGasNotEnoughForL1, l1Gas, msg.GasLimit)
			}

			// Remove the gas overhead for l1 publishing from the gas limit in order to define
			// the actual gas limit for execution
			msg.GasLimit -= l1Gas.Uint64()

			// Remove the l1 cost from the sender
			// and pay it to the coinbase of the batch
			statedb.SubBalance(msg.From, uint256.MustFromBig(l1cost), BalanceDecreaseL1Payment)
			statedb.AddBalance(header.Coinbase, uint256.MustFromBig(l1cost), BalanceIncreaseL1Payment)
		}

		// Create a new context to be used in the EVM environment
		blockContext := gethcore.NewEVMBlockContext(header, bc, author)
		vmenv = vm.NewEVM(blockContext, vm.TxContext{BlobHashes: tx.Tx.BlobHashes(), GasPrice: header.BaseFee}, statedb, config, cfg)
		var receipt *types.Receipt
		receipt, err = gethcore.ApplyTransactionWithEVM(msg, config, gp, statedb, header.Number, header.Hash(), tx.Tx, usedGas, vmenv)
		if err != nil {
			// If the transaction has l1 cost, then revert the funds exchange
			// as it will not be published on error (no receipt condition)
			if hasL1Cost {
				statedb.SubBalance(header.Coinbase, uint256.MustFromBig(l1cost), BalanceRevertIncreaseL1Payment)
				statedb.AddBalance(msg.From, uint256.MustFromBig(l1cost), BalanceRevertDecreaseL1Payment)
			}
			return receipt, err
		}

		// Synthetic transactions and ten zen are free. Do not increase the balancec of the coinbase.
		isPaidProcessing := !cfg.NoBaseFee

		// Do not increase the balance of zero address as it is the contract deployment address.
		// Doing so might cause weird interactions.
		if header.Coinbase.Big().Cmp(gethcommon.Big0) != 0 && isPaidProcessing {
			gasUsed := big.NewInt(0).SetUint64(receipt.GasUsed)
			executionGasCost := big.NewInt(0).Mul(gasUsed, header.BaseFee)
			// As the baseFee is burned, we add it back to the coinbase.
			// Geth should automatically add the tips.
			statedb.AddBalance(header.Coinbase, uint256.MustFromBig(executionGasCost), tracing.BalanceDecreaseGasBuy)
		}
		receipt.GasUsed += l1Gas.Uint64()

		return receipt, err
	}

	receipt, err := applyTx(cc, chain, nil, gp, s, header, tx, usedGas, vmCfg)
	// adjust the receipt to point to the right batch hash
	if receipt != nil {
		receipt.Logs = s.GetLogs(tx.Tx.Hash(), header.Number.Uint64(), header.Hash())
	}
	if err != nil {
		s.RevertToSnapshot(snap)
		return &core.TxExecResult{
			Receipt:      receipt,
			TxWithSender: &core.TxWithSender{Tx: tx.Tx, Sender: &from},
			Err:          err,
		}
	}

	contractsWithVisibility := make(map[gethcommon.Address]*core.ContractVisibilityConfig)
	for _, contractAddress := range createdContracts {
		contractsWithVisibility[*contractAddress] = readVisibilityConfig(vmenv, contractAddress)
	}

	return &core.TxExecResult{
		Receipt:          receipt,
		TxWithSender:     &core.TxWithSender{Tx: tx.Tx, Sender: &from},
		CreatedContracts: contractsWithVisibility,
	}
}

const (
	maxGasForVisibility = 30_000 // hardcode at 30k gas.
)

func readVisibilityConfig(vmenv *vm.EVM, contractAddress *gethcommon.Address) *core.ContractVisibilityConfig {
	cc, err := NewTransparencyConfigCaller(*contractAddress, &localContractCaller{evm: vmenv, maxGasForVisibility: maxGasForVisibility})
	if err != nil {
		// unrecoverable error. should not happen
		panic(fmt.Sprintf("could not create transparency config caller. %v", err))
	}
	visibilityRules, err := cc.VisibilityRules(nil)
	if err != nil {
		// there is no visibility defined, so we return auto
		return &core.ContractVisibilityConfig{AutoConfig: true}
	}

	transp := false
	if visibilityRules.ContractCfg == transparent {
		transp = true
	}

	cfg := &core.ContractVisibilityConfig{
		AutoConfig:   false,
		Transparent:  &transp,
		EventConfigs: make(map[gethcommon.Hash]*core.EventVisibilityConfig),
	}

	if transp {
		return cfg
	}

	// only check the config for non-transparent contracts
	for i := range visibilityRules.EventLogConfigs {
		logConfig := visibilityRules.EventLogConfigs[i]
		cfg.EventConfigs[logConfig.EventSignature] = eventCfg(logConfig)
	}

	return cfg
}

func eventCfg(logConfig ContractTransparencyConfigEventLogConfig) *core.EventVisibilityConfig {
	relevantToMap := make(map[uint8]bool)
	for _, field := range logConfig.VisibleTo {
		relevantToMap[field] = true
	}
	isPublic := relevantToMap[everyone]

	if isPublic {
		return &core.EventVisibilityConfig{AutoConfig: false, Public: true}
	}

	t1 := relevantToMap[topic1]
	t2 := relevantToMap[topic2]
	t3 := relevantToMap[topic3]
	s := relevantToMap[sender]
	return &core.EventVisibilityConfig{
		AutoConfig:    false,
		Public:        false,
		Topic1CanView: &t1,
		Topic2CanView: &t2,
		Topic3CanView: &t3,
		SenderCanView: &s,
	}
}

/*
func logReceipt(r *types.Receipt, logger gethlog.Logger) {
	if logger.Enabled(context.Background(), gethlog.LevelTrace) {
		logger.Trace("Receipt", log.TxKey, r.TxHash, "Result", receiptToString(r))
	}
}

func receiptToString(r *types.Receipt) string {
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
}
*/
// ExecuteCall - executes the eth_call call
func ExecuteCall(
	ctx context.Context,
	msg *gethcore.Message,
	s *state.StateDB,
	header *common.BatchHeader,
	storage storage.Storage,
	gethEncodingService gethencoding.EncodingService,
	chainConfig *params.ChainConfig,
	gasEstimationCap uint64,
	config enclaveconfig.EnclaveConfig,
	logger gethlog.Logger,
) (*gethcore.ExecutionResult, error) {
	noBaseFee := true
	if header.BaseFee != nil && header.BaseFee.Cmp(gethcommon.Big0) != 0 && msg.GasPrice.Cmp(gethcommon.Big0) != 0 {
		noBaseFee = false
	}

	ethHeader, err := gethEncodingService.CreateEthHeaderForBatch(ctx, header)
	if err != nil {
		return nil, err
	}

	snapshot := s.Snapshot()
	defer s.RevertToSnapshot(snapshot) // Always revert after simulation
	defer core.LogMethodDuration(logger, measure.NewStopwatch(), "evm_facade.go:Call()")

	gp := gethcore.GasPool(gasEstimationCap)
	gp.SetGas(gasEstimationCap)

	cleanState := createCleanState(s, msg, ethHeader, chainConfig)

	chain, vmCfg := initParams(storage, gethEncodingService, config, noBaseFee, nil)
	blockContext := gethcore.NewEVMBlockContext(ethHeader, chain, nil)
	// sets TxKey.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, cleanState, chainConfig, vmCfg)
	result, err := gethcore.ApplyMessage(vmenv, msg, &gp)
	// Follow the same error check structure as in geth
	// 1 - vmError / stateDB err check
	// 2 - evm.Cancelled()  todo (#1576) - support the ability to cancel function call if it takes too long
	// 3 - error check the ApplyMessage

	// Read the error stored in the database.
	if vmerr := cleanState.Error(); vmerr != nil {
		return nil, vmerr
	}

	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		logger.Debug(fmt.Sprintf("Error applying msg %v:", msg), log.CtrErrKey, err)
		return result, fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit)
	}

	return result, nil
}

func createCleanState(s *state.StateDB, msg *gethcore.Message, ethHeader *types.Header, chainConfig *params.ChainConfig) *state.StateDB {
	cleanState := s.Copy()
	cleanState.Prepare(chainConfig.Rules(ethHeader.Number, true, 0), msg.From, ethHeader.Coinbase, msg.To, nil, msg.AccessList)
	return cleanState
}

func initParams(storage storage.Storage, gethEncodingService gethencoding.EncodingService, config enclaveconfig.EnclaveConfig, noBaseFee bool, l gethlog.Logger) (*TenChainContext, vm.Config) {
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
	}
	return NewTenChainContext(storage, gethEncodingService, config, l), vmCfg
}

// used as a wrapper around the vm.EVM to allow for easier calling of smart contract view functions
type localContractCaller struct {
	evm                 *vm.EVM
	maxGasForVisibility uint64
}

// CodeAt - not implemented because it's not needed for our use case. It just has to return something non-nil
func (cc *localContractCaller) CodeAt(_ context.Context, _ gethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{0}, nil
}

func (cc *localContractCaller) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	ret, _, err := cc.evm.Call(vm.AccountRef(call.From), *call.To, call.Data, cc.maxGasForVisibility, uint256.NewInt(0))
	return ret, err
}
