// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Structs

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

// StructsTransaction is an auto generated low-level Go binding around an user-defined struct.
type StructsTransaction struct {
	TxType               uint8
	Nonce                *big.Int
	GasPrice             *big.Int
	GasLimit             *big.Int
	To                   common.Address
	Value                *big.Int
	Data                 []byte
	V                    uint8
	R                    [32]byte
	S                    [32]byte
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.Transaction\",\"name\":\"txData\",\"type\":\"tuple\"}],\"name\":\"recoverSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523461001e576040516105e56100248239308150506105e590f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c638314c06f0361003657610085565b90816101809103126100365790565b600080fd5b9060208282031261003657813567ffffffffffffffff8111610036576100619201610027565b90565b6001600160a01b031690565b6001600160a01b03909116815260200190565b565b6100a761009b61009636600461003b565b610212565b60405191829182610070565b0390f35b60ff81165b0361003657565b35610061816100ab565b806100b0565b35610061816100c1565b6001600160a01b0381166100b0565b35610061816100d1565b903590601e193682900301821215610036570180359067ffffffffffffffff8211610036576020019136829003831361003657565b90826000939282370152565b919061014981610142816101539560209181520190565b809561011f565b601f01601f191690565b0190565b956101c0906101b961019e976101a96100839d9f9e9c96986101a26101009d979960408e61019e6101cd9e61019761012084019f600085019060ff169052565b6020830152565b0152565b60608c0152565b6001600160a01b031660808a0152565b60a0880152565b85830360c087015261012b565b9660e0830152565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761020d57604052565b6101d5565b610061906102f7610222826100b7565b6102b5610231602085016100c7565b61023d604086016100c7565b926102a961024d606088016100c7565b8761025a608082016100e0565b61026660a083016100c7565b9061027460c08401846100ea565b92909161029161016061028a61014088016100c7565b96016100c7565b9561029b60405190565b9b8c9a60208c019a8b610157565b908103825203826101eb565b6102c76102c0825190565b9160200190565b207f19457468657265756d205369676e6564204d6573736167653a0a333200000000600052601c52603c60002090565b9061030460e082016100b7565b61031e61012061031761010085016100c7565b93016100c7565b6100619361032e9391929061039f565b90929192610497565b6100616100616100619290565b61006190610337565b6100646100616100619290565b6100619061034d565b61019e6100839461038c606094989795610382608086019a6000870152565b60ff166020850152565b6040830152565b6040513d6000823e3d90fd5b90916103aa84610344565b6103da6103d67f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0610337565b9190565b1161045457906103fc602094600094936103f360405190565b94859485610363565b838052039060015afa1561044f576000516000916104198361035a565b6001600160a01b0381166001600160a01b03841614610442575061043c83610337565b91929190565b915061043c600193610337565b610393565b505050610461600061035a565b9160039190565b634e487b7160e01b600052602160045260246000fd5b6004111561048857565b610468565b906100838261047e565b6104a1600061048d565b6104aa8261048d565b036104b3575050565b6104bd600161048d565b6104c68261048d565b036104fa576040517ff645eedf000000000000000000000000000000000000000000000000000000008152600490fd5b0390fd5b610504600261048d565b61050d8261048d565b03610554576104f661051e83610344565b6040519182917ffce698f70000000000000000000000000000000000000000000000000000000083526004830190815260200190565b610567610561600361048d565b9161048d565b1461056f5750565b6104f69061057c60405190565b9182917fd78bce0c000000000000000000000000000000000000000000000000000000008352600483019081526020019056fea264697066735822122094db7b9a0a1dbf0c452203de6f7e97cba3e05fc1def5e3ae5efbdebc932f7da564736f6c63430008140033",
}

// StructsABI is the input ABI used to generate the binding from.
// Deprecated: Use StructsMetaData.ABI instead.
var StructsABI = StructsMetaData.ABI

// StructsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StructsMetaData.Bin instead.
var StructsBin = StructsMetaData.Bin

// DeployStructs deploys a new Ethereum contract, binding an instance of Structs to it.
func DeployStructs(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Structs, error) {
	parsed, err := StructsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StructsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Structs{StructsCaller: StructsCaller{contract: contract}, StructsTransactor: StructsTransactor{contract: contract}, StructsFilterer: StructsFilterer{contract: contract}}, nil
}

// Structs is an auto generated Go binding around an Ethereum contract.
type Structs struct {
	StructsCaller     // Read-only binding to the contract
	StructsTransactor // Write-only binding to the contract
	StructsFilterer   // Log filterer for contract events
}

// StructsCaller is an auto generated read-only Go binding around an Ethereum contract.
type StructsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StructsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StructsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StructsSession struct {
	Contract     *Structs          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StructsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StructsCallerSession struct {
	Contract *StructsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StructsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StructsTransactorSession struct {
	Contract     *StructsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StructsRaw is an auto generated low-level Go binding around an Ethereum contract.
type StructsRaw struct {
	Contract *Structs // Generic contract binding to access the raw methods on
}

// StructsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StructsCallerRaw struct {
	Contract *StructsCaller // Generic read-only contract binding to access the raw methods on
}

// StructsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StructsTransactorRaw struct {
	Contract *StructsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStructs creates a new instance of Structs, bound to a specific deployed contract.
func NewStructs(address common.Address, backend bind.ContractBackend) (*Structs, error) {
	contract, err := bindStructs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Structs{StructsCaller: StructsCaller{contract: contract}, StructsTransactor: StructsTransactor{contract: contract}, StructsFilterer: StructsFilterer{contract: contract}}, nil
}

// NewStructsCaller creates a new read-only instance of Structs, bound to a specific deployed contract.
func NewStructsCaller(address common.Address, caller bind.ContractCaller) (*StructsCaller, error) {
	contract, err := bindStructs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StructsCaller{contract: contract}, nil
}

// NewStructsTransactor creates a new write-only instance of Structs, bound to a specific deployed contract.
func NewStructsTransactor(address common.Address, transactor bind.ContractTransactor) (*StructsTransactor, error) {
	contract, err := bindStructs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StructsTransactor{contract: contract}, nil
}

// NewStructsFilterer creates a new log filterer instance of Structs, bound to a specific deployed contract.
func NewStructsFilterer(address common.Address, filterer bind.ContractFilterer) (*StructsFilterer, error) {
	contract, err := bindStructs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StructsFilterer{contract: contract}, nil
}

// bindStructs binds a generic wrapper to an already deployed contract.
func bindStructs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StructsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Structs *StructsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Structs.Contract.StructsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Structs *StructsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Structs.Contract.StructsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Structs *StructsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Structs.Contract.StructsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Structs *StructsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Structs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Structs *StructsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Structs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Structs *StructsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Structs.Contract.contract.Transact(opts, method, params...)
}

// RecoverSender is a free data retrieval call binding the contract method 0x0cb57856.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256) txData) pure returns(address sender)
func (_Structs *StructsCaller) RecoverSender(opts *bind.CallOpts, txData StructsTransaction) (common.Address, error) {
	var out []interface{}
	err := _Structs.contract.Call(opts, &out, "recoverSender", txData)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSender is a free data retrieval call binding the contract method 0x0cb57856.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256) txData) pure returns(address sender)
func (_Structs *StructsSession) RecoverSender(txData StructsTransaction) (common.Address, error) {
	return _Structs.Contract.RecoverSender(&_Structs.CallOpts, txData)
}

// RecoverSender is a free data retrieval call binding the contract method 0x0cb57856.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256) txData) pure returns(address sender)
func (_Structs *StructsCallerSession) RecoverSender(txData StructsTransaction) (common.Address, error) {
	return _Structs.Contract.RecoverSender(&_Structs.CallOpts, txData)
}
