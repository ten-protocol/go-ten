// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Logger

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LoggerMetaData contains all meta data concerning the Logger contract.
var LoggerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"LogMessage\",\"type\":\"event\"}]",
	Bin: "0x608060405234601f5760405161017d610024823930816023015261017d90f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c632e3c2a4d0361004857307f0000000000000000000000000000000000000000000000000000000000000000036100a6575b5f80fd5b909182601f830112156100485781359167ffffffffffffffff831161004857602001926001830284011161004857565b9060208282031261004857813567ffffffffffffffff8111610048576100a2920161004c565b9091565b6100ba6100b436600461007c565b9061010a565b604051005b90825f939282370152565b91906100e8816100e1816100f29560209181520190565b80956100bf565b601f01601f191690565b0190565b6020808252610107939101916100ca565b90565b907f96561394bac381230de4649200e8831afcab1f451881bbade9ef209f6dd304809161014261013960405190565b928392836100f6565b0390a156fea2646970667358221220e73bedb654343bc33af4e06ef84b0481b6d0fce9f562fbbf938e6422c90cb45564736f6c634300081c0033",
}

// LoggerABI is the input ABI used to generate the binding from.
// Deprecated: Use LoggerMetaData.ABI instead.
var LoggerABI = LoggerMetaData.ABI

// LoggerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LoggerMetaData.Bin instead.
var LoggerBin = LoggerMetaData.Bin

// DeployLogger deploys a new Ethereum contract, binding an instance of Logger to it.
func DeployLogger(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Logger, error) {
	parsed, err := LoggerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoggerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Logger{LoggerCaller: LoggerCaller{contract: contract}, LoggerTransactor: LoggerTransactor{contract: contract}, LoggerFilterer: LoggerFilterer{contract: contract}}, nil
}

// Logger is an auto generated Go binding around an Ethereum contract.
type Logger struct {
	LoggerCaller     // Read-only binding to the contract
	LoggerTransactor // Write-only binding to the contract
	LoggerFilterer   // Log filterer for contract events
}

// LoggerCaller is an auto generated read-only Go binding around an Ethereum contract.
type LoggerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoggerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LoggerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoggerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LoggerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoggerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LoggerSession struct {
	Contract     *Logger           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LoggerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LoggerCallerSession struct {
	Contract *LoggerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LoggerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LoggerTransactorSession struct {
	Contract     *LoggerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LoggerRaw is an auto generated low-level Go binding around an Ethereum contract.
type LoggerRaw struct {
	Contract *Logger // Generic contract binding to access the raw methods on
}

// LoggerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LoggerCallerRaw struct {
	Contract *LoggerCaller // Generic read-only contract binding to access the raw methods on
}

// LoggerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LoggerTransactorRaw struct {
	Contract *LoggerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLogger creates a new instance of Logger, bound to a specific deployed contract.
func NewLogger(address common.Address, backend bind.ContractBackend) (*Logger, error) {
	contract, err := bindLogger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Logger{LoggerCaller: LoggerCaller{contract: contract}, LoggerTransactor: LoggerTransactor{contract: contract}, LoggerFilterer: LoggerFilterer{contract: contract}}, nil
}

// NewLoggerCaller creates a new read-only instance of Logger, bound to a specific deployed contract.
func NewLoggerCaller(address common.Address, caller bind.ContractCaller) (*LoggerCaller, error) {
	contract, err := bindLogger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LoggerCaller{contract: contract}, nil
}

// NewLoggerTransactor creates a new write-only instance of Logger, bound to a specific deployed contract.
func NewLoggerTransactor(address common.Address, transactor bind.ContractTransactor) (*LoggerTransactor, error) {
	contract, err := bindLogger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LoggerTransactor{contract: contract}, nil
}

// NewLoggerFilterer creates a new log filterer instance of Logger, bound to a specific deployed contract.
func NewLoggerFilterer(address common.Address, filterer bind.ContractFilterer) (*LoggerFilterer, error) {
	contract, err := bindLogger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LoggerFilterer{contract: contract}, nil
}

// bindLogger binds a generic wrapper to an already deployed contract.
func bindLogger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LoggerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Logger *LoggerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Logger.Contract.LoggerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Logger *LoggerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Logger.Contract.LoggerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Logger *LoggerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Logger.Contract.LoggerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Logger *LoggerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Logger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Logger *LoggerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Logger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Logger *LoggerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Logger.Contract.contract.Transact(opts, method, params...)
}

// LoggerLogMessageIterator is returned from FilterLogMessage and is used to iterate over the raw logs and unpacked data for LogMessage events raised by the Logger contract.
type LoggerLogMessageIterator struct {
	Event *LoggerLogMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoggerLogMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoggerLogMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoggerLogMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoggerLogMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoggerLogMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoggerLogMessage represents a LogMessage event raised by the Logger contract.
type LoggerLogMessage struct {
	Message string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogMessage is a free log retrieval operation binding the contract event 0x96561394bac381230de4649200e8831afcab1f451881bbade9ef209f6dd30480.
//
// Solidity: event LogMessage(string message)
func (_Logger *LoggerFilterer) FilterLogMessage(opts *bind.FilterOpts) (*LoggerLogMessageIterator, error) {

	logs, sub, err := _Logger.contract.FilterLogs(opts, "LogMessage")
	if err != nil {
		return nil, err
	}
	return &LoggerLogMessageIterator{contract: _Logger.contract, event: "LogMessage", logs: logs, sub: sub}, nil
}

// WatchLogMessage is a free log subscription operation binding the contract event 0x96561394bac381230de4649200e8831afcab1f451881bbade9ef209f6dd30480.
//
// Solidity: event LogMessage(string message)
func (_Logger *LoggerFilterer) WatchLogMessage(opts *bind.WatchOpts, sink chan<- *LoggerLogMessage) (event.Subscription, error) {

	logs, sub, err := _Logger.contract.WatchLogs(opts, "LogMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoggerLogMessage)
				if err := _Logger.contract.UnpackLog(event, "LogMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLogMessage is a log parse operation binding the contract event 0x96561394bac381230de4649200e8831afcab1f451881bbade9ef209f6dd30480.
//
// Solidity: event LogMessage(string message)
func (_Logger *LoggerFilterer) ParseLogMessage(log types.Log) (*LoggerLogMessage, error) {
	event := new(LoggerLogMessage)
	if err := _Logger.contract.UnpackLog(event, "LogMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
