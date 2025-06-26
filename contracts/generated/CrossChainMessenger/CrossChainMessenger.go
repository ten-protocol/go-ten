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
	_ = abi.ConvertType
)

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint64
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// CrossChainMessengerMetaData contains all meta data concerning the CrossChainMessenger contract.
var CrossChainMessengerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"crossChainSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainTarget\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"encodeCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messageBusAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContract\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"messageConsumed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"messageConsumed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"relayMessageWithProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506110848061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c806363012de511610072578063a1a227fa11610058578063a1a227fa14610162578063b859ce8314610172578063c4d66de814610185575f5ffd5b806363012de514610123578063772c655214610143575f5ffd5b80634c81bd20146100a357806350676272146100b8578063530c1e40146100cb5780635b76f28b14610103575b5f5ffd5b6100b66100b13660046107ba565b610198565b005b6100b66100c6366004610862565b6102b2565b6100ed6100d93660046108e5565b60036020525f908152604090205460ff1681565b6040516100fa919061090c565b60405180910390f35b610116610111366004610985565b6103a7565b6040516100fa9190610a16565b600154610136906001600160a01b031681565b6040516100fa9190610a30565b5f54610155906001600160a01b031681565b6040516100fa9190610a5b565b5f546001600160a01b0316610136565b600254610136906001600160a01b031681565b6100b6610193366004610a69565b610427565b6101a18161056f565b6101ae6020820182610a69565b600180546001600160a01b0319166001600160a01b03929092169190911790555f6101dc6080830183610a86565b8101906101e99190610c49565b8051600280546001600160a01b0319166001600160a01b0390921691821790559091505f9081905a84602001516040516102239190610ca2565b5f604051808303815f8787f1925050503d805f811461025d576040519150601f19603f3d011682016040523d82523d5f602084013e610262565b606091505b509150915081610290578060405163a5fa8d2b60e01b81526004016102879190610a16565b60405180910390fd5b5050600180546001600160a01b03199081169091556002805490911690555050565b6102be84848484610686565b6102cb6020850185610a69565b600180546001600160a01b0319166001600160a01b03929092169190911790555f6102f96080860186610a86565b8101906103069190610c49565b90505f5f825f01516001600160a01b03165a846020015160405161032a9190610ca2565b5f604051808303815f8787f1925050503d805f8114610364576040519150601f19603f3d011682016040523d82523d5f602084013e610369565b606091505b50915091508161038e578060405163a5fa8d2b60e01b81526004016102879190610a16565b5050600180546001600160a01b03191690555050505050565b60606040518060600160405280856001600160a01b0316815260200184848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92018290525093855250505060209182015260405161040e929101610cfa565b60405160208183030381529060405290505b9392505050565b5f610430610777565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f8115801561045c5750825b90505f8267ffffffffffffffff1660011480156104785750303b155b905081158015610486575080155b156104bd576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156104f157845468ff00000000000000001916680100000000000000001785555b5f80546001600160a01b0388166001600160a01b031991821617909155600180549091169055831561056757845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061055e90600190610d25565b60405180910390a15b505050505050565b5f546040517f91643fdd0000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906391643fdd906105b7908490600401610ef0565b602060405180830381865afa1580156105d2573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105f69190610f14565b6106125760405162461bcd60e51b815260040161028790610f65565b5f816040516020016106249190610ef0565b60408051601f1981840301815291815281516020928301205f818152600390935291205490915060ff161561066b5760405162461bcd60e51b815260040161028790610fa7565b5f908152600360205260409020805460ff1916600117905550565b5f546040517fce0d7db30000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063ce0d7db3906106d4908790879087908790600401611010565b5f6040518083038186803b1580156106ea575f5ffd5b505afa1580156106fc573d5f5f3e3d5ffd5b505050505f846040516020016107129190610ef0565b60408051601f1981840301815291815281516020928301205f818152600390935291205490915060ff16156107595760405162461bcd60e51b815260040161028790610fa7565b5f908152600360205260409020805460ff1916600117905550505050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b5f60c082840312156107b4576107b45f5ffd5b50919050565b5f602082840312156107cd576107cd5f5ffd5b813567ffffffffffffffff8111156107e6576107e65f5ffd5b6107f2848285016107a1565b949350505050565b5f5f83601f84011261080d5761080d5f5ffd5b50813567ffffffffffffffff811115610827576108275f5ffd5b602083019150836020820283011115610841576108415f5ffd5b9250929050565b805b8114610854575f5ffd5b50565b803561079b81610848565b5f5f5f5f60608587031215610878576108785f5ffd5b843567ffffffffffffffff811115610891576108915f5ffd5b61089d878288016107a1565b945050602085013567ffffffffffffffff8111156108bc576108bc5f5ffd5b6108c8878288016107fa565b93509350506108da8660408701610857565b905092959194509250565b5f602082840312156108f8576108f85f5ffd5b6104208383610857565b8015155b82525050565b6020810161079b8284610902565b5f6001600160a01b03821661079b565b61084a8161091a565b803561079b8161092a565b5f5f83601f840112610951576109515f5ffd5b50813567ffffffffffffffff81111561096b5761096b5f5ffd5b602083019150836001820283011115610841576108415f5ffd5b5f5f5f6040848603121561099a5761099a5f5ffd5b6109a48585610933565b9250602084013567ffffffffffffffff8111156109c2576109c25f5ffd5b6109ce8682870161093e565b92509250509250925092565b8281835e505f910152565b5f6109ee825190565b808452602084019350610a058185602086016109da565b601f01601f19169290920192915050565b6020808252810161042081846109e5565b6109068161091a565b6020810161079b8284610a27565b5f61079b8261091a565b5f61079b82610a3e565b61090681610a48565b6020810161079b8284610a52565b5f60208284031215610a7c57610a7c5f5ffd5b6104208383610933565b5f808335601e1936859003018112610a9f57610a9f5f5ffd5b8301915050803567ffffffffffffffff811115610abd57610abd5f5ffd5b602082019150600181023603821315610841576108415f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715610b2a57610b2a610ad7565b6040525050565b5f610b3b60405190565b9050610b478282610b04565b919050565b5f67ffffffffffffffff821115610b6557610b65610ad7565b601f19601f83011660200192915050565b82818337505f910152565b5f610b93610b8e84610b4c565b610b31565b9050828152838383011115610ba957610ba95f5ffd5b610420836020830184610b76565b5f82601f830112610bc957610bc95f5ffd5b61042083833560208501610b81565b5f60608284031215610beb57610beb5f5ffd5b610bf56060610b31565b9050610c018383610933565b8152602082013567ffffffffffffffff811115610c1f57610c1f5f5ffd5b610c2b84828501610bb7565b602083015250610c3e8360408401610857565b604082015292915050565b5f60208284031215610c5c57610c5c5f5ffd5b813567ffffffffffffffff811115610c7557610c755f5ffd5b6107f284828501610bd8565b5f610c8a825190565b610c988185602086016109da565b9290920192915050565b61079b8183610c81565b80610906565b80515f906060840190610cc58582610a27565b5060208301518482036020860152610cdd82826109e5565b9150506040830151610cf26040860182610cac565b509392505050565b602080825281016104208184610cb2565b5f67ffffffffffffffff821661079b565b61090681610d0b565b6020810161079b8284610d1c565b505f61079b6020830183610933565b67ffffffffffffffff811661084a565b803561079b81610d42565b505f61079b6020830183610d52565b67ffffffffffffffff8116610906565b63ffffffff811661084a565b803561079b81610d7c565b505f61079b6020830183610d88565b63ffffffff8116610906565b5f808335601e1936859003018112610dc757610dc75f5ffd5b830160208101925035905067ffffffffffffffff811115610de957610de95f5ffd5b36819003821315610841576108415f5ffd5b818352602083019250610e0f828483610b76565b50601f01601f19160190565b60ff811661084a565b803561079b81610e1b565b505f61079b6020830183610e24565b60ff8116610906565b5f60c08301610e568380610d33565b610e608582610a27565b50610e6e6020840184610d5d565b610e7b6020860182610d6c565b50610e896040840184610d5d565b610e966040860182610d6c565b50610ea46060840184610d93565b610eb16060860182610da2565b50610ebf6080840184610dae565b8583036080870152610ed2838284610dfb565b92505050610ee360a0840184610e2f565b610cf260a0860182610e3e565b602080825281016104208184610e47565b80151561084a565b805161079b81610f01565b5f60208284031215610f2757610f275f5ffd5b6104208383610f09565b601f8152602081017f4d657373616765206e6f7420666f756e64206f722066696e616c697a65642e00815290505b60200190565b6020808252810161079b81610f31565b60198152602081017f4d65737361676520616c726561647920636f6e73756d65642e0000000000000081529050610f5f565b6020808252810161079b81610f75565b82818337505050565b8183526020830192505f7f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610ff957610ff95f5ffd5b60208302925061100a838584610fb7565b50500190565b606080825281016110218187610e47565b90508181036020830152611036818587610fc0565b90506110456040830184610cac565b9594505050505056fea2646970667358221220a37c93eed5a175346bb40eb5882ee7d12fb5b3bb3ce99b9e433255ffcb6f64e264736f6c634300081c0033",
}

// CrossChainMessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainMessengerMetaData.ABI instead.
var CrossChainMessengerABI = CrossChainMessengerMetaData.ABI

// CrossChainMessengerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainMessengerMetaData.Bin instead.
var CrossChainMessengerBin = CrossChainMessengerMetaData.Bin

// DeployCrossChainMessenger deploys a new Ethereum contract, binding an instance of CrossChainMessenger to it.
func DeployCrossChainMessenger(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainMessenger, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainMessengerBin), backend)
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
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainTarget(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainTarget")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
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

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBusContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBusContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageConsumed(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageConsumed", messageHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) Initialize(opts *bind.TransactOpts, messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "initialize", messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessage(opts *bind.TransactOpts, message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessage", message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessageWithProof(opts *bind.TransactOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessageWithProof", message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// CrossChainMessengerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CrossChainMessenger contract.
type CrossChainMessengerInitializedIterator struct {
	Event *CrossChainMessengerInitialized // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerInitialized)
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
		it.Event = new(CrossChainMessengerInitialized)
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
func (it *CrossChainMessengerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerInitialized represents a Initialized event raised by the CrossChainMessenger contract.
type CrossChainMessengerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterInitialized(opts *bind.FilterOpts) (*CrossChainMessengerInitializedIterator, error) {

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerInitializedIterator{contract: _CrossChainMessenger.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerInitialized) (event.Subscription, error) {

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerInitialized)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseInitialized(log types.Log) (*CrossChainMessengerInitialized, error) {
	event := new(CrossChainMessengerInitialized)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
