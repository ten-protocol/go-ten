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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"callbackRefundees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"callbackRefundees\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleAllCallbacksRan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"handleRefund\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackRefundees\",\"type\":\"address\"}],\"name\":\"pendingRefunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"pendingRefunds\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405261001461000f6100bb565b6101a8565b6040516105b061053282396105b090f35b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761005a57604052565b610025565b9061007361006c60405190565b9283610039565b565b6001600160a01b031690565b90565b6001600160a01b0381165b0361009657565b5f80fd5b9050519061007382610084565b90602082820312610096576100819161009a565b610081610ae2803803806100ce8161005f565b9283398101906100a7565b9060ff60a01b9060a01b5b9181191691161790565b906100fe61008161010592151590565b82546100d9565b9055565b9060ff60a81b9060a81b6100e4565b9061012861008161010592151590565b8254610109565b905f19906100e4565b6100816100816100819290565b9061015561008161010592610138565b825461012f565b61008190610075906001600160a01b031682565b6100819061015c565b61008190610170565b906001600160a01b03906100e4565b906101a161008161010592610179565b8254610182565b6101d36101d9916101b95f806100ee565b6101c35f80610118565b6101ce5f6003610145565b610179565b5f610191565b6101e35f806100ee565b6101ed5f80610118565b6100736102ed565b634e487b7160e01b5f52601260045260245ffd5b8115610213570490565b6101f5565b61023161022b6100819263ffffffff1690565b60e01b90565b6001600160e01b03191690565b61008190610075565b610081905461023e565b8061008f565b9050519061007382610251565b906020828203126100965761008191610257565b90825f9392825e0152565b6102a46102ad6020936102b793610298815190565b80835293849260200190565b95869101610278565b601f01601f191690565b0190565b602080825261008192910190610283565b6040513d5f823e3d90fd5b906102e190610138565b5f5260205260405f2090565b6004600361034d61031061030961030384610138565b34610209565b4890610209565b61033e61032063a072d7b0610218565b9161032a60405190565b958693602085019081520190815260200190565b60208201810382520383610039565b600461038a61035f63a4c016fb610218565b61037b61036b60405190565b93849260208401908152015f0190565b60208201810382520382610039565b600461039c61035f639e79db00610218565b6103a86101ce5f610247565b9260206103e56382fbdc9c956103c061030385610138565b976103ca60405190565b9889809481936103da8c60e01b90565b8352600483016102bb565b03925af19485156104ee575f95610510575b5061040e610407600196876102d7565b3390610191565b602061044561041f6101ce5f610247565b61042b61030385610138565b9561043560405190565b9687809481936103da8c60e01b90565b03925af19283156104ee5761049e9461046f6104076020966103da945f916104f3575b50896102d7565b6104876103036104816101ce5f610247565b94610138565b9061049160405190565b9687958694859360e01b90565b03925af180156104ee5761007392610407925f926104bd575b506102d7565b6104e091925060203d6020116104e7575b6104d88183610039565b810190610264565b905f6104b7565b503d6104ce565b6102cc565b61050a9150883d8a116104e7576104d88183610039565b5f610468565b61052a91955060203d6020116104e7576104d88183610039565b935f6103f756fe60806040526004361015610011575f80fd5b5f3560e01c8062b127831461008f5780635ea395581461008a5780638103ab13146100855780639e79db0014610080578063a072d7b01461007b578063a4c016fb14610076578063b613b114146100715763ee1d58720361009e576102e5565b6102ba565b610244565b61022c565b610214565b6101ed565b610185565b610130565b5f91031261009e57565b5f80fd5b6100c4916008021c5b73ffffffffffffffffffffffffffffffffffffffff1690565b90565b906100c491546100a2565b6100c45f806100c7565b6100ab6100c46100c49273ffffffffffffffffffffffffffffffffffffffff1690565b6100c4906100dc565b6100c4906100ff565b61011a90610108565b9052565b60208101929161012e9190610111565b565b3461009e57610140366004610094565b61015761014b6100d2565b6040519182918261011e565b0390f35b805b0361009e57565b9050359061012e8261015b565b9060208282031261009e576100c491610164565b610198610193366004610171565b6103ad565b604051005b6100c46100c46100c49290565b906101b49061019d565b5f5260205260405f2090565b5f6101cf6100c49260016101aa565b6100c7565b61011a906100ab565b60208101929161012e91906101d4565b3461009e57610157610208610203366004610171565b6101c0565b604051918291826101dd565b3461009e57610224366004610094565b61019861043a565b3461009e5761019861023f366004610171565b610489565b3461009e57610254366004610094565b610198610522565b61015d816100ab565b9050359061012e8261025c565b9060208282031261009e576100c491610265565b906101b490610108565b6100c4916008021c81565b906100c49154610290565b5f6102b56100c4926002610286565b61029b565b3461009e576101576102d56102d0366004610272565b6102a6565b6040519182918290815260200190565b3461009e576102f5366004610094565b61015761030061054d565b60405191829182901515815260200190565b6100c4906100ab565b6100c49054610312565b6100c49081565b6100c49054610325565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b9190820180921161037057565b610336565b905f19905b9181191691161790565b906103946100c461039b9261019d565b8254610375565b9055565b5f1981146103705760010190565b6103cd6103c66103c16103e69360016101aa565b61031b565b6002610286565b6103e06103d98261032c565b3490610363565b90610384565b61012e6103fb6103f6600361032c565b61039f565b6003610384565b9075ff0000000000000000000000000000000000000000009060a81b61037a565b906104336100c461039b92151590565b8254610402565b61012e60015f610423565b9190820391821161037057565b9074ff00000000000000000000000000000000000000009060a01b61037a565b906104826100c461039b92151590565b8254610452565b6104a96104a56100c45a9361049f61083461019d565b90610445565b9190565b10156104b157565b61012e60015f610472565b156104c357565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f5468697320697320612074657374206661696c757265000000000000000000006044820152606490fd5b61052d60015f610472565b61012e5f6104bc565b6100c49060a81c60ff1690565b6100c49054610536565b6105565f610543565b8061055e5790565b50610569600361032c565b6105766104a5600361019d565b149056fea264697066735822122060396c01382bf488c7445073d8d40969a152c18a7df0f356dec906e0a04ec96964736f6c634300081c0033",
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
// Solidity: function callbackRefundees(uint256 callbackId) view returns(address callbackRefundees)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) CallbackRefundees(opts *bind.CallOpts, callbackId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "callbackRefundees", callbackId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 callbackId) view returns(address callbackRefundees)
func (_PublicCallbacksTest *PublicCallbacksTestSession) CallbackRefundees(callbackId *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, callbackId)
}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 callbackId) view returns(address callbackRefundees)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) CallbackRefundees(callbackId *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, callbackId)
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
// Solidity: function pendingRefunds(address callbackRefundees) view returns(uint256 pendingRefunds)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) PendingRefunds(opts *bind.CallOpts, callbackRefundees common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "pendingRefunds", callbackRefundees)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address callbackRefundees) view returns(uint256 pendingRefunds)
func (_PublicCallbacksTest *PublicCallbacksTestSession) PendingRefunds(callbackRefundees common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, callbackRefundees)
}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address callbackRefundees) view returns(uint256 pendingRefunds)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) PendingRefunds(callbackRefundees common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, callbackRefundees)
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
