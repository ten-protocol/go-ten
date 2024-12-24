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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleAllCallbacksRan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"handleRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060408190526000805461ffff60a01b191681556001556106f638819003908190833981016040819052610033916102b8565b600080546001600160b01b0319166001600160a01b03831617905561005661005c565b506103ab565b60004861006a6003346102f4565b61007491906102f4565b9050600063a072d7b060e01b826040516024016100919190610308565b60408051601f19818403018152918152602080830180516001600160e01b039081166001600160e01b031990961695909517905281516004808252602480830185528284018051881663a4c016fb60e01b179052845191825281019093529082018051909416629e79db60e81b17909352600054919350906001600160a01b03166382fbdc9c6101226003346102f4565b856040518363ffffffff1660e01b815260040161013f919061036a565b60206040518083038185885af115801561015d573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190610182919061038c565b506000546001600160a01b03166382fbdc9c61019f6003346102f4565b846040518363ffffffff1660e01b81526004016101bc919061036a565b60206040518083038185885af11580156101da573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906101ff919061038c565b506000546001600160a01b03166382fbdc9c61021c6003346102f4565b836040518363ffffffff1660e01b8152600401610239919061036a565b60206040518083038185885af1158015610257573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061027c919061038c565b5050505050565b60006001600160a01b0382165b92915050565b61029f81610283565b81146102aa57600080fd5b50565b805161029081610296565b6000602082840312156102cd576102cd600080fd5b6102d783836102ad565b9392505050565b634e487b7160e01b600052601260045260246000fd5b600082610303576103036102de565b500490565b81815260208101610290565b60005b8381101561032f578181015183820152602001610317565b50506000910152565b6000610342825190565b808452602084019350610359818560208601610314565b601f01601f19169290920192915050565b602080825281016102d78184610338565b8061029f565b80516102908161037b565b6000602082840312156103a1576103a1600080fd5b6102d78383610381565b61033c806103ba6000396000f3fe608060405234801561001057600080fd5b50600436106100715760003560e01c8063a072d7b011610050578063a072d7b0146100f4578063a4c016fb14610107578063ee1d58721461010f57600080fd5b8062b12783146100765780635ea39558146100ac5780639e79db00146100c1575b600080fd5b6000546100969073ffffffffffffffffffffffffffffffffffffffff1681565b6040516100a3919061021f565b60405180910390f35b6100bf6100ba366004610234565b610124565b005b6100bf600080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16600160a81b179055565b6100bf610102366004610234565b61013c565b6100bf61016b565b6101176101b9565b6040516100a39190610262565b600180549060006101348361029f565b919050555050565b60005a905061014d610834836102b8565b8110610167576000805460ff60a01b1916600160a01b1790555b5050565b6000805460ff60a01b1916600160a01b1790556040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101b0906102cb565b60405180910390fd5b60008054600160a81b900460ff1680156101d557506001546003145b905090565b600073ffffffffffffffffffffffffffffffffffffffff82165b92915050565b60006101f4826101da565b60006101f4826101fa565b61021981610205565b82525050565b602081016101f48284610210565b80356101f4565b60006020828403121561024957610249600080fd5b610253838361022d565b9392505050565b801515610219565b602081016101f4828461025a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000600182016102b1576102b1610270565b5060010190565b818103818111156101f4576101f4610270565b602080825281016101f481601681527f5468697320697320612074657374206661696c7572650000000000000000000060208201526040019056fea2646970667358221220f0c046ffc2a73ce46670fd0aaf571dade3f6ef12a9b5739f85851b905f8adc4264736f6c634300081c0033",
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

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleRefund(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleRefund", callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}
