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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newKey\",\"type\":\"string\"}],\"name\":\"handle_set_key_with_require\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newKey\",\"type\":\"string\"}],\"name\":\"set_key_with_require\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060408190526000805460ff60a01b19169055610aa63881900390819083398101604081905261002f9161021a565b600080546001600160a81b0319166001600160a01b038316179055610052610058565b506102dd565b60006100644834610256565b9050600063a072d7b060e01b82604051602401610081919061026a565b60408051601f198184030181529181526020820180516001600160e01b03166001600160e01b031990941693909317909252905190915060009063a4c016fb60e01b906100d290859060240161026a565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526000549091506001600160a01b03166382fbdc9c610123600234610256565b846040518363ffffffff1660e01b815260040161014091906102cc565b6000604051808303818588803b15801561015957600080fd5b505af115801561016d573d6000803e3d6000fd5b50506000546001600160a01b031692506382fbdc9c91506101919050600234610256565b836040518363ffffffff1660e01b81526004016101ae91906102cc565b6000604051808303818588803b1580156101c757600080fd5b505af11580156101db573d6000803e3d6000fd5b5050505050505050565b60006001600160a01b0382165b92915050565b610201816101e5565b811461020c57600080fd5b50565b80516101f2816101f8565b60006020828403121561022f5761022f600080fd5b610239838361020f565b9392505050565b634e487b7160e01b600052601260045260246000fd5b60008261026557610265610240565b500490565b818152602081016101f2565b60005b83811015610291578181015183820152602001610279565b50506000910152565b60006102a4825190565b8084526020840193506102bb818560208601610276565b601f01601f19169290920192915050565b60208082528101610239818461029a565b6107ba806102ec6000396000f3fe6080604052600436106100645760003560e01c8063a4c016fb11610043578063a4c016fb146100ee578063cefca92514610103578063ee1d58721461011657600080fd5b8062b127831461006957806363aaaaf6146100ac578063a072d7b0146100ce575b600080fd5b34801561007557600080fd5b506000546100969073ffffffffffffffffffffffffffffffffffffffff1681565b6040516100a3919061034a565b60405180910390f35b3480156100b857600080fd5b506100cc6100c7366004610450565b61013d565b005b3480156100da57600080fd5b506100cc6100e936600461049a565b610177565b3480156100fa57600080fd5b506100cc6101c0565b6100cc610111366004610450565b6101d8565b34801561012257600080fd5b50600054600160a01b900460ff166040516100a391906104c1565b80516000036101675760405162461bcd60e51b815260040161015e90610503565b60405180910390fd5b600161017382826105f2565b5050565b60005a9050610188610834836106c8565b811061017357600080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16600160a01b1790555050565b60405162461bcd60e51b815260040161015e9061070d565b60006363aaaaf660e01b826040516024016101f39190610773565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925260005491517f82fbdc9c00000000000000000000000000000000000000000000000000000000815290925073ffffffffffffffffffffffffffffffffffffffff909116906382fbdc9c9034906102ac908590600401610773565b6000604051808303818588803b1580156102c557600080fd5b505af11580156102d9573d6000803e3d6000fd5b50505050505050565b600061031f73ffffffffffffffffffffffffffffffffffffffff8316610306565b90565b73ffffffffffffffffffffffffffffffffffffffff1690565b92915050565b600061031f826102e2565b600061031f82610325565b61034481610330565b82525050565b6020810161031f828461033b565b634e487b7160e01b600052604160045260246000fd5b601f19601f830116810181811067ffffffffffffffff8211171561039457610394610358565b6040525050565b60006103a660405190565b90506103b2828261036e565b919050565b600067ffffffffffffffff8211156103d1576103d1610358565b601f19601f83011660200192915050565b82818337506000910152565b60006104016103fc846103b7565b61039b565b905082815283838301111561041857610418600080fd5b6104268360208301846103e2565b9392505050565b600082601f83011261044157610441600080fd5b610426838335602085016103ee565b60006020828403121561046557610465600080fd5b813567ffffffffffffffff81111561047f5761047f600080fd5b61048b8482850161042d565b949350505050565b803561031f565b6000602082840312156104af576104af600080fd5b6104268383610493565b801515610344565b6020810161031f82846104b9565b60178152602081017f4e6577206b65792063616e6e6f7420626520656d707479000000000000000000815290505b60200190565b6020808252810161031f816104cf565b634e487b7160e01b600052602260045260246000fd5b60028104600182168061053d57607f821691505b60208210810361054f5761054f610513565b50919050565b600061031f6103038381565b61056a83610555565b815460001960089490940293841b1916921b91909117905550565b6000610592818484610561565b505050565b81811015610173576105aa600082610585565b600101610597565b601f821115610592576000818152602090206020601f850104810160208510156105d95750805b6105eb6020601f860104830182610597565b5050505050565b815167ffffffffffffffff81111561060c5761060c610358565b6106168254610529565b6106218282856105b2565b506020601f821160018114610656576000831561063e5750848201515b600019600885021c19811660028502178555506105eb565b600084815260208120601f198516915b828110156106865787850151825560209485019460019092019101610666565b50848210156106a35783870151600019601f87166008021c191681555b50505050600202600101905550565b634e487b7160e01b600052601160045260246000fd5b8181038181111561031f5761031f6106b2565b60168152602081017f5468697320697320612074657374206661696c75726500000000000000000000815290506104fd565b6020808252810161031f816106db565b60005b83811015610738578181015183820152602001610720565b50506000910152565b600061074b825190565b80845260208401935061076281856020860161071d565b601f01601f19169290920192915050565b60208082528101610426818461074156fea26469706673582212202f333fdc1fd5f6718e5524fa3851b7d24d51b25e1ac51a65f313eb4904350f4464736f6c634300081c0033",
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

// HandleSetKeyWithRequire is a paid mutator transaction binding the contract method 0x63aaaaf6.
//
// Solidity: function handle_set_key_with_require(string newKey) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleSetKeyWithRequire(opts *bind.TransactOpts, newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handle_set_key_with_require", newKey)
}

// HandleSetKeyWithRequire is a paid mutator transaction binding the contract method 0x63aaaaf6.
//
// Solidity: function handle_set_key_with_require(string newKey) returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleSetKeyWithRequire(newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleSetKeyWithRequire(&_PublicCallbacksTest.TransactOpts, newKey)
}

// HandleSetKeyWithRequire is a paid mutator transaction binding the contract method 0x63aaaaf6.
//
// Solidity: function handle_set_key_with_require(string newKey) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleSetKeyWithRequire(newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleSetKeyWithRequire(&_PublicCallbacksTest.TransactOpts, newKey)
}

// SetKeyWithRequire is a paid mutator transaction binding the contract method 0xcefca925.
//
// Solidity: function set_key_with_require(string newKey) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) SetKeyWithRequire(opts *bind.TransactOpts, newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "set_key_with_require", newKey)
}

// SetKeyWithRequire is a paid mutator transaction binding the contract method 0xcefca925.
//
// Solidity: function set_key_with_require(string newKey) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) SetKeyWithRequire(newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.SetKeyWithRequire(&_PublicCallbacksTest.TransactOpts, newKey)
}

// SetKeyWithRequire is a paid mutator transaction binding the contract method 0xcefca925.
//
// Solidity: function set_key_with_require(string newKey) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) SetKeyWithRequire(newKey string) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.SetKeyWithRequire(&_PublicCallbacksTest.TransactOpts, newKey)
}
