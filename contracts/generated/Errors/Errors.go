// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Errors

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

// ErrorsMetaData contains all meta data concerning the Errors contract.
var ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedDeployment\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"MissingPrecompile\",\"type\":\"error\"}]",
	Bin: "0x608060405234601a57604051603e601f823930815050603e90f35b5f80fdfe60806040525f80fdfea264697066735822122010c977e2e8be7bb7a7423100318ef4af0b269cf012d4c151f3c238931e2a309364736f6c634300081c0033",
}

// ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use ErrorsMetaData.ABI instead.
var ErrorsABI = ErrorsMetaData.ABI

// ErrorsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ErrorsMetaData.Bin instead.
var ErrorsBin = ErrorsMetaData.Bin

// DeployErrors deploys a new Ethereum contract, binding an instance of Errors to it.
func DeployErrors(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Errors, error) {
	parsed, err := ErrorsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ErrorsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Errors{ErrorsCaller: ErrorsCaller{contract: contract}, ErrorsTransactor: ErrorsTransactor{contract: contract}, ErrorsFilterer: ErrorsFilterer{contract: contract}}, nil
}

// Errors is an auto generated Go binding around an Ethereum contract.
type Errors struct {
	ErrorsCaller     // Read-only binding to the contract
	ErrorsTransactor // Write-only binding to the contract
	ErrorsFilterer   // Log filterer for contract events
}

// ErrorsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ErrorsSession struct {
	Contract     *Errors           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ErrorsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ErrorsCallerSession struct {
	Contract *ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ErrorsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ErrorsTransactorSession struct {
	Contract     *ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ErrorsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ErrorsRaw struct {
	Contract *Errors // Generic contract binding to access the raw methods on
}

// ErrorsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ErrorsCallerRaw struct {
	Contract *ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// ErrorsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ErrorsTransactorRaw struct {
	Contract *ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErrors creates a new instance of Errors, bound to a specific deployed contract.
func NewErrors(address common.Address, backend bind.ContractBackend) (*Errors, error) {
	contract, err := bindErrors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Errors{ErrorsCaller: ErrorsCaller{contract: contract}, ErrorsTransactor: ErrorsTransactor{contract: contract}, ErrorsFilterer: ErrorsFilterer{contract: contract}}, nil
}

// NewErrorsCaller creates a new read-only instance of Errors, bound to a specific deployed contract.
func NewErrorsCaller(address common.Address, caller bind.ContractCaller) (*ErrorsCaller, error) {
	contract, err := bindErrors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ErrorsCaller{contract: contract}, nil
}

// NewErrorsTransactor creates a new write-only instance of Errors, bound to a specific deployed contract.
func NewErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*ErrorsTransactor, error) {
	contract, err := bindErrors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ErrorsTransactor{contract: contract}, nil
}

// NewErrorsFilterer creates a new log filterer instance of Errors, bound to a specific deployed contract.
func NewErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*ErrorsFilterer, error) {
	contract, err := bindErrors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ErrorsFilterer{contract: contract}, nil
}

// bindErrors binds a generic wrapper to an already deployed contract.
func bindErrors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Errors *ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Errors.Contract.ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Errors *ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Errors.Contract.ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Errors *ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Errors.Contract.ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Errors *ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Errors *ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Errors *ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Errors.Contract.contract.Transact(opts, method, params...)
}
