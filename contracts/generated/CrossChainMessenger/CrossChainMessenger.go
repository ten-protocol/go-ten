// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CrossChainMessenger

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

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// CrossChainMessengerMetaData contains all meta data concerning the CrossChainMessenger contract.
var CrossChainMessengerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messageBusAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"crossChainSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"encodeCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600180546001600160a01b031916905534801561002057600080fd5b506040516109ad3803806109ad83398101604081905261003f91610064565b600080546001600160a01b0319166001600160a01b0392909216919091179055610094565b60006020828403121561007657600080fd5b81516001600160a01b038116811461008d57600080fd5b9392505050565b61090a806100a36000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80635b76f28b1461005157806363012de51461007a5780639b7cf1ee146100a5578063a1a227fa146100ba575b600080fd5b61006461005f366004610425565b6100cb565b6040516100719190610504565b60405180910390f35b60015461008d906001600160a01b031681565b6040516001600160a01b039091168152602001610071565b6100b86100b336600461051e565b61014b565b005b6000546001600160a01b031661008d565b60606040518060600160405280856001600160a01b0316815260200184848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505050602091820152604051610133929101610559565b60405160208183030381529060405290509392505050565b6101548161027d565b610161602082018261059e565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055600061019d60808301836105b9565b8101906101aa9190610690565b905060008082600001516001600160a01b03165a84602001516040516101d09190610764565b60006040518083038160008787f1925050503d806000811461020e576040519150601f19603f3d011682016040523d82523d6000602084013e610213565b606091505b50915091508161025a57806040517fa5fa8d2b0000000000000000000000000000000000000000000000000000000081526004016102519190610504565b60405180910390fd5b50506001805473ffffffffffffffffffffffffffffffffffffffff191690555050565b6000546040517f33a88c720000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906333a88c72906102c69084906004016107ce565b60206040518083038186803b1580156102de57600080fd5b505afa1580156102f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031691906108b2565b6103625760405162461bcd60e51b815260206004820152601f60248201527f4d657373616765206e6f7420666f756e64206f722066696e616c697a65642e006044820152606401610251565b60008160405160200161037591906107ce565b60408051601f1981840301815291815281516020928301206000818152600290935291205490915060ff16156103ed5760405162461bcd60e51b815260206004820152601960248201527f4d65737361676520616c726561647920636f6e73756d65642e000000000000006044820152606401610251565b6000908152600260205260409020805460ff1916600117905550565b80356001600160a01b038116811461042057600080fd5b919050565b60008060006040848603121561043a57600080fd5b61044384610409565b9250602084013567ffffffffffffffff8082111561046057600080fd5b818601915086601f83011261047457600080fd5b81358181111561048357600080fd5b87602082850101111561049557600080fd5b6020830194508093505050509250925092565b60005b838110156104c35781810151838201526020016104ab565b838111156104d2576000848401525b50505050565b600081518084526104f08160208601602086016104a8565b601f01601f19169290920160200192915050565b60208152600061051760208301846104d8565b9392505050565b60006020828403121561053057600080fd5b813567ffffffffffffffff81111561054757600080fd5b820160c0818503121561051757600080fd5b602081526001600160a01b038251166020820152600060208301516060604084015261058860808401826104d8565b9050604084015160608401528091505092915050565b6000602082840312156105b057600080fd5b61051782610409565b6000808335601e198436030181126105d057600080fd5b83018035915067ffffffffffffffff8211156105eb57600080fd5b60200191503681900382131561060057600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff8111828210171561065957610659610607565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561068857610688610607565b604052919050565b600060208083850312156106a357600080fd5b823567ffffffffffffffff808211156106bb57600080fd5b90840190606082870312156106cf57600080fd5b6106d7610636565b6106e083610409565b815283830135828111156106f357600080fd5b8301601f8101881361070457600080fd5b80358381111561071657610716610607565b610728601f8201601f1916870161065f565b9350808452888682840101111561073e57600080fd5b808683018786013760009084018601525092830152604090810135908201529392505050565b600082516107768184602087016104a8565b9190910192915050565b803563ffffffff8116811461042057600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b803560ff8116811461042057600080fd5b602081526001600160a01b036107e383610409565b1660208201526000602083013567ffffffffffffffff80821680831461080857600080fd5b8060408601525061081b60408601610780565b915063ffffffff80831660608601528061083760608801610780565b1660808601525060808501359150601e1985360301821261085757600080fd5b9084019081358181111561086a57600080fd5b80360386131561087957600080fd5b60c060a086015261089160e086018260208601610794565b925050506108a160a085016107bd565b60ff811660c0850152509392505050565b6000602082840312156108c457600080fd5b8151801515811461051757600080fdfea264697066735822122071fa1493675e54327867ea735dcb5aade210e1ddfbd9c82187a8baab562a574064736f6c63430008090033",
}

// CrossChainMessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainMessengerMetaData.ABI instead.
var CrossChainMessengerABI = CrossChainMessengerMetaData.ABI

// CrossChainMessengerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainMessengerMetaData.Bin instead.
var CrossChainMessengerBin = CrossChainMessengerMetaData.Bin

// DeployCrossChainMessenger deploys a new Ethereum contract, binding an instance of CrossChainMessenger to it.
func DeployCrossChainMessenger(auth *bind.TransactOpts, backend bind.ContractBackend, messageBusAddr common.Address) (common.Address, *types.Transaction, *CrossChainMessenger, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainMessengerBin), backend, messageBusAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// CrossChainMessenger is an auto generated Go binding around an Ethereum contract.
type CrossChainMessenger struct {
	CrossChainMessengerCaller     // Read-only binding to the contract
	CrossChainMessengerTransactor // Write-only binding to the contract
	CrossChainMessengerFilterer   // Log filterer for contract events
}

// CrossChainMessengerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainMessengerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainMessengerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainMessengerSession struct {
	Contract     *CrossChainMessenger // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CrossChainMessengerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainMessengerCallerSession struct {
	Contract *CrossChainMessengerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CrossChainMessengerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainMessengerTransactorSession struct {
	Contract     *CrossChainMessengerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CrossChainMessengerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainMessengerRaw struct {
	Contract *CrossChainMessenger // Generic contract binding to access the raw methods on
}

// CrossChainMessengerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainMessengerCallerRaw struct {
	Contract *CrossChainMessengerCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainMessengerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactorRaw struct {
	Contract *CrossChainMessengerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainMessenger creates a new instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessenger(address common.Address, backend bind.ContractBackend) (*CrossChainMessenger, error) {
	contract, err := bindCrossChainMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// NewCrossChainMessengerCaller creates a new read-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerCaller(address common.Address, caller bind.ContractCaller) (*CrossChainMessengerCaller, error) {
	contract, err := bindCrossChainMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerCaller{contract: contract}, nil
}

// NewCrossChainMessengerTransactor creates a new write-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainMessengerTransactor, error) {
	contract, err := bindCrossChainMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerTransactor{contract: contract}, nil
}

// NewCrossChainMessengerFilterer creates a new log filterer instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainMessengerFilterer, error) {
	contract, err := bindCrossChainMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerFilterer{contract: contract}, nil
}

// bindCrossChainMessenger binds a generic wrapper to an already deployed contract.
func bindCrossChainMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CrossChainMessengerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.CrossChainMessengerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transact(opts, method, params...)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCaller) EncodeCall(opts *bind.CallOpts, target common.Address, payload []byte) ([]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "encodeCall", target, payload)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x9b7cf1ee.
//
// Solidity: function relayMessage((address,uint64,uint32,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessage(opts *bind.TransactOpts, message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessage", message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x9b7cf1ee.
//
// Solidity: function relayMessage((address,uint64,uint32,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x9b7cf1ee.
//
// Solidity: function relayMessage((address,uint64,uint32,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}
