// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TenSystemCalls

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

// TenSystemCallsMetaData contains all meta data concerning the TenSystemCalls contract.
var TenSystemCallsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getRandomNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransactionTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523460195760405161056661001e823961056690f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c80638129fc1c1461004057806381ccf4c41461003b5763dbdff2c10361004f576100a6565b61007b565b610053565b5f91031261004f57565b5f80fd5b3461004f57610063366004610045565b61006b610331565b60405180805b0390f35b9052565b565b3461004f5761008b366004610045565b6100716100966103bf565b6040519182918290815260200190565b3461004f576100b6366004610045565b61007161009661048e565b6100ce9060401c60ff1690565b90565b6100ce90546100c1565b6100ce905b67ffffffffffffffff1690565b6100ce90546100db565b6100e06100ce6100ce9290565b6101276100ce6100ce9273ffffffffffffffffffffffffffffffffffffffff1690565b73ffffffffffffffffffffffffffffffffffffffff1690565b6100ce90610104565b6100ce90610140565b6100ce6100ce6100ce9290565b9067ffffffffffffffff905b9181191691161790565b6100e06100ce6100ce9267ffffffffffffffff1690565b9061019c6100ce6101a392610175565b825461015f565b9055565b9068ff00000000000000009060401b61016b565b906101cb6100ce6101a392151590565b82546101a7565b610075906100f7565b60208101929161007991906101d2565b6101f36104f7565b8061020d610207610203836100d1565b1590565b926100ed565b916102175f6100f7565b67ffffffffffffffff8416148061032a575b600193610246610238866100f7565b9167ffffffffffffffff1690565b149081610302575b155b90816102f9575b506102cf57806102735f61026a866100f7565b9401938461018c565b6102c0575b610280575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916102af5f6102bb936101bb565b604051918291826101db565b0390a1565b6102ca83836101bb565b610278565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f610257565b905061025061031030610149565b3b61032161031d5f610152565b9190565b1491905061024e565b5080610229565b6100796101eb565b6103466100ce6100ce9290565b63ffffffff1690565b6103626100ce6100ce9263ffffffff1690565b60030b90565b634e487b7160e01b5f52601160045260245ffd5b8181029291811591840414171561038f57565b610368565b6100ce6100ce6100ce9260030b90565b9190808303925f909112801582851316918412161761038f57565b6100ce6103f06103d66103d144610339565b61034f565b6103fb6103f56103f06103ea6103e8610152565b4261037c565b610152565b91610394565b906103a4565b6100ce9081565b6100ce9054610401565b5f19811461038f5760010190565b905f199061016b565b906104396100ce6101a392610152565b8254610420565b01918252565b0190565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761048057604052565b61044a565b6100ce90610152565b6100ce61049a5f610408565b6104ac6104a682610412565b5f610429565b6104df6104b860405190565b6020808201938452909283916104d390446104468285610440565b9081038252038261045e565b6104f16104ea825190565b9160200190565b20610485565b6100ce610528565b6100ce7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610152565b6100ce6104ff56fea264697066735822122047ffec79f206ba66333500cdbbfbf403e3b5f6fbcd578e81f145541044ad915364736f6c634300081c0033",
}

// TenSystemCallsABI is the input ABI used to generate the binding from.
// Deprecated: Use TenSystemCallsMetaData.ABI instead.
var TenSystemCallsABI = TenSystemCallsMetaData.ABI

// TenSystemCallsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TenSystemCallsMetaData.Bin instead.
var TenSystemCallsBin = TenSystemCallsMetaData.Bin

// DeployTenSystemCalls deploys a new Ethereum contract, binding an instance of TenSystemCalls to it.
func DeployTenSystemCalls(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TenSystemCalls, error) {
	parsed, err := TenSystemCallsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TenSystemCallsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TenSystemCalls{TenSystemCallsCaller: TenSystemCallsCaller{contract: contract}, TenSystemCallsTransactor: TenSystemCallsTransactor{contract: contract}, TenSystemCallsFilterer: TenSystemCallsFilterer{contract: contract}}, nil
}

// TenSystemCalls is an auto generated Go binding around an Ethereum contract.
type TenSystemCalls struct {
	TenSystemCallsCaller     // Read-only binding to the contract
	TenSystemCallsTransactor // Write-only binding to the contract
	TenSystemCallsFilterer   // Log filterer for contract events
}

