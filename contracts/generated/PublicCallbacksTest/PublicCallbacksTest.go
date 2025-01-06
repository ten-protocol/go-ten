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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"callbackRefundees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleAllCallbacksRan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"handleRefund\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingRefunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060408190526000805461ffff60a01b191681556003556108ce3881900390819083398101604081905261003391610314565b600080546001600160b01b0319166001600160a01b03831617905561005661005c565b50610407565b60004861006a600334610350565b6100749190610350565b9050600063a072d7b060e01b826040516024016100919190610364565b60408051601f19818403018152918152602080830180516001600160e01b039081166001600160e01b031990961695909517905281516004808252602480830185528284018051881663a4c016fb60e01b179052845191825281019093529082018051909416629e79db60e81b179093526000805492945090916001600160a01b03166382fbdc9c610124600334610350565b866040518363ffffffff1660e01b815260040161014191906103c6565b60206040518083038185885af115801561015f573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061018491906103e8565b600081815260016020526040812080546001600160a01b03191633179055549091506001600160a01b03166382fbdc9c6101bf600334610350565b856040518363ffffffff1660e01b81526004016101dc91906103c6565b60206040518083038185885af11580156101fa573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061021f91906103e8565b600081815260016020526040812080546001600160a01b03191633179055549091506001600160a01b03166382fbdc9c61025a600334610350565b846040518363ffffffff1660e01b815260040161027791906103c6565b60206040518083038185885af1158015610295573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906102ba91906103e8565b600090815260016020526040902080546001600160a01b031916331790555050505050565b60006001600160a01b0382165b92915050565b6102fb816102df565b811461030657600080fd5b50565b80516102ec816102f2565b60006020828403121561032957610329600080fd5b6103338383610309565b9392505050565b634e487b7160e01b600052601260045260246000fd5b60008261035f5761035f61033a565b500490565b818152602081016102ec565b60005b8381101561038b578181015183820152602001610373565b50506000910152565b600061039e825190565b8084526020840193506103b5818560208601610370565b601f01601f19169290920192915050565b602080825281016103338184610394565b806102fb565b80516102ec816103d7565b6000602082840312156103fd576103fd600080fd5b61033383836103dd565b6104b8806104166000396000f3fe60806040526004361061007a5760003560e01c8063a072d7b01161004e578063a072d7b01461014d578063a4c016fb1461016d578063b613b11414610182578063ee1d5872146101bc57600080fd5b8062b127831461007f5780635ea39558146100b55780638103ab13146100ca5780639e79db001461010d575b600080fd5b34801561008b57600080fd5b5060005461009f906001600160a01b031681565b6040516100ac9190610305565b60405180910390f35b6100c86100c336600461032e565b6101de565b005b3480156100d657600080fd5b506101006100e536600461032e565b6001602052600090815260409020546001600160a01b031681565b6040516100ac919061036e565b34801561011957600080fd5b506100c8600080547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16600160a81b179055565b34801561015957600080fd5b506100c861016836600461032e565b61022f565b34801561017957600080fd5b506100c861025e565b34801561018e57600080fd5b506101af61019d366004610390565b60026020526000908152604090205481565b6040516100ac91906103b5565b3480156101c857600080fd5b506101d16102ac565b6040516100ac91906103cb565b6000818152600160209081526040808320546001600160a01b03168352600290915281208054349290610212908490610408565b9091555050600380549060006102278361041b565b919050555050565b60005a905061024061083483610434565b811061025a576000805460ff60a01b1916600160a01b1790555b5050565b6000805460ff60a01b1916600160a01b1790556040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a390610447565b60405180910390fd5b60008054600160a81b900460ff1680156102c857506003546003145b905090565b60006001600160a01b0382165b92915050565b60006102da826102cd565b60006102da826102e0565b6102ff816102eb565b82525050565b602081016102da82846102f6565b805b811461032057600080fd5b50565b80356102da81610313565b60006020828403121561034357610343600080fd5b61034d8383610323565b9392505050565b60006001600160a01b0382166102da565b6102ff81610354565b602081016102da8284610365565b61031581610354565b80356102da8161037c565b6000602082840312156103a5576103a5600080fd5b61034d8383610385565b806102ff565b602081016102da82846103af565b8015156102ff565b602081016102da82846103c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156102da576102da6103d9565b60006001820161042d5761042d6103d9565b5060010190565b818103818111156102da576102da6103d9565b602080825281016102da81601681527f5468697320697320612074657374206661696c7572650000000000000000000060208201526040019056fea26469706673582212203b1d650e5c734dbc548990dd95e5d4d6510b65d45d4f00c97321d65ef54ab1a464736f6c634300081c0033",
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

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) CallbackRefundees(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "callbackRefundees", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestSession) CallbackRefundees(arg0 *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, arg0)
}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) CallbackRefundees(arg0 *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, arg0)
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

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) PendingRefunds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "pendingRefunds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestSession) PendingRefunds(arg0 common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, arg0)
}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) PendingRefunds(arg0 common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, arg0)
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
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleRefund(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleRefund", callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}
