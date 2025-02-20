// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ManagementContract

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

// ManagementContractMetaData contains all meta data concerning the ManagementContract contract.
var ManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"ImportantContractAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"messageBusAddress\",\"type\":\"address\"}],\"name\":\"LogManagementContractCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GetImportantContractKeys\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"SetImportantContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"importantContractAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"importantContractKeys\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610c2e806100975f395ff3fe608060405234801561000f575f5ffd5b5060043610610085575f3560e01c80638129fc1c116100585780638129fc1c146101055780638da5cb5b1461010d57806398077e861461013d578063f2fde38b1461015d575f5ffd5b806303e72e48146100895780633e60a22f1461009e5780636a30d26c146100e8578063715018a6146100fd575b5f5ffd5b61009c61009736600461082c565b610170565b005b6100d26100ac36600461087e565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b6040516100df91906108cd565b60405180910390f35b6100f0610273565b6040516100df9190610984565b61009c610346565b61009c610359565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166100d2565b61015061014b3660046109a6565b610498565b6040516100df91906109c3565b61009c61016b3660046109d4565b61053d565b61017861059c565b5f6001600160a01b03166001836040516101929190610a12565b908152604051908190036020019020546001600160a01b0316036101ea575f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563016101e88382610af4565b505b806001836040516101fb9190610a12565b90815260405190819003602001812080546001600160a01b039390931673ffffffffffffffffffffffffffffffffffffffff19909316929092179091557f17b2f9f5748931099ffee882b5b64f4a560b5c55da9b4f4e396dae3bb9f98cb5906102679084908490610bb0565b60405180910390a15050565b60605f805480602002602001604051908101604052809291908181526020015f905b8282101561033d578382905f5260205f200180546102b290610a30565b80601f01602080910402602001604051908101604052809291908181526020018280546102de90610a30565b80156103295780601f1061030057610100808354040283529160200191610329565b820191905f5260205f20905b81548152906001019060200180831161030c57829003601f168201915b505050505081526020019060010190610295565b50505050905090565b61034e61059c565b6103575f610610565b565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156103a35750825b90505f8267ffffffffffffffff1660011480156103bf5750303b155b9050811580156103cd575080155b15610404576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561043857845468ff00000000000000001916680100000000000000001785555b6104413361068d565b831561049157845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061048890600190610bea565b60405180910390a15b5050505050565b5f81815481106104a6575f80fd5b905f5260205f20015f9150905080546104be90610a30565b80601f01602080910402602001604051908101604052809291908181526020018280546104ea90610a30565b80156105355780601f1061050c57610100808354040283529160200191610535565b820191905f5260205f20905b81548152906001019060200180831161051857829003601f168201915b505050505081565b61054561059c565b6001600160a01b038116610590575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161058791906108cd565b60405180910390fd5b61059981610610565b50565b336105ce7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461035757336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161058791906108cd565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b61069561069e565b61059981610705565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610357576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61054561069e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff821117156107475761074761070d565b6040525050565b5f61075860405190565b90506107648282610721565b919050565b5f67ffffffffffffffff8211156107825761078261070d565b601f19601f83011660200192915050565b82818337505f910152565b5f6107b06107ab84610769565b61074e565b90508281528383830111156107c6576107c65f5ffd5b6107d4836020830184610793565b9392505050565b5f82601f8301126107ed576107ed5f5ffd5b6107d48383356020850161079e565b5f6001600160a01b0382165b92915050565b610817816107fc565b8114610599575f5ffd5b80356108088161080e565b5f5f60408385031215610840576108405f5ffd5b823567ffffffffffffffff811115610859576108595f5ffd5b610865858286016107db565b9250506108758460208501610821565b90509250929050565b5f60208284031215610891576108915f5ffd5b813567ffffffffffffffff8111156108aa576108aa5f5ffd5b6108b6848285016107db565b949350505050565b6108c7816107fc565b82525050565b6020810161080882846108be565b8281835e505f910152565b5f6108ef825190565b8084526020840193506109068185602086016108db565b601f01601f19169290920192915050565b5f6107d483836108e6565b5f61092b825190565b808452602084019350836020820285016109458560200190565b5f5b8481101561097857838303885281516109608482610917565b93505060208201602098909801979150600101610947565b50909695505050505050565b602080825281016107d48184610922565b80610817565b803561080881610995565b5f602082840312156109b9576109b95f5ffd5b6107d4838361099b565b602080825281016107d481846108e6565b5f602082840312156109e7576109e75f5ffd5b6107d48383610821565b5f6109fa825190565b610a088185602086016108db565b9290920192915050565b61080881836109f1565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610a4457607f821691505b602082108103610a5657610a56610a1c565b50919050565b5f610808610a678381565b90565b610a7383610a5c565b81545f1960089490940293841b1916921b91909117905550565b5f610a99818484610a6a565b505050565b81811015610ab857610ab05f82610a8d565b600101610a9e565b5050565b601f821115610a99575f818152602090206020601f85010481016020851015610ae25750805b6104916020601f860104830182610a9e565b815167ffffffffffffffff811115610b0e57610b0e61070d565b610b188254610a30565b610b23828285610abc565b506020601f821160018114610b56575f8315610b3f5750848201515b5f19600885021c1981166002850217855550610491565b5f84815260208120601f198516915b82811015610b855787850151825560209485019460019092019101610b65565b5084821015610ba157838701515f19601f87166008021c191681555b50505050600202600101905550565b60408082528101610bc181856108e6565b90506107d460208301846108be565b5f67ffffffffffffffff8216610808565b6108c781610bd0565b602081016108088284610be156fea2646970667358221220f457e8fde59a04e9f66ccf9d012844eaca4a678e02191968ed54708f7e046ff664736f6c634300081c0033",
}

// ManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ManagementContractMetaData.ABI instead.
var ManagementContractABI = ManagementContractMetaData.ABI

// ManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ManagementContractMetaData.Bin instead.
var ManagementContractBin = ManagementContractMetaData.Bin

// DeployManagementContract deploys a new Ethereum contract, binding an instance of ManagementContract to it.
func DeployManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ManagementContract, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// ManagementContract is an auto generated Go binding around an Ethereum contract.
type ManagementContract struct {
	ManagementContractCaller     // Read-only binding to the contract
	ManagementContractTransactor // Write-only binding to the contract
	ManagementContractFilterer   // Log filterer for contract events
}

// ManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ManagementContractSession struct {
	Contract     *ManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ManagementContractCallerSession struct {
	Contract *ManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ManagementContractTransactorSession struct {
	Contract     *ManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ManagementContractRaw struct {
	Contract *ManagementContract // Generic contract binding to access the raw methods on
}

// ManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ManagementContractCallerRaw struct {
	Contract *ManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// ManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ManagementContractTransactorRaw struct {
	Contract *ManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewManagementContract creates a new instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContract(address common.Address, backend bind.ContractBackend) (*ManagementContract, error) {
	contract, err := bindManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// NewManagementContractCaller creates a new read-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractCaller(address common.Address, caller bind.ContractCaller) (*ManagementContractCaller, error) {
	contract, err := bindManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractCaller{contract: contract}, nil
}

// NewManagementContractTransactor creates a new write-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ManagementContractTransactor, error) {
	contract, err := bindManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractTransactor{contract: contract}, nil
}

// NewManagementContractFilterer creates a new log filterer instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ManagementContractFilterer, error) {
	contract, err := bindManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ManagementContractFilterer{contract: contract}, nil
}

// bindManagementContract binds a generic wrapper to an already deployed contract.
func bindManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.ManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transact(opts, method, params...)
}

