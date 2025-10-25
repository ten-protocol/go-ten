package gethapi

// This file is a direct copy of geth @ go-ethereum/internal/ethapi/transaction_args.go
//
import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	tencommon "github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/log"
)

// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

// String returns a human-readable representation of the transaction arguments.
// This is necessary for printing the transaction arguments in SGX mode
func (args TransactionArgs) String() string {
	var parts []string
	if args.From != nil {
		parts = append(parts, fmt.Sprintf("From:%s", args.From.Hex()))
	}
	if args.To != nil {
		parts = append(parts, fmt.Sprintf("To:%s", args.To.Hex()))
	}
	if args.Gas != nil {
		parts = append(parts, fmt.Sprintf("Gas:%d", *args.Gas))
	}
	if args.GasPrice != nil {
		parts = append(parts, fmt.Sprintf("GasPrice:%s", args.GasPrice.String()))
	}
	if args.MaxFeePerGas != nil {
		parts = append(parts, fmt.Sprintf("MaxFeePerGas:%s", args.MaxFeePerGas.String()))
	}
	if args.MaxPriorityFeePerGas != nil {
		parts = append(parts, fmt.Sprintf("MaxPriorityFeePerGas:%s", args.MaxPriorityFeePerGas.String()))
	}
	if args.Value != nil {
		parts = append(parts, fmt.Sprintf("Value:%s", args.Value.String()))
	}
	if args.Nonce != nil {
		parts = append(parts, fmt.Sprintf("Nonce:%d", *args.Nonce))
	}
	if args.Data != nil {
		parts = append(parts, fmt.Sprintf("Data:0x%x", *args.Data))
	}
	if args.Input != nil {
		parts = append(parts, fmt.Sprintf("Input:0x%x", *args.Input))
	}
	if args.AccessList != nil {
		parts = append(parts, fmt.Sprintf("AccessList:%s", accessListToString(*args.AccessList)))
	}
	if args.ChainID != nil {
		parts = append(parts, fmt.Sprintf("ChainID:%s", args.ChainID.String()))
	}

	return fmt.Sprintf("TransactionArgs{%s}", strings.Join(parts, " "))
}

// Helper function to convert AccessList to string
func accessListToString(list types.AccessList) string {
	var accessListParts []string
	for _, tuple := range list {
		storageKeys := make([]string, len(tuple.StorageKeys))
		for i, key := range tuple.StorageKeys {
			storageKeys[i] = key.Hex()
		}
		accessListParts = append(accessListParts, fmt.Sprintf("{%s: [%s]}", tuple.Address.Hex(), strings.Join(storageKeys, ", ")))
	}
	return fmt.Sprintf("[%s]", strings.Join(accessListParts, ", "))
}

