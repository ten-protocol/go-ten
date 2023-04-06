package debugger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/l2chain"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtracers "github.com/ethereum/go-ethereum/eth/tracers"
)

// Context contains some contextual infos for a transaction execution that is not
// available from within the EVM object.
type Context struct {
	BlockHash gethcommon.Hash // Hash of the block the tx is contained within (zero if dangling tx or call)
	TxIndex   int             // Index of the transaction within a block (zero if dangling tx or call)
	TxHash    gethcommon.Hash // Hash of the transaction being traced (zero if dangling call)
}

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*logger.Config
	Tracer  *string
	Timeout *string
	Reexec  *uint64
	// Config specific to given tracer. Note struct logger
	// config are historically embedded in main object.
	TracerConfig json.RawMessage
}

type Debugger struct {
	chain       *l2chain.ObscuroChain
	storage     db.Storage
	chainConfig *params.ChainConfig
}

func New(chain *l2chain.ObscuroChain, storage db.Storage, config *params.ChainConfig) *Debugger {
	return &Debugger{
		chain:       chain,
		chainConfig: config,
		storage:     storage,
	}
}

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)

	// defaultTracechainMemLimit is the size of the triedb, at which traceChain
	// switches over and tries to use a disk-backed database instead of building
	// on top of memory.
	// For non-archive nodes, this limit _will_ be overblown, as disk-backed tries
	// will only be found every ~15K blocks or so.
	defaultTracechainMemLimit = gethcommon.StorageSize(500 * 1024 * 1024)

	// maximumPendingTraceStates is the maximum number of states allowed waiting
	// for tracing. The creation of trace state will be paused if the unused
	// trace states exceed this limit.
	maximumPendingTraceStates = 128
)

// Tracer interface extends vm.EVMLogger and additionally
// allows collecting the tracing result.
type Tracer interface {
	vm.EVMLogger
	GetResult() (json.RawMessage, error)
	// Stop terminates execution of the tracer at the first opportune moment.
	Stop(err error)
}

func (d *Debugger) DebugTraceTransaction(ctx context.Context, txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, error) {
	_, blockHash, blockNumber, index, err := d.storage.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}
	// It shouldn't happen in practice.
	if blockNumber == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	batch, err := d.storage.FetchBatch(blockHash)
	if err != nil {
		return nil, err
	}

	msg, vmctx, statedb, err := d.chain.GetChainStateAtTransaction(batch, int(index), reexec)
	if err != nil {
		return nil, err
	}

	txctx := &gethtracers.Context{
		BlockHash: blockHash,
		TxIndex:   int(index),
		TxHash:    txHash,
	}
	return d.traceTx(ctx, msg, txctx, vmctx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (d *Debugger) traceTx(ctx context.Context, message core.Message, txctx *gethtracers.Context, vmctx vm.BlockContext, statedb *state.StateDB, config *tracers.TraceConfig) (json.RawMessage, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer    vm.EVMLogger
		err       error
		txContext = core.NewEVMTxContext(message)
	)
	switch {
	case config == nil:
		tracer = logger.NewStructLogger(nil)
	case config.Tracer != nil:
		// Define a meaningful timeout of a single transaction trace
		timeout := defaultTraceTimeout
		if config.Timeout != nil {
			if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
				return nil, err
			}
		}
		if t, err := gethtracers.New(*config.Tracer, txctx); err != nil {
			return nil, err
		} else {
			deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
			go func() {
				<-deadlineCtx.Done()
				if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
					t.Stop(errors.New("execution timeout"))
				}
			}()
			defer cancel()
			tracer = t
		}
	default:
		tracer = logger.NewStructLogger(config.Config)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, txContext, statedb, d.chainConfig, vm.Config{Debug: true, Tracer: tracer, NoBaseFee: true})

	// Call Prepare to clear out the statedb access list
	statedb.Prepare(txctx.TxHash, txctx.TxIndex)

	result, err := core.ApplyMessage(vmenv, message, new(core.GasPool).AddGas(message.Gas()))
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %w", err)
	}

	// Depending on the tracer type, format and return the output.
	switch tracer := tracer.(type) {
	case *logger.StructLogger:
		// If the result contains a revert reason, return it.
		returnVal := fmt.Sprintf("%x", result.Return())
		if len(result.Revert()) > 0 {
			returnVal = fmt.Sprintf("%x", result.Revert())
		}
		exec := &gethapi.ExecutionResult{
			Gas:         result.UsedGas,
			Failed:      result.Failed(),
			ReturnValue: returnVal,
			StructLogs:  gethapi.FormatLogs(tracer.StructLogs()),
		}
		jsonRaw, err := json.Marshal(exec)
		if err != nil {
			return nil, err
		}
		return jsonRaw, nil

	case Tracer:
		return tracer.GetResult()

	default:
		panic(fmt.Sprintf("bad tracer type %T", tracer))
	}
}
