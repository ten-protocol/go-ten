// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Fees

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

// FeesMetaData contains all meta data concerning the Fees contract.
var FeesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"flatFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"eoaOwner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFeeForMessage\",\"type\":\"uint256\"}],\"name\":\"setMessageFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalCollectedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506015601f565b601b601f565b60cf565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff1615606e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460cc5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b61098c806100dc5f395ff3fe6080604052600436106100b0575f3560e01c80639003adfe11610066578063da35a26f1161004c578063da35a26f14610170578063e30c39781461018f578063f2fde38b146101a3575f5ffd5b80639003adfe1461014a578063afe997ea1461015c575f5ffd5b8063715018a611610096578063715018a61461010157806379ba5097146101155780638da5cb5b14610129575f5ffd5b80631a90a219146100bb57806323aa2a9d146100e0575f5ffd5b366100b757005b5f5ffd5b3480156100c6575f5ffd5b505f545b6040516100d79190610775565b60405180910390f35b3480156100eb575f5ffd5b506100ff6100fa3660046107a0565b6101c2565b005b34801561010c575f5ffd5b506100ff61020f565b348015610120575f5ffd5b506100ff610238565b348015610134575f5ffd5b5061013d610277565b6040516100d791906107dd565b348015610155575f5ffd5b50476100ca565b348015610167575f5ffd5b506100ff6102ab565b34801561017b575f5ffd5b506100ff61018a3660046107ff565b610359565b34801561019a575f5ffd5b5061013d6104d6565b3480156101ae575f5ffd5b506100ff6101bd366004610835565b6104fe565b6101ca610590565b5f8054908290556040517f5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1906102039083908590610852565b60405180910390a15050565b610217610590565b60405162461bcd60e51b815260040161022f9061086d565b60405180910390fd5b33806102426104d6565b6001600160a01b03161461026b578060405163118cdaa760e01b815260040161022f91906107dd565b610274816105c4565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b6102b3610590565b475f6102bd610277565b6001600160a01b0316826040515f6040518083038185875af1925050503d805f8114610304576040519150601f19603f3d011682016040523d82523d5f602084013e610309565b606091505b505090508061032a5760405162461bcd60e51b815260040161022f906108ce565b7fb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752826040516102039190610775565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156103a35750825b90505f8267ffffffffffffffff1660011480156103bf5750303b155b9050811580156103cd575080155b15610404576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561043857845468ff00000000000000001916680100000000000000001785555b6104418661060d565b5f8781556040517f5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f191610475918a90610920565b60405180910390a183156104cd57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906104c490600190610948565b60405180910390a15b50505050505050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0061029b565b610506610590565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081178255610557610277565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b33610599610277565b6001600160a01b0316146105c2573360405163118cdaa760e01b815260040161022f91906107dd565b565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff1916815561060982610626565b5050565b6106156106a3565b61061e8161070a565b61027461071b565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166105c2576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107126106a3565b61027481610723565b6105c26106a3565b61072b6106a3565b6001600160a01b03811661026b575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161022f91906107dd565b805b82525050565b60208101610783828461076d565b92915050565b805b8114610274575f5ffd5b803561078381610789565b5f602082840312156107b3576107b35f5ffd5b6107bd8383610795565b9392505050565b5f6001600160a01b038216610783565b61076f816107c4565b6020810161078382846107d4565b61078b816107c4565b8035610783816107eb565b5f5f60408385031215610813576108135f5ffd5b61081d8484610795565b915061082c84602085016107f4565b90509250929050565b5f60208284031215610848576108485f5ffd5b6107bd83836107f4565b60408101610860828561076d565b6107bd602083018461076d565b6020808252810161078381603481527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60208201527f742072656e6f756e6365206f776e657273686970000000000000000000000000604082015260600190565b6020808252810161078381601481527f4661696c656420746f2073656e64204574686572000000000000000000000000602082015260400190565b5f6107836109148381565b90565b61076f81610909565b604081016108608285610917565b5f67ffffffffffffffff8216610783565b61076f8161092e565b60208101610783828461093f56fea2646970667358221220212495b5d86ef0f1bc2714da702d2ba0fb3891ba148b80286ebc9f6395de610c64736f6c634300081c0033",
}

// FeesABI is the input ABI used to generate the binding from.
// Deprecated: Use FeesMetaData.ABI instead.
var FeesABI = FeesMetaData.ABI

// FeesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeesMetaData.Bin instead.
var FeesBin = FeesMetaData.Bin

// DeployFees deploys a new Ethereum contract, binding an instance of Fees to it.
func DeployFees(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Fees, error) {
	parsed, err := FeesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Fees{FeesCaller: FeesCaller{contract: contract}, FeesTransactor: FeesTransactor{contract: contract}, FeesFilterer: FeesFilterer{contract: contract}}, nil
}

