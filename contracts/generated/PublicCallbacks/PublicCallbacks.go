// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacks

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

// PublicCallbacksMetaData contains all meta data concerning the PublicCallbacks contract.
var PublicCallbacksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeNextCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"reattemptCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callback\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506016601a565b60ca565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560695760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c75780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b610e16806100d96000396000f3fe6080604052600436106100595760003560e01c80638129fc1c116100435780638129fc1c146100ae57806382fbdc9c146100c3578063929d34e9146100d657600080fd5b8062e0d3b51461005e578063349e7eca14610097575b600080fd5b34801561006a57600080fd5b5061007e61007936600461086e565b6100f6565b60405161008e9493929190610910565b60405180910390f35b3480156100a357600080fd5b506100ac6101b1565b005b3480156100ba57600080fd5b506100ac61037e565b6100ac6100d13660046109a7565b6104c0565b3480156100e257600080fd5b506100ac6100f136600461086e565b6104d0565b600060208190529081526040902080546001820180546001600160a01b03909216929161012290610a05565b80601f016020809104026020016040519081016040528092919081815260200182805461014e90610a05565b801561019b5780601f106101705761010080835404028352916020019161019b565b820191906000526020600020905b81548152906001019060200180831161017e57829003601f168201915b5050505050908060020154908060030154905084565b60006101be600130610a47565b9050336001600160a01b038216146101f15760405162461bcd60e51b81526004016101e890610a9e565b60405180910390fd5b6002546001541161037b57600280546000918261020d83610aae565b91905055905060025481106102345760405162461bcd60e51b81526004016101e890610af9565b6000818152602081905260408120600381015460028201549192909161025b908390610b1f565b905060005a84546040519192506000916001600160a01b03909116908490610287906001890190610ba5565b60006040518083038160008787f1925050503d80600081146102c5576040519150601f19603f3d011682016040523d82523d6000602084013e6102ca565b606091505b50509050801561031c576000868152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff191681559061030b6001830182610817565b506000600282018190556003909101555b60005a905060008561032e8386610baf565b6103389190610bc2565b905080876002015461034a9190610baf565b600288015586546103659082906001600160a01b0316610675565b6103728760020154610709565b50505050505050505b50565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156103c95750825b905060008267ffffffffffffffff1660011480156103e65750303b155b9050811580156103f4575080155b1561042b576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561045f57845468ff00000000000000001916680100000000000000001785555b6000600181905560025583156104b957845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906104b090600190610bfd565b60405180910390a15b5050505050565b6104cc33838334610736565b5050565b60008181526020818152604080832081516080810190925280546001600160a01b03168252600181018054929391929184019161050c90610a05565b80601f016020809104026020016040519081016040528092919081815260200182805461053890610a05565b80156105855780601f1061055a57610100808354040283529160200191610585565b820191906000526020600020905b81548152906001019060200180831161056857829003601f168201915b50505050508152602001600282015481526020016003820154815250509050600081600001516001600160a01b031682602001516040516105c69190610c2d565b6000604051808303816000865af19150503d8060008114610603576040519150601f19603f3d011682016040523d82523d6000602084013e610608565b606091505b50509050806106295760405162461bcd60e51b81526004016101e890610c69565b6000838152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff19168155906106606001830182610817565b50600060028201819055600390910155505050565b6000816001600160a01b0316836155f090604051600060405180830381858888f193505050503d80600081146106c7576040519150601f19603f3d011682016040523d82523d6000602084013e6106cc565b606091505b505090508061070457604051419084156108fc029085906000818181858888f19350505050158015610702573d6000803e3d6000fd5b505b505050565b604051419082156108fc029083906000818181858888f193505050501580156104cc573d6000803e3d6000fd5b6040518060800160405280856001600160a01b0316815260200184848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509385525050506020820184905248604090920191909152600180548291826107a983610aae565b9091555081526020808201929092526040016000208251815473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039091161781559082015160018201906107fa9082610d20565b506040820151600282015560609091015160039091015550505050565b50805461082390610a05565b6000825580601f10610833575050565b601f01602090049060005260206000209081019061037b91905b80821115610861576000815560010161084d565b5090565b80355b92915050565b60006020828403121561088357610883600080fd5b61088d8383610865565b9392505050565b60006001600160a01b038216610868565b6108ae81610894565b82525050565b60005b838110156108cf5781810151838201526020016108b7565b50506000910152565b60006108e2825190565b8084526020840193506108f98185602086016108b4565b601f01601f19169290920192915050565b806108ae565b6080810161091e82876108a5565b818103602083015261093081866108d8565b905061093f604083018561090a565b61094c606083018461090a565b95945050505050565b60008083601f84011261096a5761096a600080fd5b50813567ffffffffffffffff81111561098557610985600080fd5b6020830191508360018202830111156109a0576109a0600080fd5b9250929050565b600080602083850312156109bd576109bd600080fd5b823567ffffffffffffffff8111156109d7576109d7600080fd5b6109e385828601610955565b92509250509250929050565b634e487b7160e01b600052602260045260246000fd5b600281046001821680610a1957607f821691505b602082108103610a2b57610a2b6109ef565b50919050565b634e487b7160e01b600052601160045260246000fd5b6001600160a01b0391821691908116908282039081111561086857610868610a31565b60088152602081017f4e6f742073656c66000000000000000000000000000000000000000000000000815290505b60200190565b6020808252810161086881610a6a565b600060018201610ac057610ac0610a31565b5060010190565b60168152602081017f506172616e6f69612d20746f646f3a2064656c6574650000000000000000000081529050610a98565b6020808252810161086881610ac7565b634e487b7160e01b600052601260045260246000fd5b600082610b2e57610b2e610b09565b500490565b60008154610b4081610a05565b600182168015610b575760018114610b6c57610b9c565b60ff1983168652811515820286019350610b9c565b60008581526020902060005b83811015610b9457815488820152600190910190602001610b78565b505081860193505b50505092915050565b6108688183610b33565b8181038181111561086857610868610a31565b818102811582820484141761086857610868610a31565b600061086882610be7565b90565b67ffffffffffffffff1690565b6108ae81610bd9565b602081016108688284610bf4565b6000610c15825190565b610c238185602086016108b4565b9290920192915050565b6108688183610c0b565b60198152602081017f43616c6c6261636b20657865637574696f6e206661696c65640000000000000081529050610a98565b6020808252810161086881610c37565b634e487b7160e01b600052604160045260246000fd5b6000610868610be48381565b610ca483610c8f565b815460001960089490940293841b1916921b91909117905550565b6000610704818484610c9b565b818110156104cc57610cdf600082610cbf565b600101610ccc565b601f821115610704576000818152602090206020601f85010481016020851015610d0e5750805b6104b96020601f860104830182610ccc565b815167ffffffffffffffff811115610d3a57610d3a610c79565b610d448254610a05565b610d4f828285610ce7565b506020601f821160018114610d845760008315610d6c5750848201515b600019600885021c19811660028502178555506104b9565b600084815260208120601f198516915b82811015610db45787850151825560209485019460019092019101610d94565b5084821015610dd15783870151600019601f87166008021c191681555b5050505060020260010190555056fea26469706673582212200965892fbc5f12605a53e4a80870ad07779988f628ada36400c03ff6f6f73b2d64736f6c634300081c0033",
}

// PublicCallbacksABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksMetaData.ABI instead.
var PublicCallbacksABI = PublicCallbacksMetaData.ABI

// PublicCallbacksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksMetaData.Bin instead.
var PublicCallbacksBin = PublicCallbacksMetaData.Bin

// DeployPublicCallbacks deploys a new Ethereum contract, binding an instance of PublicCallbacks to it.
func DeployPublicCallbacks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PublicCallbacks, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// PublicCallbacks is an auto generated Go binding around an Ethereum contract.
type PublicCallbacks struct {
	PublicCallbacksCaller     // Read-only binding to the contract
	PublicCallbacksTransactor // Write-only binding to the contract
	PublicCallbacksFilterer   // Log filterer for contract events
}

// PublicCallbacksCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksSession struct {
	Contract     *PublicCallbacks  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicCallbacksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksCallerSession struct {
	Contract *PublicCallbacksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// PublicCallbacksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTransactorSession struct {
	Contract     *PublicCallbacksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PublicCallbacksRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksRaw struct {
	Contract *PublicCallbacks // Generic contract binding to access the raw methods on
}

// PublicCallbacksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksCallerRaw struct {
	Contract *PublicCallbacksCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactorRaw struct {
	Contract *PublicCallbacksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacks creates a new instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacks(address common.Address, backend bind.ContractBackend) (*PublicCallbacks, error) {
	contract, err := bindPublicCallbacks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// NewPublicCallbacksCaller creates a new read-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksCaller, error) {
	contract, err := bindPublicCallbacks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCaller{contract: contract}, nil
}

// NewPublicCallbacksTransactor creates a new write-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTransactor, error) {
	contract, err := bindPublicCallbacks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTransactor{contract: contract}, nil
}

// NewPublicCallbacksFilterer creates a new log filterer instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksFilterer, error) {
	contract, err := bindPublicCallbacks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksFilterer{contract: contract}, nil
}

// bindPublicCallbacks binds a generic wrapper to an already deployed contract.
func bindPublicCallbacks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.PublicCallbacksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transact(opts, method, params...)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksCaller) Callbacks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbacks", arg0)

	outstruct := new(struct {
		Target  common.Address
		Data    []byte
		Value   *big.Int
		BaseFee *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Data = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Value = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.BaseFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksSession) Callbacks(arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksCallerSession) Callbacks(arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, arg0)
}

// ExecuteNextCallback is a paid mutator transaction binding the contract method 0x349e7eca.
//
// Solidity: function executeNextCallback() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ExecuteNextCallback(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "executeNextCallback")
}

// ExecuteNextCallback is a paid mutator transaction binding the contract method 0x349e7eca.
//
// Solidity: function executeNextCallback() returns()
func (_PublicCallbacks *PublicCallbacksSession) ExecuteNextCallback() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallback(&_PublicCallbacks.TransactOpts)
}

// ExecuteNextCallback is a paid mutator transaction binding the contract method 0x349e7eca.
//
// Solidity: function executeNextCallback() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ExecuteNextCallback() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallback(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ReattemptCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "reattemptCallback", callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns()
func (_PublicCallbacks *PublicCallbacksTransactor) Register(opts *bind.TransactOpts, callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "register", callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns()
func (_PublicCallbacks *PublicCallbacksSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// PublicCallbacksInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PublicCallbacks contract.
type PublicCallbacksInitializedIterator struct {
	Event *PublicCallbacksInitialized // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksInitialized)
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
		it.Event = new(PublicCallbacksInitialized)
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
func (it *PublicCallbacksInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksInitialized represents a Initialized event raised by the PublicCallbacks contract.
type PublicCallbacksInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterInitialized(opts *bind.FilterOpts) (*PublicCallbacksInitializedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksInitializedIterator{contract: _PublicCallbacks.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PublicCallbacksInitialized) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksInitialized)
				if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PublicCallbacks *PublicCallbacksFilterer) ParseInitialized(log types.Log) (*PublicCallbacksInitialized, error) {
	event := new(PublicCallbacksInitialized)
	if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