// GetImportantContractKeys is a free data retrieval call binding the contract method 0x6a30d26c.
//
// Solidity: function GetImportantContractKeys() view returns(string[])
func (_ManagementContract *ManagementContractCaller) GetImportantContractKeys(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetImportantContractKeys")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetImportantContractKeys is a free data retrieval call binding the contract method 0x6a30d26c.
//
// Solidity: function GetImportantContractKeys() view returns(string[])
func (_ManagementContract *ManagementContractSession) GetImportantContractKeys() ([]string, error) {
	return _ManagementContract.Contract.GetImportantContractKeys(&_ManagementContract.CallOpts)
}

// GetImportantContractKeys is a free data retrieval call binding the contract method 0x6a30d26c.
//
// Solidity: function GetImportantContractKeys() view returns(string[])
func (_ManagementContract *ManagementContractCallerSession) GetImportantContractKeys() ([]string, error) {
	return _ManagementContract.Contract.GetImportantContractKeys(&_ManagementContract.CallOpts)
}

// ImportantContractAddresses is a free data retrieval call binding the contract method 0x3e60a22f.
//
// Solidity: function importantContractAddresses(string ) view returns(address)
func (_ManagementContract *ManagementContractCaller) ImportantContractAddresses(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "importantContractAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ImportantContractAddresses is a free data retrieval call binding the contract method 0x3e60a22f.
//
// Solidity: function importantContractAddresses(string ) view returns(address)
func (_ManagementContract *ManagementContractSession) ImportantContractAddresses(arg0 string) (common.Address, error) {
	return _ManagementContract.Contract.ImportantContractAddresses(&_ManagementContract.CallOpts, arg0)
}

// ImportantContractAddresses is a free data retrieval call binding the contract method 0x3e60a22f.
//
// Solidity: function importantContractAddresses(string ) view returns(address)
func (_ManagementContract *ManagementContractCallerSession) ImportantContractAddresses(arg0 string) (common.Address, error) {
	return _ManagementContract.Contract.ImportantContractAddresses(&_ManagementContract.CallOpts, arg0)
}

// ImportantContractKeys is a free data retrieval call binding the contract method 0x98077e86.
//
// Solidity: function importantContractKeys(uint256 ) view returns(string)
func (_ManagementContract *ManagementContractCaller) ImportantContractKeys(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "importantContractKeys", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ImportantContractKeys is a free data retrieval call binding the contract method 0x98077e86.
//
// Solidity: function importantContractKeys(uint256 ) view returns(string)
func (_ManagementContract *ManagementContractSession) ImportantContractKeys(arg0 *big.Int) (string, error) {
	return _ManagementContract.Contract.ImportantContractKeys(&_ManagementContract.CallOpts, arg0)
}

// ImportantContractKeys is a free data retrieval call binding the contract method 0x98077e86.
//
// Solidity: function importantContractKeys(uint256 ) view returns(string)
func (_ManagementContract *ManagementContractCallerSession) ImportantContractKeys(arg0 *big.Int) (string, error) {
	return _ManagementContract.Contract.ImportantContractKeys(&_ManagementContract.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ManagementContract *ManagementContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ManagementContract *ManagementContractSession) Owner() (common.Address, error) {
	return _ManagementContract.Contract.Owner(&_ManagementContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ManagementContract *ManagementContractCallerSession) Owner() (common.Address, error) {
	return _ManagementContract.Contract.Owner(&_ManagementContract.CallOpts)
}

// SetImportantContractAddress is a paid mutator transaction binding the contract method 0x03e72e48.
//
// Solidity: function SetImportantContractAddress(string key, address newAddress) returns()
func (_ManagementContract *ManagementContractTransactor) SetImportantContractAddress(opts *bind.TransactOpts, key string, newAddress common.Address) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "SetImportantContractAddress", key, newAddress)
}

// SetImportantContractAddress is a paid mutator transaction binding the contract method 0x03e72e48.
//
// Solidity: function SetImportantContractAddress(string key, address newAddress) returns()
func (_ManagementContract *ManagementContractSession) SetImportantContractAddress(key string, newAddress common.Address) (*types.Transaction, error) {
	return _ManagementContract.Contract.SetImportantContractAddress(&_ManagementContract.TransactOpts, key, newAddress)
}

// SetImportantContractAddress is a paid mutator transaction binding the contract method 0x03e72e48.
//
// Solidity: function SetImportantContractAddress(string key, address newAddress) returns()
func (_ManagementContract *ManagementContractTransactorSession) SetImportantContractAddress(key string, newAddress common.Address) (*types.Transaction, error) {
	return _ManagementContract.Contract.SetImportantContractAddress(&_ManagementContract.TransactOpts, key, newAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ManagementContract *ManagementContractTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ManagementContract *ManagementContractSession) Initialize() (*types.Transaction, error) {
	return _ManagementContract.Contract.Initialize(&_ManagementContract.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ManagementContract *ManagementContractTransactorSession) Initialize() (*types.Transaction, error) {
	return _ManagementContract.Contract.Initialize(&_ManagementContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ManagementContract *ManagementContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ManagementContract *ManagementContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _ManagementContract.Contract.RenounceOwnership(&_ManagementContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ManagementContract *ManagementContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ManagementContract.Contract.RenounceOwnership(&_ManagementContract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ManagementContract *ManagementContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ManagementContract *ManagementContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ManagementContract.Contract.TransferOwnership(&_ManagementContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ManagementContract *ManagementContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ManagementContract.Contract.TransferOwnership(&_ManagementContract.TransactOpts, newOwner)
}

// ManagementContractImportantContractAddressUpdatedIterator is returned from FilterImportantContractAddressUpdated and is used to iterate over the raw logs and unpacked data for ImportantContractAddressUpdated events raised by the ManagementContract contract.
type ManagementContractImportantContractAddressUpdatedIterator struct {
	Event *ManagementContractImportantContractAddressUpdated // Event containing the contract specifics and raw log

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
func (it *ManagementContractImportantContractAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagementContractImportantContractAddressUpdated)
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
		it.Event = new(ManagementContractImportantContractAddressUpdated)
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
func (it *ManagementContractImportantContractAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagementContractImportantContractAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagementContractImportantContractAddressUpdated represents a ImportantContractAddressUpdated event raised by the ManagementContract contract.
type ManagementContractImportantContractAddressUpdated struct {
	Key        string
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterImportantContractAddressUpdated is a free log retrieval operation binding the contract event 0x17b2f9f5748931099ffee882b5b64f4a560b5c55da9b4f4e396dae3bb9f98cb5.
//
// Solidity: event ImportantContractAddressUpdated(string key, address newAddress)
func (_ManagementContract *ManagementContractFilterer) FilterImportantContractAddressUpdated(opts *bind.FilterOpts) (*ManagementContractImportantContractAddressUpdatedIterator, error) {

	logs, sub, err := _ManagementContract.contract.FilterLogs(opts, "ImportantContractAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &ManagementContractImportantContractAddressUpdatedIterator{contract: _ManagementContract.contract, event: "ImportantContractAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchImportantContractAddressUpdated is a free log subscription operation binding the contract event 0x17b2f9f5748931099ffee882b5b64f4a560b5c55da9b4f4e396dae3bb9f98cb5.
//
// Solidity: event ImportantContractAddressUpdated(string key, address newAddress)
func (_ManagementContract *ManagementContractFilterer) WatchImportantContractAddressUpdated(opts *bind.WatchOpts, sink chan<- *ManagementContractImportantContractAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _ManagementContract.contract.WatchLogs(opts, "ImportantContractAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagementContractImportantContractAddressUpdated)
				if err := _ManagementContract.contract.UnpackLog(event, "ImportantContractAddressUpdated", log); err != nil {
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

// ParseImportantContractAddressUpdated is a log parse operation binding the contract event 0x17b2f9f5748931099ffee882b5b64f4a560b5c55da9b4f4e396dae3bb9f98cb5.
//
// Solidity: event ImportantContractAddressUpdated(string key, address newAddress)
func (_ManagementContract *ManagementContractFilterer) ParseImportantContractAddressUpdated(log types.Log) (*ManagementContractImportantContractAddressUpdated, error) {
	event := new(ManagementContractImportantContractAddressUpdated)
	if err := _ManagementContract.contract.UnpackLog(event, "ImportantContractAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ManagementContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ManagementContract contract.
type ManagementContractInitializedIterator struct {
	Event *ManagementContractInitialized // Event containing the contract specifics and raw log

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
func (it *ManagementContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagementContractInitialized)
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
		it.Event = new(ManagementContractInitialized)
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
func (it *ManagementContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagementContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagementContractInitialized represents a Initialized event raised by the ManagementContract contract.
type ManagementContractInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ManagementContract *ManagementContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*ManagementContractInitializedIterator, error) {

	logs, sub, err := _ManagementContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ManagementContractInitializedIterator{contract: _ManagementContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ManagementContract *ManagementContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ManagementContractInitialized) (event.Subscription, error) {

	logs, sub, err := _ManagementContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagementContractInitialized)
				if err := _ManagementContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ManagementContract *ManagementContractFilterer) ParseInitialized(log types.Log) (*ManagementContractInitialized, error) {
	event := new(ManagementContractInitialized)
	if err := _ManagementContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ManagementContractLogManagementContractCreatedIterator is returned from FilterLogManagementContractCreated and is used to iterate over the raw logs and unpacked data for LogManagementContractCreated events raised by the ManagementContract contract.
type ManagementContractLogManagementContractCreatedIterator struct {
	Event *ManagementContractLogManagementContractCreated // Event containing the contract specifics and raw log

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
func (it *ManagementContractLogManagementContractCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagementContractLogManagementContractCreated)
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
		it.Event = new(ManagementContractLogManagementContractCreated)
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
func (it *ManagementContractLogManagementContractCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagementContractLogManagementContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagementContractLogManagementContractCreated represents a LogManagementContractCreated event raised by the ManagementContract contract.
type ManagementContractLogManagementContractCreated struct {
	MessageBusAddress common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogManagementContractCreated is a free log retrieval operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) FilterLogManagementContractCreated(opts *bind.FilterOpts) (*ManagementContractLogManagementContractCreatedIterator, error) {

	logs, sub, err := _ManagementContract.contract.FilterLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return &ManagementContractLogManagementContractCreatedIterator{contract: _ManagementContract.contract, event: "LogManagementContractCreated", logs: logs, sub: sub}, nil
}

// WatchLogManagementContractCreated is a free log subscription operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) WatchLogManagementContractCreated(opts *bind.WatchOpts, sink chan<- *ManagementContractLogManagementContractCreated) (event.Subscription, error) {

	logs, sub, err := _ManagementContract.contract.WatchLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagementContractLogManagementContractCreated)
				if err := _ManagementContract.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
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

// ParseLogManagementContractCreated is a log parse operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) ParseLogManagementContractCreated(log types.Log) (*ManagementContractLogManagementContractCreated, error) {
	event := new(ManagementContractLogManagementContractCreated)
	if err := _ManagementContract.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ManagementContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ManagementContract contract.
type ManagementContractOwnershipTransferredIterator struct {
	Event *ManagementContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ManagementContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagementContractOwnershipTransferred)
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
		it.Event = new(ManagementContractOwnershipTransferred)
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
func (it *ManagementContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagementContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagementContractOwnershipTransferred represents a OwnershipTransferred event raised by the ManagementContract contract.
type ManagementContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ManagementContract *ManagementContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ManagementContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ManagementContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ManagementContractOwnershipTransferredIterator{contract: _ManagementContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ManagementContract *ManagementContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ManagementContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ManagementContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagementContractOwnershipTransferred)
				if err := _ManagementContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ManagementContract *ManagementContractFilterer) ParseOwnershipTransferred(log types.Log) (*ManagementContractOwnershipTransferred, error) {
	event := new(ManagementContractOwnershipTransferred)
	if err := _ManagementContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
