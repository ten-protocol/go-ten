// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Fees

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

// FeesMetaData contains all meta data concerning the Fees contract.
var FeesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"FeeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collectedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"flatFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"eoaOwner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFeeForMessage\",\"type\":\"uint256\"}],\"name\":\"setMessageFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalCollectedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523461002257610011610026565b604051610c1d61018a8239610c1d90f35b5f80fd5b61002a5b6100326100ad565b565b6100419060401c60ff1690565b90565b6100419054610034565b610041905b6001600160401b031690565b610041905461004e565b61004190610053906001600160401b031682565b9061008d6100416100a992610069565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100b6610143565b016100c081610044565b610132576100cd8161005f565b6001600160401b03919082908116036100e4575050565b816101137fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361012d9361007d565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b610041610181565b6100416100416100419290565b6100417ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061014b565b61004161015856fe6080604052600436101561001a575b3615610018575f80fd5b005b5f3560e01c80631a90a219146100b957806323aa2a9d146100b4578063715018a6146100af57806379ba5097146100aa5780638da5cb5b146100a55780639003adfe146100a0578063afe997ea1461009b578063da35a26f14610096578063e30c3978146100915763f2fde38b0361000e57610273565b610244565b61022b565b6101d4565b6101b9565b610184565b610160565b61014b565b61012e565b6100d2565b5f9103126100c857565b5f80fd5b9052565b565b346100c8576100e23660046100be565b6100fd6100ed61029c565b6040519182918290815260200190565b0390f35b805b036100c857565b905035906100d082610101565b906020828203126100c85761012b9161010a565b90565b346100c857610146610141366004610117565b610353565b604051005b346100c85761015b3660046100be565b6103cf565b346100c8576101703660046100be565b6101466103d4565b6001600160a01b031690565b346100c8576101943660046100be565b6100fd61019f61043f565b604051918291826001600160a01b03909116815260200190565b346100c8576101c93660046100be565b6100fd6100ed61049b565b346100c8576101e43660046100be565b61014661061b565b6001600160a01b038116610103565b905035906100d0826101ec565b91906040838203126100c85761012b906020610224828661010a565b94016101fb565b346100c85761014661023e366004610208565b90610898565b346100c8576102543660046100be565b6100fd61019f6108a2565b906020828203126100c85761012b916101fb565b346100c85761014661028636600461025f565b610971565b61012b9081565b61012b905461028b565b61012b5f610292565b6100d0906102b161097a565b610304565b905f19905b9181191691161790565b61012b61012b61012b9290565b906102e261012b6102e9926102c5565b82546102b6565b9055565b9081526040810192916100d09160200152565b0152565b7f5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f19061032f5f610292565b610339825f6102d2565b61034e61034560405190565b928392836102ed565b0390a1565b6100d0906102a5565b61036461097a565b60405162461bcd60e51b815260206004820152603460248201527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60448201527f742072656e6f756e6365206f776e6572736869700000000000000000000000006064820152608490fd5b61035c565b336103dd6108a2565b6103f86001600160a01b0383165b916001600160a01b031690565b03610406576100d0906109d8565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b61012b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b01546001600160a01b031690565b61017861012b61012b926001600160a01b031690565b61012b90610473565b61012b90610489565b6104a430610492565b3190565b6104b061097a565b6100d06105be565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761050757604052565b6104b8565b906100d061051960405190565b92836104e5565b67ffffffffffffffff811161050757602090601f01601f19160190565b9061054f61054a83610520565b61050c565b918252565b3d1561056d576105633d61053d565b903d5f602084013e565b606090565b1561057957565b60405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e642045746865720000000000000000000000006044820152606490fd5b7fb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba075261034e6105eb30610492565b316100ed5f806106046105ff6105ff61043f565b610492565b6040515f9186905af1610615610554565b50610572565b6100d06104a8565b61012b9060401c60ff1690565b61012b9054610623565b61012b905b67ffffffffffffffff1690565b61012b905461063a565b61063f61012b61012b9290565b9067ffffffffffffffff906102bb565b61063f61012b61012b9267ffffffffffffffff1690565b9061069a61012b6102e992610673565b8254610663565b9068ff00000000000000009060401b6102bb565b906106c561012b6102e992151590565b82546106a1565b6100cc90610656565b6020810192916100d091906106cc565b6106ed610a0b565b9081906107096107036106ff84610630565b1590565b9261064c565b936107135f610656565b67ffffffffffffffff86161480610829575b60019561074261073488610656565b9167ffffffffffffffff1690565b149081610801575b155b90816107f8575b506107ce5761077c91836107735f61076a89610656565b9701968761068a565b6107bf57610853565b610784575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916107b35f61034e936106b5565b604051918291826106d5565b6107c986866106b5565b610853565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f610753565b905061074c61080f30610492565b3b61082061081c5f6102c5565b9190565b1491905061074a565b5082610725565b6100cc906102c5565b9160206100d092949361030060408201965f830190610830565b61034e906108817f5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f193610a30565b61088b815f6102d2565b6040519182915f83610839565b906100d0916106e5565b61012b5f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610465565b6100d0906108d861097a565b610903565b906001600160a01b03906102bb565b906108fc61012b6102e992610492565b82546108dd565b61092d817f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c006108ec565b61094161093b6105ff61043f565b91610492565b907f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270061096c60405190565b5f90a3565b6100d0906108cc565b61098261043f565b339061098d826103eb565b036104065750565b919060086102bb9102916109af6001600160a01b03841b90565b921b90565b91906109c561012b6102e993610492565b908354610995565b6100d0915f916109b4565b6100d090610a065f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c006109cd565b610a39565b61012b610acf565b6100d090610a1f610ad7565b610a2890610b27565b6100d0610b38565b6100d090610a13565b610a7b61093b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993006105ff84610a7583546001600160a01b031690565b926108ec565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e061096c60405190565b61012b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006102c5565b61012b610aa6565b610ae26106ff610b40565b610ae857565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b6100d090610b1e610ad7565b6100d090610bde565b6100d090610b12565b6100d0610ad7565b6100d0610b30565b61012b5f610b4c610a0b565b01610630565b6100d090610b5e610ad7565b610b79565b61017861012b61012b9290565b61012b90610b63565b610b825f610b70565b6001600160a01b0381166001600160a01b03831614610ba557506100d0906109d8565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6100d090610b5256fea26469706673582212201af5618100e730c17b36a81efc23779f7d0d317dd6eb3770a3d454d420f76a0064736f6c634300081c0033",
}

