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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"collectedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialMessageFeePerByte\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"eoaOwner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"messageSize\",\"type\":\"uint256\"}],\"name\":\"messageFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMessageFeePerByte\",\"type\":\"uint256\"}],\"name\":\"setMessageFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalCollectedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506016601a565b60ca565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560695760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c75780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6106ed806100d96000396000f3fe60806040526004361061007f5760003560e01c8063afe997ea1161004e578063afe997ea1461012c578063da35a26f14610141578063f1d44d5114610161578063f2fde38b1461018157600080fd5b806323aa2a9d1461008b578063715018a6146100ad5780638da5cb5b146100c25780639003adfe1461011057600080fd5b3661008657005b600080fd5b34801561009757600080fd5b506100ab6100a6366004610575565b6101a1565b005b3480156100b957600080fd5b506100ab6101ae565b3480156100ce57600080fd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031660405161010791906105bb565b60405180910390f35b34801561011c57600080fd5b50475b60405161010791906105cf565b34801561013857600080fd5b506100ab6101c2565b34801561014d57600080fd5b506100ab61015c3660046105f1565b610225565b34801561016d57600080fd5b5061011f61017c366004610575565b61036d565b34801561018d57600080fd5b506100ab61019c366004610629565b610383565b6101a96103e0565b600055565b6101b66103e0565b6101c06000610454565b565b6101ca6103e0565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546040516001600160a01b03909116904780156108fc02916000818181858888f19350505050158015610222573d6000803e3d6000fd5b50565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156102705750825b905060008267ffffffffffffffff16600114801561028d5750303b155b90508115801561029b575080155b156102d2576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561030657845468ff00000000000000001916680100000000000000001785555b61030f866104dd565b6000879055831561036457845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061035b90600190610663565b60405180910390a15b50505050505050565b60008160005461037d91906106a0565b92915050565b61038b6103e0565b6001600160a01b0381166103d75760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016103ce91906105bb565b60405180910390fd5b61022281610454565b336104127f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146101c057336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016103ce91906105bb565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff000000000000000000000000000000000000000081166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b6104e56104ee565b61022281610555565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166101c0576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61038b6104ee565b805b811461022257600080fd5b803561037d8161055d565b60006020828403121561058a5761058a600080fd5b610594838361056a565b9392505050565b60006001600160a01b03821661037d565b6105b58161059b565b82525050565b6020810161037d82846105ac565b806105b5565b6020810161037d82846105c9565b61055f8161059b565b803561037d816105dd565b6000806040838503121561060757610607600080fd5b610611848461056a565b915061062084602085016105e6565b90509250929050565b60006020828403121561063e5761063e600080fd5b61059483836105e6565b600067ffffffffffffffff821661037d565b6105b581610648565b6020810161037d828461065a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818102811582820484141761037d5761037d61067156fea26469706673582212208276a452ac5cb40d5c69cf89c3994ef02371b090b64009d242ea210d12394acb64736f6c634300081c0033",
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

// MessageFee is a free data retrieval call binding the contract method 0xf1d44d51.
//
// Solidity: function messageFee(uint256 messageSize) view returns(uint256)
func (_Fees *FeesCaller) MessageFee(opts *bind.CallOpts, messageSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "messageFee", messageSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MessageFee is a free data retrieval call binding the contract method 0xf1d44d51.
//
// Solidity: function messageFee(uint256 messageSize) view returns(uint256)
func (_Fees *FeesSession) MessageFee(messageSize *big.Int) (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts, messageSize)
}

// MessageFee is a free data retrieval call binding the contract method 0xf1d44d51.
//
// Solidity: function messageFee(uint256 messageSize) view returns(uint256)
func (_Fees *FeesCallerSession) MessageFee(messageSize *big.Int) (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts, messageSize)
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

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 initialMessageFeePerByte, address eoaOwner) returns()
func (_Fees *FeesTransactor) Initialize(opts *bind.TransactOpts, initialMessageFeePerByte *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "initialize", initialMessageFeePerByte, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 initialMessageFeePerByte, address eoaOwner) returns()
func (_Fees *FeesSession) Initialize(initialMessageFeePerByte *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, initialMessageFeePerByte, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 initialMessageFeePerByte, address eoaOwner) returns()
func (_Fees *FeesTransactorSession) Initialize(initialMessageFeePerByte *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, initialMessageFeePerByte, eoaOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fees *FeesTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fees *FeesSession) RenounceOwnership() (*types.Transaction, error) {
	return _Fees.Contract.RenounceOwnership(&_Fees.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fees *FeesTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Fees.Contract.RenounceOwnership(&_Fees.TransactOpts)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newMessageFeePerByte) returns()
func (_Fees *FeesTransactor) SetMessageFee(opts *bind.TransactOpts, newMessageFeePerByte *big.Int) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "setMessageFee", newMessageFeePerByte)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newMessageFeePerByte) returns()
func (_Fees *FeesSession) SetMessageFee(newMessageFeePerByte *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newMessageFeePerByte)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newMessageFeePerByte) returns()
func (_Fees *FeesTransactorSession) SetMessageFee(newMessageFeePerByte *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newMessageFeePerByte)
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
