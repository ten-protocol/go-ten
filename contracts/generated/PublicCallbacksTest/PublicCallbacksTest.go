// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacksTest

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

// PublicCallbacksTestMetaData contains all meta data concerning the PublicCallbacksTest contract.
var PublicCallbacksTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleAllCallbacksRan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060408190526000805461ffff60a01b1916905561068838819003908190833981016040819052610030916102b5565b600080546001600160b01b0319166001600160a01b038316179055610053610059565b506103a8565b6000486100676003346102f1565b61007191906102f1565b9050600063a072d7b060e01b8260405160240161008e9190610305565b60408051601f19818403018152918152602080830180516001600160e01b039081166001600160e01b031990961695909517905281516004808252602480830185528284018051881663a4c016fb60e01b179052845191825281019093529082018051909416629e79db60e81b17909352600054919350906001600160a01b03166382fbdc9c61011f6003346102f1565b856040518363ffffffff1660e01b815260040161013c9190610367565b60206040518083038185885af115801561015a573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061017f9190610389565b506000546001600160a01b03166382fbdc9c61019c6003346102f1565b846040518363ffffffff1660e01b81526004016101b99190610367565b60206040518083038185885af11580156101d7573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906101fc9190610389565b506000546001600160a01b03166382fbdc9c6102196003346102f1565b836040518363ffffffff1660e01b81526004016102369190610367565b60206040518083038185885af1158015610254573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906102799190610389565b5050505050565b60006001600160a01b0382165b92915050565b61029c81610280565b81146102a757600080fd5b50565b805161028d81610293565b6000602082840312156102ca576102ca600080fd5b6102d483836102aa565b9392505050565b634e487b7160e01b600052601260045260246000fd5b600082610300576103006102db565b500490565b8181526020810161028d565b60005b8381101561032c578181015183820152602001610314565b50506000910152565b600061033f825190565b808452602084019350610356818560208601610311565b601f01601f19169290920192915050565b602080825281016102d48184610335565b8061029c565b805161028d81610378565b60006020828403121561039e5761039e600080fd5b6102d4838361037e565b6102d1806103b76000396000f3fe608060405234801561001057600080fd5b50600436106100665760003560e01c8063a072d7b011610050578063a072d7b0146100d6578063a4c016fb146100e9578063ee1d5872146100f157600080fd5b8062b127831461006b5780639e79db00146100a1575b600080fd5b60005461008b9073ffffffffffffffffffffffffffffffffffffffff1681565b60405161009891906101cd565b60405180910390f35b6100d4600080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16600160a81b179055565b005b6100d46100e43660046101e2565b61010b565b6100d461013a565b600054600160a81b900460ff166040516100989190610210565b60005a905061011c6108348361024d565b8110610136576000805460ff60a01b1916600160a01b1790555b5050565b6000805460ff60a01b1916600160a01b1790556040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161017f90610260565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff82165b92915050565b60006101a282610188565b60006101a2826101a8565b6101c7816101b3565b82525050565b602081016101a282846101be565b80356101a2565b6000602082840312156101f7576101f7600080fd5b61020183836101db565b9392505050565b8015156101c7565b602081016101a28284610208565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156101a2576101a261021e565b602080825281016101a281601681527f5468697320697320612074657374206661696c7572650000000000000000000060208201526040019056fea264697066735822122034157dc839226d5d21d57d82e72a5696aa65217d03e1d0b1b1457fd4d75d977964736f6c634300081c0033",
}

// PublicCallbacksTestABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksTestMetaData.ABI instead.
var PublicCallbacksTestABI = PublicCallbacksTestMetaData.ABI

// PublicCallbacksTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksTestMetaData.Bin instead.
var PublicCallbacksTestBin = PublicCallbacksTestMetaData.Bin

// DeployPublicCallbacksTest deploys a new Ethereum contract, binding an instance of PublicCallbacksTest to it.
func DeployPublicCallbacksTest(auth *bind.TransactOpts, backend bind.ContractBackend, _callbacks common.Address) (common.Address, *types.Transaction, *PublicCallbacksTest, error) {
	parsed, err := PublicCallbacksTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksTestBin), backend, _callbacks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacksTest{PublicCallbacksTestCaller: PublicCallbacksTestCaller{contract: contract}, PublicCallbacksTestTransactor: PublicCallbacksTestTransactor{contract: contract}, PublicCallbacksTestFilterer: PublicCallbacksTestFilterer{contract: contract}}, nil
}

