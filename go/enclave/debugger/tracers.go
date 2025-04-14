// Package debugger: This file was copied/adapted from geth - go-ethereum/eth/tracers
package debugger

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/l2chain"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
// defaultTraceTimeout is the amount of time a single transaction can execute
// by default before being forcefully aborted.
// defaultTraceTimeout = 5 * time.Second

// defaultTraceReexec is the number of blocks the tracer is willing to go back
// and reexecute to produce missing historical state necessary to run a specific
// trace.
// defaultTraceReexec = uint64(128)
)

type Debugger struct {
	chain       l2chain.ObscuroChain
	storage     storage.Storage
	chainConfig *params.ChainConfig
}

func New(chain l2chain.ObscuroChain, storage storage.Storage, config *params.ChainConfig) *Debugger {
	return &Debugger{
		chain:       chain,
		chainConfig: config,
		storage:     storage,
	}
}

func (d *Debugger) DebugTraceTransaction(context.Context, gethcommon.Hash, *tracers.TraceConfig) (json.RawMessage, error) {
	return nil, nil
	//_, blockHash, blockNumber, index, err := d.storage.GetTransaction(txHash)
	//if err != nil {
	//	return nil, err
	//}
	//// It shouldn't happen in practice.
	//if blockNumber == 0 {
	//	return nil, errors.New("genesis is not traceable")
	//}
	//reexec := defaultTraceReexec
	//if config != nil && config.Reexec != nil {
	//	reexec = *config.Reexec
	//}
	//batch, err := d.storage.FetchBatch(blockHash)
	//if err != nil {
	//	return nil, err
	//}

	//msg, vmctx, statedb, err := d.chain.GetChainStateAtTransaction(batch, int(index), reexec)
	//if err != nil {
	//	return nil, err
	//}
	//
	//txctx := &gethtracers.Context{
	//	BlockHash: blockHash,
	//	TxIndex:   int(index),
	//	TxHash:    txHash,
	//}
	//return d.traceTx(ctx, msg, txctx, vmctx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
//
//func (d *Debugger) traceTx(ctx context.Context, message *core.Message, txctx *gethtracers.Context, vmctx vm.BlockContext, statedb *state.StateDB, config *tracers.TraceConfig) (json.RawMessage, error) {
//	// Assemble the structured logger or the JavaScript tracer
//	var (
//		tracer    vm.EVMLogger
//		err       error
//		txContext = core.NewEVMTxContext(message)
//	)
//	switch {
//	case config == nil:
//		tracer = logger.NewStructLogger(nil)
//	case config.Tracer != nil:
//		// Define a meaningful timeout of a single transaction trace
//		timeout := defaultTraceTimeout
//		if config.Timeout != nil {
//			if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
//				return nil, err
//			}
//		}
//		//if t, err := tracers.New(*config.Tracer, (*tracers.Context)(txctx)); err != nil {
//		//	return nil, err
//		//} else {
//		//	deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
//		//	go func() {
//		//		<-deadlineCtx.Done()
//		//		if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
//		//			t.Stop(errors.New("execution timeout"))
//		//		}
//		//	}()
//		//	defer cancel()
//		//	tracer = t
//		//}
//	default:
//		tracer = logger.NewStructLogger(config.Config)
//	}
//	// Run the transaction with tracing enabled.
//	vmenv := vm.NewEVM(vmctx, txContext, statedb, d.chainConfig, vm.Config{Debug: true, Tracer: tracer, NoBaseFee: true})
//
//	// Call Prepare to clear out the statedb access list
//	statedb.Prepare(txctx.TxHash, txctx.TxIndex)
//
//	result, err := core.ApplyMessage(vmenv, message, new(core.GasPool).AddGas(message.Gas()))
//	if err != nil {
//		return nil, fmt.Errorf("tracing failed: %w", err)
//	}
//
//	// Depending on the tracer type, format and return the output.
//	switch tracer := tracer.(type) {
//	case *logger.StructLogger:
//		// If the result contains a revert reason, return it.
//		returnVal := fmt.Sprintf("%x", result.Return())
//		if len(result.Revert()) > 0 {
//			returnVal = fmt.Sprintf("%x", result.Revert())
//		}
//		exec := &gethapi.ExecutionResult{
//			Gas:         result.UsedGas,
//			Failed:      result.Failed(),
//			ReturnValue: returnVal,
//			StructLogs:  gethapi.FormatLogs(tracer.StructLogs()),
//		}
//		jsonRaw, err := json.Marshal(exec)
//		if err != nil {
//			return nil, err
//		}
//		return jsonRaw, nil
//
//	case tracers.Tracer:
//		return tracer.GetResult()
//
//	default:
//		panic(fmt.Sprintf("bad tracer type %T", tracer))
//	}
//}
