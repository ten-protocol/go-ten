// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package UnrenouncableOwnable2Step

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

// UnrenouncableOwnable2StepMetaData contains all meta data concerning the UnrenouncableOwnable2Step contract.
var UnrenouncableOwnable2StepMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b5061001861001d565b6100fc565b5f6100266100bd565b805490915068010000000000000000900460ff16156100585760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100ba5780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b1916100e7565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e1565b610425806101095f395ff3fe608060405234801561000f575f5ffd5b5060043610610064575f3560e01c80638da5cb5b1161004d5780638da5cb5b1461007a578063e30c397814610098578063f2fde38b146100a0575f5ffd5b8063715018a61461006857806379ba509714610072575b5f5ffd5b6100706100b3565b005b6100706100f6565b610082610135565b60405161008f919061033e565b60405180910390f35b610082610169565b6100706100ae36600461036a565b610191565b6100bb610223565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100ed9061038e565b60405180910390fd5b3380610100610169565b6001600160a01b031614610129578060405163118cdaa760e01b81526004016100ed919061033e565b61013281610257565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610159565b610199610223565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03831690811782556101ea610135565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b3361022c610135565b6001600160a01b031614610255573360405163118cdaa760e01b81526004016100ed919061033e565b565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff1916815561029c826102a0565b5050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f6001600160a01b0382165b92915050565b6103388161031d565b82525050565b60208101610329828461032f565b6103558161031d565b8114610132575f5ffd5b80356103298161034c565b5f6020828403121561037d5761037d5f5ffd5b610387838361035f565b9392505050565b6020808252810161032981603481527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60208201527f742072656e6f756e6365206f776e65727368697000000000000000000000000060408201526060019056fea2646970667358221220182cb9eb75d8f4271501b5c095cc84d9ae383ddaf14468b472104f4442dc6fb564736f6c634300081c0033",
}

// UnrenouncableOwnable2StepABI is the input ABI used to generate the binding from.
// Deprecated: Use UnrenouncableOwnable2StepMetaData.ABI instead.
var UnrenouncableOwnable2StepABI = UnrenouncableOwnable2StepMetaData.ABI

// UnrenouncableOwnable2StepBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UnrenouncableOwnable2StepMetaData.Bin instead.
var UnrenouncableOwnable2StepBin = UnrenouncableOwnable2StepMetaData.Bin

// DeployUnrenouncableOwnable2Step deploys a new Ethereum contract, binding an instance of UnrenouncableOwnable2Step to it.
func DeployUnrenouncableOwnable2Step(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UnrenouncableOwnable2Step, error) {
	parsed, err := UnrenouncableOwnable2StepMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UnrenouncableOwnable2StepBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UnrenouncableOwnable2Step{UnrenouncableOwnable2StepCaller: UnrenouncableOwnable2StepCaller{contract: contract}, UnrenouncableOwnable2StepTransactor: UnrenouncableOwnable2StepTransactor{contract: contract}, UnrenouncableOwnable2StepFilterer: UnrenouncableOwnable2StepFilterer{contract: contract}}, nil
}

// UnrenouncableOwnable2Step is an auto generated Go binding around an Ethereum contract.
type UnrenouncableOwnable2Step struct {
	UnrenouncableOwnable2StepCaller     // Read-only binding to the contract
	UnrenouncableOwnable2StepTransactor // Write-only binding to the contract
	UnrenouncableOwnable2StepFilterer   // Log filterer for contract events
}

// UnrenouncableOwnable2StepCaller is an auto generated read-only Go binding around an Ethereum contract.
type UnrenouncableOwnable2StepCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnrenouncableOwnable2StepTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UnrenouncableOwnable2StepTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnrenouncableOwnable2StepFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UnrenouncableOwnable2StepFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UnrenouncableOwnable2StepSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UnrenouncableOwnable2StepSession struct {
	Contract     *UnrenouncableOwnable2Step // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UnrenouncableOwnable2StepCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UnrenouncableOwnable2StepCallerSession struct {
	Contract *UnrenouncableOwnable2StepCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// UnrenouncableOwnable2StepTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UnrenouncableOwnable2StepTransactorSession struct {
	Contract     *UnrenouncableOwnable2StepTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// UnrenouncableOwnable2StepRaw is an auto generated low-level Go binding around an Ethereum contract.
type UnrenouncableOwnable2StepRaw struct {
	Contract *UnrenouncableOwnable2Step // Generic contract binding to access the raw methods on
}

// UnrenouncableOwnable2StepCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UnrenouncableOwnable2StepCallerRaw struct {
	Contract *UnrenouncableOwnable2StepCaller // Generic read-only contract binding to access the raw methods on
}