// PublicCallbacksTest is an auto generated Go binding around an Ethereum contract.
type PublicCallbacksTest struct {
	PublicCallbacksTestCaller     // Read-only binding to the contract
	PublicCallbacksTestTransactor // Write-only binding to the contract
	PublicCallbacksTestFilterer   // Log filterer for contract events
}

// PublicCallbacksTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksTestSession struct {
	Contract     *PublicCallbacksTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PublicCallbacksTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksTestCallerSession struct {
	Contract *PublicCallbacksTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// PublicCallbacksTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTestTransactorSession struct {
	Contract     *PublicCallbacksTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// PublicCallbacksTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksTestRaw struct {
	Contract *PublicCallbacksTest // Generic contract binding to access the raw methods on
}

// PublicCallbacksTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksTestCallerRaw struct {
	Contract *PublicCallbacksTestCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTestTransactorRaw struct {
	Contract *PublicCallbacksTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacksTest creates a new instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTest(address common.Address, backend bind.ContractBackend) (*PublicCallbacksTest, error) {
	contract, err := bindPublicCallbacksTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTest{PublicCallbacksTestCaller: PublicCallbacksTestCaller{contract: contract}, PublicCallbacksTestTransactor: PublicCallbacksTestTransactor{contract: contract}, PublicCallbacksTestFilterer: PublicCallbacksTestFilterer{contract: contract}}, nil
}

// NewPublicCallbacksTestCaller creates a new read-only instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksTestCaller, error) {
	contract, err := bindPublicCallbacksTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestCaller{contract: contract}, nil
}

// NewPublicCallbacksTestTransactor creates a new write-only instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTestTransactor, error) {
	contract, err := bindPublicCallbacksTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestTransactor{contract: contract}, nil
}

// NewPublicCallbacksTestFilterer creates a new log filterer instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksTestFilterer, error) {
	contract, err := bindPublicCallbacksTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestFilterer{contract: contract}, nil
}

// bindPublicCallbacksTest binds a generic wrapper to an already deployed contract.
func bindPublicCallbacksTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacksTest *PublicCallbacksTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacksTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacksTest *PublicCallbacksTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacksTest *PublicCallbacksTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.contract.Transact(opts, method, params...)
}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) Callbacks(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "callbacks")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestSession) Callbacks() (common.Address, error) {
	return _PublicCallbacksTest.Contract.Callbacks(&_PublicCallbacksTest.CallOpts)
}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) Callbacks() (common.Address, error) {
	return _PublicCallbacksTest.Contract.Callbacks(&_PublicCallbacksTest.CallOpts)
}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) IsLastCallSuccess(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "isLastCallSuccess")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestSession) IsLastCallSuccess() (bool, error) {
	return _PublicCallbacksTest.Contract.IsLastCallSuccess(&_PublicCallbacksTest.CallOpts)
}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) IsLastCallSuccess() (bool, error) {
	return _PublicCallbacksTest.Contract.IsLastCallSuccess(&_PublicCallbacksTest.CallOpts)
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleAllCallbacksRan(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleAllCallbacksRan")
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleAllCallbacksRan() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleAllCallbacksRan(&_PublicCallbacksTest.TransactOpts)
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleAllCallbacksRan() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleAllCallbacksRan(&_PublicCallbacksTest.TransactOpts)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleCallback(opts *bind.TransactOpts, expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleCallback", expectedGas)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleCallback(expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallback(&_PublicCallbacksTest.TransactOpts, expectedGas)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleCallback(expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallback(&_PublicCallbacksTest.TransactOpts, expectedGas)
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleCallbackFail(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleCallbackFail")
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleCallbackFail() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallbackFail(&_PublicCallbacksTest.TransactOpts)
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleCallbackFail() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallbackFail(&_PublicCallbacksTest.TransactOpts)
}