// Fees is an auto generated Go binding around an Ethereum contract.
type Fees struct {
	FeesCaller     // Read-only binding to the contract
	FeesTransactor // Write-only binding to the contract
	FeesFilterer   // Log filterer for contract events
}

// FeesCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeesSession struct {
	Contract     *Fees             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeesCallerSession struct {
	Contract *FeesCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FeesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeesTransactorSession struct {
	Contract     *FeesTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeesRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeesRaw struct {
	Contract *Fees // Generic contract binding to access the raw methods on
}

// FeesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeesCallerRaw struct {
	Contract *FeesCaller // Generic read-only contract binding to access the raw methods on
}

// FeesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeesTransactorRaw struct {
	Contract *FeesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFees creates a new instance of Fees, bound to a specific deployed contract.
func NewFees(address common.Address, backend bind.ContractBackend) (*Fees, error) {
	contract, err := bindFees(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fees{FeesCaller: FeesCaller{contract: contract}, FeesTransactor: FeesTransactor{contract: contract}, FeesFilterer: FeesFilterer{contract: contract}}, nil
}

// NewFeesCaller creates a new read-only instance of Fees, bound to a specific deployed contract.
func NewFeesCaller(address common.Address, caller bind.ContractCaller) (*FeesCaller, error) {
	contract, err := bindFees(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeesCaller{contract: contract}, nil
}

// NewFeesTransactor creates a new write-only instance of Fees, bound to a specific deployed contract.
func NewFeesTransactor(address common.Address, transactor bind.ContractTransactor) (*FeesTransactor, error) {
	contract, err := bindFees(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeesTransactor{contract: contract}, nil
}

// NewFeesFilterer creates a new log filterer instance of Fees, bound to a specific deployed contract.
func NewFeesFilterer(address common.Address, filterer bind.ContractFilterer) (*FeesFilterer, error) {
	contract, err := bindFees(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeesFilterer{contract: contract}, nil
}

// bindFees binds a generic wrapper to an already deployed contract.
func bindFees(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fees *FeesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fees.Contract.FeesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fees *FeesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.Contract.FeesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fees *FeesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fees.Contract.FeesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fees *FeesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fees.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fees *FeesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fees *FeesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fees.Contract.contract.Transact(opts, method, params...)
}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesCaller) CollectedFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "collectedFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesSession) CollectedFees() (*big.Int, error) {
	return _Fees.Contract.CollectedFees(&_Fees.CallOpts)
}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesCallerSession) CollectedFees() (*big.Int, error) {
	return _Fees.Contract.CollectedFees(&_Fees.CallOpts)
}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesCaller) MessageFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "messageFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesSession) MessageFee() (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts)
}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesCallerSession) MessageFee() (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesSession) Owner() (common.Address, error) {
	return _Fees.Contract.Owner(&_Fees.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesCallerSession) Owner() (common.Address, error) {
	return _Fees.Contract.Owner(&_Fees.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesSession) PendingOwner() (common.Address, error) {
	return _Fees.Contract.PendingOwner(&_Fees.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesCallerSession) PendingOwner() (common.Address, error) {
	return _Fees.Contract.PendingOwner(&_Fees.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesSession) RenounceOwnership() error {
	return _Fees.Contract.RenounceOwnership(&_Fees.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesCallerSession) RenounceOwnership() error {
	return _Fees.Contract.RenounceOwnership(&_Fees.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesSession) AcceptOwnership() (*types.Transaction, error) {
	return _Fees.Contract.AcceptOwnership(&_Fees.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Fees.Contract.AcceptOwnership(&_Fees.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesTransactor) Initialize(opts *bind.TransactOpts, flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "initialize", flatFee, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesSession) Initialize(flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, flatFee, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesTransactorSession) Initialize(flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, flatFee, eoaOwner)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesTransactor) SetMessageFee(opts *bind.TransactOpts, newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "setMessageFee", newFeeForMessage)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesSession) SetMessageFee(newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newFeeForMessage)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesTransactorSession) SetMessageFee(newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newFeeForMessage)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.TransferOwnership(&_Fees.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.TransferOwnership(&_Fees.TransactOpts, newOwner)
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesTransactor) WithdrawalCollectedFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "withdrawalCollectedFees")
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesSession) WithdrawalCollectedFees() (*types.Transaction, error) {
	return _Fees.Contract.WithdrawalCollectedFees(&_Fees.TransactOpts)
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesTransactorSession) WithdrawalCollectedFees() (*types.Transaction, error) {
	return _Fees.Contract.WithdrawalCollectedFees(&_Fees.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesSession) Receive() (*types.Transaction, error) {
	return _Fees.Contract.Receive(&_Fees.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesTransactorSession) Receive() (*types.Transaction, error) {
	return _Fees.Contract.Receive(&_Fees.TransactOpts)
}

// FeesFeeChangedIterator is returned from FilterFeeChanged and is used to iterate over the raw logs and unpacked data for FeeChanged events raised by the Fees contract.
type FeesFeeChangedIterator struct {
	Event *FeesFeeChanged // Event containing the contract specifics and raw log

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
func (it *FeesFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesFeeChanged)
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
		it.Event = new(FeesFeeChanged)
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
func (it *FeesFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesFeeChanged represents a FeeChanged event raised by the Fees contract.
type FeesFeeChanged struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeChanged is a free log retrieval operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) FilterFeeChanged(opts *bind.FilterOpts) (*FeesFeeChangedIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "FeeChanged")
	if err != nil {
		return nil, err
	}
	return &FeesFeeChangedIterator{contract: _Fees.contract, event: "FeeChanged", logs: logs, sub: sub}, nil
}

// WatchFeeChanged is a free log subscription operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) WatchFeeChanged(opts *bind.WatchOpts, sink chan<- *FeesFeeChanged) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "FeeChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesFeeChanged)
				if err := _Fees.contract.UnpackLog(event, "FeeChanged", log); err != nil {
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

// ParseFeeChanged is a log parse operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) ParseFeeChanged(log types.Log) (*FeesFeeChanged, error) {
	event := new(FeesFeeChanged)
	if err := _Fees.contract.UnpackLog(event, "FeeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesFeeWithdrawnIterator is returned from FilterFeeWithdrawn and is used to iterate over the raw logs and unpacked data for FeeWithdrawn events raised by the Fees contract.
type FeesFeeWithdrawnIterator struct {
	Event *FeesFeeWithdrawn // Event containing the contract specifics and raw log

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
func (it *FeesFeeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesFeeWithdrawn)
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
		it.Event = new(FeesFeeWithdrawn)
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
func (it *FeesFeeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesFeeWithdrawn represents a FeeWithdrawn event raised by the Fees contract.
type FeesFeeWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeWithdrawn is a free log retrieval operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) FilterFeeWithdrawn(opts *bind.FilterOpts) (*FeesFeeWithdrawnIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return &FeesFeeWithdrawnIterator{contract: _Fees.contract, event: "FeeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFeeWithdrawn is a free log subscription operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) WatchFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *FeesFeeWithdrawn) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesFeeWithdrawn)
				if err := _Fees.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
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

// ParseFeeWithdrawn is a log parse operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) ParseFeeWithdrawn(log types.Log) (*FeesFeeWithdrawn, error) {
	event := new(FeesFeeWithdrawn)
	if err := _Fees.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Fees contract.
type FeesInitializedIterator struct {
	Event *FeesInitialized // Event containing the contract specifics and raw log

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
func (it *FeesInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesInitialized)
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
		it.Event = new(FeesInitialized)
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
func (it *FeesInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesInitialized represents a Initialized event raised by the Fees contract.
type FeesInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fees *FeesFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeesInitializedIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeesInitializedIterator{contract: _Fees.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fees *FeesFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeesInitialized) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesInitialized)
				if err := _Fees.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Fees *FeesFilterer) ParseInitialized(log types.Log) (*FeesInitialized, error) {
	event := new(FeesInitialized)
	if err := _Fees.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Fees contract.
type FeesOwnershipTransferStartedIterator struct {
	Event *FeesOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *FeesOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesOwnershipTransferStarted)
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
		it.Event = new(FeesOwnershipTransferStarted)
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
func (it *FeesOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Fees contract.
type FeesOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeesOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeesOwnershipTransferStartedIterator{contract: _Fees.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *FeesOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesOwnershipTransferStarted)
				if err := _Fees.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Fees *FeesFilterer) ParseOwnershipTransferStarted(log types.Log) (*FeesOwnershipTransferStarted, error) {
	event := new(FeesOwnershipTransferStarted)
	if err := _Fees.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Fees contract.
type FeesOwnershipTransferredIterator struct {
	Event *FeesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FeesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesOwnershipTransferred)
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
		it.Event = new(FeesOwnershipTransferred)
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
func (it *FeesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesOwnershipTransferred represents a OwnershipTransferred event raised by the Fees contract.
type FeesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeesOwnershipTransferredIterator{contract: _Fees.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesOwnershipTransferred)
				if err := _Fees.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Fees *FeesFilterer) ParseOwnershipTransferred(log types.Log) (*FeesOwnershipTransferred, error) {
	event := new(FeesOwnershipTransferred)
	if err := _Fees.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