// from retrieves the transaction sender address.
func (args *TransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *TransactionArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   string             `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

// FormatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []logger.StructLog) []StructLogRes {
	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.ErrorString(),
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = stackValue.Hex()
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}

//// setDefaults fills in default values for unspecified tx fields.
//func (args *TransactionArgs) setDefaults(ctx context.Context, b Backend) error {
//	if err := args.setFeeDefaults(ctx, b); err != nil {
//		return err
//	}
//	if args.Value == nil {
//		args.Value = new(hexutil.Big)
//	}
//	if args.Nonce == nil {
//		nonce, err := b.GetPoolNonce(ctx, args.from())
//		if err != nil {
//			return err
//		}
//		args.Nonce = (*hexutil.Uint64)(&nonce)
//	}
//	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
//		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
//	}
//	if args.To == nil && len(args.data()) == 0 {
//		return errors.New(`contract creation without any data provided`)
//	}
//	// Estimate the gas usage if necessary.
//	if args.Gas == nil {
//		// These fields are immutable during the estimation, safe to
//		// pass the pointer directly.
//		data := args.data()
//		callArgs := TransactionArgs{
//			From:                 args.From,
//			To:                   args.To,
//			GasPrice:             args.GasPrice,
//			MaxFeePerGas:         args.MaxFeePerGas,
//			MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
//			Value:                args.Value,
//			Data:                 (*hexutil.Bytes)(&data),
//			AccessList:           args.AccessList,
//		}
//		pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
//		estimated, err := DoEstimateGas(ctx, b, callArgs, pendingBlockNr, b.RPCGasCap())
//		if err != nil {
//			return err
//		}
//		args.Gas = &estimated
//		log.Trace("Estimate gas usage automatically", "gas", args.Gas)
//	}
//	// If chain id is provided, ensure it matches the local chain id. Otherwise, set the local
//	// chain id as the default.
//	want := b.ChainConfig().ChainID
//	if args.ChainID != nil {
//		if have := (*big.Int)(args.ChainID); have.Cmp(want) != 0 {
//			return fmt.Errorf("chainId does not match node's (have=%v, want=%v)", have, want)
//		}
//	} else {
//		args.ChainID = (*hexutil.Big)(want)
//	}
//	return nil
//}
//
//// setFeeDefaults fills in default fee values for unspecified tx fields.
//func (args *TransactionArgs) setFeeDefaults(ctx context.Context, b Backend) error {
//	// If both gasPrice and at least one of the EIP-1559 fee parameters are specified, error.
//	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
//		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
//	}
//	// If the tx has completely specified a fee mechanism, no default is needed. This allows users
//	// who are not yet synced past London to get defaults for other tx values. See
//	// https://github.com/ethereum/go-ethereum/pull/23274 for more information.
//	eip1559ParamsSet := args.MaxFeePerGas != nil && args.MaxPriorityFeePerGas != nil
//	if (args.GasPrice != nil && !eip1559ParamsSet) || (args.GasPrice == nil && eip1559ParamsSet) {
//		// Sanity check the EIP-1559 fee parameters if present.
//		if args.GasPrice == nil && args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
//			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
//		}
//		return nil
//	}
//	// Now attempt to fill in default value depending on whether London is active or not.
//	head := b.CurrentHeader()
//	if b.ChainConfig().IsLondon(head.Number) {
//		// London is active, set maxPriorityFeePerGas and maxFeePerGas.
//		if err := args.setLondonFeeDefaults(ctx, head, b); err != nil {
//			return err
//		}
//	} else {
//		if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
//			return fmt.Errorf("maxFeePerGas and maxPriorityFeePerGas are not valid before London is active")
//		}
//		// London not active, set gas price.
//		price, err := b.SuggestGasTipCap(ctx)
//		if err != nil {
//			return err
//		}
//		args.GasPrice = (*hexutil.Big)(price)
//	}
//	return nil
//}
//
//// setLondonFeeDefaults fills in reasonable default fee values for unspecified fields.
//func (args *TransactionArgs) setLondonFeeDefaults(ctx context.Context, head *types.Header, b Backend) error {
//	// Set maxPriorityFeePerGas if it is missing.
//	if args.MaxPriorityFeePerGas == nil {
//		tip, err := b.SuggestGasTipCap(ctx)
//		if err != nil {
//			return err
//		}
//		args.MaxPriorityFeePerGas = (*hexutil.Big)(tip)
//	}
//	// Set maxFeePerGas if it is missing.
//	if args.MaxFeePerGas == nil {
//		// Set the max fee to be 2 times larger than the previous block's base fee.
//		// The additional slack allows the tx to not become invalidated if the base
//		// fee is rising.
//		val := new(big.Int).Add(
//			args.MaxPriorityFeePerGas.ToInt(),
//			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
//		)
//		args.MaxFeePerGas = (*hexutil.Big)(val)
//	}
//	// Both EIP-1559 fee parameters are now set; sanity check them.
//	if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
//		return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
//	}
//	return nil
//}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) toTransaction() *types.Transaction {
	var data types.TxData
	switch {
	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		data = &types.DynamicFeeTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasFeeCap:  (*big.Int)(args.MaxFeePerGas),
			GasTipCap:  (*big.Int)(args.MaxPriorityFeePerGas),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: al,
		}
	case args.AccessList != nil:
		data = &types.AccessListTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasPrice:   (*big.Int)(args.GasPrice),
			Value:      (*big.Int)(args.Value),
			Data:       args.data(),
			AccessList: *args.AccessList,
		}
	default:
		data = &types.LegacyTx{
			To:       args.To,
			Nonce:    uint64(*args.Nonce),
			Gas:      uint64(*args.Gas),
			GasPrice: (*big.Int)(args.GasPrice),
			Value:    (*big.Int)(args.Value),
			Data:     args.data(),
		}
	}
	return types.NewTx(data)
}

// ToTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) ToTransaction() *types.Transaction {
	return args.toTransaction()
}

func (args *TransactionArgs) ToMessage(globalGasCap uint64, baseFee *big.Int) (*core.Message, error) {
	// Reject invalid combinations of pre- and post-1559 fee styles
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Set sender address or use zero address if none specified.
	addr := args.from()

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = uint64(math.MaxUint64 / 2)
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		log.Debug("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}
	var (
		gasPrice  *big.Int
		gasFeeCap *big.Int
		gasTipCap *big.Int
	)
	if baseFee == nil { //nolint: nestif
		// If there's no basefee, then it must be a non-1559 execution
		gasPrice = new(big.Int)
		if args.GasPrice != nil {
			gasPrice = args.GasPrice.ToInt()
		}
		gasFeeCap, gasTipCap = gasPrice, gasPrice
	} else {
		// A basefee is provided, necessitating 1559-type execution
		if args.GasPrice != nil {
			// User specified the legacy gas field, convert to 1559 gas typing
			gasPrice = args.GasPrice.ToInt()
			gasFeeCap, gasTipCap = gasPrice, gasPrice
		} else {
			// User specified 1559 gas fields (or none), use those
			gasFeeCap = new(big.Int)
			if args.MaxFeePerGas != nil {
				gasFeeCap = args.MaxFeePerGas.ToInt()
			}
			gasTipCap = new(big.Int)
			if args.MaxPriorityFeePerGas != nil {
				gasTipCap = args.MaxPriorityFeePerGas.ToInt()
			}
			// Backfill the legacy gasPrice for EVM execution, unless we're all zeroes
			gasPrice = new(big.Int)
			if gasFeeCap.BitLen() > 0 || gasTipCap.BitLen() > 0 {
				gasPrice = tencommon.BigMin(new(big.Int).Add(gasTipCap, baseFee), gasFeeCap)
			}
		}
	}
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}
	data := args.data()
	var accessList types.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}

	nonce := uint64(0)
	if args.Nonce != nil {
		nonce = uint64(*args.Nonce)
	}
	msg := &core.Message{
		To:              args.To,
		From:            addr,
		Nonce:           nonce,
		Value:           value,
		GasLimit:        gas,
		GasPrice:        gasPrice,
		GasFeeCap:       gasFeeCap,
		GasTipCap:       gasTipCap,
		Data:            data,
		AccessList:      accessList,
		BlobGasFeeCap:   nil,
		BlobHashes:      nil,
		SkipNonceChecks: true,
	}
	return msg, nil
}
