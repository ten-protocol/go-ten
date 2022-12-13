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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messageBusAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"crossChainSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600180546001600160a01b031916905534801561002057600080fd5b5060405161090e38038061090e83398101604081905261003f91610064565b600080546001600160a01b0319166001600160a01b0392909216919091179055610094565b60006020828403121561007657600080fd5b81516001600160a01b038116811461008d57600080fd5b9392505050565b61086b806100a36000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806363012de5146100465780639b7cf1ee14610075578063a1a227fa1461008a575b600080fd5b600154610059906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b6100886100833660046103aa565b61009b565b005b6000546001600160a01b0316610059565b6100a4816101b9565b6100b16020820182610408565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905560006100ed6080830183610423565b8101906100fa9190610522565b905060008082600001516001600160a01b0316836020015160405161011f919061061b565b6000604051808303816000865af19150503d806000811461015c576040519150601f19603f3d011682016040523d82523d6000602084013e610161565b606091505b5091509150816101965761017481610345565b60405162461bcd60e51b815260040161018d9190610637565b60405180910390fd5b50506001805473ffffffffffffffffffffffffffffffffffffffff191690555050565b6000546040517f33a88c720000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906333a88c72906102029084906004016106b8565b60206040518083038186803b15801561021a57600080fd5b505afa15801561022e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610252919061079c565b61029e5760405162461bcd60e51b815260206004820152601f60248201527f4d657373616765206e6f7420666f756e64206f722066696e616c697a65642e00604482015260640161018d565b6000816040516020016102b191906106b8565b60408051601f1981840301815291815281516020928301206000818152600290935291205490915060ff16156103295760405162461bcd60e51b815260206004820152601960248201527f4d65737361676520616c726561647920636f6e73756d65642e00000000000000604482015260640161018d565b6000908152600260205260409020805460ff1916600117905550565b606060448251101561038a57505060408051808201909152601d81527f5472616e73616374696f6e2072657665727465642073696c656e746c79000000602082015290565b600482019150818060200190518101906103a491906107be565b92915050565b6000602082840312156103bc57600080fd5b813567ffffffffffffffff8111156103d357600080fd5b820160c081850312156103e557600080fd5b9392505050565b80356001600160a01b038116811461040357600080fd5b919050565b60006020828403121561041a57600080fd5b6103e5826103ec565b6000808335601e1984360301811261043a57600080fd5b83018035915067ffffffffffffffff82111561045557600080fd5b60200191503681900382131561046a57600080fd5b9250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156104c3576104c3610471565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156104f2576104f2610471565b604052919050565b600067ffffffffffffffff82111561051457610514610471565b50601f01601f191660200190565b6000602080838503121561053557600080fd5b823567ffffffffffffffff8082111561054d57600080fd5b908401906060828703121561056157600080fd5b6105696104a0565b610572836103ec565b8152838301358281111561058557600080fd5b83019150601f8201871361059857600080fd5b81356105ab6105a6826104fa565b6104c9565b81815288868386010111156105bf57600080fd5b818685018783013760009181018601919091529381019390935250604090810135908201529392505050565b60005b838110156106065781810151838201526020016105ee565b83811115610615576000848401525b50505050565b6000825161062d8184602087016105eb565b9190910192915050565b60208152600082518060208401526106568160408501602087016105eb565b601f01601f19169190910160400192915050565b803563ffffffff8116811461040357600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b803560ff8116811461040357600080fd5b602081526001600160a01b036106cd836103ec565b1660208201526000602083013567ffffffffffffffff8082168083146106f257600080fd5b806040860152506107056040860161066a565b915063ffffffff8083166060860152806107216060880161066a565b1660808601525060808501359150601e1985360301821261074157600080fd5b9084019081358181111561075457600080fd5b80360386131561076357600080fd5b60c060a086015261077b60e08601826020860161067e565b9250505061078b60a085016106a7565b60ff811660c0850152509392505050565b6000602082840312156107ae57600080fd5b815180151581146103e557600080fd5b6000602082840312156107d057600080fd5b815167ffffffffffffffff8111156107e757600080fd5b8201601f810184136107f857600080fd5b80516108066105a6826104fa565b81815285602083850101111561081b57600080fd5b61082c8260208301602086016105eb565b9594505050505056fea26469706673582212204cd6092dd7776c554d259ad2a9e54d530369ed8ab89307d9056218f4be0c064c64736f6c63430008090033",
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
