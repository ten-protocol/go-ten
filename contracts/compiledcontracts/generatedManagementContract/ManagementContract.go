// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

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

// GeneratedManagementContractMetaData contains all meta data concerning the GeneratedManagementContract contract.
var GeneratedManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes20\",\"name\":\"AggregatorID\",\"type\":\"bytes20\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"requesterID\",\"type\":\"bytes20\"},{\"internalType\":\"string\",\"name\":\"pubKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"inputSecret\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes20\",\"name\":\"AggregatorID\",\"type\":\"bytes20\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610c98806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80638ef74f8914610067578063b040f26b14610097578063bdbf47e0146100c7578063e0643dfc146100e3578063e34fbfc814610116578063fbd66cbb14610132575b600080fd5b610081600480360381019061007c919061063a565b61014e565b60405161008e9190610700565b60405180910390f35b6100b160048036038101906100ac919061077a565b6101ee565b6040516100be9190610700565b60405180910390f35b6100e160048036038101906100dc919061093c565b61028e565b005b6100fd60048036038101906100f89190610a32565b6102db565b60405161010d9493929190610aa9565b60405180910390f35b610130600480360381019061012b9190610aee565b610335565b005b61014c60048036038101906101479190610b67565b610388565b005b6001602052806000526040600020600091509050805461016d90610c30565b80601f016020809104026020016040519081016040528092919081815260200182805461019990610c30565b80156101e65780601f106101bb576101008083540402835291602001916101e6565b820191906000526020600020905b8154815290600101906020018083116101c957829003601f168201915b505050505081565b6002602052806000526040600020600091509050805461020d90610c30565b80601f016020809104026020016040519081016040528092919081815260200182805461023990610c30565b80156102865780601f1061025b57610100808354040283529160200191610286565b820191906000526020600020905b81548152906001019060200180831161026957829003601f168201915b505050505081565b8360026000876bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002090805190602001906102d392919061049f565b505050505050565b600060205281600052604060002081815481106102f757600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900460601b908060020154908060030154905084565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209190610383929190610525565b505050565b600060026000876bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200190815260200160002080546103c690610c30565b9050116103d257600080fd5b60006040518060800160405280888152602001876bffffffffffffffffffffffff191681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908360601c02179055506040820151816002015560608201518160030155505050505050505050565b8280546104ab90610c30565b90600052602060002090601f0160209004810192826104cd5760008555610514565b82601f106104e657805160ff1916838001178555610514565b82800160010185558215610514579182015b828111156105135782518255916020019190600101906104f8565b5b50905061052191906105ab565b5090565b82805461053190610c30565b90600052602060002090601f016020900481019282610553576000855561059a565b82601f1061056c57803560ff191683800117855561059a565b8280016001018555821561059a579182015b8281111561059957823582559160200191906001019061057e565b5b5090506105a791906105ab565b5090565b5b808211156105c45760008160009055506001016105ac565b5090565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610607826105dc565b9050919050565b610617816105fc565b811461062257600080fd5b50565b6000813590506106348161060e565b92915050565b6000602082840312156106505761064f6105d2565b5b600061065e84828501610625565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156106a1578082015181840152602081019050610686565b838111156106b0576000848401525b50505050565b6000601f19601f8301169050919050565b60006106d282610667565b6106dc8185610672565b93506106ec818560208601610683565b6106f5816106b6565b840191505092915050565b6000602082019050818103600083015261071a81846106c7565b905092915050565b60007fffffffffffffffffffffffffffffffffffffffff00000000000000000000000082169050919050565b61075781610722565b811461076257600080fd5b50565b6000813590506107748161074e565b92915050565b6000602082840312156107905761078f6105d2565b5b600061079e84828501610765565b91505092915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6107e9826106b6565b810181811067ffffffffffffffff82111715610808576108076107b1565b5b80604052505050565b600061081b6105c8565b905061082782826107e0565b919050565b600067ffffffffffffffff821115610847576108466107b1565b5b610850826106b6565b9050602081019050919050565b82818337600083830152505050565b600061087f61087a8461082c565b610811565b90508281526020810184848401111561089b5761089a6107ac565b5b6108a684828561085d565b509392505050565b600082601f8301126108c3576108c26107a7565b5b81356108d384826020860161086c565b91505092915050565b600080fd5b600080fd5b60008083601f8401126108fc576108fb6107a7565b5b8235905067ffffffffffffffff811115610919576109186108dc565b5b602083019150836001820283011115610935576109346108e1565b5b9250929050565b600080600080600060808688031215610958576109576105d2565b5b600061096688828901610765565b955050602086013567ffffffffffffffff811115610987576109866105d7565b5b610993888289016108ae565b945050604086013567ffffffffffffffff8111156109b4576109b36105d7565b5b6109c0888289016108ae565b935050606086013567ffffffffffffffff8111156109e1576109e06105d7565b5b6109ed888289016108e6565b92509250509295509295909350565b6000819050919050565b610a0f816109fc565b8114610a1a57600080fd5b50565b600081359050610a2c81610a06565b92915050565b60008060408385031215610a4957610a486105d2565b5b6000610a5785828601610a1d565b9250506020610a6885828601610a1d565b9150509250929050565b6000819050919050565b610a8581610a72565b82525050565b610a9481610722565b82525050565b610aa3816109fc565b82525050565b6000608082019050610abe6000830187610a7c565b610acb6020830186610a8b565b610ad86040830185610a7c565b610ae56060830184610a9a565b95945050505050565b60008060208385031215610b0557610b046105d2565b5b600083013567ffffffffffffffff811115610b2357610b226105d7565b5b610b2f858286016108e6565b92509250509250929050565b610b4481610a72565b8114610b4f57600080fd5b50565b600081359050610b6181610b3b565b92915050565b60008060008060008060a08789031215610b8457610b836105d2565b5b6000610b9289828a01610b52565b9650506020610ba389828a01610765565b9550506040610bb489828a01610b52565b9450506060610bc589828a01610a1d565b935050608087013567ffffffffffffffff811115610be657610be56105d7565b5b610bf289828a016108e6565b92509250509295509295509295565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610c4857607f821691505b60208210811415610c5c57610c5b610c01565b5b5091905056fea2646970667358221220afffa69a8add163598bd81b9d077577d2b0dfe619aa10bd7d8808516951753cc64736f6c634300080c0033",
}

// GeneratedManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GeneratedManagementContractMetaData.ABI instead.
var GeneratedManagementContractABI = GeneratedManagementContractMetaData.ABI

// GeneratedManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GeneratedManagementContractMetaData.Bin instead.
var GeneratedManagementContractBin = GeneratedManagementContractMetaData.Bin

// DeployGeneratedManagementContract deploys a new Ethereum contract, binding an instance of GeneratedManagementContract to it.
func DeployGeneratedManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GeneratedManagementContract, error) {
	parsed, err := GeneratedManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GeneratedManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// GeneratedManagementContract is an auto generated Go binding around an Ethereum contract.
type GeneratedManagementContract struct {
	GeneratedManagementContractCaller     // Read-only binding to the contract
	GeneratedManagementContractTransactor // Write-only binding to the contract
	GeneratedManagementContractFilterer   // Log filterer for contract events
}

// GeneratedManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedManagementContractSession struct {
	Contract     *GeneratedManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GeneratedManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedManagementContractCallerSession struct {
	Contract *GeneratedManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// GeneratedManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedManagementContractTransactorSession struct {
	Contract     *GeneratedManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// GeneratedManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedManagementContractRaw struct {
	Contract *GeneratedManagementContract // Generic contract binding to access the raw methods on
}

// GeneratedManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCallerRaw struct {
	Contract *GeneratedManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactorRaw struct {
	Contract *GeneratedManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeneratedManagementContract creates a new instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContract(address common.Address, backend bind.ContractBackend) (*GeneratedManagementContract, error) {
	contract, err := bindGeneratedManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// NewGeneratedManagementContractCaller creates a new read-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractCaller(address common.Address, caller bind.ContractCaller) (*GeneratedManagementContractCaller, error) {
	contract, err := bindGeneratedManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractCaller{contract: contract}, nil
}

// NewGeneratedManagementContractTransactor creates a new write-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedManagementContractTransactor, error) {
	contract, err := bindGeneratedManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractTransactor{contract: contract}, nil
}

// NewGeneratedManagementContractFilterer creates a new log filterer instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedManagementContractFilterer, error) {
	contract, err := bindGeneratedManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractFilterer{contract: contract}, nil
}

// bindGeneratedManagementContract binds a generic wrapper to an already deployed contract.
func bindGeneratedManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transact(opts, method, params...)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) AttestationRequests(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attestationRequests", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xb040f26b.
//
// Solidity: function attested(bytes20 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Attested(opts *bind.CallOpts, arg0 [20]byte) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attested", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0xb040f26b.
//
// Solidity: function attested(bytes20 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Attested(arg0 [20]byte) (string, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xb040f26b.
//
// Solidity: function attested(bytes20 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Attested(arg0 [20]byte) (string, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID [20]byte
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		AggregatorID [20]byte
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[1], new([20]byte)).(*[20]byte)
	outstruct.L1Block = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID [20]byte
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID [20]byte
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// AddRollup is a paid mutator transaction binding the contract method 0xfbd66cbb.
//
// Solidity: function AddRollup(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID [20]byte, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xfbd66cbb.
//
// Solidity: function AddRollup(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID [20]byte, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xfbd66cbb.
//
// Solidity: function AddRollup(bytes32 ParentHash, bytes20 AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID [20]byte, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbdbf47e0.
//
// Solidity: function RespondNetworkSecret(bytes20 requesterID, string pubKey, string inputSecret, string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, requesterID [20]byte, pubKey string, inputSecret string, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", requesterID, pubKey, inputSecret, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbdbf47e0.
//
// Solidity: function RespondNetworkSecret(bytes20 requesterID, string pubKey, string inputSecret, string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(requesterID [20]byte, pubKey string, inputSecret string, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, requesterID, pubKey, inputSecret, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbdbf47e0.
//
// Solidity: function RespondNetworkSecret(bytes20 requesterID, string pubKey, string inputSecret, string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(requesterID [20]byte, pubKey string, inputSecret string, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, requesterID, pubKey, inputSecret, requestReport)
}