// TenSystemCallsCaller is an auto generated read-only Go binding around an Ethereum contract.
type TenSystemCallsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenSystemCallsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TenSystemCallsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenSystemCallsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TenSystemCallsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenSystemCallsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TenSystemCallsSession struct {
	Contract     *TenSystemCalls   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TenSystemCallsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TenSystemCallsCallerSession struct {
	Contract *TenSystemCallsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TenSystemCallsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TenSystemCallsTransactorSession struct {
	Contract     *TenSystemCallsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TenSystemCallsRaw is an auto generated low-level Go binding around an Ethereum contract.
type TenSystemCallsRaw struct {
	Contract *TenSystemCalls // Generic contract binding to access the raw methods on
}

// TenSystemCallsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TenSystemCallsCallerRaw struct {
	Contract *TenSystemCallsCaller // Generic read-only contract binding to access the raw methods on
}

// TenSystemCallsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TenSystemCallsTransactorRaw struct {
	Contract *TenSystemCallsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTenSystemCalls creates a new instance of TenSystemCalls, bound to a specific deployed contract.
func NewTenSystemCalls(address common.Address, backend bind.ContractBackend) (*TenSystemCalls, error) {
	contract, err := bindTenSystemCalls(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TenSystemCalls{TenSystemCallsCaller: TenSystemCallsCaller{contract: contract}, TenSystemCallsTransactor: TenSystemCallsTransactor{contract: contract}, TenSystemCallsFilterer: TenSystemCallsFilterer{contract: contract}}, nil
}

// NewTenSystemCallsCaller creates a new read-only instance of TenSystemCalls, bound to a specific deployed contract.
func NewTenSystemCallsCaller(address common.Address, caller bind.ContractCaller) (*TenSystemCallsCaller, error) {
	contract, err := bindTenSystemCalls(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TenSystemCallsCaller{contract: contract}, nil
}

// NewTenSystemCallsTransactor creates a new write-only instance of TenSystemCalls, bound to a specific deployed contract.
func NewTenSystemCallsTransactor(address common.Address, transactor bind.ContractTransactor) (*TenSystemCallsTransactor, error) {
	contract, err := bindTenSystemCalls(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TenSystemCallsTransactor{contract: contract}, nil
}

// NewTenSystemCallsFilterer creates a new log filterer instance of TenSystemCalls, bound to a specific deployed contract.
func NewTenSystemCallsFilterer(address common.Address, filterer bind.ContractFilterer) (*TenSystemCallsFilterer, error) {
	contract, err := bindTenSystemCalls(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TenSystemCallsFilterer{contract: contract}, nil
}

// bindTenSystemCalls binds a generic wrapper to an already deployed contract.
func bindTenSystemCalls(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TenSystemCallsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenSystemCalls *TenSystemCallsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenSystemCalls.Contract.TenSystemCallsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenSystemCalls *TenSystemCallsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenSystemCalls.Contract.TenSystemCallsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenSystemCalls *TenSystemCallsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenSystemCalls.Contract.TenSystemCallsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenSystemCalls *TenSystemCallsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenSystemCalls.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenSystemCalls *TenSystemCallsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenSystemCalls.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenSystemCalls *TenSystemCallsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenSystemCalls.Contract.contract.Transact(opts, method, params...)
}

// GetTransactionTimestamp is a free data retrieval call binding the contract method 0x81ccf4c4.
//
// Solidity: function getTransactionTimestamp() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsCaller) GetTransactionTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TenSystemCalls.contract.Call(opts, &out, "getTransactionTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTransactionTimestamp is a free data retrieval call binding the contract method 0x81ccf4c4.
//
// Solidity: function getTransactionTimestamp() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsSession) GetTransactionTimestamp() (*big.Int, error) {
	return _TenSystemCalls.Contract.GetTransactionTimestamp(&_TenSystemCalls.CallOpts)
}

// GetTransactionTimestamp is a free data retrieval call binding the contract method 0x81ccf4c4.
//
// Solidity: function getTransactionTimestamp() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsCallerSession) GetTransactionTimestamp() (*big.Int, error) {
	return _TenSystemCalls.Contract.GetTransactionTimestamp(&_TenSystemCalls.CallOpts)
}

// GetRandomNumber is a paid mutator transaction binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() returns(uint256)
func (_TenSystemCalls *TenSystemCallsTransactor) GetRandomNumber(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenSystemCalls.contract.Transact(opts, "getRandomNumber")
}

// GetRandomNumber is a paid mutator transaction binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() returns(uint256)
func (_TenSystemCalls *TenSystemCallsSession) GetRandomNumber() (*types.Transaction, error) {
	return _TenSystemCalls.Contract.GetRandomNumber(&_TenSystemCalls.TransactOpts)
}

// GetRandomNumber is a paid mutator transaction binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() returns(uint256)
func (_TenSystemCalls *TenSystemCallsTransactorSession) GetRandomNumber() (*types.Transaction, error) {
	return _TenSystemCalls.Contract.GetRandomNumber(&_TenSystemCalls.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TenSystemCalls *TenSystemCallsTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenSystemCalls.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TenSystemCalls *TenSystemCallsSession) Initialize() (*types.Transaction, error) {
	return _TenSystemCalls.Contract.Initialize(&_TenSystemCalls.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TenSystemCalls *TenSystemCallsTransactorSession) Initialize() (*types.Transaction, error) {
	return _TenSystemCalls.Contract.Initialize(&_TenSystemCalls.TransactOpts)
}

// TenSystemCallsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TenSystemCalls contract.
type TenSystemCallsInitializedIterator struct {
	Event *TenSystemCallsInitialized // Event containing the contract specifics and raw log

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
func (it *TenSystemCallsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenSystemCallsInitialized)
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
		it.Event = new(TenSystemCallsInitialized)
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
func (it *TenSystemCallsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenSystemCallsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenSystemCallsInitialized represents a Initialized event raised by the TenSystemCalls contract.
type TenSystemCallsInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenSystemCalls *TenSystemCallsFilterer) FilterInitialized(opts *bind.FilterOpts) (*TenSystemCallsInitializedIterator, error) {

	logs, sub, err := _TenSystemCalls.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TenSystemCallsInitializedIterator{contract: _TenSystemCalls.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenSystemCalls *TenSystemCallsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TenSystemCallsInitialized) (event.Subscription, error) {

	logs, sub, err := _TenSystemCalls.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenSystemCallsInitialized)
				if err := _TenSystemCalls.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenSystemCalls *TenSystemCallsFilterer) ParseInitialized(log types.Log) (*TenSystemCallsInitialized, error) {
	event := new(TenSystemCallsInitialized)
	if err := _TenSystemCalls.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
