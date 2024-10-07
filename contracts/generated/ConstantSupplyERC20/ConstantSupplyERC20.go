// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConstantSupplyERC20

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

// ConstantSupplyERC20MetaData contains all meta data concerning the ConstantSupplyERC20 contract.
var ConstantSupplyERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000ebb38038062000ebb833981016040819052620000349162000312565b8282600362000044838262000497565b50600462000053828262000497565b5050506200006833826200007160201b60201c565b50505062000607565b6001600160a01b038216620000a757600060405163ec442f0560e01b81526004016200009e919062000586565b60405180910390fd5b620000b560008383620000b9565b5050565b6001600160a01b038316620000e8578060026000828254620000dc9190620005ac565b909155506200014a9050565b6001600160a01b038316600090815260208190526040902054818110156200012b5783818360405163391434e360e21b81526004016200009e93929190620005c9565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b038216620001685760028054829003905562000187565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001cc9190620005f7565b60405180910390a3505050565b634e487b7160e01b600052604160045260246000fd5b601f19601f83011681016001600160401b0381118282101715620002175762000217620001d9565b6040525050565b60006200022a60405190565b9050620002388282620001ef565b919050565b60006001600160401b03821115620002595762000259620001d9565b601f19601f83011660200192915050565b60005b83811015620002875781810151838201526020016200026d565b50506000910152565b6000620002a7620002a1846200023d565b6200021e565b905082815260208101848484011115620002c457620002c4600080fd5b620002d18482856200026a565b509392505050565b600082601f830112620002ef57620002ef600080fd5b81516200030184826020860162000290565b949350505050565b80515b92915050565b6000806000606084860312156200032c576200032c600080fd5b83516001600160401b03811115620003475762000347600080fd5b6200035586828701620002d9565b602086015190945090506001600160401b03811115620003785762000378600080fd5b6200038686828701620002d9565b9250506040620003998682870162000309565b9150509250925092565b634e487b7160e01b600052602260045260246000fd5b600281046001821680620003ce57607f821691505b602082108103620003e357620003e3620003a3565b50919050565b60006200030c620003f78381565b90565b6200040583620003e9565b815460001960089490940293841b1916921b91909117905550565b60006200042f818484620003fa565b505050565b81811015620000b5576200044a60008262000420565b60010162000434565b601f8211156200042f576000818152602090206020601f850104810160208510156200047c5750805b620004906020601f86010483018262000434565b5050505050565b81516001600160401b03811115620004b357620004b3620001d9565b620004bf8254620003b9565b620004cc82828562000453565b506020601f821160018114620005045760008315620004eb5750848201515b600019600885021c198116600285021785555062000490565b600084815260208120601f198516915b8281101562000536578785015182556020948501946001909201910162000514565b5084821015620005545783870151600019601f87166008021c191681555b50505050600202600101905550565b60006001600160a01b0382166200030c565b620005808162000563565b82525050565b602081016200030c828462000575565b634e487b7160e01b600052601160045260246000fd5b808201808211156200030c576200030c62000596565b8062000580565b60608101620005d9828662000575565b620005e86020830185620005c2565b620003016040830184620005c2565b602081016200030c8284620005c2565b6108a480620006176000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063313ce5671161007657806395d89b411161005b57806395d89b4114610142578063a9059cbb1461014a578063dd62ed3e1461015d57600080fd5b8063313ce5671461010a57806370a082311461011957600080fd5b806306fdde03146100a8578063095ea7b3146100c657806318160ddd146100e657806323b872dd146100f7575b600080fd5b6100b0610196565b6040516100bd919061063c565b60405180910390f35b6100d96100d4366004610698565b610228565b6040516100bd91906106df565b6002545b6040516100bd91906106f3565b6100d9610105366004610701565b610242565b60126040516100bd919061075a565b6100ea610127366004610768565b6001600160a01b031660009081526020819052604090205490565b6100b0610266565b6100d9610158366004610698565b610275565b6100ea61016b366004610791565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6060600380546101a5906107da565b80601f01602080910402602001604051908101604052809291908181526020018280546101d1906107da565b801561021e5780601f106101f35761010080835404028352916020019161021e565b820191906000526020600020905b81548152906001019060200180831161020157829003601f168201915b5050505050905090565b600033610236818585610283565b60019150505b92915050565b600033610250858285610295565b61025b858585610322565b506001949350505050565b6060600480546101a5906107da565b600033610236818585610322565b61029083838360016103b3565b505050565b6001600160a01b03838116600090815260016020908152604080832093861683529290522054600019811461031c578181101561030d578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016103049392919061080f565b60405180910390fd5b61031c848484840360006103b3565b50505050565b6001600160a01b0383166103655760006040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016103049190610837565b6001600160a01b0382166103a85760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016103049190610837565b6102908383836104b8565b6001600160a01b0384166103f65760006040517fe602df050000000000000000000000000000000000000000000000000000000081526004016103049190610837565b6001600160a01b0383166104395760006040517f94280d620000000000000000000000000000000000000000000000000000000081526004016103049190610837565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561031c57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516104aa91906106f3565b60405180910390a350505050565b6001600160a01b0383166104e35780600260008282546104d8919061085b565b9091555061055b9050565b6001600160a01b0383166000908152602081905260409020548181101561053c578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016103049392919061080f565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b03821661057757600280548290039055610596565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105d991906106f3565b60405180910390a3505050565b60005b838110156106015781810151838201526020016105e9565b50506000910152565b6000610614825190565b80845260208401935061062b8185602086016105e6565b601f01601f19169290920192915050565b6020808252810161064d818461060a565b9392505050565b60006001600160a01b03821661023c565b61066e81610654565b811461067957600080fd5b50565b803561023c81610665565b8061066e565b803561023c81610687565b600080604083850312156106ae576106ae600080fd5b60006106ba858561067c565b92505060206106cb8582860161068d565b9150509250929050565b8015155b82525050565b6020810161023c82846106d5565b806106d9565b6020810161023c82846106ed565b60008060006060848603121561071957610719600080fd5b6000610725868661067c565b93505060206107368682870161067c565b92505060406107478682870161068d565b9150509250925092565b60ff81166106d9565b6020810161023c8284610751565b60006020828403121561077d5761077d600080fd5b6000610789848461067c565b949350505050565b600080604083850312156107a7576107a7600080fd5b60006107b3858561067c565b92505060206106cb8582860161067c565b634e487b7160e01b600052602260045260246000fd5b6002810460018216806107ee57607f821691505b602082108103610800576108006107c4565b50919050565b6106d981610654565b6060810161081d8286610806565b61082a60208301856106ed565b61078960408301846106ed565b6020810161023c8284610806565b634e487b7160e01b600052601160045260246000fd5b8082018082111561023c5761023c61084556fea26469706673582212202d779dff7fe9ec7ab0604366da8ad9497916e439b37454dbdbc1f22c4e6aa1fb64736f6c63430008140033",
}

// ConstantSupplyERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ConstantSupplyERC20MetaData.ABI instead.
var ConstantSupplyERC20ABI = ConstantSupplyERC20MetaData.ABI

// ConstantSupplyERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConstantSupplyERC20MetaData.Bin instead.
var ConstantSupplyERC20Bin = ConstantSupplyERC20MetaData.Bin

// DeployConstantSupplyERC20 deploys a new Ethereum contract, binding an instance of ConstantSupplyERC20 to it.
func DeployConstantSupplyERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int) (common.Address, *types.Transaction, *ConstantSupplyERC20, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConstantSupplyERC20Bin), backend, name, symbol, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// ConstantSupplyERC20 is an auto generated Go binding around an Ethereum contract.
type ConstantSupplyERC20 struct {
	ConstantSupplyERC20Caller     // Read-only binding to the contract
	ConstantSupplyERC20Transactor // Write-only binding to the contract
	ConstantSupplyERC20Filterer   // Log filterer for contract events
}

// ConstantSupplyERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstantSupplyERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstantSupplyERC20Session struct {
	Contract     *ConstantSupplyERC20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstantSupplyERC20CallerSession struct {
	Contract *ConstantSupplyERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ConstantSupplyERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstantSupplyERC20TransactorSession struct {
	Contract     *ConstantSupplyERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ConstantSupplyERC20Raw struct {
	Contract *ConstantSupplyERC20 // Generic contract binding to access the raw methods on
}

// ConstantSupplyERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20CallerRaw struct {
	Contract *ConstantSupplyERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ConstantSupplyERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20TransactorRaw struct {
	Contract *ConstantSupplyERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewConstantSupplyERC20 creates a new instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20(address common.Address, backend bind.ContractBackend) (*ConstantSupplyERC20, error) {
	contract, err := bindConstantSupplyERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// NewConstantSupplyERC20Caller creates a new read-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Caller(address common.Address, caller bind.ContractCaller) (*ConstantSupplyERC20Caller, error) {
	contract, err := bindConstantSupplyERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Caller{contract: contract}, nil
}

// NewConstantSupplyERC20Transactor creates a new write-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ConstantSupplyERC20Transactor, error) {
	contract, err := bindConstantSupplyERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Transactor{contract: contract}, nil
}

// NewConstantSupplyERC20Filterer creates a new log filterer instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ConstantSupplyERC20Filterer, error) {
	contract, err := bindConstantSupplyERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Filterer{contract: contract}, nil
}

// bindConstantSupplyERC20 binds a generic wrapper to an already deployed contract.
func bindConstantSupplyERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// ConstantSupplyERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20ApprovalIterator struct {
	Event *ConstantSupplyERC20Approval // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Approval)
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
		it.Event = new(ConstantSupplyERC20Approval)
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
func (it *ConstantSupplyERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Approval represents a Approval event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ConstantSupplyERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20ApprovalIterator{contract: _ConstantSupplyERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Approval)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseApproval(log types.Log) (*ConstantSupplyERC20Approval, error) {
	event := new(ConstantSupplyERC20Approval)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConstantSupplyERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20TransferIterator struct {
	Event *ConstantSupplyERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Transfer)
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
		it.Event = new(ConstantSupplyERC20Transfer)
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
func (it *ConstantSupplyERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Transfer represents a Transfer event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConstantSupplyERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20TransferIterator{contract: _ConstantSupplyERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Transfer)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseTransfer(log types.Log) (*ConstantSupplyERC20Transfer, error) {
	event := new(ConstantSupplyERC20Transfer)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