// UnrenouncableOwnable2StepTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UnrenouncableOwnable2StepTransactorRaw struct {
	Contract *UnrenouncableOwnable2StepTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUnrenouncableOwnable2Step creates a new instance of UnrenouncableOwnable2Step, bound to a specific deployed contract.
func NewUnrenouncableOwnable2Step(address common.Address, backend bind.ContractBackend) (*UnrenouncableOwnable2Step, error) {
	contract, err := bindUnrenouncableOwnable2Step(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2Step{UnrenouncableOwnable2StepCaller: UnrenouncableOwnable2StepCaller{contract: contract}, UnrenouncableOwnable2StepTransactor: UnrenouncableOwnable2StepTransactor{contract: contract}, UnrenouncableOwnable2StepFilterer: UnrenouncableOwnable2StepFilterer{contract: contract}}, nil
}

// NewUnrenouncableOwnable2StepCaller creates a new read-only instance of UnrenouncableOwnable2Step, bound to a specific deployed contract.
func NewUnrenouncableOwnable2StepCaller(address common.Address, caller bind.ContractCaller) (*UnrenouncableOwnable2StepCaller, error) {
	contract, err := bindUnrenouncableOwnable2Step(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepCaller{contract: contract}, nil
}

// NewUnrenouncableOwnable2StepTransactor creates a new write-only instance of UnrenouncableOwnable2Step, bound to a specific deployed contract.
func NewUnrenouncableOwnable2StepTransactor(address common.Address, transactor bind.ContractTransactor) (*UnrenouncableOwnable2StepTransactor, error) {
	contract, err := bindUnrenouncableOwnable2Step(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepTransactor{contract: contract}, nil
}

// NewUnrenouncableOwnable2StepFilterer creates a new log filterer instance of UnrenouncableOwnable2Step, bound to a specific deployed contract.
func NewUnrenouncableOwnable2StepFilterer(address common.Address, filterer bind.ContractFilterer) (*UnrenouncableOwnable2StepFilterer, error) {
	contract, err := bindUnrenouncableOwnable2Step(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepFilterer{contract: contract}, nil
}

// bindUnrenouncableOwnable2Step binds a generic wrapper to an already deployed contract.
func bindUnrenouncableOwnable2Step(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UnrenouncableOwnable2StepMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnrenouncableOwnable2Step.Contract.UnrenouncableOwnable2StepCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.UnrenouncableOwnable2StepTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.UnrenouncableOwnable2StepTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UnrenouncableOwnable2Step.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UnrenouncableOwnable2Step.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepSession) Owner() (common.Address, error) {
	return _UnrenouncableOwnable2Step.Contract.Owner(&_UnrenouncableOwnable2Step.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCallerSession) Owner() (common.Address, error) {
	return _UnrenouncableOwnable2Step.Contract.Owner(&_UnrenouncableOwnable2Step.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UnrenouncableOwnable2Step.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepSession) PendingOwner() (common.Address, error) {
	return _UnrenouncableOwnable2Step.Contract.PendingOwner(&_UnrenouncableOwnable2Step.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCallerSession) PendingOwner() (common.Address, error) {
	return _UnrenouncableOwnable2Step.Contract.PendingOwner(&_UnrenouncableOwnable2Step.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _UnrenouncableOwnable2Step.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepSession) RenounceOwnership() error {
	return _UnrenouncableOwnable2Step.Contract.RenounceOwnership(&_UnrenouncableOwnable2Step.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepCallerSession) RenounceOwnership() error {
	return _UnrenouncableOwnable2Step.Contract.RenounceOwnership(&_UnrenouncableOwnable2Step.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepSession) AcceptOwnership() (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.AcceptOwnership(&_UnrenouncableOwnable2Step.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.AcceptOwnership(&_UnrenouncableOwnable2Step.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.TransferOwnership(&_UnrenouncableOwnable2Step.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UnrenouncableOwnable2Step.Contract.TransferOwnership(&_UnrenouncableOwnable2Step.TransactOpts, newOwner)
}

// UnrenouncableOwnable2StepInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepInitializedIterator struct {
	Event *UnrenouncableOwnable2StepInitialized // Event containing the contract specifics and raw log

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
func (it *UnrenouncableOwnable2StepInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UnrenouncableOwnable2StepInitialized)
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
		it.Event = new(UnrenouncableOwnable2StepInitialized)
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
func (it *UnrenouncableOwnable2StepInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UnrenouncableOwnable2StepInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UnrenouncableOwnable2StepInitialized represents a Initialized event raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) FilterInitialized(opts *bind.FilterOpts) (*UnrenouncableOwnable2StepInitializedIterator, error) {

	logs, sub, err := _UnrenouncableOwnable2Step.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepInitializedIterator{contract: _UnrenouncableOwnable2Step.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UnrenouncableOwnable2StepInitialized) (event.Subscription, error) {

	logs, sub, err := _UnrenouncableOwnable2Step.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UnrenouncableOwnable2StepInitialized)
				if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) ParseInitialized(log types.Log) (*UnrenouncableOwnable2StepInitialized, error) {
	event := new(UnrenouncableOwnable2StepInitialized)
	if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UnrenouncableOwnable2StepOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepOwnershipTransferStartedIterator struct {
	Event *UnrenouncableOwnable2StepOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *UnrenouncableOwnable2StepOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UnrenouncableOwnable2StepOwnershipTransferStarted)
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
		it.Event = new(UnrenouncableOwnable2StepOwnershipTransferStarted)
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
func (it *UnrenouncableOwnable2StepOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UnrenouncableOwnable2StepOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UnrenouncableOwnable2StepOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UnrenouncableOwnable2StepOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UnrenouncableOwnable2Step.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepOwnershipTransferStartedIterator{contract: _UnrenouncableOwnable2Step.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *UnrenouncableOwnable2StepOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UnrenouncableOwnable2Step.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UnrenouncableOwnable2StepOwnershipTransferStarted)
				if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) ParseOwnershipTransferStarted(log types.Log) (*UnrenouncableOwnable2StepOwnershipTransferStarted, error) {
	event := new(UnrenouncableOwnable2StepOwnershipTransferStarted)
	if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UnrenouncableOwnable2StepOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepOwnershipTransferredIterator struct {
	Event *UnrenouncableOwnable2StepOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *UnrenouncableOwnable2StepOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UnrenouncableOwnable2StepOwnershipTransferred)
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
		it.Event = new(UnrenouncableOwnable2StepOwnershipTransferred)
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
func (it *UnrenouncableOwnable2StepOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UnrenouncableOwnable2StepOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UnrenouncableOwnable2StepOwnershipTransferred represents a OwnershipTransferred event raised by the UnrenouncableOwnable2Step contract.
type UnrenouncableOwnable2StepOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UnrenouncableOwnable2StepOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UnrenouncableOwnable2Step.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UnrenouncableOwnable2StepOwnershipTransferredIterator{contract: _UnrenouncableOwnable2Step.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UnrenouncableOwnable2StepOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UnrenouncableOwnable2Step.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UnrenouncableOwnable2StepOwnershipTransferred)
				if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UnrenouncableOwnable2Step *UnrenouncableOwnable2StepFilterer) ParseOwnershipTransferred(log types.Log) (*UnrenouncableOwnable2StepOwnershipTransferred, error) {
	event := new(UnrenouncableOwnable2StepOwnershipTransferred)
	if err := _UnrenouncableOwnable2Step.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
