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
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getRandomNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransactionTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506102f28061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80638129fc1c1461004357806381ccf4c41461004d578063dbdff2c11461006b575b5f5ffd5b61004b610073565b005b6100556101a9565b6040516100629190610203565b60405180910390f35b6100556101cc565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156100bd5750825b90505f8267ffffffffffffffff1660011480156100d95750303b155b9050811580156100e7575080155b1561011e576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561015257845468ff00000000000000001916680100000000000000001785555b83156101a257845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061019990600190610231565b60405180910390a15b5050505050565b5f44600381900b6101bc426103e861026c565b6101c69190610283565b91505090565b5f446040516020016101de91906102aa565b604051602081830303815290604052805190602001205f1c905090565b805b82525050565b6020810161021182846101fb565b92915050565b5f67ffffffffffffffff8216610211565b6101fd81610217565b602081016102118284610228565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b81810281158282048414176102115761021161023f565b8181035f83128015838313168383129190911617156102115761021161023f565b806101fd565b6102b481836102a4565b60200191905056fea26469706673582212201a281cfe3c185531f289b8e269853225e67c00015204139c5b505700011a6a0e64736f6c634300081c0033",
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

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsCaller) GetRandomNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TenSystemCalls.contract.Call(opts, &out, "getRandomNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsSession) GetRandomNumber() (*big.Int, error) {
	return _TenSystemCalls.Contract.GetRandomNumber(&_TenSystemCalls.CallOpts)
}

// GetRandomNumber is a free data retrieval call binding the contract method 0xdbdff2c1.
//
// Solidity: function getRandomNumber() view returns(uint256)
func (_TenSystemCalls *TenSystemCallsCallerSession) GetRandomNumber() (*big.Int, error) {
	return _TenSystemCalls.Contract.GetRandomNumber(&_TenSystemCalls.CallOpts)
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
