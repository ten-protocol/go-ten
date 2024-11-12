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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBefore\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfter\",\"type\":\"uint256\"}],\"name\":\"CallbackExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeNextCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"reattemptCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callback\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506016601a565b60ca565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560695760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c75780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b610e7f806100d96000396000f3fe6080604052600436106100595760003560e01c80638129fc1c116100435780638129fc1c146100ae57806382fbdc9c146100c3578063929d34e9146100d657600080fd5b8062e0d3b51461005e578063349e7eca14610097575b600080fd5b34801561006a57600080fd5b5061007e610079366004610829565b6100f6565b60405161008e94939291906108cb565b60405180910390f35b3480156100a357600080fd5b506100ac6101b1565b005b3480156100ba57600080fd5b506100ac610374565b6100ac6100d1366004610962565b6104b6565b3480156100e257600080fd5b506100ac6100f1366004610829565b61050f565b600060208190529081526040902080546001820180546001600160a01b039092169291610122906109c0565b80601f016020809104026020016040519081016040528092919081815260200182805461014e906109c0565b801561019b5780601f106101705761010080835404028352916020019161019b565b820191906000526020600020905b81548152906001019060200180831161017e57829003601f168201915b5050505050908060020154908060030154905084565b60006101be600130610a02565b9050336001600160a01b038216146101f15760405162461bcd60e51b81526004016101e890610a59565b60405180910390fd5b600254600154146103715760028054908190600061020e83610a69565b919050555060025481106102345760405162461bcd60e51b81526004016101e890610ab4565b6000818152602081905260408120600381015460028201549192909161025b908390610ada565b905060005a84546040519192506000916001600160a01b03909116908490610287906001890190610b60565b60006040518083038160008787f1925050503d80600081146102c5576040519150601f19603f3d011682016040523d82523d6000602084013e6102ca565b606091505b50509050801561031c576000868152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff191681559061030b60018301826107d4565b506000600282018190556003909101555b60005a90507f79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b287848360405161035493929190610b6a565b60405180910390a161036986600201546106b4565b505050505050505b50565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156103bf5750825b905060008267ffffffffffffffff1660011480156103dc5750303b155b9050811580156103ea575080155b15610421576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561045557845468ff00000000000000001916680100000000000000001785555b6000600181905560025583156104af57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906104a690600190610bbe565b60405180910390a15b5050505050565b600034116104d65760405162461bcd60e51b81526004016101e890610bfe565b6152086104e2346106e1565b116104ff5760405162461bcd60e51b81526004016101e890610c0e565b61050b338383346106f3565b5050565b60008181526020818152604080832081516080810190925280546001600160a01b03168252600181018054929391929184019161054b906109c0565b80601f0160208091040260200160405190810160405280929190818152602001828054610577906109c0565b80156105c45780601f10610599576101008083540402835291602001916105c4565b820191906000526020600020905b8154815290600101906020018083116105a757829003601f168201915b50505050508152602001600282015481526020016003820154815250509050600081600001516001600160a01b031682602001516040516106059190610c91565b6000604051808303816000865af19150503d8060008114610642576040519150601f19603f3d011682016040523d82523d6000602084013e610647565b606091505b50509050806106685760405162461bcd60e51b81526004016101e890610ccd565b6000838152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff191681559061069f60018301826107d4565b50600060028201819055600390910155505050565b604051419082156108fc029083906000818181858888f1935050505015801561050b573d6000803e3d6000fd5b60006106ed4883610ada565b92915050565b6040518060800160405280856001600160a01b0316815260200184848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505060208201849052486040909201919091526001805482918261076683610a69565b9091555081526020808201929092526040016000208251815473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039091161781559082015160018201906107b79082610d89565b506040820151600282015560609091015160039091015550505050565b5080546107e0906109c0565b6000825580601f106107f0575050565b601f01602090049060005260206000209081019061037191905b8082111561081e576000815560010161080a565b5090565b80356106ed565b60006020828403121561083e5761083e600080fd5b6108488383610822565b9392505050565b60006001600160a01b0382166106ed565b6108698161084f565b82525050565b60005b8381101561088a578181015183820152602001610872565b50506000910152565b600061089d825190565b8084526020840193506108b481856020860161086f565b601f01601f19169290920192915050565b80610869565b608081016108d98287610860565b81810360208301526108eb8186610893565b90506108fa60408301856108c5565b61090760608301846108c5565b95945050505050565b60008083601f84011261092557610925600080fd5b50813567ffffffffffffffff81111561094057610940600080fd5b60208301915083600182028301111561095b5761095b600080fd5b9250929050565b6000806020838503121561097857610978600080fd5b823567ffffffffffffffff81111561099257610992600080fd5b61099e85828601610910565b92509250509250929050565b634e487b7160e01b600052602260045260246000fd5b6002810460018216806109d457607f821691505b6020821081036109e6576109e66109aa565b50919050565b634e487b7160e01b600052601160045260246000fd5b6001600160a01b039182169190811690828203908111156106ed576106ed6109ec565b60088152602081017f4e6f742073656c66000000000000000000000000000000000000000000000000815290505b60200190565b602080825281016106ed81610a25565b600060018201610a7b57610a7b6109ec565b5060010190565b60168152602081017f506172616e6f69612d20746f646f3a2064656c6574650000000000000000000081529050610a53565b602080825281016106ed81610a82565b634e487b7160e01b600052601260045260246000fd5b600082610ae957610ae9610ac4565b500490565b60008154610afb816109c0565b600182168015610b125760018114610b2757610b57565b60ff1983168652811515820286019350610b57565b60008581526020902060005b83811015610b4f57815488820152600190910190602001610b33565b505081860193505b50505092915050565b6106ed8183610aee565b60608101610b7882866108c5565b610b8560208301856108c5565b610b9260408301846108c5565b949350505050565b60006106ed82610ba8565b90565b67ffffffffffffffff1690565b61086981610b9a565b602081016106ed8284610bb5565b600d8152602081017f4e6f2076616c75652073656e740000000000000000000000000000000000000081529050610a53565b602080825281016106ed81610bcc565b602080825281016106ed81602481527f47617320746f6f206c6f7720636f6d706172656420746f20636f7374206f662060208201527f63616c6c00000000000000000000000000000000000000000000000000000000604082015260600190565b6000610c79825190565b610c8781856020860161086f565b9290920192915050565b6106ed8183610c6f565b60198152602081017f43616c6c6261636b20657865637574696f6e206661696c65640000000000000081529050610a53565b602080825281016106ed81610c9b565b634e487b7160e01b600052604160045260246000fd5b60006106ed610ba58381565b610d0883610cf3565b815460001960089490940293841b1916921b91909117905550565b6000610d30818484610cff565b505050565b8181101561050b57610d48600082610d23565b600101610d35565b601f821115610d30576000818152602090206020601f85010481016020851015610d775750805b6104af6020601f860104830182610d35565b815167ffffffffffffffff811115610da357610da3610cdd565b610dad82546109c0565b610db8828285610d50565b506020601f821160018114610ded5760008315610dd55750848201515b600019600885021c19811660028502178555506104af565b600084815260208120601f198516915b82811015610e1d5787850151825560209485019460019092019101610dfd565b5084821015610e3a5783870151600019601f87166008021c191681555b5050505060020260010190555056fea2646970667358221220e3662b6b78c1b72daff09d3ae96862c0c38b5ffadafe12aeafefb4fd761ff4c164736f6c634300081c0033",
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

// PublicCallbacksCallbackExecutedIterator is returned from FilterCallbackExecuted and is used to iterate over the raw logs and unpacked data for CallbackExecuted events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecutedIterator struct {
	Event *PublicCallbacksCallbackExecuted // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackExecuted)
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
		it.Event = new(PublicCallbacksCallbackExecuted)
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
func (it *PublicCallbacksCallbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackExecuted represents a CallbackExecuted event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecuted struct {
	CallbackId *big.Int
	GasBefore  *big.Int
	GasAfter   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackExecuted is a free log retrieval operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackExecuted(opts *bind.FilterOpts) (*PublicCallbacksCallbackExecutedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackExecutedIterator{contract: _PublicCallbacks.contract, event: "CallbackExecuted", logs: logs, sub: sub}, nil
}

// WatchCallbackExecuted is a free log subscription operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackExecuted(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackExecuted) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackExecuted)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
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

// ParseCallbackExecuted is a log parse operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackExecuted(log types.Log) (*PublicCallbacksCallbackExecuted, error) {
	event := new(PublicCallbacksCallbackExecuted)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