// FeesABI is the input ABI used to generate the binding from.
// Deprecated: Use FeesMetaData.ABI instead.
var FeesABI = FeesMetaData.ABI

// FeesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeesMetaData.Bin instead.
var FeesBin = FeesMetaData.Bin

// DeployFees deploys a new Ethereum contract, binding an instance of Fees to it.
func DeployFees(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Fees, error) {
	parsed, err := FeesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Fees{FeesCaller: FeesCaller{contract: contract}, FeesTransactor: FeesTransactor{contract: contract}, FeesFilterer: FeesFilterer{contract: contract}}, nil
}

// Fees is an auto generated Go binding around an Ethereum contract.
type Fees struct {
	FeesCaller     // Read-only binding to the contract
	FeesTransactor // Write-only binding to the contract
	FeesFilterer   // Log filterer for contract events
}

// FeesCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeesSession struct {
	Contract     *Fees             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeesCallerSession struct {
	Contract *FeesCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FeesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeesTransactorSession struct {
	Contract     *FeesTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeesRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeesRaw struct {
	Contract *Fees // Generic contract binding to access the raw methods on
}

// FeesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeesCallerRaw struct {
	Contract *FeesCaller // Generic read-only contract binding to access the raw methods on
}

// FeesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeesTransactorRaw struct {
	Contract *FeesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFees creates a new instance of Fees, bound to a specific deployed contract.
func NewFees(address common.Address, backend bind.ContractBackend) (*Fees, error) {
	contract, err := bindFees(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fees{FeesCaller: FeesCaller{contract: contract}, FeesTransactor: FeesTransactor{contract: contract}, FeesFilterer: FeesFilterer{contract: contract}}, nil
}

// NewFeesCaller creates a new read-only instance of Fees, bound to a specific deployed contract.
func NewFeesCaller(address common.Address, caller bind.ContractCaller) (*FeesCaller, error) {
	contract, err := bindFees(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeesCaller{contract: contract}, nil
}

// NewFeesTransactor creates a new write-only instance of Fees, bound to a specific deployed contract.
func NewFeesTransactor(address common.Address, transactor bind.ContractTransactor) (*FeesTransactor, error) {
	contract, err := bindFees(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeesTransactor{contract: contract}, nil
}

// NewFeesFilterer creates a new log filterer instance of Fees, bound to a specific deployed contract.
func NewFeesFilterer(address common.Address, filterer bind.ContractFilterer) (*FeesFilterer, error) {
	contract, err := bindFees(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeesFilterer{contract: contract}, nil
}

// bindFees binds a generic wrapper to an already deployed contract.
func bindFees(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fees *FeesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fees.Contract.FeesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fees *FeesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.Contract.FeesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fees *FeesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fees.Contract.FeesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fees *FeesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fees.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fees *FeesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fees *FeesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fees.Contract.contract.Transact(opts, method, params...)
}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesCaller) CollectedFees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "collectedFees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesSession) CollectedFees() (*big.Int, error) {
	return _Fees.Contract.CollectedFees(&_Fees.CallOpts)
}

// CollectedFees is a free data retrieval call binding the contract method 0x9003adfe.
//
// Solidity: function collectedFees() view returns(uint256)
func (_Fees *FeesCallerSession) CollectedFees() (*big.Int, error) {
	return _Fees.Contract.CollectedFees(&_Fees.CallOpts)
}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesCaller) MessageFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "messageFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesSession) MessageFee() (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts)
}

// MessageFee is a free data retrieval call binding the contract method 0x1a90a219.
//
// Solidity: function messageFee() view returns(uint256)
func (_Fees *FeesCallerSession) MessageFee() (*big.Int, error) {
	return _Fees.Contract.MessageFee(&_Fees.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesSession) Owner() (common.Address, error) {
	return _Fees.Contract.Owner(&_Fees.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fees *FeesCallerSession) Owner() (common.Address, error) {
	return _Fees.Contract.Owner(&_Fees.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesSession) PendingOwner() (common.Address, error) {
	return _Fees.Contract.PendingOwner(&_Fees.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fees *FeesCallerSession) PendingOwner() (common.Address, error) {
	return _Fees.Contract.PendingOwner(&_Fees.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Fees.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesSession) RenounceOwnership() error {
	return _Fees.Contract.RenounceOwnership(&_Fees.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_Fees *FeesCallerSession) RenounceOwnership() error {
	return _Fees.Contract.RenounceOwnership(&_Fees.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesSession) AcceptOwnership() (*types.Transaction, error) {
	return _Fees.Contract.AcceptOwnership(&_Fees.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fees *FeesTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Fees.Contract.AcceptOwnership(&_Fees.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesTransactor) Initialize(opts *bind.TransactOpts, flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "initialize", flatFee, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesSession) Initialize(flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, flatFee, eoaOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 flatFee, address eoaOwner) returns()
func (_Fees *FeesTransactorSession) Initialize(flatFee *big.Int, eoaOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.Initialize(&_Fees.TransactOpts, flatFee, eoaOwner)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesTransactor) SetMessageFee(opts *bind.TransactOpts, newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "setMessageFee", newFeeForMessage)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesSession) SetMessageFee(newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newFeeForMessage)
}

// SetMessageFee is a paid mutator transaction binding the contract method 0x23aa2a9d.
//
// Solidity: function setMessageFee(uint256 newFeeForMessage) returns()
func (_Fees *FeesTransactorSession) SetMessageFee(newFeeForMessage *big.Int) (*types.Transaction, error) {
	return _Fees.Contract.SetMessageFee(&_Fees.TransactOpts, newFeeForMessage)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.TransferOwnership(&_Fees.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fees *FeesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fees.Contract.TransferOwnership(&_Fees.TransactOpts, newOwner)
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesTransactor) WithdrawalCollectedFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.Transact(opts, "withdrawalCollectedFees")
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesSession) WithdrawalCollectedFees() (*types.Transaction, error) {
	return _Fees.Contract.WithdrawalCollectedFees(&_Fees.TransactOpts)
}

// WithdrawalCollectedFees is a paid mutator transaction binding the contract method 0xafe997ea.
//
// Solidity: function withdrawalCollectedFees() returns()
func (_Fees *FeesTransactorSession) WithdrawalCollectedFees() (*types.Transaction, error) {
	return _Fees.Contract.WithdrawalCollectedFees(&_Fees.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fees.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesSession) Receive() (*types.Transaction, error) {
	return _Fees.Contract.Receive(&_Fees.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fees *FeesTransactorSession) Receive() (*types.Transaction, error) {
	return _Fees.Contract.Receive(&_Fees.TransactOpts)
}

// FeesFeeChangedIterator is returned from FilterFeeChanged and is used to iterate over the raw logs and unpacked data for FeeChanged events raised by the Fees contract.
type FeesFeeChangedIterator struct {
	Event *FeesFeeChanged // Event containing the contract specifics and raw log

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
func (it *FeesFeeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesFeeChanged)
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
		it.Event = new(FeesFeeChanged)
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
func (it *FeesFeeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesFeeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesFeeChanged represents a FeeChanged event raised by the Fees contract.
type FeesFeeChanged struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeChanged is a free log retrieval operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) FilterFeeChanged(opts *bind.FilterOpts) (*FeesFeeChangedIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "FeeChanged")
	if err != nil {
		return nil, err
	}
	return &FeesFeeChangedIterator{contract: _Fees.contract, event: "FeeChanged", logs: logs, sub: sub}, nil
}

// WatchFeeChanged is a free log subscription operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) WatchFeeChanged(opts *bind.WatchOpts, sink chan<- *FeesFeeChanged) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "FeeChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesFeeChanged)
				if err := _Fees.contract.UnpackLog(event, "FeeChanged", log); err != nil {
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

// ParseFeeChanged is a log parse operation binding the contract event 0x5fc463da23c1b063e66f9e352006a7fbe8db7223c455dc429e881a2dfe2f94f1.
//
// Solidity: event FeeChanged(uint256 oldFee, uint256 newFee)
func (_Fees *FeesFilterer) ParseFeeChanged(log types.Log) (*FeesFeeChanged, error) {
	event := new(FeesFeeChanged)
	if err := _Fees.contract.UnpackLog(event, "FeeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesFeeWithdrawnIterator is returned from FilterFeeWithdrawn and is used to iterate over the raw logs and unpacked data for FeeWithdrawn events raised by the Fees contract.
type FeesFeeWithdrawnIterator struct {
	Event *FeesFeeWithdrawn // Event containing the contract specifics and raw log

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
func (it *FeesFeeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesFeeWithdrawn)
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
		it.Event = new(FeesFeeWithdrawn)
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
func (it *FeesFeeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesFeeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesFeeWithdrawn represents a FeeWithdrawn event raised by the Fees contract.
type FeesFeeWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeWithdrawn is a free log retrieval operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) FilterFeeWithdrawn(opts *bind.FilterOpts) (*FeesFeeWithdrawnIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return &FeesFeeWithdrawnIterator{contract: _Fees.contract, event: "FeeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFeeWithdrawn is a free log subscription operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) WatchFeeWithdrawn(opts *bind.WatchOpts, sink chan<- *FeesFeeWithdrawn) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "FeeWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesFeeWithdrawn)
				if err := _Fees.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
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

// ParseFeeWithdrawn is a log parse operation binding the contract event 0xb7eeacba6b133788365610e83d3f130d07b6ef6e78877961f25b3f61fcba0752.
//
// Solidity: event FeeWithdrawn(uint256 amount)
func (_Fees *FeesFilterer) ParseFeeWithdrawn(log types.Log) (*FeesFeeWithdrawn, error) {
	event := new(FeesFeeWithdrawn)
	if err := _Fees.contract.UnpackLog(event, "FeeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Fees contract.
type FeesInitializedIterator struct {
	Event *FeesInitialized // Event containing the contract specifics and raw log

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
func (it *FeesInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesInitialized)
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
		it.Event = new(FeesInitialized)
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
func (it *FeesInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesInitialized represents a Initialized event raised by the Fees contract.
type FeesInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fees *FeesFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeesInitializedIterator, error) {

	logs, sub, err := _Fees.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeesInitializedIterator{contract: _Fees.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fees *FeesFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeesInitialized) (event.Subscription, error) {

	logs, sub, err := _Fees.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesInitialized)
				if err := _Fees.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Fees *FeesFilterer) ParseInitialized(log types.Log) (*FeesInitialized, error) {
	event := new(FeesInitialized)
	if err := _Fees.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Fees contract.
type FeesOwnershipTransferStartedIterator struct {
	Event *FeesOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *FeesOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesOwnershipTransferStarted)
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
		it.Event = new(FeesOwnershipTransferStarted)
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
func (it *FeesOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Fees contract.
type FeesOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeesOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeesOwnershipTransferStartedIterator{contract: _Fees.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *FeesOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesOwnershipTransferStarted)
				if err := _Fees.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) ParseOwnershipTransferStarted(log types.Log) (*FeesOwnershipTransferStarted, error) {
	event := new(FeesOwnershipTransferStarted)
	if err := _Fees.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Fees contract.
type FeesOwnershipTransferredIterator struct {
	Event *FeesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FeesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeesOwnershipTransferred)
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
		it.Event = new(FeesOwnershipTransferred)
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
func (it *FeesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeesOwnershipTransferred represents a OwnershipTransferred event raised by the Fees contract.
type FeesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeesOwnershipTransferredIterator{contract: _Fees.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fees.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeesOwnershipTransferred)
				if err := _Fees.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fees *FeesFilterer) ParseOwnershipTransferred(log types.Log) (*FeesOwnershipTransferred, error) {
	event := new(FeesOwnershipTransferred)
	if err := _Fees.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
