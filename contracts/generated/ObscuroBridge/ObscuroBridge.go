// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObscuroBridge

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
)

// ObscuroBridgeMetaData contains all meta data concerning the ObscuroBridge contract.
var ObscuroBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"busAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161097138038061097183398101604081905261002f91610053565b600080546001600160c01b0319166001600160a01b03909216919091179055610083565b60006020828403121561006557600080fd5b81516001600160a01b038116811461007c57600080fd5b9392505050565b6108df806100926000396000f3fe6080604052600436106100345760003560e01c80631888d7121461003957806383bece4d1461004e578063c432a46f1461006e575b600080fd5b61004c6100473660046106fa565b61008e565b005b34801561005a57600080fd5b5061004c610069366004610715565b61011c565b34801561007a57600080fd5b5061004c610089366004610715565b6101ac565b600034116100e35760405162461bcd60e51b815260206004820152600f60248201527f456d707479207472616e736665722e000000000000000000000000000000000060448201526064015b60405180910390fd5b610119604051806060016040528060006001600160a01b03168152602001348152602001836001600160a01b0316815250610239565b50565b6000546001600160a01b0316331461019c5760405162461bcd60e51b815260206004820152602f60248201527f46756e6374696f6e2063616e206f6e6c792062652063616c6c6564206279207460448201527f6865206d6573736167652062757321000000000000000000000000000000000060648201526084016100da565b6101a783828461035f565b505050565b600082116101fc5760405162461bcd60e51b815260206004820152601a60248201527f417474656d7074696e6720656d707479207472616e736665722e00000000000060448201526064016100da565b61020883333085610408565b6101a76040518060600160405280856001600160a01b03168152602001848152602001836001600160a01b03168152505b600080546001600160a01b0381169163b1454caa9174010000000000000000000000000000000000000000900463ffffffff1690601461027883610751565b91906101000a81548163ffffffff021916908363ffffffff160217905550600060018111156102a9576102a9610783565b6040805186516001600160a01b039081166020808401919091528801518284015291870151909116606082015260800160405160208183030381529060405260006040518563ffffffff1660e01b815260040161030994939291906107f1565b602060405180830381600087803b15801561032357600080fd5b505af1158015610337573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061035b919061082e565b5050565b6040516001600160a01b0383166024820152604481018290526101a79084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915261045f565b6040516001600160a01b03808516602483015283166044820152606481018290526104599085907f23b872dd00000000000000000000000000000000000000000000000000000000906084016103a4565b50505050565b60006104b4826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166105449092919063ffffffff16565b8051909150156101a757808060200190518101906104d29190610858565b6101a75760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016100da565b6060610553848460008561055d565b90505b9392505050565b6060824710156105d55760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016100da565b6001600160a01b0385163b61062c5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016100da565b600080866001600160a01b03168587604051610648919061087a565b60006040518083038185875af1925050503d8060008114610685576040519150601f19603f3d011682016040523d82523d6000602084013e61068a565b606091505b509150915061069a8282866106a5565b979650505050505050565b606083156106b4575081610556565b8251156106c45782518084602001fd5b8160405162461bcd60e51b81526004016100da9190610896565b80356001600160a01b03811681146106f557600080fd5b919050565b60006020828403121561070c57600080fd5b610556826106de565b60008060006060848603121561072a57600080fd5b610733846106de565b925060208401359150610748604085016106de565b90509250925092565b600063ffffffff8083168181141561077957634e487b7160e01b600052601160045260246000fd5b6001019392505050565b634e487b7160e01b600052602160045260246000fd5b60005b838110156107b457818101518382015260200161079c565b838111156104595750506000910152565b600081518084526107dd816020860160208601610799565b601f01601f19169290920160200192915050565b600063ffffffff80871683528086166020840152506080604083015261081a60808301856107c5565b905060ff8316606083015295945050505050565b60006020828403121561084057600080fd5b815167ffffffffffffffff8116811461055657600080fd5b60006020828403121561086a57600080fd5b8151801515811461055657600080fd5b6000825161088c818460208701610799565b9190910192915050565b60208152600061055660208301846107c556fea264697066735822122053a25eee8c7ffdfa503448ed0284bfe42141fe0825aa65c133d7a7b29cd3d73a64736f6c63430008090033",
}

// ObscuroBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroBridgeMetaData.ABI instead.
var ObscuroBridgeABI = ObscuroBridgeMetaData.ABI

// ObscuroBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroBridgeMetaData.Bin instead.
var ObscuroBridgeBin = ObscuroBridgeMetaData.Bin

// DeployObscuroBridge deploys a new Ethereum contract, binding an instance of ObscuroBridge to it.
func DeployObscuroBridge(auth *bind.TransactOpts, backend bind.ContractBackend, busAddress common.Address) (common.Address, *types.Transaction, *ObscuroBridge, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroBridgeBin), backend, busAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// ObscuroBridge is an auto generated Go binding around an Ethereum contract.
type ObscuroBridge struct {
	ObscuroBridgeCaller     // Read-only binding to the contract
	ObscuroBridgeTransactor // Write-only binding to the contract
	ObscuroBridgeFilterer   // Log filterer for contract events
}

// ObscuroBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ObscuroBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObscuroBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObscuroBridgeSession struct {
	Contract     *ObscuroBridge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObscuroBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObscuroBridgeCallerSession struct {
	Contract *ObscuroBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ObscuroBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObscuroBridgeTransactorSession struct {
	Contract     *ObscuroBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ObscuroBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ObscuroBridgeRaw struct {
	Contract *ObscuroBridge // Generic contract binding to access the raw methods on
}

// ObscuroBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObscuroBridgeCallerRaw struct {
	Contract *ObscuroBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ObscuroBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactorRaw struct {
	Contract *ObscuroBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewObscuroBridge creates a new instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridge(address common.Address, backend bind.ContractBackend) (*ObscuroBridge, error) {
	contract, err := bindObscuroBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// NewObscuroBridgeCaller creates a new read-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeCaller(address common.Address, caller bind.ContractCaller) (*ObscuroBridgeCaller, error) {
	contract, err := bindObscuroBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeCaller{contract: contract}, nil
}

// NewObscuroBridgeTransactor creates a new write-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ObscuroBridgeTransactor, error) {
	contract, err := bindObscuroBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeTransactor{contract: contract}, nil
}

// NewObscuroBridgeFilterer creates a new log filterer instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ObscuroBridgeFilterer, error) {
	contract, err := bindObscuroBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeFilterer{contract: contract}, nil
}

// bindObscuroBridge binds a generic wrapper to an already deployed contract.
func bindObscuroBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ObscuroBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.ObscuroBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transact(opts, method, params...)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendAssets", asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address target) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendNative(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendNative", target)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address target) payable returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendNative(target common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, target)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address target) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendNative(target common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, target)
}
